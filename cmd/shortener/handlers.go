package main

import (
	"io"
	"log"
	"net/http"
)

type MainHandler struct{}

func (h MainHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		// Read the raw body since the client sends plain text
		body, err := io.ReadAll(req.Body)

		if err != nil {
			log.Println("Failed to read request body: ", err)
			res.WriteHeader(400)
			return
		}
		defer req.Body.Close()

		originalUrl := string(body)

		if len(originalUrl) == 0 {
			log.Println("URL should't be empty. Return error / 400 ")
			res.WriteHeader(400)
			return
		}

		urlKey := store(originalUrl)
		log.Println("For URL: ", originalUrl, "the urlKey: ", urlKey, " is generated")

		res.WriteHeader(201)
		res.Write([]byte(host + port + "/" + urlKey))
		return
	}

	if req.Method == http.MethodGet {
		urlKey := req.URL.Path

		log.Println("Received GET request for urlKey: ", urlKey)

		//No incoming parameters
		if len(urlKey) < 2 {
			log.Println("urlKey should't be empty. Return error / 400 ")
			res.WriteHeader(400)
			return
		}

		originalUrl := retrieve(urlKey[1:])

		log.Println("Found stored originalUrl: ", originalUrl)

		//The URL is not found in the storage
		if len(originalUrl) == 0 {
			res.WriteHeader(400)
			return
		}
		// Set the Location header
		res.Header().Set("Location", originalUrl)
		res.WriteHeader(307)
		return
	}
	// No other methods are allowed
	res.WriteHeader(400)
	return
}
