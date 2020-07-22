package client

import (
	nodePb "com/aliyun/serverless/nodeservice/proto"
	resPb "com/aliyun/serverless/resourcemanager/proto"
	global "com/aliyun/serverless/scheduler/utils/groble"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

//resource manager grpc client
var resClient resPb.ResourceManagerClient

// node service grpc client
var nodeClients map[string]nodePb.NodeServiceClient

//连接到资源管理器服务
func ConnectResourceManagerService(endpoint string) {
	fmt.Printf("connect to resource manager service,Address: %s\n", endpoint)
	//连接到grpc服务
	resConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	//defer resConn.Close()

	//初始化客户端
	resClient = resPb.NewResourceManagerClient(resConn)
}

//预定节点
func ReserveNode(requestId string, accountId string) *resPb.ReserveNodeReply {
	req := resPb.ReserveNodeRequest{RequestId: requestId, AccountId: accountId}
	res, err := resClient.ReserveNode(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return res
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

func Test() {
	ConnectResourceManagerService(global.ResourceManagerEndpoint)
	ReserveNode("request_id_0001", "account_id_0001")
	ReleaseNode("request_id_0002", "id_0002")
	GetNodesUsage("request_id_0003")

	c := ConnectNodeService("id_0004", "127.0.0.1", 30000)
	CreateContainer(c, "request_id_0005", "name_0005", "function_name_0005", "handler_0005", 10, 10)
	RemoveContainer(c, "request_id_0006", "container_id_0006")
	GetStats(c, "request_id_0007")
}

//连接到节点服务
func ConnectNodeService(id string, address string, port int64) nodePb.NodeServiceClient {
	//连接到grpc服务
	endpoint := fmt.Sprintf("%s:%d", address, port)
	fmt.Printf("connect to node service,Address: %s\n", endpoint)
	nodeConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}
	//defer nodeConn.Close()

	//初始化客户端
	var nodeClient = nodePb.NewNodeServiceClient(nodeConn)
	return nodeClient
}

//创建容器
func CreateContainer(nodeClient nodePb.NodeServiceClient, requestId string, name string, functionName string, handler string, timeoutInMs int64, memoryInBytes int64) *nodePb.CreateContainerReply {
	req := nodePb.CreateContainerRequest{RequestId: requestId, Name: name, FunctionMeta: &nodePb.FunctionMeta{FunctionName: functionName, Handler: handler, TimeoutInMs: timeoutInMs, MemoryInBytes: memoryInBytes}}
	res, err := nodeClient.CreateContainer(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	return res
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
