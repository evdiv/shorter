package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Config
const (
	host = "http://localhost"
	port = ":8080"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", PostURL)
	r.Get("/{urlKey}", GetURL)

	// Fallback for empty key
	r.Get("/", GetURL)

	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}
