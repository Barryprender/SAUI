package statestore

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Event is an immutable record of something that happened.
type Event struct {
	ID         int64
	SessionID  string
	Type       string
	Payload    json.RawMessage
	OccurredAt time.Time
}

func (s *Store) migrate() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS events (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id  TEXT    NOT NULL,
		type        TEXT    NOT NULL,
		payload     TEXT    NOT NULL DEFAULT '{}',
		occurred_at INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_events_session ON events(session_id);
	CREATE INDEX IF NOT EXISTS idx_events_type    ON events(type);
	`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) AppendEvent(ctx context.Context, sessionID, eventType string, payload any) (Event, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return Event{}, err
	}
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO events (session_id, type, payload, occurred_at) VALUES (?, ?, ?, ?)`,
		sessionID, eventType, string(p), now.UnixMilli(),
	)
	if err != nil {
		return Event{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Event{}, fmt.Errorf("last insert id: %w", err)
	}
	return Event{ID: id, SessionID: sessionID, Type: eventType, Payload: p, OccurredAt: now}, nil
}

// eventRetention is how long an event stays attributable to the session that
// produced it. The only personal data in the log is session_id — an online
// identifier under GDPR Recital 30 — and the session cookie expires after 24
// hours, so beyond this window nobody can reach their own rows anyway.
// Article 5(1)(e) does not permit keeping identifiers with no purpose left.
const eventRetention = 30 * 24 * time.Hour

// PurgeExpired enforces that window. Events of the retained types are detached
// from their session rather than removed, because the retained type is visitor
// feedback and the argument it carries outlives the identifier attached to it;
// everything else past the window is deleted outright.
//
// Callers pass the retained types, so this package stays ignorant of the
// domain's event names. Only the placeholder count is interpolated into SQL;
// every value is bound.
// The demo fixtures share this table. The healthcare demo preseeds booked
// slots as events under hcPreseedSession, which are scenery rather than a
// visitor's activity: they carry no personal data and the demo is wrong
// without them, so retention leaves them alone.
func (s *Store) PurgeExpired(ctx context.Context, now time.Time, retain ...string) (anonymised, deleted int64, err error) {
	args := []any{now.Add(-eventRetention).UnixMilli(), hcPreseedSession}
	var list string
	for i, t := range retain {
		if i > 0 {
			list += ", "
		}
		list += "?"
		args = append(args, t)
	}
	keep, drop := "", ""
	if len(retain) > 0 {
		keep = " AND type IN (" + list + ")"
		drop = " AND type NOT IN (" + list + ")"
	}

	// One transaction, so a failure between the two statements cannot leave
	// feedback detached while the events around it survive.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("begin retention purge: %w", err)
	}
	defer tx.Rollback()

	// session_id NOT IN ('', fixtures) keeps the pass idempotent and keeps the
	// fixtures out of it.
	const expired = ` WHERE occurred_at < ? AND session_id NOT IN ('', ?)`

	if len(retain) > 0 {
		res, err := tx.ExecContext(ctx, `UPDATE events SET session_id = ''`+expired+keep, args...)
		if err != nil {
			return 0, 0, fmt.Errorf("detach expired events: %w", err)
		}
		anonymised, _ = res.RowsAffected()
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM events`+expired+drop, args...)
	if err != nil {
		return 0, 0, fmt.Errorf("delete expired events: %w", err)
	}
	deleted, _ = res.RowsAffected()

	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("commit retention purge: %w", err)
	}
	return anonymised, deleted, nil
}

func (s *Store) EventsBySession(ctx context.Context, sessionID string) ([]Event, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, session_id, type, payload, occurred_at FROM events WHERE session_id = ? ORDER BY id ASC`,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		var ms int64
		var payload string
		if err := rows.Scan(&e.ID, &e.SessionID, &e.Type, &payload, &ms); err != nil {
			return nil, err
		}
		e.Payload = json.RawMessage(payload)
		e.OccurredAt = time.UnixMilli(ms).UTC()
		events = append(events, e)
	}
	return events, rows.Err()
}

// CountEventsBySessionAndType returns the count of events for a session matching the given type.
func (s *Store) CountEventsBySessionAndType(ctx context.Context, sessionID, eventType string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM events WHERE session_id = ? AND type = ?`,
		sessionID, eventType,
	).Scan(&count)
	return count, err
}
