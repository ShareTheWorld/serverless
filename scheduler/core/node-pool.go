package core

import (
	"fmt"
	"sync"
)

type NC struct {
	Node      *Node
	Container *Container
}

//用于存放所有node,使用内存越小的放在越后面
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)
var RequestMapLock sync.Mutex

//对nodes进行插入排序
func InsertSort(p int, forward bool) {
	if p < 0 || p > len(nodes) {
		return
	}
	if forward { //向前插入，说明增加了使用内存
		for ; p > 0; p-- {
			//如果是正确顺序就直接返回
			if nodes[p-1].UsedMem < nodes[p].UsedMem {
				nodes[p], nodes[p-1] = nodes[p-1], nodes[p]
			} else {
				return
			}
		}
	} else { //向后插入，说明减少了使用内存
		for ; p < len(nodes)-1; p++ {
			//如果是正确顺序就直接返回
			if nodes[p+1].UsedMem > nodes[p].UsedMem {
				nodes[p], nodes[p+1] = nodes[p+1], nodes[p]
			} else {
				return
			}
		}
	}

}

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
	//对node进行排序
	InsertSort(len(nodes)-1, true)
}

//获取第i个位置的节点
func GetNode(i int) *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return nodes[i]
}

//获得内存最大的node
func GetMemMaxNode() *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	node := nodes[len(nodes)-1]
	return node
}

//获得nodes的数量
func NodeCount() int {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return len(nodes)
}

//放入一个请求
func PutRequestNC(requestId string, nc *NC) {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	RequestMap[requestId] = nc
}

//移除一个请求
func RemoveRequestNC(requestId string) {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	delete(RequestMap, requestId)
}

//得到请求
func GetRequestNC(requestId string) *NC {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	nc := RequestMap[requestId]
	return nc
}

func PrintNodes(tag string) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	fmt.Printf("****************************%v*******************************\n", tag)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		fmt.Printf("No:%v, NodeId:%v, Mem:%v/%v, UserCount:%v, containerCount:%v,  %v\n",
			i, node.NodeID, node.UsedMem/1024/1024,
			node.MaxMem/1024/1024, node.UserCount,
			len(node.Containers), node.Containers)
	}
	fmt.Printf("**************************************************************\n\n")

}
