package auth

import (
	"database/sql"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// ToUserDA converts a User business object to a UserDA data access object
func ToUserDA(user User) UserDA {
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

// ToUser converts a UserDA data access object to a User business object
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

// ToUsers converts a slice of UserDA to a slice of User business objects
func ToUsers(das []UserDA) []User {
	users := make([]User, len(das))
	for i, da := range das {
		users[i] = ToUser(da)
	}
	return users
}

// ToUserExt converts UserExtDA to User including roles and permissions
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

// ToRoleDA converts a Role business object to a RoleDA data access object
func ToRoleDA(role Role) RoleDA {
	return RoleDA{
		ID:          role.ID(),
		Name:        sql.NullString{String: role.Name, Valid: role.Name != ""},
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		Slug:        sql.NullString{String: role.Slug(), Valid: role.Slug() != ""},
		Status:      sql.NullString{String: role.Status, Valid: role.Status != ""},
		Permissions: toPermissionIDs(role.Permissions),
		CreatedBy:   sql.NullString{String: role.Model.CreatedBy().String(), Valid: role.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.Model.UpdatedBy().String(), Valid: role.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.Model.CreatedAt(), Valid: !role.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.Model.UpdatedAt(), Valid: !role.Model.UpdatedAt().IsZero()},
	}
}

// ToRole converts a RoleDA data access object to a Role business object
func ToRole(da RoleDA) Role {
	return Role{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
		Status:      da.Status.String,
	}
}

// ToRoles converts a slice of RoleDA to a slice of Role business objects
func ToRoles(das []RoleDA) []Role {
	roles := make([]Role, len(das))
	for i, da := range das {
		roles[i] = ToRole(da)
	}
	return roles
}

// ToRoleExt converts RoleExtDA to Role including permissions
func ToRoleExt(da RoleExtDA) Role {
	return Role{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// ToPermissionDA converts a Permission business object to a PermissionDA data access object
func ToPermissionDA(permission Permission) PermissionDA {
	return PermissionDA{
		ID:          permission.ID(),
		Name:        sql.NullString{String: permission.Name, Valid: permission.Name != ""},
		Description: sql.NullString{String: permission.Description, Valid: permission.Description != ""},
		Slug:        sql.NullString{String: permission.Slug(), Valid: permission.Slug() != ""},
		CreatedBy:   sql.NullString{String: permission.Model.CreatedBy().String(), Valid: permission.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: permission.Model.UpdatedBy().String(), Valid: permission.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: permission.Model.CreatedAt(), Valid: !permission.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: permission.Model.UpdatedAt(), Valid: !permission.Model.UpdatedAt().IsZero()},
	}
}

// ToPermission converts a PermissionDA data access object to a Permission business object
func ToPermission(da PermissionDA) Permission {
	return Permission{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// ToPermissions converts a slice of PermissionDA to a slice of Permission business objects
func ToPermissions(das []PermissionDA) []Permission {
	permissions := make([]Permission, len(das))
	for i, da := range das {
		permissions[i] = ToPermission(da)
	}
	return permissions
}

// ToResourceDA converts a Resource business object to a ResourceDA data access object
func ToResourceDA(resource Resource) ResourceDA {
	return ResourceDA{
		ID:          resource.ID(),
		Name:        sql.NullString{String: resource.Name, Valid: resource.Name != ""},
		Description: sql.NullString{String: resource.Description, Valid: resource.Description != ""},
		Label:       sql.NullString{String: resource.Label, Valid: resource.Label != ""},
		Type:        sql.NullString{String: resource.Type, Valid: resource.Type != ""},
		URI:         sql.NullString{String: resource.URI, Valid: resource.URI != ""},
		Slug:        sql.NullString{String: resource.Slug(), Valid: resource.Slug() != ""},
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   sql.NullString{String: resource.Model.CreatedBy().String(), Valid: resource.Model.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: resource.Model.UpdatedBy().String(), Valid: resource.Model.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: resource.Model.CreatedAt(), Valid: !resource.Model.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: resource.Model.UpdatedAt(), Valid: !resource.Model.UpdatedAt().IsZero()},
	}
}

// ToResource converts a ResourceDA data access object to a Resource business object
func ToResource(da ResourceDA) Resource {
	return Resource{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
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

// ToResources converts a slice of ResourceDA to a slice of Resource business objects
func ToResources(das []ResourceDA) []Resource {
	resources := make([]Resource, len(das))
	for i, da := range das {
		resources[i] = ToResource(da)
	}
	return resources
}

// ToResourceExt converts ResourceExtDA to Resource including permissions
func ToResourceExt(da ResourceExtDA) Resource {
	return Resource{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
	}
}

// Helper functions for role conversions
func toRoles(roleIDs []uuid.UUID) []Role {
	var roles []Role
	for _, id := range roleIDs {
		roles = append(roles, Role{
			Model: am.NewModel(am.WithID(id)),
		})
	}
	return roles
}

func toRoleIDs(roles []Role) []uuid.UUID {
	var roleIDs []uuid.UUID
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID())
	}
	return roleIDs
}

// Helper functions for permission conversions
func toPermissions(permissionIDs []uuid.UUID) []Permission {
	var permissions []Permission
	for _, id := range permissionIDs {
		permissions = append(permissions, Permission{
			Model: am.NewModel(am.WithID(id)),
		})
	}
	return permissions
}

func toPermissionIDs(permissions []Permission) []uuid.UUID {
	var ids []uuid.UUID
	for _, permission := range permissions {
		ids = append(ids, permission.ID())
	}
	return ids
}
