package todo

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetUserBySlug(ctx context.Context, slug string) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteUserBySlug(ctx context.Context, slug string) error
	AddRole(ctx context.Context, userSlug string, role Role) error
	EditRole(ctx context.Context, userSlug, roleID, name, description string) error
	DeleteRole(ctx context.Context, userSlug, roleID string) error
}

type BaseService struct {
	core am.Core
	repo Repo
}

func NewService(repo Repo, opts ...am.Option) *BaseService {
	return &BaseService{
		core: am.NewCore("", opts...),
		repo: repo,
	}
}

func (svc *BaseService) GetUsers(ctx context.Context) ([]User, error) {
	return svc.repo.GetUserAll(ctx)
}

func (svc *BaseService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return svc.repo.GetUserByID(ctx, id)
}

func (svc *BaseService) GetUserBySlug(ctx context.Context, slug string) (User, error) {
	return svc.repo.GetUserBySlug(ctx, slug)
}

func (svc *BaseService) CreateUser(ctx context.Context, user User) error {
	user.GenSlug()
	user.SetCreateValues()
	err := svc.repo.CreateUser(ctx, user)
	svc.repo.Debug()
	return err
}

func (svc *BaseService) UpdateUser(ctx context.Context, user User) error {
	return svc.repo.UpdateUser(ctx, user)
}

func (svc *BaseService) DeleteUser(ctx context.Context, slug string) error {
	return svc.repo.DeleteUser(ctx, slug)
}

func (svc *BaseService) DeleteUserBySlug(ctx context.Context, slug string) error {
	user, err := svc.repo.GetUserBySlug(ctx, slug)
	if err != nil {
		return err
	}
	return svc.repo.DeleteUser(ctx, user.Slug())
}

func (svc *BaseService) AddRole(ctx context.Context, userSlug string, role Role) error {
	user, err := svc.repo.GetUserBySlug(ctx, userSlug)
	if err != nil {
		return err
	}
	role.SetCreateValues()
	return svc.repo.AddRole(ctx, user.Slug(), role)
}

func (svc *BaseService) EditRole(ctx context.Context, userSlug, roleID, name, description string) error {
	user, err := svc.repo.GetUserBySlug(ctx, userSlug)
	if err != nil {
		return err
	}
	role, err := svc.repo.GetRoleByID(ctx, user.ID(), roleID)
	if err != nil {
		return err
	}
	role.Description = description
	return svc.repo.UpdateRole(ctx, user.Slug(), role)
}

func (svc *BaseService) DeleteRole(ctx context.Context, userSlug, roleID string) error {
	user, err := svc.repo.GetUserBySlug(ctx, userSlug)
	if err != nil {
		return err
	}
	return svc.repo.DeleteRole(ctx, user.Slug(), roleID)
}

func (svc *BaseService) Name() string {
	return svc.core.Name()
}

func (svc *BaseService) SetName(name string) {
	svc.core.SetName(name)
}

func (svc *BaseService) Log() am.Logger {
	return svc.core.Log()
}

func (svc *BaseService) SetLog(log am.Logger) {
	svc.core.SetLog(log)
}

func (svc *BaseService) Cfg() *am.Config {
	return svc.core.Cfg()
}

func (svc *BaseService) SetCfg(cfg *am.Config) {
	svc.core.SetCfg(cfg)
}

func (svc *BaseService) Setup(ctx context.Context) error {
	return svc.core.Setup(ctx)
}

func (svc *BaseService) Start(ctx context.Context) error {
	return svc.core.Start(ctx)
}

func (svc *BaseService) Stop(ctx context.Context) error {
	return svc.core.Stop(ctx)
}
