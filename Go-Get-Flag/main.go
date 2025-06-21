package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Response interface {
	GetResponse() string
}

type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

func (w Words) GetResponse() string {
	return fmt.Sprintf("Words: %s", strings.Join(w.Words, ", "))
}

type Occurrence struct {
	Words map[string]int `json:"words"`
}

func (o Occurrence) GetResponse() string {
	words := []string{}
	for word, occurrence := range o.Words {
		words = append(words, fmt.Sprintf("%s: %d", word, occurrence))
	}
	return fmt.Sprintf("Words: %s", strings.Join(words, ", "))
}

func main() {
	var (
		requestURL string
		password   string
		parsedURL  *url.URL
		err        error
	)

	flag.StringVar(&requestURL, "url", "", "url to access")
	flag.StringVar(&password, "password", "", "use a password to access our api")

	flag.Parse()

	if requestURL == "" {
		fmt.Println("Please provide a URL using -url flag")
		os.Exit(1)
	}

	if parsedURL, err = url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Validation error: URL is not valid: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	res, err := doRequest(parsedURL.String())
	if err != nil {
		if reqErr, ok := err.(RequestError); ok {
			fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", reqErr.Err, reqErr.HTTPCode, reqErr.Body)
			os.Exit(1)
		}
	}

	if res == nil {
		fmt.Printf("No response received.\n")
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", res.GetResponse())
}

func doRequest(requestURL string) (Response, error) {

	if _, err := url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("URL is in invalid format: %s\n", err)
		os.Exit(1)
	}

	response, err := http.Get(requestURL)

	if err != nil {
		return nil, fmt.Errorf("get error: %s", err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (http code: %d): %s", response.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "no valid json returned",
		}
	}

	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %s", err)
	}

	switch page.Name {
	// curl 'http://localhost:8080/words?input=word1'
	// Raw return example: {"page":"words","input":"word3","words":["word1","word2","word2","word3","word3","word3","word3"]}
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error for words: %s", err)
		}

		return words, nil

	// curl 'http://localhost:8080/occurrence'
	// Raw return example: {"page":"occurrence","words":{"word1":1,"word2":2,"word3":3}}
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error for occurrence: %s", err)
		}

		return occurrence, nil
	}

	return nil, nil
}
