package client

import (
	nodePb "com/aliyun/serverless/nodeservice/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
)

// node service grpc client
var nodeClients map[string]nodePb.NodeServiceClient

//连接到节点服务
func ConnectNodeService(id string, address string, port int64) (nodePb.NodeServiceClient, error) {
	//连接到grpc服务
	endpoint := fmt.Sprintf("%s:%d", address, port)
	fmt.Printf("connect to node service,Address: %s\n", endpoint)
	nodeConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	//初始化客户端
	var nodeClient = nodePb.NewNodeServiceClient(nodeConn)
	return nodeClient, nil
}

//创建容器
func CreateContainer(nodeClient nodePb.NodeServiceClient, requestId string, name string, functionName string, handler string, timeoutInMs int64, memoryInBytes int64) (*nodePb.CreateContainerReply, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := nodePb.CreateContainerRequest{RequestId: requestId, Name: name, FunctionMeta: &nodePb.FunctionMeta{FunctionName: functionName, Handler: handler, TimeoutInMs: timeoutInMs, MemoryInBytes: memoryInBytes}}
	res, err := nodeClient.CreateContainer(ctx, &req)
	return res, err
}

//销毁容器
func RemoveContainer(nodeClient nodePb.NodeServiceClient, requestId string, containerId string) *nodePb.RemoveContainerReply {
	req := nodePb.RemoveContainerRequest{RequestId: requestId, ContainerId: containerId}
	res, err := nodeClient.RemoveContainer(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

//获取状态
func GetStats(nodeClient nodePb.NodeServiceClient, requestId string) *nodePb.GetStatsReply {
	req := nodePb.GetStatsRequest{RequestId: requestId}
	res, err := nodeClient.GetStats(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return res
}
