package healthcare

import (
	"context"
	"encoding/json"

	"saui/statestore"
)

const ScheduleProjectionName = "healthcare.schedule"

// SlotStatus is the booking state of a slot relative to the current session.
type SlotStatus int

const (
	SlotAvailable    SlotStatus = iota
	SlotBookedByMe
	SlotBookedByOther
)

// ScheduleSlot is a rendered slot in the schedule grid.
type ScheduleSlot struct {
	ID          string
	DoctorName  string
	Specialty   string
	AvatarClass string
	Day         string
	Time        string
	Status      SlotStatus
	EventID     int64 // set when Status == SlotBookedByMe
}

// DayGroup groups slots by day for rendering.
type DayGroup struct {
	Day   string
	Slots []ScheduleSlot
}

// BookingRecord is a confirmed appointment in the booking panel.
type BookingRecord struct {
	SlotID     string
	DoctorName string
	Specialty  string
	Day        string
	Time       string
	EventID    int64
}

// ScheduleProjection is the read model for the healthcare schedule page.
type ScheduleProjection struct {
	Days       []DayGroup
	MyBookings []BookingRecord
}

func (p ScheduleProjection) Name() string { return ScheduleProjectionName }

// FindSlot looks up a slot by ID across all day groups.
func (p ScheduleProjection) FindSlot(id string) (ScheduleSlot, bool) {
	for _, dg := range p.Days {
		for _, s := range dg.Slots {
			if s.ID == id {
				return s, true
			}
		}
	}
	return ScheduleSlot{}, false
}

func init() {
	statestore.RegisterProjection(ScheduleProjectionName, buildSchedule)
}

func buildSchedule(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
	globalBooked, err := s.HealthcareBookedSlots(ctx)
	if err != nil {
		return nil, err
	}

	myEvents, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Most-recent booking event ID per slot for the current session.
	myBookingEventID := make(map[string]int64)
	for _, ev := range myEvents {
		if ev.Type == SlotBooked {
			var p struct {
				SlotID string `json:"slot_id"`
			}
			if json.Unmarshal(ev.Payload, &p) == nil {
				myBookingEventID[p.SlotID] = ev.ID
			}
		}
	}

	// Build day groups in slot definition order (days then times then doctors).
	dayMap := make(map[string]*DayGroup)
	var dayOrder []string

	for _, slot := range AllSlots() {
		if _, exists := dayMap[slot.Day]; !exists {
			dayMap[slot.Day] = &DayGroup{Day: slot.Day}
			dayOrder = append(dayOrder, slot.Day)
		}

		ss := ScheduleSlot{
			ID:          slot.ID,
			DoctorName:  slot.Doctor.Name,
			Specialty:   slot.Doctor.Specialty,
			AvatarClass: slot.Doctor.AvatarClass,
			Day:         slot.Day,
			Time:        slot.Time,
			Status:      SlotAvailable,
		}

		if holderSession, taken := globalBooked[slot.ID]; taken {
			if holderSession == sessionID {
				ss.Status = SlotBookedByMe
				ss.EventID = myBookingEventID[slot.ID]
			} else {
				ss.Status = SlotBookedByOther
			}
		}

		dayMap[slot.Day].Slots = append(dayMap[slot.Day].Slots, ss)
	}

	proj := ScheduleProjection{}
	for _, day := range dayOrder {
		proj.Days = append(proj.Days, *dayMap[day])
	}

	for _, dg := range proj.Days {
		for _, ss := range dg.Slots {
			if ss.Status == SlotBookedByMe {
				proj.MyBookings = append(proj.MyBookings, BookingRecord{
					SlotID:     ss.ID,
					DoctorName: ss.DoctorName,
					Specialty:  ss.Specialty,
					Day:        ss.Day,
					Time:       ss.Time,
					EventID:    ss.EventID,
				})
			}
		}
	}

	return proj, nil
}
