package actions_test

import (
	"context"
	"strings"
	"testing"

	"saui/actions"
	"saui/internal/testutil"
)

func TestSubmitFeedback_Validate_RejectsEmptyMessage(t *testing.T) {
	store, _ := testutil.NewStore(t)
	a := actions.SubmitFeedback{Message: "   "}
	if err := a.Validate(context.Background(), store, "sess-1"); err == nil {
		t.Fatal("expected error for empty message, got nil")
	}
}

func TestSubmitFeedback_Validate_RejectsOversizedMessage(t *testing.T) {
	store, _ := testutil.NewStore(t)
	a := actions.SubmitFeedback{Message: strings.Repeat("a", 2001)}
	if err := a.Validate(context.Background(), store, "sess-1"); err == nil {
		t.Fatal("expected error for oversized message, got nil")
	}
}

func TestSubmitFeedback_Validate_AcceptsValidMessage(t *testing.T) {
	store, _ := testutil.NewStore(t)
	a := actions.SubmitFeedback{Message: "This is valid feedback."}
	if err := a.Validate(context.Background(), store, "sess-1"); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestSubmitFeedback_Validate_EnforcesRateLimit(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()
	sess := "sess-ratelimit"

	// Submit up to the limit.
	for i := 0; i < 3; i++ {
		a := actions.SubmitFeedback{Message: "message"}
		if err := a.Validate(ctx, store, sess); err != nil {
			t.Fatalf("submission %d: unexpected validation error: %v", i+1, err)
		}
		if _, err := a.Apply(ctx, store, sess); err != nil {
			t.Fatalf("submission %d: Apply failed: %v", i+1, err)
		}
	}

	// Fourth submission must be rejected.
	a := actions.SubmitFeedback{Message: "one too many"}
	if err := a.Validate(ctx, store, sess); err == nil {
		t.Fatal("expected rate limit error on 4th submission, got nil")
	}
}

func TestSubmitFeedback_RateLimitIsPerSession(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	// Exhaust session A.
	for i := 0; i < 3; i++ {
		a := actions.SubmitFeedback{Message: "msg"}
		if _, err := a.Apply(ctx, store, "sess-a"); err != nil {
			t.Fatalf("Apply sess-a: %v", err)
		}
	}

	// Session B must still be accepted.
	a := actions.SubmitFeedback{Message: "independent session"}
	if err := a.Validate(ctx, store, "sess-b"); err != nil {
		t.Fatalf("sess-b should not be rate limited: %v", err)
	}
}

func TestSubmitFeedback_Apply_AppendsEvent(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()
	a := actions.SubmitFeedback{Message: "  hello  "}

	ev, err := a.Apply(ctx, store, "sess-1")
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if ev.Type != actions.FeedbackSubmittedType {
		t.Errorf("got event type %q, want %q", ev.Type, actions.FeedbackSubmittedType)
	}
	if ev.SessionID != "sess-1" {
		t.Errorf("got session %q, want %q", ev.SessionID, "sess-1")
	}
}

func TestSubmitFeedback_Apply_TrimsMessage(t *testing.T) {
	store, _ := testutil.NewStore(t)
	ctx := context.Background()

	a := actions.SubmitFeedback{Message: "  trimmed  "}
	ev, err := a.Apply(ctx, store, "sess-1")
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if !strings.Contains(string(ev.Payload), "trimmed") {
		t.Errorf("expected trimmed message in payload, got: %s", ev.Payload)
	}
	if strings.Contains(string(ev.Payload), "  trimmed  ") {
		t.Errorf("message was not trimmed in payload: %s", ev.Payload)
	}
}
