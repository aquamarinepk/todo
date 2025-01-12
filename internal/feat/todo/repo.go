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
	GetListAll(ctx context.Context) ([]List, error)
	GetListByID(ctx context.Context, id uuid.UUID) (List, error)
	GetListBySlug(ctx context.Context, slug string) (List, error)
	CreateList(ctx context.Context, list List) error
	UpdateList(ctx context.Context, list List) error
	DeleteList(ctx context.Context, id uuid.UUID) error
	Debug()
}

type BaseRepo struct {
	core  *am.Repo
	mu    sync.Mutex
	lists map[uuid.UUID]ListDA
	order []uuid.UUID
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		core:  am.NewRepo("todo-repo", qm, opts...),
		lists: make(map[uuid.UUID]ListDA),
		order: []uuid.UUID{},
	}

	repo.addSampleData() // NOTE: Used for testing purposes only.

	return repo
}

func (repo *BaseRepo) addSampleData() {
	for i := 1; i <= 5; i++ {
		id := uuid.New()
		list := NewList(fmt.Sprintf("Sample ListLists %d", i), fmt.Sprintf("This is the description for sample list %d", i))
		list.GenSlug()
		list.SetCreateValues()
		listDA := toListDA(list)
		listDA.ID = id
		repo.lists[id] = listDA
		repo.order = append(repo.order, id)
		repo.Log().Info("Created list with ID: ", id)
	}
}

func (repo *BaseRepo) GetListAll(ctx context.Context) ([]List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result []List
	for _, id := range repo.order {
		result = append(result, toList(repo.lists[id]))
	}
	return result, nil
}

func (repo *BaseRepo) GetListByID(ctx context.Context, id uuid.UUID) (List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	listDA, exists := repo.lists[id]
	if !exists {
		return List{}, errors.New("list not found")
	}
	return toList(listDA), nil
}

func (repo *BaseRepo) GetListBySlug(ctx context.Context, slug string) (List, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, listDA := range repo.lists {
		if listDA.Slug.String == slug {
			return toList(listDA), nil
		}
	}
	return List{}, errors.New("list not found")
}

func (repo *BaseRepo) CreateList(ctx context.Context, list List) error {
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

func (repo *BaseRepo) UpdateList(ctx context.Context, list List) error {
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

func (repo *BaseRepo) DeleteList(ctx context.Context, id uuid.UUID) error {
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
	result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s %-50s\n", "Type", "ID", "NameID", "Slug", "Name", "Description")
	for _, id := range repo.order {
		listDA := repo.lists[id]
		result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s %-50s\n",
			listDA.Type, listDA.ID.String(), listDA.NameID.String, listDA.Slug.String, listDA.Name.String, listDA.Description.String)
	}
	result = fmt.Sprintf("%s state:\n%s", repo.Name(), result)
	repo.Log().Info(result)
}

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
