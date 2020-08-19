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
    


###framework.jpeg 设计的开发架构图，主要分为六大模块
    NodePool:是数据核型模块，主要用于存放对应的数据结构，和真实环境中会保持数据的对应，通过api对外提供接口，其他模块通过api操作数据
    AcquireHandler:请求处理者，进来的请求会放入到请求表中作为记录，然后将请求放如acquire-queue，如果获取到了就返回，如果没有获取到就放入队列尾部
    ReturnHandler:返回请求处理者，会将函数的运行信息写入请求表中，然后归还使用完的container
    NodeHandler:负责node的管理，以及同步node节点的状态，根据node使用率进行动态的申请和释放node，做node的伸缩
    ContainerHandler:负责container的管理，对container进行多实例创建和选择不同的节点进行创建，以及对container的销毁
    RequestTable和Scheduler:scheduler会根据最近的请求记录和归还信息去对未来进行预测，然后指挥node和container进行提前准备工作
    
    注意：node里面包含多个collection，每个collection是同一个函数的实例集合(和数学上的概念不一样，这里的一个集合里面的元素是重复的，只能是同一个函数的实例）
    因为一个函数的使用内存大小是又用户定义的，cpu的使用也是又内存按比例分配的，所以一个node中需要一个需要一个collection去装同一个函数实例，以更多的使用cpu资源

###最后
    整体思路是按照上面去实现的，因为比赛环境不同于正式环境，有些地方完全按照设计架构去实现，很难取得好的分数，所以大体上是这个结构，细节上很多地方会和比赛环境耦合
