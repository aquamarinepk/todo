package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key         = am.Key
	featAuth    = "auth"
	resUser     = "user"
	resRole     = "role"
	resPerm     = "permission"
	resRes      = "resource"
	resUserRole = "user_role"
	resUserPerm = "user_permission"
	resRolePerm = "role_permission"
	resResPerm  = "resource_permission"
)

type AuthRepo struct {
	*am.Repo
	db *sqlx.DB
}

func NewAuthRepo(qm *am.QueryManager, opts ...am.Option) *AuthRepo {
	return &AuthRepo{
		Repo: am.NewRepo("sqlite-auth-repo", qm, opts...),
	}
}

// Start opens the database connection.
func (repo *AuthRepo) Start(ctx context.Context) error {
	dsn, ok := repo.Cfg().StrVal(key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	repo.db = db
	return nil
}

// Stop closes the database connection.
func (repo *AuthRepo) Stop(ctx context.Context) error {
	if repo.db != nil {
		return repo.db.Close()
	}
	return nil
}

func (repo *AuthRepo) GetAllUsers(ctx context.Context) ([]auth.User, error) {
	query, err := repo.Query.Get(featAuth, resUser, "GetAll")
	if err != nil {
		return nil, err
	}

	var users []auth.UserDA
	err = repo.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return auth.ToUsers(users), nil
}

func (repo *AuthRepo) GetUser(ctx context.Context, id uuid.UUID, preload ...bool) (auth.User, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getUserPreload(ctx, id)
	}
	return repo.getUser(ctx, id)
}

func (repo *AuthRepo) getUser(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query.Get(featAuth, resUser, "Get")
	if err != nil {
		return auth.User{}, err
	}

	var user auth.UserDA
	err = repo.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return auth.User{}, err
	}

	return auth.ToUser(user), nil
}

func (repo *AuthRepo) getUserPreload(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query.Get(featAuth, resUser, "GetPreload")
	if err != nil {
		return auth.User{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.User{}, err
	}
	defer rows.Close()

	var userDA auth.UserExtDA
	var roles []uuid.UUID
	var permissions []uuid.UUID
	userMap := make(map[uuid.UUID]auth.User)

	for rows.Next() {
		if err := rows.StructScan(&userDA); err != nil {
			return auth.User{}, err
		}

		if _, exists := userMap[userDA.ID]; !exists {
			userMap[userDA.ID] = auth.ToUserExt(userDA)
		}

		if userDA.RoleID.Valid {
			roleID, err := uuid.Parse(userDA.RoleID.String)
			if err == nil {
				roles = append(roles, roleID)
			}
		}

		if userDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(userDA.PermissionID.String)
			if err == nil {
				permissions = append(permissions, permissionID)
			}
		}
	}

	user := userMap[userDA.ID]
	user.RoleIDs = roles
	user.PermissionIDs = permissions

	return user, nil
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query.Get(featAuth, resUser, "Create")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	_, err = repo.db.ExecContext(ctx, query, userDA.ID, userDA.Username, userDA.Email, userDA.Slug)
	return err
}

func (repo *AuthRepo) UpdateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query.Get(featAuth, resUser, "Update")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	_, err = repo.db.ExecContext(ctx, query, userDA.Username, userDA.Email, userDA.Name, 
		userDA.EncPassword, userDA.Slug, userDA.UpdatedBy, userDA.UpdatedAt, userDA.ID)
	return err
}

func (repo *AuthRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resUser, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) GetAllRoles(ctx context.Context) ([]auth.Role, error) {
	query, err := repo.Query.Get(featAuth, resRole, "GetAll")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}

// GetRole retrieves a role by its ID, optionally preloading its associated permissions.
func (repo *AuthRepo) GetRole(ctx context.Context, id uuid.UUID, preload ...bool) (auth.Role, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getRolePreload(ctx, id)
	}
	return repo.getRole(ctx, id)
}

func (repo *AuthRepo) getRole(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	query, err := repo.Query.Get("auth", "role", "Get")
	if err != nil {
		return auth.Role{}, err
	}

	var roleDA auth.RoleDA
	err = repo.db.GetContext(ctx, &roleDA, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return auth.ToRole(roleDA), nil
}

func (repo *AuthRepo) getRolePreload(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	query, err := repo.Query.Get("auth", "role", "GetPreload")
	if err != nil {
		return auth.Role{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.Role{}, err
	}
	defer rows.Close()

	var roleDA auth.RoleExtDA
	roleMap := make(map[uuid.UUID]auth.Role)

	for rows.Next() {
		if err := rows.StructScan(&roleDA); err != nil {
			return auth.Role{}, err
		}

		role, exists := roleMap[roleDA.ID]
		if !exists {
			role = auth.ToRoleExt(roleDA)
		}

		if roleDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(roleDA.PermissionID.String)
			if err == nil {
				role.PermissionIDs = append(role.PermissionIDs, permissionID)
			}
		}

		roleMap[roleDA.ID] = role
	}

	return roleMap[roleDA.ID], nil
}

func (repo *AuthRepo) CreateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query.Get(featAuth, resRole, "Create")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	_, err = repo.db.ExecContext(ctx, query, roleDA.ID, roleDA.Name, roleDA.Description, roleDA.Slug)
	return err
}

func (repo *AuthRepo) UpdateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query.Get(featAuth, resRole, "Update")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	_, err = repo.db.ExecContext(ctx, query, roleDA.Name, roleDA.Description, roleDA.Slug, roleDA.ID)
	return err
}

func (repo *AuthRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resRole, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, roleID)
	return err
}

func (repo *AuthRepo) GetAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	query, err := repo.Query.Get(featAuth, resPerm, "GetAll")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToPermissions(permissionsDA), nil
}

func (repo *AuthRepo) GetPermission(ctx context.Context, id uuid.UUID) (auth.Permission, error) {
	query, err := repo.Query.Get(featAuth, resPerm, "Get")
	if err != nil {
		return auth.Permission{}, err
	}

	var permissionDA auth.PermissionDA
	err = repo.db.GetContext(ctx, &permissionDA, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Permission{}, errors.New("permission not found")
		}
		return auth.Permission{}, err
	}
	return auth.ToPermission(permissionDA), nil
}

func (repo *AuthRepo) CreatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query.Get(featAuth, resPerm, "Create")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	_, err = repo.db.ExecContext(ctx, query, permissionDA.ID, permissionDA.Name, permissionDA.Description, permissionDA.Slug)
	return err
}

func (repo *AuthRepo) UpdatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query.Get(featAuth, resPerm, "Update")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	_, err = repo.db.ExecContext(ctx, query, permissionDA.Name, permissionDA.Description, permissionDA.Slug, permissionDA.ID)
	return err
}

func (repo *AuthRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resPerm, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) GetAllResources(ctx context.Context) ([]auth.Resource, error) {
	query, err := repo.Query.Get(featAuth, resRes, "GetAll")
	if err != nil {
		return nil, err
	}

	var resourcesDA []auth.ResourceDA
	err = repo.db.SelectContext(ctx, &resourcesDA, query)
	if err != nil {
		return nil, err
	}
	return auth.ToResources(resourcesDA), nil
}

// GetResource retrieves a resource by its ID, optionally preloading its associated permissions.
func (repo *AuthRepo) GetResource(ctx context.Context, id uuid.UUID, preload ...bool) (auth.Resource, error) {
	if len(preload) > 0 && preload[0] {
		return repo.getResourcePreload(ctx, id)
	}
	return repo.getResource(ctx, id)
}

func (repo *AuthRepo) getResource(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query.Get("auth", "resource", "Get")
	if err != nil {
		return auth.Resource{}, err
	}

	var resourceDA auth.ResourceDA
	err = repo.db.GetContext(ctx, &resourceDA, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Resource{}, errors.New("resource not found")
		}
		return auth.Resource{}, err
	}
	return auth.ToResource(resourceDA), nil
}

func (repo *AuthRepo) getResourcePreload(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query.Get("auth", "resource", "GetPreload")
	if err != nil {
		return auth.Resource{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, id)
	if err != nil {
		return auth.Resource{}, err
	}
	defer rows.Close()

	var resourceDA auth.ResourceExtDA
	resourceMap := make(map[uuid.UUID]auth.Resource)

	for rows.Next() {
		if err := rows.StructScan(&resourceDA); err != nil {
			return auth.Resource{}, err
		}

		resource, exists := resourceMap[resourceDA.ID]
		if !exists {
			resource = auth.ToResourceExt(resourceDA)
		}

		if resourceDA.PermissionID.Valid {
			permissionID, err := uuid.Parse(resourceDA.PermissionID.String)
			if err == nil {
				resource.PermissionIDs = append(resource.PermissionIDs, permissionID)
			}
		}

		resourceMap[resourceDA.ID] = resource
	}

	return resourceMap[resourceDA.ID], nil
}

func (repo *AuthRepo) CreateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query.Get(featAuth, resRes, "Create")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	_, err = repo.db.ExecContext(ctx, query, resourceDA.ID, resourceDA.Name, resourceDA.Description, resourceDA.Slug)
	return err
}

func (repo *AuthRepo) UpdateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query.Get(featAuth, resRes, "Update")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	_, err = repo.db.ExecContext(ctx, query, resourceDA.Name, resourceDA.Description, resourceDA.Slug, resourceDA.ID)
	return err
}

func (repo *AuthRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resRes, "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]auth.Role, error) {
	query, err := repo.Query.Get(featAuth, resUserRole, "GetUserRoles")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query, userID)
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}

func (repo *AuthRepo) AddRole(ctx context.Context, userID uuid.UUID, role auth.Role) error {
	query, err := repo.Query.Get(featAuth, resUserRole, "AddRole")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	_, err = repo.db.ExecContext(ctx, query, userID, roleDA.ID)
	return err
}

func (repo *AuthRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resUserRole, "RemoveRole")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, userID, roleID)
	return err
}

func (repo *AuthRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query.Get(featAuth, resUserPerm, "AddPermissionToUser")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	_, err = repo.db.ExecContext(ctx, query, userID, permissionDA.ID)
	return err
}

func (repo *AuthRepo) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resUserPerm, "RemovePermissionFromUser")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, userID, permissionID)
	return err
}

func (repo *AuthRepo) GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (auth.Role, error) {
	query, err := repo.Query.Get(featAuth, resUserRole, "GetUserRole")
	if err != nil {
		return auth.Role{}, err
	}

	var roleDA auth.RoleDA
	err = repo.db.GetContext(ctx, &roleDA, query, userID, roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return auth.Role{}, errors.New("role not found")
		}
		return auth.Role{}, err
	}
	return auth.ToRole(roleDA), nil
}

func (repo *AuthRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query.Get(featAuth, resRolePerm, "AddPermissionToRole")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	_, err = repo.db.ExecContext(ctx, query, roleID, permissionDA.ID)
	return err
}

func (repo *AuthRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resRolePerm, "RemovePermissionFromRole")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

func (repo *AuthRepo) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query.Get(featAuth, resResPerm, "AddPermissionToResource")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	_, err = repo.db.ExecContext(ctx, query, resourceID, permissionDA.ID)
	return err
}

func (repo *AuthRepo) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query.Get(featAuth, resResPerm, "RemovePermissionFromResource")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, resourceID, permissionID)
	return err
}

func (repo *AuthRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query.Get(featAuth, resResPerm, "GetResourcePermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, resourceID)
	if err != nil {
		return nil, err
	}
	return auth.ToPermissions(permissionsDA), nil
}
