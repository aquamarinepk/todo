package todo

import (
	chi "github.com/go-chi/chi/v5"
)

// NewRouter creates a new router for the todo feature.
func NewRouter(handler *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.List)          // GET /todo
	r.Get("/new", handler.New)        // GET /todo/new
	r.Post("/", handler.Create)       // POST /todo
	r.Get("/{id}", handler.Show)      // GET /todo/{id}
	r.Get("/{id}/edit", handler.Edit) // GET /todo/{id}/edit
	r.Put("/{id}", handler.Update)    // PUT /todo/{id}
	r.Delete("/{id}", handler.Delete) // DELETE /todo/{id}

	return r
}
