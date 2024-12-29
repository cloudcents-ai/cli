package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudcents",
	Short: "Cloud command-line tool for price fetching and demo",
	Long:  `cloudcents is a CLI tool to fetch prices and show video demos using stylish terminal output.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Any flags or configuration settings can be added here
}
