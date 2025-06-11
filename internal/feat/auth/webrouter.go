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
	core.Post("/create-user", handler.CreateUser)
	core.Get("/show-user", handler.ShowUser)
	core.Get("/edit-user", handler.EditUser)
	core.Post("/update-user", handler.UpdateUser)
	core.Post("/delete-user", handler.DeleteUser)
	// User relationships
	core.Get("/list-user-roles", handler.ListUserRoles)
	core.Get("/list-user-permissions", handler.ListUserPermissions)
	core.Post("/add-role-to-user", handler.AddRoleToUser)
	core.Post("/remove-role-from-user", handler.RemoveRoleFromUser)
	core.Post("/add-permission-to-user", handler.AddPermissionToUser)
	core.Post("/remove-permission-from-user", handler.RemovePermissionFromUser)
	core.Get("/list-user-contextual-roles", handler.ListUserContextualRoles)
	core.Post("/add-contextual-role", handler.AddContextualRole)
	core.Post("/remove-contextual-role", handler.RemoveContextualRole)

	// Role routes
	core.Get("/list-roles", handler.ListRoles)
	core.Get("/new-role", handler.NewRole)
	core.Post("/create-role", handler.CreateRole)
	core.Get("/show-role", handler.ShowRole)
	core.Get("/edit-role", handler.EditRole)
	core.Post("/update-role", handler.UpdateRole)
	core.Post("/delete-role", handler.DeleteRole)
	// Role relationships
	core.Get("/list-role-permissions", handler.ListRolePermissions)
	core.Post("/add-permission-to-role", handler.AddPermissionToRole)
	core.Post("/remove-permission-from-role", handler.RemovePermissionFromRole)

	// Permission routes
	core.Get("/list-permissions", handler.ListPermissions)
	core.Get("/new-permission", handler.NewPermission)
	core.Post("/create-permission", handler.CreatePermission)
	core.Get("/show-permission", handler.ShowPermission)
	core.Get("/edit-permission", handler.EditPermission)
	core.Post("/update-permission", handler.UpdatePermission)
	core.Post("/delete-permission", handler.DeletePermission)

	// Resource routes
	core.Get("/list-resources", handler.ListResources)
	core.Get("/new-resource", handler.NewResource)
	core.Post("/create-resource", handler.CreateResource)
	core.Get("/show-resource", handler.ShowResource)
	core.Get("/edit-resource", handler.EditResource)
	core.Post("/update-resource", handler.UpdateResource)
	core.Post("/delete-resource", handler.DeleteResource)
	// Resource relationships
	core.Get("/list-resource-permissions", handler.ListResourcePermissions)
	core.Post("/add-permission-to-resource", handler.AddPermissionToResource)
	core.Post("/remove-permission-from-resource", handler.RemovePermissionFromResource)

	// Organization routes
	core.Get("/show-org", handler.ShowOrg)

	// Team routes
	core.Get("/list-teams", handler.ListTeams)
	core.Get("/new-team", handler.NewTeam)
	core.Post("/create-team", handler.CreateTeam)
	core.Get("/show-team", handler.ShowTeam)
	core.Get("/edit-team", handler.EditTeam)
	core.Post("/update-team", handler.UpdateTeam)
	core.Post("/delete-team", handler.DeleteTeam)
	// Team relationships
	core.Get("/list-team-members", handler.ListTeamMembers)
	core.Post("/assign-user-to-team", handler.AssignUserToTeam)
	core.Post("/remove-user-from-team", handler.RemoveUserFromTeam)

	return core
}
