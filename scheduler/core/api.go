package core

/*
	提供对外的接口
	Acquire: 获取想要个的container
	Return: 归还container
*/

//获取一个container
func Acquire(funcName string) *Container {
	var container *Container

	Lock.Lock()
	defer Lock.Unlock()

	m := FunMap[funcName]
	if m == nil { //说明没有这个函数
		return nil
	}

	//挑选一个最优的container
	for _, c := range m {
		if c.UseCount >= c.ConcurrencyCount {
			continue
		}

		if container == nil {
			container = c
			continue
		}

		if c.UseCount < container.UseCount {
			container = c
			continue
		}

		if c.UseCount == container.UseCount && c.UsageMem < c.UsageMem {
			container = c
			continue
		}

	}

	if container == nil {
		return nil
	}

	//修改container的使用情况
	container.UseCount++
	container.Node.UseCount++

	return container
}

//归还container
func Return(container *Container, usageMem int64, runTime int64) {
	if container == nil {
		return
	}
	container.ConcurrencyCount = 2 * 1024 * 1024 * 1024 / usageMem
	container.UseCount--
	container.Node.UseCount--
}
