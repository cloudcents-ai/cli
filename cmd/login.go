package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login using the stored API key",
	Run: func(cmd *cobra.Command, args []string) {
		// Perform the login using the stored API key
		err := checkStoredAPIKey()
		if err != nil {
			showLoginError(fmt.Sprintf("Error logging in: %v", err))
			return
		}
		showLoginSuccess("Login successful! You are now authenticated.")
	},
}

// checkStoredAPIKey reads the stored API key and logs in
func checkStoredAPIKey() error {
	apiKey, err := loadAPIKey()
	if err != nil {
		return err
	}

	// Simulate login success if API key is found
	if apiKey != "" {
		return nil
	}

	return fmt.Errorf("API key not found")
}

// loadAPIKey reads the stored API key
func loadAPIKey() (string, error) {
	filePath := filepath.Join(getLoginConfigDir(), "api_key.txt")
	apiKey, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read API key: %v", err)
	}

	return string(apiKey), nil
}

// getLoginConfigDir determines the OS-specific configuration directory
func getLoginConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "cloudcent")
	default:
		return filepath.Join(os.Getenv("HOME"), ".config", "cloudcent")
	}
}

// showLoginSuccess formats and displays a success message with a border
func showLoginSuccess(message string) {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Foreground(lipgloss.Color("118")). // Lime green for success
		Padding(1, 2)

	fmt.Println(borderStyle.Render(message))
}

// showLoginError formats and displays an error message with a border
func showLoginError(message string) {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		Foreground(lipgloss.Color("9")). // Red for errors
		Padding(1, 2)

	fmt.Println(borderStyle.Render(message))
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
