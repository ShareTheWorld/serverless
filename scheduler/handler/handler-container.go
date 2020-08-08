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
//var funcQueue = make(chan *pb.AcquireContainerRequest, 40000)
var NodeMaxContainerCount = 15 //node加载container最大数量

//func AddAcquireContainerToContainerHandler(req *pb.AcquireContainerRequest) {
//	funcQueue <- req
//}

var CreateContainerWG sync.WaitGroup

var FuncNameMap = make(map[string]*pb.AcquireContainerRequest) // New empty set
var FuncNameMapLock sync.Mutex

var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)
var IsChange = false

//添加req到map中
func AddReq(req *pb.AcquireContainerRequest) {
	FuncNameMapLock.Lock()
	if FuncNameMap[req.FunctionName] == nil {
		FuncNameMap[req.FunctionName] = req
		//cond.L.Lock()
		//IsChange = true
		//cond.Signal()
		//cond.L.Unlock()
		go CreateContainerHandler(req)
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
	FuncNameMapLock.Unlock()
	return tmp
}
func ContainerHandler() {
	for {
		cond.L.Lock()
		if !IsChange { //如果没有改变就等待
			cond.Wait()
		}
		IsChange = false
		cond.L.Unlock()
		nodes := core.GetNodes()
		reqMap := GetReq()

		//未每个node加载
		for i := 0; i < core.CollectionCapacity; i++ {
			for _, req := range reqMap {
				for _, node := range nodes {
					CreateContainerWG.Add(1)
					go HandleFuncName(node, req)
				}
			}
			CreateContainerWG.Wait()
			core.PrintNodes(" create container ")

		}
	}
}

func CreateContainerHandler(req *pb.AcquireContainerRequest) {
	nodes := core.GetNodes()
	for i := 0; i < core.CollectionCapacity; i++ {
		for _, node := range nodes {
			CreateContainerWG.Add(1)
			go HandleFuncName(node, req)
		}
		CreateContainerWG.Wait()
		core.PrintNodes(" create container ")
	}
}

func HandleFuncName(node *core.Node, req *pb.AcquireContainerRequest) {
	//判断这个node是否缺乏这个函数实例
	if node.Lack(req.FunctionName) {
		container := CreateContainer(node, req)
		node.AddContainer(container)
	}
	CreateContainerWG.Done()
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
		container := &core.Container{FunName: req.FunctionName, Id: reply.ContainerId, MaxUsedMem: req.FunctionConfig.MemoryInBytes, MaxUsedCount: core.DefaultMaxUsedCount}
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
//			CreateContainerWG.Add(1)
//			go HandleFuncName(resMap[funcName], req)
//		}
//		CreateContainerWG.Wait()
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
