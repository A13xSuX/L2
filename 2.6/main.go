package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"} //slice [1,2,3] len 3 cap 3
	modifySlice(s)
	fmt.Println(s) //slice [3,2,3] len 3 cap 3
}

func modifySlice(i []string) { //slice is passed by reference to the structure
	i[0] = "3"         //[3,2,3] len 3 cap 3, two slices refers on 1 basic massive
	i = append(i, "4") // [3,2,3,4] len 4 cap 6 new slice,because of the realocation
	i[1] = "5"         //[3,5,3,4] len 4 cap 6
	i = append(i, "6") //[3,5,3,4,6] len 5 cap 6
}
