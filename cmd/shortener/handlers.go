package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func PostURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}
	defer req.Body.Close()
	originalURL := string(body)

	if originalURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain URL"))
		return
	}

	urlKey := store(originalURL)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(host + port + "/" + urlKey))
}

func GetURL(res http.ResponseWriter, req *http.Request) {
	urlKey := chi.URLParam(req, "urlKey")

	if urlKey == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing URL key"))
		return
	}

	originalURL := retrieve(urlKey)

	if originalURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("URL is not found"))
		return
	}

	// Set the Location header and return a 307 Temporary Redirect
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
