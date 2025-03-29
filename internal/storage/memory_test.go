package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage_Set(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	originalURL := "https://practicum.yandex.ru/"
	userID := "111222333abc"
	key, _ := storage.Set(ctx, originalURL, userID)
	assert.NotEmpty(t, key, "Expected a non-empty key, got an empty string")

	retrievedURL, _ := storage.Get(ctx, key)
	assert.Equal(t, originalURL, retrievedURL, "Stored and retrieved URLs should match")
}

func TestMemoryStorage_Get_NonExistentKey(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryStorage()

	nonExistentKey := "random"
	result, err := storage.Get(ctx, nonExistentKey)

	assert.NotEmpty(t, err, "Expected non-empty error")
	assert.Empty(t, result, "Expected empty string for non-existent key")
}
