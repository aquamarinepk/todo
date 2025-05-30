package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ResourceDA represents the data access layer for the Resource model.
type ResourceDA struct {
	ID          uuid.UUID      `db:"id"`
	ShortID     sql.NullString `db:"short_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Label       sql.NullString `db:"label"`
	Type        sql.NullString `db:"type"`
	URI         sql.NullString `db:"uri"`
	Permissions []uuid.UUID    `db:"permissions"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert ResourceDA to Resource
// toModel methods do not preload relationships
func toResource(da ResourceDA) Resource {
	return Resource{
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithType(resourceEntityType),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:         da.Name.String,
		Description:  da.Description.String,
		Label:        da.Label.String,
		ResourceType: da.Type.String,
		URI:          da.URI.String,
	}
}

// Convert Resource to ResourceDA
func toResourceDA(resource Resource) ResourceDA {
	return ResourceDA{
		ID:          resource.ID(),
		Name:        sql.NullString{String: resource.Name, Valid: resource.Name != ""},
		Description: sql.NullString{String: resource.Description, Valid: resource.Description != ""},
		Label:       sql.NullString{String: resource.Label, Valid: resource.Label != ""},
		Type:        sql.NullString{String: resource.ResourceType, Valid: resource.ResourceType != ""},
		URI:         sql.NullString{String: resource.URI, Valid: resource.URI != ""},
		ShortID:     sql.NullString{String: resource.ShortID(), Valid: resource.ShortID() != ""},
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   sql.NullString{String: resource.CreatedBy().String(), Valid: resource.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: resource.UpdatedBy().String(), Valid: resource.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: resource.CreatedAt(), Valid: !resource.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: resource.UpdatedAt(), Valid: !resource.UpdatedAt().IsZero()},
	}
}

// ResourceExtDA represents the data access layer for the Resource with associated permissions.
type ResourceExtDA struct {
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
