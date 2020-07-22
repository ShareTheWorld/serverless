package core

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
	} else { //如果存在就直接更新
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
