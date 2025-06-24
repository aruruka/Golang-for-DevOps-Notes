# Assignment 1: JSON Parsing Implementation

This directory contains the implementation for Assignment 1 - JSON Parsing from the Golang for DevOps course.

## Project Structure

```
assignments/assignment-1-json-parsing/
├── go.mod                          # Go module definition
├── requirements.md                 # Assignment requirements
├── README.md                      # This file
├── cmd/
│   └── assignment1/
│       └── main.go                # Main application entry point
└── pkg/
    └── api/
        ├── init.go                # Interfaces and initialization
        ├── error.go               # Custom error types
        ├── assignment.go          # Core assignment logic
        └── assignment_test.go     # Unit tests
```

## Features

### Core Implementation
- **HTTP Client Interface**: Uses dependency injection for testability
- **JSON Parsing**: Handles complex JSON structures with mixed data types
- **Error Handling**: Custom error types with detailed context
- **Null Value Handling**: Properly handles null values in JSON arrays
- **Mixed Type Arrays**: Supports arrays with different data types

### Data Structure Support
The implementation correctly parses the following JSON structure:
- `page`: String field
- `words`: Array of strings
- `percentages`: Map of string to float64
- `special`: Array with null values (using pointers)
- `extraSpecial`: Array with mixed data types (interface{})

### Testing
- **Unit Tests**: Comprehensive test coverage using mock HTTP clients
- **Error Scenarios**: Tests for HTTP errors and invalid JSON
- **Data Validation**: Verifies correct parsing of all data types
- **Mock Implementation**: Follows established patterns from the course

## Usage

### Running the Application
```bash
cd assignments/assignment-1-json-parsing
go run cmd/assignment1/main.go
```

### Running Tests
```bash
cd assignments/assignment-1-json-parsing
go test ./pkg/api -v
```

### Starting the Test Server
Before running the application, ensure the test server is running:
```bash
cd test-server
go run main.go
```

## Implementation Details

### Interface Design
The implementation follows the established patterns from the course:
- `ClientIface`: HTTP client interface for dependency injection
- `APIIface`: Main API interface
- `Response`: Interface for different response types

### Mock Testing
The tests use mock HTTP clients that implement `ClientIface`:
- `MockClient`: Provides controllable HTTP responses
- No actual network calls during testing
- Predictable test data for validation

### Error Handling
Custom `RequestError` type provides:
- HTTP status codes
- Response body content
- Detailed error messages
- Context for debugging

## Key Learning Objectives Achieved

1. **HTTP Client Usage**: Making GET requests and handling responses
2. **JSON Unmarshaling**: Parsing complex JSON into Go structs
3. **Interface Design**: Creating testable code with dependency injection
4. **Error Handling**: Comprehensive error handling and custom error types
5. **Unit Testing**: Writing effective tests with mocks
6. **Mixed Data Types**: Handling arrays with different data types and null values

## Course Integration

This implementation demonstrates the core concepts from the "Golang for DevOps and Cloud Engineers" course:
- Clean architecture with package separation
- Interface-based design for testability
- Proper error handling patterns
- Comprehensive unit testing
- Real-world HTTP client usage
