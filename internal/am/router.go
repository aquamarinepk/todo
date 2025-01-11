package am

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
	Core Core
}

func NewRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)

	r := &Router{
		Router: chi.NewRouter(),
		Core:   core,
	}

	csrf := CSRFMw(core.Cfg())

	r.Use(MethodOverrideMw)
	r.Use(csrf)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Log().Info(req.Method, " ", req.URL.Path)
	r.Router.ServeHTTP(w, req)
}

func (r *Router) Name() string {
	return r.Core.Name()
}

func (r *Router) SetName(name string) {
	r.Core.SetName(name)
}

func (r *Router) Log() Logger {
	return r.Core.Log()
}

func (r *Router) SetLog(log Logger) {
	r.Core.SetLog(log)
}

func (r *Router) Cfg() *Config {
	return r.Core.Cfg()
}

func (r *Router) SetCfg(cfg *Config) {
	r.Core.SetCfg(cfg)
}

func (r *Router) Setup(ctx context.Context) error {
	return r.Core.Setup(ctx)
}

func (r *Router) Start(ctx context.Context) error {
	return r.Core.Start(ctx)
}

func (r *Router) Stop(ctx context.Context) error {
	return r.Core.Stop(ctx)
}
