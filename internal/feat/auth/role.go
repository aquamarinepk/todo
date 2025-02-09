package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Role represents a todo user role.
type Role struct {
	am.Model
	UserID      uuid.UUID
	UserSlug    string
	Name        string
	Description string
	Status      string
}

// NewRole creates a new Role.
func NewRole(name, description, status string) Role {
	return Role{
		Model:       am.NewModel(),
		Name:        name,
		Description: description,
		Status:      status,
	}
}
