package main

import (
	"log"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
	"shorter/internal/router"
	"shorter/internal/storage"
)

func main() {

	// Load configuration
	config.InitConfig()

	// Initialize storage
	dataStorage, err := storage.NewFileStorage()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers
	h := handlers.NewHandlers(dataStorage)

	// Initialize router
	r := router.NewRouter(h)

	err = http.ListenAndServe(config.AppConfig.LocalPort, r)
	if err != nil {

		panic(err)
	}
}
