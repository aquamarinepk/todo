package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	orgEntityType = "org"
)

type Org struct {
	am.Model
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	OwnerID          uuid.UUID `json:"owner_id"`
}

func NewOrg(name, shortDescription, description string, ownerID uuid.UUID) Org {
	model := am.NewModel(am.WithType(orgEntityType))
	model.GenCreationValues()
	return Org{
		Model:            model,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		OwnerID:          ownerID,
	}
}
