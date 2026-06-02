package mfe

import (
	"context"
	"encoding/json"
	"time"

	"saui/statestore"
)

const WorkspaceProjectionName = "mfe.workspace"

type Task struct {
	ID        string
	Title     string
	Completed bool
	Seq       int
}

type FeedEntry struct {
	EventID    int64
	TaskTitle  string
	Label      string
	ColorClass string
	OccurredAt time.Time
}

type WorkspaceStats struct {
	Total       int
	Completed   int
	Open        int
	CompletePct int
}

type WorkspaceProjection struct {
	Tasks       []Task
	Feed        []FeedEntry
	Stats       WorkspaceStats
	JustToggled string
	JustCreated string
}

func (p WorkspaceProjection) Name() string { return WorkspaceProjectionName }

func init() {
	statestore.RegisterProjection(WorkspaceProjectionName, func(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
		events, err := s.EventsBySession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		return buildWorkspaceFromEvents(events), nil
	})
}

func buildWorkspaceFromEvents(events []statestore.Event) WorkspaceProjection {
	type taskState struct {
		id        string
		title     string
		completed bool
		seq       int
	}

	taskMap := map[string]*taskState{}
	taskOrder := []string{}
	var feed []FeedEntry
	var justToggled, justCreated string

	for _, ev := range events {
		var p map[string]any
		_ = json.Unmarshal(ev.Payload, &p)
		taskIDVal, _ := p["task_id"].(string)
		titleVal, _ := p["title"].(string)

		switch ev.Type {
		case TaskCreated:
			if taskIDVal != "" && taskMap[taskIDVal] == nil {
				seq := len(taskOrder)
				taskMap[taskIDVal] = &taskState{id: taskIDVal, title: titleVal, seq: seq}
				taskOrder = append(taskOrder, taskIDVal)
			}
			feed = append(feed, FeedEntry{
				EventID:    ev.ID,
				TaskTitle:  titleVal,
				Label:      "Task created",
				ColorClass: "feed--created",
				OccurredAt: ev.OccurredAt,
			})
			justCreated = taskIDVal
			justToggled = ""
		case TaskCompleted:
			if t := taskMap[taskIDVal]; t != nil {
				t.completed = true
				titleVal = t.title
			}
			feed = append(feed, FeedEntry{
				EventID:    ev.ID,
				TaskTitle:  titleVal,
				Label:      "Task completed",
				ColorClass: "feed--completed",
				OccurredAt: ev.OccurredAt,
			})
			justToggled = taskIDVal
			justCreated = ""
		case TaskReopened:
			if t := taskMap[taskIDVal]; t != nil {
				t.completed = false
				titleVal = t.title
			}
			feed = append(feed, FeedEntry{
				EventID:    ev.ID,
				TaskTitle:  titleVal,
				Label:      "Task reopened",
				ColorClass: "feed--reopened",
				OccurredAt: ev.OccurredAt,
			})
			justToggled = taskIDVal
			justCreated = ""
		}
	}

	tasks := make([]Task, 0, len(taskOrder))
	for _, id := range taskOrder {
		if t := taskMap[id]; t != nil {
			tasks = append(tasks, Task{
				ID:        t.id,
				Title:     t.title,
				Completed: t.completed,
				Seq:       t.seq,
			})
		}
	}

	var completed int
	for _, t := range tasks {
		if t.Completed {
			completed++
		}
	}
	total := len(tasks)
	open := total - completed
	pct := 0
	if total > 0 {
		pct = (completed * 100) / total
	}

	// Reverse feed: newest first, cap at 10.
	for i, j := 0, len(feed)-1; i < j; i, j = i+1, j-1 {
		feed[i], feed[j] = feed[j], feed[i]
	}
	if len(feed) > 10 {
		feed = feed[:10]
	}

	return WorkspaceProjection{
		Tasks:       tasks,
		Feed:        feed,
		Stats:       WorkspaceStats{Total: total, Completed: completed, Open: open, CompletePct: pct},
		JustToggled: justToggled,
		JustCreated: justCreated,
	}
}

// currentTaskCompletedFromStore returns the current completed state of a task.
func currentTaskCompletedFromStore(ctx context.Context, s *statestore.Store, sessionID, taskID string) (bool, error) {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return false, err
	}
	proj := buildWorkspaceFromEvents(events)
	for _, t := range proj.Tasks {
		if t.ID == taskID {
			return t.Completed, nil
		}
	}
	return false, nil
}
