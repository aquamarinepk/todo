package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewAPIRouter creates a new API router for the todo feature.
// GET requests will be mounted to the app's API router that handles `/query` requests,
// and POST requests will be mounted to the app's router that handles `/cmd` requests.
func NewAPIRouter(handler *APIHandler, opts ...am.Option) *am.Router {
	r := am.NewRouter("api-router", opts...)

	r.Get("/", handler.ListLists)
	r.Get("/{slug}", handler.ShowList)
	r.Post("/create-list", handler.CreateList)
	r.Post("/update-list", handler.UpdateList)
	r.Post("/delete-list", handler.DeleteList)
	r.Post("/add-item", handler.AddItem)
	r.Post("/update-item", handler.UpdateItem)
	r.Post("/delete-item", handler.DeleteItem)

	return r
}
