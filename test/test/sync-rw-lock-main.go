package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	count int64
	lock  sync.RWMutex
)

func read() {
	for {
		lock.Lock()
		time.Sleep(time.Millisecond * 200)
		lock.Unlock()
	}
}

func write() {
	st := time.Now().UnixNano()
	for i := 0; i < 1; i++ {
		lock.Lock()
		time.Sleep(time.Millisecond * 1)
		lock.Unlock()
	}
	et := time.Now().UnixNano()
	diff := et - st
	fmt.Println(diff / 1000 / 1000)
}

func main() {
	go read()
	time.Sleep(time.Millisecond * 50)
	go read()
	time.Sleep(time.Millisecond * 50)
	go read()
	time.Sleep(time.Millisecond * 50)
	go read()

	go write()

	time.Sleep(time.Second * 100000)
}
