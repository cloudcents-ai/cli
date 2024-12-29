package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// Data structure to hold the pricing information
type PricingData struct {
	AWS struct {
		Compute map[string]float64 `json:"compute"`
		Storage map[string]float64 `json:"storage"`
	} `json:"aws"`
	GCP struct {
		Compute map[string]float64 `json:"compute"`
		Storage map[string]float64 `json:"storage"`
	} `json:"gcp"`
	Azure struct {
		Compute map[string]float64 `json:"compute"`
		Storage map[string]float64 `json:"storage"`
	} `json:"azure"`
}

// Pricing data
var prices PricingData

// Define styles (same as before)
var (
	bestPriceColor   = lipgloss.Color("#0000FF")
	highPriceColor   = lipgloss.Color("#FF4500")
	mediumPriceColor = lipgloss.Color("#FFD700")
	lowPriceColor    = lipgloss.Color("#00FFFF")
	cellStyle        = lipgloss.NewStyle().Padding(0, 2)
	headerStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	lineStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// getPricesCmd represents the getPrices command
var getPricesCmd = &cobra.Command{
	Use:   "prices",
	Short: "Get pricing for AWS, GCP, and Azure from a JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		loadPricingData()
		printLegend()
		printPricingTable()
	},
}

// Load the pricing data from data.json
func loadPricingData() {
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Error opening data.json:", err)
		return
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue, &prices)
}

// printPricingTable prints the pricing data in a styled table with a heatmap
func printPricingTable() {
	// Table header
	header := fmt.Sprintf("%-10s %-15s %-10s %-10s %-10s", "Service", "Size", "AWS ($)", "GCP ($)", "Azure ($)")
	fmt.Println(headerStyle.Render(header))
	fmt.Println(lineStyle.Render(strings.Repeat("-", 60)))

	// Services and sizes
	services := []string{"compute", "storage"}
	sizes := []string{"small", "medium", "large"}

	// Iterate through services and sizes, and print prices for each provider with heatmap color-coding
	for _, service := range services {
		for _, size := range sizes {
			awsPrice := getPrice("aws", service, size)
			gcpPrice := getPrice("gcp", service, size)
			azurePrice := getPrice("azure", service, size)

			// Apply heatmap color based on price for each provider
			awsStyledPrice := stylePriceCell(awsPrice, service, size)
			gcpStyledPrice := stylePriceCell(gcpPrice, service, size)
			azureStyledPrice := stylePriceCell(azurePrice, service, size)

			// Format the row with the styled prices
			row := fmt.Sprintf("%-10s %-15s %-10s %-10s %-10s", service, size, awsStyledPrice, gcpStyledPrice, azureStyledPrice)
			fmt.Println(row)
		}
		fmt.Println(lineStyle.Render(strings.Repeat("-", 60))) // separator line after each service block
	}
}

// getPrice returns the price for a given provider, service, and size
func getPrice(provider, service, size string) float64 {
	switch provider {
	case "aws":
		if service == "compute" {
			return prices.AWS.Compute[size]
		}
		return prices.AWS.Storage[size]
	case "gcp":
		if service == "compute" {
			return prices.GCP.Compute[size]
		}
		return prices.GCP.Storage[size]
	case "azure":
		if service == "compute" {
			return prices.Azure.Compute[size]
		}
		return prices.Azure.Storage[size]
	}
	return 0
}

// stylePriceCell styles a price cell based on its value and applies heatmap colors
func stylePriceCell(price float64, service, size string) string {
	bestPrice := findBestPrice(service, size)

	var color lipgloss.Color
	switch {
	case price == bestPrice:
		color = bestPriceColor
	case price <= bestPrice*1.2:
		color = lowPriceColor
	case price <= bestPrice*1.5:
		color = mediumPriceColor
	default:
		color = highPriceColor
	}

	return cellStyle.Copy().Background(color).Render(fmt.Sprintf("%.3f", price))
}

// findBestPrice finds the best (lowest) price across AWS, GCP, and Azure for a given service and size
func findBestPrice(service, size string) float64 {
	awsPrice := getPrice("aws", service, size)
	gcpPrice := getPrice("gcp", service, size)
	azurePrice := getPrice("azure", service, size)
	return min(awsPrice, gcpPrice, azurePrice)
}

// min returns the minimum value from a list of floats
func min(vals ...float64) float64 {
	minVal := vals[0]
	for _, val := range vals[1:] {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}

// printLegend prints a color-coded legend for the heatmap
func printLegend() {
	fmt.Println("\nLegend:")
	legend := cellStyle.Copy().Background(bestPriceColor).Render(" Best Price ") + " " +
		cellStyle.Copy().Background(lowPriceColor).Render(" Lower Price ") + " " +
		cellStyle.Copy().Background(mediumPriceColor).Render(" Moderate Price ") + " " +
		cellStyle.Copy().Background(highPriceColor).Render(" Higher Price ")

	fmt.Println(legend)
	fmt.Println("\n")
}

func init() {
	rootCmd.AddCommand(getPricesCmd)
}
