package main

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"com/aliyun/serverless/scheduler/server"
	global "com/aliyun/serverless/scheduler/utils/groble"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"

	//"com/aliyun/serverless/scheduler/server"
)

func main() {
	InitResourceMainEndpoint()
	StartSchedulerService()

	//client.Test()

}

//获取环境变量，资源管理器的地址
func InitResourceMainEndpoint() {
	endpoint := os.Getenv("RESOURCE_MANAGER_ENDPOINT")
	if endpoint == "" {
		panic("environment variable RESOURCE_MANAGER_ENDPOINT is not set")
	}
	fmt.Println("get resource manager endpoint is :" + endpoint)
	global.ResourceManagerEndpoint = endpoint
}

//启动Scheduler服务
func StartSchedulerService() {
	fmt.Println("Hello GoLang")
	listen, err := net.Listen("tcp", global.SchedulerServerAddress)
	if err != nil {
		fmt.Println(err)
	}

	//实现gRPC服务
	s := grpc.NewServer()
	//注册HelloServer为客户端提供服务
	pb.RegisterSchedulerServer(s, server.Server{})
	fmt.Println("Listen on " + global.SchedulerServerAddress)
	s.Serve(listen)
}
