package main

import (
	pb "com/aliyun/serverless/test/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

const (
	Address = "127.0.0.1:50052"
)

func main() {
	//连接到grpc服务
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	//初始化客户端
	c := pb.NewHelloClient(conn)

	//调用方法
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "golang_grpc"
	res, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Message)

}
