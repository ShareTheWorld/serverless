package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"sync"
)

//Node结构
//锁：用于控制同步
//连接信息
//节点状态
//本地使用状态
//Container信息
type Node struct {
	lock sync.RWMutex

	NodeID  string               //节点id
	Address string               //节点地址
	Port    int64                //节点端口
	Client  pb.NodeServiceClient //节点连接

	TotalMem     int64   //总内存
	UsageMem     int64   //使用内存
	AvailableMem int64   //可用内存
	CpuUsagePct  float64 //cpu使用百分比

	UseCount         int64 //当前正在使用的人数
	ConcurrencyCount int64 //并发数量

	FuncNameMap    map[string]*Container //存放所有的Container K:V=containerId:Container
	ContainerIdMap map[string]*Container //存放所有的Container K:V=containerId:Container
}

//更新node的状态
func (n *Node) updateNodeStats(stats *pb.NodeStats) {
	n.TotalMem = stats.TotalMemoryInBytes
	n.UsageMem = stats.MemoryUsageInBytes
	n.AvailableMem = stats.AvailableMemoryInBytes
	n.CpuUsagePct = stats.CpuUsagePct

}

//更新所有container的状态
func (n *Node) updateContainer(stats []*pb.ContainerStats) {
	for _, s := range stats {
		if s == nil {
			continue
		}
		container := n.ContainerIdMap[s.ContainerId]
		if container == nil {
			continue
		}
		container.UpdateContainerStats(s)

	}
}

//添加container
func (n *Node) addContainer(container *Container) {
	n.FuncNameMap[container.FuncName] = container
	n.ContainerIdMap[container.ContainerId] = container
}

//根据函数名字移除container
func (n *Node) removeContainerByFuncName(funcName string) *Container {
	container := n.FuncNameMap[funcName]
	if container == nil {
		return nil
	}

	delete(n.FuncNameMap, funcName)
	delete(n.ContainerIdMap, container.ContainerId)
	return container
}

//根据containerId移除container
func (n *Node) removeContainerByContainerId(containerId string) *Container {
	container := n.ContainerIdMap[containerId]
	if container == nil {
		return nil
	}

	delete(n.FuncNameMap, container.FuncName)
	delete(n.ContainerIdMap, containerId)
	return container
}

//
//////实例化一个node
//func NewNode(nodeId string, address string, port int64, maxMem int64, usedMem int64, client pb.NodeServiceClient, collectionMaxCapacity int64) *Node {
//	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem,
//		UsedMem: usedMem, Client: client, CollectionMaxCapacity: collectionMaxCapacity}
//	node.CollectionMap = make(map[string]*Collection)
//	return node
//}

//
////计算节点压力
//func (node *Node) CalcNodePress() float64 {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	press := float64(node.UserCount) / float64(5)
//	return press
//}
//
////判断节点是否满足container的要求,和这个container的使用人数
//func (node *Node) Satisfy(funcName string, reqMem int64) (bool, int64) {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	cs := node.CollectionMap[funcName]
//	if cs == nil {
//		return false, 0
//	}
//	bool, usedCount := cs.Satisfy(reqMem)
//	return bool, usedCount
//}
//
////获取container
//func (node *Node) Acquire(funcName string, reqMem int64) *Container {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	cs := node.CollectionMap[funcName]
//	if cs == nil {
//		return nil
//	}
//	container := cs.Acquire(reqMem)
//	if container == nil {
//		return nil
//	}
//	node.UserCount++
//	return container
//}
//
////归还container
//func (node *Node) Return(container *Container, actualUseMem int64) {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	cs := node.CollectionMap[container.FuncName]
//	if cs == nil {
//		return
//	}
//	node.UserCount--
//	cs.Return(container, actualUseMem)
//}
//
//////得到node中container的数量
////func (node *Node) GetContainerCount() int {
////	node.lock.Lock()
////	defer node.lock.Unlock()
////	return len(node.Containers)
////}
//
////得到node内存
//func (node *Node) AddContainer(container *Container) {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	cs := node.CollectionMap[container.FuncName]
//	if cs == nil {
//		cs = &Collection{FuncName: container.FuncName, UsedCount: 0, UsedMem: container.UsedMem,
//			MaxUsedMem: container.MaxUsedMem, MaxUsedCount: container.MaxUsedCount, Capacity: node.CollectionMaxCapacity}
//		node.CollectionMap[container.FuncName] = cs
//	}
//	cs.AddContainer(container)
//}
//
////判断是否缺乏某个函数实例
//func (node *Node) Lack(funcName string) bool {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	cs := node.CollectionMap[funcName]
//	if cs == nil {
//		return true
//	}
//	b := cs.Lack()
//	return b
//}
//
//////获得Container
////func (node *Node) GetContainer(funcName string) *Container {
////	node.lock.Lock()
////	defer node.lock.Unlock()
////	container := node.Containers[funcName]
////	return container
////}
//
//func (node *Node) ToString() string {
//	node.lock.Lock()
//	defer node.lock.Unlock()
//	var mapStr string
//
//	for _, cs := range node.CollectionMap {
//
//		mapStr += cs.ToString() + ", "
//	}
//	return mapStr
//}

////得到node中container的数量
//func (node *Node) GetContainerCount() int {
//	node.lock.RLock()
//	defer node.lock.RUnlock()
//	return len(node.Containers)
//}
