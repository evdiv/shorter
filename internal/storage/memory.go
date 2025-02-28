package storage

import (
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
func (m *MemoryStorage) Set(url string) string {
	key := urlkey.GenerateSlug(url)
	if key == "" {
		return ""
	}
	m.data[key] = url
	return key
}

// Get - retrieves a value from memory
func (m *MemoryStorage) Get(key string) string {
	key = strings.ToLower(key)
	return m.data[key]
}

func (m *MemoryStorage) IsAvailable() bool {
	return m.data != nil
}
