package auth

import (
	"bytes"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Org handlers
func (h *WebHandler) ShowOrg(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show org")
	ctx := r.Context()

	org, err := h.service.GetDefaultOrg(ctx)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	owners, err := h.service.GetOrgOwners(ctx, org.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, struct {
		Org    Org
		Owners []User
	}{
		Org:    org,
		Owners: owners,
	})

	menu := page.NewMenu(authPath)
	menu.AddGenericItem("list-org-owners", org.ID().String(), "Owners")

	tmpl, err := h.tm.Get("auth", "show-org")
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

// Organization relationships
func (h *WebHandler) ListOrgOwners(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	h.Log().Info("List org owners", "id", id)
	ctx := r.Context()

	org, err := h.service.GetDefaultOrg(ctx)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	owners, err := h.service.GetOrgOwners(ctx, org.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassigned, err := h.service.GetOrgUnassignedOwners(ctx, org.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, struct {
		Org        Org
		Owners     []User
		Unassigned []User
	}{
		Org:        org,
		Owners:     owners,
		Unassigned: unassigned,
	})

	menu := page.NewMenu(authPath)
	menu.AddShowItem(org, "Back")

	tmpl, err := h.tm.Get("auth", "list-org-owners")
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

func (h *WebHandler) AddOrgOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgIDStr := r.FormValue("org_id")
	orgID, err := uuid.Parse(orgIDStr)
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

	err = h.service.AddOrgOwner(ctx, orgID, userID)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/auth/list-org-owners?id="+orgID.String(), http.StatusSeeOther)
}

func (h *WebHandler) RemoveOrgOwner(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orgIDStr := r.FormValue("org_id")
	orgID, err := uuid.Parse(orgIDStr)
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

	err = h.service.RemoveOrgOwner(ctx, orgID, userID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/auth/list-org-owners?id="+orgID.String(), http.StatusSeeOther)
}
