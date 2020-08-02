package core

import (
	"sync"
)

type NC struct {
	Node      *Node
	Container *Container
}

//用于存放所有node
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)
var RequestMapLock sync.Mutex

//对nodes进行插入排序
func InsertSort(p int, inc bool) {
	if p < 0 || p > len(nodes) {
		return
	}
	if inc { //增序
		for ; p < len(nodes)-1; p++ {
			//如果是正确顺序就直接返回
			if nodes[p].UsedMem > nodes[p+1].UsedMem {
				nodes[p], nodes[p+1] = nodes[p+1], nodes[p]
			} else {
				return
			}
			//交换位置
		}
	} else { //逆序
		for ; p > 0; p-- {
			//如果是正确顺序就直接返回
			if nodes[p].UsedMem < nodes[p-1].UsedMem {
				nodes[p], nodes[p-1] = nodes[p-1], nodes[p]
			} else {
				return
			}
			//交换位置
		}
	}

}

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
	//对node进行排序
	InsertSort(len(nodes)-1, false)
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
