package core

import pb "com/aliyun/serverless/nodeservice/proto"

/*
	提供对外的接口
	Acquire: 获取想要个的container
	Return: 归还container
*/

/*********************************** nodes 相关api *************************************/
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

	avgUsageMem := float64(TotalUsageMem) / float64(TotalTotalMem)
	avgCpuUsagePct := TotalCpuUsagePct / float64(len(nodes))

	return avgUsageMem, avgCpuUsagePct
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

/*********************************** node 相关api *************************************/

//更新node的状态
func UpdateNodeStats(node *Node, stats *pb.NodeStats) {
	if node == nil || stats == nil {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.updateNodeStats(stats)
}

//更新所有container的状态
func UpdateContainer(node *Node, stats []*pb.ContainerStats) {
	if node == nil || stats == nil {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.updateContainer(stats)
}

//添加container
func AddContainer(node *Node, container *Container) {
	if node == nil || container == nil {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()
	node.addContainer(container)

	m := FunMap[container.FuncName]
	if m == nil {
		m = make(map[string]*Container)
		FunMap[container.FuncName] = m
	}
	m[container.ContainerId] = container
}

//根据函数名字移除container
func RemoveContainerByFuncName(node *Node, funcName string) {
	if node == nil {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()

	container := node.removeContainerByFuncName(funcName)
	if container == nil {
		return
	}
	m := FunMap[container.FuncName]
	if m == nil {
		return
	}
	delete(m, container.ContainerId)
}

//根据containerId移除container
func RemoveContainerByContainerId(node *Node, containerId string) {
	if node == nil {
		return
	}
	node.lock.Lock()
	defer node.lock.Unlock()

	container := node.removeContainerByContainerId(containerId)
	if container == nil {
		return
	}
	m := FunMap[container.FuncName]
	if m == nil {
		return
	}
	delete(m, container.ContainerId)
}

//
////node和container的关系结构，在申请和归还的时候会用上
//type NC struct {
//	Node      *Node
//	Container *Container
//}
//
////请求表，用于存放所有的请求
////var RequestMap = make(map[string]*NC)
//
////var RequestMapLock sync.Mutex
//
////获取一个node里面的container
//func Acquire(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
//	requestId := req.RequestId
//	funcName := req.FunctionName
//	reqMem := req.FunctionConfig.MemoryInBytes
//
//	var node *Node
//	var count int64 //表示container的会用数量
//	var container *Container
//
//	NodesLock.Lock()
//	//发现一个满足要求的container，且使用人数是最少的container
//	size := len(nodes)
//	s := rand.Intn(size) //随机选择一个开始位置
//	for i := 0; i < size; i++ {
//		p := (i + s) % size
//		n := nodes[p]
//		satisfy, usedCount := n.Satisfy(funcName, reqMem)
//		if !satisfy { //如果不满足直接返回
//			continue
//		}
//
//		//如果node为null，就直接赋值
//		if node == nil {
//			node, count = n, usedCount
//			continue
//		}
//
//		//如果使用数量少，就替换
//		if usedCount < count {
//			node, count = n, usedCount
//		}
//	}
//	NodesLock.Unlock()
//	//如果没有找到合适的node，就返回nil
//	if node == nil {
//		return nil
//	}
//
//	//获取container
//	container = node.Acquire(funcName, reqMem)
//
//	//记录请求
//	RequestMapLock.Lock()
//	RequestMap[requestId] = &NC{node, container}
//	RequestMapLock.Unlock()
//
//	return &pb.AcquireContainerReply{
//		NodeId:          node.NodeID,
//		NodeAddress:     node.Address,
//		NodeServicePort: node.Port,
//		ContainerId:     container.Id,
//	}
//}
//
//var count = 0
//
////归还node中的container
//func Return(req *pb.ReturnContainerRequest) {
//	requestId := req.RequestId
//
//	RequestMapLock.Lock()
//	nc := RequestMap[requestId]
//	delete(RequestMap, requestId)
//	RequestMapLock.Unlock()
//
//	if nc == nil {
//		return
//	}
//
//	node := nc.Node
//	container := nc.Container
//	//******************log*************************
//	//if count%100 == 0 {
//	//	PrintNodes("timer")
//	//}
//	//count++
//	//******************log*************************
//	actualUseMem := req.MaxMemoryUsageInBytes
//	node.Return(container, actualUseMem)
//}
