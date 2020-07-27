package client

import (
	resPb "com/aliyun/serverless/resourcemanager/proto"
	global "com/aliyun/serverless/scheduler/utils/groble"
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
		ConnectResourceManagerService(global.ResourceManagerEndpoint)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := resPb.ReserveNodeRequest{RequestId: requestId, AccountId: accountId}
	res, err := resClient.ReserveNode(ctx, &req)
	return res, err
}

//释放节点
func ReleaseNode(requestId string, id string) *resPb.ReleaseNodeReply {
	req := resPb.ReleaseNodeRequest{RequestId: requestId, Id: id}
	res, err := resClient.ReleaseNode(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
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
