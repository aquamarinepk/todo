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
	key = am.Key
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

func (h *WebHandler) List(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("List of lists")
	ctx := r.Context()

	lists, err := h.service.GetLists(ctx)
	if err != nil {
		http.Error(w, "Failed to get lists", http.StatusInternalServerError)
		return
	}

	page := am.NewPage(lists)
	page.SetFormAction("/todo")

	tmpl, err := h.tm.Get("todo", "list")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *WebHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New todo form")

	page := am.NewPage(List{})
	page.SetFormAction("/todo")

	tmpl, err := h.tm.Get("todo", "new")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *WebHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create todo")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	list := NewList(name, description)

	err := h.service.CreateList(ctx, list)
	if err != nil {
		http.Error(w, "Failed to create list", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (h *WebHandler) Show(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Show todo ", slug)
	ctx := r.Context()

	list, err := h.service.GetListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	cfg := h.Cfg()
	gray, _ := cfg.StrVal(key.ButtonStyleGray)
	blue, _ := cfg.StrVal(key.ButtonStyleBlue)
	red, _ := cfg.StrVal(key.ButtonStyleRed)

	page := am.NewPage(list)
	page.SetActions([]am.Action{ // NOTE: This is a WIP, it will be improved.
		{URL: "/todo", Text: "Back to List", Style: gray, IsForm: false},
		{URL: fmt.Sprintf("/todo/%s/edit", slug), Text: "Edit", Style: blue, IsForm: false},
		{URL: fmt.Sprintf("/todo/%s/delete", slug), Text: "Delete", Style: red, IsForm: true},
	})

	tmpl, err := h.tm.Get("todo", "show")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *WebHandler) Edit(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Edit todo ", slug)
	ctx := r.Context()

	list, err := h.service.GetListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	page := am.NewPage(list)
	page.SetFormAction(fmt.Sprintf("/todo/%s", slug))

	tmpl, err := h.tm.Get("todo", "edit")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, page)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *WebHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Update todo ", slug)
	ctx := r.Context()

	list, err := h.service.GetListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "List not found", http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	list.Name = name
	list.Description = description

	err = h.service.UpdateList(ctx, list)
	if err != nil {
		http.Error(w, "Failed to update list", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *WebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	h.Log().Info("Delete todo ", slug)
	ctx := r.Context()

	err := h.service.DeleteListBySlug(ctx, slug)
	if err != nil {
		http.Error(w, "Failed to delete list", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
