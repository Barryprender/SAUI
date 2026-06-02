package healthcare

import (
	"context"
	"errors"

	"saui/statestore"
)

// BookSlot reserves an appointment slot for the session.
type BookSlot struct {
	SlotID string
}

func (a BookSlot) Type() string { return SlotBooked }

func (a BookSlot) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if _, ok := SlotByID(a.SlotID); !ok {
		return errors.New("unknown slot")
	}
	booked, err := s.HealthcareBookedSlots(ctx)
	if err != nil {
		return err
	}
	if _, taken := booked[a.SlotID]; taken {
		return errors.New("slot already taken")
	}
	return nil
}

func (a BookSlot) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	slot, _ := SlotByID(a.SlotID)
	return s.AppendEvent(ctx, sessionID, SlotBooked, map[string]any{
		"slot_id":     slot.ID,
		"doctor_name": slot.Doctor.Name,
		"day":         slot.Day,
		"time":        slot.Time,
	})
}

// CancelBooking releases a slot the session currently holds.
type CancelBooking struct {
	SlotID string
}

func (a CancelBooking) Type() string { return BookingCancelled }

func (a CancelBooking) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	booked, err := s.HealthcareBookedSlots(ctx)
	if err != nil {
		return err
	}
	holder, taken := booked[a.SlotID]
	if !taken || holder != sessionID {
		return errors.New("no active booking to cancel")
	}
	return nil
}

func (a CancelBooking) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, BookingCancelled, map[string]any{
		"slot_id": a.SlotID,
	})
}
