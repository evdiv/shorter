package app

import (
	"log"
	"net/http"
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
	config.InitConfig()

	// Initialize storage
	dataStorage, err := storage.NewFileStorage()
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
	log.Println("Server running on", config.AppConfig.LocalPort)
	return http.ListenAndServe(config.AppConfig.LocalPort, a.Router)
}
