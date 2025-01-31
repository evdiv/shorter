package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func PostUrl(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}
	defer req.Body.Close()
	originalUrl := string(body)

	if originalUrl == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain URL"))
		return
	}

	urlKey := store(originalUrl)
	res.WriteHeader(201)
	res.Write([]byte(host + port + "/" + urlKey))
	return
}

func GetUrl(res http.ResponseWriter, req *http.Request) {
	urlKey := chi.URLParam(req, "urlKey")

	if urlKey == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing URL key"))
		return
	}

	originalUrl := retrieve(urlKey)

	if len(originalUrl) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("URL is not found"))
		return
	}

	// Set the Location header
	res.WriteHeader(307)
	res.Header().Set("Location", originalUrl)
	return
}
