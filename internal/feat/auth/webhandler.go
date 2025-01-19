package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
)

var (
	todoResPath = "/todo"
	key         = am.Key
	method      = am.HTTPMethod
)

type WebHandler struct {
	core    *am.Handler
	service Service
	tm      *am.TemplateManager
}

func NewWebHandler(tm *am.TemplateManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		core:    handler,
		service: service,
		tm:      tm,
	}
}

func (h *WebHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of users")
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		http.Error(w, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(users)
	page.SetFormAction(todoResPath)
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("auth", "list-user")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) ShowUser(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Show user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)
	blue, _ := cfg.StrVal(key.ButtonStyleBlue)
	red, _ := cfg.StrVal(key.ButtonStyleRed)

	page := am.NewPage(user)
	page.SetActions([]am.Action{
		{URL: todoResPath, Text: "Back to User", Style: gray},
		{URL: fmt.Sprintf("%s/%s/edit", todoResPath, slug), Text: "Edit User", Style: blue},
		{URL: fmt.Sprintf("%s/%s/delete", todoResPath, slug), Text: "Delete User", Style: red},
	})

	tmpl, err := h.tm.Get("auth", "show-user")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
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
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Edit user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(user)
	page.SetFormAction(fmt.Sprintf("%s/%s/edit", todoResPath, slug))
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("auth", "edit-user")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	slug := r.FormValue("slug")
	h.Log().Info("Update user", slug)
	ctx := r.Context()

	user, err := h.service.GetUser(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	user.Username = name
	user.EncPassword = description

	err = h.service.UpdateUser(ctx, user)
	if err != nil {
		http.Error(w, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	slug := r.FormValue("slug")
	h.Log().Info("Delete user", slug)
	ctx := r.Context()

	err := h.service.DeleteUser(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
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
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, userSlug), http.StatusSeeOther)
}

func (h *WebHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Remove role from user")
	ctx := r.Context()

	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")

	err := h.service.RemoveRole(ctx, userSlug, roleSlug)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, userSlug), http.StatusSeeOther)
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
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) EditRole(w http.ResponseWriter, r *http.Request) {
	userSlug := chi.URLParam(r, "user_slug")
	roleSlug := chi.URLParam(r, "role_slug")
	h.Log().Info("Edit role", userSlug, roleSlug)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, userSlug, roleSlug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(role)
	page.SetFormAction(fmt.Sprintf("%s/%s/roles/%s/edit", todoResPath, userSlug, roleSlug))
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("auth", "edit-role")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")
	h.Log().Info("Update role", userSlug, roleSlug)
	ctx := r.Context()

	role, err := h.service.GetRole(ctx, userSlug, roleSlug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	role.Name = r.FormValue("name")
	role.Description = r.FormValue("description")

	err = h.service.UpdateRole(ctx, role)
	if err != nil {
		http.Error(w, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, userSlug), http.StatusSeeOther)
}

func (h *WebHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	userSlug := r.FormValue("user_slug")
	roleSlug := r.FormValue("role_slug")
	h.Log().Info("Delete role", userSlug, roleSlug)
	ctx := r.Context()

	err := h.service.DeleteRole(ctx, userSlug, roleSlug)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, userSlug), http.StatusSeeOther)
}

// Name returns the name in WebHandler.
func (h *WebHandler) Name() string {
	return h.core.Name()
}

// SetName sets the name in WebHandler.
func (h *WebHandler) SetName(name string) {
	h.core.SetName(name)
}

// Log returns the Logger in WebHandler.
func (h *WebHandler) Log() am.Logger {
	return h.core.Log()
}

// SetLog sets the Logger in WebHandler.
func (h *WebHandler) SetLog(log am.Logger) {
	h.core.SetLog(log)
}

// Cfg returns the Config in WebHandler.
func (h *WebHandler) Cfg() *am.Config {
	return h.core.Cfg()
}

// SetCfg sets the Config in WebHandler.
func (h *WebHandler) SetCfg(cfg *am.Config) {
	h.core.SetCfg(cfg)
}

// Setup is the default implementation for the Setup method in WebHandler.
func (h *WebHandler) Setup(ctx context.Context) error {
	return h.core.Setup(ctx)
}

// Start is the default implementation for the Start method in WebHandler.
func (h *WebHandler) Start(ctx context.Context) error {
	return h.core.Start(ctx)
}

// Stop is the default implementation for the Stop method in WebHandler.
func (h *WebHandler) Stop(ctx context.Context) error {
	return h.core.Stop(ctx)
}
