package core

import (
	"com/aliyun/serverless/scheduler/client"
	pb "com/aliyun/serverless/scheduler/proto"
	uuid "github.com/satori/go.uuid"
	"sync"
)

//初始化，开头盛情一定数量的节点
func Init() {

}

//请求
func AcquireContainer(req *pb.AcquireContainerRequest) (*pb.AcquireContainerReply, error) {
	//node和node里面的container信息
	node, container := QueryNodeAndContainer(req.FunctionName, req.FunctionConfig.MemoryInBytes)

	//如果node为nil，就实力化创建一个一个
	if node == nil {
		var err error
		node, err = GetNode(req.AccountId, req.FunctionName, req.FunctionConfig.MemoryInBytes)
		if err != nil {
			return nil, err
		}
	}

	if container == nil {
		//创建一个container
		reply, err := client.CreateContainer(
			node.Client,
			req.RequestId,                          //demo是这样
			req.FunctionName+uuid.NewV4().String(), //demo是这样
			req.FunctionName,
			req.FunctionConfig.Handler,
			req.FunctionConfig.TimeoutInMs,
			req.FunctionConfig.MemoryInBytes,
		)
		if err != nil {
			return nil, err
		}

		//将container添加到node中
		container = &Container{FunName: req.FunctionName, Id: reply.ContainerId, UsedMem: req.FunctionConfig.MemoryInBytes}
		AddNC(node, container)
	}

	container, err := RentNC(req.RequestId, node, container)
	if err != nil { //租用container出错
		return nil, err
	}

	r := pb.AcquireContainerReply{
		NodeId:          node.NodeID,
		NodeAddress:     node.Address,
		NodeServicePort: node.Port,
		ContainerId:     container.Id,
	}
	return &r, nil
}

//归还
func ReturnContainer(req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {
	//req{RequestId,ContainerId,DurationInNanos,MaxMemoryUsageInBytes,ErrorCode,ErrorMessage}
	ReturnNC(req.RequestId)
	return &pb.ReturnContainerReply{}, nil
}

var NodeLock sync.Mutex

//保证所有的node申请是排队完成的，
func GetNode(accountId string, funcName string, reqMem int64) (*Node, error) {
	NodeLock.Lock()
	defer NodeLock.Unlock()

	//在申请之前先查询一次，是否前面已经申请了，保证一次只申请一个
	node, _ := QueryNodeAndContainer(funcName, reqMem)
	if node != nil {
		return node, nil
	}

	//预约一个node
	reply, err := client.ReserveNode("", accountId)
	if err != nil {
		return nil, err
	}

	//ReservedTimeTimestampMs ReleasedTimeTimestampMs
	node = NewNode(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort, reply.Node.MemoryInBytes)
	nodeClient, err := client.ConnectNodeService(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort)
	if err != nil {
		//TODO 由于连接错误，需要释放Node
		return nil, err
	}

	//创建成功node并且连接成功，进行节点添加
	node.Client = nodeClient
	AddNode(node)
	return node, nil
}
