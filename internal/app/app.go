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
	Router  http.Handler
	Config  *config.Config
	Storage storage.Storer
}

func NewApp() (*App, error) {
	// Load configuration
	appConfig := config.NewConfig(config.LoadFromEnv, config.LoadFromFlags, config.LoadDefault)

	// Initialize storage
	appStorage, err := storage.NewStorage(*appConfig)

	if err != nil {
		return nil, err
	}

	// Initialize handlers
	h := handlers.NewHandlers(appStorage)

	// Initialize router
	r := router.NewRouter(h)

	return &App{
		Router:  r,
		Config:  appConfig,
		Storage: appStorage,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Local Host: " + a.Config.LocalHost)
	log.Println("Result Host: " + a.Config.ResultHost)
	log.Println("File Storage Path: " + a.Config.StoragePath)
	log.Println("Db Connection String: " + a.Config.DbConnection)

	go func() {
		_ = http.ListenAndServe(config.GetPort("Local"), a.Router)
	}()

	<-sigChan
	log.Println("Shutdown signal received")
	a.Storage.Close()
	return nil
}
