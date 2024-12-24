package main

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/todo"
)

func main() {
	log := am.NewLogger("info")
	flagDefs := map[string]interface{}{
		"port": "8080",
		// TODO: Add any other required flags.
	}
	cfg := am.LoadCfg("TODO", flagDefs)

	repo := todo.NewRepo()
	service := todo.NewService(repo)

	opts := opts(log, cfg)

	webHandler := todo.NewWebHandler(service, opts...)
	apiHandler := todo.NewAPIHandler(service, opts...)

	webRouter := am.NewRouter(opts...)
	webRouter.Mount("/todo", todo.NewWebRouter(webHandler))

	apiRouter := am.NewRouter(opts...)
	apiRouter.Mount("/api/todo", todo.NewAPIRouter(apiHandler))

	webServer := am.NewServer(cfg.StrValOrDef("port", "8080"), webRouter, opts...)
	apiServer := am.NewServer(cfg.StrValOrDef("api_port", "8081"), apiRouter, opts...)

	go webServer.Start()
	apiServer.Start()
}

func opts(log am.Logger, cfg *am.Config) []am.Option {
	return []am.Option{
		am.WithLog(log),
		am.WithCfg(cfg),
	}
}
