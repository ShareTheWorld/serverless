package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

//创建一个容器，成功返回true，失败或者已经存在返回false
func CreateContainer(funcName string, handler string, timeoutInMs int64, memoryInBytes int64) bool {
	node := core.GetSuitableNode(funcName, memoryInBytes)
	if node == nil {
		return false
	}
	b := CreateContainerForNode(node, funcName, handler, timeoutInMs, memoryInBytes)
	return b
}

func CreateContainerForNode(node *core.Node, funcName string, handler string, timeoutInMs int64, memoryInBytes int64) bool {
	st := time.Now().UnixNano()
	for {
		//创建一个container
		reply, err := client.CreateContainer(node.Client, uuid.NewV4().String(), funcName+uuid.NewV4().String(),
			funcName, handler, timeoutInMs, 4*1024*1024*1024)

		if err != nil {
			fmt.Printf("FuncName:%v, Mem:%v, error: %v", funcName, memoryInBytes/1048576, err)
			return false
		}

		//将container添加到node中
		container := &core.Container{
			ContainerId: reply.ContainerId,
			TotalMem:    4 * 1024 * 1024 * 1024,
			UsageMem:    128 * 1024 * 1024,
			CpuUsagePct: 0,

			FuncName:      funcName,
			Handler:       handler,
			TimeoutInMs:   timeoutInMs,
			MemoryInBytes: memoryInBytes,

			UseCount:         0,
			ConcurrencyCount: 4,
			Node:             node,
		}

		node.AddContainer(container)
		core.AddContainer(container)

		et := time.Now().UnixNano()
		fmt.Printf("create container,FuncName:%v, Mem:%v, time=%v, nodeId=%v\n", funcName, memoryInBytes/1024/1024, (et-st)/1000000, node.NodeID)
		return true
	}

}
