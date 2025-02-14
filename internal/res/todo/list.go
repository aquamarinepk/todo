package todo

import (
	"github.com/aquamarinepk/todo/internal/am"
)

const (
	listType = "list"
)

type List struct {
	am.Model
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
