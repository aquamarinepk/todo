package am

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
	log Logger
}

func NewRouter(logger Logger) *Router {
	return &Router{
		Router: chi.NewRouter(),
		log:    logger,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.log.Info(req.Method, " ", req.URL.Path)
	r.Router.ServeHTTP(w, req)
}
