package handlers

import (
	"log/slog"

	"saui/statestore"
)

type Handler struct {
	gateway *statestore.Gateway
	logger  *slog.Logger
}

func New(gateway *statestore.Gateway, logger *slog.Logger) *Handler {
	return &Handler{gateway: gateway, logger: logger}
}
