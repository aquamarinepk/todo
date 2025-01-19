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
	core.Get("/{slug}/edit", handler.EditUser)
	core.Post("/create-user", handler.CreateUser)
	core.Post("/{slug}/update", handler.UpdateUser)
	core.Post("/{slug}/delete", handler.DeleteUser)

	core.Get("/{user_slug}/roles/{role_slug}/edit", handler.EditRole)
	core.Post("/{user_slug}/roles/{role_slug}/update", handler.UpdateRole)
	core.Post("/{user_slug}/roles/{role_slug}/delete", handler.DeleteRole)
	core.Post("/create-role", handler.CreateRole)
	core.Post("/add-role", handler.AddRole)
	core.Post("/remove-role", handler.RemoveRole)

	return core
}
