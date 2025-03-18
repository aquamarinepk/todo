package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
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
	Permissions []Permission
}

func NewResource(name, description, label, resourceType, uri string) Resource {
	return Resource{
		Model:       am.NewModel(am.WithType(resourceType)),
		Name:        name,
		Description: description,
		Label:       label,
		Type:        resourceType,
		URI:         uri,
		Permissions: []Permission{},
	}
}
