package auth

import (
	"errors"
)

// FormToUser converts a UserForm to a User entity.
// It validates that the password and password confirmation match.
func FormToUser(form UserForm, encKey []byte) (User, error) {
	// TODO: A configurable validator will take care of this briefly
	if form.Password != form.PasswordConf {
		return User{}, errors.New("passwords do not match")
	}

	return NewUserSec(form.Username, form.Email, form.Password, form.Name, encKey)
}

// FormToRole converts a RoleForm to a Role entity.
func FormToRole(form RoleForm) Role {
	return NewRole(form.Name, form.Description, form.Status)
}

// FormToPermission converts a PermissionForm to a Permission entity.
func FormToPermission(form PermissionForm) Permission {
	return NewPermission(form.Name, form.Description)
}

// FormToResource converts a ResourceForm to a Resource entity.
func FormToResource(form ResourceForm) Resource {
	return NewResource(form.Name, form.Description, form.Type)
}
