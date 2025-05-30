package todo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type Repo interface {
	GetAll(ctx context.Context) ([]List, error)
	Get(ctx context.Context, id uuid.UUID) (List, error)
	Create(ctx context.Context, list List) error
	Update(ctx context.Context, list List) error
	Delete(ctx context.Context, id uuid.UUID) error
	Debug()
}

type BaseRepo struct {
	*am.BaseRepo
	mu    sync.Mutex
	lists map[uuid.UUID]ListDA
	order []uuid.UUID
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		BaseRepo: am.NewRepo("todo-repo", qm, opts...),
		lists:    make(map[uuid.UUID]ListDA),
		order:    []uuid.UUID{},
	}

	return repo
}

func (repo *BaseRepo) GetAll(ctx context.Context) ([]List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result []List
	for _, id := range repo.order {
		result = append(result, toList(repo.lists[id]))
	}
	return result, nil
}

func (repo *BaseRepo) Get(ctx context.Context, id uuid.UUID) (List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA, exists := repo.lists[id]
	if !exists {
		return List{}, errors.New("list not found")
	}
	return toList(listDA), nil
}

func (repo *BaseRepo) Create(ctx context.Context, list List) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA := toListDA(list)
	if _, exists := repo.lists[listDA.ID]; exists {
		return errors.New("list already exists")
	}
	repo.lists[listDA.ID] = listDA
	repo.order = append(repo.order, listDA.ID)
	return nil
}

func (repo *BaseRepo) Update(ctx context.Context, list List) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA := toListDA(list)
	if _, exists := repo.lists[listDA.ID]; !exists {
		msg := fmt.Sprintf("list not found for ID: %s", listDA.ID)
		return errors.New(msg)
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
	for i, oid := range repo.order {
		if oid == id {
			repo.order = append(repo.order[:i], repo.order[i+1:]...)
			break
		}
	}
	return nil
}

func (repo *BaseRepo) Debug() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result string
	result += fmt.Sprintf("%-10s %-20s %-50s\n", "Type", "Action", "Description")
	for _, id := range repo.order {
		listDA := repo.lists[id]
		result += fmt.Sprintf("%-10s %-20s %-50s\n",
			listDA.Type, listDA.Name.String, listDA.Description.String)
	}
	result = fmt.Sprintf("%s state:\n%s", repo.Name(), result)
	repo.Log().Info(result)
}
