package handler

import (
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"sync"
	"time"
)

var RequestMap map[string]*core.Container
var RequestMapLock sync.Mutex

func AcquireContainer(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
	for {
		container := core.Acquire(req.FunctionName)
		
		if container == nil {
			//触发缺失
			time.Sleep(time.Millisecond * 1)
			continue
		}

		//记录请求
		RequestMapLock.Lock()
		RequestMap[req.RequestId] = container
		RequestMapLock.Unlock()

		res := &pb.AcquireContainerReply{
			NodeId:          container.Node.NodeID,
			NodeAddress:     container.Node.Address,
			NodeServicePort: container.Node.Port,
			ContainerId:     container.ContainerId,
		}
		return res
	}
}

func ReturnContainer(req *pb.ReturnContainerRequest) {
	core.Return(req)
}
