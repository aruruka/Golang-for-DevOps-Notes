# I don't understand the benifit of using Generics.

The original code has Response interface been implemented by different struct types, i.e. Words, Occurrences, etc.
So, the returned variable would essentially be a Response type.
What is so different with using Generics? It looks harder to understand.

That's an excellent question. It gets to the very heart of why generics were introduced to Go. You are right that the interface approach works, but it has a hidden danger that the generics approach solves.

The key difference is **when** a type error is caught: **at runtime vs. at compile time.**

Let's break it down.

### The "Old Way" (with Interfaces)

```go
// Returns a general, abstract type: Response
func (a api) DoGetRequest(requestURL string) (Response, error) {
    // ...
}
```

1.  **The Promise:** This function promises to return *something* that satisfies the `Response` interface. It could be `Words`, or it could be `Occurrence`.
2.  **The Burden on the Caller:** When you call this function, you get back a `Response` variable. You, the programmer, believe you know what the *actual* concrete type is based on the URL you provided. To use it, you **must** perform a type assertion.

    ```go
    // We hope the response is of type Words
    resp, err := apiClient.DoGetRequest("http://example.com/words")
    if err != nil { /* ... */ }

    // This is a RUNTIME check.
    // What if the API returned an Occurrence by mistake? The program would PANIC.
    wordsResponse := resp.(Words)

    fmt.Println(wordsResponse.Words) // Only works if the type assertion was correct
    ```

The critical problem here is the `resp.(Words)` line. This is a check that happens **while your program is running**. If the API for some reason changes or returns an unexpected payload, your code will crash with a runtime panic. The compiler cannot help you prevent this bug.

---

### The "New Way" (with Generics)

```go
// Returns the EXACT concrete type T that you ask for
func DoGetRequest[T any](a api, requestURL string) (T, error) {
    // ...
}
```

1.  **The Promise:** This function makes a much stronger promise. It says, "You tell me what concrete type `T` you expect, and I will give you exactly that type back."
2.  **The Power of the Caller:** The caller specifies the exact type they want when they call the function.

    ```go
    // We are explicitly telling the compiler we want a `Words` struct back.
    // This is a COMPILE-TIME instruction.
    wordsResponse, err := DoGetRequest[Words](apiClient, "http://example.com/words")
    if err != nil { /* ... */ }

    // No type assertion needed!
    // The variable `wordsResponse` is already of type `Words`.
    // The compiler guarantees this.
    fmt.Println(wordsResponse.Words)
    ```

The `DoGetRequest[Words]` syntax is not just a hint; it's a command to the compiler. The compiler ensures that the `wordsResponse` variable is of type `Words`. You don't need to guess or assert the type.

### Summary of Benefits (Why Generics is Better Here)

| Feature | Interface Method | Generic Method |
| :--- | :--- | :--- |
| **Type Safety** | **Runtime check.** You can have a hidden bug that only appears when the program runs. | **Compile-time check.** The compiler guarantees the return type. Errors are caught before you even run the program. |
| **Code Clarity** | The caller has to perform a type assertion (`resp.(Words)`), which adds clutter and a point of failure. | The call itself (`DoGetRequest[Words]`) clearly states the intent. The resulting code is cleaner and safer. |
| **Refactoring** | If you add a new response type, you have to modify the `switch` statement inside `DoGetRequest`. | You don't have to modify `DoGetRequest` at all. It works for any future type automatically. |

While generics might look a little more complex at first glance, they create code that is significantly **safer, more explicit, and easier to maintain** by moving type checking from a risky runtime operation to a guaranteed compile-time one.