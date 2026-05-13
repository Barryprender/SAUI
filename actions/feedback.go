package actions

import (
	"context"
	"errors"
	"strings"

	"saui/statestore"
)

const FeedbackSubmittedType = "feedback.submitted"

const MaxFeedbackPerSession = 3

var ErrFeedbackLimitReached = errors.New("feedback limit reached for this session")

type SubmitFeedback struct {
	Message string
}

func (a SubmitFeedback) Type() string { return FeedbackSubmittedType }

func (a SubmitFeedback) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if strings.TrimSpace(a.Message) == "" {
		return errors.New("message is required")
	}
	if len([]rune(a.Message)) > 2000 {
		return errors.New("message exceeds 2000 characters")
	}

	count, err := s.CountEventsBySessionAndType(ctx, sessionID, FeedbackSubmittedType)
	if err != nil {
		return err
	}
	if count >= MaxFeedbackPerSession {
		return ErrFeedbackLimitReached
	}
	return nil
}

func (a SubmitFeedback) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, FeedbackSubmittedType, map[string]string{
		"message": strings.TrimSpace(a.Message),
	})
}
