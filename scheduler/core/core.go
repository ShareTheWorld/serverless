package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"math/rand"
	"sync"
)


//请求相关

//请求表，用于存放所有的请求
//var RequestMap = make(map[string]*NC)
//var RequestMapLock sync.Mutex

//获取一个node里面的container
func Acquire(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
	requestId := req.RequestId
	funcName := req.FunctionName
	reqMem := req.FunctionConfig.MemoryInBytes

	var node *Node
	var count int64 //表示container的会用数量
	var container *Container

	NodesLock.Lock()
	//发现一个满足要求的container，且使用人数是最少的container
	size := len(nodes)
	s := rand.Intn(size) //随机选择一个开始位置
	for i := 0; i < size; i++ {
		p := (i + s) % size
		n := nodes[p]
		satisfy, usedCount := n.Satisfy(funcName, reqMem)
		if !satisfy { //如果不满足直接返回
			continue
		}

		//如果node为null，就直接赋值
		if node == nil {
			node, count = n, usedCount
			continue
		}

		//如果使用数量少，就替换
		if usedCount < count {
			node, count = n, usedCount
		}
	}
	NodesLock.Unlock()
	//如果没有找到合适的node，就返回nil
	if node == nil {
		return nil
	}

	//获取container
	container = node.Acquire(funcName, reqMem)

	//记录请求
	RequestMapLock.Lock()
	RequestMap[requestId] = &NC{node, container}
	RequestMapLock.Unlock()

	return &pb.AcquireContainerReply{
		NodeId:          node.NodeID,
		NodeAddress:     node.Address,
		NodeServicePort: node.Port,
		ContainerId:     container.Id,
	}
}

var count = 0

//归还node中的container
func Return(req *pb.ReturnContainerRequest) {
	requestId := req.RequestId

	RequestMapLock.Lock()
	nc := RequestMap[requestId]
	delete(RequestMap, requestId)
	RequestMapLock.Unlock()

	if nc == nil {
		return
	}

	node := nc.Node
	container := nc.Container
	//******************log*************************
	//if count%100 == 0 {
	//	PrintNodes("timer")
	//}
	//count++
	//******************log*************************
	actualUseMem := req.MaxMemoryUsageInBytes
	node.Return(container, actualUseMem)
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
//	node.Containers[container.FuncName] = container
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

////打印node方便调试的时候查看node-pool的信息
//func PrintNodes(tag string) {
//	NodesLock.Lock()
//	defer NodesLock.Unlock()
//	fmt.Printf("****************************%v*******************************\n", tag)
//	for i := 0; i < len(nodes); i++ {
//		node := nodes[i]
//		mapStr := node.ToString()
//		fmt.Printf("No:%v, NodeId:%v, Mem:%v/%v, UserCount:%v, containerCount:%v,  %v\n",
//			i, node.NodeID, node.UsedMem/1024/1024,
//			node.MaxMem/1024/1024, node.UserCount,
//			len(node.CollectionMap), mapStr)
//	}
//	fmt.Printf("**************************************************************\n\n")
//
//}
