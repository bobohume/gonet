package betree

type(
	Decorator struct {
		BaseNode
		child IBaseNode
	}

	IDecorator interface {
		IBaseNode
		SetChild(child IBaseNode)
		GetChild() IBaseNode
	}

	DecoratorN struct {
		Decorator
	}
)

func (this *Decorator) Init() {
	this.Type = DECORATOR
}

//GetChild
func (this *Decorator) GetChild() IBaseNode {
	return this.child
}

func (this *Decorator) SetChild(child IBaseNode) {
	this.child = child
}

//  它将它的Child Node执行
//  后返回的结果值做额外处理后，再返回给它的Parent Node
func (this *DecoratorN) OnExec(tick int64) bool{
	if this.GetChild() == nil {
		return false
	}

	return this.GetChild().OnExec(tick)
}


