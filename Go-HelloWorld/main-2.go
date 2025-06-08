package main

import (
	"fmt"
	"os"
)

func main2() {
	// os.Args[0] is the program name
	// os.Args[1:] contains the actual arguments
	fmt.Println("Program name:", os.Args[0])
	if len(os.Args) > 1 {
		fmt.Println("Arguments:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("  %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("No arguments provided")
	}
}
