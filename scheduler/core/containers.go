package core

import (
	cmap "github.com/orcaman/concurrent-map"
)

/**
全局的container
*/
var FunCtnMap = cmap.New() //function_name -> ContainerMap (container_id -> ContainerInfo)
//var FuncMapLock sync.RWMutex

//添加container
func AddContainer(container *Container) {
	if container == nil {
		return
	}

	FunCtnMap.SetIfAbsent(container.FuncName, cmap.New())

	//如果Map
	obj, _ := FunCtnMap.Get(container.FuncName)
	containerMap := obj.(cmap.ConcurrentMap) //转为对应的map

	containerMap.Set(container.ContainerId, container)
}

//移除container
func RemoveContainer(container *Container) {
	if container == nil {
		return
	}
	obj, _ := FunCtnMap.Get(container.FuncName)
	if obj == nil {
		return
	}
	containerMap := obj.(cmap.ConcurrentMap)
	containerMap.Remove(container.ContainerId)
}

func GetContainerMap(funcName string) cmap.ConcurrentMap {
	//如果Map
	obj, _ := FunCtnMap.Get(funcName)
	if obj == nil {
		return nil
	}
	containerMap := obj.(cmap.ConcurrentMap) //转为对应的map
	return containerMap
}

//获取内存使用最多的函数
func GetFuncByMaxMem() string {
	var funcName string = ""
	var funcMem int64 = 0
	fns := FunCtnMap.Keys()
	for _, fn := range fns {
		obj1, _ := FunCtnMap.Get(fn)
		if obj1 == nil {
			continue
		}
		ctnMap := obj1.(cmap.ConcurrentMap) //转为对应的map
		ctnIds := ctnMap.Keys()

		var mem int64 = 0
		for _, ctnId := range ctnIds {
			obj2, _ := ctnMap.Get(ctnId)
			if obj2 == nil {
				continue
			}
			ctn := obj2.(*Container)
			mem += ctn.UsageMem
		}

		if funcName == "" {
			funcName = fn
			funcMem = mem
			continue
		}

		if mem > funcMem {
			funcName = fn
			funcMem = mem
		}
	}
	return funcName
}

//获取用户使用最多的函数
func GetFuncByMaxUseCount() string {
	var funcName string = ""
	var funcMem int64 = 0
	fns := FunCtnMap.Keys()
	for _, fn := range fns {
		obj1, _ := FunCtnMap.Get(fn)
		if obj1 == nil {
			continue
		}
		ctnMap := obj1.(cmap.ConcurrentMap) //转为对应的map
		ctnIds := ctnMap.Keys()

		var mem int64 = 0
		for _, ctnId := range ctnIds {
			obj2, _ := ctnMap.Get(ctnId)
			if obj2 == nil {
				continue
			}
			ctn := obj2.(*Container)
			mem += ctn.UsageMem
		}

		if funcName == "" {
			funcName = fn
			funcMem = mem
			continue
		}

		if mem > funcMem {
			funcName = fn
			funcMem = mem
		}
	}
	return funcName
}

//获取cpu使用最多的
func GetFuncByMaxCpu() *Container {
	var container *Container
	var funcCpu float64 = 0
	fns := FunCtnMap.Keys()
	for _, fn := range fns {
		obj1, _ := FunCtnMap.Get(fn)
		if obj1 == nil {
			continue
		}
		ctnMap := obj1.(cmap.ConcurrentMap) //转为对应的map
		ctnIds := ctnMap.Keys()

		var cpu float64 = 0
		var ctn *Container
		for _, ctnId := range ctnIds {
			obj2, _ := ctnMap.Get(ctnId)
			if obj2 == nil {
				continue
			}
			ctn = obj2.(*Container)
			cpu += ctn.CpuUsagePct
		}

		if container == nil {
			container = ctn
			continue
		}

		if cpu > funcCpu {
			container = ctn
		}
	}
	return container
}
