package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// {"page":"words","input":"word1","words":["word1"]}
type Words struct {
	Page  string   `json:"page"`
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func fetchWordsFromAPI() {
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./api-client <url>\n")
		os.Exit(1)
	}
	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(args[1])
	if err != nil {
		log.Fatalf("Failed to make HTTP GET request: %s\n", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Failed to read response body: %s\n", err)
	}

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Received (HTTP Code %d) response: %s\n", response.StatusCode, body)
	}

	// Process the response...
	var words Words

	err = json.Unmarshal(body, &words)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
	}

	fmt.Printf("JSON Parsed:\nPage: %s\nWords: %s\n", words.Page, strings.Join(words.Words, ", "))
}
