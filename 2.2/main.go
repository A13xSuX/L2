package main

import "fmt"

func test() (x int) { //named return parametr
	defer func() { //defer function has access to this parameter.
		x++
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() { //defer changed local parametr, not return value
		x++
	}()
	x = 1    //saved for return
	return x //defer does not affect the return value
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
