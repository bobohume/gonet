package betree

type(
	Composite struct {
		BaseNode
		BehaviorList
	}

	IComposite interface {
		IBaseNode
		IBehaviorList
	}

	Sequence struct {
		Composite
	}

	Selector struct {
		Composite
	}

	//并发
	PSequence struct {
		Composite
	}

	//并发
	PSelector struct {
		Composite
	}
)

func (this *Composite) Init() {
	this.Type = COMPOSITE
}

//	  当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回False，那停止迭代，
//    本Node向自己的Parent Node也返回False；否则所有Child Node都返回True，
//    那本Node向自己的Parent Node返回True。
func (this *Sequence) OnExec(tick int64) bool {
	for _,v := range this.BehaviorList.Values() {
		if !v.(IBaseNode).OnExec(tick){
			return false
		}
	}
	return true
}

//    当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回True，那停止迭代，
//    本Node向自己的Parent Node也返回True；否则所有Child Node都返回False，
//    那本Node向自己的Parent Node返回False。
func (this *Selector) OnExec(tick int64) bool {
	for _,v := range this.BehaviorList.Values() {
		if v.(IBaseNode).OnExec(tick){
			return true
		}
	}
	return false
}

//	  当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回False，那停止迭代，
//    本Node向自己的Parent Node也返回False；否则所有Child Node都返回True，
//    那本Node向自己的Parent Node返回True。
func (this *PSequence) OnExec(tick int64) bool {
	bScuess := false
	for _,v := range this.BehaviorList.Values() {
		if v.(IBaseNode).OnExec(tick){
			bScuess = true
		}
	}
	return bScuess
}

//    当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//    如遇到一个Child Node执行后返回True，那停止迭代，
//    本Node向自己的Parent Node也返回True；否则所有Child Node都返回False，
//    那本Node向自己的Parent Node返回False。
func (this *PSelector) OnExec(tick int64) bool {
	bScuess := true
	for _,v := range this.BehaviorList.Values() {
		if !v.(IBaseNode).OnExec(tick){
			bScuess = false
		}
	}
	return bScuess
}