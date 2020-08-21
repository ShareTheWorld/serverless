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
	"sync"
	"time"
)

var client pb.SchedulerClient

const (
	Address = "127.0.0.1:10600"
)

func main() {
	Init()
	//test()
	//testTimer()

	//test1()

	//function_1  	    2.5		5		2610	21.75
	//function_2		5		20		4284
	//function_2 		4		8		1440
	//function_4 		2.5		18		4750
	//function_5 		3		7		1130
	//function_6 		2.5		5		1920
	//function_7		3 		3		653
	//function_8		8		7		822
	//function_9 		49		10		240
	//function_10		50		20		460
	//function_11		270		10		40
	//function_12		290		10		40
	//function_13		50		5		1100
	//function_14		2.5		5		1216
	//function_15		360		10		40
	//function_16

	go call("func_name_001", 128*1024*1024, 1000, 20, 500)
	go call("func_name_002", 128*1024*1024, 1000, 20, 500)
	go call("func_name_003", 256*1024*1024, 500, 10, 2000)
	go call("func_name_004", 256*1024*1024, 300, 20, 2500)
	go call("func_name_005", 128*1024*1024, 10, 10, 3000)
	go call("func_name_006", 256*1024*1024, 3000, 5, 2500)
	go call("func_name_007", 512*1024*1024, 4000, 5, 3000)
	go call("func_name_008", 512*1024*1024, 8000, 10, 8000)
	go call("func_name_009", 256*1024*1024, 3000, 10, 49000)
	go call("func_name_010", 256*1024*1024, 12000, 20, 50000)
	go call("func_name_011", 512*1024*1024, 10, 10, 270000)
	go call("func_name_012", 256*1024*1024, 300, 10, 290000)
	go call("func_name_013", 256*1024*1024, 1500, 5, 50000)
	go call("func_name_014", 256*1024*1024, 100, 5, 2500)
	go call("func_name_015", 512*1024*1024, 1000, 10, 360000)
	go call("func_name_016", 512*1024*1024, 1000, 10, 4000)
	time.Sleep(time.Second * 100000)
}
func call(funcName string, reqMem int64, execTime int64, concurrentCount int, intervalTime int64) {
	for {
		var group sync.WaitGroup
		for i := 0; i < concurrentCount; i++ {
			group.Add(1)
			go callOneFunction(funcName, reqMem, execTime, &group)
		}

		group.Wait()
		pauseTime := 1000 * 1000 * intervalTime
		time.Sleep(time.Duration(pauseTime))
	}
}

func callOneFunction(funcName string, reqMem int64, execTime int64, group *sync.WaitGroup) {
	startTime := time.Now().UnixNano()
	id := uuid.NewV4().String()
	req := pb.AcquireContainerRequest{
		RequestId:    id,
		AccountId:    "1317891723692367",
		FunctionName: funcName,
		FunctionConfig: &pb.FunctionConfig{
			TimeoutInMs:   60000,
			MemoryInBytes: reqMem,
			Handler:       "pre_handler_14",
		},
	}
	client.AcquireContainer(context.Background(), &req)

	//fmt.Println(reply)
	pauseTime := 1000 * 1000 * execTime
	time.Sleep(time.Duration(pauseTime))

	req2 := pb.ReturnContainerRequest{
		RequestId:             id,
		ContainerId:           "3f08d03bba4217a96abce7dc72131035e8d24730862a7",
		DurationInNanos:       pauseTime,
		MaxMemoryUsageInBytes: 7278592,
	}
	client.ReturnContainer(context.Background(), &req2)
	endTime := time.Now().UnixNano()
	fmt.Printf("%v\t%v\n", funcName, (endTime-startTime)/1000/1000)
	group.Done()
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

//线上测试用例
func onlineTest() {
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
		go func() {
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
		}()

	}
	fmt.Println(time.Now())
}
