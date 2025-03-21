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
	r.Get("/{id}", handler.Show)
	r.Get("/{id}/edit", handler.Edit)
	r.Put("/{id}", handler.Update)
	r.Delete("/{id}", handler.Delete)

	return r
}
