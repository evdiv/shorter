package storage

import (
	"shorter/internal/urlkey"
)

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Set(url string) string {
	key := urlkey.GenerateSlug(url)
	s.data[key] = url
	return key
}

func (s *MemoryStorage) Get(key string) string {
	return s.data[key]
}
