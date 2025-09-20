package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	for n := range ch {
		println(n)
	}
}

//program print numbers from 0 to 9 and end with deadlock
//main goroutine doesnt know about end of sending message
