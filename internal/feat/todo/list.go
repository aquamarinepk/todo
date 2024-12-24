package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

type List struct {
	Model       am.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewList creates a new list.
func NewList(name, description string) List {
	return List{
		Model:       am.NewModel(),
		Name:        name,
		Description: description,
	}
}

// ID returns the unique identifier of the list.
func (l *List) ID() uuid.UUID {
	return l.Model.ID()
}

// SetID sets the unique identifier of the list.
func (l *List) SetID(id uuid.UUID) {
	l.Model.GenID(id)
}

// Slug returns the slug of the list.
func (l *List) Slug() string {
	return l.Model.Slug()
}

// SetSlug generates and sets the slug of the list.
func (l *List) SetSlug(slug string) {
	l.Model.GenSlug(slug)
}
