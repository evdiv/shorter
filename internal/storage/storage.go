package storage

import (
	"shorter/internal/urlkey"
	"strings"
)

var urlStore = make(map[string]string)

// store - stores the URL in the storage and returns the corresponding key
func Store(originalURL string) string {
	originalURL = strings.ToLower(originalURL)
	urlKey := urlkey.GenerateSlug(originalURL)

	if urlKey == "" {
		return ""
	}
	urlStore[urlKey] = originalURL
	return urlKey
}

// retrieve extracts the stored URL by its key
func Retrieve(urlKey string) string {
	urlKey = strings.ToLower(urlKey)
	return urlStore[urlKey]
}
