package am

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	Core
}

func NewHandler(name string, opts ...Option) *Handler {
	core := NewCore(name, opts...)
	return &Handler{
		Core: core,
	}
}

// Err writes an error response to the client.
// This method can be overridden by types that embed Handler.
func (h *Handler) Err(w http.ResponseWriter, err error, msg string, code int) {
	h.Log().Error(msg, err)
	http.Error(w, msg, code)
}

// Render writes a template to the response writer.
// This method can be overridden by types that embed Handler.
func (h *Handler) Render(w http.ResponseWriter, r *http.Request, template string, data any) error {
	return nil // Default implementation does nothing
}

// ID parses a resource ID from the URL query parameters.
// Returns the parsed UUID and true if successful, or uuid.Nil and false if not.
func (h *Handler) ID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	return h.ParseUUIDFromQuery(w, r, "id")
}

// UserID parses a user ID from the URL query parameters.
// Returns the parsed UUID and true if successful, or uuid.Nil and false if not.
func (h *Handler) UserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	return h.ParseUUIDFromQuery(w, r, "id")
}

// RoleID parses a role ID from the URL query parameters.
// Returns the parsed UUID and true if successful, or uuid.Nil and false if not.
func (h *Handler) RoleID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	return h.ParseUUIDFromQuery(w, r, "id")
}

// PermissionID parses a permission ID from the URL query parameters.
// Returns the parsed UUID and true if successful, or uuid.Nil and false if not.
func (h *Handler) PermissionID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	return h.ParseUUIDFromQuery(w, r, "id")
}

// ShowItem is a helper method for handlers that show a single item.
// It parses the ID from the request, gets the item using the provided getter function,
// and renders it using the provided template.
// If any step fails, it writes an error response and returns false.
func (h *Handler) ShowItem(w http.ResponseWriter, r *http.Request, getter func(uuid.UUID) (any, error), template string) bool {
	id, ok := h.ID(w, r)
	if !ok {
		return false
	}

	item, err := getter(id)
	if err != nil {
		h.Err(w, err, "Failed to get item", http.StatusInternalServerError)
		return false
	}

	if err := h.Render(w, r, template, item); err != nil {
		h.Err(w, err, "Failed to render template", http.StatusInternalServerError)
		return false
	}

	return true
}

// ParseUUIDFromQuery parses a UUID from the URL query parameters.
// If the parameter is not found or is invalid, it writes an error response and returns false.
// Returns the parsed UUID and true if successful.
func (h *Handler) ParseUUIDFromQuery(w http.ResponseWriter, r *http.Request, paramName string) (uuid.UUID, bool) {
	idStr := r.URL.Query().Get(paramName)
	if idStr == "" {
		h.Err(w, nil, "Missing "+paramName, http.StatusBadRequest)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid "+paramName, http.StatusBadRequest)
		return uuid.Nil, false
	}

	return id, true
}

// ParseUUIDsFromQuery parses multiple UUIDs from the URL query parameters.
// The UUIDs should be comma-separated in the query parameter.
// If any UUID is invalid, it writes an error response and returns false.
// Returns the parsed UUIDs and true if successful.
func (h *Handler) ParseUUIDsFromQuery(w http.ResponseWriter, r *http.Request, paramName string) ([]uuid.UUID, bool) {
	idsStr := r.URL.Query().Get(paramName)
	if idsStr == "" {
		h.Err(w, nil, "Missing "+paramName, http.StatusBadRequest)
		return nil, false
	}

	ids := strings.Split(idsStr, ",")
	uuids := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(strings.TrimSpace(idStr))
		if err != nil {
			h.Err(w, err, "Invalid "+paramName, http.StatusBadRequest)
			return nil, false
		}
		uuids = append(uuids, id)
	}

	return uuids, true
}
