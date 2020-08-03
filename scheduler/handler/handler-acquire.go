package handler

import (
	pb "com/aliyun/serverless/scheduler/proto"
	"fmt"
	"time"
)

/*
	当有AcquireContainer请求过来的时候，就将请求放到acquireQueue队列中
	1、获取一个请求，去检索某个node中是否有需要的container，如果有就直接返回
	2、如果没有需要的container，就检查是否有空闲node，有就去node中创建一个container，然后将node放入到队列末尾
	3、如果没有空闲node，将节点放入到队列后面，
	4、就等待一段时间去重复1-2或者2步骤
	注意：等待的这段时间，node-manager回去申请node或者handle-return归还了资源/container
*/
type Pkg struct {
	req *pb.AcquireContainerRequest
	ch  chan *pb.AcquireContainerReply
}

var acquireQueue = make(chan Pkg, 10000)

//添加请求容器的请求到队列中
func AddAcquireContainerToAcquireHandler(req *pb.AcquireContainerRequest, ch chan *pb.AcquireContainerReply) {
	acquireQueue <- Pkg{req, ch}
}

//请求失败次数，如果连续失败超过一定次数，就会等待一定时间再处理
var RepeatAcquireFailCount = 0

//容器请求处理者
func AcquireContainerHandler() {
	fmt.Println("start handle acquire container")
	for {
		pkg := <-acquireQueue

		req := pkg.req
		ch := pkg.ch

		res := Acquire(req)
		if res != nil {
			ch <- res
			RepeatAcquireFailCount = 0
			//core.PrintNodes("acquire")
			continue
		}

		//如果没有请求到，就将请求放入到队列后面重新排队
		acquireQueue <- pkg
		RepeatAcquireFailCount++

		//当失败次数大于队列的长度时，就暂停一定时间，避免空耗cpu
		if RepeatAcquireFailCount > len(acquireQueue) { //代表队列循环了(len(acquireQueue)+1)个依然没有成功获取
			RepeatAcquireFailCount = 0
			time.Sleep(time.Millisecond * 1)
		}
	}
}
