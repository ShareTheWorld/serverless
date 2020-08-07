package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

/*
	处理container的加载
*/
var funcQueue = make(chan *pb.AcquireContainerRequest, 40000)
var NodeMaxContainerCount = 15 //node加载container最大数量

func AddAcquireContainerToContainerHandler(req *pb.AcquireContainerRequest) {
	funcQueue <- req
}

var wg sync.WaitGroup

var FuncNameMap = make(map[string]*pb.AcquireContainerRequest) // New empty set
var FuncNameMapLock sync.Mutex

var isChange = false

//添加req到map中
func AddReq(req *pb.AcquireContainerRequest) {
	FuncNameMapLock.Lock()
	if FuncNameMap[req.FunctionName] == nil {
		FuncNameMap[req.FunctionName] = req
		isChange = true
	}
	FuncNameMapLock.Unlock()
}

//复制一份出来
func GetReq() map[string]*pb.AcquireContainerRequest {
	tmp := make(map[string]*pb.AcquireContainerRequest) // New empty set
	FuncNameMapLock.Lock()
	for k, v := range FuncNameMap {
		tmp[k] = v
	}
	isChange = false
	FuncNameMapLock.Unlock()
	return tmp
}
func ContainerHandler() {
	for {
		if !isChange {
			time.Sleep(time.Millisecond * 1000) //随眠100毫秒
			continue
		}
		nodes := core.GetNodes()
		reqMap := GetReq()

		//未每个node加载
		for _, req := range reqMap {
			for _, node := range nodes {
				wg.Add(1)
				go HandleFuncName(node, req)
			}
		}
		wg.Wait()
		//fmt.Printf("create finsih")
	}

}
func HandleFuncName(node *core.Node, req *pb.AcquireContainerRequest) {
	//如果这个node中有这个container就不创建了
	container := node.GetContainer(req.FunctionName)
	if container == nil {
		container = CreateContainer(node, req)
		node.AddContainer(container)
		core.PrintNodes(fmt.Sprintf("create container fn:%v, mem:%v", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576))
	}
	wg.Done()
}

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
		container := &core.Container{FunName: req.FunctionName, Id: reply.ContainerId, MaxUsedMem: req.FunctionConfig.MemoryInBytes, MaxUsedCount: 1}
		et := time.Now().UnixNano()
		fmt.Printf("create container,FuncName:%v, Mem:%v, time=%v, nodeId=%v\n", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576, (et-st)/1000000, node.NodeID)
		return container
	}
}

//func ContainerHandler() {
//	//fmt.Println("start handle create container")
//	for {
//		reqMap := GetWaitFuncName()
//		if len(reqMap) == 0 {
//			time.Sleep(time.Millisecond * 100)
//			continue
//		}
//		resMap := core.GetSuitableNodes(reqMap)
//		for funcName, req := range reqMap {
//			wg.Add(1)
//			go HandleFuncName(resMap[funcName], req)
//		}
//		wg.Wait()
//		LoadFinishContainer()
//	}
//}

//
//func ContainerHandler() {
//	fmt.Println("start handle create container")
//	for {
//		req := <-funcQueue
//		funcName := req.FunctionName
//
//		//将方法名字放入到set中
//		//FuncNameMap[funcName] = true
//
//		//获取容器实例话最小的少的节点
//		node := core.GetMinContainerNode()
//
//		containerCount := node.GetContainerCount()
//		if containerCount >= NodeMaxContainerCount {
//			continue
//		}
//
//		//如果这个node中有这个container就不创建了
//		container := node.GetContainer(req.FunctionName)
//		if container != nil {
//			continue
//		}
//
//		//等待有足够的内存了就去创建容器
//		container = CreateContainer(node, req)
//		node.AddContainer(container)
//		core.PrintNodes(fmt.Sprintf("create container fn:%v, mem:%v", req.FunctionName, req.FunctionConfig.MemoryInBytes/1048576))
//	}
//}
