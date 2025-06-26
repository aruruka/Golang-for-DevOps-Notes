# Assignment 2: Rate Limiting

## Overview
In this assignment, you'll learn how to implement rate limiting in Go to control the number of requests made to an API endpoint. You will need to hit an endpoint at a specific rate to get a successful response.

## Requirements

### Task
Hit the `http://localhost:8080/ratelimit` endpoint at a constant rate of 5 requests per second.

### Server Behavior
- The server will respond with your current request rate.
- If you successfully hit the endpoint at 5 requests per second, the response will be prefixed with "DONE!".
- If you exceed 5 requests per second, the server will return a `429 Too Many Requests` HTTP error for the next 10 seconds.

## Success Criteria
- Successfully make HTTP GET requests to the endpoint.
- Achieve a constant request rate of 5 requests per second.
- Receive a "DONE!" message from the server.
- Handle the `429` error code gracefully if the rate limit is exceeded.

## Implementation Notes
- Use goroutines and channels to manage concurrent requests.
- Use `time.Ticker` to control the rate of requests.

## Setup
1. Start the test server (located in the `test-server` directory). Run command - `cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/test-server && go run *.go`
2. Ensure the server is running on `http://localhost:8080`.
3. Your Go program should make requests to `http://localhost:8080/ratelimit`.

## Example Output
```plaintext
Hitting API at 3 requests in a given second
Hitting API at 3 requests in a given second
Hitting API at 3 requests in a given second
DONE! You did it! Hitting API at 5 requests in a given second
DONE! You did it! Hitting API at 5 requests in a given second
DONE! You did it! Hitting API at 5 requests in a given second
```

## Validation
You'll know your implementation works when:
- You can consistently hit the endpoint at 5 requests per second.
- The server responds with the "DONE!" prefix.
- Your program doesn't crash and handles potential `429` errors.
