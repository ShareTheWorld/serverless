package main

import "fmt"

var s = make([]int, 0, 10)

func main() {
	//str := ""

	s = append(s, 1)
	s = append(s, 2)
	s = append(s, 3)
	s = append(s, 4)
	s = s[:len(s)-1]
	fmt.Println(s)

}
