package auth

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	// User-related methods
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUser(ctx context.Context, id uuid.UUID) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	GetUserUnassignedRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)
	GetAllUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error)
	AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error
	AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission Permission) error
	RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error

	// Role-related methods
	GetAllRoles(ctx context.Context) ([]Role, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (Role, error)
	CreateRole(ctx context.Context, role Role) error
	UpdateRole(ctx context.Context, role Role) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error)
	AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error
	RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error

	// Permission-related methods
	GetAllPermissions(ctx context.Context) ([]Permission, error)
	GetPermission(ctx context.Context, id uuid.UUID) (Permission, error)
	CreatePermission(ctx context.Context, permission Permission) error
	UpdatePermission(ctx context.Context, permission Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error

	// Resource-related methods
	GetAllResources(ctx context.Context) ([]Resource, error)
	GetResource(ctx context.Context, id uuid.UUID) (Resource, error)
	CreateResource(ctx context.Context, resource Resource) error
	UpdateResource(ctx context.Context, resource Resource) error
	DeleteResource(ctx context.Context, id uuid.UUID) error
	GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error)
	AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission Permission) error
	RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error
}

type BaseService struct {
	*am.Service
	repo Repo
}

func NewService(repo Repo, opts ...am.Option) *BaseService {
	return &BaseService{
		Service: am.NewService("", opts...),
		repo:    repo,
	}
}

func (svc *BaseService) GetAllUsers(ctx context.Context) ([]User, error) {
	return svc.repo.GetAllUsers(ctx)
}

func (svc *BaseService) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	return svc.repo.GetUser(ctx, id)
}

func (svc *BaseService) CreateUser(ctx context.Context, user User) error {
	return svc.repo.CreateUser(ctx, user)
}

func (svc *BaseService) UpdateUser(ctx context.Context, user User) error {
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *BaseService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteUser(ctx, id)
}

func (svc *BaseService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	return svc.repo.GetUserRoles(ctx, userID)
}

func (svc *BaseService) GetUserUnassignedRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	return svc.repo.GetUserUnassignedRoles(ctx, userID)
}

func (svc *BaseService) CreateRole(ctx context.Context, role Role) error {
	return svc.repo.CreateRole(ctx, role)
}

func (svc *BaseService) GetRole(ctx context.Context, roleID uuid.UUID) (Role, error) {
	return svc.repo.GetRole(ctx, roleID)
}

func (svc *BaseService) UpdateRole(ctx context.Context, role Role) error {
	return svc.repo.UpdateRole(ctx, role)
}

func (svc *BaseService) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return svc.repo.DeleteRole(ctx, roleID)
}

func (svc *BaseService) GetAllRoles(ctx context.Context) ([]Role, error) {
	return svc.repo.GetAllRoles(ctx)
}

func (svc *BaseService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]Permission, error) {
	return svc.repo.GetRolePermissions(ctx, roleID)
}

func (svc *BaseService) AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return svc.repo.AddRole(ctx, userID, roleID)
}

func (svc *BaseService) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	return svc.repo.RemoveRole(ctx, userID, roleID)
}

func (svc *BaseService) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	return svc.repo.GetAllPermissions(ctx)
}

func (svc *BaseService) CreatePermission(ctx context.Context, permission Permission) error {
	return svc.repo.CreatePermission(ctx, permission)
}

func (svc *BaseService) GetPermission(ctx context.Context, id uuid.UUID) (Permission, error) {
	return svc.repo.GetPermission(ctx, id)
}

func (svc *BaseService) UpdatePermission(ctx context.Context, permission Permission) error {
	return svc.repo.UpdatePermission(ctx, permission)
}

func (svc *BaseService) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeletePermission(ctx, id)
}

func (svc *BaseService) GetAllUserPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	return svc.repo.GetAllUserPermissions(ctx, userID)
}

func (svc *BaseService) GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	return svc.repo.GetUserDirectPermissions(ctx, userID)
}

func (svc *BaseService) GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]Permission, error) {
	return svc.repo.GetUserUnassignedPermissions(ctx, userID)
}

func (svc *BaseService) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission Permission) error {
	return svc.repo.AddPermissionToUser(ctx, userID, permission)
}

func (svc *BaseService) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	return svc.repo.RemovePermissionFromUser(ctx, userID, permissionID)
}

func (svc *BaseService) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	permission, err := svc.GetPermission(ctx, permissionID)
	if err != nil {
		return err
	}
	return svc.repo.AddPermissionToRole(ctx, roleID, permission)
}

func (svc *BaseService) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	return svc.repo.RemovePermissionFromRole(ctx, roleID, permissionID)
}

func (svc *BaseService) GetAllResources(ctx context.Context) ([]Resource, error) {
	return svc.repo.GetAllResources(ctx)
}

func (svc *BaseService) GetResource(ctx context.Context, id uuid.UUID) (Resource, error) {
	return svc.repo.GetResource(ctx, id)
}

func (svc *BaseService) CreateResource(ctx context.Context, resource Resource) error {
	return svc.repo.CreateResource(ctx, resource)
}

func (svc *BaseService) UpdateResource(ctx context.Context, resource Resource) error {
	return svc.repo.UpdateResource(ctx, resource)
}

func (svc *BaseService) DeleteResource(ctx context.Context, id uuid.UUID) error {
	return svc.repo.DeleteResource(ctx, id)
}

func (svc *BaseService) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error) {
	return svc.repo.GetResourcePermissions(ctx, resourceID)
}

func (svc *BaseService) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission Permission) error {
	return svc.repo.AddPermissionToResource(ctx, resourceID, permission)
}

func (svc *BaseService) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
	return svc.repo.RemovePermissionFromResource(ctx, resourceID, permissionID)
}
