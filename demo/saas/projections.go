package saas

import (
	"context"
	"encoding/json"

	"saui/statestore"
)

const TeamProjectionName = "saas.team"

// TeamMember is a single member in the team roster.
type TeamMember struct {
	UserID      string
	Name        string
	Email       string
	Role        string
	Archived    bool
	EventID     int64
	AvatarClass string
	Initials    string
}

// teamState holds all members before filtering — used internally during projection build and action validation.
type teamState struct {
	allMembers []TeamMember
}

// TeamProjection is the read model for the SaaS dashboard.
type TeamProjection struct {
	Members     []TeamMember // active, filtered
	RoleFilter  string
	CountAll    int
	CountAdmin  int
	CountMember int
	CountViewer int
}

func (p TeamProjection) Name() string { return TeamProjectionName }

func init() {
	statestore.RegisterProjection(TeamProjectionName, func(ctx context.Context, s *statestore.Store, sessionID string) (statestore.Projection, error) {
		events, err := s.EventsBySession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		state := buildTeamFromEvents(events)

		proj := TeamProjection{}
		for _, m := range state.allMembers {
			if m.Archived {
				continue
			}
			proj.CountAll++
			switch m.Role {
			case "admin":
				proj.CountAdmin++
			case "member":
				proj.CountMember++
			case "viewer":
				proj.CountViewer++
			}
			proj.Members = append(proj.Members, m)
		}
		return proj, nil
	})
}

func buildTeamFromEvents(events []statestore.Event) teamState {
	members := make(map[string]TeamMember)
	var order []string

	for _, ev := range events {
		switch ev.Type {
		case UserInvited:
			var p struct {
				UserID string `json:"user_id"`
				Name   string `json:"name"`
				Email  string `json:"email"`
				Role   string `json:"role"`
			}
			if json.Unmarshal(ev.Payload, &p) != nil {
				continue
			}
			if _, exists := members[p.UserID]; !exists {
				order = append(order, p.UserID)
			}
			members[p.UserID] = TeamMember{
				UserID:      p.UserID,
				Name:        p.Name,
				Email:       p.Email,
				Role:        p.Role,
				Archived:    false,
				EventID:     ev.ID,
				AvatarClass: RoleAvatarClass(p.Role),
				Initials:    Initials(p.Name),
			}
		case UserArchived:
			var p struct {
				UserID string `json:"user_id"`
			}
			if json.Unmarshal(ev.Payload, &p) != nil {
				continue
			}
			if m, exists := members[p.UserID]; exists {
				m.Archived = true
				members[p.UserID] = m
			}
		}
	}

	var out []TeamMember
	for _, id := range order {
		out = append(out, members[id])
	}
	return teamState{allMembers: out}
}
