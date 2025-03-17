package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewWebRouter creates a new web router for the todo feature.
func NewWebRouter(handler *WebHandler, opts ...am.Option) *am.Router {
	core := am.NewRouter("web-router", opts...)

	// User routes
	core.Get("/list-users", handler.ListUsers)
	core.Get("/new-user", handler.NewUser)
	core.Get("/show-user", handler.ShowUser)
	core.Get("/edit-user", handler.EditUser)
	core.Post("/create-user", handler.CreateUser)
	core.Post("/update-user", handler.UpdateUser)
	core.Post("/delete-user", handler.DeleteUser)
	core.Get("/list-user-roles", handler.ListUserRoles)
	core.Get("/add-role-to-user", handler.AddRoleToUser)

	// Role routes
	core.Get("/list-roles", handler.ListRoles)
	core.Get("/new-role", handler.NewRole)
	core.Get("/show-role", handler.ShowRole)
	core.Get("/edit-role", handler.EditRole)
	core.Post("/create-role", handler.CreateRole)
	core.Post("/update-role", handler.UpdateRole)
	core.Post("/delete-role", handler.DeleteRole)
	core.Post("/add-role", handler.AddRole)
	core.Post("/remove-role", handler.RemoveRole)
	core.Post("/add-role-to-user", handler.AddRoleToUser)
	core.Post("/remove-role-from-user", handler.RemoveRoleFromUser)

	// Permission routes
	core.Get("/list-permissions", handler.ListPermissions)
	core.Get("/new-permission", handler.NewPermission)
	core.Get("/show-permission", handler.ShowPermission)
	core.Get("/edit-permission", handler.EditPermission)
	core.Post("/create-permission", handler.CreatePermission)
	core.Post("/update-permission", handler.UpdatePermission)
	core.Post("/delete-permission", handler.DeletePermission)
	core.Post("/add-permission-to-role", handler.AddPermissionToRole)
	core.Post("/remove-permission-from-role", handler.RemovePermissionFromRole)
	core.Post("/add-permission-to-user", handler.AddPermissionToUser)
	core.Post("/remove-permission-from-user", handler.RemovePermissionFromUser)

	// Resource routes
	core.Get("/list-resources", handler.ListResources)
	core.Get("/new-resource", handler.NewResource)
	core.Get("/show-resource", handler.ShowResource)
	core.Get("/edit-resource", handler.EditResource)
	core.Post("/create-resource", handler.CreateResource)
	core.Post("/update-resource", handler.UpdateResource)
	core.Post("/delete-resource", handler.DeleteResource)
	core.Get("/list-resource-permissions", handler.ListResourcePermissions)
	core.Post("/add-permission-to-resource", handler.AddPermissionToResource)
	core.Post("/remove-permission-from-resource", handler.RemovePermissionFromResource)

	return core
}
