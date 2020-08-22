package server

import (
	"com/aliyun/serverless/scheduler/core"
	//"com/aliyun/serverless/scheduler/core"
	"com/aliyun/serverless/scheduler/handler"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"fmt"
	//"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type Server struct {
}

var logMap = make(map[string]*Log)
var lock sync.Mutex

type Log struct {
	st     int64
	mt     int64
	fn     string
	mem    int64
	nodeId string
}

func (s Server) AcquireContainer(ctx context.Context, req *pb.AcquireContainerRequest) (*pb.AcquireContainerReply, error) {
	//st := time.Now().UnixNano()
	req.FunctionConfig.MemoryInBytes = 4 * 1024 * 1024 * 1024
	if req.AccountId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "account ID cannot be empty")
	}

	if req.FunctionConfig == nil {
		return nil, status.Errorf(codes.InvalidArgument, "function config cannot be nil")
	}

	res := handler.AcquireContainer(req)

	if res == nil {
		return &pb.AcquireContainerReply{}, nil
	}
	//et := time.Now().UnixNano()
	//fmt.Println((et - st) / 1000 / 1000)
	return res, nil
}

var count = 0

func (s Server) ReturnContainer(ctx context.Context, req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {

	handler.ReturnContainer(req)
	count++
	if count%1000 == 0 {
		core.PrintNodes(" timer ")
	}
	if req.ErrorMessage != "" {
		fmt.Println(req.ErrorMessage)
		core.PrintNodes(" error ")
	}
	return &pb.ReturnContainerReply{}, nil
}
