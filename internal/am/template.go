package am

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"sync"
)

const (
	layoutPath    = "assets/template/layout"
	handlerPath   = "assets/template/handler"
	partialDir    = "partial"
	defaultLayout = "layout.html"
	mainTemplate  = "page"
)

type TemplateManager struct {
	core      Core
	assetsFS  embed.FS
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
	tm.loadTemplates()
}

func (tm *TemplateManager) loadTemplates() {
	entries, err := tm.assetsFS.ReadDir(handlerPath)
	if err != nil {
		tm.Log().Error("Failed to read handler directory: ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			handler := strings.ToLower(entry.Name())
			tm.loadTemplatesFromDir(handler, filepath.Join(handlerPath, handler))
		}
	}
}

func (tm *TemplateManager) loadTemplatesFromDir(handler, path string) {
	entries, err := tm.assetsFS.ReadDir(path)
	if err != nil {
		tm.Log().Error("Failed to read handler directory: ", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == partialDir {
			continue
		}
		name := strings.ToLower(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
		key := handler + ":" + name
		tm.loadTemplate(key, filepath.Join(path, entry.Name()), handler)
	}
}

func (tm *TemplateManager) loadTemplate(key, path, handler string) {
	tm.Log().Info(header("=", 120))
	defer tm.Log().Info(header("=", 120))

	tm.Log().Infof("Loading template: key=%s, path=%s, handler=%s", key, path, handler)

	partials, err := tm.assetsFS.ReadDir(filepath.Join(handlerPath, handler, partialDir))
	if err != nil {
		tm.Log().Error("Failed to read partials directory: ", err)
		return
	}

	partialPaths := []string{}
	for _, partial := range partials {
		partialPath := filepath.Join(handlerPath, handler, partialDir, partial.Name())
		partialPaths = append(partialPaths, partialPath)
		tm.Log().Infof("Found partial: %s", partialPath)
	}

	layoutPath := tm.findLayoutPath(handler, filepath.Base(path))
	tm.Log().Infof("Using layout: %s", layoutPath)

	allPaths := append([]string{layoutPath, path}, partialPaths...)
	tm.Log().Infof("All template paths: %v", allPaths)

	tmpl, err := template.New(mainTemplate).ParseFS(tm.assetsFS, allPaths...)
	if err != nil {
		tm.Log().Error("Failed to load template: ", err)
		return
	}

	tm.Log().Infof("Successfully loaded template: %s", key)
	if _, loaded := tm.templates.LoadOrStore(key, tmpl); loaded {
		tm.Log().Infof("Template key %s already exists, skipping", key)
	}
}

func (tm *TemplateManager) findLayoutPath(handler, action string) string {
	actionLayout := filepath.Join(layoutPath, handler, action)
	tm.Log().Infof("Evaluating specific action layout path: %s", actionLayout)
	if _, err := tm.assetsFS.Open(actionLayout); err == nil {
		tm.Log().Infof("Found specific action layout: %s", actionLayout)
		return actionLayout
	}

	handlerLayout := filepath.Join(layoutPath, handler, defaultLayout)
	tm.Log().Infof("Evaluating handler layout path: %s", handlerLayout)
	if _, err := tm.assetsFS.Open(handlerLayout); err == nil {
		tm.Log().Infof("Found handler layout: %s", handlerLayout)
		return handlerLayout
	}

	globalLayout := filepath.Join(layoutPath, defaultLayout)
	tm.Log().Infof("Evaluating global layout path: %s", globalLayout)
	if _, err := tm.assetsFS.Open(globalLayout); err == nil {
		tm.Log().Infof("Found global layout: %s", globalLayout)
		return globalLayout
	}

	tm.Log().Info("No specific, handler, or global layout found")
	return ""
}

func (tm *TemplateManager) Get(handler, action string) (*template.Template, error) {
	key := handler + ":" + action
	if tmpl, ok := tm.templates.Load(key); ok {
		return tmpl.(*template.Template), nil
	}
	return nil, errors.New("template not found")
}

func debugTemplate(key string, tmpl *template.Template) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Template key: %s\n", key))
	sb.WriteString(fmt.Sprintf("  Template name: %s\n", tmpl.Name()))
	sb.WriteString("  Defined templates:\n")
	for _, tmpl := range tmpl.Templates() {
		sb.WriteString(fmt.Sprintf("    %s\n", tmpl.Name()))
	}
	return sb.String()
}

func header(char string, count int) string {
	return strings.Repeat(char, count)
}

func (tm *TemplateManager) Debug() {
	tm.templates.Range(func(key, value interface{}) bool {
		tmpl := value.(*template.Template)
		tm.Log().Info(debugTemplate(key.(string), tmpl))
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
	return nil
}

func (tm *TemplateManager) Start(ctx context.Context) error {
	return tm.core.Start(ctx)
}

func (tm *TemplateManager) Stop(ctx context.Context) error {
	return tm.core.Stop(ctx)
}
