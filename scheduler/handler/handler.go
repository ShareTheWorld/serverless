package handler

import (
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"sync"
	"time"
)

var RequestMap = make(map[string]*core.Container)
var RequestMapLock sync.Mutex

//获取一个container
func AcquireContainer(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
	var isTriggerCreateContainer bool = false
	for {
		container := core.Acquire(req.FunctionName)

		if container == nil {
			if !isTriggerCreateContainer { //如果没有触发创建容器，就去创建容器
				isTriggerCreateContainer = true
				AsyncCreateContainer(req.FunctionName, req.FunctionConfig.Handler,
					req.FunctionConfig.TimeoutInMs, req.FunctionConfig.MemoryInBytes)
			}
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

//返回一个container
func ReturnContainer(req *pb.ReturnContainerRequest) {
	if req == nil {
		return
	}
	container := RequestMap[req.RequestId]
	if container == nil {
		return
	}
	core.Return(container, req.MaxMemoryUsageInBytes, req.DurationInNanos)
}
