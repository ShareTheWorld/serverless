package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"strconv"
	"sync"
)

//存放节点信息
type Node struct {
	lock       sync.Mutex
	NodeID     string                //节点id TODO sync
	Address    string                //节点地址 TODO sync
	Port       int64                 //节点端口 TODO sync
	MaxMem     int64                 //最大内存
	UsedMem    int64                 //使用内存
	UserCount  int                   //使用者数量 TODO sync
	Client     pb.NodeServiceClient  //节点连接 TODO sync
	Containers map[string]*Container //存放所有的Container
}

////实例化一个node
func NewNode(nodeId string, address string, port int64, maxMem int64, usedMem int64, client pb.NodeServiceClient) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem, UsedMem: usedMem, Client: client}
	node.Containers = make(map[string]*Container)
	return node
}

//计算节点压力
func (node *Node) CalcNodePress() float64 {
	node.lock.Lock()
	defer node.lock.Unlock()
	press := float64(node.UserCount) / float64(5)
	return press
}

//判断节点是否满足container的要求,和这个container的使用人数
func (node *Node) Satisfy(funcName string, reqMem int64) (bool, int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	container := node.Containers[funcName]
	if container == nil {
		return false, 0
	}
	//但前使用人数必须小于最大使用人数
	bool := node.Containers[funcName].UsedCount < node.Containers[funcName].MaxUsedCount
	return bool, node.Containers[funcName].UsedCount
}

//获取container
func (node *Node) Acquire(funcName string, reqMem int64) *Container {
	node.lock.Lock()
	defer node.lock.Unlock()
	container := node.Containers[funcName]
	node.UserCount++
	container.UsedCount++ //增加一个使用人数
	return container
}

//归还container
func (node *Node) Return(container *Container, actualUseMem int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UserCount--
	container.UsedCount-- //减少使用人数
	container.MaxUsedCount = container.MaxUsedMem / actualUseMem
}

//得到node中container的数量
func (node *Node) GetContainerCount() int {
	node.lock.Lock()
	defer node.lock.Unlock()
	return len(node.Containers)
}

//得到node内存
func (node *Node) AddContainer(container *Container) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.Containers[container.FunName] = container
}

//获得Container
func (node *Node) GetContainer(funcName string) *Container {
	node.lock.Lock()
	defer node.lock.Unlock()
	container := node.Containers[funcName]
	return container
}

func (node *Node) ToString() string {
	node.lock.Lock()
	defer node.lock.Unlock()
	var mapStr string
	for _, v := range node.Containers {
		mapStr += v.FunName + " " +
			strconv.FormatInt(v.UsedCount, 10) + "/" + strconv.FormatInt(v.MaxUsedCount, 10) + " " +
			strconv.FormatInt(v.MaxUsedMem/1024/1024, 10) + ", "
	}
	return mapStr
}

////得到node中container的数量
//func (node *Node) GetContainerCount() int {
//	node.lock.RLock()
//	defer node.lock.RUnlock()
//	return len(node.Containers)
//}
