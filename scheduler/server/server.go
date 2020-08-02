package server

import (
	"com/aliyun/serverless/scheduler/handler"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type Server struct {
}

var StartMap = make(map[string]int64)
var MidMap = make(map[string]int64)
var lock sync.Mutex

func (s Server) AcquireContainer(ctx context.Context, req *pb.AcquireContainerRequest) (*pb.AcquireContainerReply, error) {
	//startTime := time.Now().UnixNano()
	//str, _ := json.Marshal(req)
	//fmt.Println(startTime, string(str))
	//fmt.Printf("%v\t%v\t%v\t", "acquire", time.Now().UnixNano(), req.RequestId)
	lock.Lock()
	StartMap[req.RequestId] = time.Now().UnixNano()
	lock.Unlock()
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
	lock.Lock()
	MidMap[req.RequestId] = time.Now().UnixNano()
	lock.Unlock()
	if res == nil {
		return &pb.AcquireContainerReply{}, nil
	}
	return res, nil
}

func (s Server) ReturnContainer(ctx context.Context, req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {
	//fmt.Printf("%v\t%v\t%v\t", "return", time.Now().UnixNano(), req.RequestId)
	et := time.Now().UnixNano()
	id := req.RequestId
	lock.Lock()
	mt := MidMap[id]
	st := StartMap[id]
	lock.Unlock()
	fmt.Printf("SL:%v\t,FD:%v\t,RT:%v\t|\tmem:%v\ttime:%v\terr:%v\n", (mt-st)/1000000, (et-mt)/1000000, (et-st)/1000000,
		req.MaxMemoryUsageInBytes/1048576, req.DurationInNanos/1000000, req.ErrorMessage)
	handler.AddReturnContainerToQueue(req)
	return &pb.ReturnContainerReply{}, nil
}
