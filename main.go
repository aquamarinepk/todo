package main

import (
	"context"
	"embed"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/core"
	"github.com/aquamarinepk/todo/internal/feat/auth"
	"github.com/aquamarinepk/todo/internal/res/todo"
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

	authRepo := auth.NewRepo(queryManager, opts...)
	authService := auth.NewService(authRepo)
	authWebHandler := auth.NewWebHandler(templateManager, authService, opts...)
	authWebRouter := auth.NewWebRouter(authWebHandler, opts...)
	authAPIHandler := auth.NewAPIHandler(authService, opts...)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler, opts...)

	// TODO: Auth related functionality not fully implemented yet.
	app.MountFeatWeb("/auth", authWebRouter)
	app.MountFeatAPI(version, "/auth", authAPIRouter)

	todoRepo := todo.NewRepo(queryManager, opts...)
	todoService := todo.NewService(todoRepo)
	todoWebHandler := todo.NewWebHandler(templateManager, todoService, opts...)
	todoWebRouter := todo.NewWebRouter(todoWebHandler, opts...)
	todoAPIHandler := todo.NewAPIHandler(todoService, opts...)
	todoAPIRouter := todo.NewAPIRouter(todoAPIHandler, opts...)

	app.MountWeb("/todo", todoWebRouter)
	app.MountAPI(version, "/todo", todoAPIRouter)

	app.Add(templateManager)
	app.Add(todoRepo)
	app.Add(todoService)
	app.Add(todoWebHandler)
	app.Add(todoAPIHandler)
	app.Add(todoWebRouter)
	app.Add(todoAPIRouter)

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
