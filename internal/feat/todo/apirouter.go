package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewAPIRouter creates a new API router for the todo feature.
func NewAPIRouter(handler *APIHandler, opts ...am.Option) *am.Router {
	r := am.NewRouter("api-router", opts...)

	r.Get("/", handler.ListLists)
	r.Get("/{slug}", handler.ShowList)
	r.Post("/create-list", handler.CreateList)
	r.Post("/edit-list", handler.EditList)
	r.Post("/delete-list", handler.DeleteList)
	r.Post("/add-item", handler.AddItem)
	r.Post("/edit-item", handler.EditItem)
	r.Post("/remove-item", handler.RemoveItem)

	return r
}
