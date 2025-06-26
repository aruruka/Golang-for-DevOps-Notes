package ratelimiter

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter controls the rate of HTTP requests.
type RateLimiter struct {
	Client      *http.Client
	Rate        int
	StopChannel chan bool
	stopOnce    sync.Once
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(rate int) *RateLimiter {
	return &RateLimiter{
		Client:      &http.Client{},
		Rate:        rate,
		StopChannel: make(chan bool),
	}
}

// Start sends requests at a specified rate.
func (rl *RateLimiter) Start() {
	ticker := time.NewTicker(time.Second / time.Duration(rl.Rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			req, err := http.NewRequest("GET", "http://localhost:8080/ratelimit", nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				continue
			}
			rl.MakeRequest(req)
		case <-rl.StopChannel:
			return
		}
	}
}

// MakeRequest sends an HTTP request and handles the response.
func (rl *RateLimiter) MakeRequest(req *http.Request) {
	resp, err := rl.Client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Println(string(body))
		if strings.HasPrefix(string(body), "DONE!") {
			rl.Stop()
		}
	case http.StatusTooManyRequests:
		fmt.Println("Rate limit exceeded. Backing off...")
		time.Sleep(10 * time.Second)
	default:
		fmt.Printf("Received status code %d: %s\n", resp.StatusCode, string(body))
	}
}

// Stop stops the rate limiter.
func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.StopChannel)
	})
}
