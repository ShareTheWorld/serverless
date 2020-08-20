package core

import (
	"sync"
)

//用于存放所有node
var nodes = make([]*Node, 0, 100)
var NodesLock sync.RWMutex

var FunMap map[string]map[string]*Container
var FuncMapLock sync.RWMutex

var RequestMap map[string]*Container
var RequestMapLock sync.Mutex
