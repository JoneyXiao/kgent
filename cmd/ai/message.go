package ai

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"

	promptTpl "kgent/cmd/prompt"
	"kgent/cmd/utils"
)

// Global variables to store configuration and message history
var MessageStore ChatMessages
var ModelName string
var Token string
var DashScopeURL string

// Role constants for chat messages
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
	RoleTool      = "tool"
)

// ChatMessages represents a collection of chat messages
type ChatMessages []*ChatMessage

// ChatMessage represents a single message in the chat
type ChatMessage struct {
	Msg openai.ChatCompletionMessage
}

// Clear initializes or resets the chat history
func (cm *ChatMessages) Clear() {
	*cm = make([]*ChatMessage, 0)
	cm.AddSystem(promptTpl.SystemPrompt)
}

func init() {
	// Load .env file if it exists
	workDir, _ := os.Getwd()
	envPath := filepath.Join(workDir, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables directly")
	}

	// Set environment variables
	Token = utils.GetEnv("DASH_SCOPE_API_KEY", "")
	ModelName = utils.GetEnv("DASH_SCOPE_MODEL", "qwen-max")
	DashScopeURL = utils.GetEnv("DASH_SCOPE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1")

	// Check for missing configuration
	if Token == "" {
		log.Println("Error: DASH_SCOPE_API_KEY is not set. Please set this environment variable.")
		os.Exit(1)
	}

	// Initialize message store
	MessageStore = make(ChatMessages, 0)
	MessageStore.Clear()
}

// NewOpenAiClient creates a new client with the configured API key and URL
func NewOpenAiClient() *openai.Client {
	config := openai.DefaultConfig(Token)
	config.BaseURL = DashScopeURL

	return openai.NewClientWithConfig(config)
}

// AppendMessage appends a message with the specified role
func (cm *ChatMessages) AppendMessage(msg string, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:    role,
			Content: msg,
		},
	})
}

// GetMessage converts the chat history to a format suitable for the API
func (cm *ChatMessages) GetMessage() []openai.ChatCompletionMessage {
	ret := make([]openai.ChatCompletionMessage, len(*cm))
	for i, msg := range *cm {
		ret[i] = msg.Msg
	}
	return ret
}

// AddToolCall adds a tool call to the chat history
func (cm *ChatMessages) AddToolCall(rsp openai.ChatCompletionMessage, role string) {
	*cm = append(*cm, &ChatMessage{
		Msg: openai.ChatCompletionMessage{
			Role:         role,
			Content:      rsp.Content,
			FunctionCall: rsp.FunctionCall,
			ToolCalls:    rsp.ToolCalls,
		},
	})
}

// AddSystem adds a system message
func (cm *ChatMessages) AddSystem(msg string) {
	cm.AppendMessage(msg, RoleSystem)
}

// AddAssistant adds an assistant message
func (cm *ChatMessages) AddAssistant(rsp openai.ChatCompletionMessage) {
	cm.AddToolCall(rsp, RoleAssistant)
}

// AddUser adds a user message
func (cm *ChatMessages) AddUser(msg string) {
	cm.AppendMessage(msg, RoleUser)
}

// GetLast returns the last message in the chat history
func (cm *ChatMessages) GetLast() string {
	if len(*cm) == 0 {
		return "No messages in the chat history"
	}

	return (*cm)[len(*cm)-1].Msg.Content
}

// Chat sends a message to the AI API and returns the response
func Chat(message []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c := NewOpenAiClient()
	rsp, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    ModelName,
		Messages: message,
	})

	if err != nil {
		log.Printf("Error calling AI API: %v\n", err)
		return openai.ChatCompletionMessage{
			Role:    RoleAssistant,
			Content: fmt.Sprintf("Sorry, I encountered an error when processing your request: %v", err),
		}
	}

	if len(rsp.Choices) == 0 {
		log.Println("Error: No response choices received from API")
		return openai.ChatCompletionMessage{
			Role:    RoleAssistant,
			Content: "Sorry, I received an empty response. Please try again.",
		}
	}

	return rsp.Choices[0].Message
}
