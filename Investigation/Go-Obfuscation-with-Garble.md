# Go Stealth: Obfuscating Your Code with Garble

**Original:** Kong Lingfei, Lingfei Programming  
**Date:** July 5, 2025, 13:30

---

In project development, we often need to provide a tool or service to external parties, but we don't want them to reverse-engineer the source code from the provided files. This is where code obfuscation becomes necessary. This article introduces how to use Garble to obfuscate your Go code.

[Garble](https://github.com/burrowers/garble) is an open-source tool developed by the burrowers community. It wraps the Go compiler to provide a one-stop solution for generating highly obfuscated Go binaries. It aims to maximize the difficulty of source code restoration and reverse engineering while maintaining binary compatibility as much as possible.

### Key Features of Garble:

*   **Identifier/Package Path Obfuscation:** Renames functions, variables, and struct names, and removes most metadata.
*   **String Literal Encryption:** The `-literals` flag ensures that each string is only decrypted at runtime.
*   **Minimal File Size:** The `-tiny` flag removes debugging symbols, filenames, and line numbers, increasing the difficulty for attackers.
*   **Reproducible Builds:** The `-seed` flag fixes the seed to ensure that the same obfuscation result can be reproduced.
*   **Stack Trace De-obfuscation:** Use `garble reverse` with a known seed to restore obfuscated stack symbols.

---

### Installing Garble

```shell
$ go install github.com/burrowers/garble@latest
```

---

### Obfuscating a Simple Program

Here is an example code snippet:

```go
// File: main.go
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

#### Normal Build:

```shell
$ go build -o normal_app main.go
$ strings normal_app | grep process
# >> process
```

#### Now, using Garble for obfuscation:

```shell
$ garble build -o garbled_app main.go
$ strings garbled_app | grep process
# >> no "process" found
```

---

### Literal Encryption (Strings Not Visible)

Encrypt every string literal:

```shell
$ garble -literals build -o garbled_lit main.go
$ strings garbled_lit | grep Hello
# >> (nothing - strings scrambled at runtime)
```

---

### Deterministic Builds and Reverse Engineering Support

#### 1. Deterministic Obfuscation

Modify the `main.go` file as follows:

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

Use a fixed seed to get a unique binary (convenient for bug reproduction and location):

```shell
$ garble -seed=random build -o deterministic_app main.go
# -seed chosen at random: 75MYDgjSJGFJT7ktvUROYW
```

#### 2. Restoring Stack Symbols

When a program crashes and needs investigation, developers can reverse-parse the stack symbols:

```shell
$ ./deterministic_app &> panic-output.txt
$ garble -seed=75MYDgjSJGFJT7ktvUROYW reverse main.go panic-output.txt
Hello, Obfuscation!-processed
panic: panic me

goroutine 1 [running]:
main.main()
    command-line-arguments/main.go:8 +0x7c
```

---

### Caveats and Experimental Features

*   Exported symbols (used for reflection/interfaces) will not be obfuscated.
*   Go plugins are not currently supported.
*   Control flow obfuscation can be enabled with an experimental variable: `GARBLE_EXPERIMENTAL_CONTROLFLOW=1 garble build ...`
*   Source code information is cleaned, but some Go runtime strings may still be visible.

---

### Why Recommend Garble?

*   **Greatly increases reverse engineering and analysis difficulty** (hard to restore function names/algorithms/business logic).
*   **Perfectly compatible** with Go modules, caching, stack traces, and automated builds.
*   **High performance:** Only 1-2 times slower than a standard `go build`.

---

### Practical Suggestions

1.  **CI Integration:** Automatically/manually build both normal and obfuscated versions.
2.  **`-tiny` mode:** Output extremely small executable files.
3.  **Security Enhancement:** Combine with `-ldflags="-s -w"` and `-trimpath` to clean up the symbol table and absolute paths.
4.  **High-Security Needs:** Experimentally enable control flow obfuscation.

---

### Conclusion

Garble significantly increases the cost of decompilation and restoration. However, "obfuscation ≠ absolute security." Tools like `GoStringUngarbler` exist to counter literal obfuscation, and runtime debugging can bypass it. A determined adversary will always find a way, but Garble is a vital link in enhancing security, suitable as the "last line of defense" in the development and release process.

---

### References

[1] Garble: https://github.com/burrowers/garble
