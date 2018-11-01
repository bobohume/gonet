package betree

type(
	Action struct {
		BaseNode
	}

	IAction interface {
		IBaseNode
	}
)

func (this *Action) Init(){
	this.Type = ACTION
}