package core

import "sync"

//用于存放所有node
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
}

//移出最后一个node
func RemoveLastNode() *Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	node := nodes[len(nodes)-1]
	nodes = nodes[:len(nodes)-1]
	return node
}

//移除第i个位置的node
func RemoveNode(i int) *Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	node := nodes[i]
	nodes = nodes[i : i+1]
	return node
}

//得到Nodes数量
func GetNodeCount() int {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return len(nodes)
}

//计算nodes的压力,返回内存和cpu的使用
func CalcNodesPress() (float64, float64) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	var TotalTotalMem int64
	var TotalUsageMem int64
	var TotalCpuUsagePct float64

	for _, n := range nodes {
		TotalTotalMem += n.TotalMem
		TotalUsageMem += n.UsageMem
		TotalCpuUsagePct += n.CpuUsagePct
	}

	avgMemUsagePct := float64(TotalUsageMem) / float64(TotalTotalMem)
	avgCpuUsagePct := TotalCpuUsagePct / float64(len(nodes)) / 100.0

	return avgMemUsagePct, avgCpuUsagePct
}

//得到所有的node
func GetNodes() []*Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	var ns = make([]*Node, 0, 100)
	for _, n := range nodes {
		ns = append(ns, n)
	}
	return ns
}

//根据函数名字和需要内存获取n个node,返回的个数小于等于n
func GetSuitableNodes(funcName string, reqMem int64, n int) []*Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	//size := len(nodes)
	//s := rand.Intn(size)
	//resMap := make(map[string]*Node)
	//for k, _ := range reqMap {
	//	i := s % size
	//	resMap[k] = nodes[i]
	//	s++
	//}
	//return resMap
	return nil
}
