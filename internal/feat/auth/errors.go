package auth

import "errors"

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrResourceNotFound   = errors.New("resource not found")
)
