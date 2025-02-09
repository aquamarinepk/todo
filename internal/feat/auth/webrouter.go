package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewWebRouter creates a new web router for the todo feature.
func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	core.Get("/list-users", handler.ListUsers)
	core.Get("/new-user", handler.NewUser)
	core.Get("/show-user", handler.ShowUser)
	core.Get("/edit-user", handler.EditUser)
	core.Post("/create-user", handler.CreateUser)
	core.Post("/update-user", handler.UpdateUser)
	core.Post("/delete-user", handler.DeleteUser)
	core.Get("/list-user-roles", handler.ListUserRoles)
	core.Get("/add-role-to-user", handler.AddRoleToUser)

	core.Get("/edit-role", handler.EditRole)
	core.Post("/update-role", handler.UpdateRole)
	core.Post("/delete-role", handler.DeleteRole)
	core.Post("/create-role", handler.CreateRole)
	core.Post("/add-role", handler.AddRole)
	core.Post("/remove-role", handler.RemoveRole)

	return core
}
