package tools

import (
	"strings"

	"kgent/cmd/utils"
)

type ListToolParam struct {
	Resource  string `json:"resource"`
	Namespace string `json:"namespace"`
}

// ListTool represents a tool that lists Kubernetes resources in a specified namespace.
type ListTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewListTool creates a new ListTool instance.
func NewListTool() *ListTool {
	return &ListTool{
		Name:        "ListTool",
		Description: "Used to list the specified Kubernetes resources in a specified namespace, such as pod list etc.",
		ArgsSchema:  `{"type":"object","properties":{"resource":{"type":"string", "description": "The specified Kubernetes resource type, such as pod, service etc."}, "namespace":{"type":"string", "description": "The specified Kubernetes namespace"}}`,
	}
}

// Run executes the command and returns the output.
func (l *ListTool) Run(resource string, ns string) (string, error) {
	// Set default namespace if not provided
	if ns == "" {
		ns = "default"
	}

	resource = strings.ToLower(resource)

	// Get the API URL from environment variable with fallback
	apiURL := utils.GetEnv("KGENT_API_URL", "http://localhost:8000/api/v1/resources")

	// Ensure URL ends with a slash
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}

	url := apiURL + resource + "?ns=" + ns

	s, err := utils.GetHTTP(url)

	return s, err
}
