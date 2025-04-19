package auth

import (
	"database/sql"
	"time"

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
		EmailEnc:      user.EmailEnc,
		PasswordEnc:   user.PasswordEnc,
		RoleIDs:       toRoleIDs(user.Roles),
		PermissionIDs: toPermissionIDs(user.Permissions),
		CreatedBy:     sql.NullString{String: user.CreatedBy().String(), Valid: user.CreatedBy() != uuid.Nil},
		UpdatedBy:     sql.NullString{String: user.UpdatedBy().String(), Valid: user.UpdatedBy() != uuid.Nil},
		CreatedAt:     sql.NullTime{Time: user.CreatedAt(), Valid: !user.CreatedAt().IsZero()},
		UpdatedAt:     sql.NullTime{Time: user.UpdatedAt(), Valid: !user.UpdatedAt().IsZero()},
		LastLoginAt:   sql.NullTime{Time: derefTime(user.LastLoginAt), Valid: user.LastLoginAt != nil},
		LastLoginIP:   sql.NullString{String: user.LastLoginIP, Valid: user.LastLoginIP != ""},
		IsActive:      sql.NullBool{Bool: user.IsActive, Valid: true},
	}
}

// ToUser converts a UserDA data access object to a User business object
func ToUser(da UserDA) User {
	var lastLoginAt *time.Time
	if da.LastLoginAt.Valid {
		t := da.LastLoginAt.Time
		lastLoginAt = &t
	}

	return User{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(userType),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Username:    da.Username.String,
		EmailEnc:    da.EmailEnc,
		PasswordEnc: da.PasswordEnc,
		LastLoginAt: lastLoginAt,
		LastLoginIP: da.LastLoginIP.String,
		IsActive:    da.IsActive.Bool,
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
	user := ToUser(UserDA{
		ID:          da.ID,
		Slug:        da.Slug,
		Name:        da.Name,
		Username:    da.Username,
		EmailEnc:    da.EmailEnc,
		PasswordEnc: da.PasswordEnc,
		CreatedBy:   da.CreatedBy,
		UpdatedBy:   da.UpdatedBy,
		CreatedAt:   da.CreatedAt,
		UpdatedAt:   da.UpdatedAt,
		LastLoginAt: da.LastLoginAt,
		LastLoginIP: da.LastLoginIP,
		IsActive:    da.IsActive,
	})

	// Add role if present
	if da.RoleID.Valid {
		user.RoleIDs = append(user.RoleIDs, am.ParseUUID(da.RoleID))
		user.Roles = append(user.Roles, Role{
			Model: am.NewModel(
				am.WithID(am.ParseUUID(da.RoleID)),
				am.WithType(roleType),
			),
			Name: da.RoleName.String,
		})
	}

	// Add permission if present
	if da.PermissionID.Valid {
		user.PermissionIDs = append(user.PermissionIDs, am.ParseUUID(da.PermissionID))
		user.Permissions = append(user.Permissions, Permission{
			Model: am.NewModel(
				am.WithID(am.ParseUUID(da.PermissionID)),
				am.WithType(permissionType),
			),
			Name: da.PermissionName.String,
		})
	}

	return user
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
		CreatedBy:   sql.NullString{String: role.CreatedBy().String(), Valid: role.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: role.UpdatedBy().String(), Valid: role.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: role.CreatedAt(), Valid: !role.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: role.UpdatedAt(), Valid: !role.UpdatedAt().IsZero()},
	}
}

// ToRole converts a RoleDA data access object to a Role business object
func ToRole(da RoleDA) Role {
	return Role{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(roleType),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		Status:        da.Status.String,
		PermissionIDs: da.Permissions,
		Permissions:   []Permission{},
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
	permission := Permission{
		Model: am.NewModel(
			am.WithID(am.ParseUUID(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: da.PermissionName.String,
	}

	return Role{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(roleType),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:        da.Name.String,
		Description: da.Description.String,
		Status:      "active", // Default status since it's not in RoleExtDA
		Permissions: []Permission{permission},
	}
}

// ToPermissionDA converts a Permission business object to a PermissionDA data access object
func ToPermissionDA(permission Permission) PermissionDA {
	return PermissionDA{
		ID:          permission.ID(),
		Name:        sql.NullString{String: permission.Name, Valid: permission.Name != ""},
		Description: sql.NullString{String: permission.Description, Valid: permission.Description != ""},
		Slug:        sql.NullString{String: permission.Slug(), Valid: permission.Slug() != ""},
		CreatedBy:   sql.NullString{String: permission.CreatedBy().String(), Valid: permission.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: permission.UpdatedBy().String(), Valid: permission.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: permission.CreatedAt(), Valid: !permission.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: permission.UpdatedAt(), Valid: !permission.UpdatedAt().IsZero()},
	}
}

// ToPermission converts a PermissionDA data access object to a Permission business object
func ToPermission(da PermissionDA) Permission {
	return Permission{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(permissionType),
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
		Type:        sql.NullString{String: resource.ResourceType, Valid: resource.ResourceType != ""},
		URI:         sql.NullString{String: resource.URI, Valid: resource.URI != ""},
		Slug:        sql.NullString{String: resource.Slug(), Valid: resource.Slug() != ""},
		Permissions: toPermissionIDs(resource.Permissions),
		CreatedBy:   sql.NullString{String: resource.CreatedBy().String(), Valid: resource.CreatedBy() != uuid.Nil},
		UpdatedBy:   sql.NullString{String: resource.UpdatedBy().String(), Valid: resource.UpdatedBy() != uuid.Nil},
		CreatedAt:   sql.NullTime{Time: resource.CreatedAt(), Valid: !resource.CreatedAt().IsZero()},
		UpdatedAt:   sql.NullTime{Time: resource.UpdatedAt(), Valid: !resource.UpdatedAt().IsZero()},
	}
}

// ToResource converts a ResourceDA data access object to a Resource business object
func ToResource(da ResourceDA) Resource {
	return Resource{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(resourceEntityType),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		Label:         da.Label.String,
		ResourceType:  da.Type.String,
		URI:           da.URI.String,
		PermissionIDs: da.Permissions,
		Permissions:   []Permission{},
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
	permission := Permission{
		Model: am.NewModel(
			am.WithID(am.ParseUUID(da.PermissionID)),
			am.WithType(permissionType),
		),
		Name: da.PermissionName.String,
	}

	return Resource{
		Model: am.NewModel(
			am.WithID(da.ID),
			am.WithType(resourceEntityType),
			am.WithSlug(da.Slug.String),
			am.WithCreatedBy(am.ParseUUID(da.CreatedBy)),
			am.WithUpdatedBy(am.ParseUUID(da.UpdatedBy)),
			am.WithCreatedAt(da.CreatedAt.Time),
			am.WithUpdatedAt(da.UpdatedAt.Time),
		),
		Name:          da.Name.String,
		Description:   da.Description.String,
		ResourceType:  "entity",
		PermissionIDs: []uuid.UUID{am.ParseUUID(da.PermissionID)},
		Permissions:   []Permission{permission},
	}
}

func toRoleIDs(roles []Role) []uuid.UUID {
	var ids []uuid.UUID
	for _, r := range roles {
		ids = append(ids, r.ID())
	}
	return ids
}

func toPermissionIDs(perms []Permission) []uuid.UUID {
	var ids []uuid.UUID
	for _, p := range perms {
		ids = append(ids, p.ID())
	}
	return ids
}

func derefTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
