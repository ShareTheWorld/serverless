package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	"sync"
)

type NC struct {
	Node      *Node
	Container *Container
}

//用于存放所有node,使用内存越小的放在越后面
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)
var RequestMapLock sync.Mutex

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

//得到Nodes数量
func GetNodeCount() int {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	return len(nodes)
}

//得到nodes的压力
func GetNodesPress() float64 {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	var totalUserCount = 0 //总的使用数
	for _, n := range nodes {
		totalUserCount += n.UserCount
	}
	press := float64(totalUserCount) / float64(10*len(nodes))
	return press
}

//得到最小使用内存的节点
func GetMinUsedMemNode() *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	var node = nodes[0]
	for _, n := range nodes {
		if n.UsedMem < node.UsedMem {
			node = n
		}
	}
	return node
}

//获取一个使用最少的节点
func GetMinUseNode() *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	var node = nodes[0]
	for _, n := range nodes {
		if n.UserCount < node.UserCount {
			node = n
		}
	}
	return node
}

//得到一个container最少的节点
func GetMinContainerNode() *Node {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	var node = nodes[0]
	for _, n := range nodes {
		if len(n.Containers) < len(node.Containers) {
			node = n
		}
	}
	return node
}

//获取一个node里面的container
func Acquire(req *pb.AcquireContainerRequest) *pb.AcquireContainerReply {
	requestId := req.RequestId
	funcName := req.FunctionName
	reqMem := req.FunctionConfig.MemoryInBytes

	var node *Node
	var container *Container
	NodesLock.RLock()
	//遍历node，找到一个满足要求，且连接数最少的
	for _, n := range nodes {
		c := n.GetContainer(funcName)
		//如果不包含要找的方法
		if c == nil {
			continue
		}

		//如果内存不够
		if n.MaxMem-n.UsedMem < reqMem {
			continue
		}

		//函数存在，且内存足够
		if node == nil || container == nil {
			node, container = n, c
			continue
		}

		//如果n的使用数 > node的使用数
		if n.UserCount > node.UserCount {
			continue
		}

		//这个连接数是最少的，就替换成当前的这个
		node, container = n, c
	}
	NodesLock.RUnlock()

	if node == nil || container == nil {
		return nil
	}

	//修改Node的数据
	node.Acquire(container)

	//记录请求
	RequestMapLock.Lock()
	RequestMap[requestId] = &NC{node, container}
	RequestMapLock.Unlock()

	return &pb.AcquireContainerReply{
		NodeId:          node.NodeID,
		NodeAddress:     node.Address,
		NodeServicePort: node.Port,
		ContainerId:     container.Id,
	}
}

//归还node中的container
func Return(req *pb.ReturnContainerRequest) {
	requestId := req.RequestId

	RequestMapLock.Lock()
	nc := RequestMap[requestId]
	delete(RequestMap, requestId)
	RequestMapLock.Unlock()

	if nc == nil {
		return
	}

	node := nc.Node
	container := nc.Container

	node.Return(container)

}

func PrintNodes(tag string) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	fmt.Printf("****************************%v*******************************\n", tag)
	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		fmt.Printf("No:%v, NodeId:%v, Mem:%v/%v, UserCount:%v, containerCount:%v,  %v\n",
			i, node.NodeID, node.UsedMem/1024/1024,
			node.MaxMem/1024/1024, node.UserCount,
			len(node.Containers), node.Containers)
	}
	fmt.Printf("**************************************************************\n\n")

}
