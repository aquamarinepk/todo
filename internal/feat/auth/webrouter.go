package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

type WebRouter struct {
	core    *am.Router
	handler *WebHandler
}

// NewWebRouter creates a new web router for the todo feature.
func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	core.Get("/", handler.ListUsers)
	core.Get("/{slug}", handler.ShowUser)
	core.Post("/create-user", handler.CreateUser)
	core.Post("/edit-user", handler.EditUser)
	core.Post("/delete-user", handler.DeleteUser)
	core.Post("/create-role", handler.CreateRole)
	core.Post("/edit-role", handler.EditRole)
	core.Post("/delete-role", handler.DeleteRole)
	core.Post("/add-role", handler.AddRoleToUser)
	core.Post("/remove-role", handler.RemoveRoleFromUser)

	return core
}
