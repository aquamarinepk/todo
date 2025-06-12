package am

import (
	"fmt"
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
	errMsg := fmt.Sprintf("%s: '%v'", msg, err)
	h.Log().Error(errMsg)
	http.Error(w, errMsg, code)
}

// Render writes a template to the response writer.
// This method can be overridden by types that embed Handler.
func (h *Handler) Render(w http.ResponseWriter, r *http.Request, template string, data any) error {
	return nil // Default implementation does nothing
}

// Redir redirects to the specified path with http.StatusSeeOther
func (h *Handler) Redir(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, path, http.StatusSeeOther)
}

// ID returns the ID from the request query parameters.
// If the ID is missing or invalid, it writes an error response and returns an error.
func (h *Handler) ID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	id, err := h.ParseUUIDFromQuery(r, "id")
	if err != nil {
		h.Err(w, err, "Invalid ID", http.StatusBadRequest)
		return uuid.Nil, err
	}
	return id, nil
}

// ShowItem is a helper method that shows an item by ID.
// It handles ID parsing, item retrieval, and rendering.
// If any step fails, it writes an error response and returns false.
func (h *Handler) ShowItem(w http.ResponseWriter, r *http.Request, getter func(uuid.UUID) (any, error), template string) bool {
	id, err := h.ID(w, r)
	if err != nil {
		h.Err(w, err, "Failed to parse ID", http.StatusBadRequest)
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
// Returns the parsed UUID and an error if the ID is missing or invalid.
func (h *Handler) ParseUUIDFromQuery(r *http.Request, paramName string) (uuid.UUID, error) {
	idStr := r.URL.Query().Get(paramName)
	if idStr == "" {
		return uuid.Nil, fmt.Errorf("missing %s", paramName)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %w", paramName, err)
	}

	return id, nil
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
