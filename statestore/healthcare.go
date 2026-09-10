package statestore

import (
	"context"
)

const (
	hcSlotBooked       = "healthcare.slot.booked"
	hcBookingCancelled = "healthcare.booking.cancelled"
	hcPreseedSession   = "__preseed__"
)

var hcPreseedSlots = []struct{ id, doctor, day, time string }{
	{"patel-wed-1430", "Dr. Priya Patel", "Wednesday", "14:30"},
	{"patel-fri-1430", "Dr. Priya Patel", "Friday", "14:30"},
	{"clarke-mon-0900", "Dr. James Clarke", "Monday", "09:00"},
}

func (s *Store) migrateHealthcareDemo() error {
	var count int
	if err := s.db.QueryRow(
		`SELECT COUNT(*) FROM events WHERE session_id = ? AND type = ?`,
		hcPreseedSession, hcSlotBooked,
	).Scan(&count); err != nil {
		return err
	}
	if count >= len(hcPreseedSlots) {
		return nil
	}
	for _, slot := range hcPreseedSlots {
		if _, err := s.AppendEvent(context.Background(), hcPreseedSession, hcSlotBooked, map[string]any{
			"slot_id":     slot.id,
			"doctor_name": slot.doctor,
			"day":         slot.day,
			"time":        slot.time,
		}); err != nil {
			return err
		}
	}
	return nil
}

// HealthcareBookedSlots returns map[slotID]sessionID for currently-booked slots.
// A slot is booked if its most-recent event across all sessions is a booking, not a cancellation.
func (s *Store) HealthcareBookedSlots(ctx context.Context) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT slot_id, session_id FROM (
			SELECT
				json_extract(payload, '$.slot_id') AS slot_id,
				session_id,
				type,
				ROW_NUMBER() OVER (PARTITION BY json_extract(payload, '$.slot_id') ORDER BY id DESC) AS rn
			FROM events
			WHERE type IN (?, ?)
		) WHERE rn = 1 AND type = ?
	`, hcSlotBooked, hcBookingCancelled, hcSlotBooked)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var slotID, sessionID string
		if err := rows.Scan(&slotID, &sessionID); err != nil {
			return nil, err
		}
		result[slotID] = sessionID
	}
	return result, rows.Err()
}
