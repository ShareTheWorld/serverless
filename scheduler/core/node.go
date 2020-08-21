package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	rmpb "com/aliyun/serverless/resourcemanager/proto"
	"sync"
)

//Node结构
//锁：用于控制同步 ,连接信息 ,节点状态 ,本地使用状态 ,Container信息
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

	ContainerIdMap map[string]*Container //存放所有的Container K:V=containerId:Container
}

func NewNode(reply *rmpb.ReserveNodeReply, client pb.NodeServiceClient) *Node {
	node := &Node{
		NodeID:           reply.Node.Id,
		Address:          reply.Node.Address,
		Port:             reply.Node.NodeServicePort,
		Client:           client,
		TotalMem:         reply.Node.MemoryInBytes,
		UsageMem:         1 * 1024 * 1024 * 1024,
		AvailableMem:     3 * 1024 * 1024 * 102,
		CpuUsagePct:      1,
		UseCount:         0,
		ConcurrencyCount: 1,
		ContainerIdMap:   make(map[string]*Container),
	}
	return node
}

//更新node的状态
func (n *Node) UpdateNodeStats(stats *pb.NodeStats) {
	if stats == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	n.TotalMem = stats.TotalMemoryInBytes
	n.UsageMem = stats.MemoryUsageInBytes
	n.AvailableMem = stats.AvailableMemoryInBytes
	n.CpuUsagePct = stats.CpuUsagePct

}

//更新所有container的状态
func (n *Node) UpdateContainer(stats []*pb.ContainerStats) {
	if stats == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	for _, s := range stats {
		if s == nil {
			continue
		}
		container := n.ContainerIdMap[s.ContainerId]
		if container == nil {
			continue
		}
		container.updateContainerStats(s)
	}
}

//添加container
func (n *Node) AddContainer(container *Container) {
	if container == nil {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	n.ContainerIdMap[container.ContainerId] = container
}

//根据containerId移除container  TODO 需要在全局中移除
func (n *Node) RemoveContainer(containerId string) *Container {
	n.lock.Lock()
	defer n.lock.Unlock()

	container := n.ContainerIdMap[containerId]
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
