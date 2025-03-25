package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	userType = "user"
)

type User struct {
	am.Model
	Username      string        `json:"username"`
	Email         string        `json:"email"`
	Name          string        `json:"name"`
	EncPassword   string
	RoleIDs       []uuid.UUID
	PermissionIDs []uuid.UUID
	Roles         []Role
	Permissions   []Permission
}

// NewUser creates a new user.
func NewUser(username, email, name string) User {
	return User{
		Model:         am.NewModel(am.WithType(userType)),
		Username:      username,
		Email:         email,
		Name:          name,
		Roles:         []Role{},
		Permissions:   []Permission{},
		RoleIDs:       []uuid.UUID{},
		PermissionIDs: []uuid.UUID{},
	}
}

// AddRole adds a role to the user.
func (l *User) AddRole(role Role) {
	l.Roles = append(l.Roles, role)
}

// RemoveRole removes a role from the user.
func (l *User) RemoveRole(roleID uuid.UUID) {
	for i, role := range l.Roles {
		if role.ID() == roleID {
			l.Roles = append(l.Roles[:i], l.Roles[i+1:]...)
			break
		}
	}
}
