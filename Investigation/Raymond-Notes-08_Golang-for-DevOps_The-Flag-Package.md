# Go Flag Package Reference for DevOps and Cloud Engineers

## Table of Contents
1. [Overview](#overview)
2. [Core Concepts](#core-concepts)
3. [Basic Flag Types](#basic-flag-types)
4. [Flag Definition Methods](#flag-definition-methods)
5. [Parsing and Validation](#parsing-and-validation)
6. [Best Practices](#best-practices)
7. [Real-World Examples](#real-world-examples)
8. [Common Patterns in DevOps](#common-patterns-in-devops)

## Overview

The Go `flag` package is a built-in library for parsing command-line arguments in Go applications. It provides a clean, standardized way to handle command-line flags, making it essential for DevOps tools, CLI utilities, and cloud engineering applications.

### Key Benefits
- **Type Safety**: Automatic type conversion and validation
- **Self-Documenting**: Built-in help generation
- **POSIX Compliance**: Follows standard Unix flag conventions
- **Minimal Dependencies**: Part of Go's standard library

## Core Concepts

### Flag vs Argument Terminology

It's important to distinguish between different command-line elements:

- **Flag**: A named parameter that modifies program behavior (e.g., `-verbose`, `--port=8080`)
- **Argument**: Positional parameters passed to the program (e.g., file names)
- **Option**: Another term for flag, commonly used interchangeably
- **Parameter**: The value associated with a flag

### Flag Syntax Patterns

Go's flag package supports several syntax patterns:
- Short flags: `-v` (single dash, single character)
- Long flags: `--verbose` (double dash, descriptive name)
- Value assignment: `-port=8080` or `-port 8080`
- Boolean flags: `-debug` (presence indicates true)

## Basic Flag Types

The flag package provides built-in support for common data types:

### String Flags
```go
var name = flag.String("name", "default", "Help text for name flag")
var config string
flag.StringVar(&config, "config", "/etc/app.conf", "Configuration file path")
```

### Integer Flags
```go
var port = flag.Int("port", 8080, "Server port number")
var timeout int
flag.IntVar(&timeout, "timeout", 30, "Request timeout in seconds")
```

### Boolean Flags
```go
var verbose = flag.Bool("verbose", false, "Enable verbose logging")
var debug bool
flag.BoolVar(&debug, "debug", false, "Enable debug mode")
```

### Duration Flags
```go
var interval = flag.Duration("interval", 5*time.Second, "Polling interval")
```

## Flag Definition Methods

### Method 1: Using Flag Functions (Pointer Return)
```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // Define flags that return pointers
    name := flag.String("name", "world", "Name to greet")
    count := flag.Int("count", 1, "Number of greetings")
    verbose := flag.Bool("verbose", false, "Enable verbose output")
    
    flag.Parse()
    
    // Dereference pointers to get values
    fmt.Printf("Hello %s!\n", *name)
    if *verbose {
        fmt.Printf("Greeting count: %d\n", *count)
    }
}
```

**Why this approach?** The instructor uses pointer returns because it allows the flag package to manage memory allocation and provides a clean separation between flag definition and usage.

### Method 2: Using Var Functions (Reference Assignment)
```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    var (
        name    string
        count   int
        verbose bool
    )
    
    // Define flags using existing variables
    flag.StringVar(&name, "name", "world", "Name to greet")
    flag.IntVar(&count, "count", 1, "Number of greetings")
    flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")
    
    flag.Parse()
    
    // Use variables directly (no dereferencing needed)
    fmt.Printf("Hello %s!\n", name)
    if verbose {
        fmt.Printf("Greeting count: %d\n", count)
    }
}
```

**Why this approach?** This method is preferred when you want to use the variables throughout your application without constant dereferencing, leading to cleaner code.

## Parsing and Validation

### The Parse() Function
```go
flag.Parse()
```

The `Parse()` function must be called after all flags are defined but before accessing their values. This function:
- Processes `os.Args[1:]` (command-line arguments excluding program name)
- Validates flag syntax and types
- Assigns values to flag variables
- Generates help text when `-h` or `--help` is used

### Accessing Non-Flag Arguments
```go
flag.Parse()
args := flag.Args()  // Returns []string of non-flag arguments
nargs := flag.NArg() // Returns count of non-flag arguments
```

### Error Handling
```go
import (
    "flag"
    "fmt"
    "os"
)

func main() {
    port := flag.Int("port", 8080, "Server port")
    flag.Parse()
    
    // Validate parsed values
    if *port < 1 || *port > 65535 {
        fmt.Fprintf(os.Stderr, "Error: Port must be between 1 and 65535\n")
        os.Exit(1)
    }
}
```

## Best Practices

### 1. Consistent Naming Conventions
```go
// Good: Use kebab-case for multi-word flags
flag.String("log-level", "info", "Logging level")
flag.String("config-file", "app.yaml", "Configuration file")

// Avoid: Inconsistent naming
flag.String("logLevel", "info", "Logging level")  // camelCase
flag.String("config_file", "app.yaml", "Configuration file")  // snake_case
```

### 2. Meaningful Default Values
```go
// Good: Sensible defaults that work in most environments
flag.Int("port", 8080, "HTTP server port")
flag.String("log-level", "info", "Log level (debug, info, warn, error)")
flag.Duration("timeout", 30*time.Second, "Request timeout")

// Poor: Defaults that require user intervention
flag.String("database-url", "", "Database connection string (required)")
```

### 3. Clear, Descriptive Help Text
```go
// Good: Clear, actionable help text
flag.String("config", "/etc/app/config.yaml", 
    "Path to YAML configuration file. Must be readable by the application.")

// Poor: Vague or unhelpful descriptions
flag.String("config", "", "config file")
```

### 4. Validation After Parsing
```go
func validateFlags() error {
    flag.Parse()
    
    if *port < 1 || *port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535, got %d", *port)
    }
    
    if *logLevel != "debug" && *logLevel != "info" && *logLevel != "warn" && *logLevel != "error" {
        return fmt.Errorf("invalid log level: %s", *logLevel)
    }
    
    return nil
}
```

## Real-World Examples

### DevOps Tool Example: Log Analyzer
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "time"
)

type Config struct {
    LogFile     string
    OutputDir   string
    StartTime   time.Time
    EndTime     time.Time
    Verbose     bool
    Format      string
    Workers     int
}

func main() {
    var config Config
    
    // Define flags with meaningful defaults
    flag.StringVar(&config.LogFile, "log-file", "", 
        "Path to log file to analyze (required)")
    flag.StringVar(&config.OutputDir, "output-dir", "./reports", 
        "Directory to write analysis reports")
    flag.BoolVar(&config.Verbose, "verbose", false, 
        "Enable detailed progress output")
    flag.StringVar(&config.Format, "format", "json", 
        "Output format: json, csv, or text")
    flag.IntVar(&config.Workers, "workers", 4, 
        "Number of concurrent workers for processing")
    
    // Custom flag type for time parsing
    var startTimeStr = flag.String("start-time", "", 
        "Start time for analysis (RFC3339 format)")
    var endTimeStr = flag.String("end-time", "", 
        "End time for analysis (RFC3339 format)")
    
    flag.Parse()
    
    // Validate required flags
    if config.LogFile == "" {
        fmt.Fprintf(os.Stderr, "Error: -log-file is required\n")
        flag.Usage()
        os.Exit(1)
    }
    
    // Parse time flags with validation
    if *startTimeStr != "" {
        var err error
        config.StartTime, err = time.Parse(time.RFC3339, *startTimeStr)
        if err != nil {
            log.Fatalf("Invalid start-time format: %v", err)
        }
    }
    
    if *endTimeStr != "" {
        var err error
        config.EndTime, err = time.Parse(time.RFC3339, *endTimeStr)
        if err != nil {
            log.Fatalf("Invalid end-time format: %v", err)
        }
    }
    
    // Validate format option
    validFormats := map[string]bool{"json": true, "csv": true, "text": true}
    if !validFormats[config.Format] {
        log.Fatalf("Invalid format: %s. Must be json, csv, or text", config.Format)
    }
    
    if config.Verbose {
        fmt.Printf("Starting log analysis with config: %+v\n", config)
    }
    
    // Process remaining arguments as additional log files
    additionalFiles := flag.Args()
    if len(additionalFiles) > 0 && config.Verbose {
        fmt.Printf("Additional log files: %v\n", additionalFiles)
    }
}
```

**Why this structure?** The instructor organizes flags into a configuration struct to:
- Group related configuration options
- Make the code more maintainable
- Enable easy serialization of configuration
- Simplify testing by allowing config injection

### Cloud Service Configuration Example
```go
package main

import (
    "flag"
    "fmt"
    "log"
    "time"
)

type CloudConfig struct {
    Region          string
    InstanceType    string
    MinInstances    int
    MaxInstances    int
    HealthCheckURL  string
    HealthInterval  time.Duration
    AutoScale       bool
    DryRun          bool
}

func main() {
    config := &CloudConfig{}
    
    // Infrastructure flags
    flag.StringVar(&config.Region, "region", "us-west-2", 
        "AWS region for deployment")
    flag.StringVar(&config.InstanceType, "instance-type", "t3.micro", 
        "EC2 instance type")
    
    // Scaling configuration
    flag.IntVar(&config.MinInstances, "min-instances", 1, 
        "Minimum number of instances")
    flag.IntVar(&config.MaxInstances, "max-instances", 10, 
        "Maximum number of instances")
    flag.BoolVar(&config.AutoScale, "auto-scale", true, 
        "Enable automatic scaling based on metrics")
    
    // Health check configuration
    flag.StringVar(&config.HealthCheckURL, "health-check-url", "/health", 
        "HTTP endpoint for health checks")
    flag.DurationVar(&config.HealthInterval, "health-interval", 30*time.Second, 
        "Interval between health checks")
    
    // Operational flags
    flag.BoolVar(&config.DryRun, "dry-run", false, 
        "Show what would be done without making changes")
    
    flag.Parse()
    
    // Validation logic
    if err := validateCloudConfig(config); err != nil {
        log.Fatalf("Configuration error: %v", err)
    }
    
    if config.DryRun {
        fmt.Printf("DRY RUN: Would deploy with config: %+v\n", config)
        return
    }
    
    fmt.Printf("Deploying with configuration: %+v\n", config)
}

func validateCloudConfig(config *CloudConfig) error {
    if config.MinInstances < 1 {
        return fmt.Errorf("min-instances must be at least 1")
    }
    
    if config.MaxInstances < config.MinInstances {
        return fmt.Errorf("max-instances (%d) must be >= min-instances (%d)", 
            config.MaxInstances, config.MinInstances)
    }
    
    if config.HealthInterval < time.Second {
        return fmt.Errorf("health-interval must be at least 1 second")
    }
    
    return nil
}
```

## Common Patterns in DevOps

### 1. Configuration Hierarchy Pattern
```go
// Priority: CLI flags > Environment variables > Config file > Defaults
func loadConfig() *Config {
    config := &Config{
        // Set defaults
        Port:     8080,
        LogLevel: "info",
    }
    
    // Load from config file if specified
    if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
        loadFromFile(config, configFile)
    }
    
    // Override with environment variables
    if port := os.Getenv("PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Port = p
        }
    }
    
    // CLI flags have highest priority
    flag.IntVar(&config.Port, "port", config.Port, "Server port")
    flag.StringVar(&config.LogLevel, "log-level", config.LogLevel, "Log level")
    flag.Parse()
    
    return config
}
```

### 2. Subcommand Pattern
```go
package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "Commands: deploy, scale, status\n")
        os.Exit(1)
    }
    
    switch os.Args[1] {
    case "deploy":
        deployCmd()
    case "scale":
        scaleCmd()
    case "status":
        statusCmd()
    default:
        fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func deployCmd() {
    deployFlags := flag.NewFlagSet("deploy", flag.ExitOnError)
    region := deployFlags.String("region", "us-west-2", "Deployment region")
    instances := deployFlags.Int("instances", 3, "Number of instances")
    
    deployFlags.Parse(os.Args[2:])
    
    fmt.Printf("Deploying %d instances to %s\n", *instances, *region)
}

func scaleCmd() {
    scaleFlags := flag.NewFlagSet("scale", flag.ExitOnError)
    replicas := scaleFlags.Int("replicas", 1, "Target replica count")
    service := scaleFlags.String("service", "", "Service name to scale")
    
    scaleFlags.Parse(os.Args[2:])
    
    if *service == "" {
        fmt.Fprintf(os.Stderr, "Error: -service is required\n")
        scaleFlags.Usage()
        os.Exit(1)
    }
    
    fmt.Printf("Scaling service %s to %d replicas\n", *service, *replicas)
}

func statusCmd() {
    statusFlags := flag.NewFlagSet("status", flag.ExitOnError)
    verbose := statusFlags.Bool("verbose", false, "Show detailed status")
    
    statusFlags.Parse(os.Args[2:])
    
    fmt.Printf("Service status (verbose: %t)\n", *verbose)
}
```

**Why subcommands?** The instructor demonstrates this pattern because modern DevOps tools often need multiple operations (deploy, scale, monitor) with different flag sets for each operation.

### 3. Environment Integration Pattern
```go
func envOrFlag(envVar, flagValue, defaultValue string) string {
    if flagValue != defaultValue {
        return flagValue  // Flag was explicitly set
    }
    if env := os.Getenv(envVar); env != "" {
        return env  // Use environment variable
    }
    return defaultValue  // Use default
}

func main() {
    apiKey := flag.String("api-key", "", "API key for authentication")
    region := flag.String("region", "us-west-2", "AWS region")
    
    flag.Parse()
    
    // Integrate with environment variables
    finalAPIKey := envOrFlag("API_KEY", *apiKey, "")
    finalRegion := envOrFlag("AWS_REGION", *region, "us-west-2")
    
    if finalAPIKey == "" {
        log.Fatal("API key required via -api-key flag or API_KEY environment variable")
    }
    
    fmt.Printf("Using region: %s\n", finalRegion)
}
```

This comprehensive reference covers the essential aspects of Go's flag package for DevOps and cloud engineering applications, with practical examples and best practices that reflect real-world usage patterns.