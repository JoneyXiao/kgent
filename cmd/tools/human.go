package tools

import (
	"fmt"

	"kgent/cmd/utils"
)

type HumanToolParam struct {
	Prompt string `json:"prompt"`
}

// HumanTool represents a tool that asks for human confirmation before performing dangerous operations.
type HumanTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewHumanTool creates a new HumanTool instance.
func NewHumanTool() *HumanTool {
	return &HumanTool{
		Name:        "HumanTool",
		Description: "When you determine that you need to perform dangerous operations, such as deletion, you need to use this tool to initiate a confirmation request to humans first.",
		ArgsSchema:  `{"type":"object","properties":{"prompt":{"type":"string", "description": "The action you want to perform, such as deleting a pod", "example": "Please confirm whether to delete the foo-app pod in the default namespace"}}}`,
	}
}

// Run executes the command and returns the output.
func (h *HumanTool) Run(prompt string) string {
	utils.PrintYellowNoNewline(prompt + " (yes/no): ")
	var input string
	fmt.Scanln(&input)

	// Provide more context in the response
	if input == "y" || input == "yes" {
		return "Human confirmed! Do I need to use a tool? Yes"
	} else if input == "n" || input == "no" {
		return "Human declined! Do I need to use a tool? No"
	} else {
		return "Human response: " + input
	}
}
