package core

import (
	spb "com/aliyun/serverless/scheduler/proto"
	"sync"
)

/*
	提供对外的接口
	Acquire: 获取想要个的container
	Return: 归还container
*/
var RequestMap map[string]*Container
var RequestMapLock sync.Mutex

/*********************************** Acquire and Return 相关api *************************************/

//获取一个container
func Acquire(req *spb.AcquireContainerRequest) *spb.AcquireContainerReply {
	requestId := req.RequestId
	funcName := req.FunctionName
	//TotalMem := req.FunctionConfig.MemoryInBytes

	var container *Container

	m := FunMap[funcName]
	if m == nil { //说明没有这个函数
		return nil
	}

	//挑选一个最优的container
	for _, c := range m {
		if c.UseCount >= c.ConcurrencyCount {
			continue
		}

		if container == nil {
			container = c
			continue
		}

		if c.UseCount < container.UseCount {
			container = c
			continue
		}

		if c.UseCount == container.UseCount && c.UsageMem < c.UsageMem {
			container = c
			continue
		}

	}
	NodesLock.RUnlock()

	if container == nil {
		return nil
	}

	//修改container的使用情况
	NodesLock.Lock()
	container.UseCount++
	container.node.UseCount++
	NodesLock.Unlock()

	//记录请求
	RequestMapLock.Lock()
	RequestMap[requestId] = container
	RequestMapLock.Unlock()

	return &spb.AcquireContainerReply{
		NodeId:          container.node.NodeID,
		NodeAddress:     container.node.Address,
		NodeServicePort: container.node.Port,
		ContainerId:     container.ContainerId,
	}
}

//container
func Return(req *spb.ReturnContainerRequest) {
	requestId := req.RequestId

	RequestMapLock.Lock()
	container := RequestMap[requestId]
	delete(RequestMap, requestId)
	RequestMapLock.Unlock()

	container.ConcurrencyCount = 2 * 1024 * 1024 * 1024 / req.MaxMemoryUsageInBytes
}

/*********************************** other 相关api *************************************/

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
