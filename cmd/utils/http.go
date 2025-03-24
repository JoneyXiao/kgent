package utils

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

// HTTP client with timeout
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Health check
func init() {
	healthCheckURL := GetEnv("KGENT_API_URL", "http://localhost:8000/health")
	_, err := GetHTTP(healthCheckURL)
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
}

// GetHTTP executes a GET HTTP request to the specified URL and returns the response body.
func GetHTTP(url string) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create HTTP GET request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// PostHTTP executes a POST HTTP request to the specified URL with the given body and returns the response.
func PostHTTP(url string, body []byte) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create HTTP POST request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// DeleteHTTP executes a DELETE HTTP request to the specified URL and returns the response.
func DeleteHTTP(url string) (string, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create HTTP DELETE request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return "", err
	}

	// Execute request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
