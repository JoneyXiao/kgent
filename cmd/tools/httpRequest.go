package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// RequestsToolParam represents the input for the RequestsTool
type RequestsToolParam struct {
	Url string `json:"url"`
}

// RequestsTool provides functionality to make HTTP requests and process the responses.
// It extracts text content from HTML responses.
type RequestsTool struct {
	Name        string
	Description string
	ArgsSchema  string
	client      *http.Client
}

// NewRequestsTool creates and returns a new RequestsTool with default configuration.
func NewRequestsTool() *RequestsTool {
	return &RequestsTool{
		Name:        "RequestsTool",
		Description: `A portal to the internet. Use this when you need to get specific content from a website. Input should be a url (i.e. https://www.kubernetes.io/releases). The output will be the text response of the GET request.`,
		ArgsSchema:  `{"type":"object","properties":{"url":{"type":"string", "description": "the url to be accessed, e.g. https://www.kubernetes.io/releases"}}}`,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run makes a GET request to the specified URL and returns the text content.
// It uses context with timeout for proper cancellation support.
func (r *RequestsTool) Run(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	text, err := r.parseHTML(string(body))
	if err != nil {
		return string(body), fmt.Errorf("failed to parse HTML (returning raw content): %w", err)
	}

	return text, nil
}

// parseHTML processes HTML content and extracts clean text.
// It removes scripts, styles, headers, and footers before extracting text.
func (r *RequestsTool) parseHTML(htmlContent string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML document: %w", err)
	}

	// Remove unwanted tags
	doc.Find("header, footer, script, style, nav, iframe, noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// Get the pure text after processing
	text := strings.TrimSpace(doc.Find("body").Text())
	return text, nil
}

// RunWithHeaders makes a GET request with custom headers to the specified URL.
// This is a more flexible version of Run that allows setting HTTP headers.
func (r *RequestsTool) RunWithHeaders(url string, headers map[string]string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d - %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	text, err := r.parseHTML(string(body))
	if err != nil {
		return string(body), fmt.Errorf("failed to parse HTML (returning raw content): %w", err)
	}

	return text, nil
}
