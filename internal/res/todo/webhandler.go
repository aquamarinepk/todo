package todo

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	todoResPath = "/todo"
	key         = am.Key
	method      = am.HTTPMethod
)

type WebHandler struct {
	*am.Handler
	service Service
	tm      *am.TemplateManager
}

func NewWebHandler(tm *am.TemplateManager, service Service, options ...am.Option) *WebHandler {
	handler := am.NewHandler("web-handler", options...)
	return &WebHandler{
		Handler: handler,
		service: service,
		tm:      tm,
	}
}

func (h *WebHandler) List(w http.ResponseWriter, r *http.Request) {
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

func (h *WebHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New todo form")

	page := am.NewPage(List{})
	page.SetFormAction(todoResPath)
	page.SetFormMethod(method.POST)
	page.SetFormButtonText("Create")
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("todo", "new")
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

func (h *WebHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create todo")
	ctx := r.Context()

	name := r.FormValue("name")
	description := r.FormValue("description")
	list := NewList(name, description)

	err := h.service.Create(ctx, list)
	if err != nil {
		http.Error(w, am.ErrCannotCreateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Show todo ", id)
	ctx := r.Context()

	listID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	list, err := h.service.Get(ctx, listID)
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
		{URL: fmt.Sprintf("%s/%s/edit", todoResPath, id), Text: "Edit", Style: blue},
		{URL: fmt.Sprintf("%s/%s/delete", todoResPath, id), Text: "Delete", Style: red},
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

func (h *WebHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Edit todo ", id)
	ctx := r.Context()

	listID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	list, err := h.service.Get(ctx, listID)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	page := am.NewPage(list)
	page.SetFormAction(fmt.Sprintf("%s/%s", todoResPath, id))
	page.SetFormMethod(method.PUT)
	page.SetFormButtonText("Update")
	page.GenCSRFToken(r)

	tmpl, err := h.tm.Get("todo", "edit")
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

func (h *WebHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Update todo ", id)
	ctx := r.Context()

	listID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	list, err := h.service.Get(ctx, listID)
	if err != nil {
		http.Error(w, am.ErrResourceNotFound, http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	list.Name = name
	list.Description = description

	err = h.service.Update(ctx, list)
	if err != nil {
		http.Error(w, am.ErrCannotUpdateResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}

func (h *WebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Delete todo ", id)
	ctx := r.Context()

	listID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, am.ErrInvalidID, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(ctx, listID)
	if err != nil {
		http.Error(w, am.ErrCannotDeleteResource, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, todoResPath, http.StatusSeeOther)
}
