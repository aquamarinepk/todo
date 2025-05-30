package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// RoleDA represents the data access layer for the Role model.
type RoleDA struct {
	ID          uuid.UUID      `db:"id"`
	ShortID     sql.NullString `db:"short_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Status      sql.NullString `db:"status"`
	Permissions []uuid.UUID
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert RoleDA to Role
// toModel methods do not preload relationships
func toRole(da RoleDA) Role {
	return Role{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(roleType),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		Status:        da.Status.String,
		PermissionIDs: da.Permissions,
		Permissions:   []Permission{},
	}
}

// Convert Role to RoleDA
// toModel methods do not preload relationships
func toRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID(),
		Name:        sql.NullString{String: role.Name, Valid: role.Name != ""},
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		Permissions: toPermissionIDs(role.Permissions),
		CreatedBy:   sql.NullString{String: role.CreatedBy().String(), Valid: role.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.UpdatedBy().String(), Valid: role.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.CreatedAt(), Valid: !role.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.UpdatedAt(), Valid: !role.UpdatedAt().IsZero()},
	}
}

// RoleExtDA represents the data access layer for the Role with associated permissions.
type RoleExtDA struct {
	ID             uuid.UUID      `db:"id"`
	Name           sql.NullString `db:"name"`
	Description    sql.NullString `db:"description"`
	ShortID        sql.NullString `db:"short_id"`
	PermissionID   sql.NullString `db:"permission_id"`
	PermissionName sql.NullString `db:"permission_name"`
	CreatedBy      sql.NullString `db:"created_by"`
	UpdatedBy      sql.NullString `db:"updated_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}
