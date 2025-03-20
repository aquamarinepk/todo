package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// PermissionDA represents the data access layer for the Permission model.
type PermissionDA struct {
	ID          uuid.UUID      `db:"id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert PermissionDA to Permission
func toPermission(da PermissionDA) Permission {
	return Permission{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithCreatedBy(uuid.MustParse(da.CreatedBy.String)),
			am.WithUpdatedBy(uuid.MustParse(da.UpdatedBy.String)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// Convert Permission to PermissionDA
func toPermissionDA(permission Permission) PermissionDA {
	return PermissionDA{
		ID:          permission.ID(),
		Name:        sql.NullString{String: permission.Name, Valid: permission.Name != ""},
		Description: sql.NullString{String: permission.Description, Valid: permission.Description != ""},
		CreatedBy:   sql.NullString{String: permission.Model.CreatedBy().String(), Valid: permission.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: permission.Model.UpdatedBy().String(), Valid: permission.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: permission.Model.CreatedAt(), Valid: !permission.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: permission.Model.UpdatedAt(), Valid: !permission.Model.UpdatedAt().IsZero()},
	}
}

func toPermissionIDs(permissions []Permission) []uuid.UUID {
	var permissionIDs []uuid.UUID
	for _, permission := range permissions {
		permissionIDs = append(permissionIDs, permission.ID())
	}
	return permissionIDs
}
