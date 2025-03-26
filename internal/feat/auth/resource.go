package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	resourceType = "resource"
)

type Resource struct {
	am.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
	PermissionIDs []uuid.UUID
	Permissions []Permission
}

func NewResource(name, description string) Resource {
	return Resource{
		Model:       am.NewModel(am.WithType(resourceType)),
		Name:        name,
		Description: description,
		PermissionIDs: []uuid.UUID{},
		Permissions: []Permission{},
	}
}
