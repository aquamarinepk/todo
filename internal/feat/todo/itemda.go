package todo

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ItemDA represents the data access layer for the Item model.
type ItemDA struct {
	ID          string         `db:"id"`
	ListID      uuid.UUID      `db:"list_id"`
	Description sql.NullString `db:"description"`
	Status      sql.NullString `db:"status"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert ItemDA to Item
func toItem(da ItemDA) Item {
	return Item{
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

// Convert Item to ItemDA
func toItemDA(item Item) ItemDA {
	return ItemDA{
		ID:          item.ID().String(),
		ListID:      item.ListID,
		Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
		Status:      sql.NullString{String: item.Status, Valid: item.Status != ""},
		CreatedBy:   sql.NullString{String: item.Model.CreatedBy().String(), Valid: item.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: item.Model.UpdatedBy().String(), Valid: item.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: item.Model.CreatedAt(), Valid: !item.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: item.Model.UpdatedAt(), Valid: !item.Model.UpdatedAt().IsZero()},
	}
}
