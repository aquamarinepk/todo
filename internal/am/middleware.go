package am

import (
	"net/http"
	"sync"

	"github.com/gorilla/csrf"
)

const defaultCSRFKey = "set-a-csrf-key!"

var (
	csrfMiddleware func(http.Handler) http.Handler
	once           sync.Once
)

// MethodOverrideMw is a middleware that checks for a _method form field and overrides the request method.
func MethodOverrideMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if override := r.FormValue("_method"); override != "" {
				r.Method = override
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CSRFMw is a middleware that protects against CSRF attacks.
func CSRFMw(cfg *Config) func(http.Handler) http.Handler {
	if cfg == nil {
		return passThroughMw
	}

	initCSRF(cfg)

	return func(next http.Handler) http.Handler {
		return csrfMiddleware(next)
	}
}

func passThroughMw(next http.Handler) http.Handler {
	return next
}

func initCSRF(cfg *Config) {
	once.Do(func() {
		key := cfg.StrValOrDef(Key.SecCSRFKey, defaultCSRFKey)
		to := cfg.StrValOrDef(Key.SecCSRFRedirect, "/csrf-error")

		csrfMiddleware = csrf.Protect(
			[]byte(key),
			csrf.FieldName(CSRFFieldName), 
			csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, to, http.StatusFound)
			})),
		)
	})
}
