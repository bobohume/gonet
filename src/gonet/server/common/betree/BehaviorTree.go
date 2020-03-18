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
	BehaviorList struct {
		base.Vector
	}

	IBehaviorList interface {
		AddChild(string, IBaseNode)
		DelChild(string)
		GetChild(string) IBaseNode
		GetChildCount() int
	}

	BehaviorTree struct {
		BehaviorList
	}

	IBehaviorTree interface {
		IBehaviorList
		Init()
		OnExec(int64)
	}
)

func (this *BehaviorList) Less(i, j int) bool{
	return this.Get(i).(IBaseNode).GetName() < this.Get(j).(IBaseNode).GetName()
}

func (this *BehaviorList) AddChild(name string, pNode IBaseNode){
	if pNode.GetType() != COMPOSITE && pNode.GetType() != DECORATOR && pNode.GetType() != ACTION && pNode.GetType() != CONDITION{
		return
	}

	//pNode.Init()
	pCurNode := this.GetChild(name)
	if pCurNode != nil{
		if pCurNode.GetType() == COMPOSITE && pNode.(IComposite) != nil{
			pCurNode.(IComposite).AddChild(name, pNode)
		}else if pCurNode.GetType() == DECORATOR && pNode.(IDecorator) != nil{
			pCurNode.(IDecorator).SetChild(pNode)
		}
	}else{
		pNode.SetName(name)
		this.Push_front(pNode)
		sort.Sort(this)
	}
}

func (this *BehaviorList) DelChild(name string){
	nIndex := sort.Search(this.Len(), func(i int) bool {
		return this.Get(i).(IBaseNode).GetName() >= name
	})
	if nIndex < this.Len() && this.Get(nIndex).(IBaseNode).GetName() == name{
		this.Erase(nIndex)
	}
}

func (this *BehaviorList) GetChild(name string) IBaseNode{
	nIndex := sort.Search(this.Len(), func(i int) bool {
		return this.Get(i).(IBaseNode).GetName() >= name
	})
	if nIndex < this.Len() && this.Get(nIndex).(IBaseNode).GetName() == name{
		return this.Get(nIndex).(IBaseNode)
	}
	return nil
}

func (this *BehaviorList) GetChildCount() int {
	return this.Len()
}

func (this *BehaviorTree) Init(){
}

func (this *BehaviorTree) OnExec(tick int64){
	for _,v := range this.Array() {
		v.(IBaseNode).OnExec(tick)
	}
}

func NewBeTree() *BehaviorTree {
	tree := &BehaviorTree{}
	tree.Init()
	return tree
}