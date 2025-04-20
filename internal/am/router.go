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
	router := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	for _, opt := range opts {
		opt(router)
	}

	return router
}

func NewWebRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	cfg := core.Cfg()
	//hashKey := cfg.ByteSliceVal(Key.SecHashKey)
	csrf := CSRFMw(cfg)

	r.Use(MethodOverrideMw)
	r.Use(csrf)

	//flash := NewFlashMiddleware(hashKey)
	//r.Use(flash.Middleware)

	return r
}

func NewAPIRouter(name string, opts ...Option) *Router {
	core := NewCore(name, opts...)

	r := &Router{
		Core:   core,
		Router: chi.NewRouter(),
	}

	r.Use(MethodOverrideMw)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			r.Log().Error("Error serving request: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()
	r.Router.ServeHTTP(w, req)
}
