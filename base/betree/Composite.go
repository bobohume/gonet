package betree

type (
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

func (c *Composite) Init() {
	c.Type = COMPOSITE
}

//		  当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//	   如遇到一个Child Node执行后返回False，那停止迭代，
//	   本Node向自己的Parent Node也返回False；否则所有Child Node都返回True，
//	   那本Node向自己的Parent Node返回True。
func (s *Sequence) OnExec(tick int64) bool {
	for i := 0; i < s.BehaviorList.Len(); i++ {
		if !s.BehaviorList.Get(i).OnExec(tick) {
			return false
		}
	}
	return true
}

// 当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
// 如遇到一个Child Node执行后返回True，那停止迭代，
// 本Node向自己的Parent Node也返回True；否则所有Child Node都返回False，
// 那本Node向自己的Parent Node返回False。
func (s *Selector) OnExec(tick int64) bool {
	for i := 0; i < s.BehaviorList.Len(); i++ {
		if s.BehaviorList.Get(i).OnExec(tick) {
			return true
		}
	}
	return false
}

//		  当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
//	   如遇到一个Child Node执行后返回False，那停止迭代，
//	   本Node向自己的Parent Node也返回False；否则所有Child Node都返回True，
//	   那本Node向自己的Parent Node返回True。
func (p *PSequence) OnExec(tick int64) bool {
	bScuess := false
	for i := 0; i < p.BehaviorList.Len(); i++ {
		if p.BehaviorList.Get(i).OnExec(tick) {
			bScuess = true
		}
	}
	return bScuess
}

// 当执行本类型Node时，它将从begin到end迭代执行自己的Child Node：
// 如遇到一个Child Node执行后返回True，那停止迭代，
// 本Node向自己的Parent Node也返回True；否则所有Child Node都返回False，
// 那本Node向自己的Parent Node返回False。
func (p *PSelector) OnExec(tick int64) bool {
	bScuess := true
	for i := 0; i < p.BehaviorList.Len(); i++ {
		if !p.BehaviorList.Get(i).OnExec(tick) {
			bScuess = false
		}
	}
	return bScuess
}
