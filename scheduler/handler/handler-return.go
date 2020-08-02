package handler

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
)

/*
	当有ReturnContainer请求过来的时候，就先将请求放到returnQueue队列中，由于不需要返回值，所以不需要返回
	HandleReturnContainer方法负责处理所有的归还任务
*/

var returnQueue = make(chan *pb.ReturnContainerRequest, 100)

//添加归还容器的请求到队列中
func AddReturnContainerToQueue(req *pb.ReturnContainerRequest) {
	returnQueue <- req
}

//容器归还处理者
func ReturnContainerHandler() {
	fmt.Println("start handle return container")
	for {
		req := <-returnQueue
		Return(req)
		//core.PrintNodes("return")
	}
}
