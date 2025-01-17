package todo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
)

var (
	todoResPath = "/todo"
	key         = am.Key
	method      = am.HTTPMethod
)

type WebHandler struct {
	core    *am.Handler
	service Service
	tm      *am.TemplateManager
}

func NewWebHandler(tm *am.TemplateManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		core:    handler,
		service: service,
		tm:      tm,
	}
}

func (h *WebHandler) ListLists(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of lists")
	ctx := r.Context()

	lists, err := h.service.GetLists(ctx)
	if err != nil {
		http.Error(w, am.ErrCannotGetResources, http.StatusInternalServerError)
		return
	}

	page := am.NewPage(lists)
	page.SetFormAction(todoResPath)
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("todo", "list")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) ShowList(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Show todo list ", slug)
	ctx := r.Context()

	list, err := h.service.GetListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)
	blue, _ := cfg.StrVal(key.ButtonStyleBlue)
	red, _ := cfg.StrVal(key.ButtonStyleRed)

	page := am.NewPage(list)
	page.SetActions([]am.Action{
		{URL: todoResPath, Text: "Back to List", Style: gray},
		{URL: fmt.Sprintf("%s/%s/edit", todoResPath, slug), Text: "Edit List", Style: blue},
		{URL: fmt.Sprintf("%s/%s/delete", todoResPath, slug), Text: "Delete List", Style: red},
	})

	tmpl, err := h.tm.Get("todo", "show")
	if err != nil {
		http.Error(w, am.ErrTemplateNotFound, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, am.ErrCannotRenderTemplate, http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, am.ErrCannotWriteResponse, http.StatusInternalServerError)
	}
}

func (h *WebHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create todo list")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	list := NewList(name, description)

	err := h.service.CreateList(ctx, list)
	if err != nil {
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) EditList(w http.ResponseWriter, r *http.Request) {
	slug := r.FormValue("slug")
	h.Log().Info("Edit todo list ", slug)
	ctx := r.Context()

	list, err := h.service.GetListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	list.Name = name
	list.Description = description

	err = h.service.UpdateList(ctx, list)
	if err != nil {
		http.Error(w, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	slug := r.FormValue("slug")
	h.Log().Info("Delete todo list ", slug)
	ctx := r.Context()

	err := h.service.DeleteListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Add item to list")
	ctx := r.Context()

	listSlug := r.FormValue("list_slug")
	description := r.FormValue("description")
	status := r.FormValue("status")
	item := NewItem(description, status)

	err := h.service.AddItem(ctx, listSlug, item)
	if err != nil {
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, listSlug), http.StatusSeeOther)
}

func (h *WebHandler) EditItem(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Edit item in list")
	ctx := r.Context()

	listSlug := r.FormValue("list_slug")
	itemID := r.FormValue("item_id")
	description := r.FormValue("description")
	status := r.FormValue("status")

	err := h.service.EditItem(ctx, listSlug, itemID, description, status)
	if err != nil {
		http.Error(w, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, listSlug), http.StatusSeeOther)
}

func (h *WebHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Delete item from list")
	ctx := r.Context()

	listSlug := r.FormValue("list_slug")
	itemID := r.FormValue("item_id")

	err := h.service.DeleteItem(ctx, listSlug, itemID)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%s/%s", todoResPath, listSlug), http.StatusSeeOther)
}

// Name returns the name in WebHandler.
func (h *WebHandler) Name() string {
	return h.core.Name()
}

// SetName sets the name in WebHandler.
func (h *WebHandler) SetName(name string) {
	h.core.SetName(name)
}

// Log returns the Logger in WebHandler.
func (h *WebHandler) Log() am.Logger {
	return h.core.Log()
}

// SetLog sets the Logger in WebHandler.
func (h *WebHandler) SetLog(log am.Logger) {
	h.core.SetLog(log)
}

// Cfg returns the Config in WebHandler.
func (h *WebHandler) Cfg() *am.Config {
	return h.core.Cfg()
}

// SetCfg sets the Config in WebHandler.
func (h *WebHandler) SetCfg(cfg *am.Config) {
	h.core.SetCfg(cfg)
}

// Setup is the default implementation for the Setup method in WebHandler.
func (h *WebHandler) Setup(ctx context.Context) error {
	return h.core.Setup(ctx)
}

// Start is the default implementation for the Start method in WebHandler.
func (h *WebHandler) Start(ctx context.Context) error {
	return h.core.Start(ctx)
}

// Stop is the default implementation for the Stop method in WebHandler.
func (h *WebHandler) Stop(ctx context.Context) error {
	return h.core.Stop(ctx)
}
