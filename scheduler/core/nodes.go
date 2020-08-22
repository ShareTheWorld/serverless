package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
)

//用于存放所有node
var nodes = make([]*Node, 0, 100)
var Lock sync.RWMutex //整个数据的一把锁

//添加一个Node
func AddNode(node *Node) {
	Lock.Lock()
	defer Lock.Unlock()
	nodes = append(nodes, node)
}

//移出最后一个node
func RemoveLastNode() *Node {
	Lock.Lock()
	defer Lock.Unlock()
	node := nodes[len(nodes)-1]
	nodes = nodes[:len(nodes)-1]
	return node
}

//移除第i个位置的node
func RemoveNode(i int) *Node {
	Lock.Lock()
	defer Lock.Unlock()
	node := nodes[i]
	nodes = nodes[i : i+1]
	return node
}

//得到Nodes数量
func GetNodeCount() int {
	Lock.RLock()
	defer Lock.RUnlock()
	return len(nodes)
}

//计算nodes的压力,返回内存和cpu的使用
func CalcNodesPress() (float64, float64) {
	Lock.RLock()
	defer Lock.RUnlock()

	var TotalTotalMem int64
	var TotalUsageMem int64
	var TotalCpuUsagePct float64

	for _, n := range nodes {
		TotalTotalMem += n.TotalMem
		TotalUsageMem += n.UsageMem
		TotalCpuUsagePct += n.CpuUsagePct
	}

	//fmt.Printf("calc press: TotalTotalMem:%v, TotalUsageMem:%v, TotalCpuUsagePct:%v\n", TotalTotalMem, TotalUsageMem, TotalCpuUsagePct)
	avgMemUsagePct := float64(TotalUsageMem) / float64(TotalTotalMem)
	avgCpuUsagePct := TotalCpuUsagePct / float64(len(nodes)) / 100.0

	return avgMemUsagePct, avgCpuUsagePct
}

//得到所有的node
func GetNodes() []*Node {
	Lock.RLock()
	defer Lock.RUnlock()
	var ns = make([]*Node, 0, 100)
	for _, n := range nodes {
		ns = append(ns, n)
	}
	return ns
}

//得到一个合适的node
func GetSuitableNode(funcName string, reqMem int64) *Node {
	Lock.Lock()
	defer Lock.Unlock()
	size := len(nodes)
	s := rand.Intn(size)
	var node *Node
	for i := 0; i < size; i++ {
		n := nodes[(s+i)%size]

		if n.FuncNameMap[funcName] { //如果这个node中已经有这个函数了，就不是合适的node
			continue
		}

		if node == nil {
			node = n
			continue
		}
		if n.AvailableMem > node.AvailableMem {
			node = n
		}
	}
	if node != nil {
		node.FuncNameMap[funcName] = true
		node.AvailableMem -= 128 * 1024 * 1024 //减少一点可用内存，避免下次再选中它，状态同步任务最后会自动修复它
	}
	return node
}

func PrintNodes(tag string) {
	Lock.RLock()
	defer Lock.RUnlock()
	size := len(nodes)
	fmt.Printf("*******************************%v**************************\n", tag)
	for i := 0; i < size; i++ {
		node := nodes[i]
		jsonBytes1, _ := json.Marshal(node)
		jsonBytes2, _ := json.Marshal(node.ContainerIdMap)
		fmt.Println(string(jsonBytes1) + "," + string(jsonBytes2))
	}
	fmt.Printf("*******************************%v**************************\n", tag)

}



//
////计算nodes的压力,返回内存和cpu的使用
//func CalcNodesPress() (float64, float64) {
//	Lock.RLock()
//	defer Lock.RUnlock()
//
//	var TotalTotalMem int64
//	var TotalUsageMem int64
//	//var TotalCpuUsagePct float64
//	var MaxCpuUseAgePct float64
//
//	for _, n := range nodes {
//		TotalTotalMem += n.TotalMem
//		TotalUsageMem += n.UsageMem
//		if n.CpuUsagePct > MaxCpuUseAgePct {
//			MaxCpuUseAgePct = n.CpuUsagePct
//		}
//	}
//
//	//fmt.Printf("calc press: TotalTotalMem:%v, TotalUsageMem:%v, TotalCpuUsagePct:%v\n", TotalTotalMem, TotalUsageMem, TotalCpuUsagePct)
//	avgMemUsagePct := float64(TotalUsageMem) / float64(TotalTotalMem)
//	//avgCpuUsagePct := TotalCpuUsagePct / float64(len(nodes)) / 100.0
//	avgCpuUsagePct := MaxCpuUseAgePct / 100.0 //把最大值作为cpu的平局值
//
//	return avgMemUsagePct, avgCpuUsagePct
//}