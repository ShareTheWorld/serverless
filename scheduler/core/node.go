package core

import pb "com/aliyun/serverless/nodeservice/proto"

type Node struct {
	NodeID  string               //节点id
	Address string               //节点地址
	Port    int64                //节点端口
	MaxMem  int64                //最大内存
	FreeMem int64                //使用内存
	UsedMem int64                //使用内存
	Client  pb.NodeServiceClient //节点连接
}
