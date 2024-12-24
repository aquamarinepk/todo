package todo

import (
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	base    *am.Handler
	service Service
}

func NewWebHandler(service Service, log am.Logger) *Handler {
	handler := am.NewHandler(log)
	return &Handler{
		base:    handler,
		service: service,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of lists")
	w.Write([]byte("List of lists"))
}

func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New todo form")
	w.Write([]byte("New todo form"))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create todo")
	w.Write([]byte("Create todo"))
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Show todo ", id)
	w.Write([]byte("Show todo " + id))
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Edit todo ", id)
	w.Write([]byte("Edit todo " + id))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Update todo ", id)
	w.Write([]byte("Update todo " + id))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Delete todo ", id)
	w.Write([]byte("Delete todo " + id))
}

func (h *Handler) Log() am.Logger {
	return h.base.Log()
}
