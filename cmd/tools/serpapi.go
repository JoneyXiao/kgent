package tools

import (
	"encoding/json"
	"fmt"
	"os"

	g "github.com/serpapi/google-search-results-golang"
)

// SerpApiToolParam represents the input for the SerpApiTool
type SerpApiToolParam struct {
	Query string `json:"query"`
}

// SerpApiTool represents a tool for searching the web via SerpAPI
type SerpApiTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

type FinalResult struct {
	Title string
	Link  string
}

// NewSerpApiTool creates a new instance of SerpApiTool
func NewSerpApiTool() *SerpApiTool {
	return &SerpApiTool{
		Name:        "serpapi_search",
		Description: "Search the web for information using DuckDuckGo search engine via SerpAPI",
		ArgsSchema:  `{"type":"object","properties":{"query":{"type":"string", "description": "the search query to be used"}}}`,
	}
}

// GetName returns the tool name
func (t *SerpApiTool) GetName() string {
	return t.Name
}

// GetDescription returns the tool description
func (t *SerpApiTool) GetDescription() string {
	return t.Description
}

// ToJSON converts the tool to JSON format
func (t *SerpApiTool) ToJSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ExecuteSearch performs a web search using SerpAPI
func (t *SerpApiTool) ExecuteSearch(query string, engine string) (map[string]interface{}, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("SERPAPI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SERPAPI_API_KEY environment variable not set")
	}

	// Set search parameters
	parameter := map[string]string{
		"engine": engine,
		"q":      query,
		"kl":     "us-en",
	}

	// Create search client and execute search
	search := g.NewGoogleSearch(parameter, apiKey)
	results, err := search.GetJSON()
	if err != nil {
		return nil, fmt.Errorf("SerpAPI search failed: %v", err)
	}

	return results, nil
}

// Run executes the tool with the given arguments
func (t *SerpApiTool) Run(query string) ([]FinalResult, error) {
	// Extract search query from arguments
	if query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	// Get search engine (default to bing if not specified)
	engine := "duckduckgo"

	// Execute search
	results, err := t.ExecuteSearch(query, engine)
	if err != nil {
		return nil, err
	}

	// get all the title and link
	finalResult := make([]FinalResult, 0)
	organicResults := results["organic_results"].([]interface{})
	for _, result := range organicResults {
		finalResult = append(finalResult, FinalResult{
			Title: "title: " + result.(map[string]interface{})["title"].(string),
			Link:  "link: " + result.(map[string]interface{})["link"].(string),
		})
	}

	// Format results as JSON string
	// resultBytes, err := json.MarshalIndent(results, "", "  ")
	// if err != nil {
	// 	return "", fmt.Errorf("failed to marshal search results: %v", err)
	// }

	return finalResult, nil
}
