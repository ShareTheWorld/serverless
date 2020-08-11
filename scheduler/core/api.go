package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"math/rand"
	"sync"
)

/*
	提供对外的接口
	Acquire: 获取想要个的container
	Return: 归还container
*/
type NC struct {
	Node      *Node
	Container *Container
}

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)

var RequestMapLock sync.Mutex

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
