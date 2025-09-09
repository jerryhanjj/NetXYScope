package main

import (
	"fmt"
	"os"

	"github.com/jerryhanjj/NetXYScope/internal/models"
	"github.com/jerryhanjj/NetXYScope/internal/search"
)

func main() {
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	searchTerm := os.Args[1]
	directory := os.Args[2]

	engine := search.NewEngine()
	results, err := engine.SearchFiles(directory, searchTerm)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(results) == 0 {
		fmt.Printf("No matches found for '%s'\n", searchTerm)
		return
	}

	printResults(results, searchTerm)
}

func printUsage() {
	fmt.Printf("Usage: %s <search-term> <directory-path>\n", os.Args[0])
	fmt.Println("Example: NetXYScope interface /path/to/yang/files")
	fmt.Println("\nSupported file types:")
	fmt.Println("  - .xml (NETCONF XML configuration files)")
	fmt.Println("  - .yang (YANG model files)")
	fmt.Println("  - .yin (YANG in XML format)")
}

func printResults(results []models.SearchResult, searchTerm string) {
	fmt.Printf("Found %d matches for '%s':\n\n", len(results), searchTerm)

	// Group results by file for better readability
	currentFile := ""
	for _, result := range results {
		if result.FilePath != currentFile {
			if currentFile != "" {
				fmt.Println()
			}
			fmt.Printf("=== %s ===\n", result.FilePath)
			currentFile = result.FilePath
		}

		// Highlight the search term in the output
		highlighted := search.HighlightSearchTerm(result.LineContent, searchTerm)
		fmt.Printf("%4d | %s\n", result.LineNumber, highlighted)
	}
}
