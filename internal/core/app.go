package core

import (
	"net/http"

	"github.com/aquamarinepk/todo/internal/am"
)

type App struct {
	base      *am.App
	webServer *http.Server
	apiServer *http.Server
}

type Option func(*App)

func NewApp(options ...Option) *App {
	base := am.NewApp()
	app := &App{
		base: base,
	}
	for _, option := range options {
		option(app)
	}
	return app
}

func WithWebServer(server *http.Server) Option {
	return func(app *App) {
		app.webServer = server
	}
}

func WithAPIServer(server *http.Server) Option {
	return func(app *App) {
		app.apiServer = server
	}
}

func (app *App) SetWebServer(server *http.Server) {
	app.webServer = server
}

func (app *App) SetAPIServer(server *http.Server) {
	app.apiServer = server
}

func (app *App) Start() {
	app.base.StartServer(app.webServer, ":8080")
	app.base.StartServer(app.apiServer, ":8081")
}
