package am

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
	Core Core
}

func NewRouter(log Logger) *Router {
	core := NewCore(log)
	return &Router{
		Router: chi.NewRouter(),
		Core:   core,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Log().Info(req.Method, " ", req.URL.Path)
	r.Router.ServeHTTP(w, req)
}

func (r *Router) Log() Logger {
	return r.Core.Log()
}
