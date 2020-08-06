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
	"time"
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
	//startTime := time.Now().UnixNano()
	//str, _ := json.Marshal(req)
	//fmt.Println(startTime, string(str))
	//fmt.Printf("%v\t%v\t%v\t", "acquire", time.Now().UnixNano(), req.RequestId)
	st := time.Now().UnixNano()
	if req.AccountId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "account ID cannot be empty")
	}

	if req.FunctionConfig == nil {
		return nil, status.Errorf(codes.InvalidArgument, "function config cannot be nil")
	}

	////container handler负责创建container
	//handler.AddAcquireContainerToContainerHandler(req)

	//acquire handler负责取走container
	var ch = make(chan *pb.AcquireContainerReply)
	handler.AddAcquireContainerToAcquireHandler(req, ch)

	res := <-ch

	mt := time.Now().UnixNano()
	log := Log{st, mt, req.FunctionName, req.FunctionConfig.MemoryInBytes / 1048576, res.NodeId}
	lock.Lock()
	logMap[req.RequestId] = &log
	lock.Unlock()

	if res == nil {
		return &pb.AcquireContainerReply{}, nil
	}

	fmt.Printf("Call Acquire Container, RequestId:%v, NodeId:%v, FN:%v, MEM:%v, SL:%v, reqMem:%v\n",
		req.RequestId,
		log.nodeId,
		log.fn,
		log.mem,
		(log.mt-log.st)/1000000,
		req.FunctionConfig.MemoryInBytes/1048576)

	return res, nil
}

func (s Server) ReturnContainer(ctx context.Context, req *pb.ReturnContainerRequest) (*pb.ReturnContainerReply, error) {
	//fmt.Printf("%v\t%v\t%v\t", "return", time.Now().UnixNano(), req.RequestId)
	et := time.Now().UnixNano()
	id := req.RequestId

	lock.Lock()
	log := logMap[id]
	lock.Unlock()
	fmt.Printf("RequestId:%v, NodeId:%v, FN:%v, MEM:%v, SL:%v, FD:%v, RT:%v, mem:%v, time:%v, err:%v\n",
		req.RequestId,
		log.nodeId,
		log.fn,
		log.mem,
		(log.mt-log.st)/1000000,
		(et-log.mt)/1000000,
		(et-log.st)/1000000,
		req.MaxMemoryUsageInBytes/1048576,
		req.DurationInNanos/1000000, req.ErrorMessage)
	if req.ErrorMessage != "" {
		core.PrintNodes(" error ")
	}
	handler.AddReturnContainerToQueue(req)
	return &pb.ReturnContainerReply{}, nil
}
