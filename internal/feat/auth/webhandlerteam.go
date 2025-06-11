package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Team handlers
func (h *WebHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, err := h.service.GetDefaultOrg(ctx)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	teams, err := h.service.GetAllTeams(ctx, org.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		Org   Org
		Teams []Team
	}{
		Org:   org,
		Teams: teams,
	})
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddNewItem("team")
	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "list-teams")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) NewTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	org, err := h.service.GetDefaultOrg(ctx)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	team := NewTeam(org.ID(), "", "", "")
	page := am.NewPage(team)
	page.SetFormAction("/auth/create-team")
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(team)
	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "new-team")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	org, err := h.service.GetDefaultOrg(ctx)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	name := r.Form.Get("name")
	shortDescription := r.Form.Get("short_description")
	description := r.Form.Get("description")

	team := NewTeam(org.ID(), name, shortDescription, description)
	team.BaseModel = am.NewModel(
		am.WithID(team.ID()),
		am.WithCreatedBy(team.CreatedBy()),
		am.WithUpdatedBy(team.UpdatedBy()),
		am.WithCreatedAt(team.CreatedAt()),
		am.WithUpdatedAt(team.UpdatedAt()),
	)

	err = h.service.CreateTeam(ctx, team)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/auth/list-teams", http.StatusSeeOther)
}

func (h *WebHandler) ShowTeam(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	// You will need to implement GetTeam in the service/repo layer
	team, err := h.service.GetTeam(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(team)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(team)
	menu.AddEditItem(team)
	menu.AddGenericItem(ActionListTeamMembers, team.ID().String(), TextMembers)
	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "show-team")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) EditTeam(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	team, err := h.service.GetTeam(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(team)
	page.SetFormAction("/auth/update-team")
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(team)
	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "edit-team")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	idStr := r.Form.Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid team ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Update team ", id)
	ctx := r.Context()

	team, err := h.service.GetTeam(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	team.Name = r.Form.Get("name")
	team.ShortDescription = r.Form.Get("short_description")
	team.Description = r.Form.Get("description")
	team.BaseModel = am.NewModel(
		am.WithID(team.ID()),
		am.WithCreatedBy(team.CreatedBy()),
		am.WithUpdatedBy(uuid.New()),
		am.WithCreatedAt(team.CreatedAt()),
		am.WithUpdatedAt(time.Now()),
	)

	err = h.service.UpdateTeam(ctx, team)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "team"), http.StatusSeeOther)
}

func (h *WebHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteTeam(ctx, id); err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "team"), http.StatusSeeOther)
}

func (h *WebHandler) ListTeamMembers(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	team, err := h.service.GetTeam(ctx, id)
	if err != nil {
		h.Log().Error("Failed to get team", "id", id, "error", err)
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	members, err := h.service.GetTeamMembers(ctx, id)
	if err != nil {
		h.Log().Error("Failed to get team members", "team_id", id, "error", err)
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassigned, err := h.service.GetTeamUnassignedUsers(ctx, id)
	if err != nil {
		h.Log().Error("Failed to get unassigned users for team", "team_id", id, "error", err)
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		Team       Team
		Members    []User
		Unassigned []User
	}{
		Team:       team,
		Members:    members,
		Unassigned: unassigned,
	})

	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)

	menu.AddListItem(team)
	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "list-team-members")
	if err != nil {
		h.Err(w, err, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		h.Err(w, err, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		h.Err(w, err, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) AssignUserToTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamIDStr := r.FormValue("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	// For now we'll use "member" as the default relation type
	err = h.service.AddUserToTeam(ctx, teamID, userID, "member")
	if err != nil {
		h.Log().Error("Failed to assign user to team", "team_id", teamID, "user_id", userID, "error", err)
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	// Add success flash message
	err = h.AddFlash(w, r, am.NotificationType.Success, "User assigned to team successfully")
	if err != nil {
		h.Log().Error("Failed to add flash message", err)
	}

	// Redirect back to team members page
	http.Redirect(w, r, fmt.Sprintf("/auth/list-team-members?id=%s", teamID), http.StatusSeeOther)
}

func (h *WebHandler) RemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamIDStr := r.FormValue("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	err = h.service.RemoveUserFromTeam(ctx, teamID, userID)
	if err != nil {
		h.Log().Error("Failed to remove user from team", "team_id", teamID, "user_id", userID, "error", err)
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	// Add success flash message
	err = h.AddFlash(w, r, am.NotificationType.Success, "User removed from team successfully")
	if err != nil {
		h.Log().Error("Failed to add flash message", err)
	}

	// Redirect back to team members page
	http.Redirect(w, r, fmt.Sprintf("/auth/list-team-members?id=%s", teamID), http.StatusSeeOther)
}
