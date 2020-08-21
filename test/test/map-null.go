package main

import "fmt"

type Container struct {
	Type int
}

var ContainerIdMap = make(map[string]Container)

func main() {
	container := ContainerIdMap["a"]
	fmt.Println(container)
}
