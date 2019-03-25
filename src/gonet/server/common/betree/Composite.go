package betree

import "sort"

type(
	Composite struct {
		BaseNode
		children map[string] IBaseNode
		childrenNameList NameList
	}

	IComposite interface {
		IBaseNode
		GetChildCount() int
		GetChild(name string) IBaseNode
		AddChild(name string, child IBaseNode)
		DelChild(name string)
	}

	Sequence struct {
		Composite
	}

	Selector struct {
		Composite
	}

	//并发
	/*PSequence struct {
		Composite
	}

	//并发
	PSelector struct {
		Composite
	}*/
)

func (this *Composite) Init() {
	this.children = make(map[string] IBaseNode)
	this.Type = COMPOSITE
}

func (this *Composite) GetChildCount() int {
	return len(this.children)
}

func (this *Composite) GetChild(name string) IBaseNode {
	return this.children[name]
}

func (this *Composite) AddChild(name string, child IBaseNode) {
	this.children[name] = child
	this.childrenNameList.List.Push_front(name)
	sort.Sort(&this.childrenNameList)
}

func (this *Composite) DelChild(name string) {
	delete(this.children, name)
	for i,v := range this.childrenNameList.List.Array(){
		if v.(string) == name{
			this.childrenNameList.List.Erase(i)
			break
		}
	}
}

//    当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回True，那停止迭代，
//    本Node向自己的Parent Node也返回True；否则所有Child Node都返回False，
//    那本Node向自己的Parent Node返回False。
func (this *Sequence) OnExec(tick int64) bool {
	for i := 0; i < this.childrenNameList.List.Len(); i++{
		if this.children[this.childrenNameList.List.Get(i).(string)].OnExec(tick){
			return true
		}
	}
	return false
	/*for _,v := range this.children{
		if v.OnExec(tick){
			return true
		}
	}
	return false*/
}

//	  当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回False，那停止迭代，
//    本Node向自己的Parent Node也返回False；否则所有Child Node都返回True，
//    那本Node向自己的Parent Node返回True。
func (this *Selector) OnExec(tick int64) bool {
	for i := 0; i < this.childrenNameList.List.Len(); i++{
		if !this.children[this.childrenNameList.List.Get(i).(string)].OnExec(tick){
			return false
		}
	}
	return true
	/*for _,v := range this.children{
		if !v.OnExec(tick){
			return false
		}
	}
	return true*/
}

/*func (this *PSequence) OnExec(tick int64) bool {
	bScuess := true
	for _,v := range this.children{
		if !v.OnExec(tick){
			bScuess = false
		}
	}
	return bScuess
}

func (this *PSelector) OnExec(tick int64) bool {
	bScuess := false
	for _,v := range this.children{
		if v.OnExec(tick){
			bScuess = true
		}
	}
	return bScuess
}*/