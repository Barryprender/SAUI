package mfe

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"saui/statestore"
)

// SeedWorkspace populates 5 preseed tasks on first visit.
type SeedWorkspace struct{}

func (a SeedWorkspace) Type() string { return WorkspaceSeeded }

func (a SeedWorkspace) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	n, err := s.CountEventsBySessionAndType(ctx, sessionID, WorkspaceSeeded)
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("already seeded")
	}
	return nil
}

func (a SeedWorkspace) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	if _, err := s.AppendEvent(ctx, sessionID, WorkspaceSeeded, map[string]any{}); err != nil {
		return statestore.Event{}, err
	}
	var last statestore.Event
	for i, title := range preseedTitles {
		ev, err := s.AppendEvent(ctx, sessionID, TaskCreated, map[string]any{
			"task_id": taskID(i + 1),
			"title":   title,
		})
		if err != nil {
			return statestore.Event{}, err
		}
		last = ev
	}
	return last, nil
}

// CreateTask adds a new task.
type CreateTask struct {
	TaskID string
	Title  string
}

func (a CreateTask) Type() string { return TaskCreated }

func (a CreateTask) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if strings.TrimSpace(a.Title) == "" {
		return errors.New("title is required")
	}
	return nil
}

func (a CreateTask) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, TaskCreated, map[string]any{
		"task_id": a.TaskID,
		"title":   a.Title,
	})
}

// ToggleTask flips the completion state of a task. Type() returns a sentinel
// value used only for gateway logging; the actual event type is decided in Apply.
type ToggleTask struct {
	TaskID string
}

func (a ToggleTask) Type() string { return "task.toggle" }

func (a ToggleTask) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if strings.TrimSpace(a.TaskID) == "" {
		return errors.New("task_id is required")
	}
	return nil
}

func (a ToggleTask) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	completed, err := currentTaskCompletedFromStore(ctx, s, sessionID, a.TaskID)
	if err != nil {
		return statestore.Event{}, err
	}
	eventType := TaskCompleted
	if completed {
		eventType = TaskReopened
	}
	return s.AppendEvent(ctx, sessionID, eventType, map[string]any{
		"task_id": a.TaskID,
	})
}

// NextTaskID generates a task ID based on timestamp to avoid collisions.
func NextTaskID() string {
	return fmt.Sprintf("task-u%d", time.Now().UnixMilli())
}
