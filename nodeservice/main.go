package main

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"net"
)

type NodeService struct {
}

//预定
func (s NodeService) Reserve(ctx context.Context, in *pb.ReserveRequest) (*pb.ReserveReply, error) {
	fmt.Println("call function: NodeService.Reserve")
	fmt.Println(in)
	res := new(pb.ReserveReply)
	return res, nil
}

//创建容器
func (s NodeService) CreateContainer(ctx context.Context, in *pb.CreateContainerRequest) (*pb.CreateContainerReply, error) {
	fmt.Println("call function: NodeService.CreateContainer")
	fmt.Println(in)
	res := new(pb.CreateContainerReply)
	res.ContainerId = uuid.NewV4().String()
	return res, nil
}

//销毁容器
func (s NodeService) RemoveContainer(ctx context.Context, in *pb.RemoveContainerRequest) (*pb.RemoveContainerReply, error) {
	fmt.Println("call function: NodeService.RemoveContainer")
	fmt.Println(in)
	res := new(pb.RemoveContainerReply)
	return res, nil
}

//调用函数
func (s NodeService) InvokeFunction(in *pb.InvokeFunctionRequest, out pb.NodeService_InvokeFunctionServer) error {
	fmt.Println("call function: NodeService.InvokeFunction")
	fmt.Println(in)
	return nil
}

//得到容器状态
func (s NodeService) GetStats(ctx context.Context, in *pb.GetStatsRequest) (*pb.GetStatsReply, error) {
	fmt.Println("call function: NodeService.GetStats")
	fmt.Println(in)
	res := new(pb.GetStatsReply)
	return res, nil
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println(err)
	}

	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterNodeServiceServer(s, NodeService{})

	fmt.Println("Node Service Listen on 127.0.0.1:30000")
	s.Serve(listen)
}
