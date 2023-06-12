package betree

type (
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

func (d *Decorator) Init() {
	d.Type = DECORATOR
}

//GetChild
func (d *Decorator) GetChild() IBaseNode {
	return d.child
}

func (d *Decorator) SetChild(child IBaseNode) {
	d.child = child
}

//  它将它的Child Node执行
//  后返回的结果值做额外处理后，再返回给它的Parent Node
func (d *DecoratorN) OnExec(tick int64) bool {
	if d.GetChild() == nil {
		return false
	}

	return d.GetChild().OnExec(tick)
}
