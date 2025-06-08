package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func httpGet() {
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
		log.Fatalf("Received non-200 response: %d\n", response.StatusCode)
	}

	// Process the response...
	fmt.Printf("HTTP Status Code: %d\nBody: %v\n", response.StatusCode, string(body))
}
