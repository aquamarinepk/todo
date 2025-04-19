package auth

// Error messages specific to auth domain
const (
	ErrInvalidUserID          = "Invalid user ID"
	ErrInvalidRoleID          = "Invalid role ID"
	ErrInvalidPermissionID    = "Invalid permission ID"
	ErrInvalidResourceID      = "Invalid resource ID"
	ErrInvalidFormData        = "Invalid form data"
	ErrValidationFailed       = "Validation failed"
	ErrCannotCreateUser       = "Failed to create user"
	ErrCannotUpdateUser       = "Failed to update user"
	ErrCannotDeleteUser       = "Failed to delete user"
	ErrCannotCreateRole       = "Failed to create role"
	ErrCannotUpdateRole       = "Failed to update role"
	ErrCannotDeleteRole       = "Failed to delete role"
	ErrCannotCreatePermission = "Failed to create permission"
	ErrCannotUpdatePermission = "Failed to update permission"
	ErrCannotDeletePermission = "Failed to delete permission"
	ErrCannotCreateResource   = "Failed to create resource"
	ErrCannotUpdateResource   = "Failed to update resource"
	ErrCannotDeleteResource   = "Failed to delete resource"
)
