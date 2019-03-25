package betree

import (
	"gonet/base"
	"sort"
)

//  它只有4大类型的Node：
//  * Composite Node
//  * Decorator Node
//  * Condition Node
//  * Action Node
type(
	NameList struct {
		List base.Vector
	}

	BehaviorTree struct {
		BehaviorMap map[string] IBaseNode
		BehaviorNameList NameList//just for sort,because range  BehaviorMap is not inorder
	}

	IBehaviorTree interface {
		Init()
		AddNode(string, IBaseNode)
		DelNode(string)
		GetNode(string) IBaseNode
		OnExec(int64)
	}
)

func (this *NameList) Len() int{
	return this.List.Len()
}

func (this *NameList) Swap(i, j int) {
	this.List.Swap(i,j)
}

func (this *NameList) Less(i, j int) bool{
	return this.List.Get(i).(string) < this.List.Get(j).(string)
}

func (this *BehaviorTree) Init(){
	this.BehaviorMap = make(map[string] IBaseNode)
}

func (this *BehaviorTree) AddNode(name string, pNode IBaseNode){
	if pNode.GetType() != COMPOSITE && pNode.GetType() != DECORATOR && pNode.GetType() != ACTION && pNode.GetType() != CONDITION{
		return
	}

	//pNode.Init()
	pCurNode, exist := this.BehaviorMap[name]
	if exist{
		if pCurNode.GetType() == COMPOSITE && pNode.(IComposite) != nil{
			pCurNode.(IComposite).AddChild(name, pNode)
		}else if pCurNode.GetType() == DECORATOR && pNode.(IDecorator) != nil{
			pCurNode.(IDecorator).SetChild(pNode)
		}
	}else{
		this.BehaviorMap[name] = pNode
		this.BehaviorNameList.List.Push_front(name)
		sort.Sort(&this.BehaviorNameList)
	}
}

func (this *BehaviorTree) DelNode(name string){
	delete(this.BehaviorMap, name)
	for i,v := range this.BehaviorNameList.List.Array(){
		if v.(string) == name{
			this.BehaviorNameList.List.Erase(i)
			break
		}
	}
}

func (this *BehaviorTree) GetNode(name string) IBaseNode{
	return this.BehaviorMap[name]
}

func (this *BehaviorTree) OnExec(tick int64){
	for i := 0; i < this.BehaviorNameList.List.Len(); i++{
		this.BehaviorMap[this.BehaviorNameList.List.Get(i).(string)].OnExec(tick)
	}
	/*for _, v := range this.BehaviorMap{
		v.OnExec(tick)
	}*/
}

func NewBeTree() *BehaviorTree {
	tree := &BehaviorTree{}
	tree.Init()
	return tree
}