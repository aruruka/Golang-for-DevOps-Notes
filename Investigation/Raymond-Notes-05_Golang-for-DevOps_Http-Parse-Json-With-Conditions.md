# Advanced JSON Parsing in API Clients

## Introduction

When developing API clients, we often encounter scenarios where the structure of the JSON response varies based on certain conditions or endpoints. This document explores techniques for conditionally parsing JSON responses in both Go and Python, with a focus on handling different data structures efficiently.

## Partial JSON Parsing

Partial JSON parsing is a technique where we initially parse only a portion of the JSON response to determine its structure or type before fully processing it. This approach is particularly useful when dealing with APIs that return different data structures for different endpoints.

### Why Use Partial Parsing?

1. **Efficiency**: Only parse what you need
2. **Flexibility**: Handle different response structures
3. **Error prevention**: Avoid unmarshaling errors when response structures vary

## JSON Parsing with Conditions in Go

### Defining Struct Types for Different Responses

In Go, we define different struct types to match different JSON response structures:

```go
type Page struct {
    Name string `json:"page"`
}

type Words struct {
    Input string   `json:"input"`
    Words []string `json:"words"`
}

type Occurrence struct {
    Words map[string]int `json:"words"`
}
```

Each struct corresponds to a specific JSON structure we expect from the API:
- `Page` represents the common field in all responses
- `Words` handles responses with arrays of strings
- `Occurrence` handles responses with maps of word frequencies

### Using Switch Statements for Conditional Processing

Go's switch statement is used to handle different response types based on the value of the "page" field:

```go
var page Page
err = json.Unmarshal(body, &page)
if err != nil {
    log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
}

switch page.Name {
case "words":
    var words Words
    err = json.Unmarshal(body, &words)
    if err != nil {
        log.Fatalf("Failed to unmarshal JSON response for words: %s\n", err)
    }
    fmt.Printf("JSON Parsed:\nPage: %s\nWords: %s\n", page.Name, strings.Join(words.Words, ", "))
case "occurrence":
    var occurrence Occurrence
    err = json.Unmarshal(body, &occurrence)
    if err != nil {
        log.Fatalf("Failed to unmarshal JSON response for occurrence: %s\n", err)
    }
    
    // Processing map data...
default:
    fmt.Printf("Page not found\n")
}
```

The instructor coded this way because:
1. It first parses only the common `page` field to determine the response type
2. Based on the response type, it unmarshals the full JSON into the appropriate struct
3. This prevents errors when trying to unmarshal JSON with incompatible structures

## Working with Maps in Go

### Iterating Over Maps

In Go, maps are unordered collections of key-value pairs. You can iterate over a map using a for-range loop:

```go
for word, occurrence := range occurrence.Words {
    fmt.Printf("%s: %d\n", word, occurrence)
}
```

It's important to note that unlike arrays, maps in Go do not guarantee order. When iterating over a map, the order of elements may vary between executions.

### Checking if Map Elements Exist

Go provides a special syntax for checking if a key exists in a map:

```go
if val, ok := occurrence.Words["word1"]; ok {
    fmt.Printf("Found word1: %d\n", val)
}
```

This returns two values:
- `val`: The value associated with the key (if it exists)
- `ok`: A boolean indicating whether the key exists in the map

## Error Handling Considerations

The instructor notes that this approach only works if the response is valid JSON. If the API returns non-JSON data, the unmarshal operation will fail. It's important to handle these errors appropriately:

```go
err = json.Unmarshal(body, &page)
if err != nil {
    log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
}
```

## Python Equivalent Implementation

Python handles conditional JSON parsing differently than Go, primarily because Python uses dictionaries rather than strictly defined structs.

### Python Classes for Different Response Types

```python
class Page:
    def __init__(self, data):
        self.name = data.get('page', '')

class Words:
    def __init__(self, data):
        self.input = data.get('input', '')
        self.words = data.get('words', [])

class Occurrence:
    def __init__(self, data):
        self.words = data.get('words', {})
```

### Conditional Processing in Python

Python doesn't have a switch statement (prior to Python 3.10), so we use if-elif-else chains:

```python
import json
import requests
import sys

def main():
    if len(sys.argv) < 2:
        print("Usage: python api_client.py <url>")
        sys.exit(1)
    
    url = sys.argv[1]
    
    try:
        response = requests.get(url)
        response.raise_for_status()  # Raise exception for 4XX/5XX responses
        
        data = response.json()
        page = Page(data)
        
        if page.name == "words":
            words_data = Words(data)
            print(f"JSON Parsed:\nPage: {page.name}\nWords: {', '.join(words_data.words)}")
        elif page.name == "occurrence":
            occurrence_data = Occurrence(data)
            
            # Check if specific element exists
            if "word1" in occurrence_data.words:
                print(f"Found word1: {occurrence_data.words['word1']}")
            
            # Iterate over dictionary
            for word, count in occurrence_data.words.items():
                print(f"{word}: {count}")
        else:
            print("Page not found")
            
    except requests.exceptions.RequestException as e:
        print(f"HTTP request failed: {e}")
    except json.JSONDecodeError as e:
        print(f"Failed to parse JSON: {e}")

if __name__ == "__main__":
    main()
```

### Python 3.10+ Switch Statement (match-case)

Python 3.10 introduced the match-case statement, which is similar to Go's switch:

```python
data = response.json()
page = Page(data)

match page.name:
    case "words":
        words_data = Words(data)
        print(f"JSON Parsed:\nPage: {page.name}\nWords: {', '.join(words_data.words)}")
    case "occurrence":
        occurrence_data = Occurrence(data)
        
        # Check if specific element exists
        if "word1" in occurrence_data.words:
            print(f"Found word1: {occurrence_data.words['word1']}")
        
        # Iterate over dictionary
        for word, count in occurrence_data.words.items():
            print(f"{word}: {count}")
    case _:
        print("Page not found")
```

### Working with Dictionaries in Python

Checking if a key exists in a dictionary:

```python
# Method 1: Using 'in' operator
if "word1" in occurrence_data.words:
    print(f"Found word1: {occurrence_data.words['word1']}")

# Method 2: Using get() with default value
value = occurrence_data.words.get("word1", None)
if value is not None:
    print(f"Found word1: {value}")
```

## Key Differences Between Go and Python Approaches

1. **Type Safety**:
   - Go: Static typing with explicit struct definitions
   - Python: Dynamic typing with flexible dictionaries

2. **Conditional Logic**:
   - Go: Switch statement
   - Python: If-elif-else or match-case (3.10+)

3. **Error Handling**:
   - Go: Explicit error returns
   - Python: Exception handling with try-except

4. **Map/Dictionary Access**:
   - Go: Special syntax for checking existence `val, ok := map[key]`
   - Python: `in` operator or `.get()` method with default value

## Best Practices

1. **Always validate API responses** before attempting to parse them
2. **Handle errors gracefully** to prevent application crashes
3. **Use partial parsing** when dealing with variable response structures
4. **Document expected response structures** to make your code more maintainable
5. **Consider using interfaces (Go) or inheritance (Python)** for more complex response handling

## Conclusion

Conditional JSON parsing is a powerful technique for building robust API clients that can handle various response structures. Both Go and Python offer effective ways to implement this pattern, with Go focusing on type safety and explicit error handling, while Python offers more flexibility and concise syntax.

By understanding these approaches, you can build more resilient API clients that gracefully handle diverse API responses.
