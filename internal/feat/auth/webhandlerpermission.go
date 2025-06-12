package auth

import (
	"bytes"
	"net/http"
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Permission handlers
func (h *WebHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List permissions")
	ctx := r.Context()

	permissions, err := h.service.GetAllPermissions(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, permissions)
	page.SetFormAction(authPath)

	menu := am.NewMenu(authPath)
	menu.AddNewItem("permission")

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "list-permissions")
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

func (h *WebHandler) NewPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New permission form")

	permission := NewPermission("", "")

	page := am.NewPage(r, permission)
	page.SetFormAction(am.CreatePath(authPath, "permission"))
	page.SetFormButtonText("Create")

	menu := am.NewMenu(authPath)
	menu.AddListItem(permission)

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "new-permission")
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

func (h *WebHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create permission")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	permission := NewPermission(name, description)
	permission.GenCreateValues()

	err := h.service.CreatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
}

func (h *WebHandler) ShowPermission(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Show permission ", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, permission)

	menu := am.NewMenu(authPath)

	menu.AddListItem(permission)
	menu.AddEditItem(permission)
	menu.AddDeleteItem(permission)

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "show-permission")
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

func (h *WebHandler) EditPermission(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Edit permission ", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, permission)
	page.SetFormAction(am.UpdatePath(authPath, "permission"))
	page.SetFormButtonText("Update")

	menu := am.NewMenu(authPath)
	menu.AddListItem(permission)

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "edit-permission")
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

func (h *WebHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	idStr := r.Form.Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Update permission ", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	permission.Name = name
	permission.Description = description
	permission.BaseModel = am.NewModel(
		am.WithID(permission.ID()),
		am.WithType(permissionType),
		am.WithCreatedBy(permission.CreatedBy()),
		am.WithUpdatedBy(uuid.New()),
		am.WithCreatedAt(permission.CreatedAt()),
		am.WithUpdatedAt(time.Now()),
	)

	err = h.service.UpdatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
}

func (h *WebHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Delete permission")
	ctx := r.Context()

	err = h.service.DeletePermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
}
