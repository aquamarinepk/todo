package auth

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUser(ctx context.Context, slug string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, slug string) error
	CreateRole(ctx context.Context, role Role) error
	GetUserRoles(ctx context.Context, userSlug string) ([]Role, error)
	GetRole(ctx context.Context, userSlug string, roleSlug string) (Role, error)
	UpdateRole(ctx context.Context, role Role) error
	DeleteRole(ctx context.Context, userSlug string, roleSlug string) error
	AddRole(ctx context.Context, userSlug string, roleSlug string) error
	RemoveRole(ctx context.Context, userSlug string, roleSlug string) error
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

func (svc *BaseService) GetUsers(ctx context.Context) ([]User, error) {
	return svc.repo.GetAllUsers(ctx)
}

func (svc *BaseService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return svc.repo.GetUserByID(ctx, id)
}

func (svc *BaseService) GetUser(ctx context.Context, slug string) (User, error) {
	return svc.repo.GetUserBySlug(ctx, slug)
}

func (svc *BaseService) CreateUser(ctx context.Context, user User) error {
	return svc.repo.CreateUser(ctx, user)
}

func (svc *BaseService) UpdateUser(ctx context.Context, user User) error {
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *BaseService) DeleteUser(ctx context.Context, slug string) error {
	return svc.repo.DeleteUser(ctx, slug)
}

func (svc *BaseService) GetUserRoles(ctx context.Context, userSlug string) ([]Role, error) {
	return svc.repo.GetRolesForUser(ctx, userSlug)
}

func (svc *BaseService) CreateRole(ctx context.Context, role Role) error {
	return svc.repo.CreateRole(ctx, role)
}

func (svc *BaseService) GetRole(ctx context.Context, userSlug string, roleSlug string) (Role, error) {
	_, err := svc.repo.GetUserBySlug(ctx, userSlug)
	if err != nil {
		return Role{}, err
	}
	return svc.repo.GetRoleBySlug(ctx, userSlug, roleSlug)
}

func (svc *BaseService) UpdateRole(ctx context.Context, role Role) error {
	_, err := svc.repo.GetUserBySlug(ctx, role.UserSlug)
	if err != nil {
		return err
	}
	return svc.repo.UpdateRole(ctx, role.UserSlug, role)
}

func (svc *BaseService) DeleteRole(ctx context.Context, userSlug string, roleSlug string) error {
	_, err := svc.repo.GetUserBySlug(ctx, userSlug)
	if err != nil {
		return err
	}
	return svc.repo.DeleteRole(ctx, userSlug, roleSlug)
}

func (svc *BaseService) AddRole(ctx context.Context, userSlug string, roleSlug string) error {
	role, err := svc.repo.GetRoleBySlug(ctx, userSlug, roleSlug)
	if err != nil {
		return err
	}
	return svc.repo.AddRole(ctx, userSlug, role)
}

func (svc *BaseService) RemoveRole(ctx context.Context, userSlug string, roleSlug string) error {
	return svc.repo.RemoveRole(ctx, userSlug, roleSlug)
}
