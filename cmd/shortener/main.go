package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/storage"
)

func main() {

	// Load configuration
	config.NewConfig(config.EnvConfigLoader{}, config.FlagConfigLoader{})

	// Initialize in memory storage
	memStorage := storage.NewMemoryStorage()

	// Handlers with storage dependency
	h := handlers.NewHandlers(memStorage)

	// Start HTTP router
	r := chi.NewRouter()

	r.Post("/", h.PostURL)
	r.Get("/{urlKey}", h.GetURL)

	// Fallback for empty key
	r.Get("/", h.GetURL)

	err := http.ListenAndServe(config.Local.Port, r)
	if err != nil {
		panic(err)
	}
}
