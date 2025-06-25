# Golang: Arrays and Slices

This document provides a comprehensive reference for arrays and slices in Go, based on the course material.

## 1. The Core Difference: Fixed vs. Dynamic

The fundamental difference between arrays and slices is their length:

*   **Arrays:** Have a fixed length that is defined at compile time. You cannot change the size of an array after it's created.
*   **Slices:** Are dynamic and can grow or shrink in size. They provide a more flexible way to work with sequences of data.

## 2. Arrays: The Foundation

Arrays are the building blocks of slices. While you'll often work with slices directly, it's crucial to understand arrays.

### 2.1. Declaration and Initialization

You declare an array by specifying its size and the type of its elements.

```go
// An array of 7 integers
var buffer [7]int
```

You can also initialize an array with values using an array literal:

```go
var arr1 [7]int = [7]int{7, 3, 6, 0, 4, 9, 10}
```

### 2.2. Length and Capacity

For an array, the length and capacity are always the same. The capacity is the maximum number of elements the array can hold, which is its fixed size.

```go
fmt.Printf("Length: %d, Capacity: %d\n", len(arr1), cap(arr1))
// Output: Length: 7, Capacity: 7
```

## 3. Slices: Flexible Views of Arrays

Slices are a more common and powerful tool in Go. A slice is a descriptor of a contiguous segment of an underlying array.

### 3.1. Creating a Slice

You can create a slice from an existing array by specifying a half-open range with `[low:high]`. This creates a new slice that "points" to the data in the original array.

```go
var arr1 [7]int = [7]int{7, 3, 6, 0, 4, 9, 10}
var arr2 []int = arr1[1:3] // Creates a slice from index 1 up to (but not including) index 3
```

In the example above, `arr2` will contain the elements `[3 6]`.

### 3.2. Length and Capacity of a Slice

A slice has three components:

1.  **Pointer:** A pointer to the first element of the underlying array that the slice can access.
2.  **Length:** The number of elements in the slice.
3.  **Capacity:** The number of elements in the underlying array from the start of the slice to the end of the array.

Let's look at our example:

```go
fmt.Printf("Length: %d, Capacity: %d\n", len(arr2), cap(arr2))
// Output: Length: 2, Capacity: 6
```

*   **Length is 2:** The slice `arr2` contains two elements (`3` and `6`).
*   **Capacity is 6:** The slice starts at index 1 of `arr1`. From that point, there are 6 elements remaining in the underlying array (`3, 6, 0, 4, 9, 10`).

### 3.3. Slices are References

A crucial concept to understand is that a slice is a *reference* to the underlying array. **It does not create a copy of the data.** This means that if you modify the elements of a slice, you are also modifying the elements of the original array.

Here's a demonstration:

```go
// arr2 is a slice of arr1
for k := range arr2 {
    arr2[k] += 1 // Increment each element in the slice
}

fmt.Println(arr1)
// Output: [7 4 7 0 4 9 10]
```

As you can see, the original array `arr1` has been modified.

## 4. Slice Operations

### 4.1. Creating a Slice with `make`

You can create a slice with a specified length and capacity using the built-in `make` function.

```go
// Create a slice of integers with length 3 and capacity 9
arr4 := make([]int, 3, 9)
```

This creates an underlying array of size 9 and a slice that refers to the first 3 elements. The initial values will be the zero value for the type (e.g., `0` for `int`).

### 4.2. Appending to a Slice

The `append` function is used to add elements to a slice.

```go
var arr3 []int = []int{1, 2, 3}
arr3 = append(arr3, 4)
```

`append` is a smart function:

*   **If there is enough capacity** in the underlying array, `append` will reuse the existing array and simply increase the length of the slice.
*   **If there is not enough capacity**, `append` will allocate a new, larger underlying array, copy the existing elements, and then add the new element.

Because `append` might return a new slice (with a new underlying array), you should always assign the result of `append` back to the slice variable: `mySlice = append(mySlice, newValue)`.

## 5. Code Examples from the Course

The following code from `slices-demo/cmd/array-and-slice/main.go` demonstrates the concepts discussed above.

```go
package main

import "fmt"

func main() {
	// 1. Create an array
	var arr1 [7]int = [7]int{7, 3, 6, 0, 4, 9, 10}
	fmt.Println(arr1)
	fmt.Printf("%d %d\n", len(arr1), cap(arr1))

	// 2. Create a slice from the array
	var arr2 []int = arr1[1:3]
	fmt.Println(arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))

	// 3. Extend the slice (possible because of available capacity)
	arr2 = arr2[0 : len(arr2)+2]
	fmt.Println(arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))

	// 4. Modify the slice and see the effect on the original array
	for k := range arr2 {
		arr2[k] += 1
	}
	fmt.Println(arr2)
	fmt.Printf("%d %d\n", len(arr2), cap(arr2))
	fmt.Println(arr1)

	// 5. Create a slice and append to it
	var arr3 []int = []int{1, 2, 3}
	fmt.Println(arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))
	arr3 = append(arr3, 4)
	fmt.Println(arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))
	arr3 = append(arr3, 5)
	fmt.Println(arr3)
	fmt.Printf("%d %d\n", len(arr3), cap(arr3))

	// 6. Create a slice with make
	arr4 := make([]int, 3, 9)
	fmt.Println(arr4)
	fmt.Printf("%d %d\n", len(arr4), cap(arr4))
}
```

### Why the Instructor Coded It This Way

The instructor's code is designed to be a clear and concise demonstration of the key properties of arrays and slices.

*   **Step-by-step progression:** The code builds from a simple array to more complex slice operations, making it easy to follow the logic.
*   **`fmt.Printf` for clarity:** Using `fmt.Printf` to print the length and capacity at each step provides immediate feedback and reinforces the concepts being taught.
*   **Demonstrating the reference nature:** The loop that modifies `arr2` and then prints `arr1` is a powerful and effective way to show that slices are references to underlying arrays.
*   **Illustrating `append`:** The repeated use of `append` on `arr3` shows how the length and capacity change as elements are added.
*   **Introducing `make`:** The final example with `make` introduces the standard way to create slices with a predefined capacity, which is a common practice for performance optimization.
