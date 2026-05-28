package handlers

import (
	"net/http"

	"saui/demo/food"
	"saui/middleware"
)

func (h *Handler) FoodMenu(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	menuProj, err := h.gateway.Project(ctx, sessionID, food.MenuProjectionName)
	if err != nil {
		h.logger.Error("project food menu", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	cartProj, err := h.gateway.Project(ctx, sessionID, food.CartProjectionName)
	if err != nil {
		h.logger.Error("project food cart", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	menu := menuProj.(food.MenuProjection)
	cart := cartProj.(food.CartProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := food.Menu(menu, cart, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render food menu", "err", err)
	}
}

func (h *Handler) FoodCartAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	dispatchErr := h.gateway.Dispatch(ctx, sessionID, food.AddToCart{
		ItemID: r.FormValue("item_id"),
	})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/food-ordering", http.StatusSeeOther)
		return
	}

	cartProj, err := h.gateway.Project(ctx, sessionID, food.CartProjectionName)
	if err != nil {
		h.logger.Error("project food cart", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	cart := cartProj.(food.CartProjection)

	var errMsg string
	if dispatchErr != nil {
		errMsg = dispatchErr.Error()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := food.CartPanel(cart, csrfToken, errMsg).Render(ctx, w); err != nil {
		h.logger.Error("render food cart", "err", err)
	}
}

func (h *Handler) FoodCartRemove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	dispatchErr := h.gateway.Dispatch(ctx, sessionID, food.RemoveFromCart{
		ItemID: r.FormValue("item_id"),
	})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/food-ordering", http.StatusSeeOther)
		return
	}

	cartProj, err := h.gateway.Project(ctx, sessionID, food.CartProjectionName)
	if err != nil {
		h.logger.Error("project food cart", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	cart := cartProj.(food.CartProjection)

	var errMsg string
	if dispatchErr != nil {
		errMsg = dispatchErr.Error()
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := food.CartPanel(cart, csrfToken, errMsg).Render(ctx, w); err != nil {
		h.logger.Error("render food cart", "err", err)
	}
}

func (h *Handler) FoodCheckout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)

	if err := h.gateway.Dispatch(ctx, sessionID, food.PlaceOrder{}); err != nil {
		h.logger.Warn("food place order failed", "err", err, "session", sessionID)
		http.Redirect(w, r, "/demo/food-ordering", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/demo/food-ordering/confirmation", http.StatusSeeOther)
}

func (h *Handler) FoodConfirmation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)

	orderProj, err := h.gateway.Project(ctx, sessionID, food.OrderProjectionName)
	if err != nil {
		h.logger.Error("project food order", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	order := orderProj.(food.OrderProjection)

	if !order.HasOrder {
		http.Redirect(w, r, "/demo/food-ordering", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := food.Confirmation(order).Render(ctx, w); err != nil {
		h.logger.Error("render food confirmation", "err", err)
	}
}
