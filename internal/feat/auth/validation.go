package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// ValidateUser validates a UserForm.
// It checks:
// - Username length (min 3, max 50)
// - Email format (basic check)
// - Password length (min 8)
// - Password confirmation
func ValidateUser(form UserForm) (am.Validation, error) {
	validate := am.ComposeValidators(
		am.MinLength("username", form.Username, 3),
		am.MaxLength("username", form.Username, 50),
		am.MinLength("email", form.Email, 5),
		am.MinLength("password", form.Password, 8),
		am.Equals("password", form.Password, form.PasswordConf),
	)

	return validate(form)
}

// ValidateRole validates a RoleForm.
// It checks:
// - Name length (min 3, max 50)
// - Description length (max 255)
func ValidateRole(form RoleForm) (am.Validation, error) {
	validate := am.ComposeValidators(
		am.MinLength("name", form.Name, 3),
		am.MaxLength("name", form.Name, 50),
		am.MaxLength("description", form.Description, 255),
	)

	return validate(form)
}

// ValidatePermission validates a PermissionForm.
// It checks:
// - Name length (min 3, max 50)
// - Description length (max 255)
func ValidatePermission(form PermissionForm) (am.Validation, error) {
	validate := am.ComposeValidators(
		am.MinLength("name", form.Name, 3),
		am.MaxLength("name", form.Name, 50),
		am.MaxLength("description", form.Description, 255),
	)

	return validate(form)
}

// ValidateResource validates a ResourceForm.
// It checks:
// - Name length (min 3, max 50)
// - Description length (max 255)
func ValidateResource(form ResourceForm) (am.Validation, error) {
	validate := am.ComposeValidators(
		am.MinLength("name", form.Name, 3),
		am.MaxLength("name", form.Name, 50),
		am.MaxLength("description", form.Description, 255),
	)

	return validate(form)
}
