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

const (
	layoutDir    = "assets/template/layout"
	handlerDir   = "assets/template/handler"
	layoutPrefix = layoutDir + "/"
)

type TemplateManager struct {
	core      Core
	assetsFS  embed.FS
	layouts   sync.Map
	templates sync.Map
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
	tm.loadLayoutTemplates()
	tm.loadHandlerTemplates()
}

func (tm *TemplateManager) loadLayoutTemplates() {
	tm.loadLayoutTemplatesFromDir(layoutDir)
}

func (tm *TemplateManager) loadLayoutTemplatesFromDir(path string) {
	entries, err := tm.assetsFS.ReadDir(path)
	if err != nil {
		tm.Log().Error("Failed to read layout subdirectory: ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			tm.loadLayoutTemplatesFromDir(filepath.Join(path, entry.Name()))
		} else {
			key := strings.TrimPrefix(filepath.ToSlash(filepath.Join(path, entry.Name())), layoutPrefix)
			key = strings.TrimSuffix(key, filepath.Ext(key))
			key = strings.ReplaceAll(key, "/", ":")
			tm.loadTemplate(key, filepath.Join(path, entry.Name()))
		}
	}
}

func (tm *TemplateManager) loadHandlerTemplates() {
	entries, err := tm.assetsFS.ReadDir(handlerDir)
	if err != nil {
		tm.Log().Error("Failed to read handler directory: ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			handler := strings.ToLower(entry.Name())
			tm.loadHandlerTemplatesFromDir(handler, filepath.Join(handlerDir, handler))
		}
	}
}

func (tm *TemplateManager) loadHandlerTemplatesFromDir(handler, path string) {
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
	if strings.HasPrefix(path, layoutDir) {
		if _, loaded := tm.layouts.LoadOrStore(key, tmpl); loaded {
			tm.Log().Infof("Layout template key %s already exists, skipping", key)
		}
		return
	}
	if _, loaded := tm.templates.LoadOrStore(key, tmpl); loaded {
		tm.Log().Infof("Handler template key %s already exists, skipping", key)
	}
}

func (tm *TemplateManager) Get(handler, action string) (*template.Template, error) {
	key := handler + ":" + action
	if tmpl, ok := tm.templates.Load(key); ok {
		return tmpl.(*template.Template), nil
	}
	if tmpl, ok := tm.layouts.Load(handler + ":layout"); ok {
		return tmpl.(*template.Template), nil
	}
	if tmpl, ok := tm.layouts.Load("layout"); ok {
		return tmpl.(*template.Template), nil
	}
	return nil, errors.New("template not found")
}

func (tm *TemplateManager) Debug() {
	tm.debugTemplates(&tm.layouts, "Layout template key")
	tm.debugTemplates(&tm.templates, "Handler template key")
}

func (tm *TemplateManager) debugTemplates(store *sync.Map, keyPrefix string) {
	store.Range(func(key, value interface{}) bool {
		tmpl := value.(*template.Template)
		tm.Log().Infof("%s: %s, Template name: %s, Defined templates: %v", keyPrefix, key, tmpl.Name(), tmpl.DefinedTemplates())
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
