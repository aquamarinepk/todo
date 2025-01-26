// internal/core/app.go
package core

import (
	"context"
	"os"
	"os/signal"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/res/todo"
)

type App struct {
	core      *am.App
	router    *am.Router
	webRouter *am.Router
	apiRouter *am.Router
	repo      todo.Repo
	service   todo.Service
}

func NewApp(name, version string, opts ...am.Option) *App {
	core := am.NewApp(name, version, opts...)
	app := &App{
		core: core,
	}
	return app
}

func (app *App) Setup(ctx context.Context) error {
	return app.core.Setup(ctx)
}

func (app *App) Start(ctx context.Context) error {
	err := app.core.Start(ctx)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	return app.Stop(ctx)
}

func (a *App) Add(dep am.Core) {
	a.core.Add(dep)
}

func (a *App) Dep(name string) (*am.Dep, bool) {
	return a.core.Dep(name)
}

func (app *App) Stop(ctx context.Context) error {
	return nil
}

func (app *App) SetRepo(repo todo.Repo) {
	app.repo = repo
}

func (app *App) SetService(service todo.Service) {
	app.service = service
}

func (app *App) SetWebRouter(router *am.Router) {
	app.webRouter = router
}

func (app *App) SetAPIRouter(router *am.Router) {
	app.apiRouter = router
}

func (app *App) MountWeb(path string, router *am.Router) {
	app.core.Mount(path, router)
}

func (app *App) MountAPI(version, path string, router *am.Router) {
	app.core.MountAPI(version, path, router)
}

func (app *App) MountFeatWeb(path string, router *am.Router) {
	app.core.MountFeat(path, router)
}

func (app *App) MountFeatAPI(version, path string, router *am.Router) {
	app.core.MountFeatAPI(version, path, router)
}

func (app *App) Log() am.Logger {
	return app.core.Log()
}

func (app *App) Cfg() *am.Config {
	return app.core.Cfg()
}
