package core

import (
	cmap "github.com/orcaman/concurrent-map"
)

/**
全局的container
*/
var FunMap = cmap.New() //function_name -> ContainerMap (container_id -> ContainerInfo)
//var FuncMapLock sync.RWMutex

//添加container
func AddContainer(container *Container) {
	if container == nil {
		return
	}

	FunMap.SetIfAbsent(container.FuncName, cmap.New())

	//如果Map
	obj, _ := FunMap.Get(container.FuncName)
	containerMap := obj.(cmap.ConcurrentMap) //转为对应的map

	containerMap.Set(container.ContainerId, container)
}

//移除container
func RemoveContainer(container *Container) {
	if container == nil {
		return
	}
	obj, _ := FunMap.Get(container.FuncName)
	if obj == nil {
		return
	}
	containerMap := obj.(cmap.ConcurrentMap)
	containerMap.Remove(container.ContainerId)
}

func GetContainerMap(funcName string) cmap.ConcurrentMap {
	//如果Map
	obj, _ := FunMap.Get(funcName)
	if obj == nil {
		return nil
	}
	containerMap := obj.(cmap.ConcurrentMap) //转为对应的map
	return containerMap
}
