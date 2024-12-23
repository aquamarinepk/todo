package todo

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List of lists"))
}

func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("New todo form"))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create todo"))
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Show todo " + id))
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Edit todo " + id))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Update todo " + id))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Delete todo " + id))
}
