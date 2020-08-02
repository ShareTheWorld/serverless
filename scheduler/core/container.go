package core

import "sync"

//存放container信息
type Container struct {
	FunName string //函数名字
	Id      string //容器id
	UsedMem int64  //使用内存
	lock    sync.RWMutex
}

//得到容器使用内存大小
func (container *Container) GetUsedMem() int64 {
	container.lock.RLock()
	defer container.lock.RUnlock()
	return container.UsedMem
}

//设置内存使用大小
func (container *Container) SetUsedMem(usedMem int64) {
	container.lock.Lock()
	defer container.lock.Unlock()
	container.UsedMem = usedMem
}
