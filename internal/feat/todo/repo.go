package todo

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type Repo interface {
	GetAll(ctx context.Context) ([]List, error)
	GetByID(ctx context.Context, id uuid.UUID) (List, error)
	Create(ctx context.Context, list List) error
	Update(ctx context.Context, list List) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type BaseRepo struct {
	mu    sync.Mutex
	lists map[uuid.UUID]List
}

func NewBaseRepo() *BaseRepo {
	return &BaseRepo{
		lists: make(map[uuid.UUID]List),
	}
}

func (r *BaseRepo) GetAll(ctx context.Context) ([]List, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []List
	for _, list := range r.lists {
		result = append(result, list)
	}
	return result, nil
}

func (r *BaseRepo) GetByID(ctx context.Context, id uuid.UUID) (List, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	list, exists := r.lists[id]
	if !exists {
		return List{}, errors.New("list not found")
	}
	return list, nil
}

func (r *BaseRepo) Create(ctx context.Context, list List) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.lists[list.ID()]; exists {
		return errors.New("list already exists")
	}
	r.lists[list.ID()] = list
	return nil
}

func (r *BaseRepo) Update(ctx context.Context, list List) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.lists[list.ID()]; !exists {
		return errors.New("list not found")
	}
	r.lists[list.ID()] = list
	return nil
}

func (r *BaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.lists[id]; !exists {
		return errors.New("list not found")
	}
	delete(r.lists, id)
	return nil
}
