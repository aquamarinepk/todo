package core

import (
	"context"
	"embed"
	"os"
	"os/signal"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/res/todo"
)

type App struct {
	*am.App
	repo    todo.Repo
	service todo.Service
}

func NewApp(name, version string, fs embed.FS, opts ...am.Option) *App {
	core := am.NewApp(name, version, fs, opts...)
	app := &App{
		App: core,
	}
	return app
}


func (app *App) Start(ctx context.Context) error {
	err := app.App.Start(ctx)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	return app.Core.Stop(ctx)
}

func (app *App) SetRepo(repo todo.Repo) {
	app.repo = repo
}

func (app *App) SetService(service todo.Service) {
	app.service = service
}
