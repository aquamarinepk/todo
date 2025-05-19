package auth

import (
	"encoding/json"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	resourceEntityType = "resource"
)

type Resource struct {
	*am.BaseModel
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
		BaseModel:     model,
		Name:          name,
		Description:   description,
		ResourceType:  resourceType,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (r *Resource) UnmarshalJSON(data []byte) error {
	type Alias Resource
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*r = Resource(*temp)
	if r.BaseModel == nil {
		r.BaseModel = am.NewModel(am.WithType(resourceEntityType))
	}
	return nil
}
