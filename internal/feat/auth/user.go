package auth

import (
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	userType = "user"
)

type User struct {
	am.Model

	Username    string     `json:"username"`
	Email       string     `json:"email"`
	EmailEnc    []byte     `json:"-"`
	Name        string     `json:"name"`
	Password    string     `json:"password"`
	PasswordEnc []byte     `json:"-"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP string     `json:"last_login_ip,omitempty"`

	RoleIDs       []uuid.UUID `json:"-"`
	PermissionIDs []uuid.UUID `json:"-"`

	Roles       []Role       `json:"roles,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

// NewUser creates a user with default values.
func NewUser(username, name string) User {
	return User{
		Model:         am.NewModel(am.WithType(userType)),
		Username:      username,
		Name:          name,
		IsActive:      true,
		Roles:         []Role{},
		Permissions:   []Permission{},
		RoleIDs:       []uuid.UUID{},
		PermissionIDs: []uuid.UUID{},
	}
}

// NewUserSec creates a new user and encrypts email/password.
func NewUserSec(username, email, password, name string, emailKey []byte) (User, error) {
	emailEnc, err := EncryptEmail(email, emailKey)
	if err != nil {
		return User{}, err
	}

	passwordEnc, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	u := NewUser(username, name)
	u.SetEmailEnc(emailEnc)
	u.SetPasswordEnc(passwordEnc)

	return u, nil
}

// SetEmailEnc sets the encrypted email for the user.
func (u *User) SetEmailEnc(emailEnc []byte) {
	u.EmailEnc = emailEnc
}

// SetPasswordEnc sets the encrypted password for the user.
func (u *User) SetPasswordEnc(passwordEnc []byte) {
	u.PasswordEnc = passwordEnc
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
