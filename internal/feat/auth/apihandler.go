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

func NewAPIHandler(service Service, options ...am.Option) *APIHandler {
	handler := am.NewHandler("api-handler", options...)
	return &APIHandler{
		core:    handler,
		service: service,
	}
}

func (h *APIHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *APIHandler) ShowUser(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	user, err := h.service.GetUserBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.service.GetUserBySlug(r.Context(), payload.Slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user.Username = payload.Name
	user.EncPassword = payload.Description
	if err := h.service.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteUserBySlug(r.Context(), payload.Slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug    string `json:"user_slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	role := NewRole(payload.Name, payload.Description, payload.Status)
	if err := h.service.AddRole(r.Context(), payload.UserSlug, role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug    string `json:"user_slug"`
		RoleID      string `json:"role_id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.EditRole(r.Context(), payload.UserSlug, payload.RoleID, payload.Description, payload.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserSlug string `json:"user_slug"`
		RoleID   string `json:"role_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteRole(r.Context(), payload.UserSlug, payload.RoleID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
