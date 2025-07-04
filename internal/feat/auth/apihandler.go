package auth

import (
	"encoding/json"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type APIHandler struct {
	*am.Handler
	service Service
}

// NewAPIHandler creates a new API handler.
func NewAPIHandler(service Service, options ...am.Option) *APIHandler {
	handler := am.NewHandler("api-handler", options...)
	return &APIHandler{
		Handler: handler,
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
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		res := am.NewErrorResponse("Invalid user ID", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	user, err := h.service.GetUser(r.Context(), id)
	var res am.Response
	if err != nil {
		res := am.NewErrorResponse("User not found", am.ErrorCodeNotFound, err.Error())
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
	user.GenCreateValues()
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
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	user, err := h.service.GetUser(r.Context(), payload.ID)
	if err != nil {
		res := am.NewErrorResponse("User not found", am.ErrorCodeNotFound, err.Error())
		am.Respond(w, http.StatusNotFound, res)
		return
	}
	user.Name = payload.Name
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
		ID uuid.UUID `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.DeleteUser(r.Context(), payload.ID); err != nil {
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
	role.GenCreateValues()
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
		UserID      uuid.UUID `json:"user_id"`
		RoleID      uuid.UUID `json:"role_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	role, err := h.service.GetRole(r.Context(), payload.RoleID)
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
		UserID uuid.UUID `json:"user_id"`
		RoleID uuid.UUID `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.DeleteRole(r.Context(), payload.RoleID); err != nil {
		res := am.NewErrorResponse("Failed to delete role", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role deleted successfully", nil)
	am.Respond(w, http.StatusNoContent, res)
}

func (h *APIHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID uuid.UUID `json:"user_id"`
		RoleID uuid.UUID `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.AddRole(r.Context(), payload.UserID, payload.RoleID); err != nil {
		res := am.NewErrorResponse("Failed to add role to user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role added to user successfully", nil)
	am.Respond(w, http.StatusOK, res)
}

func (h *APIHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID uuid.UUID `json:"user_id"`
		RoleID uuid.UUID `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		res := am.NewErrorResponse("Invalid request payload", am.ErrorCodeBadRequest, err.Error())
		am.Respond(w, http.StatusBadRequest, res)
		return
	}
	if err := h.service.RemoveRole(r.Context(), payload.UserID, payload.RoleID); err != nil {
		res := am.NewErrorResponse("Failed to remove role from user", am.ErrorCodeInternalError, err.Error())
		am.Respond(w, http.StatusInternalServerError, res)
		return
	}
	res := am.NewSuccessResponse("Role removed from user successfully", nil)
	am.Respond(w, http.StatusOK, res)
}
