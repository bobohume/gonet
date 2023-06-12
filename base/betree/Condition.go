package betree

type (
	Condition struct {
		BaseNode
	}

	ICondition interface {
		IBaseNode
	}
)

func (c *Condition) Init() {
	c.Type = CONDITION
}
