package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/middleware"
	"shorter/internal/storage"
)

func main() {

	// Load configuration
	config.InitConfig()

	// In memory storage
	//dataStorage := storage.NewMemoryStorage()

	// local file storage
	dataStorage, err := storage.NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}

	// Handlers with storage dependency
	h := handlers.NewHandlers(dataStorage)

	// Start HTTP router
	r := chi.NewRouter()

	//Add logging middleware
	r.Use(middleware.WithLogging)

	r.Use(middleware.WithGzip)

	r.Post("/", h.PostURL)
	r.Post("/api/shorten", h.ShortenURL)

	r.Get("/{urlKey}", h.GetURL)
	r.Get("/", h.GetURL)

	err = http.ListenAndServe(config.AppConfig.LocalPort, r)
	if err != nil {

		panic(err)
	}
}
