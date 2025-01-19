package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Role represents a todo user role.
type Role struct {
	Model       am.Model
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

// ID returns the unique identifier of the role.
func (r *Role) ID() uuid.UUID {
	return r.Model.ID()
}

// SetID sets the unique identifier of the role.
func (r *Role) SetID(id uuid.UUID, force ...bool) {
	r.Model.SetID(id, force...)
}

// Slug returns the slug of the role.
func (r *Role) Slug() string {
	return r.Model.Slug()
}

// GenSlug generates and sets the slug of the role.
func (r *Role) GenSlug() {
	r.Model.GenSlug(r.Name)
}

// GetNameID generates and sets the name ID of the role.
func (r *Role) GetNameID() {
	r.Model.GenNameID()
}

// SetCreateValues sets the creation values for the role.
func (r *Role) SetCreateValues() {
	r.Model.GenCreationValues()
}
