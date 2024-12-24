package todo

import (
	"context"

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
	repo Repo
}

func NewService(repo Repo) *BaseService {
	return &BaseService{
		repo: repo,
	}
}

func (s *BaseService) GetAllLists(ctx context.Context) ([]List, error) {
	return s.repo.GetAll(ctx)
}

func (s *BaseService) GetListByID(ctx context.Context, id uuid.UUID) (List, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BaseService) CreateList(ctx context.Context, list List) error {
	return s.repo.Create(ctx, list)
}

func (s *BaseService) UpdateList(ctx context.Context, list List) error {
	return s.repo.Update(ctx, list)
}

func (s *BaseService) DeleteList(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
