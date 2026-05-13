package handlers

import (
	"errors"
	"net/http"

	"saui/actions"
	"saui/middleware"
	"saui/projections"
	"saui/templates/partials"
)

func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)
	message := r.FormValue("message")

	err := h.gateway.Dispatch(ctx, sessionID, actions.SubmitFeedback{Message: message})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/code", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		if errors.Is(err, actions.ErrFeedbackLimitReached) {
			if renderErr := partials.FeedbackDone().Render(ctx, w); renderErr != nil {
				h.logger.Error("render feedback done", "err", renderErr)
			}
		} else {
			if renderErr := partials.FeedbackForm(csrfToken, err.Error(), "error", actions.MaxFeedbackPerSession).Render(ctx, w); renderErr != nil {
				h.logger.Error("render feedback form", "err", renderErr)
			}
		}
		return
	}

	remaining := 0
	if proj, perr := h.gateway.Project(ctx, sessionID, projections.FeedbackProjectionName); perr == nil {
		if fp, ok := proj.(projections.FeedbackProjection); ok {
			remaining = actions.MaxFeedbackPerSession - fp.Count
		}
	}

	if remaining <= 0 {
		if renderErr := partials.FeedbackDone().Render(ctx, w); renderErr != nil {
			h.logger.Error("render feedback done", "err", renderErr)
		}
	} else {
		if renderErr := partials.FeedbackForm(csrfToken, "Received. Thanks.", "success", remaining).Render(ctx, w); renderErr != nil {
			h.logger.Error("render feedback form", "err", renderErr)
		}
	}
}
