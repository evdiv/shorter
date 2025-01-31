package main

import (
	"fmt"
	"strings"
)

var urlStore = make(map[string]string)

// generateKey - generates the corresponding key for the stored URL
func generateKey(url string) string {
	if len(url) == 0 {
		return ""
	}

	var sum int
	for i, char := range url {
		// Multiply ASCII value by position (1-based)
		sum += int(char) * (i + 1)
	}
	return fmt.Sprintf("%x", sum)
}

// store - stores the URL in the storage and returns the corresponding key
func store(originalURL string) string {
	originalURL = strings.ToLower(originalURL)
	urlKey := generateKey(originalURL)

	if urlKey == "" {
		return ""
	}
	urlStore[urlKey] = originalURL
	return urlKey
}

// retrieve extracts the stored URL by its key
func retrieve(urlKey string) string {
	urlKey = strings.ToLower(urlKey)
	return urlStore[urlKey]
}
