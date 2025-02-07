package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/storage"
	"shorter/internal/urlkey"
)

// NewRouter initializes the router and storage, returning the chi router
func NewRouter() *chi.Mux {

	storer := storage.NewMemoryStorage()

	r := chi.NewRouter()

	r.Post("/", func(res http.ResponseWriter, req *http.Request) {
		PostURL(res, req, storer) // Pass storer to PostURL
	})
	r.Get("/{urlKey}", func(res http.ResponseWriter, req *http.Request) {
		GetURL(res, req, storer) // Pass storer to GetURL
	})

	// Fallback for empty key
	r.Get("/", func(res http.ResponseWriter, req *http.Request) {
		GetURL(res, req, storer) // Pass storer to GetURL
	})

	return r
}

// Update PostURL to accept a storage interface
func PostURL(res http.ResponseWriter, req *http.Request, storer storage.Storer) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}
	defer req.Body.Close()

	url, valid := urlkey.IsValidURL(string(body))

	if !valid {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain valid URL"))
		return
	}

	urlKey := storer.Set(url)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.GetHost("Result") + "/" + urlKey))
}

// Update GetURL to accept a storage interface
func GetURL(res http.ResponseWriter, req *http.Request, storer storage.Storer) {
	urlKey := chi.URLParam(req, "urlKey")

	if urlKey == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing URL key"))
		return
	}

	originalURL := storer.Get(urlKey)

	if originalURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("URL is not found"))
		return
	}

	// Set the Location header and return a 307 Temporary Redirect
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
