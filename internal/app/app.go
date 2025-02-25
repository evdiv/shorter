package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/router"
	"shorter/internal/storage"
)

type App struct {
	Router      http.Handler
	DataStorage *storage.FileStorage
}

func NewApp() (*App, error) {
	// Load configuration
	config.NewConfig()

	// Initialize storage
	dataStorage, err := storage.NewFileStorage(config.AppConfig.StoragePath)
	if err != nil {
		return nil, err
	}

	// Initialize handlers
	h := handlers.NewHandlers(dataStorage)

	// Initialize router
	r := router.NewRouter(h)

	return &App{
		Router:      r,
		DataStorage: dataStorage,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Server running on", config.GetPort("Local"))

	log.Println("Config loaded from: " + config.AppConfig.LoadedFrom)
	log.Println("Local Host: " + config.AppConfig.LocalHost)
	log.Println("Result Host: " + config.AppConfig.ResultHost)
	log.Println("File Storage Path: " + config.AppConfig.StoragePath)

	go func() {
		_ = http.ListenAndServe(config.GetPort("Local"), a.Router)
	}()

	<-sigChan
	log.Println("Shutdown signal received")
	return nil
}
