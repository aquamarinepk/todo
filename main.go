package main

import (
	"context"
	"embed"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/core"
	"github.com/aquamarinepk/todo/internal/feat/auth"
	"github.com/aquamarinepk/todo/internal/repo/sqlite"
	"github.com/aquamarinepk/todo/internal/res/todo"
)

const (
	name      = "todo"
	version   = "v1"
	namespace = "TODO"
	engine    = "sqlite"
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

	// FlashManager
	flashManager := am.NewFlashManager()

	// Append WithFlashMiddleware to opts
	opts = append(opts, am.WithFlashMiddleware(flashManager))

	// Create the app with the updated opts
	app := core.NewApp(name, version, assetsFS, opts...)

	queryManager := am.NewQueryManager(assetsFS, engine)
	templateManager := am.NewTemplateManager(assetsFS)

	// Migrator
	migrator := am.NewMigrator(assetsFS, engine)

	// Seeder
	// seeder := am.NewSeeder(assetsFS, engine)

	// FileServer
	fileServer := am.NewFileServer(assetsFS)
	app.MountFileServer("/", fileServer)

	// Auth feature
	authRepo := sqlite.NewAuthRepo(queryManager)
	authService := auth.NewService(authRepo)
	authWebHandler := auth.NewWebHandler(templateManager, flashManager, authService)
	authWebRouter := auth.NewWebRouter(authWebHandler)
	authAPIHandler := auth.NewAPIHandler(authService)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler)
	authSeeder := auth.NewSeeder(assetsFS, engine, authRepo)

	app.MountWeb("/auth", authWebRouter)
	app.MountAPI(version, "/auth", authAPIRouter)

	// Todo resource
	todoRepo := todo.NewRepo(queryManager)
	todoService := todo.NewService(todoRepo)
	todoWebHandler := todo.NewWebHandler(templateManager, todoService)
	todoWebRouter := todo.NewWebRouter(todoWebHandler)
	todoAPIHandler := todo.NewAPIHandler(todoService)
	todoAPIRouter := todo.NewAPIRouter(todoAPIHandler)

	app.MountResWeb("/todo", todoWebRouter)
	app.MountResAPI(version, "/todo", todoAPIRouter)

	// Add deps
	app.Add(migrator)
	app.Add(flashManager)
	app.Add(fileServer)
	app.Add(queryManager)
	app.Add(templateManager)
	app.Add(authRepo)
	app.Add(authService)
	app.Add(authWebHandler)
	app.Add(authAPIHandler)
	app.Add(authWebRouter)
	app.Add(authAPIRouter)
	app.Add(todoRepo)
	app.Add(todoService)
	app.Add(todoWebHandler)
	app.Add(todoAPIHandler)
	app.Add(todoWebRouter)
	app.Add(todoAPIRouter)
	app.Add(authSeeder)

	err := app.Setup(ctx)
	if err != nil {
		log.Error("Failed to setup the app: ", err)
		return
	}

	err = app.Start(ctx)
	if err != nil {
		log.Error("Failed to start the app: ", err)
	}
}
