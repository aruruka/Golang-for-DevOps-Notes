# Golang for DevOps: Modularizing an HTTP Client with JWT Authentication

This document provides a comprehensive reference for modularizing a Go application that uses a JWT token for HTTP POST requests. It is based on a lecture series that refactors a simple command-line tool into a reusable package.

## 1. From Monolith to Modular: The "Why"

Initially, all the application logic resided in the `main` package. This approach is common for small tools but quickly becomes unmanageable as the codebase grows.

**Problem:**
- **Low Reusability:** The code for handling HTTP requests and JWT authentication is tightly coupled with the main application logic.
- **Difficult to Test:** It's hard to write unit tests for specific functionalities without running the entire application.
- **Poor Maintainability:** Changes to one part of the code can have unintended consequences in other parts.

**Solution:**
Refactor the code into a separate, reusable package. This aligns with the **Single Responsibility Principle**, where the `main` package is responsible for command-line parsing and orchestrating the application flow, while the new `api` package is responsible for all API interactions.

## 2. Project Structure for a Packaged Application

A well-organized project structure is crucial for maintainability. The instructor adopts a standard Go project layout:

```
http-login-packaged/
├── cmd/
│   └── http-login/
│       └── main.go
├── go.mod
└── pkg/
    └── api/
        ├── api.go
        ├── error.go
        ├── login.go
        └── transport.go
```

- **`cmd/`**: This directory contains the entry point of the application. Each subdirectory corresponds to a single executable.
- **`pkg/`**: This directory contains the reusable packages that can be imported by other applications.
- **`go.mod`**: This file defines the module's path and its dependencies.

The module path in [`go.mod`](http-login-packaged/go.mod:1) is set to a unique identifier, typically a URL to the code repository. This ensures that the package can be uniquely identified and imported.

```go
module http-login-packaged

go 1.24.2
```

## 3. The `api` Package: A Deep Dive

The `api` package encapsulates all the logic for making API calls.

### 3.1. The `APIface` Interface and the `New` Function

To promote loose coupling and testability, the package exposes an interface (`APIface`) and a `New` function to create instances of the API client. This follows the **Dependency Inversion Principle**.

```go
// pkg/api/api.go

package api

import (
	"net/http"
)

// APIface defines the contract for our API client.
type APIface interface {
	DoGetRequest(url string) (Response, error)
}

// API is the concrete implementation of the APIface.
type API struct {
	client *http.Client
	// ... other fields
}

// Options holds the configuration for the API client.
type Options struct {
	Password string
	LoginURL string
}

// New creates a new API client.
func New(opts Options) APIface {
	// ... implementation ...
}
```

### 3.2. Handling JWT Authentication with `http.RoundTripper`

The core of the authentication mechanism is a custom `http.RoundTripper`. An `http.RoundTripper` is an interface that defines a single method, `RoundTrip`, which executes a single HTTP transaction.

By creating a custom `RoundTripper`, we can intercept outgoing requests and modify them, in this case, to add the `Authorization` header with the JWT token.

#### The `myJWTTransport` Struct

The `myJWTTransport` struct holds the state needed for authentication, including the JWT token, password, and login URL.

```go
// pkg/api/transport.go

type myJWTTransport struct {
	transport http.RoundTripper
	token     string
	password  string
	loginURL  string
}
```

#### The `RoundTrip` Method: Where the Magic Happens

The `RoundTrip` method is the heart of the custom transport. It performs the following steps:

1.  **Check for a token:** If the `token` field is empty, it means we need to authenticate.
2.  **Login:** It calls the `login` function to get a new JWT token.
3.  **Store the token:** The new token is stored in the `token` field for subsequent requests.
4.  **Add the `Authorization` header:** The `Authorization: Bearer <token>` header is added to the original request.
5.  **Execute the request:** The request is passed to the underlying (default) transport to be sent to the server.

This is where the instructor places the Bearer token logic. By encapsulating it within the `myJWTTransport`, the authentication logic is completely transparent to the code that uses the `http.Client`. This is a great example of the **Single Responsibility Principle**.

### 3.3. The Pointer Correction: A Crucial Fix

In the initial implementation, the `myJWTTransport` was passed by value. This meant that when the `RoundTrip` method updated the `token` field, it was updating a *copy* of the struct, not the original. As a result, a new login request was made for every API call.

**The Fix:**

The `myJWTTransport` must be passed as a pointer. This ensures that all method calls operate on the *same* instance of the struct, and the `token` is correctly cached between requests.

```go
// pkg/api/api.go (in the New function)
client := &http.Client{
    Transport: &myJWTTransport{ // Pass as a pointer
        // ...
    },
}

// pkg/api/transport.go
func (t *myJWTTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    // ...
}
```

This correction is vital for performance and for adhering to the API's rate limits.

## 4. The `main` Package: The Application Entry Point

The [`main.go`](http-login-packaged/cmd/http-login/main.go:1) file is now much simpler. Its responsibilities are:

1.  **Parse command-line arguments:** It uses the `flag` package to get the URL and password from the user.
2.  **Instantiate the API client:** It calls `api.New` to create a new client instance.
3.  **Make the API call:** It calls the `DoGetRequest` method on the API client.
4.  **Print the result:** It prints the response or any errors to the console.

```go
// cmd/http-login/main.go

package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"http-login-packaged/pkg/api"
)

func main() {
	// ... flag parsing ...

	apiInstance := api.New(api.Options{
		Password: password,
		LoginURL: parsedURL.Scheme + "://" + parsedURL.Host + "/login",
	})

	res, err := apiInstance.DoGetRequest(parsedURL.String())
	if err != nil {
		// ... error handling ...
	}
	fmt.Printf("Response: %s\n", res.GetResponse())
}
```

## 5. Conclusion: SOLID Principles in Action

This refactoring exercise demonstrates several SOLID principles:

-   **Single Responsibility Principle:** The `main` package handles user input, and the `api` package handles API communication. The `myJWTTransport` handles authentication.
-   **Open/Closed Principle:** The `api` package can be extended with new functionality without modifying the `main` package.
-   **Dependency Inversion Principle:** The `main` package depends on the `APIface` interface, not the concrete `API` struct, which makes the code more flexible and testable.

By following these principles, the instructor has created a more robust, maintainable, and reusable Go application.