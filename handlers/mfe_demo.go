package handlers

import (
	"net/http"
	"strings"

	"saui/demo/mfe"
	"saui/middleware"
)

func (h *Handler) MFEWorkspace(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	_ = h.gateway.Dispatch(ctx, sessionID, mfe.SeedWorkspace{})

	proj, err := h.gateway.Project(ctx, sessionID, mfe.WorkspaceProjectionName)
	if err != nil {
		h.logger.Error("project mfe workspace", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	workspace := proj.(mfe.WorkspaceProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mfe.Workspace(workspace, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render mfe workspace", "err", err)
	}
}

func (h *Handler) MFEToggle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	taskID := strings.TrimSpace(r.FormValue("task_id"))
	if taskID == "" {
		http.Error(w, "task_id required", http.StatusBadRequest)
		return
	}

	_ = h.gateway.Dispatch(ctx, sessionID, mfe.ToggleTask{TaskID: taskID})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/mfe", http.StatusSeeOther)
		return
	}

	proj, err := h.gateway.Project(ctx, sessionID, mfe.WorkspaceProjectionName)
	if err != nil {
		h.logger.Error("project mfe workspace", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	workspace := proj.(mfe.WorkspaceProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mfe.AlphaPanel(workspace, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render mfe alpha", "err", err)
		return
	}
	if err := mfe.BetaPanelOOB(workspace).Render(ctx, w); err != nil {
		h.logger.Error("render mfe beta oob", "err", err)
		return
	}
	if err := mfe.GammaPanelOOB(workspace).Render(ctx, w); err != nil {
		h.logger.Error("render mfe gamma oob", "err", err)
	}
}

func (h *Handler) MFECreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sessionID := middleware.SessionID(ctx)
	csrfToken := middleware.CSRFToken(r)

	title := strings.TrimSpace(r.FormValue("title"))

	_ = h.gateway.Dispatch(ctx, sessionID, mfe.CreateTask{
		TaskID: mfe.NextTaskID(),
		Title:  title,
	})

	if !middleware.IsHXRequest(r) {
		http.Redirect(w, r, "/demo/mfe", http.StatusSeeOther)
		return
	}

	proj, err := h.gateway.Project(ctx, sessionID, mfe.WorkspaceProjectionName)
	if err != nil {
		h.logger.Error("project mfe workspace", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	workspace := proj.(mfe.WorkspaceProjection)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := mfe.AlphaPanel(workspace, csrfToken).Render(ctx, w); err != nil {
		h.logger.Error("render mfe alpha", "err", err)
		return
	}
	if err := mfe.BetaPanelOOB(workspace).Render(ctx, w); err != nil {
		h.logger.Error("render mfe beta oob", "err", err)
		return
	}
	if err := mfe.GammaPanelOOB(workspace).Render(ctx, w); err != nil {
		h.logger.Error("render mfe gamma oob", "err", err)
	}
}
