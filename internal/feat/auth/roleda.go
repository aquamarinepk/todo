package todo

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// RoleDA represents the data access layer for the Role model.
type RoleDA struct {
	ID          string         `db:"id"`
	UserID      uuid.UUID      `db:"user_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Status      sql.NullString `db:"status"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert RoleDA to Role
func toRole(da RoleDA) Role {
	return Role{
		Model: am.NewModel(
			am.WithID(uuid.MustParse(da.ID)),
			am.WithCreatedBy(uuid.MustParse(da.CreatedBy.String)),
			am.WithUpdatedBy(uuid.MustParse(da.UpdatedBy.String)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Description: da.Description.String,
		Status:      da.Status.String,
	}
}

// Convert Role to RoleDA
func toRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID().String(),
		UserID:      role.UserID,
		Name:        sql.NullString{String: role.Name, Valid: role.Name != ""},
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		Status:      sql.NullString{String: role.Status, Valid: role.Status != ""},
		CreatedBy:   sql.NullString{String: role.Model.CreatedBy().String(), Valid: role.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.Model.UpdatedBy().String(), Valid: role.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.Model.CreatedAt(), Valid: !role.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.Model.UpdatedAt(), Valid: !role.Model.UpdatedAt().IsZero()},
	}
}
