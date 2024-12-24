package todo

import (
	chi "github.com/go-chi/chi/v5"
)

// NewAPIRouter creates a new API router for the todo feature.
func NewAPIRouter(handler *APIHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.List)          // GET /api/todo
	r.Post("/", handler.Create)       // POST /api/todo
	r.Get("/{id}", handler.Show)      // GET /api/todo/{id}
	r.Put("/{id}", handler.Update)    // PUT /api/todo/{id}
	r.Delete("/{id}", handler.Delete) // DELETE /api/todo/{id}

	return r
}
