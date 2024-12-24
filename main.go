package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aquamarinepk/todo/internal/am"
	"github.com/aquamarinepk/todo/internal/feat/todo"
)

func main() {
	log := am.NewLogger("info")
	repo := todo.NewBaseRepo()
	service := todo.NewBaseService(repo)

	webHandler := todo.NewHandler(service, log)
	apiHandler := todo.NewAPIHandler(service)

	webRouter := am.NewRouter(log)
	webRouter.Mount("/todo", todo.NewWebRouter(webHandler))

	apiRouter := am.NewRouter(log)
	apiRouter.Mount("/api/todo", todo.NewAPIRouter(apiHandler))

	webServer := &http.Server{
		Addr:    ":8080",
		Handler: webRouter,
	}

	apiServer := &http.Server{
		Addr:    ":8081",
		Handler: apiRouter,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Info("Starting web server on :8080")
		err := webServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Could not listen on :8080: %v\n", err)
		}
	}()

	go func() {
		log.Info("Starting API server on :8081")
		err := apiServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Could not listen on :8081: %v\n", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Shutting down web server...")
	if err := webServer.Shutdown(ctx); err != nil {
		log.Errorf("Web server forced to shutdown: %v", err)
	}

	log.Info("Shutting down API server...")
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Errorf("API server forced to shutdown: %v", err)
	}

	log.Info("Servers stopped gracefully")
}
