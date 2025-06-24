package api

import "fmt"

// RequestError represents an HTTP request error with context
type RequestError struct {
	HTTPCode int
	Body     string
	Err      string
}

// Error implements the error interface for RequestError
func (r RequestError) Error() string {
	return fmt.Sprintf("HTTP %d: %s - %s", r.HTTPCode, r.Err, r.Body)
}
