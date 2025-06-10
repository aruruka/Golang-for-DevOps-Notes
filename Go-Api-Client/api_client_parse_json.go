package main

/* type Page struct {
	Name string `json:"page"`
}

type Words struct {
	Input string   `json:"input"`
	Words []string `json:"words"`
}

type Occurrence struct {
	Words map[string]int `json:"words"`
} */

func main() {

	getHTTPJsonMap()
	/* args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: ./api-client-parse-json <url>\n")
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
		log.Fatalf("Received (HTTP Code %d) response: %s\n", response.StatusCode, string(body))
	}

	// Process the response...
	var page Page

	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
	}

	switch page.Name {
	// curl 'http://localhost:8080/words?input=word1'
	// Raw return example: {"page":"words","input":"word3","words":["word1","word2","word2","word3","word3","word3","word3"]}
	case "words":
		var words Words
		err = json.Unmarshal(body, &words)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON response for words: %s\n", err)
		}
		fmt.Printf("JSON Parsed:\nPage: %s\nWords: %s\n", page.Name, strings.Join(words.Words, ", "))
	// curl 'http://localhost:8080/occurrence'
	// Raw return example: {"page":"occurrence","words":{"word1":1,"word2":2,"word3":3}}
	case "occurrence":
		var occurrence Occurrence
		err = json.Unmarshal(body, &occurrence)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON response for occurrence: %s\n", err)
		}

		if val, ok := occurrence.Words["word1"]; ok {
			fmt.Printf("Found word1: %d\n", val)
		}

		for word, occurrence := range occurrence.Words {
			fmt.Printf("%s: %d\n", word, occurrence)
		}
	default:
		fmt.Printf("Page not found\n")
	} */
}
