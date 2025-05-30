package auth

import (
	"database/sql"

	"github.com/google/uuid"
)

// UserDA represents the data access layer for the User model.
type UserDA struct {
	ID            uuid.UUID      `db:"id"`
	ShortID       sql.NullString `db:"short_id"`
	Name          sql.NullString `db:"name"`
	Username      sql.NullString `db:"username"`
	EmailEnc      []byte         `db:"email_enc"`
	PasswordEnc   []byte         `db:"password_enc"`
	RoleIDs       []uuid.UUID
	PermissionIDs []uuid.UUID
	CreatedBy     sql.NullString `db:"created_by"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	LastLoginAt   sql.NullTime   `db:"last_login_at"`
	LastLoginIP   sql.NullString `db:"last_login_ip"`
	IsActive      sql.NullBool   `db:"is_active"`
}

// UserExtDA represents the data access layer for the UserRolePermission.
type UserExtDA struct {
	ID             uuid.UUID      `db:"id"`
	ShortID        sql.NullString `db:"short_id"`
	Name           sql.NullString `db:"name"`
	Username       sql.NullString `db:"username"`
	EmailEnc       []byte         `db:"email_enc"`
	PasswordEnc    []byte         `db:"password_enc"`
	RoleID         sql.NullString `db:"role_id"`
	PermissionID   sql.NullString `db:"permission_id"`
	RoleName       sql.NullString `db:"role_name"`
	PermissionName sql.NullString `db:"permission_name"`
	CreatedBy      sql.NullString `db:"created_by"`
	UpdatedBy      sql.NullString `db:"updated_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	LastLoginAt    sql.NullTime   `db:"last_login_at"`
	LastLoginIP    sql.NullString `db:"last_login_ip"`
	IsActive       sql.NullBool   `db:"is_active"`
}
