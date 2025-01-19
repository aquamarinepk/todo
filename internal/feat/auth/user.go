package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	userType = "user"
)

type User struct {
	Model       am.Model
	Username    string `json:"username"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	EncPassword string
	Roles       []Role `json:"roles"`
}

// NewUser creates a new user.
func NewUser(username, email, name string) User {
	return User{
		Model:    am.NewModel(am.WithType(userType)),
		Username: username,
		Email:    email,
		Name:     name,
	}
}

// AddRole adds a role to the user.
func (l *User) AddRole(role Role) {
	l.Roles = append(l.Roles, role)
}

// RemoveRole removes a role from the user.
func (l *User) RemoveRole(roleID uuid.UUID) {
	for i, role := range l.Roles {
		if role.ID() == roleID {
			l.Roles = append(l.Roles[:i], l.Roles[i+1:]...)
			break
		}
	}
}

// ID returns the unique identifier of the user.
func (l *User) ID() uuid.UUID {
	return l.Model.ID()
}

// SetID sets the unique identifier of the user.
func (l *User) SetID(id uuid.UUID, force ...bool) {
	l.Model.SetID(id, force...)
}

// Slug returns the slug of the user.
func (l *User) Slug() string {
	return l.Model.Slug()
}

// GenSlug generates and sets the slug of the user.
func (l *User) GenSlug() {
	l.Model.GenSlug(l.Username)
}

func (l *User) GetNameID() {
	l.Model.GenNameID()
}

func (l *User) SetCreateValues() {
	l.Model.GenCreationValues()
}
