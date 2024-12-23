package main

import (
	"net/http"

	"github.com/aquamarinepk/todo.git/internal/feat/list"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Mount("/list", list.NewRouter())

	http.ListenAndServe(":8080", r)
}
