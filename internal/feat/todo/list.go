package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	listType = "list"
)

type List struct {
	Model       am.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewList creates a new list.
func NewList(name, description string) List {
	return List{
		Model:       am.NewModel(am.WithType(listType)),
		Name:        name,
		Description: description,
	}
}

// ID returns the unique identifier of the list.
func (l *List) ID() uuid.UUID {
	return l.Model.ID()
}

// SetID sets the unique identifier of the list.
func (l *List) SetID(id uuid.UUID, force ...bool) {
	l.Model.SetID(id, force...)
}

// Slug returns the slug of the list.
func (l *List) Slug() string {
	return l.Model.Slug()
}

// GenSlug generates and sets the slug of the list.
func (l *List) GenSlug() {
	l.Model.GenSlug(l.Name)
}

func (l *List) GetNameID() {
	l.Model.GenNameID()
}

func (l *List) SetCreateValues() {
	l.Model.GenCreationValues()
}
