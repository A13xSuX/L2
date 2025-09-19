package main

import "fmt"

func main() {
	a := [5]int{76, 77, 78, 79, 80} //massive size 5
	var b []int = a[1:4]            //create slice [77,78,79] len 3 cap 4, reffers to basic massive
	fmt.Println(b)
}
