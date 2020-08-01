package core

import pb "com/aliyun/serverless/nodeservice/proto"

//存放container信息
type Container struct {
	FunName string //函数名字
	Id      string //容器id
	UsedMem int64  //使用内存
}

//存放节点信息
type Node struct {
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
var Nodes = make([]*Node, 0, 100)

//请求表，用于存放所有的请求
var RequestMap = make(map[string]*NC)
