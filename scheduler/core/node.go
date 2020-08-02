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
func NewNode(nodeId string, address string, port int64, maxMem int64) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem}
	node.Containers = make(map[string]*Container)
	node.MaxMem -= 128 * 1024 * 1024 //每个节点预留512M的空间，不使用完
	return node
}

//申请使用Node资源
func (node *Node) Acquire(reqMem int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UsedMem += reqMem
	node.UserCount++
	//对node进行排序
	InsertSort(len(nodes)-1, true)
}

//归还资源
func (node *Node) Return(reqMem int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UsedMem -= reqMem
	node.UserCount--
	//对node进行排序
	InsertSort(len(nodes)-1, false)
}

//判断内存是否足够
func (node *Node) RequireMem(reqMem int64) bool {
	node.lock.RLock()
	defer node.lock.RUnlock()
	b := node.MaxMem-node.UsedMem > reqMem
	return b
}

//得到node内存
func (node *Node) GetMem() (int64, int64) {
	node.lock.RLock()
	defer node.lock.RUnlock()
	return node.UsedMem, node.MaxMem
}

//得到node内存
func (node *Node) AddContainer(container *Container) {
	node.lock.RLock()
	defer node.lock.RUnlock()
	node.Containers[container.FunName] = container
}

//获得Container
func (node *Node) GetContainer(funcName string) *Container {
	node.lock.RLock()
	defer node.lock.RUnlock()
	container := node.Containers[funcName]
	return container
}
