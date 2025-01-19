package todo

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
)

type APIHandler struct {
	core    *am.Handler
	service Service
}

// NewAPIHandler creates a new API handler.
func NewAPIHandler(service Service, options ...am.Option) *APIHandler {
	handler := am.NewHandler("api-handler", options...)
	return &APIHandler{
		core:    handler,
		service: service,
	}
}

func (h *APIHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	var res am.Response
	if err != nil {
		res = am.NewErrorResponse("Failed to list users", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res = am.NewSuccessResponse("Users listed successfully", users)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) ShowUser(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	user, err := h.service.GetUser(r.Context(), slug)
	var res am.Response
	if err != nil {
		res = am.NewErrorResponse("User not found", am.ErrorCodeNotFound, err.Error())
		am.Respond(w, http.StatusNotFound, res)
		return
	}
	res = am.NewSuccessResponse("User retrieved successfully", user)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	user.SetCreateValues()
	if err := h.service.CreateUser(r.Context(), user); err != nil {
		res := am.NewErrorResponse("Failed to create user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("User created successfully", user)
	am.Respond(w, http.StatusCreated, res)
}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	user, err := h.service.GetUser(r.Context(), payload.Slug)
	if err != nil {
		res := am.NewErrorResponse("User not found", am.ErrorCodeNotFound, err.Error())
		am.Respond(w, http.StatusNotFound, res)
		return
	}
	user.Name = payload.Name
	user.EncPassword = payload.Description
	if err := h.service.UpdateUser(r.Context(), user); err != nil {
		res := am.NewErrorResponse("Failed to update user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("User updated successfully", user)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.DeleteUser(r.Context(), payload.Slug); err != nil {
		res := am.NewErrorResponse("Failed to delete user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("User deleted successfully", nil)
	am.Respond(w, http.StatusNoContent, res)
}

func (h *APIHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	role.SetCreateValues()
	if err := h.service.CreateRole(r.Context(), role); err != nil {
		res := am.NewErrorResponse("Failed to create role", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role created successfully", role)
	am.Respond(w, http.StatusCreated, res)
}

func (h *APIHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug    string `json:"user_slug"`
		RoleSlug    string `json:"role_slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	role, err := h.service.GetRole(r.Context(), payload.UserSlug, payload.RoleSlug)
	if err != nil {
		res := am.NewErrorResponse("Role not found", am.ErrorCodeNotFound, err.Error())
		am.Respond(w, http.StatusNotFound, res)
		return
	}
	role.Name = payload.Name
	role.Description = payload.Description
	if err := h.service.UpdateRole(r.Context(), role); err != nil {
		res := am.NewErrorResponse("Failed to update role", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role updated successfully", role)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug string `json:"user_slug"`
		RoleSlug string `json:"role_slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.DeleteRole(r.Context(), payload.UserSlug, payload.RoleSlug); err != nil {
		res := am.NewErrorResponse("Failed to delete role", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role deleted successfully", nil)
	am.Respond(w, http.StatusNoContent, res)
}

func (h *APIHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug string `json:"user_slug"`
		RoleSlug string `json:"role_slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.AddRole(r.Context(), payload.UserSlug, payload.RoleSlug); err != nil {
		res := am.NewErrorResponse("Failed to add role to user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role added to user successfully", nil)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug string `json:"user_slug"`
		RoleSlug string `json:"role_slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.RemoveRole(r.Context(), payload.UserSlug, payload.RoleSlug); err != nil {
		res := am.NewErrorResponse("Failed to remove role from user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role removed from user successfully", nil)
	am.Respond(w, http.StatusOK, res)
}

// Name returns the name in APIHandler.
func (h *APIHandler) Name() string {
	return h.core.Name()
}

// SetName sets the name in APIHandler.
func (h *APIHandler) SetName(name string) {
	h.core.SetName(name)
}

// Log returns the Logger in APIHandler.
func (h *APIHandler) Log() am.Logger {
	return h.core.Log()
}

// SetLog sets the Logger in APIHandler.
func (h *APIHandler) SetLog(log am.Logger) {
	h.core.SetLog(log)
}

// Cfg returns the Config in APIHandler.
func (h *APIHandler) Cfg() *am.Config {
	return h.core.Cfg()
}

// SetCfg sets the Config in APIHandler.
func (h *APIHandler) SetCfg(cfg *am.Config) {
	h.core.SetCfg(cfg)
}

// Setup is the default implementation for the Setup method in APIHandler.
func (h *APIHandler) Setup(ctx context.Context) error {
	return h.core.Setup(ctx)
}

// Start is the default implementation for the Start method in APIHandler.
func (h *APIHandler) Start(ctx context.Context) error {
	return h.core.Start(ctx)
}

// Stop is the default implementation for the Stop method in APIHandler.
func (h *APIHandler) Stop(ctx context.Context) error {
	return h.core.Stop(ctx)
}
