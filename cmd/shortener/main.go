package main

import (
	"net/http"
)

// Config
const (
	host = "http://localhost"
	port = ":8080"
)

func main() {
	var h MainHandler

	err := http.ListenAndServe(port, h)
	if err != nil {
		panic(err)
	}
}
