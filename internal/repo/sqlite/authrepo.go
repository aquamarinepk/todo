package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/auth"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type AuthRepo struct {
	am.Repo
	db *sql.DB
}

func NewAuthRepo(dsn string, qm *am.QueryManager, opts ...am.Option) (*AuthRepo, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return &AuthRepo{
		Repo: *am.NewRepo("sqlite-auth-repo", qm, opts...),
		db:   db,
	}, nil
}

// User methods

func (repo *AuthRepo) GetAllUsers(ctx context.Context) ([]auth.User, error) {
	query, err := repo.Query.Get("sqlite", "auth", "user", "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.User
	for rows.Next() {
		var user auth.User
		baseModel := user.Model.(*am.BaseModel)
		if err := rows.Scan(baseModel.ID(), &user.Username, &user.Email, baseModel.Slug()); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *AuthRepo) GetUser(ctx context.Context, id uuid.UUID) (auth.User, error) {
	query, err := repo.Query.Get("sqlite", "auth", "user", "Get")
	if err != nil {
		return auth.User{}, err
	}

	var user auth.User
	baseModel := user.Model.(*am.BaseModel)
	err = repo.db.QueryRowContext(ctx, query, id).Scan(baseModel.ID(), &user.Username, &user.Email, baseModel.Slug())
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (repo *AuthRepo) CreateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query.Get("sqlite", "auth", "user", "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, user.ID(), user.Username, user.Email, user.Slug())
	return err
}

func (repo *AuthRepo) UpdateUser(ctx context.Context, user auth.User) error {
	query, err := repo.Query.Get("sqlite", "auth", "user", "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, user.Username, user.Email, user.Slug(), user.ID())
	return err
}

func (repo *AuthRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get("sqlite", "auth", "user", "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Role methods

func (repo *AuthRepo) GetAllRoles(ctx context.Context) ([]auth.Role, error) {
	query, err := repo.Query.Get("sqlite", "auth", "role", "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []auth.Role
	for rows.Next() {
		var role auth.Role
		baseModel := role.Model.(*am.BaseModel)
		if err := rows.Scan(baseModel.ID(), &role.Name, &role.Description, baseModel.Slug()); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (repo *AuthRepo) GetRole(ctx context.Context, id uuid.UUID) (auth.Role, error) {
	query, err := repo.Query.Get("sqlite", "auth", "role", "Get")
	if err != nil {
		return auth.Role{}, err
	}

	var role auth.Role
	baseModel := role.Model.(*am.BaseModel)
	err = repo.db.QueryRowContext(ctx, query, id).Scan(baseModel.ID(), &role.Name, &role.Description, baseModel.Slug())
	if err != nil {
		if err == sql.ErrNoRows {
			return role, errors.New("role not found")
		}
		return role, err
	}
	return role, nil
}

func (repo *AuthRepo) CreateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query.Get("sqlite", "auth", "role", "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, role.ID(), role.Name, role.Description, role.Slug())
	return err
}

func (repo *AuthRepo) UpdateRole(ctx context.Context, role auth.Role) error {
	query, err := repo.Query.Get("sqlite", "auth", "role", "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, role.Name, role.Description, role.Slug(), role.ID())
	return err
}

func (repo *AuthRepo) DeleteRole(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get("sqlite", "auth", "role", "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Permission methods

func (repo *AuthRepo) GetAllPermissions(ctx context.Context) ([]auth.Permission, error) {
	query, err := repo.Query.Get("sqlite", "auth", "permission", "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []auth.Permission
	for rows.Next() {
		var permission auth.Permission
		baseModel := permission.Model.(*am.BaseModel)
		if err := rows.Scan(baseModel.ID(), &permission.Name, &permission.Description, baseModel.Slug()); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (repo *AuthRepo) GetPermission(ctx context.Context, id uuid.UUID) (auth.Permission, error) {
	query, err := repo.Query.Get("sqlite", "auth", "permission", "Get")
	if err != nil {
		return auth.Permission{}, err
	}

	var permission auth.Permission
	baseModel := permission.Model.(*am.BaseModel)
	err = repo.db.QueryRowContext(ctx, query, id).Scan(baseModel.ID(), &permission.Name, &permission.Description, baseModel.Slug())
	if err != nil {
		if err == sql.ErrNoRows {
			return permission, errors.New("permission not found")
		}
		return permission, err
	}
	return permission, nil
}

func (repo *AuthRepo) CreatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query.Get("sqlite", "auth", "permission", "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, permission.ID(), permission.Name, permission.Description, permission.Slug())
	return err
}

func (repo *AuthRepo) UpdatePermission(ctx context.Context, permission auth.Permission) error {
	query, err := repo.Query.Get("sqlite", "auth", "permission", "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, permission.Name, permission.Description, permission.Slug(), permission.ID())
	return err
}

func (repo *AuthRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get("sqlite", "auth", "permission", "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}

// Resource methods

func (repo *AuthRepo) GetAllResources(ctx context.Context) ([]auth.Resource, error) {
	query, err := repo.Query.Get("sqlite", "auth", "resource", "GetAll")
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []auth.Resource
	for rows.Next() {
		var resource auth.Resource
		baseModel := resource.Model.(*am.BaseModel)
		if err := rows.Scan(baseModel.ID(), &resource.Name, &resource.Description, baseModel.Slug()); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func (repo *AuthRepo) GetResource(ctx context.Context, id uuid.UUID) (auth.Resource, error) {
	query, err := repo.Query.Get("sqlite", "auth", "resource", "Get")
	if err != nil {
		return auth.Resource{}, err
	}

	var resource auth.Resource
	baseModel := resource.Model.(*am.BaseModel)
	err = repo.db.QueryRowContext(ctx, query, id).Scan(baseModel.ID(), &resource.Name, &resource.Description, baseModel.Slug())
	if err != nil {
		if err == sql.ErrNoRows {
			return resource, errors.New("resource not found")
		}
		return resource, err
	}
	return resource, nil
}

func (repo *AuthRepo) CreateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query.Get("sqlite", "auth", "resource", "Create")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, resource.ID(), resource.Name, resource.Description, resource.Slug())
	return err
}

func (repo *AuthRepo) UpdateResource(ctx context.Context, resource auth.Resource) error {
	query, err := repo.Query.Get("sqlite", "auth", "resource", "Update")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, resource.Name, resource.Description, resource.Slug(), resource.ID())
	return err
}

func (repo *AuthRepo) DeleteResource(ctx context.Context, id uuid.UUID) error {
	query, err := repo.Query.Get("sqlite", "auth", "resource", "Delete")
	if err != nil {
		return err
	}

	_, err = repo.db.ExecContext(ctx, query, id)
	return err
}
