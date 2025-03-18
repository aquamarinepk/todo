package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// UserDA represents the data access layer for the User model.
type UserDA struct {
	Type        string
	ID          uuid.UUID      `db:"id"`
	Slug        sql.NullString `db:"slug"`
	NameID      sql.NullString `db:"name_id"`
	Username    sql.NullString `db:"username"`
	Email       sql.NullString `db:"email"`
	Name        sql.NullString `db:"name"`
	EncPassword sql.NullString `db:"enc_password"`
	Roles       []string 
	Permissions []string
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

// Convert UserDA to User
func toUser(da UserDA) User {
	return User{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(uuid.MustParse(da.CreatedBy.String)),
			am.WithUpdatedBy(uuid.MustParse(da.UpdatedBy.String)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Username:    da.Name.String,
		Email:       da.Email.String,
		Name:        da.Name.String,
		EncPassword: da.EncPassword.String,
		Roles:       toRoles(da.Roles),
		Permissions: toPermissions(da.Permissions),
	}
}

// Convert User to UserDA
func toUserDA(user User) UserDA {
	return UserDA{
		ID:          user.ID(),
		Slug:        sql.NullString{String: user.Slug(), Valid: user.Slug() != ""},
		NameID:      sql.NullString{String: user.Model.NameID(), Valid: user.Model.NameID() != ""},
		Username:    sql.NullString{String: user.Username, Valid: user.Username != ""},
		Email:       sql.NullString{String: user.Email, Valid: user.Email != ""},
		Name:        sql.NullString{String: user.Username, Valid: user.Username != ""},
		Roles:       toRoleIDs(user.Roles),
		Permissions: toPermissionIDs(user.Permissions),
		CreatedBy:   sql.NullString{String: user.Model.CreatedBy().String(), Valid: user.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: user.Model.UpdatedBy().String(), Valid: user.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: user.Model.CreatedAt(), Valid: !user.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: user.Model.UpdatedAt(), Valid: !user.Model.UpdatedAt().IsZero()},
	}
}
