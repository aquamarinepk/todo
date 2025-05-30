package todo

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ListDA represents the data access layer for the List model.
type ListDA struct {
	Type        string
	ID          uuid.UUID      `db:"id"`
	ShortID     sql.NullString `db:"short_id"`
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
		BaseModel: am.NewModel(
			am.WithID(da.ID),
			am.WithShortID(da.ShortID.String),
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
		ShortID:     sql.NullString{String: list.ShortID(), Valid: list.Slug() != ""},
		Name:        sql.NullString{String: list.Name, Valid: list.Name != ""},
		Description: sql.NullString{String: list.Description, Valid: list.Description != ""},
		CreatedBy:   sql.NullString{String: list.CreatedBy().String(), Valid: list.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: list.UpdatedBy().String(), Valid: list.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: list.CreatedAt(), Valid: !list.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: list.UpdatedAt(), Valid: !list.UpdatedAt().IsZero()},
	}
}
