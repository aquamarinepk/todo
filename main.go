package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aquamarinepk/todo/internal/feat/todo"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := todo.NewBaseRepo()
	service := todo.NewBaseService(repo)

	webHandler := todo.NewHandler(service)
	apiHandler := todo.NewAPIHandler(service)

	webRouter := chi.NewRouter()
	webRouter.Mount("/todo", todo.NewWebRouter(webHandler))

	apiRouter := chi.NewRouter()
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
		log.Println("Starting web server on :8080")
		err := webServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on :8080: %v\n", err)
		}
	}()

	go func() {
		log.Println("Starting API server on :8081")
		err := apiServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on :8081: %v\n", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down web server...")
	if err := webServer.Shutdown(ctx); err != nil {
		log.Fatalf("Web server forced to shutdown: %v", err)
	}

	log.Println("Shutting down API server...")
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Fatalf("API server forced to shutdown: %v", err)
	}

	log.Println("Servers stopped gracefully")
}
