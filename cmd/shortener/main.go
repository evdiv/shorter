package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shorter/cmd/shortener/config"
)

func main() {

	// Load configuration
	config.LoadConfig()

	// Start HTTP router
	r := chi.NewRouter()

	r.Post("/", PostURL)
	r.Get("/{urlKey}", GetURL)

	// Fallback for empty key
	r.Get("/", GetURL)

	err := http.ListenAndServe(config.Local.Port, r)
	if err != nil {
		panic(err)
	}
}
