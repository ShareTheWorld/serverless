package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	"fmt"
	"time"
)

/*
	node-manager负责探测node资源的使用率，
	当使用率高的时候就去申请资源，
	当使用率低的时候就释放资源
*/
const ReservePress = 0.7                 //申请压力
const ReleasePress = 0.4                 //释放压力
const AccountId = "1317891723692367"     //TODO 线上可能会变化
const MinNodeCount = 5                   //最少节点数量
const MaxNodeCount = 20                  //最大节点数量
const SleepTime = time.Millisecond * 100 //当没有事干的时候睡眠多少毫秒

//MinNodeCount=a,MaxNodeCount=b
//(0,a)申请资源
//[a,a]只能申请资源
//(a,b)申请或者释放资源
//[b,)只能释放资源

func NodeHandler() {
	for {
		size := core.NodeCount()
		//(0,a)不满足最低要求，无条件直接申请资源
		if size < MinNodeCount {
			node := ReserveOneNode()
			core.AddNode(node)
			fmt.Println(node)
			continue
		}
		press := calcNodePress() //计算节点压力

		//[a,a]只能申请资源
		if size == MinNodeCount { //刚好是最小情况，什么也不做
			if press > ReservePress {
				node := ReserveOneNode()
				core.AddNode(node)
				fmt.Println(node)
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}

		//(a,b)申请或者释放资源
		if size > MinNodeCount && size < MaxNodeCount {
			if press > ReservePress { //当压力达到0.7就申请一个node
				node := ReserveOneNode()
				core.AddNode(node)
				fmt.Println(node)
			} else if press < ReleasePress { //当压力小于0.4就释放一个
				ReleaseOneNode()
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}

		if size >= MaxNodeCount {
			if press < ReleasePress {
				ReleaseOneNode()
			} else {
				time.Sleep(SleepTime)
			}
			continue
		}
	}
}

//计算节点的压力
func calcNodePress() float64 {
	var allNodeTotalMem int64 = 0
	var allNodeUsedMem int64 = 0

	for i := 0; i < core.NodeCount(); i++ {
		node := core.GetNode(i)
		usedMem, maxMem := node.GetMem()
		allNodeTotalMem += maxMem
		allNodeUsedMem += usedMem
	}

	if allNodeTotalMem == 0 {
		return 1
	}
	press := float64(allNodeUsedMem) / float64(allNodeTotalMem)
	return press
}

//这个方法需要保证一定要申请一个Node
func ReserveOneNode() *core.Node {
	core.PrintNodes("reserve node")
	st := time.Now().UnixNano()
	for {
		//预约一个node
		reply, err := client.ReserveNode("", AccountId)
		if err != nil || reply == nil || reply.Node == nil {
			fmt.Println("error ", err)
			time.Sleep(time.Second * 1) //一秒过后再重试
			continue
		}

		//ReservedTimeTimestampMs ReleasedTimeTimestampMs
		nodeClient, err := client.ConnectNodeService(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort)
		if err != nil {
			fmt.Println("error ", err)
			continue
		}
		//requestId := uuid.NewV4().String()
		//statsReply := client.GetStats(nodeClient, requestId)
		//totalMem := statsReply.GetNodeStats().TotalMemoryInBytes
		//usedMem := statsReply.GetNodeStats().MemoryUsageInBytes
		//创建成功node并且连接成功，进行节点添加
		node := core.NewNode(reply.Node.Id, reply.Node.Address, reply.Node.NodeServicePort, reply.Node.MemoryInBytes, 0, nodeClient)
		et := time.Now().UnixNano()
		fmt.Printf("---- reserve node, time=%v, node:%v \n", (et-st)/1000000, node)
		return node
	}
}

func ReleaseOneNode() {

}
