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
	DeleteList(ctx context.Context, slug string) error
	AddItem(ctx context.Context, listSlug string, item Item) error
	GetItemByID(ctx context.Context, listID uuid.UUID, itemID string) (Item, error)
	GetItemBySlug(ctx context.Context, listSlug, itemSlug string) (Item, error)
	UpdateItem(ctx context.Context, listSlug string, item Item) error
	DeleteItem(ctx context.Context, listSlug, itemSlug string) error
	Debug()
}

type BaseRepo struct {
	core  *am.Repo
	mu    sync.Mutex
	lists map[uuid.UUID]ListDA
	items map[string]ItemDA
	order []uuid.UUID
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		core:  am.NewRepo("todo-repo", qm, opts...),
		lists: make(map[uuid.UUID]ListDA),
		items: make(map[string]ItemDA),
		order: []uuid.UUID{},
	}

	repo.addSampleData() // NOTE: Used for testing purposes only.

	return repo
}

func (repo *BaseRepo) addSampleData() {
	for i := 1; i <= 5; i++ {
		id := uuid.New()
		list := NewList(fmt.Sprintf("Sample List %d", i), fmt.Sprintf("This is the description for sample list %d", i))
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

func (repo *BaseRepo) DeleteList(ctx context.Context, slug string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var id uuid.UUID
	for _, listDA := range repo.lists {
		if listDA.Slug.String == slug {
			id = listDA.ID
			break
		}
	}
	if id == uuid.Nil {
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

func (repo *BaseRepo) AddItem(ctx context.Context, listSlug string, item Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var listID uuid.UUID
	for _, listDA := range repo.lists {
		if listDA.Slug.String == listSlug {
			listID = listDA.ID
			break
		}
	}
	if listID == uuid.Nil {
		return errors.New("list not found")
	}

	itemDA := toItemDA(item)
	if _, exists := repo.items[itemDA.ID]; exists {
		return errors.New("item already exists")
	}
	repo.items[itemDA.ID] = itemDA
	return nil
}

func (repo *BaseRepo) GetItemByID(ctx context.Context, listID uuid.UUID, itemID string) (Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	itemDA, exists := repo.items[itemID]
	if !exists {
		return Item{}, errors.New("item not found")
	}
	return toItem(itemDA), nil
}

func (repo *BaseRepo) GetItemBySlug(ctx context.Context, listSlug, itemSlug string) (Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var listID uuid.UUID
	for _, listDA := range repo.lists {
		if listDA.Slug.String == listSlug {
			listID = listDA.ID
			break
		}
	}
	if listID == uuid.Nil {
		return Item{}, errors.New("list not found")
	}

	for _, itemDA := range repo.items {
		if itemDA.ListID == listID && itemDA.Description.String == itemSlug {
			return toItem(itemDA), nil
		}
	}
	return Item{}, errors.New("item not found")
}

func (repo *BaseRepo) UpdateItem(ctx context.Context, listSlug string, item Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var listID uuid.UUID
	for _, listDA := range repo.lists {
		if listDA.Slug.String == listSlug {
			listID = listDA.ID
			break
		}
	}
	if listID == uuid.Nil {
		return errors.New("list not found")
	}

	itemDA := toItemDA(item)
	if _, exists := repo.items[itemDA.ID]; !exists {
		msg := fmt.Sprintf("item not found for ID: %s", itemDA.ID)
		return errors.New(msg)
	}
	repo.items[itemDA.ID] = itemDA
	return nil
}

func (repo *BaseRepo) DeleteItem(ctx context.Context, listSlug, itemSlug string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var listID uuid.UUID
	for _, listDA := range repo.lists {
		if listDA.Slug.String == listSlug {
			listID = listDA.ID
			break
		}
	}
	if listID == uuid.Nil {
		return errors.New("list not found")
	}

	var itemID string
	for id, itemDA := range repo.items {
		if itemDA.ListID == listID && itemDA.Description.String == itemSlug {
			itemID = id
			break
		}
	}
	if itemID == "" {
		return errors.New("item not found")
	}
	delete(repo.items, itemID)
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
