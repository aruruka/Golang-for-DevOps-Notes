package api

import (
	"net/http"
)

// Options contains configuration for the API client
type Options struct {
	BaseURL string
}

// ClientIface defines the interface for HTTP client operations
type ClientIface interface {
	Get(url string) (resp *http.Response, err error)
}

// APIIface defines the interface for our API operations
type APIIface interface {
	GetAssignmentData(endpoint string) (Response, error)
}

// Response interface for different response types
type Response interface {
	GetResponse() string
}

// api struct implements APIIface
type api struct {
	Options Options
	Client  ClientIface
}

// New creates a new API client instance
func New(options Options) APIIface {
	return api{
		Options: options,
		Client:  &http.Client{},
	}
}
