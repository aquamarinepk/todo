package todo

import (
	"context"
	"errors"
	"sync"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Repo interface {
	GetAll(ctx context.Context) ([]List, error)
	GetByID(ctx context.Context, id uuid.UUID) (List, error)
	GetBySlug(ctx context.Context, slug string) (List, error)
	Create(ctx context.Context, list List) error
	Update(ctx context.Context, list List) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BaseRepo struct {
	core  *am.BaseCore
	mu    sync.Mutex
	lists map[uuid.UUID]ListDA
}

func NewRepo(opts ...am.Option) *BaseRepo {
	return &BaseRepo{
		core:  am.NewCore("todo-repo", opts...),
		lists: make(map[uuid.UUID]ListDA),
	}
}

func (repo *BaseRepo) GetAll(ctx context.Context) ([]List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result []List
	for _, listDA := range repo.lists {
		result = append(result, toList(listDA))
	}
	return result, nil
}

func (repo *BaseRepo) GetByID(ctx context.Context, id uuid.UUID) (List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA, exists := repo.lists[id]
	if !exists {
		return List{}, errors.New("list not found")
	}
	return toList(listDA), nil
}

func (repo *BaseRepo) GetBySlug(ctx context.Context, slug string) (List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, listDA := range repo.lists {
		if listDA.Slug.String == slug {
			return toList(listDA), nil
		}
	}
	return List{}, errors.New("list not found")
}

func (repo *BaseRepo) Create(ctx context.Context, list List) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA := toListDA(list)
	if _, exists := repo.lists[listDA.ID]; exists {
		return errors.New("list already exists")
	}
	repo.lists[listDA.ID] = listDA
	return nil
}

func (repo *BaseRepo) Update(ctx context.Context, list List) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA := toListDA(list)
	if _, exists := repo.lists[listDA.ID]; !exists {
		return errors.New("list not found")
	}
	repo.lists[listDA.ID] = listDA
	return nil
}

func (repo *BaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.lists[id]; !exists {
		return errors.New("list not found")
	}
	delete(repo.lists, id)
	return nil
}

// Implementing the am.Core interface methods by delegating to the core field
func (repo *BaseRepo) Name() string {
	return repo.core.Name()
}

func (repo *BaseRepo) SetName(name string) {
	repo.core.SetName(name)
}

func (repo *BaseRepo) Log() am.Logger {
	return repo.core.Log()
}

func (repo *BaseRepo) SetLog(log am.Logger) {
	repo.core.SetLog(log)
}

func (repo *BaseRepo) Cfg() *am.Config {
	return repo.core.Cfg()
}

func (repo *BaseRepo) SetCfg(cfg *am.Config) {
	repo.core.SetCfg(cfg)
}

func (repo *BaseRepo) Setup(ctx context.Context) error {
	return repo.core.Setup(ctx)
}

func (repo *BaseRepo) Start(ctx context.Context) error {
	return repo.core.Start(ctx)
}

func (repo *BaseRepo) Stop(ctx context.Context) error {
	return repo.core.Stop(ctx)
}
