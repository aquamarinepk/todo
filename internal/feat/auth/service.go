package todo

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	GetLists(ctx context.Context) ([]List, error)
	GetListByID(ctx context.Context, id uuid.UUID) (List, error)
	GetListBySlug(ctx context.Context, slug string) (List, error)
	CreateList(ctx context.Context, list List) error
	UpdateList(ctx context.Context, list List) error
	DeleteList(ctx context.Context, id uuid.UUID) error
	DeleteListBySlug(ctx context.Context, slug string) error
	AddItem(ctx context.Context, listSlug string, item Item) error
	EditItem(ctx context.Context, listSlug, itemID, name, description string) error
	DeleteItem(ctx context.Context, listSlug, itemID string) error
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

func (svc *BaseService) GetLists(ctx context.Context) ([]List, error) {
	return svc.repo.GetListAll(ctx)
}

func (svc *BaseService) GetListByID(ctx context.Context, id uuid.UUID) (List, error) {
	return svc.repo.GetListByID(ctx, id)
}

func (svc *BaseService) GetListBySlug(ctx context.Context, slug string) (List, error) {
	return svc.repo.GetListBySlug(ctx, slug)
}

func (svc *BaseService) CreateList(ctx context.Context, list List) error {
	list.GenSlug()
	list.SetCreateValues()
	err := svc.repo.CreateList(ctx, list)
	svc.repo.Debug()
	return err
}

func (svc *BaseService) UpdateList(ctx context.Context, list List) error {
	return svc.repo.UpdateList(ctx, list)
}

func (svc *BaseService) DeleteList(ctx context.Context, slug string) error {
	return svc.repo.DeleteList(ctx, slug)
}

func (svc *BaseService) DeleteListBySlug(ctx context.Context, slug string) error {
	list, err := svc.repo.GetListBySlug(ctx, slug)
	if err != nil {
		return err
	}
	return svc.repo.DeleteList(ctx, list.Slug())
}

func (svc *BaseService) AddItem(ctx context.Context, listSlug string, item Item) error {
	list, err := svc.repo.GetListBySlug(ctx, listSlug)
	if err != nil {
		return err
	}
	item.SetCreateValues()
	return svc.repo.AddItem(ctx, list.Slug(), item)
}

func (svc *BaseService) EditItem(ctx context.Context, listSlug, itemID, name, description string) error {
	list, err := svc.repo.GetListBySlug(ctx, listSlug)
	if err != nil {
		return err
	}
	item, err := svc.repo.GetItemByID(ctx, list.ID(), itemID)
	if err != nil {
		return err
	}
	item.Description = description
	return svc.repo.UpdateItem(ctx, list.Slug(), item)
}

func (svc *BaseService) DeleteItem(ctx context.Context, listSlug, itemID string) error {
	list, err := svc.repo.GetListBySlug(ctx, listSlug)
	if err != nil {
		return err
	}
	return svc.repo.DeleteItem(ctx, list.Slug(), itemID)
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
