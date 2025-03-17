package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/middleware"
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

func getUserIDFromContext(req *http.Request) (string, error) {
	userID, ok := req.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		return "", errors.New("userID not found in context")
	}

	return userID, nil
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

	userID, _ := getUserIDFromContext(req)

	urlKey, err := h.Storage.Set(originalURL, userID)
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
		res.Write([]byte("The incoming JSON string should contain a valid URL"))
		return
	}
	userID, _ := getUserIDFromContext(req)
	log.Println("In the ShortenURL handler. UserID: " + userID)

	urlKey, err := h.Storage.Set(jReq.URL, userID)
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

func (h *Handlers) GetUserURL(res http.ResponseWriter, req *http.Request) {

	userID, err := getUserIDFromContext(req)

	log.Println("in the GetUserURL. UserID: " + userID)

	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	jResBatch, err := h.Storage.GetUserURLs(userID)
	if err != nil {
		http.Error(res, "No content", http.StatusNoContent)
		return
	}

	out, err := json.Marshal(jResBatch)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(out))
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

// ShortenBatchURL - inserts batch records in storage
func (h *Handlers) ShortenBatchURL(res http.ResponseWriter, req *http.Request) {
	jReqBatch := []models.JSONReq{}

	if err := json.NewDecoder(req.Body).Decode(&jReqBatch); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	userID, _ := getUserIDFromContext(req)
	jResBatch, err := h.Storage.SetBatch(jReqBatch, userID)

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
