package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Send a dynamic prompt to the Cloud Cents API",
	Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided (the prompt)
	Run: func(cmd *cobra.Command, args []string) {
		prompt := args[0]
		p := tea.NewProgram(chatModel{prompt: prompt})
		if err := p.Start(); err != nil {
			fmt.Printf("Error starting program: %v\n", err)
		}
	},
}

// Bubble Tea model for chat
type chatModel struct {
	prompt  string
	output  string
	loading bool
}

// Initialize the chat model
func (m chatModel) Init() tea.Cmd {
	// When the program starts, send the API request asynchronously
	return m.sendAPIRequest(m.prompt)
}

// Update handles messages and updates the state
func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Allow the user to quit with "q" or "esc"
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}

	case string:
		// Once the API response is received, update the output
		m.output = msg
		m.loading = false
		return m, nil
	}

	return m, nil
}

// View renders the styled output
func (m chatModel) View() string {
	var builder strings.Builder

	// Styles for the chat bubbles
	userBubbleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")). // White text
		Background(lipgloss.Color("33")). // Blue background
		Padding(1, 2).
		Width(50).
		Margin(1, 0, 1, 15).
		Align(lipgloss.Right)

	botBubbleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).   // Black text
		Background(lipgloss.Color("248")). // Light grey background
		Padding(1, 2).
		Width(50).
		Margin(1, 15, 1, 0).
		Align(lipgloss.Left)

	// User prompt bubble
	userBubble := userBubbleStyle.Render(fmt.Sprintf("You: %s", m.prompt))

	// If we're still waiting for the response, show a loading message
	if m.loading {
		botBubble := botBubbleStyle.Render("Loading response...")
		builder.WriteString(botBubble + "\n" + userBubble)
	} else {
		// Bot response bubble
		botBubble := botBubbleStyle.Render(fmt.Sprintf("Response: %s", m.output))
		builder.WriteString(botBubble + "\n" + userBubble)
	}

	return builder.String() + "\n\nPress 'q' or 'esc' to exit."
}

// Function to query the API with a dynamic prompt
func (m chatModel) sendAPIRequest(prompt string) tea.Cmd {
	m.loading = true
	return func() tea.Msg {
		// Replace spaces in the prompt with '%20' to make it URL-safe
		escapedPrompt := strings.ReplaceAll(prompt, " ", "%20")

		// Define the URL with the dynamic prompt
		url := fmt.Sprintf("https://cloud-cents.onrender.com/generate/test/%s?max_length=200", escapedPrompt)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "Error creating request"
		}

		// Add the headers
		req.Header.Add("accept", "application/json")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "Error sending request"
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "Error reading response"
		}

		// Return the response as a string message
		return string(body)
	}
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
