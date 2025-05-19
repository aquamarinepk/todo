package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// BaseRepo provides an in-memory implementation of the Repo interface.
// This implementation is intended to simplify initial prototyping.
// In the future, a relational database implementation and possibly a NoSQL implementation will be provided.
type BaseRepo struct {
	*am.BaseRepo
	mu                  sync.Mutex
	users               map[uuid.UUID]UserDA
	roles               map[uuid.UUID]RoleDA
	permissions         map[uuid.UUID]PermissionDA
	resources           map[uuid.UUID]ResourceDA
	userRoles           map[uuid.UUID][]uuid.UUID
	userPermissions     map[uuid.UUID][]uuid.UUID
	rolePermissions     map[uuid.UUID][]uuid.UUID
	resourcePermissions map[uuid.UUID][]uuid.UUID
	order               []uuid.UUID
	emailKey            []byte
}

func NewRepo(qm *am.QueryManager, opts ...am.Option) *BaseRepo {
	repo := &BaseRepo{
		BaseRepo:            am.NewRepo("todo-repo", qm, opts...),
		users:               make(map[uuid.UUID]UserDA),
		roles:               make(map[uuid.UUID]RoleDA),
		permissions:         make(map[uuid.UUID]PermissionDA),
		resources:           make(map[uuid.UUID]ResourceDA),
		userRoles:           make(map[uuid.UUID][]uuid.UUID),
		userPermissions:     make(map[uuid.UUID][]uuid.UUID),
		rolePermissions:     make(map[uuid.UUID][]uuid.UUID),
		resourcePermissions: make(map[uuid.UUID][]uuid.UUID),
		order:               []uuid.UUID{},
		emailKey:            []byte{},
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
		result = append(result, ToUser(repo.users[id]))
	}
	return result, nil
}

func (repo *BaseRepo) GetUser(ctx context.Context, id uuid.UUID, preload ...bool) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if len(preload) > 0 && preload[0] {
		return repo.getUserPreload(ctx, id)
	}
	return repo.getUser(ctx, id)
}

func (repo *BaseRepo) getUser(ctx context.Context, id uuid.UUID) (User, error) {
	userDA, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return ToUser(userDA), nil
}

func (repo *BaseRepo) getUserPreload(ctx context.Context, id uuid.UUID) (User, error) {
	userDA, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}

	user := ToUser(userDA)
	user.Roles = repo.getUserRolesByID(id)
	user.Permissions = repo.getUserPermissionsByID(id)
	return user, nil
}

func (repo *BaseRepo) getUserRolesByID(userID uuid.UUID) []Role {
	var roles []Role
	for _, roleID := range repo.userRoles[userID] {
		roleDA := repo.roles[roleID]
		roles = append(roles, toRole(roleDA))
	}
	return roles
}

func (repo *BaseRepo) getUserPermissionsByID(userID uuid.UUID) []Permission {
	var permissions []Permission
	for _, permissionID := range repo.userPermissions[userID] {
		permissionDA := repo.permissions[permissionID]
		permissions = append(permissions, ToPermission(permissionDA))
	}
	return permissions
}

func toUserDA(user User) UserDA {
	return UserDA{
		ID:            user.ID(),
		Slug:          sql.NullString{String: user.Slug(), Valid: user.Slug() != ""},
		Name:          sql.NullString{String: user.Name, Valid: user.Name != ""},
		Username:      sql.NullString{String: user.Username, Valid: user.Username != ""},
		EmailEnc:      user.EmailEnc,
		PasswordEnc:   user.PasswordEnc,
		RoleIDs:       user.RoleIDs,
		PermissionIDs: user.PermissionIDs,
		CreatedBy:     sql.NullString{String: user.CreatedBy().String(), Valid: user.CreatedBy() != uuid.Nil},
		UpdatedBy:     sql.NullString{String: user.UpdatedBy().String(), Valid: user.UpdatedBy() != uuid.Nil},
		CreatedAt:     sql.NullTime{Time: user.CreatedAt(), Valid: !user.CreatedAt().IsZero()},
		UpdatedAt:     sql.NullTime{Time: user.UpdatedAt(), Valid: !user.UpdatedAt().IsZero()},
		LastLoginAt:   sql.NullTime{Time: derefTime(user.LastLoginAt), Valid: user.LastLoginAt != nil},
		LastLoginIP:   sql.NullString{String: user.LastLoginIP, Valid: user.LastLoginIP != ""},
		IsActive:      sql.NullBool{Bool: user.IsActive, Valid: true},
	}
}

func (repo *BaseRepo) CreateUser(ctx context.Context, u User) (User, error) {
	// Encrypt email and password
	emailEnc, err := EncryptEmail(string(u.EmailEnc), repo.emailKey)
	if err != nil {
		return User{}, err
	}

	passwordEnc, err := HashPassword(string(u.PasswordEnc))
	if err != nil {
		return User{}, err
	}

	user := NewUser(u.Username, u.Name)
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)
	user.RoleIDs = u.RoleIDs
	user.PermissionIDs = u.PermissionIDs

	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA := toUserDA(user)
	if _, exists := repo.users[userDA.ID]; exists {
		return User{}, errors.New("user already exists")
	}
	repo.users[userDA.ID] = userDA
	repo.order = append(repo.order, userDA.ID)
	return user, nil
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
	delete(repo.userRoles, id)
	delete(repo.userPermissions, id)
	return nil
}

func (repo *BaseRepo) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return nil, errors.New("user not found")
	}

	var roles []Role
	for _, roleID := range repo.userRoles[userID] {
		roleDA := repo.roles[roleID]
		roles = append(roles, toRole(roleDA))
	}
	return roles, nil
}

func (repo *BaseRepo) AddRole(ctx context.Context, userID uuid.UUID, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
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

func (repo *BaseRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
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

func (repo *BaseRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA, exists := repo.users[userID]
	if !exists {
		return errors.New("user not found")
	}
	userDA.PermissionIDs = append(userDA.PermissionIDs, permission.ID())
	repo.users[userDA.ID] = userDA
	repo.userPermissions[userDA.ID] = append(repo.userPermissions[userDA.ID], permission.ID())
	return nil
}

func (repo *BaseRepo) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	userDA, exists := repo.users[userID]
	if !exists {
		return errors.New("user not found")
	}

	for i, pid := range userDA.PermissionIDs {
		if pid == permissionID {
			userDA.PermissionIDs = append(userDA.PermissionIDs[:i], userDA.PermissionIDs[i+1:]...)
			repo.users[userDA.ID] = userDA
			for j, upid := range repo.userPermissions[userDA.ID] {
				if upid == permissionID {
					repo.userPermissions[userDA.ID] = append(repo.userPermissions[userDA.ID][:j], repo.userPermissions[userDA.ID][j+1:]...)
					break
				}
			}
			return nil
		}
	}
	return errors.New("permission not found")
}

// Role methods

func (repo *BaseRepo) GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return Role{}, errors.New("user not found")
	}

	for _, rid := range repo.userRoles[userID] {
		if rid == roleID {
			roleDA := repo.roles[rid]
			return toRole(roleDA), nil
		}
	}
	return Role{}, errors.New("role not found")
}

func (repo *BaseRepo) GetRole(ctx context.Context, userID, roleID uuid.UUID) (Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return Role{}, errors.New("user not found")
	}

	for _, rid := range repo.userRoles[userID] {
		if rid == roleID {
			roleDA := repo.roles[rid]
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

func (repo *BaseRepo) UpdateRole(ctx context.Context, userID uuid.UUID, role Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
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

func (repo *BaseRepo) DeleteRole(ctx context.Context, userID, roleID uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	if _, exists := repo.roles[roleID]; !exists {
		return errors.New("role not found")
	}
	delete(repo.roles, roleID)

	for i, rid := range repo.userRoles[userID] {
		if rid == roleID {
			repo.userRoles[userID] = append(repo.userRoles[userID][:i], repo.userRoles[userID][i+1:]...)
			break
		}
	}

	delete(repo.rolePermissions, roleID)
	return nil
}

// AddPermissionToRole adds a permission to a role.
func (repo *BaseRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	role, err := repo.GetRole(ctx, roleID, roleID)
	if err != nil {
		return err
	}

	permission, err := repo.GetPermission(ctx, permissionID)
	if err != nil {
		return err
	}

	role.Permissions = append(role.Permissions, permission)
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (repo *BaseRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	role, err := repo.GetRole(ctx, roleID, roleID)
	if err != nil {
		return err
	}

	for i, p := range role.Permissions {
		if p.ID() == permissionID {
			role.Permissions = append(role.Permissions[:i], role.Permissions[i+1:]...)
			return nil
		}
	}

	return errors.New(am.ErrResourceNotFound)
}

// Permission methods

func (repo *BaseRepo) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var permissions []Permission
	for _, permissionDA := range repo.permissions {
		permissions = append(permissions, ToPermission(permissionDA))
	}
	return permissions, nil
}

func (repo *BaseRepo) GetPermission(ctx context.Context, id uuid.UUID) (Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	permissionDA, exists := repo.permissions[id]
	if !exists {
		return Permission{}, errors.New("permission not found")
	}
	return ToPermission(permissionDA), nil
}

func (repo *BaseRepo) CreatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.permissions[permission.ID()]; exists {
		return errors.New("permission already exists")
	}
	repo.permissions[permission.ID()] = ToPermissionDA(permission)
	return nil
}

func (repo *BaseRepo) UpdatePermission(ctx context.Context, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	permissionDA := ToPermissionDA(permission)
	if _, exists := repo.permissions[permissionDA.ID]; !exists {
		return errors.New("permission not found")
	}
	repo.permissions[permissionDA.ID] = permissionDA
	return nil
}

func (repo *BaseRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
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

func (repo *BaseRepo) GetResource(ctx context.Context, id uuid.UUID) (Resource, error) {
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

	if _, exists := repo.resources[resource.ID()]; exists {
		return errors.New("resource already exists")
	}
	repo.resources[resource.ID()] = toResourceDA(resource)
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

func (repo *BaseRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[id]; !exists {
		return errors.New("resource not found")
	}
	delete(repo.resources, id)
	return nil
}

func (repo *BaseRepo) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission Permission) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	resourceDA, exists := repo.resources[resourceID]
	if !exists {
		return errors.New("resource not found")
	}
	resourceDA.Permissions = append(resourceDA.Permissions, permission.ID())
	repo.resources[resourceID] = resourceDA
	repo.resourcePermissions[resourceID] = append(repo.resourcePermissions[resourceID], permission.ID())
	return nil
}

func (repo *BaseRepo) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
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
			for j, rpid := range repo.resourcePermissions[resourceID] {
				if rpid == permissionID {
					repo.resourcePermissions[resourceID] = append(repo.resourcePermissions[resourceID][:j], repo.resourcePermissions[resourceID][j+1:]...)
					break
				}
			}
			return nil
		}
	}
	return errors.New("permission not found")
}

func (repo *BaseRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]Permission, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.resources[resourceID]; !exists {
		return nil, errors.New("resource not found")
	}

	var permissions []Permission
	for _, permissionID := range repo.resourcePermissions[resourceID] {
		permissionDA := repo.permissions[permissionID]
		permissions = append(permissions, ToPermission(permissionDA))
	}
	return permissions, nil
}

func (repo *BaseRepo) Debug() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var result string
	result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s\n", "Type", "ID", "Slug", "Username", "Extra")
	for _, id := range repo.order {
		userDA := repo.users[id]
		result += fmt.Sprintf("%-10s %-36s %-36s %-36s %-20s\n", "User", userDA.ID, userDA.Slug.String, userDA.Name.String, userDA.Username.String)
	}
	fmt.Println(result)
}

func (repo *BaseRepo) addSampleData() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	// Add sample users
	emailEnc, _ := EncryptEmail("john@example.com", repo.emailKey)
	passwordEnc, _ := HashPassword("password")
	user := NewUser("john", "John Doe")
	user.SetEmailEnc(emailEnc)
	user.SetPasswordEnc(passwordEnc)
	user.RoleIDs = []uuid.UUID{repo.roles[uuid.MustParse("00000000-0000-0000-0000-000000000001")].ID}
	user.PermissionIDs = []uuid.UUID{repo.permissions[uuid.MustParse("00000000-0000-0000-0000-000000000001")].ID}
	userDA := toUserDA(user)
	userDA.ID = uuid.New()
	repo.users[userDA.ID] = userDA
	repo.order = append(repo.order, userDA.ID)

	// Add sample roles
	role := NewRole("admin", "Administrator", "Administrator role with full access")
	role.PermissionIDs = []uuid.UUID{repo.permissions[uuid.MustParse("00000000-0000-0000-0000-000000000001")].ID}
	roleDA := toRoleDA(role)
	roleDA.ID = uuid.New()
	repo.roles[roleDA.ID] = roleDA
	repo.order = append(repo.order, roleDA.ID)

	// Add sample permissions
	perm := NewPermission("read", "Read permission")
	permDA := ToPermissionDA(perm)
	permDA.ID = uuid.New()
	repo.permissions[permDA.ID] = permDA
	repo.order = append(repo.order, permDA.ID)

	// Assign roles to users
	repo.userRoles[userDA.ID] = []uuid.UUID{roleDA.ID}

	// Assign permissions to roles
	repo.rolePermissions[roleDA.ID] = []uuid.UUID{permDA.ID}

	// Add sample resources
	for i := 1; i <= 3; i++ {
		id := uuid.New()
		name := fmt.Sprintf("resource%d", i)
		description := fmt.Sprintf("%s description", name)
		resource := NewResource(name, description, "entity")
		resource.GenSlug()
		resource.GenCreationValues()
		resourceDA := toResourceDA(resource)
		resourceDA.ID = id
		repo.resources[id] = resourceDA
		repo.Log().Info("Created resource with ID: ", id)
	}

	// Assign permissions to resources
	for resourceID := range repo.resources {
		repo.resourcePermissions[resourceID] = []uuid.UUID{permDA.ID}
	}
}

func (repo *BaseRepo) GetAllRoles(ctx context.Context) ([]Role, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	var roles []Role
	for _, roleDA := range repo.roles {
		roles = append(roles, toRole(roleDA))
	}
	return roles, nil
}
