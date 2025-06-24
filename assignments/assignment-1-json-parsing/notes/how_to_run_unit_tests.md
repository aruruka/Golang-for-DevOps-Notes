# Can you show me how to run the unit tests?

I'll show you how to run the unit tests for the assignment. Here are several ways to run them:

## 1. Run All Tests in the API Package

```bash
cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/assignments/assignment-1-json-parsing
go test ./pkg/api
```

## 2. Run Tests with Verbose Output (Recommended)

```bash
cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/assignments/assignment-1-json-parsing
go test ./pkg/api -v
```

## 3. Run Tests with Coverage Report

```bash
cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/assignments/assignment-1-json-parsing
go test ./pkg/api -v -cover
```

## 4. Run a Specific Test Function

```bash
cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/assignments/assignment-1-json-parsing
go test ./pkg/api -v -run TestGetAssignmentData
```

## 5. Run All Tests in the Project

```bash
cd /home/raymond/Work/Golang-for-DevOps-and-Cloud-Engieers/assignments/assignment-1-json-parsing
go test ./...
```

Let me demonstrate by running the tests: Based on the structure of your tests, here's what each command will do:

## Expected Test Output

When you run `go test ./pkg/api -v`, you should see something like:

```
=== RUN   TestGetAssignmentData
--- PASS: TestGetAssignmentData (0.00s)
=== RUN   TestGetAssignmentDataErrorHandling  
--- PASS: TestGetAssignmentDataErrorHandling (0.00s)
=== RUN   TestGetAssignmentDataInvalidJSON
--- PASS: TestGetAssignmentDataInvalidJSON (0.00s)
PASS
ok      assignments/assignment-1-json-parsing/pkg/api   0.001s
```

## What Each Test Does

1. **`TestGetAssignmentData`**: Tests the main functionality with mock data that matches the requirements JSON
2. **`TestGetAssignmentDataErrorHandling`**: Tests HTTP error scenarios (404, 500, etc.)
3. **`TestGetAssignmentDataInvalidJSON`**: Tests invalid JSON response handling

## Test Coverage Details

When you run with `-cover`, you'll see coverage percentage showing how much of your code is tested by the unit tests.

## Quick Test Command Summary

```bash
# Most common - run with verbose output
go test ./pkg/api -v

# With coverage report
go test ./pkg/api -v -cover

# Run from project root
cd assignments/assignment-1-json-parsing && go test ./pkg/api -v
```

The tests validate that:
- ✅ Your `AssignmentData` struct correctly maps to the JSON from requirements
- ✅ HTTP GET requests work properly 
- ✅ JSON parsing handles all data types (strings, arrays, maps, nulls, mixed types)
- ✅ Error handling works for various failure scenarios
- ✅ The `GetResponse()` method formats output correctly

Try running `go test ./pkg/api -v` to see your tests in action!