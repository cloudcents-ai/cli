package cmd

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Task structure with title and description
type Task struct {
	title       string
	description string
}

// Define tasks with titles and descriptions
var tasks = []Task{
	{"Configuration", "Enter your cloud service configurations like compute and storage."},
	{"Dashboard", "Visualize and analyze the usage and cost data for your services."},
	{"Forecasting", "Predict future usage and optimize cloud expenses."},
	{"Report", "Generate detailed reports on your cloud usage and expenses."},
}

// Task completion status for each task (false by default)
var checked = map[int]bool{
	0: false,
	1: false,
	2: false,
	3: false,
}

// Define styles for checklist items
var (
	titleStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("118")).Bold(true)   // Lime Green
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))                // Aqua
	statusStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Bold(true)    // Pink for status
)

// checklistCmd represents the checklist command
var checklistCmd = &cobra.Command{
	Use:   "checklist",
	Short: "View and complete a checklist of tasks",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(model{})
		if err := p.Start(); err != nil {
			fmt.Printf("Error starting program: %v\n", err)
			os.Exit(1)
		}
	},
}

// model represents the Bubble Tea model for the checklist
type model struct{}

// Init initializes the Bubble Tea program (no initial command here)
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles user input for the checklist
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses
		switch strings.ToLower(msg.String()) {
		case "q", "esc":
			return m, tea.Quit // Exit the program on 'q' or 'esc'
		case "1", "2", "3", "4":
			toggleCheck(msg.String()) // Toggle the task when 1-4 is pressed
			if allTasksComplete() {
				fmt.Println("You've completed the checklist!")
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View renders the checklist UI
func (m model) View() string {
	var sb strings.Builder
	sb.WriteString("Your Task Checklist:\n\n")

	// Iterate over tasks and build the formatted checklist
	for i, task := range tasks {
		status := "[ ]"
		if checked[i] {
			status = "[x]"
		}

		// Render task title and description with styling
		sb.WriteString(fmt.Sprintf("%s %s\n", statusStyle.Render(status), titleStyle.Render(task.title)))
		sb.WriteString(fmt.Sprintf("    %s\n", descriptionStyle.Render(task.description)))
	}

	sb.WriteString("\nPress 1, 2, 3, or 4 to toggle tasks. Press 'q' or 'esc' to quit.")
	return sb.String()
}

// toggleCheck toggles the completion status of a task based on the user input (1-4)
func toggleCheck(taskNumber string) {
	index := taskNumberToIndex(taskNumber)
	checked[index] = !checked[index]
}

// taskNumberToIndex converts user input (1-4) into array indices (0-3)
func taskNumberToIndex(taskNumber string) int {
	switch taskNumber {
	case "1":
		return 0
	case "2":
		return 1
	case "3":
		return 2
	case "4":
		return 3
	default:
		return 0
	}
}

// allTasksComplete checks if all tasks are checked off (true if all are checked)
func allTasksComplete() bool {
	for _, done := range checked {
		if !done {
			return false
		}
	}
	return true
}

// Initialize the command
func init() {
	rootCmd.AddCommand(checklistCmd)
}
