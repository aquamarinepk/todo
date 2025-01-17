package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

// Item represents a todo list item.
type Item struct {
	Model       am.Model
	ListID      uuid.UUID
	Description string
	Status      string
}

// NewItem creates a new Item.
func NewItem(description, status string) Item {
	return Item{
		Model:       am.NewModel(),
		Description: description,
		Status:      status,
	}
}

// ID returns the ID of the item.
func (i *Item) ID() uuid.UUID {
	return i.Model.ID()
}

// SetCreateValues sets the creation values for the item.
func (i *Item) SetCreateValues() {
	i.Model.GenCreationValues()
}
