package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Role represents a todo user role.
type Role struct {
	Model       am.Model
	UserID      uuid.UUID
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

// ID returns the ID of the role.
func (i *Role) ID() uuid.UUID {
	return i.Model.ID()
}

// SetCreateValues sets the creation values for the role.
func (i *Role) SetCreateValues() {
	i.Model.GenCreationValues()
}
