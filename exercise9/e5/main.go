package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	done := make(chan struct{})
	res := make(chan int)
	go func() {
		for {
			x := <-ch2
			ch1 <- x + 1
		}
	}()
	go func() {
		for {
			x := <- ch1
			select {
			case <-done:
				res <- x
				return
			default:
				ch2 <- x + 1
			}
		}
	}()
	ch1 <- 0
	time.Sleep(time.Second)
	done <- struct{}{}
	fmt.Println(<-res)
}
