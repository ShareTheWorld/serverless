package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 5)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5
		ch <- 6
		ch <- 7
	}()

	for {
		i := <-ch
		fmt.Println(i)
		//if i%2 == 0 {
		//	ch <- i
		//}
		time.Sleep(time.Second * 1)
	}
}
