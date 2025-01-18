package todo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type BaseRepo struct {
	core  *am.Repo
	mu    sync.Mutex
	users map[uuid.UUID]UserDA
	roles map[uuid.UUID]RoleDA
	order []uuid.UUID
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		core:  am.NewRepo("todo-repo", qm, opts...),
		users: make(map[uuid.UUID]UserDA),
		roles: make(map[uuid.UUID]RoleDA),
		order: []uuid.UUID{},
	}

	repo.addSampleData() // NOTE: Used for testing purposes only.

	return repo
}

func (repo *BaseRepo) addSampleData() {
	for i := 1; i <= 5; i++ {
		id := uuid.New()
		username := fmt.Sprintf("sampleuser%d", i)
		email := fmt.Sprintf("sampleuser%d@example.com", i)
		user := NewUser(username, email, username)
		user.GenSlug()
		user.SetCreateValues()
		userDA := toUserDA(user)
		userDA.ID = id
		repo.users[id] = userDA
		repo.order = append(repo.order, id)
		repo.Log().Info("Created user with ID: ", id)
	}
}

func (repo *BaseRepo) GetUserAll(ctx context.Context) ([]User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result []User
	for _, id := range repo.order {
		result = append(result, toUser(repo.users[id]))
	}
	return result, nil
}

func (repo *BaseRepo) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return toUser(userDA), nil
}

func (repo *BaseRepo) GetUserBySlug(ctx context.Context, slug string) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, userDA := range repo.users {
		if userDA.Slug.String == slug {
			return toUser(userDA), nil
		}
	}
	return User{}, errors.New("user not found")
}

func (repo *BaseRepo) CreateUser(ctx context.Context, user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA := toUserDA(user)
	if _, exists := repo.users[userDA.ID]; exists {
		return errors.New("user already exists")
	}
	repo.users[userDA.ID] = userDA
	repo.order = append(repo.order, userDA.ID)
	return nil
}

func (repo *BaseRepo) UpdateUser(ctx context.Context, user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA := toUserDA(user)
	if _, exists := repo.users[userDA.ID]; !exists {
		msg := fmt.Sprintf("user not found for ID: %s", userDA.ID)
		return errors.New(msg)
	}
	repo.users[userDA.ID] = userDA
	return nil
}

func (repo *BaseRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(repo.users, id)
	for i, oid := range repo.order {
		if oid == id {
			repo.order = append(repo.order[:i], repo.order[i+1:]...)
			break
		}
	}
	return nil
}

func (repo *BaseRepo) CreateRole(ctx context.Context, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	roleDA := toRoleDA(role)
	if _, exists := repo.roles[roleDA.ID]; exists {
		return errors.New("role already exists")
	}
	repo.roles[roleDA.ID] = roleDA
	return nil
}

func (repo *BaseRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	roleDA, exists := repo.roles[id]
	if !exists {
		return Role{}, errors.New("role not found")
	}
	return toRole(roleDA), nil
}

func (repo *BaseRepo) UpdateRole(ctx context.Context, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	roleDA := toRoleDA(role)
	if _, exists := repo.roles[roleDA.ID]; !exists {
		msg := fmt.Sprintf("role not found for ID: %s", roleDA.ID)
		return errors.New(msg)
	}
	repo.roles[roleDA.ID] = roleDA
	return nil
}

func (repo *BaseRepo) DeleteRole(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.roles[id]; !exists {
		return errors.New("role not found")
	}
	delete(repo.roles, id)
	return nil
}

func (repo *BaseRepo) AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA, userExists := repo.users[userID]
	if !userExists {
		return errors.New("user not found")
	}

	roleDA, roleExists := repo.roles[roleID]
	if !roleExists {
		return errors.New("role not found")
	}

	// TODO: This a simplified approach that maybe needs to be revisited later.
	userDA.Roles = append(userDA.Roles, roleDA)
	repo.users[userID] = userDA
	return nil
}

func (repo *BaseRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA, userExists := repo.users[userID]
	if !userExists {
		return errors.New("user not found")
	}

	// TODO: This a simplified approach that maybe needs to be revisited later.
	for i, role := range userDA.Roles {
		if role.ID == roleID {
			userDA.Roles = append(userDA.Roles[:i], userDA.Roles[i+1:]...)
			break
		}
	}
	repo.users[userID] = userDA
	return nil
}

func (repo *BaseRepo) Debug() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result string
	result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s %-50s\n", "Type", "ID", "NameID", "Slug", "Username", "EncPassword")
	for _, id := range repo.order {
		userDA := repo.users[id]
		result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s %-50s\n",
			userDA.Type, userDA.ID.String(), userDA.NameID.String, userDA.Slug.String, userDA.Name.String, userDA.Description.String)
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
