package main

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

const (
	Address = "127.0.0.1:10600"
)

func main1() {
	//连接到grpc服务
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	//初始化客户端
	c := pb.NewSchedulerClient(conn)

	//调用方法
	reqBody := new(pb.AcquireContainerRequest)
	reqBody.RequestId = "request_id_001"
	reqBody.AccountId = "account_id_0002"
	reqBody.FunctionConfig = new(pb.FunctionConfig)
	res, err := c.AcquireContainer(context.Background(), reqBody)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	reqBody2 := new(pb.ReturnContainerRequest)
	reqBody2.RequestId = "request_id_000a"
	reqBody2.ContainerId = "container_id_000b"

	res2, err := c.ReturnContainer(context.Background(), reqBody2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2)
}
