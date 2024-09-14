// pkg/api/client.go
package api

import (
	"net/http"
	"net/http/cookiejar"
	"nhs-bank-notifier/pkg/logger"
)

// NewClient creates a new HTTP client with a cookie jar
func NewClient() *http.Client {
	log := logger.GetLogger()
	log.Debug("Creating new HTTP client...")

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	return &http.Client{Jar: jar}
}
