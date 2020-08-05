package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
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
	var press float64
	var container *Container

	NodesLock.RLock()
	//发现一个满足要求，且压力最小的node
	for _, n := range nodes {
		if !n.Satisfy(funcName, reqMem) {
			continue
		}
		p := n.CalcNodePress()

		//如果node为null，就直接赋值
		if node == nil {
			node, press = n, p
			continue
		}

		//如果p的压力比选中的压力小，就使用新的
		if p < press {
			node, press = n, p
		}
	}
	NodesLock.RUnlock()
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
	defer RequestMapLock.Unlock()

	nc := RequestMap[requestId]
	delete(RequestMap, requestId)

	if nc == nil {
		return
	}

	node := nc.Node
	container := nc.Container
	if count%100 == 0 {
		PrintNodes("timer")
	}
	count++

	node.Return(container)
}
