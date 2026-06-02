package healthcare

func slotStatusClass(s SlotStatus) string {
	switch s {
	case SlotBookedByMe:
		return "mine"
	case SlotBookedByOther:
		return "taken"
	default:
		return "available"
	}
}

func justBookedClass(slotID, justBookedSlotID string) string {
	if justBookedSlotID != "" && slotID == justBookedSlotID {
		return " booking-event-id--new"
	}
	return ""
}
