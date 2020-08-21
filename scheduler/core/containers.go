package core

import "sync"

/**
全局的container
*/
var FunMap = make(map[string]map[string]*Container)
var FuncMapLock sync.RWMutex

//添加container
func AddContainer(container *Container) {
	if container == nil {
		return
	}

	FuncMapLock.Lock()
	defer FuncMapLock.Unlock()

	m := FunMap[container.FuncName]
	if m == nil {
		m = make(map[string]*Container)
		FunMap[container.FuncName] = m
	}
	m[container.ContainerId] = container
}

//移除container
func RemoveContainer(container *Container) {
	if container == nil {
		return
	}
	m := FunMap[container.FuncName]
	if m == nil {
		return
	}
	delete(m, container.ContainerId)
}
