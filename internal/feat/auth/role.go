package auth

import (
	"encoding/json"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	roleType = "role"
)

// Role represents a todo user role.
type Role struct {
	*am.BaseModel
	UserID        uuid.UUID
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string
	PermissionIDs []uuid.UUID
	Permissions   []Permission
}

// NewRole creates a new Role.
func NewRole(name, description, status string) Role {
	return Role{
		BaseModel:     am.NewModel(am.WithType(roleType)),
		Name:          name,
		Description:   description,
		Status:        status,
		PermissionIDs: []uuid.UUID{},
		Permissions:   []Permission{},
	}
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (r *Role) UnmarshalJSON(data []byte) error {
	type Alias Role
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*r = Role(*temp)
	if r.BaseModel == nil {
		r.BaseModel = am.NewModel(am.WithType(roleType))
	}
	return nil
}

func (r *Role) Slug() string {
	return am.Normalize(r.Name) + "-" + r.ShortID()
}
