package projections

import (
	"context"
	"encoding/json"

	"saui/actions"
	"saui/statestore"
)

const FeedbackProjectionName = "feedback"

type FeedbackProjection struct {
	Count    int
	Messages []string
}

func (p FeedbackProjection) Name() string { return FeedbackProjectionName }

func init() {
	statestore.RegisterProjection(FeedbackProjectionName, BuildFeedback)
}

func BuildFeedback(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	proj := FeedbackProjection{}
	for _, ev := range events {
		if ev.Type != actions.FeedbackSubmittedType {
			continue
		}
		var payload map[string]string
		if err := json.Unmarshal(ev.Payload, &payload); err != nil {
			continue
		}
		if msg, ok := payload["message"]; ok {
			proj.Messages = append(proj.Messages, msg)
			proj.Count++
		}
	}
	return proj, nil
}
