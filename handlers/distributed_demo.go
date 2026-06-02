package handlers

import (
	"net/http"

	"saui/demo/distributed"
	"saui/middleware"
	"saui/statestore"
)

func (h *Handler) DistributedPipeline(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	proj, err := h.gateway.Project(ctx, sessionID, distributed.PipelineProjectionName)
	if err != nil {
		h.logger.Error("project distributed pipeline", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	pipeline := proj.(distributed.PipelineProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := distributed.Pipeline(pipeline, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render distributed pipeline", "err", err)
	}
}

func (h *Handler) DistributedAdvance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	var action statestore.Action
	switch r.FormValue("action") {
	case "place":
		action = distributed.PlaceOrder{}
	case "reserve":
		action = distributed.ReserveInventory{}
	case "capture":
		action = distributed.CapturePayment{}
	case "decline":
		action = distributed.DeclinePayment{}
	case "fulfill":
		action = distributed.FulfillOrder{}
	case "cancel":
		action = distributed.CancelOrder{}
	default:
		http.Error(w, "unknown action", http.StatusBadRequest)
		return
	}

	_ = h.gateway.Dispatch(ctx, sessionID, action)

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/distributed", http.StatusSeeOther)
		return
	}

	proj, err := h.gateway.Project(ctx, sessionID, distributed.PipelineProjectionName)
	if err != nil {
		h.logger.Error("project distributed pipeline", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	pipeline := proj.(distributed.PipelineProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := distributed.PipelineBoard(pipeline, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render distributed board", "err", err)
	}
}
