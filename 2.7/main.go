package main

import (
	"fmt"
	"math/rand"
	"time"
)

// value to channel
func asChan(vs ...int) <-chan int {
	c := make(chan int) //initial channel
	go func() {         //goroutine
		for _, v := range vs { //run for vs
			c <- v                                                        //value to channel
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) //random delay
		}
		close(c) //close channel
	}()
	return c //return value of chan
}

// merge channels
func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for { //infinity loop while ch open
			select { //look at where the elements are
			case v, ok := <-a: //if get value from ch 'a'
				if ok { //if have then we transmit to c
					c <- v
				} else {
					a = nil //empty(block this case)
				}
			case v, ok := <-b: //analogy
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	c := merge(a, b)
	for v := range c {
		fmt.Print(v)
	}
}

//the program outputs a random sequence from 1 to 8
