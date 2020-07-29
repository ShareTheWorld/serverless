package core

import (
	"sync"
)

//存放是Node和Container的关系
type NC struct {
	Node      *Node      //租用的那个node
	Container *Container //租用的那个Container
}

//存放所有的Node，kv=nodeId:node
var nodes = NewLM()

var lock sync.RWMutex

//存放所有的租借信息，(因为归还的时候是更具请求id来归还的)kv=requestId,NC
var ncs = NewLM()

//查询node和container，如果有多个node中存在container，那么就随机选择一个
func QueryNodeAndContainer(funcName string, reqMem int64) (*Node, *Container) {
	lock.RLock()
	defer lock.RUnlock()
	//遍历node, 查询container实例子
	for e := nodes.L.Front(); nil != e; e = e.Next() {
		node := e.Value.(*Node)
		container := node.QueryContainer(funcName, reqMem)
		if container != nil {
			return node, container
		}
	}

	//找到一个可以装载这个容器的Node
	for e := nodes.L.Front(); nil != e; e = e.Next() {
		node := e.Value.(*Node)
		if node.MaxMem-node.UsedMem > reqMem { //如果内存足够
			return node, nil
		}
	}
	return nil, nil
}

//添加节点
func AddNode(node *Node) {
	lock.Lock()
	defer lock.Unlock()
	nodes.Add(node.NodeID, node)

}

//移除节点
func RemoveNode(nodeId string) {
	lock.Lock()
	defer lock.Unlock()
	nodes.Remove(nodeId)
}

//归还container，只是减少使用的内存
func ReturnNC(requestId string) {
	lock.Lock()
	defer lock.Unlock()

	nc := ncs.Get(requestId)
	if nc == nil {
		return
	}
	rent := nc.(*NC)
	if rent == nil { //没有租借信息，就直接返回
		ncs.Remove(requestId)
		return
	}

	node := rent.Node
	container := rent.Container
	if node != nil && container != nil {
		node.ReturnContainer(container)
	}
	ncs.Remove(requestId)
}

//添加一个NC
func AddNC(node *Node, container *Container) {
	lock.Lock()
	defer lock.Unlock()
	node.AddContainer(container)
}

//租用Container，会消耗cpu和内存
func RentNC(requestId string, node *Node, container *Container) (*Container, error) {
	lock.Lock()
	defer lock.Unlock()
	c, err := node.RentContainer(container)
	if err != nil {
		return nil, err
	}

	rent := &NC{Node: node, Container: c}
	ncs.Add(requestId, rent)

	return c, nil
}
