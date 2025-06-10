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

	owners, err := h.service.GetOrgOwners(ctx, org.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		Org    Org
		Owners []User
	}{
		Org:    org,
		Owners: owners,
	})
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(org)
	// menu.AddGenericItem("list-org-members", org.ID().String(), "Members")
	page.Menu = *menu

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
