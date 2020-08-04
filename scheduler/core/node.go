package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"sync"
)

//存放节点信息
type Node struct {
	lock       sync.RWMutex
	NodeID     string                //节点id TODO sync
	Address    string                //节点地址 TODO sync
	Port       int64                 //节点端口 TODO sync
	MaxMem     int64                 //最大内存
	UsedMem    int64                 //使用内存
	UserCount  int                   //使用者数量 TODO sync
	Client     pb.NodeServiceClient  //节点连接 TODO sync
	Containers map[string]*Container //存放所有的Container
}

//实例化一个node
func NewNode(nodeId string, address string, port int64, maxMem int64, usedMem int64, client pb.NodeServiceClient) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem, UsedMem: usedMem, Client: client}
	node.Containers = make(map[string]*Container)
	//node.MaxMem -= 512 * 1024 * 1024 //每个节点预留512M的空间，不使用完
	return node
}

////申请使用Node资源
func (node *Node) Acquire(container *Container) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UserCount++
	if container.UsedCount > 0 { //如果有人正在使用
		node.UsedMem += container.UsedMem
	}
	container.UsedCount++
}

//归还资源
func (node *Node) Return(container *Container) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UserCount--

	if container.UsedCount > 0 { //如果还有人在使用
		node.UsedMem -= container.UsedMem
		return
	}
	container.UsedCount--
}

//
////判断内存是否足够
//func (node *Node) RequireMem(reqMem int64) bool {
//	node.lock.RLock()
//	defer node.lock.RUnlock()
//	b := node.MaxMem-node.UsedMem > reqMem
//	return b
//}

//得到node内存
func (node *Node) GetMem() (int64, int64) {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return node.UsedMem, node.MaxMem
}

//得到node内存
func (node *Node) AddContainer(container *Container) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.Containers[container.FunName] = container
	//创建了一个实例，就减少一点容器所占用的空间
	node.UsedMem += container.UsedMem //添加一个container就会消耗这么多内存
}

//获得Container
func (node *Node) GetContainer(funcName string) *Container {
	node.lock.RLock()
	defer node.lock.RUnlock()
	container := node.Containers[funcName]
	return container
}

//得到node中container的数量
func (node *Node) GetContainerCount() int {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return len(node.Containers)
}
