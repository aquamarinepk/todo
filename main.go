package main

import (
	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/todo"
)

func main() {
	log := am.NewLogger("info")

	repo := todo.NewRepo()
	service := todo.NewService(repo)

	webHandler := todo.NewWebHandler(service, log)
	apiHandler := todo.NewAPIHandler(service)

	webRouter := am.NewRouter(log)
	webRouter.Mount("/todo", todo.NewWebRouter(webHandler))

	apiRouter := am.NewRouter(log)
	apiRouter.Mount("/api/todo", todo.NewAPIRouter(apiHandler))

	webServer := am.NewServer("8080", webRouter, log)
	apiServer := am.NewServer("8081", apiRouter, log)

	go webServer.Start()
	apiServer.Start()
}
