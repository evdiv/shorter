package main

import (
	"log"
	"shorter/internal/app"
)

func main() {
	//Initialize app
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	//Run app
	err = application.Run()
	if err != nil {
		log.Fatal(err)
	}
}
