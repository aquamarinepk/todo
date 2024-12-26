package todo

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	GetAllLists(ctx context.Context) ([]List, error)
	GetListByID(ctx context.Context, id uuid.UUID) (List, error)
	CreateList(ctx context.Context, list List) error
	UpdateList(ctx context.Context, list List) error
	DeleteList(ctx context.Context, id uuid.UUID) error
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

func (svc *BaseService) GetAllLists(ctx context.Context) ([]List, error) {
	return svc.repo.GetAll(ctx)
}

func (svc *BaseService) GetListByID(ctx context.Context, id uuid.UUID) (List, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *BaseService) CreateList(ctx context.Context, list List) error {
	return svc.repo.Create(ctx, list)
}

func (svc *BaseService) UpdateList(ctx context.Context, list List) error {
	return svc.repo.Update(ctx, list)
}

func (svc *BaseService) DeleteList(ctx context.Context, id uuid.UUID) error {
	return svc.repo.Delete(ctx, id)
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
