package handler

import (
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
)

/*
	提供对外的接口
	Acquire: 获取想要个的container
	Return: 归还container
*/

//检索想要的node和container
func Acquire(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
	requestId := req.RequestId
	funcName := req.FunctionName
	reqMem := req.FunctionConfig.MemoryInBytes

	var node *core.Node
	var container *core.Container
	for i := 0; i < core.NodesSize(); i++ {
		n := core.GetNode(i)
		c := n.Containers[funcName]
		if c == nil { //判断是否存在想要的方法
			continue
		}

		if n.MaxMem-n.UsedMem > reqMem { //判断内存是否足够
			node = n
			container = c
			break
		}
	}

	if node == nil || container == nil {
		return nil
	}

	//对node做相应的操作
	node.UsedMem -= reqMem
	node.UserCount++

	//在requestMap上做好登记
	core.PutRequestNC(requestId, &core.NC{Node: node, Container: container})

	//TODO 将方法的请求时间记录下来

	return &pb.AcquireContainerReply{
		NodeId:          node.NodeID,
		NodeAddress:     node.Address,
		NodeServicePort: node.Port,
		ContainerId:     container.Id,
	}
}

//归还container,只需要更具请求者的id归还就行
func Return(req *pb.ReturnContainerRequest) {
	requestId := req.RequestId
	nc := core.GetRequestNC(requestId)
	if nc == nil {
		return
	}

	node := nc.Node
	container := nc.Container

	node.UsedMem -= container.UsedMem
	node.UserCount--
	core.RemoveRequestNC(requestId)
}
