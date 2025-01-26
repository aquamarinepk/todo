package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewAPIRouter creates a new API router for the todo resource.
func NewAPIRouter(handler *APIHandler, opts ...am.Option) *am.Router {
	r := am.NewRouter("api-router", opts...)

	r.Get("/", handler.List)          // GET /api/todo
	r.Post("/", handler.Create)       // POST /api/todo
	r.Get("/{id}", handler.Show)      // GET /api/todo/{id}
	r.Put("/{id}", handler.Update)    // PUT /api/todo/{id}
	r.Delete("/{id}", handler.Delete) // DELETE /api/todo/{id}

	return r
}
