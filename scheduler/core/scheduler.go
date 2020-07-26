package core

import "github.com/pkg/errors"

//存放是Node和Container的关系
type NC struct {
	Node      *Node      //租用的那个node
	Container *Container //租用的那个Container
}

//存放所有的Node，kv=nodeId:node
var nodes = NewLM()

//存放所有的租借信息，(因为归还的时候是更具请求id来归还的)kv=requestId,NC
var ncs = NewLM()

//查询node和container，如果有多个node中存在container，那么就随机选择一个
func QueryNodeAndContainer(funcName string, reqMem int64) (*Node, *Container) {
	//遍历node, 查询container实例子
	for e := nodes.L.Front(); nil != e; e = e.Next() {
		node := e.Value.(*Node)
		container := node.GetContainer(funcName, reqMem)
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
	nodes.Add(node.NodeID, node)
}

//移除节点
func RemoveNode(nodeId string) {
	nodes.Remove(nodeId)
}

//归还container，只是减少使用的内存
func ReturnNC(requestId string) {
	rent := ncs.Get(requestId).(*NC)
	if rent == nil { //没有租借信息，就直接返回
		ncs.Remove(requestId)
		return
	}

	node := rent.Node
	container := rent.Container
	node.UsedMem -= container.UsedMem

	ncs.Remove(requestId)
}

//租用Container，会消耗cpu和内存
func RentNC(requestId string, node *Node, container *Container) (*Container, error) {
	if node.UsedMem+container.UsedMem > node.MaxMem {
		return nil, errors.New("The lack of memory")
	}

	//先去查询container
	c := node.Containers.Get(container.FunName).(*Container)
	if c == nil {
		return nil, errors.New("No Containers available")
	}

	node.UsedMem += c.UsedMem
	rent := &NC{Node: node, Container: c}
	ncs.Add(requestId, rent)

	return c, nil
}
