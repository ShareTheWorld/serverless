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
var NodesLock sync.Mutex

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)

//var RequestMapLock sync.Mutex

//添加一个Node
func AddNode(node *Node) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	nodes = append(nodes, node)
}

//获取第i个位置的节点
func GetNode(i int) *Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	return nodes[i]
}

//得到Nodes数量
func GetNodeCount() int {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	return len(nodes)
}

//得到nodes的压力
func GetNodesPress() float64 {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	var totalUserCount = 0 //总的使用数
	for _, n := range nodes {
		totalUserCount += n.UserCount
	}
	press := float64(totalUserCount) / float64(10*len(nodes))
	return press
}

//得到最小使用内存的节点
func GetMinUsedMemNode() *Node {
	NodesLock.Lock()
	defer NodesLock.Unlock()
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
	NodesLock.Lock()
	defer NodesLock.Unlock()
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
	NodesLock.Lock()
	defer NodesLock.Unlock()
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
	NodesLock.Lock()
	defer NodesLock.Unlock()

	requestId := req.RequestId
	funcName := req.FunctionName
	reqMem := req.FunctionConfig.MemoryInBytes

	var node *Node
	var container *Container

	//遍历node，找到一个满足要求，且连接数最少的
	for _, n := range nodes {
		c := n.Containers[funcName]
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

	if node == nil || container == nil {
		return nil
	}

	//修改Node的数据
	node.Acquire(container)

	//记录请求
	RequestMap[requestId] = &NC{node, container}

	return &pb.AcquireContainerReply{
		NodeId:          node.NodeID,
		NodeAddress:     node.Address,
		NodeServicePort: node.Port,
		ContainerId:     container.Id,
	}
}

//减少node容量
func RequireMem(node *Node, reqMem int64) bool {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	//如果内存不够
	if node.MaxMem-node.UsedMem < reqMem {
		return false
	}
	node.UsedMem += reqMem
	return true
}

//归还node中的container
func Return(req *pb.ReturnContainerRequest) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	requestId := req.RequestId

	nc := RequestMap[requestId]
	delete(RequestMap, requestId)

	if nc == nil {
		return
	}

	node := nc.Node
	container := nc.Container

	node.Return(container)

}

//得到node内存
func AddContainer(node *Node, container *Container) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	node.Containers[container.FunName] = container
	//创建了一个实例，就减少一点容器所占用的空间，申请之前就已经减少了，所以这里不再减少
}

//获得Container
func GetContainer(node *Node, funcName string) *Container {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	container := node.Containers[funcName]
	return container
}

//得到node中container的数量
func GetContainerCount(node *Node) int {
	NodesLock.Lock()
	defer NodesLock.Unlock()
	return len(node.Containers)
}

func PrintNodes(tag string) {
	NodesLock.Lock()
	defer NodesLock.Unlock()
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
