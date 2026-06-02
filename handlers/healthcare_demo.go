package handlers

import (
	"net/http"

	"saui/demo/healthcare"
	"saui/middleware"
)

func (h *Handler) HealthcareSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	proj, err := h.gateway.Project(ctx, sessionID, healthcare.ScheduleProjectionName)
	if err != nil {
		h.logger.Error("project healthcare schedule", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	schedule := proj.(healthcare.ScheduleProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := healthcare.Schedule(schedule, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render healthcare schedule", "err", err)
	}
}

func (h *Handler) HealthcareBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)
	slotID := r.FormValue("slot_id")

	dispatchErr := h.gateway.Dispatch(ctx, sessionID, healthcare.BookSlot{SlotID: slotID})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/healthcare", http.StatusSeeOther)
		return
	}

	proj, err := h.gateway.Project(ctx, sessionID, healthcare.ScheduleProjectionName)
	if err != nil {
		h.logger.Error("project healthcare schedule", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	schedule := proj.(healthcare.ScheduleProjection)

	justBooked := ""
	if dispatchErr == nil {
		justBooked = slotID
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := healthcare.BookingPanel(schedule, csrfToken, justBooked).Render(ctx, w); err != nil {
		h.logger.Error("render booking panel", "err", err)
		return
	}
	if slot, ok := schedule.FindSlot(slotID); ok {
		if err := healthcare.SlotRowOOB(slot, csrfToken).Render(ctx, w); err != nil {
			h.logger.Error("render slot row oob", "err", err)
		}
	}
}

func (h *Handler) HealthcareCancel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)
	slotID := r.FormValue("slot_id")

	_ = h.gateway.Dispatch(ctx, sessionID, healthcare.CancelBooking{SlotID: slotID})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/healthcare", http.StatusSeeOther)
		return
	}

	proj, err := h.gateway.Project(ctx, sessionID, healthcare.ScheduleProjectionName)
	if err != nil {
		h.logger.Error("project healthcare schedule", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	schedule := proj.(healthcare.ScheduleProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := healthcare.BookingPanel(schedule, csrfToken, "").Render(ctx, w); err != nil {
		h.logger.Error("render booking panel", "err", err)
		return
	}
	if slot, ok := schedule.FindSlot(slotID); ok {
		if err := healthcare.SlotRowOOB(slot, csrfToken).Render(ctx, w); err != nil {
			h.logger.Error("render slot row oob", "err", err)
		}
	}
}
