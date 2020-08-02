package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"sync"
)

//存放container信息
type Container struct {
	FunName string //函数名字
	Id      string //容器id
	UsedMem int64  //使用内存
	lock    sync.RWMutex
}

//存放节点信息
type Node struct {
	lock       sync.RWMutex
	NodeID     string                //节点id
	Address    string                //节点地址
	Port       int64                 //节点端口
	MaxMem     int64                 //最大内存
	UsedMem    int64                 //使用内存
	UserCount  int                   //使用者数量
	Client     pb.NodeServiceClient  //节点连接
	Containers map[string]*Container //存放所有的Container
}

type NC struct {
	Node      *Node
	Container *Container
}

//用于存放所有node
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)
var RequestMapLock sync.Mutex

//实例化一个node
func NewNode(nodeId string, address string, port int64, maxMem int64) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem}
	node.Containers = make(map[string]*Container)
	node.MaxMem -= 128 * 1024 * 1024 //每个节点预留512M的空间，不使用完
	return node
}

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
}

//获取第i个位置的节点
func GetNode(i int) *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return nodes[i]
}

//获得nodes的数量
func NodeCount() int {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return len(nodes)
}

//放入一个请求
func PutRequestNC(requestId string, nc *NC) {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	RequestMap[requestId] = nc
}

//移除一个请求
func RemoveRequestNC(requestId string) {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	delete(RequestMap, requestId)
}

//得到请求
func GetRequestNC(requestId string) *NC {
	RequestMapLock.Lock()
	defer RequestMapLock.Unlock()
	nc := RequestMap[requestId]
	return nc
}

//申请使用Node资源
func (node *Node) Acquire(reqMem int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UsedMem += reqMem
	node.UserCount++
}

//归还资源
func (node *Node) Return(reqMem int64) {
	node.lock.Lock()
	defer node.lock.Unlock()
	node.UsedMem -= reqMem
	node.UserCount--
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

//得到容器使用内存大小
func (container *Container) GetUsedMem() int64 {
	container.lock.RLock()
	defer container.lock.RUnlock()
	return container.UsedMem
}

//设置内存使用大小
func (container *Container) SetUsedMem(usedMem int64) {
	container.lock.Lock()
	defer container.lock.Unlock()
	container.UsedMem = usedMem
}
