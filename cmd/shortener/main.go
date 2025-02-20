package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/middleware"
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

	//Add logging middleware
	r.Use(middleware.WithLogging)

	r.Use(middleware.WithGzip)

	r.Post("/", h.PostURL)
	r.Post("/api/shorten", h.ShortenURL)

	r.Get("/{urlKey}", h.GetURL)
	r.Get("/", h.GetURL)

	err := http.ListenAndServe(config.Local.Port, r)
	if err != nil {
		panic(err)
	}
}
