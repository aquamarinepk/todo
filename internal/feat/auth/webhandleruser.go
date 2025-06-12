package auth

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// User handlers
func (h *WebHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of users")
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		h.Err(w, err, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(r, users)
	page.SetFormAction(authPath)
	//

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

// NewUser handles the creation of a new user.
// WIP: This is a work in progress. The flash message system is still not available
// to deliver notifications. Some tweaking is still needed to properly display
// flash messages in the template.
func (h *WebHandler) NewUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New user")

	// WIPT: Just testing the flash mesages.
	// Still not working
	err := h.AddInfoFlash(w, r, "Welcome to the user creation page!")
	if err != nil {
		h.Log().Error("add info flash message error", err)
	}

	err = h.AddSuccessFlash(w, r, "This is a success message!")
	if err != nil {
		h.Log().Error("add success flash message error", err)
	}

	err = h.AddWarningFlash(w, r, "Please fill in all required fields!")
	if err != nil {
		h.Log().Error("add warning flash message error", err)
	}

	err = h.AddErrorFlash(w, r, "This is an error message!")
	if err != nil {
		h.Log().Error("add error flash message error", err)
	}

	user := NewUser("", "")

	page := am.NewPage(r, user)
	page.SetFormAction(fmt.Sprintf("%s/create-user", authPath))
	page.SetFormButtonText("Create")

	// Convert auth.Flash to am.Flash
	authFlash := h.GetFlash(r)
	amFlash := am.Flash{}
	for _, n := range authFlash.Notifications {
		amFlash.Add(string(n.Type), n.Msg)
	}
	page.SetFlash(amFlash)

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

func (h *WebHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := UserForm{}

	err := am.ToForm(r, &user)
	if err != nil {
		h.Err(w, err, ErrInvalidFormData, http.StatusBadRequest)
		return
	}

	validation, err := ValidateUser(user)
	if err != nil {
		h.Err(w, err, ErrValidationFailed, http.StatusBadRequest)
		return
	}

	if validation.HasErrors() {
		for _, err := range validation.Errors {
			h.AddFlash(w, r, am.NotificationType.Error, err)
		}
		http.Redirect(w, r, am.NewPath(authPath, "user"), http.StatusSeeOther)
		return
	}

	encKey := h.Cfg().ByteSliceVal(am.Key.SecEncryptionKey)

	newUser, err := FormToUser(user, encKey)
	if err != nil {
		h.Err(w, err, ErrCannotCreateUser, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	err = h.service.CreateUser(ctx, newUser)
	if err != nil {
		h.Err(w, err, ErrCannotCreateUser, http.StatusInternalServerError)
		return
	}

	// Set success flash message
	err = h.AddFlash(w, r, am.NotificationType.Success, "User created successfully")
	if err != nil {
		h.Log().Error("Failed to add flash message", err)
	}

	http.Redirect(w, r, am.ListPath(authPath, "user"), http.StatusSeeOther)
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

	page := am.NewPage(r, user)

	menu := am.NewMenu(authPath)
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

	page := am.NewPage(r, &user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "update", am.NoSlug))
	page.SetFormButtonText("Update")

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

func (h *WebHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
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
			h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
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

// User relationships
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

	page := am.NewPage(r, struct {
		User            User
		AssignedRoles   []Role
		UnassignedRoles []Role
	}{
		User:            user,
		AssignedRoles:   assignedRoles,
		UnassignedRoles: unassignedRoles,
	})

	page.SetFormAction("/auth/add-role-to-user")

	menu := am.NewMenu(authPath)

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
	page := am.NewPage(r, struct {
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

	// Create the menu
	menu := am.NewMenu(authPath)

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

func (h *WebHandler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add role to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}
	roleIDStr := r.FormValue("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
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

func (h *WebHandler) AddPermissionToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add permission to user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
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

	http.Redirect(w, r, fmt.Sprintf("list-user-permissions?id=%s", userID), http.StatusSeeOther)
}

func (h *WebHandler) RemovePermissionFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove permission from user")
	ctx := r.Context()

	userIDStr := r.FormValue("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	permissionIDStr := r.FormValue("permission_id")
	permissionID, err := uuid.Parse(permissionIDStr)
	if err != nil {
		h.Err(w, err, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	err = h.service.RemovePermissionFromUser(ctx, userID, permissionID)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("list-user-permissions?id=%s", userID), http.StatusSeeOther)
}
