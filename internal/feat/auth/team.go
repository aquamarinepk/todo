package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	teamEntityType = "team"
)

type Team struct {
	am.Model
	OrgID            uuid.UUID `json:"org_id"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
}

func NewTeam(orgID uuid.UUID, name, shortDescription, description string) Team {
	model := am.NewModel(am.WithType(teamEntityType))
	model.GenCreationValues()
	return Team{
		Model:            model,
		OrgID:            orgID,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
	}
}
