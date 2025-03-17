package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
)

const (
	userPathFmt = "%s/%s-user%s"
)

var authPath = "/feat/auth"

var (
	key    = am.Key
	Ñmethod = am.HTTPMethod
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

	user := &User{}

	page := am.NewPage(user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "create", am.NoSlug))
	page.GenCSRFToken(r)

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
	slug := r.URL.Query().Get("slug")
	h.Log().Info("Show user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	// TODO: This is presentation logic and will be moved to a more appropriate place later.
	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)
	blue, _ := cfg.StrVal(key.ButtonStyleBlue)
	red, _ := cfg.StrVal(key.ButtonStyleRed)

	page := am.NewPage(user)
	page.SetActions([]am.Action{
		{URL: authPath, Text: "Back to User", Style: gray},
		{URL: fmt.Sprintf(userPathFmt, authPath, "edit", fmt.Sprintf(am.Slug, slug)), Text: "Edit User", Style: blue},
		{URL: fmt.Sprintf(userPathFmt, authPath, "delete", fmt.Sprintf(am.Slug, slug)), Text: "Delete User", Style: red},
		{URL: fmt.Sprintf("list-user-roles?slug=%s", slug), Text: "Manage Roles", Style: blue},
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
	slug := r.URL.Query().Get("slug")
	h.Log().Info("Edit user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(&user)
	page.SetFormAction(fmt.Sprintf(userPathFmt, authPath, "update", am.NoSlug))
	page.GenCSRFToken(r)

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

	err := h.service.CreateUser(ctx, user)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authPath, http.StatusSeeOther)
}

func (h *WebHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	slug := r.FormValue("slug")
	h.Log().Info("Update user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
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
	slug := r.FormValue("slug")
	h.Log().Info("Delete user", slug)
	ctx := r.Context()

	err := h.service.DeleteUser(ctx, slug)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "list-users", http.StatusSeeOther)
}

func (h *WebHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	h.Log().Info("List roles for user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	roles, err := h.service.GetUserRoles(ctx, slug)
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

	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")

	err := h.service.AddRole(ctx, userSlug, roleSlug)
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userSlug), http.StatusSeeOther)
}

func (h *WebHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add role to user")
	ctx := r.Context()

	userSlug := r.FormValue("user_slug")
	name := r.FormValue("name")
	description := r.FormValue("description")
	status := r.FormValue("status")
	role := NewRole(name, description, status) // TODO: This should be obtained from the DB.

	err := h.service.AddRole(ctx, userSlug, role.Slug())
	if err != nil {
		h.Err(w, err, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userSlug), http.StatusSeeOther)
}

func (h *WebHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove role from user")
	ctx := r.Context()

	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")

	err := h.service.RemoveRole(ctx, userSlug, roleSlug)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userSlug), http.StatusSeeOther)
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

func (h *WebHandler) EditRole(w http.ResponseWriter, r *http.Request) {
	userSlug := chi.URLParam(r, "user_slug")
	roleSlug := chi.URLParam(r, "role_slug")
	h.Log().Info("Edit role", userSlug, roleSlug)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, userSlug, roleSlug)
	if err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(role)
	page.SetFormAction(fmt.Sprintf("%s/%s/roles/%s/edit", authPath, userSlug, roleSlug))
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

func (h *WebHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")
	h.Log().Info("Update role", userSlug, roleSlug)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, userSlug, roleSlug)
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

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userSlug), http.StatusSeeOther)
}

func (h *WebHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")
	h.Log().Info("Delete role", userSlug, roleSlug)
	ctx := r.Context()

	err := h.service.DeleteRole(ctx, userSlug, roleSlug)
	if err != nil {
		h.Err(w, err, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", authPath, userSlug), http.StatusSeeOther)
}

// TODO Following handlers are not implemented yet.
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

func (h *WebHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement RemoveRoleFromUser")
	// TODO: Implement this handler
}

func (h *WebHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ListPermissions")
	// TODO: Implement this handler
}

func (h *WebHandler) NewPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement NewPermission")
	// TODO: Implement this handler
}

func (h *WebHandler) ShowPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ShowPermission")
	// TODO: Implement this handler
}

func (h *WebHandler) EditPermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement EditPermission")
	// TODO: Implement this handler
}

func (h *WebHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement CreatePermission")
	// TODO: Implement this handler
}

func (h *WebHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement UpdatePermission")
	// TODO: Implement this handler
}

func (h *WebHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement DeletePermission")
	// TODO: Implement this handler
}

func (h *WebHandler) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement AddPermissionToRole")
	// TODO: Implement this handler
}

func (h *WebHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement RemovePermissionFromRole")
	// TODO: Implement this handler
}

func (h *WebHandler) AddPermissionToUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement AddPermissionToUser")
	// TODO: Implement this handler
}

func (h *WebHandler) RemovePermissionFromUser(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement RemovePermissionFromUser")
	// TODO: Implement this handler
}

func (h *WebHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ListResources")
	// TODO: Implement this handler
}

func (h *WebHandler) NewResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement NewResource")
	// TODO: Implement this handler
}

func (h *WebHandler) ShowResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ShowResource")
	// TODO: Implement this handler
}

func (h *WebHandler) EditResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement EditResource")
	// TODO: Implement this handler
}

func (h *WebHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement CreateResource")
	// TODO: Implement this handler
}

func (h *WebHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement UpdateResource")
	// TODO: Implement this handler
}

func (h *WebHandler) DeleteResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement DeleteResource")
	// TODO: Implement this handler
}

func (h *WebHandler) ListResourcePermissions(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement ListResourcePermissions")
	// TODO: Implement this handler
}

func (h *WebHandler) AddPermissionToResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement AddPermissionToResource")
	// TODO: Implement this handler
}

func (h *WebHandler) RemovePermissionFromResource(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Implement RemovePermissionFromResource")
	// TODO: Implement this handler
}

// SetupPaths sets the paths for the feature.
// A default path is set by default during creation of the handler.
// This hook allows to update using the configuration during setup.
func (h *WebHandler) SetupPaths(ctx context.Context) {
	featPath := h.Cfg().StrValOrDef(key.ServerFeatPath, "/feat")
	authPath = fmt.Sprintf("%s/auth", featPath)
}

// Setup is the default implementation for the Setup method in WebHandler.
func (h *WebHandler) Setup(ctx context.Context) error {
	err := h.Handler.Setup(ctx)
	if err != nil {
		return err
	}

	h.SetupPaths(ctx)

	return nil
}

// Err logs the error and returns the message with the code.
// TODO: Move this to am.Handler to make it available to all handlers.
func (h *WebHandler) Err(w http.ResponseWriter, err error, msg string, code int) {
	h.Log().Error(err)
	renderErr := h.Cfg().BoolVal(key.RenderWebErrors, false)
	if renderErr {
		http.Error(w, fmt.Sprintf("%s: %s", msg, err.Error()), code)
		return
	}
	http.Error(w, msg, code)
}
