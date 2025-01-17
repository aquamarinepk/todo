package todo

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
}

// NewUser creates a new user.
func NewUser(username, email string) User {
	return User{
		Model:    am.NewModel(am.WithType(userType)),
		Username: username,
		Email:    email,
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
