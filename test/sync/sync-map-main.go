package main

import (
	"fmt"
	"sync"
)

//sync.Map 并发安全的map
var mapWg sync.WaitGroup

var m = make(map[int]int)
var m2 = sync.Map{}

func get(key int) int {
	return m[key]
}

func set(key int, value int) {
	m[key] = value
}

//func main() {
//	for i := 0; i < 20; i++ {
//		mapWg.Add(1)
//		go func() {
//			set(i, i+100)
//			fmt.Printf("k-v: %v - %v\n", i, get(i))
//			mapWg.Done()
//		}()
//	}
//	mapWg.Wait()
//}

func main() {
	for i := 0; i < 20; i++ {
		mapWg.Add(1)
		go func(i int) {
			m2.Store(i, i+100)
			value, _ := m2.Load(i)
			fmt.Printf("k-v: %v - %v\n", i, value)
			mapWg.Done()
		}(i)
	}
	mapWg.Wait()
}
