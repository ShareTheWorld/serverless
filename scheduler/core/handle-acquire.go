package core

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
)

/*
	当有AcquireContainer请求过来的时候，就将请求放到acquireQueue队列中
	1、跟去请求去检索是否用某个node中用需要的container，如果有就直接返回
	2、如果没有需要的container，就检查是否有空闲node，有就去node中创建一个container，然后返回
	3、如果没有空闲node，就等待一段时间去重复1-2或者2步骤
	注意：等待的这段时间，node-manager回去申请node或者handle-return归还了资源/container
*/
type Pkg struct {
	req *pb.AcquireContainerRequest
	ch  chan *pb.AcquireContainerReply
}

var acquireQueue = make(chan Pkg, 100)

//添加请求容器的请求到队列中
func AddAcquireContainerToQueue(req *pb.AcquireContainerRequest, ch chan *pb.AcquireContainerReply) {
	acquireQueue <- Pkg{req, ch}
}

//处理请求容器队列
func HandleAcquireContainer() {
	fmt.Println("start handle acquire container")
	for {
		pkg := <-acquireQueue
		fmt.Println(pkg.req)
		pkg.ch <- &pb.AcquireContainerReply{}
	}
}
