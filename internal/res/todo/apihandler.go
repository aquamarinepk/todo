package todo

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

func NewAPIHandler(service Service, options ...am.Option) *APIHandler {
	handler := am.NewHandler("api-handler", options...)
	return &APIHandler{
		Handler: handler,
		service: service,
	}
}

func (h *APIHandler) List(w http.ResponseWriter, r *http.Request) {
	lists, err := h.service.GetLists(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(lists) // TODO: handle error
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
