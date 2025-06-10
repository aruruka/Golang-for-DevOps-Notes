# Custom Error Handling in Go

## 1. Introduction to Error Handling in Go

Unlike many programming languages that use exceptions, Go handles errors as values that are returned from functions. This approach makes error handling explicit and part of the normal control flow.

The standard pattern in Go for functions that can fail is to return two values: the result and an error.

```go
func doSomething() (Result, error) {
    // Function implementation
    if somethingWentWrong {
        return nil, errors.New("something went wrong")
    }
    return result, nil
}
```

This pattern is visible throughout the Go standard library and is considered idiomatic Go.

## 2. Understanding Go's Error Interface

In Go, errors are values that implement the built-in `error` interface. This interface is defined in the language as:

```go
type error interface {
    Error() string
}
```

Any type that implements this method can be used as an error. This simple but powerful design allows for creating custom error types with additional context while still being compatible with the standard error handling pattern.

## 3. Creating Custom Error Types

While the standard error interface is sufficient for simple cases, real-world applications often require more context when errors occur. This is where custom error types become valuable.

In our example, we created a custom `RequestError` type to capture additional information about HTTP request failures:

```go
type RequestError struct {
    HTTPCode int
    Body     string
    Err      string
}

func (r RequestError) Error() string {
    return r.Err
}
```

This custom error type:
- Stores the HTTP status code
- Captures the response body
- Includes an error message
- Implements the `Error()` method to satisfy the `error` interface

The instructor chose this approach because it allows us to:
1. Return a standard error that any Go code can handle
2. Provide additional context that can be accessed when needed
3. Keep the error handling interface consistent with Go's conventions

## 4. Using Custom Error Types in API Clients

Our example demonstrates using custom error types in an HTTP API client. When an error occurs during JSON parsing, we return a `RequestError` with context:

```go
if !json.Valid(body) {
    return nil, RequestError{
        HTTPCode: response.StatusCode,
        Body:     string(body),
        Err:      "no valid json returned",
    }
}
```

Similarly, when unmarshaling JSON fails:

```go
err = json.Unmarshal(body, &page)
if err != nil {
    return nil, RequestError{
        HTTPCode: response.StatusCode,
        Body:     string(body),
        Err:      fmt.Sprintf("page unmarshal error: %s", err),
    }
}
```

## 5. Type Assertion for Enhanced Error Handling

To access the additional information in our custom error type, we use type assertion:

```go
if err != nil {
    if reqErr, ok := err.(RequestError); ok {
        fmt.Printf("Error: %s (HTTP Code: %d, Body: %s)\n", reqErr.Err, reqErr.HTTPCode, reqErr.Body)
        os.Exit(1)
    }
}
```

The type assertion `err.(RequestError)` attempts to convert the error to our custom type:
- If successful, `ok` is true and `reqErr` contains our custom error
- If not successful, `ok` is false, and we can handle it as a standard error

This pattern allows us to:
1. Handle specific error types differently
2. Access additional context when available
3. Gracefully fall back to standard error handling

## 6. Decoupling Business Logic with Interfaces

The example also demonstrates how to decouple business logic using interfaces. The `Response` interface allows different response types to implement a common method:

```go
type Response interface {
    GetResponse() string
}
```

Two different types implement this interface:

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

This approach allows the client code to work with different response types through a common interface, without needing to know the specific implementation details.

## 7. Python Alternative Implementation

Here's how we would implement the same approach in Python:

```python
from abc import ABC, abstractmethod
import json
import sys
import urllib.request
from urllib.parse import urlparse

class RequestError(Exception):
    """Custom exception for HTTP request errors with additional context"""
    
    def __init__(self, http_code, body, err):
        self.http_code = http_code
        self.body = body
        self.err = err
        super().__init__(err)
    
    def __str__(self):
        return self.err

class Response(ABC):
    """Interface for different response types"""
    
    @abstractmethod
    def get_response(self):
        pass

class Words(Response):
    """Response type for word lists"""
    
    def __init__(self, input_str, words):
        self.input = input_str
        self.words = words
    
    def get_response(self):
        return f"Words: {', '.join(self.words)}"

class Occurrence(Response):
    """Response type for word occurrences"""
    
    def __init__(self, words):
        self.words = words
    
    def get_response(self):
        word_strings = [f"{word}: {count}" for word, count in self.words.items()]
        return f"Words: {', '.join(word_strings)}"

def do_request(request_url):
    """Make HTTP request and parse JSON response"""
    try:
        # Validate URL format
        result = urlparse(request_url)
        if not all([result.scheme, result.netloc]):
            print(f"URL is in invalid format")
            sys.exit(1)
        
        # Make HTTP request
        with urllib.request.urlopen(request_url) as response:
            body = response.read().decode('utf-8')
            status_code = response.status
            
            # Check status code
            if status_code != 200:
                raise RequestError(
                    status_code, 
                    body, 
                    f"invalid output (http code: {status_code}): {body}"
                )
            
            # Validate JSON
            try:
                data = json.loads(body)
            except json.JSONDecodeError:
                raise RequestError(status_code, body, "no valid json returned")
            
            # Parse page type
            if "page" not in data:
                raise RequestError(status_code, body, "page field missing")
            
            page_name = data.get("page")
            
            # Handle different response types
            if page_name == "words":
                try:
                    return Words(data.get("input", ""), data.get("words", []))
                except Exception as e:
                    raise RequestError(status_code, body, f"words unmarshal error: {str(e)}")
            
            elif page_name == "occurrence":
                try:
                    return Occurrence(data.get("words", {}))
                except Exception as e:
                    raise RequestError(status_code, body, f"occurrence unmarshal error: {str(e)}")
            
            return None
    
    except urllib.error.URLError as e:
        raise Exception(f"get error: {str(e)}")

def get_http_json_map():
    """Main function to handle command line arguments and make request"""
    args = sys.argv
    
    if len(args) < 2:
        print("Usage: python api_client.py <url>")
        sys.exit(1)
    
    try:
        res = do_request(args[1])
        
        if res is None:
            print("No response received.")
            sys.exit(1)
        
        print(f"Response: {res.get_response()}")
    
    except RequestError as e:
        print(f"Error: {e} (HTTP Code: {e.http_code}, Body: {e.body})")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    get_http_json_map()
```

## Key Differences Between Go and Python Implementations

1. **Error Handling:**
   - Go uses return values for errors
   - Python uses exceptions (though we've designed our `RequestError` class to contain similar information)

2. **Interfaces:**
   - Go has explicit interfaces that types implicitly implement
   - Python uses abstract base classes (ABC) to define interfaces explicitly

3. **Type Assertions:**
   - Go uses type assertions (`err.(RequestError)`)
   - Python uses `isinstance()` checks (implied in the exception catch)

4. **Memory Management:**
   - Go requires explicit resource cleanup with `defer`
   - Python's context manager (`with` statement) handles resource cleanup

## Conclusion

Custom error types provide a powerful way to enhance error handling with additional context while maintaining compatibility with standard error handling patterns. Whether in Go or Python, the approach of:

1. Creating custom error types with additional fields
2. Using interfaces to decouple business logic
3. Implementing type-specific behaviors

results in more maintainable and informative code. This pattern is particularly valuable in applications that interact with external systems like HTTP APIs, where errors can occur at multiple levels and additional context is crucial for debugging.
