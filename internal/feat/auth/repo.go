package auth

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"

	"github.com/google/uuid"
)

type Repo interface {
	am.Repo

	// SECTION: User-related methods

	GetUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID, preload ...bool) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, user User) error
	GetUserAssignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]Role, error)
	GetUserUnassignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]Role, error)
	AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error
	RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error
	AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission Permission) error
	RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error
	GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (Role, error)
	GetUserAssignedPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	GetUserIndirectPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)

	// SECTION:  Role-related methods

	GetAllRoles(ctx context.Context) ([]Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID, preload ...bool) (Role, error)
	CreateRole(ctx context.Context, role Role) error
	UpdateRole(ctx context.Context, role Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	GetRoleUnassignedPermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission Permission) error
	RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error

	// SECTION: Permission-related methods

	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetPermission(ctx context.Context, id uuid.UUID) (Permission, error)
	CreatePermission(ctx context.Context, permission Permission) error
	UpdatePermission(ctx context.Context, permission Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error

	// SECTION: Resource-related methods

	GetAllResources(ctx context.Context) ([]Resource, error)
	GetResource(ctx context.Context, id uuid.UUID, preload ...bool) (Resource, error)
	CreateResource(ctx context.Context, resource Resource) error
	UpdateResource(ctx context.Context, resource Resource) error
	DeleteResource(ctx context.Context, id uuid.UUID) error
	GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error)
	GetResourceUnassignedPermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error)
	AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission Permission) error
	RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error

	// SECTION: Organization-related methods
	CreateOrg(ctx context.Context, org Org) error
	AddOrgOwner(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error
	GetDefaultOrg(ctx context.Context) (Org, error)
	GetOrgOwners(ctx context.Context, orgID uuid.UUID) ([]User, error)
	GetOrgUnassignedOwners(ctx context.Context, orgID uuid.UUID) ([]User, error)
	GetAllTeams(ctx context.Context, orgID uuid.UUID) ([]Team, error)
	GetTeam(ctx context.Context, id uuid.UUID) (Team, error)
	CreateTeam(ctx context.Context, team Team) error
	UpdateTeam(ctx context.Context, team Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
}
