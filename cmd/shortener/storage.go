package main

import "fmt"

// Storage for the compressed URLs
var urlStore = make(map[string]string)

// generateKey generates the corresponding key for the stored URL
func generateKey(url string) string {
	var sum int
	for i, char := range url {
		// Multiply ASCII value by position (1-based)
		sum += int(char) * (i + 1)
	}
	return fmt.Sprintf("%x", sum)
}

// store stores the URL in the storage nad returns the corresponding key
func store(url string) string {
	key := generateKey(url)

	if urlStore[key] == "" {
		urlStore[key] = url
	}
	return key
}

// retrieve extracts the stored URL by its key
func retrieve(key string) string {
	return urlStore[key]
}
