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
var funcQueue = make(chan *pb.AcquireContainerRequest, 10000)

func AddAcquireContainerToContainerHandler(req *pb.AcquireContainerRequest) {
	funcQueue <- req
}

func ContainerHandler() {
	fmt.Println("start handle create container")
	for {
		req := <-funcQueue
		node := core.GetMemMaxNode()
		container := node.GetContainer(req.FunctionName)
		if container != nil {
			continue
		}
		container = CreateContainer(node, req)
		node.AddContainer(container)
	}
}

//
////保证创建一个container
func CreateContainer(node *core.Node, req *pb.AcquireContainerRequest) *core.Container {
	core.PrintNodes("create container")
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
			return nil
		}

		//将container添加到node中
		container := &core.Container{FunName: req.FunctionName, Id: reply.ContainerId, UsedMem: req.FunctionConfig.MemoryInBytes}
		et := time.Now().UnixNano()
		fmt.Printf("---- create container, time=%v, node:%v \n", (et-st)/1000000, node)
		return container
	}
}
