package core

import "strconv"

var DefaultMaxUsedCount int64 = 2 //Container实例的默认最大连接数
var CollectionMaxCapacity = 2     //集合最大容量

//存放container信息
type Container struct {
	FunName      string //函数名字
	Id           string //容器id
	UsedCount    int64  //使用数量
	UsedMem      int64  //使用内存
	MaxUsedMem   int64  //最大使用内存
	MaxUsedCount int64  //最大使用数量，会根据实际内存去计算
}

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

//向集合中添加一个Container
func (cs *Collection) AddContainer(container *Container) {
	cs.Containers = append(cs.Containers, container)
}

//得到容器数量
func (cs *Collection) Lack() bool {
	return int64(len(cs.Containers)) < cs.Capacity
}

//判断节点是否满足container的要求,和这个collection的使用人数
func (cs *Collection) Satisfy(reqMem int64) (bool, int64) {
	//判断集合中是否有容器
	if len(cs.Containers) <= 0 {
		return false, 0
	}
	//如果集合的使用人数，小于集合的最大使用人数，就数名满足需要
	bool := cs.UsedCount < int64(len(cs.Containers))*cs.MaxUsedCount
	return bool, cs.UsedCount
}

//获取container
func (cs *Collection) Acquire(reqMem int64) *Container {
	//判断集合中是否有容器
	if len(cs.Containers) <= 0 {
		return nil
	}

	//获取一个使用人数最少的容器
	container := cs.Containers[0]
	for _, c := range cs.Containers {
		if c.UsedCount < container.UsedCount {
			container = c
		}
	}

	cs.UsedCount++
	container.UsedCount++
	return container
}

func (cs *Collection) Return(container *Container, actualUseMem int64) {
	cs.UsedCount--
	container.UsedCount--
	if actualUseMem == 0 {
		actualUseMem = 1 * 1024 * 1024
	}
	cs.MaxUsedCount = cs.MaxUsedMem / actualUseMem
	container.MaxUsedCount = container.MaxUsedMem / actualUseMem
	cs.UsedMem = actualUseMem
	container.UsedMem = actualUseMem
}

func (cs *Collection) ToString() string {
	info := "{" + cs.FunName + ", " + str(cs.UsedCount) + "/" + str(int64(len(cs.Containers))*cs.MaxUsedCount) +
		", " + str(cs.UsedMem/1024/1024) + "/" + str(cs.MaxUsedMem/1024/1024) + ", "

	for _, c := range cs.Containers {
		//info += "[" + c.FunName + ", " + str(c.UsedCount) + ", " + str(c.UsedMem/1024/1024) + "], "
		info += "[" + str(c.UsedCount) + "/" + str(c.MaxUsedCount) + ", " + str(c.UsedMem/1024/1024) + "], "
	}
	info += "}"
	return info
}

func str(i int64) string {
	return strconv.FormatInt(i, 10)
}
