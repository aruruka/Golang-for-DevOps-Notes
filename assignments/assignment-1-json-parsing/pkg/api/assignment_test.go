package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// MockClient implements ClientIface for testing
type MockClient struct {
	GetResponse *http.Response
}

// Get implements the ClientIface interface for testing
func (m MockClient) Get(url string) (resp *http.Response, err error) {
	return m.GetResponse, nil
}

func TestGetAssignmentData(t *testing.T) {
	// Create test data that matches the expected API response
	testData := AssignmentData{
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
			nil, // This represents the null value
		},
		ExtraSpecial: []interface{}{1, 2, "3"},
	}

	// Marshal test data to JSON
	testDataBytes, err := json.Marshal(testData)
	if err != nil {
		t.Errorf("marshal error: %s", err)
	}

	// Create API instance with mock client
	apiInstance := api{
		Options: Options{BaseURL: "http://localhost:8080"},
		Client: MockClient{
			GetResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(testDataBytes)),
			},
		},
	}

	// Test the GetAssignmentData method
	response, err := apiInstance.GetAssignmentData("/assignment1")
	if err != nil {
		t.Errorf("GetAssignmentData error: %s", err)
	}

	if response == nil {
		t.Errorf("Response is nil")
		return
	}

	// Verify the response type
	assignmentData, ok := response.(AssignmentData)
	if !ok {
		t.Errorf("Response is not of type AssignmentData")
		return
	}

	// Test individual fields
	if assignmentData.Page != "assignment1" {
		t.Errorf("Expected page 'assignment1', got '%s'", assignmentData.Page)
	}

	if len(assignmentData.Words) != 3 {
		t.Errorf("Expected 3 words, got %d", len(assignmentData.Words))
	}

	if assignmentData.Words[0] != "one" {
		t.Errorf("Expected first word 'one', got '%s'", assignmentData.Words[0])
	}

	if assignmentData.Percentages["one"] != 0.33 {
		t.Errorf("Expected percentage for 'one' to be 0.33, got %f", assignmentData.Percentages["one"])
	}

	// Test null handling in Special array
	if assignmentData.Special[2] != nil {
		t.Errorf("Expected third item in Special to be nil")
	}

	if assignmentData.Special[0] == nil || *assignmentData.Special[0] != "one" {
		t.Errorf("Expected first item in Special to be 'one'")
	}

	// Test mixed types in ExtraSpecial
	if len(assignmentData.ExtraSpecial) != 3 {
		t.Errorf("Expected 3 items in ExtraSpecial, got %d", len(assignmentData.ExtraSpecial))
	}

	// Verify the response string formatting
	responseStr := response.GetResponse()
	if responseStr == "" {
		t.Errorf("GetResponse returned empty string")
	}

	// Check if response contains expected elements
	expectedContents := []string{"Page: assignment1", "Words:", "Percentages:", "Special:", "ExtraSpecial:"}
	for _, expected := range expectedContents {
		if !containsString(responseStr, expected) {
			t.Errorf("Response does not contain '%s'", expected)
		}
	}
}

func TestGetAssignmentDataErrorHandling(t *testing.T) {
	// Test HTTP error response
	apiInstance := api{
		Options: Options{BaseURL: "http://localhost:8080"},
		Client: MockClient{
			GetResponse: &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte("Not Found"))),
			},
		},
	}

	_, err := apiInstance.GetAssignmentData("/assignment1")
	if err == nil {
		t.Errorf("Expected error for 404 response, got nil")
	}
}

func TestGetAssignmentDataInvalidJSON(t *testing.T) {
	// Test invalid JSON response
	apiInstance := api{
		Options: Options{BaseURL: "http://localhost:8080"},
		Client: MockClient{
			GetResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
			},
		},
	}

	_, err := apiInstance.GetAssignmentData("/assignment1")
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}

	// Verify it's a RequestError
	if _, ok := err.(RequestError); !ok {
		t.Errorf("Expected RequestError for invalid JSON, got %T", err)
	}
}

// Helper function to create string pointers
func stringPointer(s string) *string {
	return &s
}

// Helper function to check if a string contains a substring
func containsString(str, substr string) bool {
	return len(str) >= len(substr) && findSubstring(str, substr) != -1
}

// Helper function to find substring
func findSubstring(str, substr string) int {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
