# Some notes for running tests in this directory

1. When running the "api_client.go".

1.1 Start the test server first:

```bash
cd test-server
go run *.go # This starts the simple api server
```

1.2 Then run the test client:

```bash
go run main.go api_client.go 'http://localhost:8080/words?input=word1' # Make sure `httpGet()` is called in the main function
```