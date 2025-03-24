package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"kgent/cmd/ai"
	promptTpl "kgent/cmd/prompt"
	"kgent/cmd/tools"
	"kgent/cmd/utils"

	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Interact with the Kubernetes assistant through chat",
	Long: `Chat with the Kubernetes assistant to create, list, or delete resources.
The assistant will help you perform actions on your Kubernetes cluster using
natural language. You can ask it to create resources like pods or services,
list existing resources, or delete resources.

Simply type your query and the assistant will either answer directly or
ask for additional information if needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize tools
		createTool := tools.NewCreateTool()
		listTool := tools.NewListTool()
		deleteTool := tools.NewDeleteTool()
		humanTool := tools.NewHumanTool()

		// Get default namespace from flag
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace != "" {
			fmt.Printf("Using namespace: %s\n", namespace)
		}

		// Get debug mode flag
		debugMode, _ := cmd.Flags().GetBool("debug")

		// Get max loops flag
		maxLoops, _ := cmd.Flags().GetInt("max-loops")

		runChatLoop(cmd, createTool, listTool, deleteTool, humanTool, namespace, debugMode, maxLoops)
	},
}

// runChatLoop handles the main chat interaction loop
func runChatLoop(cmd *cobra.Command, createTool *tools.CreateTool, listTool *tools.ListTool,
	deleteTool *tools.DeleteTool, humanTool *tools.HumanTool,
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

		prompt := buildPrompt(createTool, listTool, deleteTool, humanTool, input)
		if debugMode {
			fmt.Println("User prompt:", prompt)
		}
		ai.MessageStore.AddUser(prompt)

		processConversation(createTool, listTool, deleteTool, humanTool, maxLoops, debugMode)
		ai.MessageStore.Clear()
	}
}

// processConversation handles the AI interaction and tool execution
func processConversation(createTool *tools.CreateTool, listTool *tools.ListTool,
	deleteTool *tools.DeleteTool, humanTool *tools.HumanTool,
	maxLoops int, debugMode bool) {
	loopCount := 1

	for loopCount <= maxLoops {
		if debugMode {
			fmt.Printf("---------------- Response round %d ----------------\n", loopCount)
			printDebugInfo()
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
			result := handleAction(createTool, listTool, deleteTool, humanTool,
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
func handleAction(createTool *tools.CreateTool, listTool *tools.ListTool,
	deleteTool *tools.DeleteTool, humanTool *tools.HumanTool,
	action string, actionInput string, debugMode bool) string {
	if debugMode {
		fmt.Println("# Action Debug:")
		fmt.Println("Action:", action)
		fmt.Println("Action Input:", actionInput)
	}

	var result string

	switch action {
	case createTool.Name:
		var param tools.CreateToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			result = createTool.Run(param.Prompt, param.Resource, debugMode)
		}
	case listTool.Name:
		var param tools.ListToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			var err error
			result, err = listTool.Run(param.Resource, param.Namespace)
			if err != nil {
				result = fmt.Sprintf("Error listing resources: %v", err)
			}
		}
	case deleteTool.Name:
		var param tools.DeleteToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			if err := deleteTool.Run(param.Resource, param.Name, param.Namespace); err != nil {
				result = fmt.Sprintf("Delete failed: %v", err)
			} else {
				result = "Resource deleted successfully"
			}
		}
	case humanTool.Name:
		var param tools.HumanToolParam
		if err := json.Unmarshal([]byte(actionInput), &param); err != nil {
			result = fmt.Sprintf("Error: Failed to parse action input: %v", err)
		} else {
			result = humanTool.Run(param.Prompt)
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
func printDebugInfo() {
	fmt.Println("# Message Store Debug:")
	messages := ai.MessageStore.GetMessage()
	fmt.Printf("Number of messages: %d\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("Message %d: Role=%s, Content length=%d\n", i, msg.Role, len(msg.Content))
	}
}

func buildPrompt(createTool *tools.CreateTool, listTool *tools.ListTool, deleteTool *tools.DeleteTool, humanTool *tools.HumanTool, query string) string {
	createToolDef := "Name: " + createTool.Name + "\nDescription: " + createTool.Description + "\nArgsSchema: " + createTool.ArgsSchema + "\n"
	listToolDef := "Name: " + listTool.Name + "\nDescription: " + listTool.Description + "\nArgsSchema: " + listTool.ArgsSchema + "\n"
	deleteToolDef := "Name: " + deleteTool.Name + "\nDescription: " + deleteTool.Description + "\nArgsSchema: " + deleteTool.ArgsSchema + "\n"
	humanToolDef := "Name: " + humanTool.Name + "\nDescription: " + humanTool.Description + "\nArgsSchema: " + humanTool.ArgsSchema + "\n"

	toolsList := make([]string, 0)
	toolsList = append(toolsList, createToolDef, listToolDef, deleteToolDef, humanToolDef)

	toolNames := make([]string, 0)
	toolNames = append(toolNames, createTool.Name, listTool.Name, deleteTool.Name, humanTool.Name)

	prompt := fmt.Sprintf(promptTpl.Template, toolsList, toolNames, query)

	return prompt
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Add namespace flag to the chat command
	chatCmd.Flags().StringP("namespace", "n", "", "Default namespace to use for Kubernetes operations")

	// Add debug mode flag
	chatCmd.Flags().BoolP("debug", "d", false, "Enable debug mode to see detailed processing information")

	// Add max loops flag
	chatCmd.Flags().IntP("max-loops", "m", 5, "Maximum number of reasoning loops before stopping")
}
