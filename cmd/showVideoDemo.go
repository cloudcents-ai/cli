package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// showVideoDemoCmd represents the showVideoDemo command
var showVideoDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Display a video demo",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(videoModel{})
		if err := p.Start(); err != nil {
			fmt.Printf("Error starting program: %v\n", err)
		}
	},
}

type videoModel struct{}

// Initialize model
func (m videoModel) Init() tea.Cmd {
	return nil
}

// Update model
func (m videoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// If the user presses 'q' or 'esc', quit the program
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View model - Styled with YouTube and Twitch icons, colors, and reduced padding
func (m videoModel) View() string {
	// Define styles for the title, link, and Twitch section with reduced padding
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")). // White text
		Background(lipgloss.Color("63")).  // Purple background
		Bold(true).
		Align(lipgloss.Center)

	twitchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("231")). // White text
		Background(lipgloss.Color("170")). // Magenta background
		Bold(true).
		Align(lipgloss.Center)

	// Padding at the top
	paddingStyle := lipgloss.NewStyle().Padding(0, 0) // No extra padding

	// YouTube and Twitch Icons
	youtubeIcon := "📺" // YouTube icon
	twitchIcon := "🎮"  // Twitch icon

	// Combine the title, link, and Twitch section into a formatted "table-like" output
	title := titleStyle.Render(fmt.Sprintf("%s CloudCents Demo: https://www.youtube.com/watch?v=Sx7q05HhqBA", youtubeIcon))
	twitch := twitchStyle.Render(fmt.Sprintf("%s Join me coding this live on Twitch: https://www.twitch.tv/majesticcodingtwitch", twitchIcon))

	// Create the full layout with minimal padding at the top
	view := paddingStyle.Render("") + lipgloss.JoinVertical(lipgloss.Top, title, twitch) +
		"\n\n" +
		"Use 'cloudcents demo --open' to open this video in your browser.\n" +
		"Press 'esc' or 'q' to exit."

	return view
}

func init() {
	rootCmd.AddCommand(showVideoDemoCmd)
}

// Open the video link in the browser when the "open" argument is provided
func init() {
	showVideoDemoCmd.Flags().BoolP("open", "o", false, "Open video in browser")
	showVideoDemoCmd.Run = func(cmd *cobra.Command, args []string) {
		openVideo, _ := cmd.Flags().GetBool("open")
		if openVideo {
			openBrowser("https://www.youtube.com/watch?v=Sx7q05HhqBA")
		} else {
			p := tea.NewProgram(videoModel{})
			if err := p.Start(); err != nil {
				fmt.Printf("Error starting program: %v\n", err)
			}
		}
	}
}

// Open a URL in the default browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("Unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

// escCmd represents the command to exit the CLI
var escCmd = &cobra.Command{
	Use:   "esc",
	Short: "Exit the program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Exiting the program. Have a great day!")
	},
}

func init() {
	rootCmd.AddCommand(escCmd)
}
