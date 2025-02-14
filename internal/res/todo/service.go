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

func (svc *BaseService) GetLists(ctx context.Context) ([]List, error) {
	return svc.repo.GetAll(ctx)
}

func (svc *BaseService) GetListByID(ctx context.Context, id uuid.UUID) (List, error) {
	return svc.repo.GetByID(ctx, id)
}

func (svc *BaseService) GetListBySlug(ctx context.Context, slug string) (List, error) {
	return svc.repo.GetBySlug(ctx, slug)
}

func (svc *BaseService) CreateList(ctx context.Context, list List) error {
	list.GenSlug("") // NOTE: This function will not accept any arguments later.
	list.GenCreationValues()
	err := svc.repo.Create(ctx, list) // NOTE: Remove assignment to err and return the function call directly.
	svc.repo.Debug()
	return err
}

func (svc *BaseService) UpdateList(ctx context.Context, list List) error {
	return svc.repo.Update(ctx, list)
}

func (svc *BaseService) DeleteList(ctx context.Context, id uuid.UUID) error {
	return svc.repo.Delete(ctx, id)
}

func (svc *BaseService) DeleteListBySlug(ctx context.Context, slug string) error {
	list, err := svc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	return svc.repo.Delete(ctx, list.ID())
}
