package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	authPath    = "/auth"
	userPathFmt = "%s/%s-user%s"
)

const (
	ActionListUserRoles       = "list-user-roles"
	ActionListUserPermissions = "list-user-permissions"
	TextRoles                 = "Roles"
	TextPermissions           = "Permissions"
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
	h.Log().Info("New user")

	user := NewUser("", "")

	page := struct {
		User User
	}{
		User: user,
	}

	h.Render(w, r, "new-user", page)
}

func (h *WebHandler) ShowUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Show user ", id)
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
	menu.AddGenericItem(ActionListUserPermissions, user.ID().String(), TextPermissions)

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
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Edit user ", id)
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
	password := r.FormValue("password")

	encKey := h.Cfg().ByteSliceVal(am.Key.SecEncryptionKey)
	emailEnc, err := EncryptEmail(email, encKey)
	if err != nil {
		h.Err(w, err, "Failed to encrypt email", http.StatusInternalServerError)
		return
	}

	passwordEnc, err := HashPassword(password)
	if err != nil {
		h.Err(w, err, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := NewUser(username, name)
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)

	err = h.service.CreateUser(ctx, user)
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

	h.Log().Info("Update user ", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	username := r.FormValue("username")
	name := r.FormValue("name")
	password := r.FormValue("password")

	// Update user fields
	user.Username = username
	user.Name = name

	// Update password if provided
	if password != "" {
		passwordEnc, err := HashPassword(password)
		if err != nil {
			h.Err(w, err, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.PasswordEnc = passwordEnc
	}

	err = h.service.UpdateUser(ctx, user)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "user"), http.StatusSeeOther)
}

func (h *WebHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Update user password for user ", id)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// NOTE: This will be done briefly with a to-be-implemented validator.
	password := r.FormValue("password")
	if password == "" {
		h.Err(w, fmt.Errorf("password is required"), "Password is required", http.StatusBadRequest)
		return
	}

	user.Password = password

	err = h.service.UpdateUserPassword(ctx, user)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "user"), http.StatusSeeOther)
}

func (h *WebHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Delete user ", id)
	ctx := r.Context()

	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "user"), http.StatusSeeOther)
}

func (h *WebHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	var err error
	var userID uuid.UUID
	userID, err = h.ID(w, r)
	if err != nil {
		return
	}

	var user User
	user, err = h.service.GetUser(r.Context(), userID)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	var assignedRoles []Role
	assignedRoles, err = h.service.GetUserRoles(r.Context(), userID)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	var unassignedRoles []Role
	unassignedRoles, err = h.service.GetUserUnassignedRoles(r.Context(), userID)
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

	page.SetFormAction("/auth/add-role-to-user")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddShowItem(user, "Back")

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

func (h *WebHandler) ListUserPermissions(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	h.Log().Info("List permissions for user ", "id", id)
	ctx := r.Context()

	// Get user details
	user, err := h.service.GetUser(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// Get permissions assigned to the user through roles
	permissionsFromRoles, err := h.service.GetUserIndirectPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	// Get permissions directly assigned to the user
	directPermissions, err := h.service.GetUserDirectPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	// Get permissions not assigned to the user (neither through roles nor directly)
	unassignedPermissions, err := h.service.GetUserUnassignedPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	// Prepare the page data
	page := am.NewPage(struct {
		User                  User
		PermissionsFromRoles  []Permission
		DirectPermissions     []Permission
		UnassignedPermissions []Permission
	}{
		User:                  user,
		PermissionsFromRoles:  permissionsFromRoles,
		DirectPermissions:     directPermissions,
		UnassignedPermissions: unassignedPermissions,
	})
	page.GenCSRFToken(r)

	// Create the menu
	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddShowItem(user, "Back")

	page.Menu = *menu

	// Render the template
	tmpl, err := h.tm.Get("auth", "list-user-permissions")
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

func (h *WebHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List roles")
	ctx := r.Context()

	roles, err := h.service.GetAllRoles(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(roles)
	page.SetFormAction(authPath)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddNewItem(roleType)

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "list-roles")
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

func (h *WebHandler) NewRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New role form")

	role := NewRole("", "", "active")

	page := am.NewPage(role)
	page.SetFormAction(fmt.Sprintf("%s/create-role", authPath))
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(role)

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "new-role")
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

func (h *WebHandler) ShowRole(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Show role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(role)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(role)
	menu.AddEditItem(role)
	menu.AddDeleteItem(role)
	menu.AddGenericItem("list-role-permissions", role.ID().String(), "Permissions")

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "show-role")
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

func (h *WebHandler) EditRole(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		return
	}

	h.Log().Info("Edit role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(role)
	page.SetFormAction(fmt.Sprintf("%s/update-role", authPath))
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.AddListItem(role)

	page.Menu = *menu

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

	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	status := r.Form.Get("status")
	if status == "" {
		status = "active"
	}

	if name == "" {
		h.Err(w, nil, "Name is required", http.StatusBadRequest)
		return
	}

	role := NewRole(name, description, status)

	err := h.service.CreateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "role"), http.StatusSeeOther)
}

func (h *WebHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.Err(w, err, "Invalid form data", http.StatusBadRequest)
		return
	}

	idStr := r.Form.Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Update role ", id)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	role.Name = r.Form.Get("name")
	role.Description = r.Form.Get("description")
	role.SetSlug(r.Form.Get("slug"))

	err = h.service.UpdateRole(ctx, role)
	if err != nil {
		h.Err(w, err, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "role"), http.StatusSeeOther)
}

func (h *WebHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}

	h.Log().Info("Delete role ", id)
	ctx := r.Context()

	err = h.service.DeleteRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "role"), http.StatusSeeOther)
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

// SECTION: Permission handlers

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

	page := am.NewPage(permission)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
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

func (h *WebHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create permission")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	permission := NewPermission(name, description)
	permission.GenID()
	permission.GenSlug()
	permission.GenCreationValues()

	err := h.service.CreatePermission(ctx, permission)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, am.ListPath(authPath, "permission"), http.StatusSeeOther)
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
	permission.Model = am.NewModel(
		am.WithID(permission.ID()),
		am.WithType(permissionType),
		am.WithSlug(name),
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

	http.Redirect(w, r, fmt.Sprintf("%s/list-user-permissions?id=%s", authPath, userID), http.StatusSeeOther)
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

	http.Redirect(w, r, fmt.Sprintf("%s/list-user-permissions?id=%s", authPath, userID), http.StatusSeeOther)
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

	err = h.service.AddPermissionToRole(ctx, roleID, permissionID)
	if err != nil {
		// Check if the error is due to a unique constraint violation
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// Permission is already assigned to the role, redirect back to the permissions page
			http.Redirect(w, r, fmt.Sprintf("list-role-permissions?id=%s", roleID), http.StatusSeeOther)
			return
		}
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("list-role-permissions?id=%s", roleID), http.StatusSeeOther)
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

	http.Redirect(w, r, fmt.Sprintf("list-role-permissions?id=%s", roleID), http.StatusSeeOther)
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

	page := am.NewPage(resource)
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(resource)
	menu.AddEditItem(resource)
	menu.AddDeleteItem(resource)
	menu.AddGenericItem("list-resource-permissions", resource.ID().String(), "Permissions")

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
	resourceType := r.FormValue("type")

	resource := NewResource(name, description, resourceType)
	resource.GenSlug()
	resource.Model = am.NewModel(
		am.WithID(resource.ID()),
		am.WithSlug(resource.Slug()),
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
	resource.GenSlug()
	resource.Model = am.NewModel(
		am.WithID(resource.ID()),
		am.WithType(resourceEntityType),
		am.WithSlug(resource.Slug()),
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

	http.Redirect(w, r, fmt.Sprintf("list-resource-permissions?id=%s", resourceID), http.StatusSeeOther)
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

	http.Redirect(w, r, fmt.Sprintf("list-resource-permissions?id=%s", resourceID), http.StatusSeeOther)
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
	h.Log().Info("Get role ", id)
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
	h.Log().Info("Get permission ", id)
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
	h.Log().Info("Get user ", id)
	ctx := r.Context()

	if _, err := h.service.GetUser(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) Err(w http.ResponseWriter, err error, msg string, code int) {
	h.Log().Error(msg, err)
	http.Error(w, fmt.Sprintf("%s: %v", msg, err), code)
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

	page := am.NewPage(resource)
	page.SetFormAction(am.UpdatePath(authPath, "resource"))
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(resource)
	menu.AddShowItem(resource)
	menu.AddDeleteItem(resource)

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

	page := am.NewPage(struct {
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
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(resource)
	menu.AddEditItem(resource)
	menu.AddDeleteItem(resource)
	menu.AddGenericItem("list-resource-permissions", resource.ID().String(), "Permissions")

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

func (h *WebHandler) ListRolePermissions(w http.ResponseWriter, r *http.Request) {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	h.Log().Info("Showing role permissions ", "id ", id)

	ctx := r.Context()
	role, err := h.service.GetRole(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	assigned, err := h.service.GetRolePermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	unassigned, err := h.service.GetRoleUnassignedPermissions(ctx, id)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(struct {
		ID                   uuid.UUID
		Name                 string
		Description          string
		Permissions          []Permission
		AvailablePermissions []Permission
	}{
		ID:                   role.ID(),
		Name:                 role.Name,
		Description:          role.Description,
		Permissions:          assigned,
		AvailablePermissions: unassigned,
	})
	page.GenCSRFToken(r)

	menu := am.NewMenu(authPath)
	menu.SetCSRFToken(page.Form.CSRF)
	menu.AddListItem(role)
	menu.AddEditItem(role)
	menu.AddDeleteItem(role)
	menu.AddGenericItem("list-role-permissions", role.ID().String(), "Permissions")

	page.Menu = *menu

	tmpl, err := h.tm.Get("auth", "list-role-permissions")
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
