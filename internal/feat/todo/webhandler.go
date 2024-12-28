package todo

import (
	"bytes"
	"context"
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/go-chi/chi/v5"
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
	w.Write([]byte("List of lists"))
}

func (h *WebHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New todo form")
	tmpl, err := h.tm.Get("todo", "new")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (h *WebHandler) New2(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("New todo form")
	tmpl, err := h.tm.Get("todo", "new")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func (h *WebHandler) Create(w http.ResponseWriter, r *http.Request) {
	h.Log().Info("Create todo")
	w.Write([]byte("Create todo"))
}

func (h *WebHandler) Show(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Show todo ", id)
	w.Write([]byte("Show todo " + id))
}

func (h *WebHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Edit todo ", id)
	w.Write([]byte("Edit todo " + id))
}

func (h *WebHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Update todo ", id)
	w.Write([]byte("Update todo " + id))
}

func (h *WebHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.Log().Info("Delete todo ", id)
	w.Write([]byte("Delete todo " + id))
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
