# serverless

###/docker
    存放常用的命令和docker打包需要的文件，其中的main文件是scheduler编译出来的可运行文件
    编译命令见scheduler
###/nodeservice
    实现了nodeservice的所有接口，供scheduler和/test/apiservice两个模块调用
    
###/resourcemanager
    实现了resourcemanager的接口，供scheduler创建node和释放node
    
###/scheduler
    实现了调度的核型功能，主要实现是按照framework.jpeg这张图的结构去实现的
    编译命令:CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
    
###/test
    apiservice主要是模拟线上调用AcquireContainer和ReturnContainer调用
    其他的文件主要是学习go语言的时候写的demo和基本的数据算法
    
