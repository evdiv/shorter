package storage

import (
	"fmt"
	"shorter/internal/models"
	"shorter/internal/urlkey"
	"strings"
)

type MemoryStorage struct {
	data map[string][]string
}

// NewMemoryStorage - constructor to create a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string][]string)}
}

// Set - stores a url into the memory storage
func (m *MemoryStorage) Set(OriginalURL string, userID string) (string, error) {
	urlKey := urlkey.GenerateSlug(OriginalURL)
	if urlKey == "" {
		return "", fmt.Errorf("ShortURL is empty")
	}
	existing, found := m.data[urlKey]

	if found && existing[0] != "" {
		err := fmt.Errorf("the URL: %s is already stored in the memory", m.data[urlKey])
		return urlKey, NewStorageError("already exists", OriginalURL, urlKey, err)
	}
	m.data[urlKey] = []string{OriginalURL, userID}
	return urlKey, nil
}

func (m *MemoryStorage) SetBatch(jReqBatch []models.JSONReq, userID string) ([]models.JSONRes, error) {
	jResBatch := make([]models.JSONRes, len(jReqBatch))

	for _, el := range jReqBatch {
		ShortURL, err := m.Set(el.OriginalURL, userID)
		if err != nil {
			return nil, err
		}

		row := models.JSONRes{
			CorrID:      el.CorrID,
			ShortURL:    ShortURL,
			OriginalURL: el.OriginalURL,
		}
		jResBatch = append(jResBatch, row)
	}
	return jResBatch, nil
}

// Get - retrieves a value from memory
func (m *MemoryStorage) Get(urlKey string) (string, error) {
	urlKey = strings.ToLower(urlKey)

	existing, found := m.data[urlKey]
	if !found || existing[0] == "" {
		return "", fmt.Errorf("OriginalURL is empty")
	}
	return existing[0], nil
}

func (m *MemoryStorage) GetUserURLs(userID string) ([]models.JSONUserRes, error) {
	jResBatch := make([]models.JSONUserRes, 0)

	for key, el := range m.data {
		if el[1] != userID {
			continue
		}
		row := models.JSONUserRes{
			ShortURL:    key,
			OriginalURL: el[0],
		}
		jResBatch = append(jResBatch, row)
	}
	return jResBatch, nil
}

func (m *MemoryStorage) IsAvailable() bool {
	return m.data != nil
}

// Close - ensure that the in memory storage fits the Storer interface
func (m *MemoryStorage) Close() error {
	return nil
}
