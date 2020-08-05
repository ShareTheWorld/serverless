package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

/*
	处理container的加载
*/
var FuncNameSet = make(map[string]bool) // New empty set
var funcQueue = make(chan *pb.AcquireContainerRequest, 10000)
var NodeMaxContainerCount = 15 //node加载container最大数量

func AddAcquireContainerToContainerHandler(req *pb.AcquireContainerRequest) {
	funcQueue <- req
}

func ContainerHandler() {
	fmt.Println("start handle create container")
	for {
		req := <-funcQueue
		funcName := req.FunctionName

		//将方法名字放入到set中
		FuncNameSet[funcName] = true

		//获取容器实例话最小的少的节点
		node := core.GetMinContainerNode()

		containerCount := node.GetContainerCount()
		if containerCount >= NodeMaxContainerCount {
			continue
		}

		//如果这个node中有这个container就不创建了
		container := node.GetContainer(req.FunctionName)
		if container != nil {
			continue
		}

		//等待有足够的内存了就去创建容器
		container = CreateContainer(node, req)
		node.AddContainer(container)
		core.PrintNodes(fmt.Sprintf("create container fn:%v, mem:%v", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576))
	}
}

//
////保证创建一个container
func CreateContainer(node *core.Node, req *pb.AcquireContainerRequest) *core.Container {
	st := time.Now().UnixNano()
	for {
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
			fmt.Printf("FuncName:%v, Mem:%v, error: %v", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576, err)
			return nil
		}

		//将container添加到node中
		container := &core.Container{FunName: req.FunctionName, Id: reply.ContainerId, UsedMem: req.FunctionConfig.MemoryInBytes}
		et := time.Now().UnixNano()
		fmt.Printf("create container,FuncName:%v, Mem:%v, time=%v, node:%v \n", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576, (et-st)/1000000, node)
		return container
	}
}
