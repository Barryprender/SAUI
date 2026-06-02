package saas

import (
	"context"
	"errors"
	"strings"

	"saui/statestore"
)

// SeedTeam populates the initial eight team members once per session.
type SeedTeam struct{}

func (a SeedTeam) Type() string { return TeamSeeded }

func (a SeedTeam) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	n, err := s.CountEventsBySessionAndType(ctx, sessionID, UserInvited)
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("already seeded")
	}
	return nil
}

func (a SeedTeam) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	var last statestore.Event
	for _, m := range preseededMembers {
		ev, err := s.AppendEvent(ctx, sessionID, UserInvited, map[string]any{
			"user_id": m.UserID,
			"name":    m.Name,
			"email":   m.Email,
			"role":    m.Role,
		})
		if err != nil {
			return statestore.Event{}, err
		}
		last = ev
	}
	return last, nil
}

// InviteUser adds a new team member.
type InviteUser struct {
	UserID string
	Name   string
	Email  string
	Role   string
}

func (a InviteUser) Type() string { return UserInvited }

func (a InviteUser) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("name is required")
	}
	if a.Role != "admin" && a.Role != "member" && a.Role != "viewer" {
		return errors.New("invalid role")
	}
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	proj := buildTeamFromEvents(events)
	for _, m := range proj.allMembers {
		if m.UserID == a.UserID && !m.Archived {
			return errors.New("a team member with this name already exists")
		}
	}
	return nil
}

func (a InviteUser) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, UserInvited, map[string]any{
		"user_id": a.UserID,
		"name":    a.Name,
		"email":   a.Email,
		"role":    a.Role,
	})
}

// ArchiveUser removes a team member from the active roster.
type ArchiveUser struct {
	UserID string
}

func (a ArchiveUser) Type() string { return UserArchived }

func (a ArchiveUser) Validate(ctx context.Context, s *statestore.Store, sessionID string) error {
	events, err := s.EventsBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	proj := buildTeamFromEvents(events)
	for _, m := range proj.allMembers {
		if m.UserID == a.UserID && !m.Archived {
			return nil
		}
	}
	return errors.New("member not found or already archived")
}

func (a ArchiveUser) Apply(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Event, error) {
	return s.AppendEvent(ctx, sessionID, UserArchived, map[string]any{
		"user_id": a.UserID,
	})
}
