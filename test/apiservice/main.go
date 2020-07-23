package main

import (
	"bufio"
	pb "com/aliyun/serverless/scheduler/proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var client pb.SchedulerClient

func main() {
	Init()
	test()
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
				client.AcquireContainer(context.Background(), req1)
			}
		} else {
			req2 := new(pb.ReturnContainerRequest)
			if err := json.Unmarshal([]byte(arr[1]), &req2); err == nil {
				fmt.Print(callTime)
				fmt.Print("   ")
				fmt.Println(req2)
				client.ReturnContainer(context.Background(), req2)

			}
		}
	}
	fmt.Println(time.Now())
}
