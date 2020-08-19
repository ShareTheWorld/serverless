package core

//一个collection中装的是相同的函数的实例，这里集合的概念和数学上的不一样
//container 的一个集合
type Collection struct {
	FunName      string //函数名字
	UsedCount    int64  //总的使用数量
	UsedMem      int64  //总的使用内存
	MaxUsedMem   int64  //每个最大使用内存
	MaxUsedCount int64  //每个container的最大使用数量
	Capacity     int64  //Collection的容量
	Containers   []*Container
}
