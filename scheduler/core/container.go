package core

import (
	pb "com/aliyun/serverless/nodeservice/proto"
	"sync"
)

//表示一个函数实例
//存放container信息
type Container struct {
	lock sync.RWMutex `json:"-"`

	ContainerId string  `json:"-"` //容器id
	TotalMem    int64   //容器总内存
	UsageMem    int64   //容器使用内存
	CpuUsagePct float64 //容器使用百分比

	FuncName      string //函数名字
	Handler       string
	TimeoutInMs   int64
	MemoryInBytes int64

	UseCount         int64 //使用数量
	ConcurrencyCount int64 //支持并发数量

	Node *Node `json:"-"` //所属node
}

func NewContainer() {

}

func (c *Container) updateContainerStats(stats *pb.ContainerStats) {
	c.TotalMem = stats.TotalMemoryInBytes
	c.UsageMem = stats.MemoryUsageInBytes
	c.CpuUsagePct = stats.CpuUsagePct

	//如果cpu使用大于30%，那么就把这个函数定义为cpu型函数
	if stats.CpuUsagePct > 35 {
		CpuFunc.Set(c.FuncName, c)
	}

}
