package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Config
const (
	host = "http://localhost"
	port = ":8080"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", PostUrl)
	//Needed for avoiding 405 Error according to the requirements.
	r.Get("/", GetUrl)
	r.Get("/{urlKey}", GetUrl)

	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}
