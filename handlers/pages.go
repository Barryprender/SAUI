package handlers

import (
	"net/http"

	"saui/middleware"
	"saui/templates/pages"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := pages.Home(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render home", "err", err)
	}
}

func (h *Handler) Why(w http.ResponseWriter, r *http.Request) {
	if err := pages.Why(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render why", "err", err)
	}
}

func (h *Handler) Architecture(w http.ResponseWriter, r *http.Request) {
	if err := pages.Architecture(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render architecture", "err", err)
	}
}

func (h *Handler) Stack(w http.ResponseWriter, r *http.Request) {
	if err := pages.Stack(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render stack", "err", err)
	}
}

func (h *Handler) Cases(w http.ResponseWriter, r *http.Request) {
	if err := pages.Cases(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render cases", "err", err)
	}
}

func (h *Handler) CaseFoodOrdering(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseFoodOrdering(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case food-ordering", "err", err)
	}
}

func (h *Handler) CaseBanking(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseBanking(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case banking", "err", err)
	}
}

func (h *Handler) CaseHealthcare(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseHealthcare(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case healthcare", "err", err)
	}
}

func (h *Handler) CaseSaaS(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseSaaS(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case saas", "err", err)
	}
}

func (h *Handler) CaseDistributed(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseDistributed(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case distributed", "err", err)
	}
}

func (h *Handler) CaseMFE(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseMFE(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render case mfe", "err", err)
	}
}

func (h *Handler) Testing(w http.ResponseWriter, r *http.Request) {
	if err := pages.Testing(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render testing", "err", err)
	}
}

func (h *Handler) Limits(w http.ResponseWriter, r *http.Request) {
	if err := pages.Limits(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render limits", "err", err)
	}
}

func (h *Handler) Blog(w http.ResponseWriter, r *http.Request) {
	if err := pages.Blog(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog", "err", err)
	}
}

func (h *Handler) Code(w http.ResponseWriter, r *http.Request) {
	if err := pages.Code(middleware.CSRFToken(r)).Render(r.Context(), w); err != nil {
		h.logger.Error("render code", "err", err)
	}
}
