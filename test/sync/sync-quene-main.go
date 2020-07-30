package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	myMap = make(map[int]int, 10)
	//lock是全局互斥锁,synchornized
	//lock sync.Mutex
	Lock sync.Mutex
)

func cal(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}

func main() {
	for i := 1; i <= 15; i++ {
		t := i
		go func() {
			r := test(t)
			fmt.Printf("%v!=%v\n", t, r)
		}()
	}
	time.Sleep(time.Millisecond * 1000)
}
func test(n int) int {
	c := make(chan int)
	go run(n, c)
	r := <-c
	return r
}
func run(n int, ch chan int) {
	r := cal(n)
	ch <- r
}
