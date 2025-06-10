package auth

import (
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/google/uuid"
)

const (
	authPath    = "/auth"
	userPathFmt = "%s/%s-user%s"
)

const (
	ActionListUserRoles       = "list-user-roles"
	ActionListUserPermissions = "list-user-permissions"
	ActionListTeamMembers     = "list-team-members"
	TextRoles                 = "Roles"
	TextPermissions           = "Permissions"
	TextMembers               = "Members"
)

type WebHandler struct {
	*am.Handler
	service Service
	tm      *am.TemplateManager
	crypto  *am.Crypto
	flash   *am.FlashManager
}

func NewWebHandler(tm *am.TemplateManager, flash *am.FlashManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		service: service,
		tm:      tm,
		crypto:  &am.Crypto{},
		flash:   flash,
	}
}

// WIP: The following flash-related functions were part of the initial experimentation
// phase to test ideas and concepts. They will be deprecated once the am.FlashManager
// is fully functional and integrated. These functions served as a proof of concept
// for the flash message system that is now being properly implemented in the am package.
// After the am.FlashManager is fully functional, these functions will be removed.

func (h *WebHandler) AddFlash(w http.ResponseWriter, r *http.Request, notificationType string, msg string) error {
	flash := GetFlash(r)
	flash.Add(notificationType, msg)
	return h.flash.SetFlash(w, flash)
}

// WIP:
func (h *WebHandler) GetFlash(r *http.Request) am.Flash {
	return GetFlash(r)
}

func (h *WebHandler) AddInfoFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	return h.AddFlash(w, r, am.NotificationType.Info, msg)
}

func (h *WebHandler) AddSuccessFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	return h.AddFlash(w, r, am.NotificationType.Success, msg)
}

func (h *WebHandler) AddWarningFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	return h.AddFlash(w, r, am.NotificationType.Warn, msg)
}

func (h *WebHandler) AddErrorFlash(w http.ResponseWriter, r *http.Request, msg string) error {
	return h.AddFlash(w, r, am.NotificationType.Error, msg)
}

func (h *WebHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid user ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get user ", id)
	ctx := r.Context()

	if _, err := h.service.GetUser(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid role ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get role ", id)
	ctx := r.Context()

	if _, err := h.service.GetRole(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid permission ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get permission ", id)
	ctx := r.Context()

	if _, err := h.service.GetPermission(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}

func (h *WebHandler) GetResource(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.Err(w, err, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.Log().Info("Get resource", id)
	ctx := r.Context()

	if _, err := h.service.GetResource(ctx, id); err != nil {
		h.Err(w, err, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}
}
