package client

import (
	resPb "com/aliyun/serverless/resourcemanager/proto"
	"com/aliyun/serverless/scheduler/utils/groble"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

//resource manager grpc client
var resClient resPb.ResourceManagerClient

//连接到资源管理器服务
func ConnectResourceManagerService(endpoint string) {
	fmt.Printf("connect to resource manager service,Address: %s\n", endpoint)
	//连接到grpc服务
	resConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	//defer resConn.Close()

	fmt.Println("connect to resource manager service success")
	//初始化客户端
	resClient = resPb.NewResourceManagerClient(resConn)
}

//预定节点,requestId可以不用传入
func ReserveNode(requestId string, accountId string) (*resPb.ReserveNodeReply, error) {
	if resClient == nil {
		ConnectResourceManagerService(groble.ResourceManagerEndpoint)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req := resPb.ReserveNodeRequest{RequestId: requestId, AccountId: accountId}
	res, err := resClient.ReserveNode(ctx, &req)
	fmt.Printf("reserve node: requestId:%v,accountId:%v,reply:%v,err:%v \n", requestId, accountId, res, err)
	return res, err
}

//释放节点
func ReleaseNode(requestId string, id string) *resPb.ReleaseNodeReply {
	req := resPb.ReleaseNodeRequest{RequestId: requestId, Id: id}
	res, err := resClient.ReleaseNode(context.Background(), &req)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("release node: requestId:%v, id:%v \n", requestId, id)

	return res
}

//获取节点使用情况
func GetNodesUsage(requestId string) *resPb.GetNodesUsageReply {
	req := resPb.GetNodesUsageRequest{RequestId: requestId}
	res, err := resClient.GetNodesUsage(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return res
}
