package testutil

import (
	"log/slog"
	"os"
	"testing"

	"saui/statestore"
)

// NewStore returns an in-memory Store and Gateway for use in tests.
// The store is closed automatically when the test ends.
func NewStore(t *testing.T) (*statestore.Store, *statestore.Gateway) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	store, err := statestore.New(":memory:", logger)
	if err != nil {
		t.Fatalf("testutil.NewStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	gw := statestore.NewGateway(store, logger)
	return store, gw
}
