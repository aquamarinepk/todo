package main

import (
	"context"
	"embed"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/core"
	"github.com/aquamarinepk/todo/internal/feat/todo"
)

const (
	name      = "todo"
	version   = "v1"
	namespace = "TODO"
)

var (
	//go:embed assets
	assetsFS embed.FS
)

func main() {
	ctx := context.Background()
	log := am.NewLogger("info")
	cfg := am.LoadCfg(namespace, am.Flags)
	opts := am.DefOpts(log, cfg)

	app := core.NewApp(name, version, opts...)

	queryManager := am.NewQueryManager(assetsFS, opts...)
	templateManager := am.NewTemplateManager(assetsFS, opts...)

	todoRepo := todo.NewRepo(queryManager, opts...)
	todoService := todo.NewService(todoRepo)

	todoWebHandler := todo.NewWebHandler(templateManager, todoService, opts...)
	webRouter := todo.NewWebRouter(todoWebHandler, opts...)

	todoAPIHandler := todo.NewAPIHandler(todoService, opts...)
	apiRouter := todo.NewAPIRouter(todoAPIHandler, opts...)

	app.MountWeb("/todo", webRouter)
	app.MountAPI(version, "/todo", apiRouter)

	app.Add(templateManager)
	app.Add(todoRepo)
	app.Add(todoService)
	app.Add(todoWebHandler)
	app.Add(todoAPIHandler)
	app.Add(webRouter)
	app.Add(apiRouter)

	err := app.Setup(ctx)
	if err != nil {
		log.Error("Failed to setup the app: ", err)
		return
	}

	//queryManager.Debug()
	//templateManager.Debug()

	err = app.Start(ctx)
	if err != nil {
		log.Error("Failed to start the app: ", err)
	}
}
