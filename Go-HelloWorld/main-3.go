package main

import "fmt"

type myStruct struct {
	string
}

func main3() {
	a := myStruct{"a string"}

	fmt.Println(a.string)
}
