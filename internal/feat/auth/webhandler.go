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
	featBasePath = "/feat"
	authFeatPath = "/auth"
	userPathFmt  = "%s/%s-user%s"
)

var (
	authPath = fmt.Sprintf("%s%s", featBasePath, authFeatPath)
)

const (
	ActionListUserRoles = "list-user-roles"
	TextRoles           = "Roles"
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

func (h *WebHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of users")
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(users)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddNewItem(userType)

	page.Menu = *menu

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

	page := am.NewPage(user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "create", am.NoSlug))
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(user)

	page.Menu = *menu

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

	page := am.NewPage(user)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(user)
	menu.AddEditItem(user)
	menu.AddDeleteItem(user)
	menu.AddGenericItem(ActionListUserRoles, user.ID().String(), TextRoles)

	page.Menu = *menu

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

	page := am.NewPage(&user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "update", am.NoSlug))
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(user)

	page.Menu = *menu

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
	ctx := r.Context()

	username := r.FormValue("username")
	email := r.FormValue("email")
	name := r.FormValue("name")
	user := NewUser(username, email, name)
	user.GenID()
	user.GenSlug()
	user.GenCreationValues()

	err := h.service.CreateUser(ctx, user)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "user"), http.StatusSeeOther)
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
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	assignedRoles, err := h.service.GetUserRoles(r.Context(), id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassignedRoles, err := h.service.GetUserUnassignedRoles(r.Context(), id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		User            User
		AssignedRoles   []Role
		UnassignedRoles []Role
	}{
		User:            user,
		AssignedRoles:   assignedRoles,
		UnassignedRoles: unassignedRoles,
	})

	// Set form action for adding roles
	page.SetFormAction("/auth/add-role-to-user")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	// menu.AddGenericItem("add-role-to-user", user.ID().String(), "Add Role to User")

	page.Menu = *menu

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

	http.Redirect(w, r, fmt.Sprintf("list-user-roles?id=%s", userID), http.StatusSeeOther)
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

	http.Redirect(w, r, fmt.Sprintf("list-user-roles?id=%s", userID), http.StatusSeeOther)
}

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

	http.Redirect(w, r, fmt.Sprintf("list-user-roles?id=%s", userID), http.StatusSeeOther)
}

// Permission handlers
func (h *WebHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List permissions")
	ctx := r.Context()

	permissions, err := h.service.GetAllPermissions(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(permissions)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

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

func (h *WebHandler) ShowPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show permission")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(permission)
	page.SetFormAction(fmt.Sprintf("%s/%s", authPath, id))
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddEditItem(NewPermission("", ""))
	menu.AddDeleteItem(NewPermission("", ""))

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

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
}

func (h *WebHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update permission")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	permission, err := h.service.GetPermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	permission.Name = name
	permission.Description = description

	err = h.service.UpdatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ShowPath(authPath, "permission", id), http.StatusSeeOther)
}

func (h *WebHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete permission")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePermission(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
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
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	err = h.service.AddPermissionToUser(ctx, userID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ShowPath(authPath, "user", userID), http.StatusSeeOther)
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

	http.Redirect(w, r, am.ShowPath(authPath, "user", userID), http.StatusSeeOther)
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
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	err = h.service.AddPermissionToRole(ctx, roleID, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ShowPath(authPath, "role", roleID), http.StatusSeeOther)
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

	http.Redirect(w, r, am.ShowPath(authPath, "role", roleID), http.StatusSeeOther)
}

// Resource handlers
func (h *WebHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List resources")
	ctx := r.Context()

	resources, err := h.service.GetAllResources(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(resources)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddNewItem("resource")

	page.Menu = *menu

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

func (h *WebHandler) ShowResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Show resource")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(resource)
	page.SetFormAction(fmt.Sprintf("%s/%s", authPath, id))
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddEditItem(NewResource("", "", "entity"))
	menu.AddDeleteItem(NewResource("", "", "entity"))

	page.Menu = *menu

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

func (h *WebHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create resource")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	resource := NewResource(name, description, "entity")

	err := h.service.CreateResource(ctx, resource)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "resource"), http.StatusSeeOther)
}

func (h *WebHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Update resource")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	resource, err := h.service.GetResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusInternalServerError)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	resource.Name = name
	resource.Description = description

	err = h.service.UpdateResource(ctx, resource)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ShowPath(authPath, "resource", id), http.StatusSeeOther)
}

func (h *WebHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete resource")
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteResource(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "resource"), http.StatusSeeOther)
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

	http.Redirect(w, r, am.ShowPath(authPath, "resource", resourceID), http.StatusSeeOther)
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

	http.Redirect(w, r, am.ShowPath(authPath, "resource", resourceID), http.StatusSeeOther)
}

func (h *WebHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get resource", id)
	ctx := r.Context()

	if _, err := h.service.GetResource(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get role", id)
	ctx := r.Context()

	if _, err := h.service.GetRole(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get permission", id)
	ctx := r.Context()

	if _, err := h.service.GetPermission(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get user", id)
	ctx := r.Context()

	if _, err := h.service.GetUser(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) Err(w http.ResponseWriter, err error, msg string, code int) {
	h.Log().Error(msg, err)
	http.Error(w, msg, code)
}

func (h *WebHandler) addSampleData() {
	// Sample resources will be added in a future implementation
}

func (h *WebHandler) NewPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New permission form")

	permission := NewPermission("", "")

	page := am.NewPage(permission)
	page.SetFormAction(am.CreatePath(authPath, "permission"))
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

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

func (h *WebHandler) EditPermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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

	page := am.NewPage(permission)
	page.SetFormAction(am.UpdatePath(authPath, "permission"))
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

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

func (h *WebHandler) NewResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New resource form")

	resource := NewResource("", "", "entity")

	page := am.NewPage(resource)
	page.SetFormAction(am.CreatePath(authPath, "resource"))
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(resource)

	page.Menu = *menu

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

func (h *WebHandler) EditResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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

	page := am.NewPage(resource)
	page.SetFormAction(am.UpdatePath(authPath, "resource"))
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(resource)

	page.Menu = *menu

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

	menu := am.NewMenu(authPath)
	menu.AddListItem(resource)

	page.Menu = *menu

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
