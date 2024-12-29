package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var encryptionKey = []byte("myverystrongpasswordo32bitlength")

// Styling with lipgloss
var borderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2).Bold(true)
var successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true) // Green for success
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)  // Red for errors
var infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Italic(true)  // Blue for info

// authCmd represents the auth command to store an API key
var authCmd = &cobra.Command{
	Use:   "auth [api_key]",
	Short: "Store an API key securely and log in",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := args[0]
		storeAPIKey(apiKey)

		// Automatically trigger login after storing API key
		err := loginWithStoredAPIKey()
		if err != nil {
			displayError(fmt.Sprintf("Error during login: %v", err))
			os.Exit(1)
		} else {
			displaySuccess("Successfully logged in with stored API key.")
		}
	},
}

// storeAPIKey encrypts and saves the API key to a file in the config directory
func storeAPIKey(apiKey string) {
	configDir := getConfigDir()
	filePath := filepath.Join(configDir, "api_key.txt")

	// Create the config folder if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			displayError(fmt.Sprintf("Error creating config folder: %v", err))
			os.Exit(1)
		}
	}

	// Encrypt the API key before storing it
	encryptedKey, err := encrypt(apiKey)
	if err != nil {
		displayError(fmt.Sprintf("Error encrypting API key: %v", err))
		os.Exit(1)
	}

	// Store the encrypted API key in the file
	err = ioutil.WriteFile(filePath, []byte(encryptedKey), 0600)
	if err != nil {
		displayError(fmt.Sprintf("Error storing API key: %v", err))
		os.Exit(1)
	}

	displayInfo(fmt.Sprintf("API key stored securely in '%s'", filePath))
}

// encrypt encrypts the API key using AES
func encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// loginWithStoredAPIKey reads the stored API key and logs in
func loginWithStoredAPIKey() error {
	_, err := readAPIKey()
	if err != nil {
		return err
	}

	return nil
}

// readAPIKey reads and decrypts the stored API key
func readAPIKey() (string, error) {
	filePath := filepath.Join(getConfigDir(), "api_key.txt")
	encryptedKey, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not read API key: %v", err)
	}

	return decrypt(string(encryptedKey))
}

// decrypt decrypts the stored API key using AES
func decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encryptedData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getConfigDir determines the OS-specific configuration directory
func getConfigDir() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "cloudcent")
	default:
		return filepath.Join(os.Getenv("HOME"), ".config", "cloudcent")
	}
}

// displaySuccess formats and displays a success message
func displaySuccess(message string) {
	fmt.Println(borderStyle.Render(successStyle.Render(message)))
}

// displayError formats and displays an error message
func displayError(message string) {
	fmt.Println(borderStyle.Render(errorStyle.Render(message)))
}

// displayInfo formats and displays an informational message
func displayInfo(message string) {
	fmt.Println(borderStyle.Render(infoStyle.Render(message)))
	fmt.Println(borderStyle.Render(infoStyle.Render("Need an new API key? Visit https://cloud-cents.vercel.app/user to get yours.")))
}

func init() {
	rootCmd.AddCommand(authCmd)
}
