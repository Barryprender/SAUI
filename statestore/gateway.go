package statestore

import (
	"context"
	"fmt"
	"log/slog"
)

// Action is a validated intent to change state.
type Action interface {
	Type() string
	Validate(ctx context.Context, s *Store, sessionID string) error
	Apply(ctx context.Context, s *Store, sessionID string) (Event, error)
}

// Gateway is the sole point of access for all state reads and writes.
// No handler or middleware touches the Store directly.
type Gateway struct {
	store     *Store
	ephemeral *ephemeralStore
	logger    *slog.Logger
}

func NewGateway(store *Store, logger *slog.Logger) *Gateway {
	return &Gateway{
		store:     store,
		ephemeral: newEphemeralStore(),
		logger:    logger,
	}
}

// Project returns the named projection for the given session.
func (g *Gateway) Project(ctx context.Context, sessionID, name string) (Projection, error) {
	fn, ok := projections[name]
	if !ok {
		return nil, fmt.Errorf("unknown projection: %s", name)
	}
	return fn(ctx, g.store, sessionID)
}

// Dispatch validates and applies an action, appending the resulting event to the durable log.
func (g *Gateway) Dispatch(ctx context.Context, sessionID string, action Action) error {
	if err := action.Validate(ctx, g.store, sessionID); err != nil {
		return fmt.Errorf("action %s invalid: %w", action.Type(), err)
	}
	event, err := action.Apply(ctx, g.store, sessionID)
	if err != nil {
		return fmt.Errorf("action %s failed: %w", action.Type(), err)
	}
	g.logger.Info("event appended", "type", event.Type, "session", sessionID, "id", event.ID)
	return nil
}

// Ephemeral returns the ephemeral state for a session (never used in Dispatch).
func (g *Gateway) Ephemeral(sessionID string) *EphemeralState {
	return g.ephemeral.Get(sessionID)
}

// SetEphemeral updates ephemeral state for a session.
func (g *Gateway) SetEphemeral(sessionID string, state *EphemeralState) {
	g.ephemeral.Set(sessionID, state)
}

// ClearEphemeral removes ephemeral state on session end.
func (g *Gateway) ClearEphemeral(sessionID string) {
	g.ephemeral.Delete(sessionID)
}
