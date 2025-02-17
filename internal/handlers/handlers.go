package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"shorter/internal/config"
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

	urlKey := h.Storage.Set(originalURL)

	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.GetHost("Result") + "/" + urlKey))
}

func (h *Handlers) ShortenURL(res http.ResponseWriter, req *http.Request) {

	var buf bytes.Buffer
	type JsonReq struct {
		Url string `json:"url"`
	}

	var jReq JsonReq

	type JsonResp struct {
		Result string `json:"result"`
	}

	var jRes JsonResp

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &jReq); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unable to handle the request"))
		return
	}

	originalURL, valid := urlkey.IsValidURL(jReq.Url)

	if !valid {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("The body should contain a valid URL"))
		return
	}

	urlKey := h.Storage.Set(originalURL)
	jRes.Result = config.GetHost("Result") + "/" + urlKey

	jsonRes, err := json.Marshal(jRes)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(jsonRes))
}

func (h *Handlers) GetURL(res http.ResponseWriter, req *http.Request) {
	urlKey := chi.URLParam(req, "urlKey")

	if urlKey == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing URL key"))
		return
	}

	originalURL := h.Storage.Get(urlKey)

	if originalURL == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("URL is not found"))
		return
	}

	// Set the Location header and return a 307 Temporary Redirect
	res.Header().Set("Location", originalURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
