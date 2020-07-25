package core

import (
	"container/list"
	"fmt"
)

//ListMap工具，数据放在List，使用Map建立索引
type LM struct {
	L *list.List
	M map[string]*list.Element
}

func NewLM() *LM {
	lm := new(LM)
	lm.L = list.New()
	lm.M = make(map[string]*list.Element)
	return lm
}

func (lm *LM) Add(id string, node interface{}) {
	oldNode := lm.Get(id)
	if oldNode == nil {
		element := lm.L.PushBack(node)
		lm.M[id] = element
	} else { //如果存在就直接更新
		lm.M[id].Value = node
	}
}

func (lm *LM) Get(id string) interface{} {
	a := lm.M[id]
	if a == nil {
		return nil
	}
	node := a.Value.(*interface{})
	return node
}

func (lm *LM) Remove(id string) {
	element := lm.M[id]
	if element == nil {
		return
	}
	lm.L.Remove(element)
	delete(lm.M, id)
}

func (lm *LM) Println() {
	for i := lm.L.Front(); i != nil; i = i.Next() {
		fmt.Print(i.Value.(interface{}), ", ")
	}
	fmt.Println()
	for k := range lm.M {
		fmt.Print(k, ":", lm.M[k].Value.(interface{}), ", ")
	}
	fmt.Println()
}

//func main() {
//	lm := NewLM()
//	lm.Add("1", Node{NodeID: "1111"})
//	lm.Add("2", Node{NodeID: "2222"})
//	lm.Add("3", Node{NodeID: "3333"})
//	lm.Remove("2")
//	lm.Add("2", Node{NodeID: "2222"})
//	lm.Println()
//
//}
