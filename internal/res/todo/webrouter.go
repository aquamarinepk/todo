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

	core.Get("/", handler.List)
	core.Get("/new", handler.New)
	core.Post("/", handler.Create)
	core.Get("/{slug}", handler.Show)
	core.Get("/{slug}/edit", handler.Edit)
	core.Put("/{slug}", handler.Update)
	core.Delete("/{slug}", handler.Delete)

	return core
}
