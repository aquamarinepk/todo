package main

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/todo"
)

const (
	WebHostKey = "server.web.host"
	WebPortKey = "server.web.port"
	APIHostKey = "server.api.host"
	AIPortKey  = "server.api.port"
)

func main() {
	log := am.NewLogger("info")
	flags := map[string]interface{}{
		WebPortKey: "8080",
		AIPortKey:  "8081",
		WebHostKey: "localhost",
		APIHostKey: "localhost",
	}
	cfg := am.LoadCfg("TODO", flags)

	repo := todo.NewRepo()
	service := todo.NewService(repo)

	opts := opts(log, cfg)

	webHandler := todo.NewWebHandler(service, opts...)
	apiHandler := todo.NewAPIHandler(service, opts...)

	webRouter := am.NewRouter(opts...)
	webRouter.Mount("/todo", todo.NewWebRouter(webHandler))

	apiRouter := am.NewRouter(opts...)
	apiRouter.Mount("/api/todo", todo.NewAPIRouter(apiHandler))

	webServer := am.NewServer(WebHostKey, WebPortKey, webRouter, opts...)
	apiServer := am.NewServer(APIHostKey, AIPortKey, apiRouter, opts...)

	go webServer.Start()
	apiServer.Start()
}

func opts(log am.Logger, cfg *am.Config) []am.Option {
	return []am.Option{
		am.WithLog(log),
		am.WithCfg(cfg),
	}
}
