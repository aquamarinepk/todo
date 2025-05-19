package auth

import (
	"encoding/json"

	"github.com/aquamarinepk/todo/internal/am"
)

const (
	permissionType = "permission"
)

type Permission struct {
	*am.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewPermission(name, description string) Permission {
	return Permission{
		BaseModel:   am.NewModel(am.WithType(permissionType)),
		Name:        name,
		Description: description,
	}
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (p *Permission) UnmarshalJSON(data []byte) error {
	type Alias Permission
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*p = Permission(*temp)
	if p.BaseModel == nil {
		p.BaseModel = am.NewModel(am.WithType(permissionType))
	}
	return nil
}
