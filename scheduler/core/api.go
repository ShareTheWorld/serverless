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

//归还container
func Return(req *spb.ReturnContainerRequest) {
	requestId := req.RequestId

	RequestMapLock.Lock()
	container := RequestMap[requestId]
	delete(RequestMap, requestId)
	RequestMapLock.Unlock()

	container.ConcurrencyCount = 2 * 1024 * 1024 * 1024 / req.MaxMemoryUsageInBytes
}
