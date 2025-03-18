package auth

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type BaseRepo struct {
	*am.Repo
	mu          sync.Mutex
	users       map[uuid.UUID]UserDA
	roles       map[string]RoleDA
	permissions map[string]PermissionDA
	resources   map[string]ResourceDA
	userRoles   map[uuid.UUID][]string
	order       []uuid.UUID
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		Repo:        am.NewRepo("todo-repo", qm, opts...),
		users:       make(map[uuid.UUID]UserDA),
		roles:       make(map[string]RoleDA),
		permissions: make(map[string]PermissionDA),
		resources:   make(map[string]ResourceDA),
		userRoles:   make(map[uuid.UUID][]string),
		order:       []uuid.UUID{},
	}

	repo.addSampleData() // NOTE: Used for testing purposes only.

	return repo
}

// User methods

func (repo *BaseRepo) GetAllUsers(ctx context.Context) ([]User, error) {
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

func (repo *BaseRepo) DeleteUser(ctx context.Context, slug string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var id uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == slug {
			id = userDA.ID
			break
		}
	}
	if id == uuid.Nil {
		return errors.New("user not found")
	}
	delete(repo.users, id)
	for i, oid := range repo.order {
		if oid == id {
			repo.order = append(repo.order[:i], repo.order[i+1:]...)
			break
		}
	}
	delete(repo.userRoles, id) // Remove user roles
	return nil
}

func (repo *BaseRepo) GetRolesForUser(ctx context.Context, userSlug string) ([]Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return nil, errors.New("user not found")
	}

	var roles []Role
	for _, roleID := range repo.userRoles[userID] {
		roleDA := repo.roles[roleID]
		roles = append(roles, toRole(roleDA))
	}
	return roles, nil
}

func (repo *BaseRepo) AddRole(ctx context.Context, userSlug string, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return errors.New("user not found")
	}

	roleDA := toRoleDA(role)
	if _, exists := repo.roles[roleDA.ID]; exists {
		return errors.New("role already exists")
	}
	repo.roles[roleDA.ID] = roleDA
	repo.userRoles[userID] = append(repo.userRoles[userID], roleDA.ID) // Add role to user
	return nil
}

func (repo *BaseRepo) RemoveRole(ctx context.Context, userSlug string, roleID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return errors.New("user not found")
	}

	// Remove role from userRoles
	for i, rid := range repo.userRoles[userID] {
		if rid == roleID {
			repo.userRoles[userID] = append(repo.userRoles[userID][:i], repo.userRoles[userID][i+1:]...)
			break
		}
	}
	return nil
}

func (repo *BaseRepo) AddPermissionToUser(ctx context.Context, userSlug string, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userDA.Permissions = append(userDA.Permissions, permission.ID().String())
			repo.users[userDA.ID] = userDA
			return nil
		}
	}
	return errors.New("user not found")
}

func (repo *BaseRepo) RemovePermissionFromUser(ctx context.Context, userSlug string, permissionID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			for i, pid := range userDA.Permissions {
				if pid == permissionID {
					userDA.Permissions = append(userDA.Permissions[:i], userDA.Permissions[i+1:]...)
					repo.users[userDA.ID] = userDA
					return nil
				}
			}
		}
	}
	return errors.New("user or permission not found")
}

// Role methods

func (repo *BaseRepo) GetRoleByID(ctx context.Context, userID uuid.UUID, roleID string) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	roleDA, exists := repo.roles[roleID]
	if !exists {
		return Role{}, errors.New("role not found")
	}
	return toRole(roleDA), nil
}

func (repo *BaseRepo) GetRoleBySlug(ctx context.Context, userSlug, roleSlug string) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return Role{}, errors.New("user not found")
	}

	for _, roleID := range repo.userRoles[userID] {
		roleDA := repo.roles[roleID]
		if roleDA.Description.String == roleSlug {
			return toRole(roleDA), nil
		}
	}
	return Role{}, errors.New("role not found")
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

func (repo *BaseRepo) UpdateRole(ctx context.Context, userSlug string, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return errors.New("user not found")
	}

	roleDA := toRoleDA(role)
	if _, exists := repo.roles[roleDA.ID]; !exists {
		msg := fmt.Sprintf("role not found for ID: %s", roleDA.ID)
		return errors.New(msg)
	}
	repo.roles[roleDA.ID] = roleDA
	return nil
}

func (repo *BaseRepo) DeleteRole(ctx context.Context, userSlug, roleSlug string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var userID uuid.UUID
	for _, userDA := range repo.users {
		if userDA.Slug.String == userSlug {
			userID = userDA.ID
			break
		}
	}
	if userID == uuid.Nil {
		return errors.New("user not found")
	}

	var roleID string
	for _, rid := range repo.userRoles[userID] {
		roleDA := repo.roles[rid]
		if roleDA.Description.String == roleSlug {
			roleID = rid
			break
		}
	}
	if roleID == "" {
		return errors.New("role not found")
	}
	delete(repo.roles, roleID)
	// Remove role from userRoles
	for i, rid := range repo.userRoles[userID] {
		if rid == roleID {
			repo.userRoles[userID] = append(repo.userRoles[userID][:i], repo.userRoles[userID][i+1:]...)
			break
		}
	}
	return nil
}

func (repo *BaseRepo) AddPermissionToRole(ctx context.Context, roleSlug string, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, roleDA := range repo.roles {
		if roleDA.Description.String == roleSlug {
			roleDA.Permissions = append(roleDA.Permissions, permission.ID().String())
			repo.roles[roleDA.ID] = roleDA
			return nil
		}
	}
	return errors.New("role not found")
}

func (repo *BaseRepo) RemovePermissionFromRole(ctx context.Context, roleSlug string, permissionID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, roleDA := range repo.roles {
		if roleDA.Description.String == roleSlug {
			for i, pid := range roleDA.Permissions {
				if pid == permissionID {
					roleDA.Permissions = append(roleDA.Permissions[:i], roleDA.Permissions[i+1:]...)
					repo.roles[roleDA.ID] = roleDA
					return nil
				}
			}
		}
	}
	return errors.New("role or permission not found")
}

// Permission methods

func (repo *BaseRepo) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var permissions []Permission
	for _, permissionDA := range repo.permissions {
		permissions = append(permissions, toPermission(permissionDA))
	}
	return permissions, nil
}

func (repo *BaseRepo) GetPermissionByID(ctx context.Context, id string) (Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	permissionDA, exists := repo.permissions[id]
	if !exists {
		return Permission{}, errors.New("permission not found")
	}
	return toPermission(permissionDA), nil
}

func (repo *BaseRepo) CreatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[permission.ID().String()]; exists {
		return errors.New("permission already exists")
	}
	repo.permissions[permission.ID().String()] = toPermissionDA(permission)
	return nil
}

func (repo *BaseRepo) UpdatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	permissionDA := toPermissionDA(permission)
	if _, exists := repo.permissions[permissionDA.ID]; !exists {
		return errors.New("permission not found")
	}
	repo.permissions[permissionDA.ID] = permissionDA
	return nil
}

func (repo *BaseRepo) DeletePermission(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[id]; !exists {
		return errors.New("permission not found")
	}
	delete(repo.permissions, id)
	return nil
}

// Resource methods

func (repo *BaseRepo) GetAllResources(ctx context.Context) ([]Resource, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var resources []Resource
	for _, resourceDA := range repo.resources {
		resources = append(resources, toResource(resourceDA))
	}
	return resources, nil
}

func (repo *BaseRepo) GetResourceByID(ctx context.Context, id string) (Resource, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resourceDA, exists := repo.resources[id]
	if !exists {
		return Resource{}, errors.New("resource not found")
	}
	return toResource(resourceDA), nil
}

func (repo *BaseRepo) CreateResource(ctx context.Context, resource Resource) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[resource.ID().String()]; exists {
		return errors.New("resource already exists")
	}
	repo.resources[resource.ID().String()] = toResourceDA(resource)
	return nil
}

func (repo *BaseRepo) UpdateResource(ctx context.Context, resource Resource) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resourceDA := toResourceDA(resource)
	if _, exists := repo.resources[resourceDA.ID]; !exists {
		return errors.New("resource not found")
	}
	repo.resources[resourceDA.ID] = resourceDA
	return nil
}

func (repo *BaseRepo) DeleteResource(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[id]; !exists {
		return errors.New("resource not found")
	}
	delete(repo.resources, id)
	return nil
}

func (repo *BaseRepo) AddPermissionToResource(ctx context.Context, resourceID string, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resourceDA, exists := repo.resources[resourceID]
	if !exists {
		return errors.New("resource not found")
	}
	resourceDA.Permissions = append(resourceDA.Permissions, permission.ID().String())
	repo.resources[resourceID] = resourceDA
	return nil
}

func (repo *BaseRepo) RemovePermissionFromResource(ctx context.Context, resourceID string, permissionID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resourceDA, exists := repo.resources[resourceID]
	if !exists {
		return errors.New("resource not found")
	}
	for i, pid := range resourceDA.Permissions {
		if pid == permissionID {
			resourceDA.Permissions = append(resourceDA.Permissions[:i], resourceDA.Permissions[i+1:]...)
			repo.resources[resourceID] = resourceDA
			return nil
		}
	}
	return errors.New("permission not found")
}

// Debug method

func (repo *BaseRepo) Debug() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result string
	result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s %-50s\n", "Type", "ID", "NameID", "Slug", "Username", "EncPassword")
	for _, id := range repo.order {
		userDA := repo.users[id]
		result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s\n",
			userDA.Type, userDA.ID.String(), userDA.NameID.String, userDA.Slug.String, userDA.Name.String)
	}
	result = fmt.Sprintf("%s state:\n%s", repo.Name(), result)
	repo.Log().Info(result)
}

// Sample data for testing purposes

func (repo *BaseRepo) addSampleData() {
	for i := 1; i <= 5; i++ {
		id := uuid.New()
		username := fmt.Sprintf("sampleuser%d", i)
		email := fmt.Sprintf("sampleuser%d@example.com", i)
		user := NewUser(username, email, username) // Provide the correct number of arguments
		user.GenSlug("")                           // TODO: This function should be called without arguments later
		user.GenCreationValues()
		userDA := toUserDA(user)
		userDA.ID = id
		repo.users[id] = userDA
		repo.order = append(repo.order, id)
		repo.Log().Info("Created user with ID: ", id)
	}
}
