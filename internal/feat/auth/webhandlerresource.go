package auth

import (
	"bytes"
	"net/http"
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Resource handlers
func (h *WebHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List resources")
	ctx := r.Context()

	resources, err := h.service.GetAllResources(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, resources)
	page.SetFormAction(authPath)

	menu := page.NewMenu(authPath)
	menu.AddNewItem("resource")

	tmpl, err := h.tm.Get("auth", "list-resources")
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

func (h *WebHandler) NewResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New resource form")

	resource := NewResource("", "", "entity")

	page := am.NewPage(r, resource)
	page.SetFormAction(am.CreatePath(authPath, "resource"))
	page.SetFormButtonText("Create")

	menu := page.NewMenu(authPath)
	menu.AddListItem(resource)

	tmpl, err := h.tm.Get("auth", "new-resource")
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

func (h *WebHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create resource")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	resourceType := r.FormValue("type")

	resource := NewResource(name, description, resourceType)
	resource.BaseModel = am.NewModel(
		am.WithID(resource.ID()),
		am.WithCreatedBy(uuid.New()),
		am.WithUpdatedBy(uuid.New()),
		am.WithCreatedAt(time.Now()),
		am.WithUpdatedAt(time.Now()),
	)

	err := h.service.CreateResource(ctx, resource)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "resource"), http.StatusSeeOther)
}

func (h *WebHandler) ShowResource(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Show resource ", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, resource)

	menu := page.NewMenu(authPath)

	menu.AddListItem(resource)
	menu.AddEditItem(resource)
	menu.AddDeleteItem(resource)
	menu.AddGenericItem("list-resource-permissions", resource.ID().String(), "Permissions")

	tmpl, err := h.tm.Get("auth", "show-resource")
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

func (h *WebHandler) EditResource(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Edit resource ", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(r, resource)
	page.SetFormAction(am.UpdatePath(authPath, "resource"))

	menu := page.NewMenu(authPath)

	menu.AddListItem(resource)
	menu.AddShowItem(resource)
	menu.AddDeleteItem(resource)

	tmpl, err := h.tm.Get("auth", "edit-resource")
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

func (h *WebHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.Form.Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	resource, err := h.service.GetResource(r.Context(), id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")

	resource.Name = name
	resource.Description = description
	resource.BaseModel = am.NewModel(
		am.WithID(resource.ID()),
		am.WithType(resourceEntityType),
		am.WithCreatedBy(resource.CreatedBy()),
		am.WithUpdatedBy(uuid.New()),
		am.WithCreatedAt(resource.CreatedAt()),
		am.WithUpdatedAt(time.Now()),
	)

	if err := h.service.UpdateResource(r.Context(), resource); err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "resource"), http.StatusSeeOther)
}

func (h *WebHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if err := h.service.DeleteResource(ctx, id); err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "resource"), http.StatusSeeOther)
}

// Resource relationships
func (h *WebHandler) ListResourcePermissions(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	h.Log().Info("Showing resource permissions ", "id", id)

	ctx := r.Context()
	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	assigned, err := h.service.GetResourcePermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassigned, err := h.service.GetResourceUnassignedPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, struct {
		ID                   uuid.UUID
		Name                 string
		Description          string
		Permissions          []Permission
		AvailablePermissions []Permission
	}{
		ID:                   resource.ID(),
		Name:                 resource.Name,
		Description:          resource.Description,
		Permissions:          assigned,
		AvailablePermissions: unassigned,
	})

	menu := page.NewMenu(authPath)

	menu.AddListItem(resource)
	menu.AddEditItem(resource)
	menu.AddDeleteItem(resource)
	menu.AddGenericItem("list-resource-permissions", resource.ID().String(), "Permissions")

	tmpl, err := h.tm.Get("auth", "list-resource-permissions")
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

func (h *WebHandler) AddPermissionToResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add permission to resource")
	ctx := r.Context()

	resourceIDStr := r.FormValue("resource_id")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermission(ctx, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	err = h.service.AddPermissionToResource(ctx, resourceID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listResourcePermissionsPath) + "?id=" + resourceID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (h *WebHandler) RemovePermissionFromResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove permission from resource")
	ctx := r.Context()

	resourceIDStr := r.FormValue("resource_id")
	resourceID, err := uuid.Parse(resourceIDStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemovePermissionFromResource(ctx, resourceID, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	path := am.ListPath(authPath, listResourcePermissionsPath) + "?id=" + resourceID.String()
	http.Redirect(w, r, path, http.StatusSeeOther)
}
