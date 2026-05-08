package projections_test

import (
	"context"
	"testing"

	"saui/actions"
	"saui/internal/testutil"
	"saui/projections"
)

func TestFeedbackProjection_EmptyForNewSession(t *testing.T) {
	_, gw := testutil.NewStore(t)
	ctx := context.Background()

	p, err := gw.Project(ctx, "sess-new", projections.FeedbackProjectionName)
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	fp, ok := p.(projections.FeedbackProjection)
	if !ok {
		t.Fatalf("unexpected projection type %T", p)
	}
	if fp.Count != 0 {
		t.Errorf("got count %d, want 0", fp.Count)
	}
	if len(fp.Messages) != 0 {
		t.Errorf("got %d messages, want 0", len(fp.Messages))
	}
}

func TestFeedbackProjection_ReflectsSubmissions(t *testing.T) {
	store, gw := testutil.NewStore(t)
	ctx := context.Background()
	sess := "sess-proj"

	msgs := []string{"first message", "second message"}
	for _, msg := range msgs {
		a := actions.SubmitFeedback{Message: msg}
		if _, err := a.Apply(ctx, store, sess); err != nil {
			t.Fatalf("Apply %q: %v", msg, err)
		}
	}

	p, err := gw.Project(ctx, sess, projections.FeedbackProjectionName)
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	fp, ok := p.(projections.FeedbackProjection)
	if !ok {
		t.Fatalf("unexpected projection type %T", p)
	}
	if fp.Count != 2 {
		t.Errorf("got count %d, want 2", fp.Count)
	}
	for i, want := range msgs {
		if fp.Messages[i] != want {
			t.Errorf("message[%d] = %q, want %q", i, fp.Messages[i], want)
		}
	}
}

func TestFeedbackProjection_IsolatedBySession(t *testing.T) {
	store, gw := testutil.NewStore(t)
	ctx := context.Background()

	a := actions.SubmitFeedback{Message: "only for sess-a"}
	if _, err := a.Apply(ctx, store, "sess-a"); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	p, err := gw.Project(ctx, "sess-b", projections.FeedbackProjectionName)
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	fp, ok := p.(projections.FeedbackProjection)
	if !ok {
		t.Fatalf("unexpected projection type %T", p)
	}
	if fp.Count != 0 {
		t.Errorf("sess-b projection should be empty, got count %d", fp.Count)
	}
}

func TestGateway_Dispatch_RejectsInvalidAction(t *testing.T) {
	_, gw := testutil.NewStore(t)
	ctx := context.Background()

	err := gw.Dispatch(ctx, "sess-1", actions.SubmitFeedback{Message: ""})
	if err == nil {
		t.Fatal("expected Dispatch to return error for invalid action, got nil")
	}
}

func TestGateway_Dispatch_AppliesValidAction(t *testing.T) {
	_, gw := testutil.NewStore(t)
	ctx := context.Background()

	if err := gw.Dispatch(ctx, "sess-1", actions.SubmitFeedback{Message: "valid"}); err != nil {
		t.Fatalf("Dispatch: %v", err)
	}

	p, err := gw.Project(ctx, "sess-1", projections.FeedbackProjectionName)
	if err != nil {
		t.Fatalf("Project: %v", err)
	}
	fp, ok := p.(projections.FeedbackProjection)
	if !ok {
		t.Fatalf("unexpected projection type %T", p)
	}
	if fp.Count != 1 {
		t.Errorf("got count %d after dispatch, want 1", fp.Count)
	}
}
