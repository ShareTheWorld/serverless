package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	count       int64
	waitGroup   sync.WaitGroup
	mutexLock   sync.Mutex
	rwMutexLock sync.RWMutex
)

func read() {
	//mutexLock.Lock()
	rwMutexLock.RLock()
	time.Sleep(time.Millisecond)
	//mutexLock.Unlock()
	rwMutexLock.RUnlock()
	waitGroup.Done()
}

func write() {
	//mutexLock.Lock()
	rwMutexLock.Lock()
	count += 1
	time.Sleep(time.Millisecond * 10)
	//mutexLock.Unlock()
	rwMutexLock.Unlock()
	waitGroup.Done()
}

func main2() {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		waitGroup.Add(1)
		go read()
	}

	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go write()
	}
	waitGroup.Wait()
	fmt.Println(time.Now().Sub(start))
}
