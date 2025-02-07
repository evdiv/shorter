package main

import (
	"net/http"
	"shorter/internal/config"
	"shorter/internal/handlers"
)

func main() {

	// Load configuration
	config.NewConfig(config.EnvConfigLoader{}, config.FlagConfigLoader{})

	// Start HTTP router
	r := handlers.NewRouter()

	err := http.ListenAndServe(config.Local.Port, r)
	if err != nil {
		panic(err)
	}
}
