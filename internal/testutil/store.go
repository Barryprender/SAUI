package testutil

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"saui/statestore"
)

// NewStore returns a Store and Gateway backed by a private database file, and
// closes them when the test ends.
//
// Not ":memory:". database/sql pools connections and SQLite gives every
// connection its own anonymous in-memory database, so a write and the read
// that checks it could land on different databases — a test passed or failed
// depending on which connection the pool handed out. A file under t.TempDir()
// is one database, isolated per test, and removed with the test.
func NewStore(t *testing.T) (*statestore.Store, *statestore.Gateway) {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	store, err := statestore.New(filepath.Join(t.TempDir(), "state.db"), logger)
	if err != nil {
		t.Fatalf("testutil.NewStore: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	gw := statestore.NewGateway(store, logger)
	return store, gw
}
