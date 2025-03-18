package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// RoleDA represents the data access layer for the Role model.
type RoleDA struct {
	ID          string         `db:"id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Status      sql.NullString `db:"status"`
	Permissions []string
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
		Permissions: toPermissions(da.Permissions),
	}
}

// Convert Role to RoleDA
func toRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID().String(),
		Name:        sql.NullString{String: role.Name, Valid: role.Name != ""},
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		Permissions: toPermissionIDs(role.Permissions),
		CreatedBy:   sql.NullString{String: role.Model.CreatedBy().String(), Valid: role.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.Model.UpdatedBy().String(), Valid: role.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.Model.CreatedAt(), Valid: !role.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.Model.UpdatedAt(), Valid: !role.Model.UpdatedAt().IsZero()},
	}
}

func toRoles(roleIDs []string) []Role {
	var roles []Role
	for _, id := range roleIDs {
		roles = append(roles, Role{
			Model: am.NewModel(am.WithID(uuid.MustParse(id))),
		})
	}
	return roles
}

func toRoleIDs(roles []Role) []string {
	var roleIDs []string
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID().String())
	}
	return roleIDs
}
