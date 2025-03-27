package storage

import (
	"errors"
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

func (m *MemoryStorage) DeleteBatch(keysToDelete []models.KeysToDelete) (bool, error) {
	if len(keysToDelete) == 0 {
		return false, errors.New("no URLs provided for deletion")
	}

	// Flag that indicates if any record was deleted
	deleted := false

	// Iterate over each KeysToDelete entry
	for _, item := range keysToDelete {
		for _, key := range item.Keys {
			urlKey := strings.ToLower(key)
			if existing, found := m.data[urlKey]; found && existing[1] == item.UserID {
				// Delete the record
				delete(m.data, urlKey)
				deleted = true
			}
		}
	}

	// Return true if at least one record was deleted
	return deleted, nil
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
