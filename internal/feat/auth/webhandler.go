package auth

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	userPathFmt = "%s/%s-user%s"
)

var authPath = "/feat/auth"

var (
	key    = am.Key
	method = am.HTTPMethod
)

type WebHandler struct {
	*am.Handler
	service Service
	tm      *am.TemplateManager
}

func NewWebHandler(tm *am.TemplateManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		service: service,
		tm:      tm,
	}
}

// User handlers
func (h *WebHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of users")
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	blue, _ := h.Cfg().StrVal(key.ButtonStyleBlue)

	page := am.NewPage(users)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)
	page.SetActions([]am.Action{
		am.NewAction(fmt.Sprintf("%s/new-user", authPath), "New User", blue),
	})

	tmpl, err := h.tm.Get("auth", "list-users")
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

func (h *WebHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New user form")

	user := NewUser("", "", "")

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)

	page := am.NewPage(user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "create", am.NoSlug))
	page.GenCSRFToken(r)
	page.SetActions([]am.Action{
		am.NewListAction(authPath, "user", gray),
	})

	tmpl, err := h.tm.Get("auth", "new-user")
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

func (h *WebHandler) ShowUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Show user", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	cfg := h.Cfg()

	// TODO: When a proper asset building pipeline for the client is implemented,
	// we can remove this non-business style configuration from the controller.
	gray, _ := cfg.StrVal(key.ButtonStyleGray)
	blue, _ := cfg.StrVal(key.ButtonStyleBlue)
	red, _ := cfg.StrVal(key.ButtonStyleRed)

	page := am.NewPage(user)

	page.GenCSRFToken(r)
	page.SetActions([]am.Action{
		am.NewListAction(authPath, "user", gray),
		am.NewEditAction(authPath, "user", id, blue),
		am.NewDeleteAction(authPath, "user", id, red),
		am.NewAction(fmt.Sprintf("%s/list-user-roles?id=%s", authPath, id), "Roles", blue),
	})

	tmpl, err := h.tm.Get("auth", "show-user")
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

func (h *WebHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Edit user", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)

	page := am.NewPage(&user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "update", am.NoSlug))
	page.GenCSRFToken(r)
	page.SetActions([]am.Action{
		am.NewListAction(authPath, "user", gray),
	})

	tmpl, err := h.tm.Get("auth", "edit-user")
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

func (h *WebHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create user")

	username := r.FormValue("username")
	email := r.FormValue("email")
	name := r.FormValue("name")
	user := NewUser(username, email, name)

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)

	page := am.NewPage(user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "create", am.NoSlug))
	page.GenCSRFToken(r)
	page.SetActions([]am.Action{
		am.NewListAction(authPath, "user", gray),
	})

	tmpl, err := h.tm.Get("auth", "new-user")
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

func (h *WebHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Update user", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Name = r.FormValue("name")

	err = h.service.UpdateUser(ctx, user)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "list-users", http.StatusSeeOther)
}

func (h *WebHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Delete user", id)
	ctx := r.Context()

	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "list-users", http.StatusSeeOther)
}

func (h *WebHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("List roles for user", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	roles, err := h.service.GetUserRoles(ctx, user.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		User  *User
		Roles []Role
	}{
		User:  &user,
		Roles: roles,
	})
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("auth", "list-user-roles")
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

func (h *WebHandler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add role to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.AddRole(ctx, userID, roleID)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove role from user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveRole(ctx, userID, roleID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) AddPermissionToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add permission to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
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
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	err = h.service.AddPermissionToUser(ctx, userID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) RemovePermissionFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove permission from user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemovePermissionFromUser(ctx, userID, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

// Role handlers
func (h *WebHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ListRoles")
	// TODO: Implement this handler
}

func (h *WebHandler) NewRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement NewRole")
	// TODO: Implement this handler
}

func (h *WebHandler) ShowRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ShowRole")
	// TODO: Implement this handler
}

func (h *WebHandler) EditRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := chi.URLParam(r, "role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Edit role", userID, roleID)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, roleID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(role)
	page.SetFormAction(fmt.Sprintf("%s/%s/roles/%s/edit", authPath, userID, roleID))
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("auth", "edit-role")
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

func (h *WebHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create role")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	status := r.FormValue("status")
	role := NewRole(name, description, status)

	err := h.service.CreateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authPath, http.StatusSeeOther)
}

func (h *WebHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Update role", userID, roleID)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, roleID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	role.Name = r.FormValue("name")
	role.Description = r.FormValue("description")

	err = h.service.UpdateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Delete role", userID, roleID)
	ctx := r.Context()

	err = h.service.DeleteRole(ctx, roleID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add role to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	description := r.FormValue("description")
	status := r.FormValue("status")
	role := NewRole(name, description, status) // TODO: This should be obtained from the DB.

	err = h.service.AddRole(ctx, userID, role.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

func (h *WebHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove role from user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemoveRole(ctx, userID, roleID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userID), http.StatusSeeOther)
}

// Permission handlers
func (h *WebHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ListPermissions")
	// TODO: Implement this handler
}

func (h *WebHandler) NewPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New permission form")

	permission := &Permission{}

	page := am.NewPage(permission)
	page.SetFormAction(fmt.Sprintf("%s/permissions/create", authPath))
	page.GenCSRFToken(r)

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

func (h *WebHandler) ShowPermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Show permission", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(permission)
	page.SetFormAction(fmt.Sprintf("%s/permissions/%s", authPath, id))
	page.GenCSRFToken(r)

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

func (h *WebHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create permission")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	permission := NewPermission(name, description)

	err := h.service.CreatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/permissions", authPath), http.StatusSeeOther)
}

func (h *WebHandler) EditPermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Edit permission", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(&permission)
	page.SetFormAction(fmt.Sprintf("%s/permissions/%s/update", authPath, id))
	page.GenCSRFToken(r)

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
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Update permission", id)
	ctx := r.Context()

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	permission.Name = r.FormValue("name")
	permission.Description = r.FormValue("description")

	err = h.service.UpdatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/permissions", authPath), http.StatusSeeOther)
}

func (h *WebHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Delete permission", id)
	ctx := r.Context()

	err = h.service.DeletePermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/permissions", authPath), http.StatusSeeOther)
}

func (h *WebHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add permission to role")
	ctx := r.Context()

	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
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
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	err = h.service.AddPermissionToRole(ctx, roleID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/roles/%s", authPath, roleID), http.StatusSeeOther)
}

func (h *WebHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove permission from role")
	ctx := r.Context()

	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.RemovePermissionFromRole(ctx, roleID, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/roles/%s", authPath, roleID), http.StatusSeeOther)
}

func (h *WebHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of resources")
	ctx := r.Context()

	resources, err := h.service.GetAllResources(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(resources)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

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

	resource := &Resource{}

	page := am.NewPage(resource)
	page.SetFormAction(fmt.Sprintf("%s/resources/create", authPath))
	page.GenCSRFToken(r)

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

func (h *WebHandler) ShowResource(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Show resource", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(resource)
	page.SetFormAction(fmt.Sprintf("%s/resources/%s", authPath, id))
	page.GenCSRFToken(r)

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
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Edit resource", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(&resource)
	page.SetFormAction(fmt.Sprintf("%s/resources/%s/update", authPath, id))
	page.GenCSRFToken(r)

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

func (h *WebHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create resource")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	resource := NewResource(name, description)

	err := h.service.CreateResource(ctx, resource)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/resources", authPath), http.StatusSeeOther)
}

func (h *WebHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Update resource", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	resource.Name = r.FormValue("name")
	resource.Description = r.FormValue("description")

	err = h.service.UpdateResource(ctx, resource)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/resources", authPath), http.StatusSeeOther)
}

func (h *WebHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Delete resource", id)
	ctx := r.Context()

	err = h.service.DeleteResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/resources", authPath), http.StatusSeeOther)
}

func (h *WebHandler) ListResourcePermissions(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("List permissions for resource", id)
	ctx := r.Context()

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	permissions, err := h.service.GetResourcePermissions(ctx, resource.ID())
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		Resource    *Resource
		Permissions []Permission
	}{
		Resource:    &resource,
		Permissions: permissions,
	})
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

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
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	err = h.service.AddPermissionToResource(ctx, resourceID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/resources/%s", authPath, resourceID), http.StatusSeeOther)
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

	http.Redirect(w, r, fmt.Sprintf("%s/resources/%s", authPath, resourceID), http.StatusSeeOther)
}

func (h *WebHandler) Err(w http.ResponseWriter, err error, msg string, code int) {
	h.Log().Error(msg, err)
	http.Error(w, msg, code)
}
