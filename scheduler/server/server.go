package server

import (
	"com/aliyun/serverless/scheduler/handler"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func (s Server) AcquireContainer(ctx context.Context, req *pb.AcquireContainerRequest) (*pb.AcquireContainerReply, error) {
	//startTime := time.Now().UnixNano()
	//str, _ := json.Marshal(req)
	//fmt.Println(startTime, string(str))
	if req.AccountId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "account ID cannot be empty")
	}

	if req.FunctionConfig == nil {
		return nil, status.Errorf(codes.InvalidArgument, "function config cannot be nil")
	}

	//acquire handler负责取走container
	var ch = make(chan *pb.AcquireContainerReply)
	handler.AddAcquireContainerToAcquireHandler(req, ch)

	res := <-ch
	if res == nil {
		return &pb.AcquireContainerReply{}, nil
	}
	return res, nil
}

func (s Server) ReturnContainer(ctx context.Context, req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {
	handler.AddReturnContainerToQueue(req)
	return &pb.ReturnContainerReply{}, nil
}
