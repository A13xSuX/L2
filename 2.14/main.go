package main

import (
	"fmt"
	"reflect"
	"time"
)

func orChan(channels ...<-chan interface{}) <-chan interface{} {
	result := make(chan interface{})
	cases := make([]reflect.SelectCase, 0, len(channels))
	if len(channels) == 0 { //if we got 0 channels
		close(result)
		return result
	}
	go func() { //merge for each channel
		for _, ch := range channels {
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
		}
		for {
			_, _, ok := reflect.Select(cases)
			if !ok { //if we got signal -> done all
				break
			}
		}
		close(result)
	}()

	return result
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-orChan(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("Done after %v\n", time.Since(start))
}
