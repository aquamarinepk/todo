package todo

import (
	"context"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Service interface {
	GetLists(ctx context.Context) ([]List, error)
	Get(ctx context.Context, id uuid.UUID) (List, error)
	Create(ctx context.Context, list List) error
	Update(ctx context.Context, list List) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (svc *BaseService) Get(ctx context.Context, id uuid.UUID) (List, error) {
	return svc.repo.Get(ctx, id)
}

func (svc *BaseService) Create(ctx context.Context, list List) error {
	list.GenSlug()
	list.GenCreationValues()
	return svc.repo.Create(ctx, list)
}

func (svc *BaseService) Update(ctx context.Context, list List) error {
	return svc.repo.Update(ctx, list)
}

func (svc *BaseService) Delete(ctx context.Context, id uuid.UUID) error {
	return svc.repo.Delete(ctx, id)
}
