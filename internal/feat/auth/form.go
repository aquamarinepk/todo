package auth

// UserForm represents the form data for creating/updating a user
type UserForm struct {
	Username     string `form:"username" required:"true"`
	Email        string `form:"email" required:"true"`
	Name         string `form:"name" required:"true"`
	Password     string `form:"password" required:"true"`
	PasswordConf string `form:"password_conf" required:"true"`
}

// RoleForm represents the form data for creating/updating a role
type RoleForm struct {
	Name        string `form:"name" required:"true"`
	Description string `form:"description"`
	Status      string `form:"status"`
}

// PermissionForm represents the form data for creating/updating a permission
type PermissionForm struct {
	Name        string `form:"name" required:"true"`
	Description string `form:"description"`
}

// ResourceForm represents the form data for creating/updating a resource
type ResourceForm struct {
	Name        string `form:"name" required:"true"`
	Description string `form:"description"`
	Type        string `form:"type" required:"true"`
	URI         string `form:"uri"`
}
