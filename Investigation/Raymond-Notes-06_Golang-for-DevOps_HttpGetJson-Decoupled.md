# HTTP JSON Request Handling in Go: Decoupled Approach

## Introduction

This document explores the decoupling of HTTP request and JSON parsing logic in Go applications. By separating responsibilities into distinct functions, we can improve code maintainability, testability, and extensibility.

## Decoupling Code in Go

Decoupling refers to the practice of separating code into logical units with specific responsibilities. In the context of our HTTP client application, this involves:

1. Separating the main function from the request logic
2. Creating a clean interface between components
3. Implementing proper error handling

### Benefits of Decoupling

- **Maintainability**: Changes to one part of the system don't require changes to others
- **Testability**: Smaller functions are easier to test in isolation
- **Reusability**: Decoupled components can be reused in different contexts
- **Readability**: Code is easier to understand when it has a single responsibility

## Using Interfaces in Go

Interfaces in Go provide a powerful way to define behavior without specifying implementation details. This allows for flexible code that can work with different types that implement the same behavior.

### The Response Interface

In our example, we use an interface to handle different types of API responses:

```go
type Response interface {
    GetResponse() string
}
```

This interface allows us to handle different struct types (Words and Occurrence) with a single return type in our function signature.

### Why Use Interfaces?

The instructor explains the rationale for using interfaces:

> "Ideally you don't really want to change your function signature every time that you add a feature or add another API call here. So we are going to choose something more generic."

By using interfaces, we can add new response types without changing the function signature of our request function. This maintains backward compatibility and adheres to the Open/Closed Principle (open for extension, closed for modification).

## Implementation Details

### Struct Types and Interface Implementation

Our application handles two types of responses:

1. **Words**: A list of words from the API
2. **Occurrence**: A map of words and their occurrence counts

Both structs implement the `Response` interface:

```go
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
```

### The Request Function

The decoupled request function handles:
- URL validation
- Making the HTTP request
- Reading the response body
- Parsing the JSON
- Returning the appropriate response type

```go
func doRequest(requestURL string) (Response, error) {
    if _, err := url.ParseRequestURI(requestURL); err != nil {
        return nil, fmt.Errorf("validation error: URL is not valid: %s", err)
    }

    response, err := http.Get(requestURL)
    if err != nil {
        return nil, fmt.Errorf("HTTP get error: %s", err)
    }
    defer response.Body.Close()

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("read all error: %s", err)
    }

    if response.StatusCode != 200 {
        return nil, fmt.Errorf("invalid output (http code: %d): %s", 
                              response.StatusCode, string(body))
    }

    var page Page
    err = json.Unmarshal(body, &page)
    if err != nil {
        return nil, fmt.Errorf("unmarshal error: %s", err)
    }

    switch page.Name {
    case "words":
        var words Words
        err = json.Unmarshal(body, &words)
        if err != nil {
            return nil, fmt.Errorf("unmarshal error for words: %s", err)
        }
        return words, nil

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
```

### Main Function

The main function is now simplified to:
1. Parse command-line arguments
2. Call the request function
3. Handle errors and display responses

```go
func getHTTPJsonMap() {
    args := os.Args

    if len(args) < 2 {
        fmt.Printf("Usage: ./api-client-parse-json <url>\n")
        os.Exit(1)
    }

    res, err := doRequest(args[1])
    if err != nil {
        log.Fatalf("Failed to make HTTP request: %s\n", err)
    }

    if res == nil {
        fmt.Printf("No response received.\n")
        os.Exit(1)
    }

    fmt.Printf("Response: %s\n", res.GetResponse())
}
```

## Error Handling Approach

The decoupled code uses Go's idiomatic error handling:

1. **Return errors, don't panic**: Instead of using `log.Fatal`, we return errors to the caller
2. **Descriptive errors**: Errors include context about where they occurred
3. **Proper resource cleanup**: Using `defer` to close resources
4. **Error propagation**: Errors bubble up to where they can be properly handled

The instructor explains:
> "So instead of exiting our program within our function we're actually going to always return an error."

This approach allows the caller to decide how to handle errors, rather than the function making that decision.

## Python Equivalent Implementation

Here's how we could implement the same approach in Python using abstract base classes (ABC) for the interface concept:

```python
from abc import ABC, abstractmethod
import requests
import json
import sys
from typing import Optional, Dict, List, Union, Tuple

# Response interface equivalent
class Response(ABC):
    @abstractmethod
    def get_response(self) -> str:
        pass

# Words type
class Words(Response):
    def __init__(self, input_str: str, words: List[str]):
        self.input = input_str
        self.words = words
    
    def get_response(self) -> str:
        return f"Words: {', '.join(self.words)}"

# Occurrence type
class Occurrence(Response):
    def __init__(self, words: Dict[str, int]):
        self.words = words
    
    def get_response(self) -> str:
        formatted_words = [f"{word}: {count}" for word, count in self.words.items()]
        return f"Words: {', '.join(formatted_words)}"

# Page type to determine response type
class Page:
    def __init__(self, name: str):
        self.name = name

def do_request(request_url: str) -> Tuple[Optional[Response], Optional[str]]:
    # Validate URL (simplified)
    if not request_url.startswith(('http://', 'https://')):
        return None, f"Validation error: URL is not valid: {request_url}"
    
    try:
        response = requests.get(request_url)
        response.raise_for_status()  # Raise exception for 4XX/5XX responses
        
        data = response.json()
        
        # Parse page type
        page = Page(data.get("page", ""))
        
        if page.name == "words":
            return Words(data.get("input", ""), data.get("words", [])), None
        elif page.name == "occurrence":
            return Occurrence(data.get("words", {})), None
        else:
            return None, None
            
    except requests.exceptions.RequestException as e:
        return None, f"HTTP request error: {str(e)}"
    except json.JSONDecodeError as e:
        return None, f"JSON decode error: {str(e)}"
    except Exception as e:
        return None, f"Unexpected error: {str(e)}"

def main():
    if len(sys.argv) < 2:
        print("Usage: python script.py <url>")
        sys.exit(1)
    
    response, error = do_request(sys.argv[1])
    
    if error:
        print(f"Error: {error}")
        sys.exit(1)
    
    if not response:
        print("No response received.")
        sys.exit(1)
    
    print(f"Response: {response.get_response()}")

if __name__ == "__main__":
    main()
```

### Key Differences Between Go and Python Implementations:

1. **Interface Implementation**:
   - Go uses implicit interfaces (types implement interfaces automatically if they have the required methods)
   - Python uses explicit abstract base classes that must be inherited

2. **Error Handling**:
   - Go uses multiple return values (value, error)
   - Python typically uses exceptions, but we've mimicked Go's approach with a tuple return

3. **Type System**:
   - Go is statically typed
   - Python is dynamically typed, but we've added type hints for clarity

4. **JSON Handling**:
   - Go requires explicit unmarshaling into structs
   - Python's `requests` library handles JSON parsing automatically

## Conclusion

Decoupling code in Go provides several benefits:

1. **Improved maintainability** through separation of concerns
2. **Enhanced flexibility** via interfaces for handling multiple response types
3. **Better error handling** by returning errors rather than terminating the program
4. **Cleaner code structure** with dedicated functions for specific tasks

This approach aligns with modern software development principles and produces code that is easier to test, maintain, and extend.
