package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
)

const (
	permissionType = "permission"
)

type Permission struct {
	am.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewPermission(name, description string) Permission {
	return Permission{
		Model:       am.NewModel(am.WithType(permissionType)),
		Name:        name,
		Description: description,
	}
}
