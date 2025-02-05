package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
)

func main() {

	// Load configuration
	config.NewConfig(config.EnvConfigLoader{}, config.FlagConfigLoader{})

	// Start HTTP router
	r := chi.NewRouter()

	r.Post("/", handlers.PostURL)
	r.Get("/{urlKey}", handlers.GetURL)

	// Fallback for empty key
	r.Get("/", handlers.GetURL)

	err := http.ListenAndServe(config.Local.Port, r)
	if err != nil {
		panic(err)
	}
}
