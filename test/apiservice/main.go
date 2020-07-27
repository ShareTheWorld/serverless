package main

import (
	"bufio"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var client pb.SchedulerClient
const (
	Address = "127.0.0.1:10600"
)

func main() {
	Init()
	test()
	//testTimer()

	//test1()
}

func Init() {
	//连接到grpc服务
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}

	//初始化客户端
	client = pb.NewSchedulerClient(conn)
}

func test1() {
	id := uuid.NewV4().String()
	req := pb.AcquireContainerRequest{
		RequestId:    id,
		AccountId:    "1317891723692367",
		FunctionName: "pre_function_15",
		FunctionConfig: &pb.FunctionConfig{
			TimeoutInMs:   60000,
			MemoryInBytes: 536870912,
			Handler:       "pre_handler_15",
		},
	}
	reply, _ := client.AcquireContainer(context.Background(), &req)
	fmt.Println(reply)
	req2 := pb.ReturnContainerRequest{
		RequestId:             id,
		ContainerId:           "3f08d03bba4217a96abce7dc72131035e8d24730862a7",
		DurationInNanos:       1005291237,
		MaxMemoryUsageInBytes: 7278592,
	}
	client.ReturnContainer(context.Background(), &req2)
}

//测试定时函数用例
func testTimer() {
	for {
		time.Sleep(time.Millisecond * 5000)
		req := pb.AcquireContainerRequest{
			RequestId:    "03decb9a-5e32-407e-9c8f-2a1390c5feb",
			AccountId:    "1317891723692367",
			FunctionName: "pre_function_15",
			FunctionConfig: &pb.FunctionConfig{
				TimeoutInMs:   60000,
				MemoryInBytes: 536870912,
				Handler:       "pre_handler_15",
			},
		}
		client.AcquireContainer(context.Background(), &req)
		time.Sleep(time.Millisecond * 500)

		req2 := pb.ReturnContainerRequest{
			RequestId:             "03decb9a-5e32-407e-9c8f-2a1390c5feb",
			ContainerId:           "3f08d03bba4217a96abce7dc72131035e8d24730862a7",
			DurationInNanos:       1005291237,
			MaxMemoryUsageInBytes: 7278592,
		}
		client.ReturnContainer(context.Background(), &req2)
	}
}

//测试内存用例
func testMemoryIntensive() {

}

//测试cpu用例
func testCpuIntensive() {

}

//线上测试用例
func test() {
	// 读取一个文件的内容
	file, err := os.Open("/Users/fht/Desktop/serverless/api-service-function-call.txt")
	if err != nil {
		fmt.Println("open file err:", err.Error())
		return
	}

	// 处理结束后关闭文件
	defer file.Close()

	// 使用bufio读取
	r := bufio.NewReader(file)
	var preTime int64 = 0
	fmt.Println(time.Now())
	i := 0
	for {
		i++
		fmt.Printf("第%v行 ", i)
		fmt.Print("   ")
		data, _, err := r.ReadLine()
		// 读取到末尾退出
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err", err.Error())
			break
		}
		// 打印出内容
		str := string(data)
		arr := strings.Split(str, "   ")
		callTime, err := strconv.ParseInt(arr[0], 10, 64)
		if preTime == 0 {
			preTime = callTime
		} else {
			diffTime := callTime - preTime
			s1 := time.Now().UnixNano()
			if diffTime > 1000 {
				time.Sleep(time.Nanosecond * time.Duration(diffTime))
			}
			preTime = callTime
			fmt.Printf("diff=%vms real=%vms ", diffTime/1000000, (time.Now().UnixNano()-s1)/1000000)
			if diffTime > 1000000000 { //如果等待时间大于1秒
				fmt.Print("\n\n\n")
			}
		}
		bool := strings.Contains(arr[1], "function_config")
		if bool {
			req1 := new(pb.AcquireContainerRequest)
			req1.FunctionConfig = new(pb.FunctionConfig)
			if err := json.Unmarshal([]byte(arr[1]), &req1); err == nil {
				fmt.Print(callTime)
				fmt.Print("   ")
				fmt.Println(req1)
				reply, err := client.AcquireContainer(context.Background(), req1)
				fmt.Println(reply)
				fmt.Println(err)
			}
		} else {
			req2 := new(pb.ReturnContainerRequest)
			if err := json.Unmarshal([]byte(arr[1]), &req2); err == nil {
				fmt.Print(callTime)
				fmt.Print("   ")
				fmt.Println(req2)
				reply, err := client.ReturnContainer(context.Background(), req2)
				fmt.Println(reply)
				fmt.Println(err)
			}
		}
	}
	fmt.Println(time.Now())
}
