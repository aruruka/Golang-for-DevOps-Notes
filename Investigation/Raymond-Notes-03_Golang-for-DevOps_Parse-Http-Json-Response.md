# Go for DevOps: Parsing HTTP JSON Responses

## Overview

This document covers the fundamental concepts and implementation patterns for parsing JSON responses from HTTP APIs in Go. It demonstrates how to transform raw JSON data into structured Go types using structs and the `encoding/json` package.

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [JSON Structure Analysis](#json-structure-analysis)
3. [Go Struct Definition](#go-struct-definition)
4. [JSON Unmarshaling Process](#json-unmarshaling-process)
5. [Error Handling Patterns](#error-handling-patterns)
6. [Data Processing and Output](#data-processing-and-output)
7. [Best Practices](#best-practices)
8. [Code Example](#complete-code-example)

## Core Concepts

### JSON vs Go Types Mapping

Understanding the relationship between JSON data types and Go types is crucial:

- **JSON String** → **Go string**
- **JSON Number** → **Go int, float64, etc.**
- **JSON Array** → **Go slice []T**
- **JSON Object** → **Go struct**
- **JSON Boolean** → **Go bool**
- **JSON null** → **Go pointer types or interface{}**

### Struct vs Other Data Structures

In Go, **structs** are composite types that group related data together, similar to:
- **Objects** in JavaScript/JSON
- **Classes** in OOP languages (but without methods by default)
- **Records** in functional programming

Key differences from other Go types:
- **Arrays/Slices**: Homogeneous collections of the same type
- **Maps**: Key-value pairs with dynamic keys
- **Structs**: Fixed fields with known names and types at compile time

## JSON Structure Analysis

Given this sample JSON response:
```json
{
  "page": "words",
  "input": "word1", 
  "words": ["word1", "word2"]
}
```

We need to identify:
- **page**: String field
- **input**: String field  
- **words**: Array of strings

## Go Struct Definition

### Basic Struct Declaration

```go
type Words struct {
    Page  string   `json:"page"`
    Input string   `json:"input"`
    Words []string `json:"words"`
}
```

### Critical Design Decisions

#### 1. Field Visibility (Exported vs Unexported)

**Why uppercase field names are required:**

```go
// ❌ WRONG - Unexported fields
type Words struct {
    page  string   `json:"page"`  // lowercase = unexported
    input string   `json:"input"` // JSON package cannot access
    words []string `json:"words"`
}

// ✅ CORRECT - Exported fields  
type Words struct {
    Page  string   `json:"page"`  // uppercase = exported
    Input string   `json:"input"` // JSON package can access
    Words []string `json:"words"`
}
```

**Explanation**: The `encoding/json` package is external to your package. It can only access **exported** (capitalized) fields due to Go's visibility rules. Unexported fields are invisible to external packages.

#### 2. JSON Struct Tags

**Purpose of struct tags:**
```go
Page string `json:"page"`
//           ^^^^^^^^^^^^
//           This maps Go field "Page" to JSON key "page"
```

The struct tag `json:"page"` tells the JSON package:
- When **unmarshaling**: Look for JSON key "page" and put its value in the "Page" field
- When **marshaling**: Use "page" as the JSON key for the "Page" field

## JSON Unmarshaling Process

### Understanding Pointers in Unmarshaling

**Why `json.Unmarshal()` requires a pointer:**

```go
var words Words
err := json.Unmarshal(body, &words)  // &words = pointer to words
//                          ^
//                          Address-of operator
```

**Technical reasoning:**
1. `json.Unmarshal()` needs to **modify** the struct
2. Go is **pass-by-value** by default
3. Without a pointer, `Unmarshal()` would receive a **copy**
4. Changes to a copy don't affect the original variable
5. Using `&words` passes the **memory address**, allowing direct modification

**Function signature analysis:**
```go
func Unmarshal(data []byte, v interface{}) error
//                          ^
//                          Must be a pointer to modify the original
```

### Complete Unmarshaling Implementation

```go
// Declare variable to hold parsed data
var words Words

// Unmarshal JSON into struct
err = json.Unmarshal(body, &words)
if err != nil {
    log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
}
```

## Error Handling Patterns

### HTTP Status Code Validation

```go
if response.StatusCode != http.StatusOK {
    log.Fatalf("Received (HTTP Code %d) response: %s\n", response.StatusCode, body)
}
```

**Why check status codes:**
- **200**: Success - safe to parse JSON
- **4xx**: Client errors - likely invalid request
- **5xx**: Server errors - temporary issues
- **Other codes**: Redirects, etc. - may not contain expected JSON

### Comprehensive Error Handling Chain

```go
// 1. HTTP request error
response, err := http.Get(args[1])
if err != nil {
    log.Fatalf("Failed to make HTTP GET request: %s\n", err)
}

// 2. Response body reading error  
body, err := io.ReadAll(response.Body)
if err != nil {
    log.Fatalf("Failed to read response body: %s\n", err)
}

// 3. HTTP status code validation
if response.StatusCode != http.StatusOK {
    log.Fatalf("Received (HTTP Code %d) response: %s\n", response.StatusCode, body)
}

// 4. JSON parsing error
err = json.Unmarshal(body, &words)
if err != nil {
    log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
}
```

## Data Processing and Output

### String Array Processing

**Converting slice to comma-separated string:**

```go
fmt.Printf("Words: %s\n", strings.Join(words.Words, ", "))
```

**Why use `strings.Join()`:**
- **Raw slice output**: `[word1 word2]` (with brackets)
- **Joined output**: `word1, word2` (clean, readable)
- **Alternative methods**: Manual iteration (more verbose)

### Accessing Struct Fields

```go
// Access individual fields
fmt.Printf("Page: %s\n", words.Page)
fmt.Printf("Input: %s\n", words.Input)

// Process array elements
for i, word := range words.Words {
    fmt.Printf("Word %d: %s\n", i+1, word)
}
```

## Best Practices

### 1. Struct Design

```go
// ✅ Good: Clear, descriptive names
type APIResponse struct {
    Status    string   `json:"status"`
    Message   string   `json:"message"`
    Data      []string `json:"data"`
    Timestamp int64    `json:"timestamp"`
}

// ❌ Avoid: Unclear abbreviations  
type Resp struct {
    St string   `json:"status"`
    Msg string  `json:"message"`
    D  []string `json:"data"`
}
```

### 2. Error Context

```go
// ✅ Provide context in error messages
if err != nil {
    log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
}

// ❌ Generic error messages
if err != nil {
    log.Fatalf("Error: %s\n", err)
}
```

### 3. Resource Management

```go
response, err := http.Get(url)
if err != nil {
    return err
}
defer response.Body.Close()  // ✅ Always close response body
```

### 4. Input Validation

```go
// Validate URL format before making request
if _, err := url.ParseRequestURI(args[1]); err != nil {
    fmt.Printf("URL is in invalid format: %s\n", err)
    os.Exit(1)
}
```

## Complete Code Example

```go
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

// Define struct to match JSON structure
type Words struct {
    Page  string   `json:"page"`
    Input string   `json:"input"`
    Words []string `json:"words"`
}

func fetchWordsFromAPI() {
    // Validate command line arguments
    args := os.Args
    if len(args) < 2 {
        fmt.Printf("Usage: ./api-client <url>\n")
        os.Exit(1)
    }
    
    // Validate URL format
    if _, err := url.ParseRequestURI(args[1]); err != nil {
        fmt.Printf("URL is in invalid format: %s\n", err)
        os.Exit(1)
    }

    // Make HTTP GET request
    response, err := http.Get(args[1])
    if err != nil {
        log.Fatalf("Failed to make HTTP GET request: %s\n", err)
    }
    defer response.Body.Close()

    // Read response body
    body, err := io.ReadAll(response.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %s\n", err)
    }

    // Check HTTP status code
    if response.StatusCode != http.StatusOK {
        log.Fatalf("Received (HTTP Code %d) response: %s\n", response.StatusCode, body)
    }

    // Parse JSON into struct
    var words Words
    err = json.Unmarshal(body, &words)
    if err != nil {
        log.Fatalf("Failed to unmarshal JSON response: %s\n", err)
    }

    // Output parsed data
    fmt.Printf("JSON Parsed:\nPage: %s\nWords: %s\n", 
        words.Page, 
        strings.Join(words.Words, ", "))
}
```

## Key Takeaways

1. **Struct fields must be exported** (capitalized) for JSON package access
2. **Struct tags map Go fields to JSON keys** using `json:"key_name"`
3. **Pointers are required for unmarshaling** to allow modification of original data
4. **Always validate HTTP status codes** before parsing JSON
5. **Handle errors at each step** of the HTTP → JSON → Struct pipeline
6. **Use appropriate string processing** for clean output formatting

## Advanced Topics

- **Custom JSON unmarshaling** with `UnmarshalJSON()` methods
- **Nested structs** for complex JSON objects
- **Interface{} usage** for dynamic JSON structures
- **JSON streaming** for large datasets
- **Marshal/Unmarshal performance** optimization techniques

## External Tools

For complex JSON structures, consider using online tools:
- **JSON-to-Go converters**: Automatically generate struct definitions
- **Search**: "JSON to struct Golang" for various web-based tools
- **Benefit**: Saves time on large, complex JSON schemas
