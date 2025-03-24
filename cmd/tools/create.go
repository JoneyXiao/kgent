package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"kgent/cmd/ai"
	promptTpl "kgent/cmd/prompt"
	"kgent/cmd/utils"

	"github.com/sashabaranov/go-openai"
)

type CreateToolParam struct {
	Prompt   string `json:"prompt"`
	Resource string `json:"resource"`
}

// define struct to parse JSON response
type response struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

// CreateTool represents a tool that creates a specified Kubernetes resource in a specified namespace.
type CreateTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewCreateTool creates a new CreateTool instance.
func NewCreateTool() *CreateTool {
	return &CreateTool{
		Name:        "CreateTool",
		Description: "Used to create a specified Kubernetes resource in a specified namespace, such as creating a pod etc.",
		ArgsSchema:  `{"type":"object","properties":{"prompt":{"type":"string", "description": "Put the user's prompt for creating a resource exactly here, without any changes"},"resource":{"type":"string", "description": "The specified Kubernetes resource type, such as pod, service etc."}}}`,
	}
}

// Run executes the command and returns the output.
func (c *CreateTool) Run(prompt string, resource string, debugMode bool) string {
	// let the large model generate yaml
	messages := make([]openai.ChatCompletionMessage, 2)

	messages[0] = openai.ChatCompletionMessage{Role: "system", Content: promptTpl.K8sAssistantPrompt}
	messages[1] = openai.ChatCompletionMessage{Role: "user", Content: prompt}

	rsp := ai.Chat(messages)

	// remove ```yaml and ``` from the response
	rsp.Content = strings.Replace(rsp.Content, "```yaml", "", -1)
	rsp.Content = strings.Replace(rsp.Content, "```", "", -1)
	rsp.Content = strings.TrimSpace(rsp.Content)

	// create JSON object {"yaml":"xxx"}
	body := map[string]string{"yaml": rsp.Content}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}

	// Get the API URL from environment variable with fallback
	apiURL := utils.GetEnv("KGENT_API_URL", "http://localhost:8000/api/v1/resources")

	// Ensure URL ends with a slash
	if !strings.HasSuffix(apiURL, "/") {
		apiURL += "/"
	}

	url := apiURL + resource
	s, err := utils.PostHTTP(url, jsonBody)
	if err != nil {
		return err.Error()
	}

	var response response
	// parse JSON response
	err = json.Unmarshal([]byte(s), &response)
	if err != nil {
		return err.Error()
	}

	if debugMode {
		fmt.Println(rsp.Content)
		fmt.Println("[CreateTool] jsonBody", string(jsonBody))
		fmt.Println("[CreateTool] url", url)
		fmt.Println("[CreateTool] response", response)
	}
	// return error if response.Data is empty
	if response.Data == "" {
		return response.Error
	}

	return response.Data
}
