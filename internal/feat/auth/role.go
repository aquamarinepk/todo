package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	roleType = "role"
)

// Role represents a todo user role.
type Role struct {
	am.Model
	UserID        uuid.UUID
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string
	PermissionIDs []uuid.UUID
	Permissions   []Permission
}

// NewRole creates a new Role.
func NewRole(name, description, status string) Role {
	return Role{
		Model:         am.NewModel(am.WithType(roleType)),
		Name:          name,
		Description:   description,
		Status:        status,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
}
