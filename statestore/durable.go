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
