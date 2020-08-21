package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"sync"
)

//表示一个函数实例
//存放container信息
type Container struct {
	lock sync.RWMutex

	ContainerId string  //容器id
	TotalMem    int64   //容器总内存
	UsageMem    int64   //容器使用内存
	CpuUsagePct float64 //容器使用百分比

	FuncName         string //函数名字
	UseCount         int64  //使用数量
	ConcurrencyCount int64  //支持并发数量

	Node *Node //所属node
}

func NewContainer() {

}

func (c *Container) updateContainerStats(stats *pb.ContainerStats) {
	c.TotalMem = stats.TotalMemoryInBytes
	c.UsageMem = stats.MemoryUsageInBytes
	c.CpuUsagePct = stats.CpuUsagePct
}
