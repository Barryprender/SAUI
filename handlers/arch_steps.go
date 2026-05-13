package handlers

import (
	"net/http"

	"saui/templates/partials"
)

func (h *Handler) ArchStep(w http.ResponseWriter, r *http.Request) {
	step := r.PathValue("step")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := partials.ArchStep(step).Render(r.Context(), w); err != nil {
		h.logger.Error("render arch step", "step", step, "err", err)
	}
}
