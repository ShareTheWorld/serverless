package server

import (
	"com/aliyun/serverless/scheduler/core"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
}

func (s Server) AcquireContainer(ctx context.Context, req *pb.AcquireContainerRequest) (*pb.AcquireContainerReply, error) {
	fmt.Println(req)
	if req.AccountId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "account ID cannot be empty")
	}

	if req.FunctionConfig == nil {
		return nil, status.Errorf(codes.InvalidArgument, "function config cannot be nil")
	}

	reply, err := core.AcquireContainer(req.RequestId, req.AccountId, req.FunctionName,
		req.FunctionConfig.TimeoutInMs, req.FunctionConfig.MemoryInBytes, req.FunctionConfig.Handler)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (s Server) ReturnContainer(ctx context.Context, req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {
	fmt.Println(req)
	reply, err := core.ReturnContainer(req.RequestId, req.ContainerId, req.DurationInNanos,
		req.MaxMemoryUsageInBytes, req.ErrorCode, req.ErrorMessage)
	return reply, err
}
