package main

import (
	"fmt"
)

func main() {
	//str := ""
	s := fmt.Sprintf("aa=%v", 1)
	if s == "aa=1" {
		fmt.Printf("****")
	}
	fmt.Printf("end")
}
