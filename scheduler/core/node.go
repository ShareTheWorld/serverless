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
	Client  pb.NodeServiceClient `json:"-"` //节点连接

	TotalMem     int64   //总内存
	UsageMem     int64   //使用内存
	AvailableMem int64   //可用内存
	CpuUsagePct  float64 //cpu使用百分比

	UseCount         int64 //当前正在使用的人数
	ConcurrencyCount int64 //并发数量
	Status           int   //1表示正常，0表示不以使用

	ContainerIdMap map[string]*Container `json:"-"` //存放所有的Container K:V=containerId:Container
	FuncNameMap    map[string]bool       `json:"-"` //用于判断函数是否存在了
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
		Status:           1,
		FuncNameMap:      make(map[string]bool),
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
	n.FuncNameMap[container.FuncName] = true
}

//根据containerId移除container  TODO 需要在全局中移除
func (n *Node) RemoveContainer(containerId string) *Container {
	n.lock.Lock()
	defer n.lock.Unlock()

	container := n.ContainerIdMap[containerId]
	delete(n.ContainerIdMap, containerId)
	delete(n.FuncNameMap, container.FuncName)
	return container
}
