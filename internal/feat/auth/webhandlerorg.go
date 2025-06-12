package auth

import (
	"bytes"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
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

	page := am.NewPage(r, org)

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
