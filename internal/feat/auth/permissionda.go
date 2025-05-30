package auth

import (
	"database/sql"

	"github.com/google/uuid"
)

// PermissionDA represents the data access layer for the Permission model.
type PermissionDA struct {
	ID          uuid.UUID      `db:"id"`
	ShortID     sql.NullString `db:"short_id"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}
