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

	core.Get("/", handler.ListLists)
	core.Get("/{slug}", handler.Show)
	core.Post("/create-list", handler.CreateList)
	core.Post("/edit-list", handler.EditList)
	core.Post("/delete-list", handler.DeleteList)
	core.Post("/add-item", handler.AddItem)
	core.Post("/edit-item", handler.EditItem)
	core.Post("/remove-item", handler.RemoveItem)

	return core
}
