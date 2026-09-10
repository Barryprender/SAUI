package handlers

import (
	"net/http"

	"saui/locale"
	"saui/middleware"
	"saui/templates/pages"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := pages.Home(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render home", "err", err)
	}
}

func (h *Handler) HomeES(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/es/" {
		http.NotFound(w, r)
		return
	}
	if err := pages.Home(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render home es", "err", err)
	}
}

func (h *Handler) Why(w http.ResponseWriter, r *http.Request) {
	if err := pages.Why(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render why", "err", err)
	}
}

func (h *Handler) WhyES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Why(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render why es", "err", err)
	}
}

func (h *Handler) Architecture(w http.ResponseWriter, r *http.Request) {
	if err := pages.Architecture(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render architecture", "err", err)
	}
}

func (h *Handler) ArchitectureES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Architecture(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render architecture es", "err", err)
	}
}

func (h *Handler) Stack(w http.ResponseWriter, r *http.Request) {
	if err := pages.Stack(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render stack", "err", err)
	}
}

func (h *Handler) StackES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Stack(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render stack es", "err", err)
	}
}

func (h *Handler) Cases(w http.ResponseWriter, r *http.Request) {
	if err := pages.Cases(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render cases", "err", err)
	}
}

func (h *Handler) CasesES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Cases(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render cases es", "err", err)
	}
}

func (h *Handler) CaseFoodOrdering(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseFoodOrdering(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case food-ordering", "err", err)
	}
}

func (h *Handler) CaseFoodOrderingES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseFoodOrdering(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case food-ordering es", "err", err)
	}
}

func (h *Handler) CaseBanking(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseBanking(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case banking", "err", err)
	}
}

func (h *Handler) CaseBankingES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseBanking(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case banking es", "err", err)
	}
}

func (h *Handler) CaseHealthcare(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseHealthcare(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case healthcare", "err", err)
	}
}

func (h *Handler) CaseHealthcareES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseHealthcare(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case healthcare es", "err", err)
	}
}

func (h *Handler) CaseSaaS(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseSaaS(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case saas", "err", err)
	}
}

func (h *Handler) CaseSaaSES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseSaaS(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case saas es", "err", err)
	}
}

func (h *Handler) CaseDistributed(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseDistributed(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case distributed", "err", err)
	}
}

func (h *Handler) CaseDistributedES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseDistributed(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case distributed es", "err", err)
	}
}

func (h *Handler) CaseMFE(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseMFE(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render case mfe", "err", err)
	}
}

func (h *Handler) CaseMFEES(w http.ResponseWriter, r *http.Request) {
	if err := pages.CaseMFE(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render case mfe es", "err", err)
	}
}

func (h *Handler) Testing(w http.ResponseWriter, r *http.Request) {
	if err := pages.Testing(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render testing", "err", err)
	}
}

func (h *Handler) TestingES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Testing(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render testing es", "err", err)
	}
}

func (h *Handler) Limits(w http.ResponseWriter, r *http.Request) {
	if err := pages.Limits(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render limits", "err", err)
	}
}

func (h *Handler) LimitsES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Limits(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render limits es", "err", err)
	}
}

func (h *Handler) Blog(w http.ResponseWriter, r *http.Request) {
	if err := pages.Blog(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog", "err", err)
	}
}

func (h *Handler) BlogES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Blog(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog es", "err", err)
	}
}

func (h *Handler) BlogPostSupplyChain(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostSupplyChain(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post", "err", err)
	}
}

func (h *Handler) BlogPostSupplyChainES(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostSupplyChain(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post es", "err", err)
	}
}

func (h *Handler) BlogPostServerResponseTime(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostServerResponseTime(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post", "err", err)
	}
}

func (h *Handler) BlogPostServerResponseTimeES(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostServerResponseTime(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post es", "err", err)
	}
}

func (h *Handler) BlogPostEUCompliance(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostEUCompliance(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post", "err", err)
	}
}

func (h *Handler) BlogPostEUComplianceES(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostEUCompliance(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post es", "err", err)
	}
}

func (h *Handler) BlogPostCRAClock(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostCRAClock(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post", "err", err)
	}
}

func (h *Handler) BlogPostCRAClockES(w http.ResponseWriter, r *http.Request) {
	if err := pages.BlogPostCRAClock(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render blog post es", "err", err)
	}
}

func (h *Handler) Code(w http.ResponseWriter, r *http.Request) {
	if err := pages.Code(middleware.CSRFToken(r), locale.EN).Render(r.Context(), w); err != nil {
		h.logger.Error("render code", "err", err)
	}
}

func (h *Handler) CodeES(w http.ResponseWriter, r *http.Request) {
	if err := pages.Code(middleware.CSRFToken(r), locale.ES).Render(r.Context(), w); err != nil {
		h.logger.Error("render code es", "err", err)
	}
}
