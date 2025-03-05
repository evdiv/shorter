package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"shorter/internal/config"
	"shorter/internal/storage"
	"shorter/internal/urlkey"
	"strings"
)

// Handlers struct holds dependencies (storage)
type Handlers struct {
	Storage storage.Storer
}

type JSONReq struct {
	URL         string `json:"url,omitempty"`
	CorrID      string `json:"correlation_id,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

type JSONRes struct {
	Result      string `json:"result,omitempty"`
	CorrID      string `json:"correlation_id,omitempty"`
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"-"`
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
		if strings.Contains(err.Error(), "already exists") {
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
	var jReq JSONReq
	var jRes JSONRes
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
		if strings.Contains(err.Error(), "already exists") {
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

func (h *Handlers) ShortenBatchURL(res http.ResponseWriter, req *http.Request) {
	jReqBatch := []JSONReq{}
	jResBatch := []JSONRes{}

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
			shortURL := config.AppConfig.ResultHost + "/" + urlKey
			jResBatch = append(jResBatch, JSONRes{CorrID: jReq.CorrID, ShortURL: shortURL, OriginalURL: jReq.OriginalURL})
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
