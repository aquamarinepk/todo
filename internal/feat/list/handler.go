package list

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct{}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	h := &Handler{}

	r.Get("/", h.List)          // GET /list
	r.Get("/new", h.New)        // GET /list/new
	r.Post("/", h.Create)       // POST /list
	r.Get("/{id}", h.Show)      // GET /list/{id}
	r.Get("/{id}/edit", h.Edit) // GET /list/{id}/edit
	r.Put("/{id}", h.Update)    // PUT /list/{id}
	r.Delete("/{id}", h.Delete) // DELETE /list/{id}

	return r
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("List of lists"))
}

func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("New list form"))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create list"))
}

func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Show list " + id))
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Edit list " + id))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Update list " + id))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	w.Write([]byte("Delete list " + id))
}
