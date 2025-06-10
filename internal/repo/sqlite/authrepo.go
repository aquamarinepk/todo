package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	key           = am.Key
	featAuth      = "auth"
	resUser       = "user"
	resRole       = "role"
	resPerm       = "permission"
	resRes        = "resource"
	resUserRole   = "user_role"
	resUserPerm   = "user_permission"
	resRolePerm   = "role_permission"
	resResPerm    = "resource_permission"
	resOrg        = "org"
	resOrgOwner   = "org_owner"
	resTeamMember = "team_member"
	resTeam       = "team"
)

type AuthRepo struct {
	*am.BaseRepo
	db *sqlx.DB
}

func NewAuthRepo(qm *am.QueryManager, opts ...am.Option) *AuthRepo {
	return &AuthRepo{
		BaseRepo: am.NewRepo("sqlite-auth-repo", qm, opts...),
	}
}

// Setup the database connection.
func (repo *AuthRepo) Setup(ctx context.Context) error {
	dsn, ok := repo.Cfg().StrVal(key.DBSQLiteDSN)
	if !ok {
		return errors.New("database DSN not found in configuration")
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return fmt.Errorf("failed to set WAL mode: %w", err)
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

// DB returns the underlying *sqlx.DB for transaction management.
func (repo *AuthRepo) DB() *sqlx.DB {
	return repo.db
}

// getExec returns the correct Execer (Tx or DB) from context.
func (repo *AuthRepo) getExec(ctx context.Context) sqlx.ExtContext {
	tx, ok := am.TxFromContext(ctx)
	if ok {
		if sqlxTx, ok := tx.(*sqlx.Tx); ok {
			return sqlxTx
		}
	}
	return repo.db
}

func (r *AuthRepo) BeginTx(ctx context.Context) (context.Context, am.Tx, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, nil, err
	}
	ctxWithTx := am.WithTx(ctx, tx)
	return ctxWithTx, tx, nil
}

func (repo *AuthRepo) GetUsers(ctx context.Context) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resUser, "GetAll")
	if err != nil {
		return nil, err
	}

	var users []auth.UserDA
	err = repo.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		repo.Log().Infof("User: %+v", user)
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
	query, err := repo.Query().Get(featAuth, resUser, "Get")
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
	query, err := repo.Query().Get(featAuth, resUser, "GetPreload")
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
	query, err := repo.Query().Get(featAuth, resUser, "Create")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userDA.ID,
		userDA.Username,
		userDA.EmailEnc,
		userDA.Name,
		userDA.PasswordEnc,
		userDA.ShortID,
		userDA.CreatedBy,
		userDA.UpdatedBy,
		userDA.CreatedAt,
		userDA.UpdatedAt,
	)
	return err
}

func (repo *AuthRepo) UpdateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "Update")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userDA.Username, userDA.EmailEnc, userDA.Name,
		userDA.ShortID, userDA.UpdatedBy, userDA.UpdatedAt, userDA.ID)
	return err
}

func (repo *AuthRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resUser, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) UpdatePassword(ctx context.Context, user auth.User) error {
	query, err := repo.Query().Get(featAuth, resUser, "UpdatePassword")
	if err != nil {
		return err
	}

	userDA := auth.ToUserDA(user)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userDA.PasswordEnc, userDA.UpdatedBy, userDA.UpdatedAt, userDA.ID)
	return err
}

func (repo *AuthRepo) GetAllRoles(ctx context.Context) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resRole, "GetAll")
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
	query, err := repo.Query().Get(featAuth, resRole, "Get")
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
	query, err := repo.Query().Get(featAuth, resRole, "GetPreload")
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
	query, err := repo.Query().Get(featAuth, resRole, "Create")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		roleDA.ID,
		roleDA.Name,
		roleDA.Description,
		roleDA.ShortID,
		roleDA.CreatedBy,
		roleDA.UpdatedBy,
		roleDA.CreatedAt,
		roleDA.UpdatedAt)
	return err
}

func (repo *AuthRepo) UpdateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query().Get(featAuth, resRole, "Update")
	if err != nil {
		return err
	}

	roleDA := auth.ToRoleDA(role)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		roleDA.Name,
		roleDA.Description,
		roleDA.ShortID,
		roleDA.UpdatedBy,
		roleDA.UpdatedAt,
		roleDA.ID)
	return err
}

func (repo *AuthRepo) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resRole, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, roleID)
	return err
}

func (repo *AuthRepo) GetAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resPerm, "GetAll")
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

// GetPermission returns a permission by ID
func (repo *AuthRepo) GetPermission(ctx context.Context, id uuid.UUID) (auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resPerm, "Get")
	if err != nil {
		return auth.Permission{}, err
	}

	var permissionDA auth.PermissionDA
	if err := repo.db.GetContext(ctx, &permissionDA, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Permission{}, auth.ErrPermissionNotFound
		}
		return auth.Permission{}, err
	}

	return auth.ToPermission(permissionDA), nil
}

func (repo *AuthRepo) CreatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Create")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permissionDA.ID,
		permissionDA.ShortID,
		permissionDA.Name,
		permissionDA.Description,
		permissionDA.CreatedBy,
		permissionDA.UpdatedBy,
		permissionDA.CreatedAt,
		permissionDA.UpdatedAt,
	)
	return err
}

func (repo *AuthRepo) UpdatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Update")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		permissionDA.ShortID,
		permissionDA.Name,
		permissionDA.Description,
		permissionDA.UpdatedBy,
		permissionDA.UpdatedAt,
		permissionDA.ID)
	return err
}

func (repo *AuthRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resPerm, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) GetAllResources(ctx context.Context) ([]auth.Resource, error) {
	query, err := repo.Query().Get(featAuth, resRes, "GetAll")
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
	query, err := repo.Query().Get(featAuth, resRes, "Get")
	if err != nil {
		return auth.Resource{}, err
	}

	var resourceDA auth.ResourceDA
	if err := repo.db.GetContext(ctx, &resourceDA, query, id); err != nil {
		if err == sql.ErrNoRows {
			return auth.Resource{}, auth.ErrResourceNotFound
		}
		return auth.Resource{}, err
	}

	return auth.ToResource(resourceDA), nil
}

func (repo *AuthRepo) getResourcePreload(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query().Get(featAuth, resRes, "GetPreload")
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
	query, err := repo.Query().Get(featAuth, resRes, "Create")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resourceDA.ID,
		resourceDA.Name,
		resourceDA.Description,
		resourceDA.ShortID,
		resourceDA.CreatedBy,
		resourceDA.UpdatedBy,
		resourceDA.CreatedAt,
		resourceDA.UpdatedAt,
	)
	return err
}

func (repo *AuthRepo) UpdateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query().Get(featAuth, resRes, "Update")
	if err != nil {
		return err
	}

	resourceDA := auth.ToResourceDA(resource)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		resourceDA.Name,
		resourceDA.Description,
		resourceDA.ShortID,
		resourceDA.UpdatedBy,
		resourceDA.UpdatedAt,
		resourceDA.ID)
	return err
}

func (repo *AuthRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resRes, "Delete")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id)
	return err
}

func (repo *AuthRepo) GetUserAssignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserAssignedRoles")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}
	return auth.ToRoles(rolesDA), nil
}

// GetUserAssignedPermissions retrieves all permissions assigned to a user, both directly and through roles.
func (repo *AuthRepo) GetUserAssignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserAssignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (repo *AuthRepo) GetUserIndirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserIndirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetUserDirectPermissions retrieves permissions directly assigned to a user.
func (repo *AuthRepo) GetUserDirectPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserDirectPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetUserUnassignedPermissions retrieves permissions not assigned to a user, either directly or through roles.
func (repo *AuthRepo) GetUserUnassignedPermissions(ctx context.Context, userID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resUserPerm, "GetUserUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	err = repo.db.SelectContext(ctx, &permissionsDA, query, userID, userID)
	if err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (repo *AuthRepo) AddPermissionToUser(ctx context.Context, userID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resUserPerm, "AddPermissionToUser")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userID, permissionDA.ID)
	return err
}

func (repo *AuthRepo) RemovePermissionFromUser(ctx context.Context, userID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resUserPerm, "RemovePermissionFromUser")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, userID, permissionID)
	return err
}

func (repo *AuthRepo) GetUserUnassignedRoles(ctx context.Context, userID uuid.UUID, contextType, contextID string) ([]auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserUnassignedRoles")
	if err != nil {
		return nil, err
	}

	var rolesDA []auth.RoleDA
	err = repo.db.SelectContext(ctx, &rolesDA, query,
		userID.String(), contextType, contextID,
	)
	if err != nil {
		return nil, err
	}

	return auth.ToRoles(rolesDA), nil
}

func (repo *AuthRepo) AddRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error {
	query, err := repo.Query().Get(featAuth, resUserRole, "AddRole")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userID.String(), roleID.String(), contextType, contextID,
		roleID.String(),
	)
	return err
}

func (repo *AuthRepo) RemoveRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID, contextType, contextID string) error {
	query, err := repo.Query().Get(featAuth, resUserRole, "RemoveRole")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		userID.String(), roleID.String(), contextType, contextID,
	)
	return err
}

func (repo *AuthRepo) GetUserRole(ctx context.Context, userID, roleID uuid.UUID) (auth.Role, error) {
	query, err := repo.Query().Get(featAuth, resUserRole, "GetUserRole")
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

// AddPermissionToRole adds a permission to a role.
func (repo *AuthRepo) AddPermissionToRole(ctx context.Context, roleID uuid.UUID, permission auth.Permission) error {
	query := `
		INSERT INTO role_permission (role_id, permission_id)
		VALUES (?, ?)
	`
	exec := repo.getExec(ctx)
	_, err := exec.ExecContext(ctx, query, roleID, permission.ID())
	if err != nil {
		return fmt.Errorf("failed to add permission to role: %w", err)
	}
	return nil
}

// RemovePermissionFromRole removes a permission from a role.
func (repo *AuthRepo) RemovePermissionFromRole(ctx context.Context, roleID uuid.UUID, permissionID uuid.UUID) error {
	query := `
		DELETE FROM role_permission
		WHERE role_id = ? AND permission_id = ?
	`
	exec := repo.getExec(ctx)
	result, err := exec.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New(am.ErrResourceNotFound)
	}

	return nil
}

func (repo *AuthRepo) AddPermissionToResource(ctx context.Context, resourceID uuid.UUID, permission auth.Permission) error {
	query, err := repo.Query().Get(featAuth, resResPerm, "AddPermissionToResource")
	if err != nil {
		return err
	}

	permissionDA := auth.ToPermissionDA(permission)
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, resourceID, permissionDA.ID)
	return err
}

func (repo *AuthRepo) RemovePermissionFromResource(ctx context.Context, resourceID uuid.UUID, permissionID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resResPerm, "RemovePermissionFromResource")
	if err != nil {
		return err
	}

	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, resourceID, permissionID)
	return err
}

// GetRolePermissions returns all permissions assigned to a role
func (repo *AuthRepo) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resRolePerm, "GetRolePermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, roleID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetResourcePermissions returns all permissions assigned to a resource
func (repo *AuthRepo) GetResourcePermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourcePermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, resourceID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetResourceUnassignedPermissions returns all permissions not assigned to a resource
func (repo *AuthRepo) GetResourceUnassignedPermissions(ctx context.Context, resourceID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resResPerm, "GetResourceUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, resourceID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

// GetRoleUnassignedPermissions returns all permissions not assigned to a role
func (repo *AuthRepo) GetRoleUnassignedPermissions(ctx context.Context, roleID uuid.UUID) ([]auth.Permission, error) {
	query, err := repo.Query().Get(featAuth, resRolePerm, "GetRoleUnassignedPermissions")
	if err != nil {
		return nil, err
	}

	var permissionsDA []auth.PermissionDA
	if err := repo.db.SelectContext(ctx, &permissionsDA, query, roleID); err != nil {
		return nil, err
	}

	return auth.ToPermissions(permissionsDA), nil
}

func (r *AuthRepo) CreateOrg(ctx context.Context, org auth.Org) error {
	da := auth.ToOrgDA(org)
	query, err := r.Query().Get(featAuth, resOrg, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		da.ID, da.ShortID, da.Name, da.ShortDescription, da.Description, da.CreatedBy, da.UpdatedBy, da.CreatedAt, da.UpdatedAt,
	)
	return err
}

func (r *AuthRepo) GetDefaultOrg(ctx context.Context) (auth.Org, error) {
	query, err := r.Query().Get(featAuth, resOrg, "GetDefault")
	if err != nil {
		return auth.Org{}, err
	}
	row := r.db.QueryRowContext(ctx, query)
	var orgDA auth.OrgDA
	err = row.Scan(
		&orgDA.ID,
		&orgDA.ShortID,
		&orgDA.Name,
		&orgDA.ShortDescription,
		&orgDA.Description,
		&orgDA.CreatedBy,
		&orgDA.UpdatedBy,
		&orgDA.CreatedAt,
		&orgDA.UpdatedAt,
	)
	if err != nil {
		return auth.Org{}, err
	}
	return auth.ToOrg(orgDA), nil
}

func (r *AuthRepo) GetOrgOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgOwners")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = r.db.SelectContext(ctx, &usersDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (r *AuthRepo) GetOrgUnassignedOwners(ctx context.Context, orgID uuid.UUID) ([]auth.User, error) {
	query, err := r.Query().Get(featAuth, resOrgOwner, "GetOrgUnassignedOwners")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = r.db.SelectContext(ctx, &usersDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (r *AuthRepo) GetAllTeams(ctx context.Context, orgID uuid.UUID) ([]auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "GetAll")
	if err != nil {
		return nil, err
	}
	var teamsDA []auth.TeamDA
	err = r.db.SelectContext(ctx, &teamsDA, query, orgID.String())
	if err != nil {
		return nil, err
	}
	teams := make([]auth.Team, len(teamsDA))
	for i, da := range teamsDA {
		teams[i] = auth.ToTeam(da)
	}
	return teams, nil
}

func (r *AuthRepo) GetTeam(ctx context.Context, id uuid.UUID) (auth.Team, error) {
	query, err := r.Query().Get(featAuth, resTeam, "Get")
	if err != nil {
		return auth.Team{}, err
	}
	var da auth.TeamDA
	err = r.db.GetContext(ctx, &da, query, id.String())
	if err != nil {
		return auth.Team{}, err
	}
	return auth.ToTeam(da), nil
}

func (r *AuthRepo) CreateTeam(ctx context.Context, team auth.Team) error {
	da := auth.TeamDA{
		ID:               sql.NullString{String: team.ID().String(), Valid: team.ID() != uuid.Nil},
		OrgID:            sql.NullString{String: team.OrgID.String(), Valid: team.OrgID != uuid.Nil},
		ShortID:          team.ShortID(),
		Name:             team.Name,
		ShortDescription: team.ShortDescription,
		Description:      team.Description,
		CreatedBy:        sql.NullString{String: team.CreatedBy().String(), Valid: team.CreatedBy() != uuid.Nil},
		UpdatedBy:        sql.NullString{String: team.UpdatedBy().String(), Valid: team.UpdatedBy() != uuid.Nil},
		CreatedAt:        team.CreatedAt(),
		UpdatedAt:        team.UpdatedAt(),
	}
	query, err := r.Query().Get(featAuth, resTeam, "Create")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		da.ID, da.OrgID, da.ShortID, da.Name, da.ShortDescription, da.Description, da.CreatedBy, da.UpdatedBy, da.CreatedAt, da.UpdatedAt,
	)
	return err
}

func (r *AuthRepo) UpdateTeam(ctx context.Context, team auth.Team) error {
	query, err := r.Query().Get(featAuth, resTeam, "Update")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query,
		team.ShortID(), team.Name, team.ShortDescription, team.Description, team.UpdatedBy().String(), team.UpdatedAt(), team.ID().String(),
	)
	return err
}

func (r *AuthRepo) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	query, err := r.Query().Get(featAuth, resTeam, "Delete")
	if err != nil {
		return err
	}
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id.String())
	return err
}

func (r *AuthRepo) AddOrgOwner(ctx context.Context, orgID uuid.UUID, userID uuid.UUID) error {
	query, err := r.Query().Get(featAuth, resOrgOwner, "Add")
	if err != nil {
		return err
	}
	id := uuid.New()
	exec := r.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id.String(), orgID.String(), userID.String())
	return err
}

func (repo *AuthRepo) GetTeamMembers(ctx context.Context, teamID uuid.UUID) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resTeamMember, "ListTeamMembers")
	repo.Log().Debugf("GetTeamMembers query: %s", query)
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = repo.db.SelectContext(ctx, &usersDA, query, teamID.String())
	repo.Log().Debugf("GetTeamMembers teamID: %s", teamID.String())
	for _, user := range usersDA {
		repo.Log().Debugf("User: %+v", user)
	}
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (repo *AuthRepo) GetTeamUnassignedUsers(ctx context.Context, teamID uuid.UUID) ([]auth.User, error) {
	query, err := repo.Query().Get(featAuth, resTeamMember, "ListUsersNotInTeam")
	if err != nil {
		return nil, err
	}
	var usersDA []auth.UserDA
	err = repo.db.SelectContext(ctx, &usersDA, query, teamID.String())
	if err != nil {
		return nil, err
	}
	return auth.ToUsers(usersDA), nil
}

func (repo *AuthRepo) AddUserToTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, relationType string) error {
	query, err := repo.Query().Get(featAuth, resTeamMember, "AddUserToTeam")
	if err != nil {
		return err
	}
	id := uuid.New().String()
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, id, teamID.String(), userID.String(), relationType)
	return err
}

func (repo *AuthRepo) RemoveUserFromTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) error {
	query, err := repo.Query().Get(featAuth, resTeamMember, "RemoveUserFromTeam")
	if err != nil {
		return err
	}
	exec := repo.getExec(ctx)
	_, err = exec.ExecContext(ctx, query, teamID.String(), userID.String())
	return err
}
