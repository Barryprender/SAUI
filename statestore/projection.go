package statestore

import "context"

// Projection is a read model derived from the event log.
type Projection interface {
	Name() string
}

// ProjectionFunc derives a Projection from the current Store for a given session.
type ProjectionFunc func(ctx context.Context, s *Store, sessionID string) (Projection, error)

var projections = map[string]ProjectionFunc{}

// RegisterProjection registers a named projection builder.
func RegisterProjection(name string, fn ProjectionFunc) {
	projections[name] = fn
}
