package todo

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (h *APIHandler) List(w http.ResponseWriter, r *http.Request) {
	lists, err := h.service.GetAllLists(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lists)
}

func (h *APIHandler) Create(w http.ResponseWriter, r *http.Request) {
	var list List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.CreateList(r.Context(), list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) Show(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	list, err := h.service.GetListByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *APIHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	var list List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list.SetID(id)
	if err := h.service.UpdateList(r.Context(), list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteList(r.Context(), id); err != nil {
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
