package handler

import pb "com/aliyun/serverless/scheduler/proto"

type Record struct {
	AcquireReq *pb.AcquireContainerRequest
	ReturnReq  *pb.ReturnContainerRequest
}

func AcquireContainerReq(req *pb.AcquireContainerRequest) {

}

func ReturnContainerReq(req *pb.ReturnContainerRequest) {

}
