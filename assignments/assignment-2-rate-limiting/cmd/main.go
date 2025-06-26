package main

import (
	"assignment-2-rate-limiting/pkg/ratelimiter"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rl := ratelimiter.NewRateLimiter(5)

	go rl.Start()

	// Wait for a signal to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down...")
	rl.Stop()
}
