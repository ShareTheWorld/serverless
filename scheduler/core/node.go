package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"github.com/pkg/errors"
)

//存放节点信息
type Node struct {
	NodeID     string               //节点id
	Address    string               //节点地址
	Port       int64                //节点端口
	MaxMem     int64                //最大内存
	UsedMem    int64                //使用内存
	Client     pb.NodeServiceClient //节点连接
	Containers *LM                  //保存node里面所有的容器信息
}

func NewNode(nodeId string, address string, port int64, maxMem int64) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem, Containers: NewLM()}
	return node
}

//添加Container,
func (node *Node) AddContainer(container *Container) {
	node.Containers.Add(container.FunName, container)
}

//移除Container，会将container销毁
func (node *Node) RemoveContainer(containerId string) {

}

//租用container，会消耗内存
func (node *Node) RentContainer(container *Container) (*Container, error) {

	c := node.QueryContainer(container.FunName, container.UsedMem)
	if c == nil {
		return nil, errors.New("No Containers available or lack of memory")
	}
	node.UsedMem += c.UsedMem
	return c, nil
}

func (node *Node) ReturnContainer(container *Container) {
	node.UsedMem -= container.UsedMem
}

//查询container，不会消耗内存
func (node *Node) QueryContainer(funcName string, reqMem int64) *Container {
	//如果内存不够就直接返回
	if node.UsedMem+reqMem > node.MaxMem {
		return nil
	}
	container := node.Containers.Get(funcName)
	if container == nil {
		return nil
	}
	return container.(*Container)
}
