package core

//ResourceManager，主要负责管理各个节点
//所有的节点会放在LM的数据结构中，
//功能点：申请node，释放node，管理node

type RM struct {
	lm LM
}

//返回一个节点对象
//
func (rm RM) getNode(accountId string, reqMem int64) *Node {

}
