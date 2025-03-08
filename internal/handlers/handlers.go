package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/models"
	"shorter/internal/storage"
	"shorter/internal/urlkey"
)

// Handlers struct holds dependencies (storage)
type Handlers struct {
	Storage storage.Storer
}

// NewHandlers initializes handlers with storage
func NewHandlers(s storage.Storer) *Handlers {
	return &Handlers{Storage: s}
}

func (h *Handlers) PostURL(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}
	defer req.Body.Close()

	originalURL, valid := urlkey.IsValidURL(string(body))

	if !valid {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain a valid URL"))
		return
	}

	urlKey, err := h.Storage.Set(originalURL)
	HeaderStatus := http.StatusCreated

	if err != nil {
		var storageErr *storage.StorageError
		if errors.As(err, &storageErr) && storageErr.Type == "already exists" {
			HeaderStatus = http.StatusConflict
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	res.WriteHeader(HeaderStatus)
	res.Write([]byte(config.AppConfig.ResultHost + "/" + urlKey))
}

func (h *Handlers) ShortenURL(res http.ResponseWriter, req *http.Request) {
	var jReq models.JSONReq
	var jRes models.JSONRes
	if err := json.NewDecoder(req.Body).Decode(&jReq); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if _, valid := urlkey.IsValidURL(jReq.URL); !valid {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain a valid URL"))
		return
	}

	urlKey, err := h.Storage.Set(jReq.URL)
	HeaderStatus := http.StatusCreated

	if err != nil {
		var storageErr *storage.StorageError
		if errors.As(err, &storageErr) && storageErr.Type == "already exists" {
			HeaderStatus = http.StatusConflict
		} else {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	jRes.Result = config.AppConfig.ResultHost + "/" + urlKey

	out, err := json.Marshal(jRes)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(HeaderStatus)
	res.Write(out)
}

func (h *Handlers) GetURL(res http.ResponseWriter, req *http.Request) {
	urlKey := chi.URLParam(req, "urlKey")

	if urlKey == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing URL key"))
		return
	}

	originalURL, err := h.Storage.Get(urlKey)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(err.Error()))
		return
	}

	// Set the Location header and return a 307 Temporary Redirect
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handlers) IsAvailable(res http.ResponseWriter, req *http.Request) {
	if h.Storage.IsAvailable() {
		res.WriteHeader(http.StatusOK)
	}
	res.WriteHeader(http.StatusInternalServerError)
}

// ShortenBatchURL - inserts records line by line
func (h *Handlers) ShortenBatchURL(res http.ResponseWriter, req *http.Request) {
	var jReqBatch []models.JSONReq
	var jResBatch []models.JSONRes

	if err := json.NewDecoder(req.Body).Decode(&jReqBatch); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	for _, jReq := range jReqBatch {
		if _, valid := urlkey.IsValidURL(jReq.OriginalURL); !valid {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("The batch contains incorrect URL: " + jReq.URL))
			return
		}

		urlKey, err := h.Storage.Set(jReq.OriginalURL)

		if err == nil {
			row := models.JSONRes{
				CorrID:      jReq.CorrID,
				ShortURL:    config.AppConfig.ResultHost + "/" + urlKey,
				OriginalURL: jReq.OriginalURL,
			}
			jResBatch = append(jResBatch, row)
		}
	}

	out, err := json.Marshal(jResBatch)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(out))
}

// ShortenBatchURLInBulk - inserts records in bulk with a single query
func (h *Handlers) ShortenBatchURLInBulk(res http.ResponseWriter, req *http.Request) {
	jReqBatch := []models.JSONReq{}

	if err := json.NewDecoder(req.Body).Decode(&jReqBatch); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	jResBatch, err := h.Storage.SetBatch(jReqBatch)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(jResBatch)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(out))
}
