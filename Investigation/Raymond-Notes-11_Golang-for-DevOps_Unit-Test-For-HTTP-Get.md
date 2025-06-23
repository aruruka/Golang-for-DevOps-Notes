# Golang for DevOps: Unit Testing for HTTP Get Requests

## Table of Contents
- [Introduction](#introduction)
- [Overview of Testing Challenges](#overview-of-testing-challenges)
- [Code Architecture for Testability](#code-architecture-for-testability)
  - [Dependency Injection with Interfaces](#dependency-injection-with-interfaces)
  - [Mock Implementation Strategy](#mock-implementation-strategy)
- [HTTP Get Request Testing](#http-get-request-testing)
  - [Original Implementation](#original-implementation)
  - [Testing-Ready Implementation](#testing-ready-implementation)
  - [Test Implementation](#test-implementation)
- [HTTP Transport Layer Testing](#http-transport-layer-testing)
  - [Transport Implementation](#transport-implementation)
  - [Transport Test Implementation](#transport-test-implementation)
- [Key Design Principles](#key-design-principles)
- [Best Practices for Unit Testing HTTP Services](#best-practices-for-unit-testing-http-services)
- [Common Patterns and Techniques](#common-patterns-and-techniques)
- [Summary](#summary)

## Introduction

Unit testing HTTP-based services in Go presents unique challenges, particularly when dealing with external API calls and authentication layers. This document explores comprehensive strategies for testing HTTP Get requests and custom HTTP transport layers in Go, based on practical implementations and instructor insights from the "Golang for DevOps and Cloud Engineers" course.

The key challenge in testing HTTP services is **avoiding actual network calls** while still testing the business logic. This requires careful design decisions that enable **dependency injection** and **mock implementations**.

## Overview of Testing Challenges

When testing HTTP-based code, developers face several challenges:

1. **External Dependencies**: Real HTTP calls depend on external services
2. **Network Reliability**: Tests should not fail due to network issues
3. **Test Performance**: Network calls slow down test execution
4. **Test Isolation**: Tests should not affect external systems
5. **Authentication Complexity**: JWT tokens and auth flows add complexity

As the instructor explains:
> "Tests are a bit difficult, a little bit more difficult to write when we are using API calls. You are connecting to a server, so we will have to kind of intercept those connections and make sure that we populate the correct variables with information that we would get if we would make a real connection to test our functions."

## Code Architecture for Testability

### Dependency Injection with Interfaces

The foundation of testable HTTP code is **interface-based dependency injection**. This allows swapping real HTTP clients with mock implementations during testing.

#### Interface Design

```go
type ClientIface interface {
    Get(url string) (resp *http.Response, err error)
    Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
```

**Key Design Decisions:**
- **Interface over Concrete Types**: Use `ClientIface` instead of `*http.Client` directly
- **Minimal Interface**: Only include methods actually used by the code
- **Standard Signatures**: Match the signatures of `http.Client` methods exactly

#### API Structure with Dependency Injection

```go
type api struct {
    Options Options
    Client  ClientIface  // Interface, not *http.Client
}

func New(options Options) APIIface {
    return api{
        Options: options,
        Client: &http.Client{  // Real implementation for production
            Transport: MyJWTTransport{
                transport:  http.DefaultTransport,
                password:   options.Password,
                loginURL:   options.LoginURL,
                HTTPClient: &http.Client{},
            },
        },
    }
}
```

**Why This Design Works:**
- **Production Use**: `New()` returns a real `http.Client` wrapped in the interface
- **Test Use**: Tests can inject mock implementations of `ClientIface`
- **Type Safety**: Interface ensures mock implementations match required methods

### Mock Implementation Strategy

#### Mock Client Implementation

```go
type MockClient struct {
    GetResponse  *http.Response
    PostResponse *http.Response
}

func (m MockClient) Get(url string) (resp *http.Response, err error) {
    if url == "http://localhost/login" {
        fmt.Printf("Login endpoint")
    }
    return m.GetResponse, nil
}

func (m MockClient) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
    return m.PostResponse, nil
}
```

**Mock Design Principles:**
- **Predictable Responses**: Return pre-configured responses
- **No Network Calls**: All responses are fabricated
- **Flexible Configuration**: Allow different responses for different test scenarios
- **Interface Compliance**: Implement all required interface methods

## HTTP Get Request Testing

### Original Implementation

The original `DoGetRequest` function handles JSON API responses with different page types:

```go
func (a api) DoGetRequest(requestURL string) (Response, error) {
    response, err := a.Client.Get(requestURL)
    
    if err != nil {
        return nil, fmt.Errorf("Get error: %s", err)
    }
    
    defer response.Body.Close()
    
    body, err := io.ReadAll(response.Body)
    
    if err != nil {
        return nil, fmt.Errorf("ReadAll error: %s", err)
    }
    
    if response.StatusCode != 200 {
        return nil, fmt.Errorf("Invalid output (HTTP Code %d): %s\n", response.StatusCode, string(body))
    }
    
    var page Page
    
    if !json.Valid(body) {
        return nil, RequestError{
            Err:      fmt.Sprintf("Response is not a json"),
            HTTPCode: response.StatusCode,
            Body:     string(body),
        }
    }
    
    err = json.Unmarshal(body, &page)
    if err != nil {
        return nil, RequestError{
            Err:      fmt.Sprintf("Page unmarshal error: %s", err),
            HTTPCode: response.StatusCode,
            Body:     string(body),
        }
    }
    
    switch page.Name {
    case "words":
        var words Words
        err = json.Unmarshal(body, &words)
        if err != nil {
            return nil, fmt.Errorf("Words unmarshal error: %s", err)
        }
        return words, nil
    case "occurrence":
        var occurrence Occurrence
        err = json.Unmarshal(body, &occurrence)
        if err != nil {
            return nil, fmt.Errorf("Occurrence unmarshal error: %s", err)
        }
        return occurrence, nil
    }
    
    return nil, nil
}
```

**Key Business Logic to Test:**
1. HTTP status code validation
2. JSON validity checking
3. Page type detection and routing
4. Type-specific JSON unmarshaling
5. Error handling at each step

### Testing-Ready Implementation

To make the code testable, the key change is using `ClientIface` instead of `*http.Client`:

```go
type api struct {
    Options Options
    Client  ClientIface  // This is the key change for testability
}
```

**Changes Made for Testability:**
1. **Interface Introduction**: `ClientIface` replaces direct `*http.Client` usage
2. **Dependency Injection**: Client can be injected during struct initialization
3. **Pointer Handling**: Use `&http.Client{}` to match interface pointer receiver requirements

### Test Implementation

#### Creating Test Data

The test requires creating realistic JSON responses that match the API contract:

```go
func TestDoGetRequest(t *testing.T) {
    // Create test data that matches expected API response
    words := WordsPage{
        Page: Page{"words"},
        Words: Words{
            Input: "abc",
            Words: []string{"a", "b"},
        },
    }
    
    // Marshal to JSON bytes
    wordsBytes, err := json.Marshal(words)
    if err != nil {
        t.Errorf("marshal error: %s", err)
    }
```

**Test Data Strategy:**
- **Composite Struct**: `WordsPage` combines `Page` and `Words` for complete JSON structure
- **Realistic Data**: Use data that would actually be returned by the API
- **Error Handling**: Handle marshaling errors in test setup

#### Mock HTTP Response Creation

```go
    apiInstance := api{
        Options: Options{},
        Client: MockClient{
            GetResponse: &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(bytes.NewReader(wordsBytes)),
            },
        },
    }
```

**Critical HTTP Response Details:**
- **StatusCode**: Must be 200 to pass validation
- **Body**: Use `io.NopCloser` to convert bytes to `io.ReadCloser`
- **bytes.NewReader**: Converts byte slice to `io.Reader`

**Why `io.NopCloser`?**
As the instructor explains:
> "To return a ReadCloser we need to use the io.NopCloser package. NopCloser returns a ReadCloser with a no-op Close method wrapping the provided Reader."

#### Test Execution and Validation

```go
    response, err := apiInstance.DoGetRequest("http://localhost/words")
    if err != nil {
        t.Errorf("DoGetRequest error: %s", err)
    }
    if response == nil {
        t.Errorf("Response is nil")
    }
    if response.GetResponse() != `Words: a, b` {
        t.Errorf("Got wrong output: %s", response.GetResponse())
    }
```

**Test Validation Strategy:**
1. **Error Checking**: Ensure no errors during execution
2. **Nil Checking**: Verify response is not nil
3. **Output Validation**: Check the formatted output matches expectations

## HTTP Transport Layer Testing

### Transport Implementation

The JWT transport layer adds authentication headers to HTTP requests:

```go
type MyJWTTransport struct {
    transport http.RoundTripper
    token     string
    password  string
    loginURL  string
}

func (m *MyJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    if m.token == "" {
        if m.password != "" {
            token, err := doLoginRequest(http.Client{}, m.loginURL, m.password)
            if err != nil {
                return nil, err
            }
            m.token = token
        }
    }
    if m.token != "" {
        req.Header.Add("Authorization", "Bearer "+m.token)
    }
    return m.transport.RoundTrip(req)
}
```

**Transport Layer Responsibilities:**
1. **Token Management**: Acquire token if not present
2. **Header Injection**: Add Authorization header to requests
3. **Request Forwarding**: Delegate to underlying transport

### Transport Test Implementation

#### Mock RoundTripper

```go
type MockRoundTripper struct {
    RoundTripOutput *http.Response
}

func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    if req.Header.Get("Authorization") != "Bearer 123" {
        return nil, fmt.Errorf("wrong Authorization header: %s", req.Header.Get("Authorization"))
    }
    return m.RoundTripOutput, nil
}
```

**Mock RoundTripper Features:**
- **Header Validation**: Verify correct Authorization header is set
- **Error Simulation**: Return errors for invalid headers
- **Configurable Response**: Return pre-configured response

#### Complete Transport Test

```go
func TestRoundtrip(t *testing.T) {
    loginResponse := LoginResponse{
        Token: "123",
    }
    loginResponseBytes, err := json.Marshal(loginResponse)
    if err != nil {
        t.Errorf("marshal error: %s", err)
    }

    jwtTransport := MyJWTTransport{
        HTTPClient: MockClient{
            PostResponse: &http.Response{
                StatusCode: 200,
                Body:       io.NopCloser(bytes.NewReader(loginResponseBytes)),
            },
        },
        transport: MockRoundTripper{
            RoundTripOutput: &http.Response{
                StatusCode: 200,
            },
        },
        password: "xyz",
    }
    
    req := &http.Request{
        Header: make(http.Header),
    }
    
    res, err := jwtTransport.RoundTrip(req)
    if err != nil {
        t.Errorf("got error: %s", err)
        t.FailNow()
    }
    if res.StatusCode != 200 {
        t.Errorf("expected status code 200, got %d", res.StatusCode)
    }
}
```

**Test Flow:**
1. **Setup**: Create mock login response with token
2. **Transport Configuration**: Configure transport with mocks
3. **Request Creation**: Create HTTP request with headers
4. **Execution**: Call RoundTrip method
5. **Validation**: Verify response and status code

**Important Header Initialization:**
```go
req := &http.Request{
    Header: make(http.Header),  // Critical: Initialize header map
}
```

As the instructor notes:
> "If a map is nil we cannot assign anything to the map. So we need to instantiate the header."

## Key Design Principles

### 1. Interface Segregation

**Principle**: Create small, focused interfaces that only expose required methods.

```go
type ClientIface interface {
    Get(url string) (resp *http.Response, err error)
    Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
```

**Why This Works:**
- **Minimal Surface Area**: Only includes methods actually used
- **Easy Mocking**: Fewer methods to implement in mocks
- **Clear Dependencies**: Explicitly shows what HTTP methods are needed

### 2. Dependency Injection

**Principle**: Accept dependencies through interfaces, not concrete types.

```go
// ❌ Hard to test
type api struct {
    Client *http.Client
}

// ✅ Easy to test
type api struct {
    Client ClientIface
}
```

**Benefits:**
- **Testability**: Can inject mocks during testing
- **Flexibility**: Can swap implementations without code changes
- **Isolation**: Tests don't depend on external services

### 3. Mock-Friendly Design

**Principle**: Design code to be easily mockable.

```go
type MockClient struct {
    GetResponse  *http.Response
    PostResponse *http.Response
}
```

**Mock Design Guidelines:**
- **Configurable Responses**: Allow different responses for different scenarios
- **Validation Logic**: Include validation in mocks to verify correct usage
- **Error Simulation**: Support error scenarios for comprehensive testing

### 4. Realistic Test Data

**Principle**: Use test data that closely resembles production data.

```go
words := WordsPage{
    Page: Page{"words"},
    Words: Words{
        Input: "abc",
        Words: []string{"a", "b"},
    },
}
```

**Why Realistic Data Matters:**
- **Integration Confidence**: Tests exercise the same code paths as production
- **JSON Validation**: Ensures JSON marshaling/unmarshaling works correctly
- **Type Safety**: Catches type mismatches between test and production data

## Best Practices for Unit Testing HTTP Services

### 1. Test Structure Organization

```go
func TestDoGetRequest(t *testing.T) {
    // 1. Setup test data
    words := WordsPage{...}
    wordsBytes, err := json.Marshal(words)
    
    // 2. Configure system under test
    apiInstance := api{
        Client: MockClient{...},
    }
    
    // 3. Execute function
    response, err := apiInstance.DoGetRequest("http://localhost/words")
    
    // 4. Validate results
    if err != nil {
        t.Errorf("DoGetRequest error: %s", err)
    }
    // Additional assertions...
}
```

### 2. Error Handling in Tests

```go
// Check for errors
if err != nil {
    t.Errorf("DoGetRequest error: %s", err)
}

// Stop execution on critical failures
if response == nil {
    t.Errorf("Response is nil")
}

// Use t.FailNow() to stop test immediately
if res.StatusCode != 200 {
    t.Errorf("expected status code 200, got %d", res.StatusCode)
    t.FailNow()
}
```

### 3. Mock Validation

```go
func (m MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    // Validate that the code under test behaves correctly
    if req.Header.Get("Authorization") != "Bearer 123" {
        return nil, fmt.Errorf("wrong Authorization header: %s", req.Header.Get("Authorization"))
    }
    return m.RoundTripOutput, nil
}
```

**Mock Validation Benefits:**
- **Behavior Verification**: Ensures code sets headers correctly
- **Input Validation**: Verifies correct data is passed to dependencies
- **Edge Case Testing**: Can simulate various error conditions

### 4. Test Coverage Strategy

As the instructor mentions:
> "We now already have 60% of our statements covered. You see, so everything in green is already covered. And so now you could actually write more tests to come close to a hundred. A hundred percent is not always possible, and I wouldn't really aim for a hundred percent. Somewhere between eighty and a hundred percent is definitely fine."

**Coverage Guidelines:**
- **80-100% Coverage**: Aim for high coverage without obsessing over 100%
- **Critical Path Focus**: Prioritize testing error handling and business logic
- **Edge Cases**: Include tests for boundary conditions and error scenarios

## Common Patterns and Techniques

### 1. Composite Test Structs

```go
type WordsPage struct {
    Page   // Embedded struct
    Words  // Embedded struct
}
```

**When to Use:**
- **API Response Simulation**: When test data needs to combine multiple JSON structures
- **Realistic Payloads**: When production responses include multiple data types
- **Type Safety**: When you want compile-time validation of test data structure

### 2. io.ReadCloser Handling

```go
// Convert byte slice to io.ReadCloser
Body: io.NopCloser(bytes.NewReader(wordsBytes))
```

**Pattern Breakdown:**
1. `bytes.NewReader(wordsBytes)` - Creates `io.Reader` from bytes
2. `io.NopCloser(...)` - Wraps reader to add `Close()` method
3. Result satisfies `io.ReadCloser` interface required by `http.Response`

### 3. HTTP Header Initialization

```go
req := &http.Request{
    Header: make(http.Header),  // Essential: Initialize the map
}
```

**Common Pitfall:** Forgetting to initialize `Header` map leads to nil pointer panics when adding headers.

### 4. Table-Driven Tests (Advanced Pattern)

For testing multiple scenarios:

```go
func TestDoGetRequest(t *testing.T) {
    tests := []struct {
        name           string
        responseData   interface{}
        expectedOutput string
        expectError    bool
    }{
        {
            name: "words response",
            responseData: WordsPage{
                Page:  Page{"words"},
                Words: Words{Input: "abc", Words: []string{"a", "b"}},
            },
            expectedOutput: "Words: a, b",
            expectError:    false,
        },
        // Additional test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation...
        })
    }
}
```

### 5. Mock Configuration Patterns

```go
// Pattern: Configurable mock responses
type MockClient struct {
    GetResponse  *http.Response
    PostResponse *http.Response
    GetError     error
    PostError    error
}

func (m MockClient) Get(url string) (*http.Response, error) {
    if m.GetError != nil {
        return nil, m.GetError
    }
    return m.GetResponse, nil
}
```

## Summary

Unit testing HTTP services in Go requires careful architectural decisions made **before** writing the actual business logic. The key insights from this comprehensive analysis are:

### Architectural Requirements

1. **Interface-Based Design**: Use interfaces instead of concrete types for all external dependencies
2. **Dependency Injection**: Accept dependencies through constructors or struct fields
3. **Separation of Concerns**: Separate HTTP transport logic from business logic

### Testing Strategies

1. **Mock Everything External**: Never make real HTTP calls in unit tests
2. **Realistic Test Data**: Use data structures that match production API responses
3. **Comprehensive Validation**: Test both success paths and error conditions

### Code Quality Principles

1. **Fail Fast**: Use `t.FailNow()` for critical test failures
2. **Clear Assertions**: Validate specific outputs, not just absence of errors
3. **Mock Validation**: Include validation logic in mocks to verify correct behavior

### Design Changes for Testability

The instructor's approach demonstrates that **testable code requires upfront design decisions**:

- **Original**: `Client *http.Client` (not testable)
- **Testable**: `Client ClientIface` (easily mockable)

This seemingly small change enables comprehensive testing without network dependencies, faster test execution, and reliable CI/CD pipelines.

The testing patterns shown here apply broadly to any Go service that interacts with external APIs, databases, or other networked services. The key is designing for testability from the beginning rather than retrofitting tests onto tightly-coupled code.

As the instructor emphasizes, the goal is not perfect test coverage but **confidence in code correctness** through well-designed, maintainable tests that accurately reflect production behavior.
