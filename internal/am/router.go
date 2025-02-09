package am

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Core
	chi.Router
}

func NewRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
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
