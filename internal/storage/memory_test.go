package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage_Set(t *testing.T) {

	storage := NewMemoryStorage()

	originalURL := "https://practicum.yandex.ru/"
	key := storage.Set(originalURL)
	assert.NotEmpty(t, key, "Expected a non-empty key, got an empty string")

	retrievedURL := storage.Get(key)
	assert.Equal(t, originalURL, retrievedURL, "Stored and retrieved URLs should match")
}

func TestMemoryStorage_Get_NonExistentKey(t *testing.T) {

	storage := NewMemoryStorage()

	nonExistentKey := "random"
	result := storage.Get(nonExistentKey)

	assert.Empty(t, result, "Expected empty string for non-existent key")
}
