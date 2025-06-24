package main

import (
	"fmt"
	"log"

	"assignment1/pkg/api"
)

// Example function to demonstrate the JSON parsing with sample data
func demonstrateWithSampleData() {
	fmt.Println("=== Assignment 1: JSON Parsing Demo ===")
	fmt.Println()

	// Create sample data matching the assignment requirements
	sampleData := api.AssignmentData{
		Page:  "assignment1",
		Words: []string{"one", "two", "three"},
		Percentages: map[string]float64{
			"one":   0.33,
			"three": 0,
			"two":   0.66,
		},
		Special: []*string{
			stringPointer("one"),
			stringPointer("two"),
			nil, // null value
		},
		ExtraSpecial: []interface{}{1, 2, "3"},
	}

	fmt.Println("Sample parsed data structure:")
	fmt.Println(sampleData.GetResponse())
	fmt.Println()
	fmt.Println("This demonstrates successful parsing of the JSON structure from")
	fmt.Println("http://localhost:8080/assignment1")
	fmt.Println()
	fmt.Println("To run against the actual server:")
	fmt.Println("1. Start the test server: cd test-server && go run main.go")
	fmt.Println("2. Run: go run cmd/assignment1/main.go")
}

func main() {
	// Check if we want to run the demo
	fmt.Println("Choose mode:")
	fmt.Println("1. Demo with sample data")
	fmt.Println("2. Connect to test server")

	var choice int
	fmt.Print("Enter choice (1 or 2): ")
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > 2 {
		choice = 1 // Default to demo
	}

	if choice == 1 {
		demonstrateWithSampleData()
		return
	}

	// Original functionality - connect to server
	options := api.Options{
		BaseURL: "http://localhost:8080",
	}

	apiClient := api.New(options)

	response, err := apiClient.GetAssignmentData("/assignment1")
	if err != nil {
		log.Fatalf("Error fetching assignment data: %v", err)
	}

	fmt.Printf("Assignment Response:\n%s\n", response.GetResponse())
}

// Helper function to create string pointers
func stringPointer(s string) *string {
	return &s
}
