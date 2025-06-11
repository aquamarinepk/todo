package auth

import (
	"github.com/aquamarinepk/todo/internal/am"
)

// NewAPIRouter creates a new API router for the todo feature.
// Both GET and POST requests will be mounted to the app's router that handles `/cq` requests.
func NewAPIRouter(handler *APIHandler, opts ...am.Option) *am.Router {
	r := am.NewRouter("api-router", opts...)

	r.Get("/", handler.ListUsers)
	r.Get("/{slug}", handler.ShowUser)
	r.Post("/create-user", handler.CreateUser)
	r.Post("/update-user", handler.UpdateUser)
	r.Post("/delete-user", handler.DeleteUser)
	r.Post("/roles", handler.CreateRole)
	r.Put("/roles", handler.UpdateRole)
	r.Delete("/roles", handler.DeleteRole)

	return r
}
