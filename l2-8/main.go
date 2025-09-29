package main

import (
	"fmt"
	"os"

	"l2-8/ntp"
)

func main() {
	time, err := ntp.GetCurrentTime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(time)
}
