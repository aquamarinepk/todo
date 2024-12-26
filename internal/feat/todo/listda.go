package todo

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ListDA represents the data access layer for the List model.
type ListDA struct {
	ID          uuid.UUID      `db:"id"`
	Slug        sql.NullString `db:"slug"`
	NameID      sql.NullString `db:"name_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert ListDA to List
func toList(da ListDA) List {
	return List{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithNameID(da.NameID.String),
			am.WithCreatedBy(uuid.MustParse(da.CreatedBy.String)),
			am.WithUpdatedBy(uuid.MustParse(da.UpdatedBy.String)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// Convert List to ListDA
func toListDA(list List) ListDA {
	return ListDA{
		ID:          list.ID(),
		Slug:        sql.NullString{String: list.Slug(), Valid: list.Slug() != ""},
		NameID:      sql.NullString{String: list.Model.NameID(), Valid: list.Model.NameID() != ""},
		Name:        sql.NullString{String: list.Name, Valid: list.Name != ""},
		Description: sql.NullString{String: list.Description, Valid: list.Description != ""},
		CreatedBy:   sql.NullString{String: list.Model.CreatedBy().String(), Valid: list.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: list.Model.UpdatedBy().String(), Valid: list.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: list.Model.CreatedAt(), Valid: !list.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: list.Model.UpdatedAt(), Valid: !list.Model.UpdatedAt().IsZero()},
	}
}
