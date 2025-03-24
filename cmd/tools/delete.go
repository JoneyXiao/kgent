package tools

import (
	"strings"

	"kgent/cmd/utils"
)

type DeleteToolParam struct {
	Resource  string `json:"resource"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// DeleteTool represents a tool that deletes a specified Kubernetes resource in a specified namespace.
type DeleteTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewDeleteTool creates a new DeleteTool instance.
func NewDeleteTool() *DeleteTool {
	return &DeleteTool{
		Name:        "DeleteTool",
		Description: "Used to delete a specified Kubernetes resource in a specified namespace, such as deleting a pod etc.",
		ArgsSchema:  `{"type":"object","properties":{"resource":{"type":"string", "description": "The specified Kubernetes resource type, such as pod, service etc."}, "name":{"type":"string", "description": "The name of the specified Kubernetes resource instance"}, "namespace":{"type":"string", "description": "The namespace of the specified Kubernetes resource"}}`,
	}
}

// Run executes the command and returns the output.
func (d *DeleteTool) Run(resource, name, ns string) error {
	resource = strings.ToLower(resource)

	// Get the API URL from environment variable with fallback
	apiURL := utils.GetEnv("KGENT_API_URL", "http://localhost:8000/api/v1/resources")

	// Ensure URL ends with a slash
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}

	url := apiURL + resource + "?ns=" + ns + "&name=" + name

	_, err := utils.DeleteHTTP(url)

	return err
}
