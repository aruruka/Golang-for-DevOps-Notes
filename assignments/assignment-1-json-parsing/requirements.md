# Assignment 1: JSON Parsing

## Overview
In this assignment you'll need to parse a JSON response in a Go struct. This exercise will help you practice HTTP requests, JSON unmarshaling, and working with Go structs.

## Requirements

### Task
Parse the JSON from the endpoint `http://localhost:8080/assignment1` in a Go struct using `http.Get` to retrieve the data.

### Expected JSON Structure
The API endpoint will return JSON data with the following structure:

```json
{
  "page": "assignment1",
  "words": [
    "one",
    "two", 
    "three"
  ],
  "percentages": {
    "one": 0.33,
    "three": 0,
    "two": 0.66
  },
  "special": [
    "one",
    "two",
    null
  ],
  "extraSpecial": [
    1,
    2,
    "3"
  ]
}
```

### Success Criteria
- Successfully make an HTTP GET request to the endpoint
- Parse the JSON response into appropriate Go structs
- Output all JSON data in a Go struct format
- Handle the mixed data types in the `extraSpecial` array appropriately
- Handle null values in the `special` array correctly

### Implementation Notes
- Use `http.Get` to retrieve the data from the test server
- Define appropriate Go structs to match the JSON structure
- Use JSON struct tags for proper field mapping
- Handle the different data types present in the response
- Test your implementation multiple times as the values will change

### Setup
1. Start the test server (located in the `test-server` directory)
2. Ensure the server is running on `http://localhost:8080`
3. Your Go program should make requests to `http://localhost:8080/assignment1`

### Validation
You'll know your implementation works when:
- You can successfully retrieve data from the endpoint
- All JSON fields are properly parsed into Go struct fields
- The program outputs the complete struct data
- The program handles changing values correctly on multiple runs

## Additional Resources
Feel free to ask any questions in the Q&A section if you need help with implementation details.
