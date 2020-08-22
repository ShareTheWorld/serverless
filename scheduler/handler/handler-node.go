package handler

import (
	"com/aliyun/serverless/scheduler/client"
	"com/aliyun/serverless/scheduler/core"
	"fmt"
	"sync"
	"time"
)

/*
	node-manager负责探测node资源的使用率，
	当使用率高的时候就去申请资源，
	当使用率低的时候就释放资源
*/
const AccountId = "1317891723692367"      //TODO 线上可能会变化
const MinNodeCount = 5                    //最少节点数量
const MaxNodeCount = 20                   //最大节点数量
const SleepTime = time.Millisecond * 2000 //睡眠时间
const ReserveNodeStep = 2                 //发现node压力过大时，每次申请多少个node

const CpuReservePress = 0.6  //预定node的cpu压力
const CpuReleasePress = 0.25 //释放node的cpu使用率
const MemReservePress = 0.7  //预定node的cpu压力
const MemReleasePress = 0.25 //释放node的cpu使用率

const NodeStatusSyncFrequency = 50 //node状态同步频率每秒多少次

const ActionReserveNode = 1  //预定node
const ActionKeepNode = 0     //保持node
const ActionReleaseNode = -1 //释放node

//const NodeSniffIntervalTime = time.Millisecond * 2000 //Node嗅探间隔时间

//MinNodeCount=a,MaxNodeCount=b
//(0,a)申请资源
//[a,a]只能申请资源
//(a,b)申请或者释放资源
//[b,)只能释放资源

func NodeHandler() {
	go NodeManager() //启动node管理协程

	go SyncNodeStats() //启动状态同步协程
}

//管理node
func NodeManager() {
	for {
		size := core.GetNodeCount()
		//(0,a)不满足最低要求，无条件直接申请资源
		if size < MinNodeCount {
			node := ReserveOneNode()
			core.AddNode(node)
			continue
		}
		//if true { //TODO 只是为了固定容器的个数
		//	return
		//}
		time.Sleep(SleepTime)

		avgMemUsagePct, avgCpuUsagePct := core.CalcNodesPress() //计算节点压力
		action := Action(avgMemUsagePct, avgCpuUsagePct)

		//[a,a]只能申请资源
		if size == MinNodeCount {
			if action == ActionReserveNode {
				DownNodesPress()
			}
			continue
		}

		//(a,b)申请或者释放资源
		if size > MinNodeCount && size < MaxNodeCount {
			if action == ActionReserveNode { //当压力达到0.7就申请一个node
				DownNodesPress()
			} else if action == ActionReleaseNode { //当压力小于0.4就释放一个
				ReleaseOneNode()
			}
			continue
		}

		if size >= MaxNodeCount {
			if action == ActionReleaseNode {
				ReleaseOneNode()
			}
			continue
		}
	}
}

//将mem和cpu平均是用率转化为压力,-1表示需要是否node，0表示不需要做任何操作，1表示需要增加node
func Action(avgMemUsagePct float64, avgCpuUsagePct float64) int {
	//mem和cpu有一个压力过大，就申请node
	if avgMemUsagePct > MemReservePress || avgCpuUsagePct > CpuReservePress {
		return 1
	}
	//mem和cpu两个都压力很小，就是释放node
	if avgMemUsagePct < MemReleasePress && avgCpuUsagePct < CpuReleasePress {
		return -1
	}
	return 0
}

//同步所有node节点的状态
func SyncNodeStats() {
	for {
		nodes := core.GetNodes()
		var wg sync.WaitGroup
		wg.Add(len(nodes))
		for _, node := range nodes {
			go func(n *core.Node) {
				res := client.GetStats(n.Client, "")
				if res != nil { //更新node节点的状态
					n.UpdateNodeStats(res.NodeStats)
					n.UpdateContainer(res.ContainerStatsList)
				}
				wg.Done()
			}(node)
		}
		wg.Wait()
		time.Sleep(time.Millisecond * 1000 / NodeStatusSyncFrequency)
	}
}

//减少node的压力
func DownNodesPress() {
	//每次添加指定步长的node，但是不能超过总量
	for i := 0; i < ReserveNodeStep; i++ {
		size := core.GetNodeCount()
		if size >= MaxNodeCount { //如果node数量已经达到限制了，就什么也不做
			break
		}
		node := ReserveOneNode()
		core.AddNode(node) //必须先添加，否则后面的计算node压力时，统计不到新增节点
		fmt.Println(node)
	}
}

//这个方法需要保证一定要申请一个Node
func ReserveOneNode() *core.Node {
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

		//创建成功node并且连接成功，进行节点添加
		node := core.NewNode(reply, nodeClient)
		et := time.Now().UnixNano()
		fmt.Printf("---- reserve node, time=%v, node:%v \n", (et-st)/1000000, node)
		return node
	}
}

//释放一个Node
func ReleaseOneNode() {
	node := core.RemoveLastNode() //这里从node池中移除了node，就不会再分配给其他节点了
	for i := 0; i < 100; i++ {    //最多等待30秒
		if node.UseCount <= 0 { //说明这个node没有使用者了
			break
		}
		time.Sleep(time.Millisecond * 300)
	}
	client.ReleaseNode("", node.NodeID)
}

//
//func PrintNodeStats() {
//	for {
//		time.Sleep(time.Millisecond * 10000) //没10秒打印一次node状态
//		nodes := core.GetNodes()
//		//******************log*************************
//		//core.PrintNodes("local node status")
//		//******************log*************************
//		fmt.Printf("****************************%v*******************************\n", "remote node stats")
//		for _, n := range nodes {
//			reply := client.GetStats(n.Client, "")
//			jsonStr, err := json.Marshal(reply)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//			fmt.Println(string(jsonStr))
//		}
//		fmt.Printf("**************************************************************\n\n")
//	}
//}
