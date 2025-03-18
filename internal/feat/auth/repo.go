package auth

import (
	"context"

	"github.com/google/uuid"
)

type Repo interface {
	// User methods
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserBySlug(ctx context.Context, slug string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, slug string) error
	GetRolesForUser(ctx context.Context, userSlug string) ([]Role, error)
	AddRole(ctx context.Context, userSlug string, role Role) error
	RemoveRole(ctx context.Context, userSlug string, roleID string) error
	AddPermissionToUser(ctx context.Context, userSlug string, permission Permission) error
	RemovePermissionFromUser(ctx context.Context, userSlug string, permissionID string) error

	GetRoleByID(ctx context.Context, userID uuid.UUID, roleID string) (Role, error)
	GetRoleBySlug(ctx context.Context, userSlug, roleSlug string) (Role, error)
	CreateRole(ctx context.Context, role Role) error
	UpdateRole(ctx context.Context, userSlug string, role Role) error
	DeleteRole(ctx context.Context, userSlug, roleSlug string) error
	AddPermissionToRole(ctx context.Context, roleSlug string, permission Permission) error
	RemovePermissionFromRole(ctx context.Context, roleSlug string, permissionID string) error

	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetPermissionByID(ctx context.Context, id string) (Permission, error)
	CreatePermission(ctx context.Context, permission Permission) error
	UpdatePermission(ctx context.Context, permission Permission) error
	DeletePermission(ctx context.Context, id string) error

	GetAllResources(ctx context.Context) ([]Resource, error)
	GetResourceByID(ctx context.Context, id string) (Resource, error)
	CreateResource(ctx context.Context, resource Resource) error
	UpdateResource(ctx context.Context, resource Resource) error
	DeleteResource(ctx context.Context, id string) error
	AddPermissionToResource(ctx context.Context, resourceID string, permission Permission) error
	RemovePermissionFromResource(ctx context.Context, resourceID string, permissionID string) error

	Debug()
}
