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
	"time"
)

type App struct {
	Router     http.Handler
	Config     *config.Config
	Storage    storage.Storer
	DeleteChan chan models.KeysToDelete
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
	deleteChan := make(chan models.KeysToDelete, appConfig.DeleteBufferSize)

	// Initialize handlers
	h := handlers.NewHandlers(appStorage, deleteChan)

	// Initialize router
	r := router.NewRouter(h)

	return &App{
		Router:     r,
		Config:     appConfig,
		Storage:    appStorage,
		DeleteChan: deleteChan,
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

	ticker := time.NewTicker(10 * time.Second)
	var keysToDelete []models.KeysToDelete

	for {
		select {
		case k := <-a.DeleteChan:
			//Add a key to the slice for deleting later
			keysToDelete = append(keysToDelete, k)
		case <-ticker.C:
			//Wait for at least one message
			if len(keysToDelete) == 0 {
				continue
			}
			//update all incoming requests at once
			_, err := a.Storage.DeleteBatch(keysToDelete)
			if err != nil {
				log.Printf("Failed to delete records: %v\n", err)
				continue
			}

			//Remove all keys that have been sent
			keysToDelete = nil

		}
	}
}
