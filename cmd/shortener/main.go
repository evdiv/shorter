package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// Config
const (
	host = "http://localhost"
	port = ":8080"
)

// Storage for the compressed URLs
var urlStore = make(map[string]string)

type MainHandler struct{}

func (h MainHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	log.Println("Request received: ", req)

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

// generateKey generates the corresponding key for the stored URL
func generateKey(url string) string {
	var sum int
	for i, char := range url {
		// Multiply ASCII value by position (1-based)
		sum += int(char) * (i + 1)
	}
	return fmt.Sprintf("%x", sum)
}

// store stores the URL in the storage nad returns the corresponding key
func store(url string) string {
	key := generateKey(url)

	if urlStore[key] == "" {
		urlStore[key] = url
	}
	return key
}

// retrieve extracts the stored URL by its key
func retrieve(key string) string {
	return urlStore[key]
}

func main() {
	var h MainHandler

	log.Println("Starting the server on port", port)
	err := http.ListenAndServe(port, h)
	if err != nil {
		panic(err)
	}
}
