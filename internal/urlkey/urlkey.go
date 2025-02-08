package urlkey

import (
	"fmt"
	"net/url"
	"strings"
)

// IsValidURL - validates the url
func IsValidURL(u string) (string, bool) {
	parsedURL, err := url.Parse(u)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return u, false
	}
	return u, true
}

// Generate - generates the corresponding key for the stored URL
func GenerateSlug(originalURL string) string {
	u := strings.ToLower(strings.TrimSpace(originalURL))

	if len(u) == 0 {
		return ""
	}

	var sum int
	for i, char := range u {
		// Multiply ASCII value by position (1-based)
		sum += int(char) * (i + 1)
	}
	return fmt.Sprintf("%x", sum)
}
