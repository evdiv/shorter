package main

import (
	"fmt"
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

	if req.Method == http.MethodPost {
		originalUrl := req.FormValue("URL")

		if len(originalUrl) == 0 {
			res.WriteHeader(400)
			return
		}

		urlKey := store(originalUrl)

		res.WriteHeader(201)
		res.Write([]byte(host + port + "/" + urlKey))
		return
	}

	if req.Method == http.MethodGet {
		urlKey := req.URL.Path

		//No incoming parameters
		if len(urlKey) < 2 {
			res.WriteHeader(400)
			return
		}

		originalUrl := retrieve(urlKey[1:])

		//The URL is not found in the storage
		if len(originalUrl) == 0 {
			res.WriteHeader(400)
			return
		}

		res.WriteHeader(307)
		res.Write([]byte(originalUrl))
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

	err := http.ListenAndServe(port, h)
	if err != nil {
		panic(err)
	}
}
