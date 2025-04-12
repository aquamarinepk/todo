package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	resourceEntityType = "resource"
)

type Resource struct {
	am.Model
	Name          string `json:"name"`
	Description   string `json:"description"`
	Label         string `json:"label"`
	ResourceType  string `json:"type"` // Type of resource (e.g., "url", "entity")
	URI           string `json:"uri"`
	PermissionIDs []uuid.UUID
	Permissions   []Permission
}

func NewResource(name, description, resourceType string) Resource {
	model := am.NewModel(am.WithType(resourceEntityType))
	model.GenCreationValues()
	return Resource{
		Model:         model,
		Name:          name,
		Description:   description,
		ResourceType:  resourceType,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
}
