package main

import (
	pb "com/aliyun/serverless/test/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

const (
	Address = "127.0.0.1:50052"
)

//定义一个helloServer并实现约定的接口
type HelloService struct{}

func (h HelloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	resp := new(pb.HelloReply)
	resp.Message = "hello " + in.Name + "."
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		fmt.Println(err)
	}

	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterHelloServer(s, HelloService{})

	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
