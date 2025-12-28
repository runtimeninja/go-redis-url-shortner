package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// RemoveDomainError blocks urls that match our own domain
func RemoveDomainError(url string) bool {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return true 
	}

	if strings.EqualFold(url, domain) {
		return false
	}

	// remove scheme and www for comparison
	newURL := strings.ToLower(url)
	newURL = strings.TrimPrefix(newURL, "http://")
	newURL = strings.TrimPrefix(newURL, "https://")
	newURL = strings.TrimPrefix(newURL, "www.")

	parts := strings.Split(newURL, "/")
	if len(parts) > 0 && parts[0] == strings.ToLower(domain) {
		return false
	}

	return true
}
