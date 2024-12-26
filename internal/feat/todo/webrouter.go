package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

type WebRouter struct {
	core    *am.Router
	handler *WebHandler
}

func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	core.Get("/", handler.List)          // GET /todo
	core.Get("/new", handler.New)        // GET /todo/new
	core.Post("/", handler.Create)       // POST /todo
	core.Get("/{id}", handler.Show)      // GET /todo/{id}
	core.Get("/{id}/edit", handler.Edit) // GET /todo/{id}/edit
	core.Put("/{id}", handler.Update)    // PUT /todo/{id}
	core.Delete("/{id}", handler.Delete) // DELETE /todo/{id}

	return core
}
