package tools

import (
	"fmt"
	"os/exec"
	"strings"
)

// KubeInput represents the input for the KubeTool.
type KubeInput struct {
	Commands string `json:"commands"`
}

// KubeTool represents a tool that runs Kubernetes commands.
type KubeTool struct {
	Name        string
	Description string
	ArgsSchema  string
}

// NewKubeTool creates a new KubeTool instance.
func NewKubeTool() *KubeTool {
	return &KubeTool{
		Name:        "KubeTool",
		Description: "A tool for running Kubernetes commands (kubectl, helm) on a Kubernetes cluster.",
		ArgsSchema:  `{"type":"object","properties":{"commands":{"type":"string", "description": "The kubectl/helm related command to run. e.g. kubectl get pods"}}}`,
	}
}

// Run executes the command and returns the output.
func (k *KubeTool) Run(commands string) (string, error) {
	parsedCommands := k.parseCommands(commands)

	splitedCommands := k.splitCommands(parsedCommands)
	fmt.Println("splitedCommands", splitedCommands)
	// You usually use the os/exec package to execute the command and return the output.
	cmd := exec.Command(splitedCommands[0], splitedCommands[1:]...)

	// Run the command and get the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}

	return fmt.Sprintf("The result of the command execution: %s", output), nil
}

// parseCommands cleans the command string.
func (k *KubeTool) parseCommands(commands string) string {
	return strings.TrimSpace(strings.Trim(commands, "\"`"))
}

// splitCommands splits the command string.
func (k *KubeTool) splitCommands(commands string) []string {
	return strings.Split(commands, " ")
}
