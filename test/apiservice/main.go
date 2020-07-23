package main

import (
	"bufio"
	pb "com/aliyun/serverless/scheduler/proto"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

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
		fmt.Print(i)
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
			e2 := time.Now().UnixNano()
			s2 := (e2 - s1) / 1000
			preTime = callTime
			fmt.Print(diffTime / 1000)
			fmt.Print("  ")
			fmt.Println(s2)
		}
		bool := strings.Contains(arr[1], "function_config")
		if bool {
			req1 := new(pb.AcquireContainerRequest)
			req1.FunctionConfig = new(pb.FunctionConfig)
			if err := json.Unmarshal([]byte(arr[1]), &req1); err == nil {
				fmt.Print(callTime)
				fmt.Print("   ")
				fmt.Println(req1)
			}
		} else {
			req2 := new(pb.ReturnContainerRequest)
			if err := json.Unmarshal([]byte(arr[1]), &req2); err == nil {
				fmt.Print(callTime)
				fmt.Print("   ")
				fmt.Println(req2)
			}
		}
		//fmt.Println(str)
	}
	fmt.Println(time.Now())
}

type Req1 struct {
	CallTime       int64        `json:"call_time"`
	RequestId      string       `json:"request_id"`
	AccountId      string       `json:"account_id"`
	FunctionName   string       `json:"function_name"`
	FunctionConfig FunctionName `json:"function_config"`
}
type FunctionName struct {
	TimeoutInMs   int64  `json:"timeout_in_ms"`
	MemoryInBytes int64  `json:"memory_in_bytes"`
	Handler       string `json:"handler"`
}

type Req2 struct {
	CallTime              int64  `json:"call_time"`
	RequestId             string `json:"request_id"`
	ContainerId           string `json:"container_id"`
	DurationInNanos       int64  `json:"duration_in_nanos"`
	MaxMemoryUsageInBytes int64  `json:"max_memory_usage_in_bytes"`
}
