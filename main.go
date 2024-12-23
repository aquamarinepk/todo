package main

import (
	"net/http"

	"github.com/aquamarinepk/todo/internal/feat/todo"
	"github.com/go-chi/chi/v5"
)

func main() {
	repo := todo.NewBaseRepo()

	service := todo.NewBaseService(repo)

	handler := todo.NewHandler(service)

	r := chi.NewRouter()
	r.Mount("/todo", todo.NewRouter(handler))

	http.ListenAndServe(":8080", r)
}
