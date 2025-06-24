package api

import (
	"encoding/json"
	"fmt"
	"io"
)

// AssignmentData represents the structure of the assignment1 JSON response
type AssignmentData struct {
	Page         string             `json:"page"`
	Words        []string           `json:"words"`
	Percentages  map[string]float64 `json:"percentages"`
	Special      []*string          `json:"special"`      // Pointer to handle null values
	ExtraSpecial []interface{}      `json:"extraSpecial"` // interface{} to handle mixed types
}

// GetResponse implements the Response interface for AssignmentData
func (a AssignmentData) GetResponse() string {
	result := fmt.Sprintf("Page: %s\n", a.Page)
	result += fmt.Sprintf("Words: %v\n", a.Words)
	result += fmt.Sprintf("Percentages: %v\n", a.Percentages)

	// Handle special array with null values
	specialStr := "["
	for i, item := range a.Special {
		if i > 0 {
			specialStr += ", "
		}
		if item == nil {
			specialStr += "null"
		} else {
			specialStr += fmt.Sprintf("\"%s\"", *item)
		}
	}
	specialStr += "]"
	result += fmt.Sprintf("Special: %s\n", specialStr)

	result += fmt.Sprintf("ExtraSpecial: %v", a.ExtraSpecial)
	return result
}

// GetAssignmentData implements the APIIface interface
func (a api) GetAssignmentData(endpoint string) (Response, error) {
	requestURL := a.Options.BaseURL + endpoint

	response, err := a.Client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP Get error: %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll error: %s", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid output (HTTP Code %d): %s", response.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      "Response is not valid JSON",
		}
	}

	var assignmentData AssignmentData
	err = json.Unmarshal(body, &assignmentData)
	if err != nil {
		return nil, RequestError{
			HTTPCode: response.StatusCode,
			Body:     string(body),
			Err:      fmt.Sprintf("JSON unmarshal error: %s", err),
		}
	}

	return assignmentData, nil
}
