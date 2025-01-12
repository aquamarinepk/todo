package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

type WebRouter struct {
	core    *am.Router
	handler *WebHandler
}

// NewWebRouter creates a new web router for the todo feature.
// GET requests will be mounted to the app's web router that handles `/query` requests,
// and POST requests will be mounted to the app's router that handles `/cmd` requests.
func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	core.Get("/", handler.ListLists)
	core.Get("/{slug}", handler.ShowList)
	core.Post("/create-list", handler.CreateList)
	core.Post("/edit-list", handler.EditList)
	core.Post("/delete-list", handler.DeleteList)
	core.Post("/add-item", handler.AddItem)
	core.Post("/edit-item", handler.EditItem)
	core.Post("/delete-item", handler.DeleteItem)

	return core
}
