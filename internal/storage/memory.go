package storage

import (
	"fmt"
	"shorter/internal/urlkey"
	"strings"
)

type MemoryStorage struct {
	data map[string]string
}

// NewMemoryStorage - constructor to create a new MemoryStorage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

// Set - stores a url into the memory storage
func (m *MemoryStorage) Set(OriginalURL string) (string, error) {
	ShortURL := urlkey.GenerateSlug(OriginalURL)
	if ShortURL == "" {
		return "", fmt.Errorf("ShortURL is empty")
	}
	m.data[ShortURL] = OriginalURL
	return ShortURL, nil
}

// Get - retrieves a value from memory
func (m *MemoryStorage) Get(ShortURL string) (string, error) {
	ShortURL = strings.ToLower(ShortURL)
	OriginalURL := m.data[ShortURL]

	if OriginalURL == "" {
		return "", fmt.Errorf("OriginalURL is empty")
	}
	return OriginalURL, nil
}

func (m *MemoryStorage) IsAvailable() bool {
	return m.data != nil
}

// Close - ensure that the in memory storage fits the Storer interface
func (m *MemoryStorage) Close() error {
	return nil
}
