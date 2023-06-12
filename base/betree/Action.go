package betree

type (
	Action struct {
		BaseNode
	}

	IAction interface {
		IBaseNode
	}
)

func (a *Action) Init() {
	a.Type = ACTION
}
