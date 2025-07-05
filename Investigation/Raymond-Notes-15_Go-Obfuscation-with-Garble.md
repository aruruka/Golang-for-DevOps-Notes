# Go Stealth Technique: Obfuscating Your Code with Garble

*Original Author: LingFei Programming, July 5, 2025*

---

In project development, we often need to provide an external tool or service, but we don't want the recipient to reverse-engineer the source code from the provided files. This is where code obfuscation becomes necessary. This article explains how to use Garble to obfuscate your Go code.

Garble is an open-source tool developed by the `burrowers` community. It wraps the Go compiler to provide a one-stop solution for generating highly obfuscated Go binaries. It significantly increases the difficulty of source code restoration and reverse engineering while maintaining binary compatibility.

### Key Features of Garble:

*   **Identifier/Package Path Obfuscation**: Renames functions, variables, and struct names, and removes most metadata.
*   **String Literal Encryption**: The `-literals` flag ensures that strings are only decrypted at runtime.
*   **Minimal File Size**: The `-tiny` flag removes debugging symbols, filenames, and line numbers, making it harder to attack.
*   **Reproducible Builds**: The `-seed` flag uses a fixed seed to ensure that the same obfuscation result can be reproduced.
*   **Stack Trace De-obfuscation**: Use `garble reverse` with a known seed to restore original symbols from a stack trace.

---

## Installing Garble

```bash
go install github.com/burrowers/garble@latest
```

---

## Obfuscating a Simple Program

Here is an example `main.go` file:

```go
// file: main.go
package main

import "fmt"

func main() {
    secret := "Hello, Obfuscation!"
    fmt.Println(process(secret))
}

func process(s string) string {
    return s + "-processed"
}
```

### Normal Build

```bash
# Build the application normally
go build -o normal_app main.go

# Search for the function name in the binary
strings normal_app | grep process
# >> process
```

### Build with Garble

```bash
# Build with Garble
garble build -o garbled_app main.go

# Search for the function name in the obfuscated binary
strings garbled_app | grep process
# >> no "process" found
```

---

## Literal Encryption (Making Strings Invisible)

To encrypt every string literal in your binary:

```bash
# Build with the -literals flag
garble -literals build -o garbled_lit main.go

# Search for a string literal in the binary
strings garbled_lit | grep Hello
# >> (nothing - strings scrambled at runtime)
```

---

## Deterministic Builds and Reverse Engineering Support

### 1. Deterministic Obfuscation

Modify the `main.go` file to include a panic, which helps in demonstrating stack trace reversal.

```go
package main

import "fmt"

func main() {
    secret := "Hello, Obfuscation!"
    fmt.Println(process(secret))
    panic("panic me")
}

func process(s string) string {
    return s + "-processed"
}
```

Use a fixed seed to produce a unique binary, which is useful for bug reproduction and debugging.

```bash
# Build with a random seed, which Garble will output
garble -seed=random build -o deterministic_app main.go
# >> -seed chosen at random: 75MYDgjSJGFJT7ktVURDYW
```

### 2. Restoring Stack Symbols

When the program crashes, developers can reverse-engineer the stack trace symbols.

```bash
# Redirect the panic output to a file
./deterministic_app &> panic-output.txt

# Use the same seed to reverse the stack trace
garble -seed=75MYDgjSJGFJT7ktVURDYW reverse main.go panic-output.txt
```

The output will show the original, de-obfuscated stack trace:

```
Hello, Obfuscation!-processed
panic: panic me

goroutine 1 [running]:
main.main()
    command-line-arguments/main.go:8 +0x7c
```

---

## Important Considerations and Experimental Features

*   **Exported Symbols**: Exported symbols (used for reflection/interfaces) are not obfuscated.
*   **Go Plugins**: Go plugins are not currently supported.
*   **Control Flow Obfuscation**: This is an experimental feature that can be enabled with an environment variable: `GARBLE_EXPERIMENTAL_CONTROLFLOW=1 garble build`.
*   **Runtime Strings**: Although source code information is cleaned, some Go runtime strings may still be visible.

---

## Why Recommend Garble?

*   **Greatly Increases Difficulty of Reverse Engineering**: Makes it hard to restore function names, algorithms, and business logic.
*   **Fully Compatible**: Works seamlessly with Go modules, caching, stack traces, and automated builds.
*   **High Performance**: Only 1-2 times slower than a standard `go build`.

---

## Practical Recommendations

1.  **CI Integration**: Set up automated or manual builds for both normal and obfuscated versions.
2.  **`-tiny` Mode**: Use this to output extremely small executable files.
3.  **Enhanced Security**: Combine with `-ldflags="-s -w"` and `-trimpath` to strip symbol tables and absolute paths.
4.  **High-Security Needs**: Enable experimental control flow obfuscation for maximum protection.

---

## Conclusion

Garble significantly increases the cost of decompilation and restoration. However, "obfuscation is not absolute security." Tools like `GoStringUngarbler` exist to counter string literal obfuscation, and runtime debugging can bypass some protections. A determined adversary will always find a way, but Garble is a crucial layer of security, suitable as a "final line of defense" in the development and release process.
