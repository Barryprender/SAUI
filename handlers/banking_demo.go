package handlers

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"saui/demo/banking"
	"saui/middleware"
)

func (h *Handler) BankingDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)

	// Seed balance on first visit; ignore error if already funded.
	_ = h.gateway.Dispatch(ctx, sessionID, banking.FundAccount{})

	proj, err := h.gateway.Project(ctx, sessionID, banking.AccountProjectionName)
	if err != nil {
		h.logger.Error("project banking account", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	account := proj.(banking.AccountProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := banking.Dashboard(account).Render(ctx, w); err != nil {
		h.logger.Error("render banking dashboard", "err", err)
	}
}

func (h *Handler) BankingTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	proj, err := h.gateway.Project(ctx, sessionID, banking.AccountProjectionName)
	if err != nil {
		h.logger.Error("project banking account", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	account := proj.(banking.AccountProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := banking.TransferForm(account, csrfToken, "").Render(ctx, w); err != nil {
		h.logger.Error("render banking transfer", "err", err)
	}
}

func (h *Handler) BankingTransferSubmit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	recipientID := r.FormValue("recipient_id")
	amountStr := r.FormValue("amount_euros")

	amountCents, parseErr := parseEuros(amountStr)
	if parseErr != nil || amountCents <= 0 {
		h.renderTransferWithErr(w, r, sessionID, csrfToken, "Enter a valid amount")
		return
	}

	dispatchErr := h.gateway.Dispatch(ctx, sessionID, banking.SendTransfer{
		RecipientID: recipientID,
		AmountCents: amountCents,
	})
	if dispatchErr != nil {
		h.renderTransferWithErr(w, r, sessionID, csrfToken, unwrapMsg(dispatchErr))
		return
	}

	http.Redirect(w, r, "/demo/banking/confirmation", http.StatusSeeOther)
}

func (h *Handler) BankingConfirmation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)

	proj, err := h.gateway.Project(ctx, sessionID, banking.AccountProjectionName)
	if err != nil {
		h.logger.Error("project banking account", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	account := proj.(banking.AccountProjection)

	if len(account.Transfers) == 0 {
		http.Redirect(w, r, "/demo/banking", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := banking.TransferConfirmation(account.Transfers[0], account).Render(ctx, w); err != nil {
		h.logger.Error("render banking confirmation", "err", err)
	}
}

func (h *Handler) renderTransferWithErr(w http.ResponseWriter, r *http.Request, sessionID, csrfToken, msg string) {
	ctx := r.Context()
	proj, err := h.gateway.Project(ctx, sessionID, banking.AccountProjectionName)
	if err != nil {
		h.logger.Error("project banking account", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	account := proj.(banking.AccountProjection)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := banking.TransferForm(account, csrfToken, msg).Render(ctx, w); err != nil {
		h.logger.Error("render banking transfer", "err", err)
	}
}

func parseEuros(s string) (int64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(f * 100)), nil
}

func unwrapMsg(err error) string {
	if uw := errors.Unwrap(err); uw != nil {
		return uw.Error()
	}
	return err.Error()
}
