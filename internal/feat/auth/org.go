package auth

import (
	"encoding/json"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	orgEntityType = "org"
)

type Org struct {
	*am.BaseModel
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	OwnerID          uuid.UUID `json:"owner_id"`
}

func NewOrg(name, shortDescription, description string, ownerID uuid.UUID) Org {
	model := am.NewModel(am.WithType(orgEntityType))
	model.GenCreationValues()
	return Org{
		BaseModel:        model,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		OwnerID:          ownerID,
	}
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (o *Org) UnmarshalJSON(data []byte) error {
	type Alias Org
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*o = Org(*temp)
	if o.BaseModel == nil {
		o.BaseModel = am.NewModel(am.WithType(orgEntityType))
	}
	return nil
}
