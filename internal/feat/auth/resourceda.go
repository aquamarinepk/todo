package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ResourceDA represents the data access layer for the Resource model.
type ResourceDA struct {
	ID          uuid.UUID      `db:"id"`
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
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithCreatedBy(uuid.MustParse(da.CreatedBy.String)),
			am.WithUpdatedBy(uuid.MustParse(da.UpdatedBy.String)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
		Label:       da.Label.String,
		Type:        da.Type.String,
		URI:         da.URI.String,
	}
}

// Convert Resource to ResourceDA
func toResourceDA(resource Resource) ResourceDA {
	return ResourceDA{
		ID:          resource.ID(),
		Name:        sql.NullString{String: resource.Name, Valid: resource.Name != ""},
		Description: sql.NullString{String: resource.Description, Valid: resource.Description != ""},
		Label:       sql.NullString{String: resource.Label, Valid: resource.Label != ""},
		Type:        sql.NullString{String: resource.Type, Valid: resource.Type != ""},
		URI:         sql.NullString{String: resource.URI, Valid: resource.URI != ""},
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   sql.NullString{String: resource.Model.CreatedBy().String(), Valid: resource.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: resource.Model.UpdatedBy().String(), Valid: resource.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: resource.Model.CreatedAt(), Valid: !resource.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: resource.Model.UpdatedAt(), Valid: !resource.Model.UpdatedAt().IsZero()},
	}
}

