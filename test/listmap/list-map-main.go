package main

import (
	"container/list"
	"fmt"
)

type Node struct {
	NodeId string
}

var L = list.New()
var M = make(map[string]*list.Element)

func Add(id string, node *Node) {
	oldNode := Get(id)
	if oldNode == nil {
		element := L.PushBack(node)
		M[id] = element
	} else {//如果存在就直接更新
		M[id].Value = node
	}
}

func Get(id string) *Node {
	a := M[id]
	if a == nil {
		return nil
	}
	node := a.Value.(*Node)
	return node
}

func Remove(id string) {
	element := M[id]
	if element == nil {
		return
	}
	L.Remove(element)
	delete(M, id)
}

func Println() {
	for i := L.Front(); i != nil; i = i.Next() {
		fmt.Print(*(i.Value.(*Node)), ", ")
	}
	fmt.Println()
	for k := range M {
		fmt.Print(k, ":", *(M[k].Value.(*Node)), ", ")
	}
	fmt.Println()
}

func main() {
	Add("1", &Node{"111"})
	Add("2", &Node{"222"})
	Add("3", &Node{"333"})
	Add("4", &Node{"444"})
	Add("5", &Node{"555"})

	Println()
	Remove("1")
	Println()

	Remove("4")
	Println()
	Add("6", &Node{"666"})
	Println()
	Remove("2")
	Remove("3")
	Remove("5")
	Println()
	Add("7", &Node{"7777"})
	Add("7", &Node{"77777777"})
	Println()
}
