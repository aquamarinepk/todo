package auth

import (
	"encoding/json"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	teamEntityType = "team"
)

type Team struct {
	*am.BaseModel
	OrgID            uuid.UUID `json:"org_id"`
	OrgRef           string    `json:"org_ref"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
}

func NewTeam(orgID uuid.UUID, name, shortDescription, description string) Team {
	model := am.NewModel(am.WithType(teamEntityType))
	model.GenCreationValues()
	return Team{
		BaseModel:        model,
		OrgID:            orgID,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
	}
}

// UnmarshalJSON ensures Model is always initialized after unmarshal.
func (t *Team) UnmarshalJSON(data []byte) error {
	type Alias Team
	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}
	*t = Team(*temp)
	if t.BaseModel == nil {
		t.BaseModel = am.NewModel(am.WithType(teamEntityType))
	}
	return nil
}
