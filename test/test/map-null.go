package main

import "fmt"

type Container struct {
	Type int
}

var ContainerIdMap = make(map[string]Container)

func main() {
	container := ContainerIdMap["a"]
	fmt.Println(container)

	var c Container = Container{Type: 1}
	fmt.Println(c)
	get(c)
	c1 := c
	fmt.Println(c1)
}

func get(container Container) {
	fmt.Println(&container)
	container.Type = 10
}
