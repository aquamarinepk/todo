package am

import (
	"fmt"

	"github.com/google/uuid"
)

// ListPath returns the path for listing resources
func ListPath(basePath, resourceType string) string {
	return fmt.Sprintf("%s/list-%ss", basePath, resourceType)
}

// NewPath returns the path for creating a new resource
func NewPath(basePath, resourceType string) string {
	return fmt.Sprintf("%s/new-%s", basePath, resourceType)
}

// CreatePath returns the path for creating a resource
func CreatePath(basePath, resourceType string) string {
	return fmt.Sprintf("%s/create-%s", basePath, resourceType)
}

// ShowPath returns the path for showing a resource
func ShowPath(basePath, resourceType string, id uuid.UUID) string {
	return fmt.Sprintf("%s/%ss/%s", basePath, resourceType, id)
}

// EditPath returns the path for editing a resource
func EditPath(basePath, resourceType string, id uuid.UUID) string {
	return fmt.Sprintf("%s/edit-%s?id=%s", basePath, resourceType, id)
}

// UpdatePath returns the path for updating a resource
func UpdatePath(basePath, resourceType string) string {
	return fmt.Sprintf("%s/update-%s", basePath, resourceType)
}

// DeletePath returns the path for deleting a resource
func DeletePath(basePath, resourceType string) string {
	return fmt.Sprintf("%s/delete-%s", basePath, resourceType)
}

// ListRelatedPath returns the path for listing related resources
func ListRelatedPath(basePath, resourceType, relatedType string, id uuid.UUID) string {
	return fmt.Sprintf("%s/list-%s-%ss?id=%s", basePath, resourceType, relatedType, id)
}

// AddRelatedPath returns the path for adding a related resource
func AddRelatedPath(basePath, resourceType, relatedType string) string {
	return fmt.Sprintf("%s/add-%s-to-%s", basePath, relatedType, resourceType)
}

// RemoveRelatedPath returns the path for removing a related resource
func RemoveRelatedPath(basePath, resourceType, relatedType string) string {
	return fmt.Sprintf("%s/remove-%s-from-%s", basePath, relatedType, resourceType)
}
