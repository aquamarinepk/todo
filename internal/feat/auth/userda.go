package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// UserDA represents the data access layer for the User model.
type UserDA struct {
	ID            uuid.UUID      `db:"id"`
	Slug          sql.NullString `db:"slug"`
	Name          sql.NullString `db:"name"`
	Username      sql.NullString `db:"username"`
	Email         sql.NullString `db:"email"`
	EncPassword   sql.NullString `db:"password"`
	RoleIDs       []uuid.UUID
	PermissionIDs []uuid.UUID
	CreatedBy     sql.NullString `db:"created_by"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
}

// Convert UserDA to User
// toModel methods do not preload relationships

func ToUser(da UserDA) User {
	return User{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Username:    da.Username.String,
		Email:       da.Email.String,
		EncPassword: da.EncPassword.String,
	}
}

func ToUsers(das []UserDA) []User {
	users := make([]User, len(das))
	for i, da := range das {
		users[i] = ToUser(da)
	}
	return users
}

// Convert User to UserDA
func toUserDA(user User) UserDA {
	return UserDA{
		ID:            user.ID(),
		Slug:          sql.NullString{String: user.Slug(), Valid: user.Slug() != ""},
		Name:          sql.NullString{String: user.Name, Valid: user.Name != ""},
		Username:      sql.NullString{String: user.Username, Valid: user.Username != ""},
		Email:         sql.NullString{String: user.Email, Valid: user.Email != ""},
		EncPassword:   sql.NullString{String: user.EncPassword, Valid: user.EncPassword != ""},
		RoleIDs:       toRoleIDs(user.Roles),
		PermissionIDs: toPermissionIDs(user.Permissions),
		CreatedBy:     sql.NullString{String: user.Model.CreatedBy().String(), Valid: user.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:     sql.NullString{String: user.Model.UpdatedBy().String(), Valid: user.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:     sql.NullTime{Time: user.Model.CreatedAt(), Valid: !user.Model.CreatedAt().IsZero()},
		UpdatedAt:     sql.NullTime{Time: user.Model.UpdatedAt(), Valid: !user.Model.UpdatedAt().IsZero()},
	}
}

// UserExtDA represents the data access layer for the UserRolePermission.
type UserExtDA struct {
	ID             uuid.UUID      `db:"id"`
	Slug           sql.NullString `db:"slug"`
	Name           sql.NullString `db:"name"`
	Username       sql.NullString `db:"username"`
	Email          sql.NullString `db:"email"`
	EncPassword    sql.NullString `db:"password"`
	RoleID         sql.NullString `db:"role_id"`
	PermissionID   sql.NullString `db:"permission_id"`
	RoleName       sql.NullString `db:"role_name"`
	PermissionName sql.NullString `db:"permission_name"`
	CreatedBy      sql.NullString `db:"created_by"`
	UpdatedBy      sql.NullString `db:"updated_by"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}

// ToUserExt converts UserExtDA to User including roles and permissions.
func ToUserExt(da UserExtDA) User {
	return User{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Username:    da.Username.String,
		Email:       da.Email.String,
		EncPassword: da.EncPassword.String,
	}
}
