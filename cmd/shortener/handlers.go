package main

import (
	"io"
	"net/http"
)

type MainHandler struct{}

func (h MainHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		// Read the raw body since the client sends plain text
		body, err := io.ReadAll(req.Body)

		if err != nil {
			res.WriteHeader(400)
			return
		}
		defer req.Body.Close()

		originalUrl := string(body)

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
		// Set the Location header
		res.Header().Set("Location", originalUrl)
		res.WriteHeader(307)
		return
	}
	// No other methods are allowed
	res.WriteHeader(400)
	return
}
