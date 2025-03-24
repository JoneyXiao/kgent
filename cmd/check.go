package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kgent/cmd/ai"
	promptTpl "kgent/cmd/prompt"
	"kgent/cmd/tools"
	"kgent/cmd/utils"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the status of the kubernetes cluster",
	Long:  `A tool to check the status of the kubernetes cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize tools
		kubeTool := tools.NewKubeTool()
		searchTool := tools.NewSerpApiTool()
		requestTool := tools.NewRequestsTool()

		// Get default namespace from flag
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace != "" {
			fmt.Printf("Using namespace: %s\n", namespace)
		}

		// Get debug mode flag
		debugMode, _ := cmd.Flags().GetBool("debug")

		// Get max loops flag
		maxLoops, _ := cmd.Flags().GetInt("max-loops")

		runCheckLoop(cmd, kubeTool, searchTool, requestTool, namespace, debugMode, maxLoops)
	},
}

// runCheckLoop handles the main chat interaction loop
func runCheckLoop(cmd *cobra.Command, kubeTool *tools.KubeTool,
	searchTool *tools.SerpApiTool, requestTool *tools.RequestsTool,
	namespace string, debugMode bool, maxLoops int) {
	scanner := bufio.NewScanner(cmd.InOrStdin())
	utils.PrintCyan("Hello, I'm k8s assistant, how can I help you today? (type 'exit' to quit)")

	for {
		utils.PrintYellowNoNewline("> ")
		if !scanner.Scan() {
			// Check for scanning errors
			if err := scanner.Err(); err != nil {
				utils.PrintRed("Error reading input: %v\n", err)
				return
			}
			break
		}

		input := scanner.Text()
		if input == "" {
			continue // Skip empty inputs
		}
		if input == "exit" {
			utils.PrintGreen("Goodbye!")
			return
		}

		// Add namespace to the input if provided
		if namespace != "" && !regexp.MustCompile(`(?i)namespace`).MatchString(input) {
			input = fmt.Sprintf("%s (in namespace %s)", input, namespace)
		}

		prompt := buildCheckPrompt(kubeTool, searchTool, requestTool, input)
		if debugMode {
			fmt.Println("User prompt:", prompt)
		}
		ai.MessageStore.AddUser(prompt)

		processCheckLoop(kubeTool, searchTool, requestTool, maxLoops, debugMode)
		ai.MessageStore.Clear()
	}
}

// processCheckLoop handles the AI interaction and tool execution
func processCheckLoop(kubeTool *tools.KubeTool, searchTool *tools.SerpApiTool,
	requestTool *tools.RequestsTool, maxLoops int, debugMode bool) {
	loopCount := 1

	for loopCount <= maxLoops {
		if debugMode {
			fmt.Printf("---------------- Response round %d ----------------\n", loopCount)
			printCheckDebugInfo()
		}

		response := ai.Chat(ai.MessageStore.GetMessage())

		if debugMode {
			fmt.Println("# Response from LLM:")
			fmt.Println(response.Content)
			fmt.Println()
		} else {
			// In non-debug mode, don't print intermediate thinking
			if !strings.Contains(response.Content, "Final Answer:") {
				utils.PrintCyan("Thinking...")
			}
		}

		// Check for final answer
		regexPattern := regexp.MustCompile(`Final Answer:\s*(.*)`)
		finalAnswer := regexPattern.FindStringSubmatch(response.Content)
		if len(finalAnswer) > 0 {
			// Print only the final answer in non-debug mode
			if !debugMode {
				utils.PrintCyan(strings.TrimSpace(finalAnswer[1]))
			} else {
				utils.PrintCyan("# Final Answer from LLM:")
				utils.PrintCyan(response.Content)
				utils.PrintCyan("-------------------------------------------------")
			}
			break
		}

		ai.MessageStore.AddAssistant(response)

		// Process action if present
		regexAction := regexp.MustCompile(`Action:\s*(.*?)(?:$|[\n\r])`)
		regexActionInput := regexp.MustCompile(`Action Input:\s*(.*?)(?:$|[\n\r])`)
		action := regexAction.FindStringSubmatch(response.Content)
		actionInput := regexActionInput.FindStringSubmatch(response.Content)

		if len(action) > 1 && len(actionInput) > 1 {
			result := handleCheckAction(kubeTool, searchTool, requestTool,
				action[1], actionInput[1], debugMode)

			// Add the observation as a user message
			observation := "Observation: " + result
			prompt := response.Content + "\n" + observation

			if debugMode {
				fmt.Printf("# Round %d user prompt:\n", loopCount)
				fmt.Println(prompt)
			}

			ai.MessageStore.AddUser(prompt)
		}

		loopCount++
	}

	if loopCount >= maxLoops {
		utils.PrintYellow("Exceeded maximum number of reasoning loops. Stopping execution.")
		return
	}
}

// handleAction executes the appropriate tool based on the action
func handleCheckAction(kubeTool *tools.KubeTool, searchTool *tools.SerpApiTool,
	requestTool *tools.RequestsTool,
	action string, actionInput string, debugMode bool) string {
	if debugMode {
		fmt.Println("# Action Debug:")
		fmt.Println("Action:", action)
		fmt.Println("Action Input:", actionInput)
	}

	var result string

	switch action {
	case kubeTool.Name:
		var param tools.KubeInput
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			result, err = kubeTool.Run(param.Commands)
			if err != nil {
				result = fmt.Sprintf("Error: Failed to run kubeTool: %v", err)
			}
		}
	case searchTool.Name:
		var param tools.SerpApiToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			output, err := searchTool.Run(param.Query)
			if err != nil {
				result = fmt.Sprintf("Error: Failed to run searchTool: %v", err)
			} else {
				result = fmt.Sprintf("Search results: %v, I need to use the tool httpRequest to get the content of the search results", output)
			}
		}
	case requestTool.Name:
		var param tools.RequestsToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			output, err := requestTool.Run(param.Url)
			if err != nil {
				result = fmt.Sprintf("Error: Failed to run requestTool: %v", err)
			} else {
				result = fmt.Sprintf("Request results: %v", output)
			}
		}
	default:
		result = fmt.Sprintf("Unknown tool: %s", action)
	}

	if debugMode {
		fmt.Println("Result:", result)
	}

	return result
}

// printDebugInfo prints debug information about the message store
func printCheckDebugInfo() {
	fmt.Println("# Message Store Debug:")
	messages := ai.MessageStore.GetMessage()
	fmt.Printf("Number of messages: %d\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("Message %d: Role=%s, Content length=%d\n", i, msg.Role, len(msg.Content))
	}
}

func buildCheckPrompt(kubeTool *tools.KubeTool, searchTool *tools.SerpApiTool, requestTool *tools.RequestsTool, query string) string {
	kubeToolDef := "Name: " + kubeTool.Name + "\nDescription: " + kubeTool.Description + "\nArgsSchema: " + kubeTool.ArgsSchema + "\n"
	searchToolDef := "Name: " + searchTool.Name + "\nDescription: " + searchTool.Description + "\nArgsSchema: " + searchTool.ArgsSchema + "\n"
	requestToolDef := "Name: " + requestTool.Name + "\nDescription: " + requestTool.Description + "\nArgsSchema: " + requestTool.ArgsSchema + "\n"

	toolsList := make([]string, 0)
	toolsList = append(toolsList, kubeToolDef, searchToolDef, requestToolDef)

	toolNames := make([]string, 0)
	toolNames = append(toolNames, kubeTool.Name, searchTool.Name, requestTool.Name)

	prompt := fmt.Sprintf(promptTpl.Template, toolsList, toolNames, query)

	return prompt
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Add namespace flag to the check command
	checkCmd.Flags().StringP("namespace", "n", "", "Default namespace to use for Kubernetes operations")

	// Add debug mode flag
	checkCmd.Flags().BoolP("debug", "d", false, "Enable debug mode to see detailed processing information")

	// Add max loops flag
	checkCmd.Flags().IntP("max-loops", "m", 5, "Maximum number of reasoning loops before stopping")
}
