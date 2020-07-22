package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
)

func AcquireContainer(requestId string, accountId string, functionName string, timeoutInMs int64, memoryInBytes int64, handler string) (*pb.AcquireContainerReply, error) {

	reply := new(pb.AcquireContainerReply)
	return reply, nil
}

func ReturnContainer(requestId string, containerId string, durationInNanos int64, maxMemoryUsageInBytes int64, errorCode string, errorMessage string) (*pb.ReturnContainerReply, error) {

	reply := new(pb.ReturnContainerReply)
	return reply, nil
}
