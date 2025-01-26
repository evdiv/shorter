package main

import (
	"log"
	"net/http"
)

// Config
const (
	host = "http://localhost"
	port = ":8080"
)

func main() {
	var h MainHandler

	log.Println("Starting the server on port", port)
	err := http.ListenAndServe(port, h)
	if err != nil {
		panic(err)
	}
}
