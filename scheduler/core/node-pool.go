package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	"math/rand"
	"sync"
)

//用于存放所有node,使用内存越小的放在越后面
var nodes = make([]*Node, 0, 100)
var NodesLock sync.Mutex

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
}

//得到Nodes数量
func GetNodeCount() int {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	return len(nodes)
}

//计算nodes的压力,TODO 只有node handler协程才会调用这个方法所以不用加锁
func CalcNodesPress() float64 {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	var totalPress float64
	for _, n := range nodes {
		totalPress += n.CalcNodePress()
	}
	avgPress := totalPress / float64(len(nodes))
	return avgPress
}

////获取container最少的node
//func GetMinContainerNode() *Node {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	var node = nodes[0]
//	for _, n := range nodes {
//		if len(n.Containers) < len(node.Containers) {
//			node = n
//		}
//	}
//	return node
//}

func GetNodes() []*Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	var ns = make([]*Node, 0, 100)
	for _, n := range nodes {
		ns = append(ns, n)
	}
	return ns
}

//根据请求，返回node
func GetSuitableNodes(reqMap map[string]*pb.AcquireContainerRequest) map[string]*Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	size := len(nodes)
	s := rand.Intn(size)
	resMap := make(map[string]*Node)
	for k, _ := range reqMap {
		i := s % size
		resMap[k] = nodes[i]
		s++
	}
	return resMap
}

//
////减少node容量
//func RequireMem(node *Node, reqMem int64) bool {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	//如果内存不够
//	if node.MaxMem-node.UsedMem < reqMem {
//		return false
//	}
//	node.UsedMem += reqMem
//	return true
//}
//
////得到node内存
//func AddContainer(node *Node, container *Container) {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	node.Containers[container.FunName] = container
//	//创建了一个实例，就减少一点容器所占用的空间，申请之前就已经减少了，所以这里不再减少
//}
//
////获得Container
//func GetContainer(node *Node, funcName string) *Container {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	container := node.Containers[funcName]
//	return container
//}
//
////得到node中container的数量
//func GetContainerCount(node *Node) int {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	return len(node.Containers)
//}

func PrintNodes(tag string) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	fmt.Printf("****************************%v*******************************\n", tag)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		mapStr := node.ToString()
		fmt.Printf("No:%v, NodeId:%v, Mem:%v/%v, UserCount:%v, containerCount:%v,  %v\n",
			i, node.NodeID, node.UsedMem/1024/1024,
			node.MaxMem/1024/1024, node.UserCount,
			len(node.CollectionMap), mapStr)
	}
	fmt.Printf("**************************************************************\n\n")

}
