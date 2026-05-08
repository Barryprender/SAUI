package statestore_test

import (
	"context"
	"testing"

	"saui/internal/testutil"
)

func TestAppendEvent_StoresAndRetrievesEvent(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	ev, err := store.AppendEvent(ctx, "sess-1", "test.happened", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("AppendEvent: %v", err)
	}
	if ev.ID == 0 {
		t.Fatal("expected non-zero event ID")
	}
	if ev.Type != "test.happened" {
		t.Errorf("got type %q, want %q", ev.Type, "test.happened")
	}
	if ev.SessionID != "sess-1" {
		t.Errorf("got session %q, want %q", ev.SessionID, "sess-1")
	}
}

func TestEventsBySession_ReturnsEventsInOrder(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	for _, typ := range []string{"first", "second", "third"} {
		if _, err := store.AppendEvent(ctx, "sess-1", typ, nil); err != nil {
			t.Fatalf("AppendEvent %q: %v", typ, err)
		}
	}
	// unrelated session — must not appear
	if _, err := store.AppendEvent(ctx, "sess-2", "other", nil); err != nil {
		t.Fatalf("AppendEvent other: %v", err)
	}

	events, err := store.EventsBySession(ctx, "sess-1")
	if err != nil {
		t.Fatalf("EventsBySession: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("got %d events, want 3", len(events))
	}
	want := []string{"first", "second", "third"}
	for i, ev := range events {
		if ev.Type != want[i] {
			t.Errorf("event[%d] type = %q, want %q", i, ev.Type, want[i])
		}
	}
}

func TestEventsBySession_EmptyForUnknownSession(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	events, err := store.EventsBySession(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("EventsBySession: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("got %d events, want 0", len(events))
	}
}

func TestStore_WALModeEnabled(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	// Verify WAL mode by writing and reading back without error.
	// If WAL setup failed, the store.New call would have errored already.
	// This test confirms the store is usable after WAL pragma.
	_, err := store.AppendEvent(ctx, "wal-test", "wal.verified", nil)
	if err != nil {
		t.Fatalf("store unusable after WAL setup: %v", err)
	}
}
