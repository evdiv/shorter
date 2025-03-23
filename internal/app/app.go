package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/models"
	"shorter/internal/router"
	"shorter/internal/storage"
)

type App struct {
	Router      http.Handler
	Config      *config.Config
	Storage     storage.Storer
	DeleteQueue chan models.KeysToDelete
}

func NewApp() (*App, error) {
	// Load configuration
	appConfig := config.NewConfig(config.LoadFromEnv, config.LoadFromFlags, config.LoadDefault)

	// Initialize storage
	appStorage, err := storage.NewStorage(*appConfig)

	if err != nil {
		return nil, err
	}

	// Create a channel for batch deleting records
	deleteQueue := make(chan models.KeysToDelete, 1024)

	// Initialize handlers
	h := handlers.NewHandlers(appStorage, deleteQueue)

	// Initialize router
	r := router.NewRouter(h)

	return &App{
		Router:      r,
		Config:      appConfig,
		Storage:     appStorage,
		DeleteQueue: deleteQueue,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Local Host: " + a.Config.LocalHost)
	log.Println("Result Host: " + a.Config.ResultHost)
	log.Println("File Storage Path: " + a.Config.StoragePath)
	log.Println("Db Connection String: " + a.Config.DBConnection)

	// Start background deletion worker
	go a.StartDeletionWorker()

	go func() {
		_ = http.ListenAndServe(config.GetPort("Local"), a.Router)
	}()

	<-sigChan
	log.Println("Shutdown signal received")
	a.Storage.Close()
	return nil
}

// StartDeletionWorker processes delete tasks from the queue.
func (a *App) StartDeletionWorker() {
	for task := range a.DeleteQueue {
		_, err := a.Storage.DeleteBatch(task.Keys, task.UserID)
		if err != nil {
			log.Printf("Failed to delete records: %v\n", err)
		}
	}
}
