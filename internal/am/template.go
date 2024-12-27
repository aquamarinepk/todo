package am

import (
	"context"
	"embed"
	"errors"
	"html/template"
	"path/filepath"
	"strings"
	"sync"
)

type TemplateManager struct {
	core      Core
	templates sync.Map
	assetsFS  embed.FS
}

func NewTemplateManager(assetsFS embed.FS, opts ...Option) *TemplateManager {
	core := NewCore("template-manager", opts...)
	tm := &TemplateManager{
		core:     core,
		assetsFS: assetsFS,
	}

	return tm
}

func (tm *TemplateManager) Load() {
	tm.loadTemplate("layout", "assets/layout/layout.html")

	entries, err := tm.assetsFS.ReadDir("assets/layout")
	if err != nil {
		tm.Log().Error("Failed to read assets/layout directory: ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			handler := strings.ToLower(entry.Name())
			tm.loadHandlerTemplates(handler, filepath.Join("assets/layout", handler))
		} else {
			name := strings.ToLower(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
			if name != "layout" {
				tm.loadTemplate(name, filepath.Join("assets/layout", entry.Name()))
			}
		}
	}
}

func (tm *TemplateManager) loadHandlerTemplates(handler, path string) {
	entries, err := tm.assetsFS.ReadDir(path)
	if err != nil {
		tm.Log().Error("Failed to read handler directory: ", err)
		return
	}

	for _, entry := range entries {
		name := strings.ToLower(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
		key := handler + ":" + name
		tm.loadTemplate(key, filepath.Join(path, entry.Name()))
	}
}

func (tm *TemplateManager) loadTemplate(key, path string) {
	tmpl, err := template.ParseFS(tm.assetsFS, path)
	if err != nil {
		tm.Log().Error("Failed to load template: ", err)
		return
	}
	if _, loaded := tm.templates.LoadOrStore(key, tmpl); loaded {
		tm.Log().Info("Template key %s already exists, skipping", key)
	}
}

func (tm *TemplateManager) Get(handler, action string) (*template.Template, error) {
	key := handler + ":" + action
	if tmpl, ok := tm.templates.Load(key); ok {
		return tmpl.(*template.Template), nil
	}
	if tmpl, ok := tm.templates.Load(handler + ":layout"); ok {
		return tmpl.(*template.Template), nil
	}
	if tmpl, ok := tm.templates.Load("layout"); ok {
		return tmpl.(*template.Template), nil
	}
	return nil, errors.New("template not found")
}

func (tm *TemplateManager) Debug() {
	tm.templates.Range(func(key, value interface{}) bool {
		tmpl := value.(*template.Template)
		tm.Log().Infof("Template key: %s, Template name: %s, Defined templates: %v", key, tmpl.Name(), tmpl.DefinedTemplates())
		return true
	})
}

func (tm *TemplateManager) Name() string {
	return tm.core.Name()
}

func (tm *TemplateManager) SetName(name string) {
	tm.core.SetName(name)
}

func (tm *TemplateManager) Log() Logger {
	return tm.core.Log()
}

func (tm *TemplateManager) SetLog(log Logger) {
	tm.core.SetLog(log)
}

func (tm *TemplateManager) Cfg() *Config {
	return tm.core.Cfg()
}

func (tm *TemplateManager) SetCfg(cfg *Config) {
	tm.core.SetCfg(cfg)
}

func (tm *TemplateManager) Setup(ctx context.Context) error {
	tm.Load()
	return tm.core.Setup(ctx)
}

func (tm *TemplateManager) Start(ctx context.Context) error {
	return tm.core.Start(ctx)
}

func (tm *TemplateManager) Stop(ctx context.Context) error {
	return tm.core.Stop(ctx)
}
