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
    
### docker镜像路径
    registry.cn-hangzhou.aliyuncs.com/aliyun_challenge/challenge:10.0.0

### 整体思路
    当一个请求过来的时候（gRPC是启动了一个新的协程），这个协程回去调用hander里面的AcquiereContainer方法获取请求，
    接着回去调用core里面提供的api中的acquire方法获取container，如果获取到就直接返回，如果没有获取到就会触发hander-container创建新的容器
    然后睡眠一毫秒会再次去尝试获取
    
    node有两个协程，一个状态同步的协程，会定时去获取所有node的状态。一个node伸缩协程，他会去计算当前所有node的压力，根据压力动态伸缩node，
    当node扩容过后，会加载一些cpu型的函数到新的node中。（node的cpu由于是两核，最高200%，申请压力设置的是60%）
    
    container的创建，当有新的函数过来的时候会去创建，container的卸载没有实现，原因是测试数据中很多函数是从一开始就到最后就会一只运行，
    虽然有一些卸载的想法，但是只要实现了，下次请求再来需要重新加载，这样rt的的分恐怕很难有好的分数，比如在20分钟内没有那么多的定时调用，而是在20分钟内，新函数的个数达到1000个，
    
    core包，主要是存放node相关的核心数据，是正式node的一个缩影
    
    handler主要是提供业务逻辑处理的函数
    
    其他的一些想法：
    比如container的创建需要五百毫秒，根据阿姆达尔定律定律，这完全应该作为一个重点需要解决的问题，比如函数创建直接加载已经初始化的函数的内存镜像，
    
    比赛环境中的cpu的利用率为什么一直不高，平均最高的时候也达不到50%，而且还是峰值情况。
    
    根据信息论的一些原理，能否考虑在node状态中提供更多的container和node的状态信息，作为scheduler判断函数类型、卸载以及node扩展的依据
    
    测试的数据中就十几个函数，为了提高响应，差不多动用了10-20个node，但是再真实的场景中，小公司的一台服务器就会部署十几个jar包，每个jar包的方法接口在几十个左右，所以觉得测试环境和真实环境还是有很大偏离
    
    阿里的函数计算实现了java、nodeJs、PHP等函数，这些语言定义一个函数加载环境都很耗时，实际上运行的逻辑没有几句代码。有没有可能提供一种新的语言，其特性就是支持动态加载函数，
    这样在每个node中初始化一套环境，用户定义好函数，就可以在这个环境下动态加载那个函数可运行的几KB的编译好的代码，这样可以解决函数加载容器和初始化运行环境耗时的难题，从而提高动态响应速度
    学习一门新语言也就也就一两天时间，如果可以让整个node初始化一套运行环境，然后动态加载一个函数，这可能是函数计算最好的一种方式，而不是靠scheduler去预测然后提前加载函数
    
    
    