package core




/*

func NewNode(nodeId string, address string, port int64, maxMem int64) *Node {
	node := &Node{NodeID: nodeId, Address: address, Port: port, MaxMem: maxMem, Containers: NewLM()}
	node.MaxMem -= 512 * 1024 * 1024 //每个节点预留512M的空间，不使用完
	return node
}

//添加Container,
func (node *Node) AddContainer(container *Container) {
	node.Containers.Add(container.FunName, container)
}

//移除Container，会将container销毁
func (node *Node) RemoveContainer(containerId string) {

}

//租用container，会消耗内存
func (node *Node) RentContainer(container *Container) (*Container, error) {
	c := node.QueryContainer(container.FunName, container.UsedMem)
	if c == nil {
		return nil, errors.New("No Containers available or lack of memory")
	}
	node.UsedMem += c.UsedMem
	return c, nil
}

func (node *Node) ReturnContainer(container *Container) {
	node.UsedMem -= container.UsedMem
}

//查询container，不会消耗内存
func (node *Node) QueryContainer(funcName string, reqMem int64) *Container {
	//如果内存不够就直接返回
	if node.UsedMem+reqMem > node.MaxMem {
		return nil
	}
	container := node.Containers.Get(funcName)
	if container == nil {
		return nil
	}
	return container.(*Container)
}
*/
