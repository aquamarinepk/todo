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

func (h *APIHandler) ListLists(w http.ResponseWriter, r *http.Request) {
	lists, err := h.service.GetLists(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lists)
}

func (h *APIHandler) ShowList(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	list, err := h.service.GetListBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(list)
}

func (h *APIHandler) CreateList(w http.ResponseWriter, r *http.Request) {
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

func (h *APIHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list, err := h.service.GetListBySlug(r.Context(), payload.Slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	list.Name = payload.Name
	list.Description = payload.Description
	if err := h.service.UpdateList(r.Context(), list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Slug string `json:"slug"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteListBySlug(r.Context(), payload.Slug); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ListSlug    string `json:"list_slug"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item := NewItem(payload.Description, payload.Status)
	if err := h.service.AddItem(r.Context(), payload.ListSlug, item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ListSlug    string `json:"list_slug"`
		ItemID      string `json:"item_id"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.EditItem(r.Context(), payload.ListSlug, payload.ItemID, payload.Description, payload.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ListSlug string `json:"list_slug"`
		ItemID   string `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteItem(r.Context(), payload.ListSlug, payload.ItemID); err != nil {
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
