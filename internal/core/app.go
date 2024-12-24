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

func NewApp(log am.Logger) *App {
	base := am.NewApp(log)
	return &App{
		base: base,
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
