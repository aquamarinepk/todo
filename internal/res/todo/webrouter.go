package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

type WebRouter struct {
	handler *WebHandler
}

func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	r := am.NewRouter("web-router", opts...)

	r.Get("/", handler.List)
	r.Get("/new", handler.New)
	r.Post("/", handler.Create)
	r.Get("/{slug}", handler.Show)
	r.Get("/{slug}/edit", handler.Edit)
	r.Put("/{slug}", handler.Update)
	r.Delete("/{slug}", handler.Delete)

	return r
}
