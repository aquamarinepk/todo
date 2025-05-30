package auth

import (
	"database/sql"
	"time"
)

type TeamDA struct {
	ID               sql.NullString `db:"id"`
	OrgID            sql.NullString `db:"org_id"`
	ShortID          string         `db:"short_id"`
	Name             string         `db:"name"`
	ShortDescription string         `db:"short_description"`
	Description      string         `db:"description"`
	CreatedBy        sql.NullString `db:"created_by"`
	UpdatedBy        sql.NullString `db:"updated_by"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}
