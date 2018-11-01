package betree

type(
	Condition struct {
		BaseNode
	}

	ICondition interface {
		IBaseNode
	}
)

func (this *Condition) Init(){
	this.Type = CONDITION
}