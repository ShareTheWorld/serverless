package main

import (
	pb "com/aliyun/serverless/resourcemanager/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

//type ResourceManagerClient interface {
//	ReserveNode(ctx context.Context, in *ReserveNodeRequest, opts ...grpc.CallOption) (*ReserveNodeReply, error)
//	ReleaseNode(ctx context.Context, in *ReleaseNodeRequest, opts ...grpc.CallOption) (*ReleaseNodeReply, error)
//	GetNodesUsage(ctx context.Context, in *GetNodesUsageRequest, opts ...grpc.CallOption) (*GetNodesUsageReply, error)
//}
//定义一个helloServer并实现约定的接口
type ResourceManagerService struct{}

var id int

func (s ResourceManagerService) ReserveNode(ctx context.Context, in *pb.ReserveNodeRequest) (*pb.ReserveNodeReply, error) {
	fmt.Println("call function: ResourceManager.ReserveNode")
	fmt.Println(in)
	res := new(pb.ReserveNodeReply)
	res.Node = new(pb.NodeDesc)
	id++
	res.Node.Id = fmt.Sprintf("node_id_00%v", id)
	res.Node.Address = "127.0.0.1"
	res.Node.NodeServicePort = 30000
	res.Node.MemoryInBytes = 4 * 1024 * 1024 * 1024
	return res, nil
}

func (s ResourceManagerService) ReleaseNode(ctx context.Context, in *pb.ReleaseNodeRequest) (*pb.ReleaseNodeReply, error) {
	fmt.Println("call function: ResourceManager.ReleaseNode")
	fmt.Println(in)
	res := new(pb.ReleaseNodeReply)
	return res, nil
}

func (s ResourceManagerService) GetNodesUsage(ctx context.Context, in *pb.GetNodesUsageRequest) (*pb.GetNodesUsageReply, error) {
	fmt.Println("call function: ResourceManager.GetNodesUsage")
	fmt.Println(in)
	res := new(pb.GetNodesUsageReply)

	//var arr [3]pb.NodeDesc = [3]pb.NodeDesc{
	//	pb.NodeDesc{Id: "1", Address: "李白", NodeServicePort: 100, MemoryInBytes: 0, ReservedTimeTimestampMs: 0, ReleasedTimeTimestampMs: 0},
	//	pb.NodeDesc{Id: "2", Address: "李白", NodeServicePort: 100, MemoryInBytes: 0, ReservedTimeTimestampMs: 0, ReleasedTimeTimestampMs: 0},
	//	pb.NodeDesc{Id: "3", Address: "李白", NodeServicePort: 100, MemoryInBytes: 0, ReservedTimeTimestampMs: 0, ReleasedTimeTimestampMs: 0},
	//}

	return res, nil
}
func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println(err)
	}

	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterResourceManagerServer(s, ResourceManagerService{})

	fmt.Println("Resource Manager Service Listen on 127.0.0.1:20000")
	s.Serve(listen)
}
