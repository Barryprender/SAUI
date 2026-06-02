package handlers

import (
	"net/http"

	"saui/demo/saas"
	"saui/middleware"
)

func (h *Handler) SAASDashboard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	// Seed eight team members on first visit; ignore error if already seeded.
	_ = h.gateway.Dispatch(ctx, sessionID, saas.SeedTeam{})

	proj, err := h.gateway.Project(ctx, sessionID, saas.TeamProjectionName)
	if err != nil {
		h.logger.Error("project saas team", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	team := proj.(saas.TeamProjection)

	roleFilter := r.URL.Query().Get("role")
	if roleFilter != "admin" && roleFilter != "member" && roleFilter != "viewer" {
		roleFilter = ""
	}
	team.RoleFilter = roleFilter
	if roleFilter != "" {
		var filtered []saas.TeamMember
		for _, m := range team.Members {
			if m.Role == roleFilter {
				filtered = append(filtered, m)
			}
		}
		team.Members = filtered
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := saas.Dashboard(team, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render saas dashboard", "err", err)
	}
}

func (h *Handler) SAASInvite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	name := r.FormValue("name")
	role := r.FormValue("role")
	userID := saas.Slugify(name)

	_ = h.gateway.Dispatch(ctx, sessionID, saas.InviteUser{
		UserID: userID,
		Name:   name,
		Email:  userID + "@demo.io",
		Role:   role,
	})

	http.Redirect(w, r, "/demo/saas", http.StatusSeeOther)
}

func (h *Handler) SAASArchive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	userID := r.FormValue("user_id")
	roleFilter := r.FormValue("role_filter")

	_ = h.gateway.Dispatch(ctx, sessionID, saas.ArchiveUser{UserID: userID})

	redirect := "/demo/saas"
	if roleFilter != "" {
		redirect += "?role=" + roleFilter
	}
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}
