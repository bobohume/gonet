package base

type(
	stNode struct{
		Value interface{}
		prev, next *stNode
	}

	Deque struct{
		m_pFrontNode, m_pBackNode *stNode
	}

	iDeque interface {
		PushBack(interface{})
		PushFront(interface{})
		PopBack()
		PopFront()
		Back()interface{}
		Front()interface{}
		Empty() bool
	}
)

func (this *Deque)PushBack(val interface{})  {
	pNode := &stNode{val, nil, nil}
	if this.m_pBackNode == nil{
		this.m_pBackNode = pNode
		this.m_pFrontNode = pNode
	}else{
		pNode.next = this.m_pBackNode
		this.m_pBackNode.prev = pNode
		this.m_pBackNode = pNode
	}
}

func (this *Deque)PushFront(val interface{})  {
	pNode := &stNode{val, nil, nil}
	if this.m_pFrontNode == nil{
		this.m_pBackNode = pNode
		this.m_pFrontNode = pNode
	}else{
		pNode.prev = this.m_pFrontNode
		this.m_pFrontNode.next = pNode
		this.m_pFrontNode = pNode
	}
}

func (this *Deque)PopBack()  {
	if this.m_pBackNode != nil{
		this.m_pBackNode = this.m_pBackNode.next
		if this.m_pBackNode != nil {
			this.m_pBackNode.prev = nil
		}else{
			this.m_pFrontNode =nil
		}
	}
}

func (this *Deque)PopFront()  {
	if this.m_pFrontNode != nil{
		this.m_pFrontNode = this.m_pFrontNode.prev
		if this.m_pFrontNode != nil {
			this.m_pFrontNode.next = nil
		} else{
			this.m_pBackNode =nil
		}
	}
}

func (this *Deque)Back()  interface{}{
	if this.m_pBackNode != nil{
		return  this.m_pBackNode.Value
	}
	return nil
}

func (this *Deque)Front()  interface{}{
	if this.m_pFrontNode != nil{
		return  this.m_pFrontNode.Value
	}
	return nil
}

func (this *Deque)Empty()  bool {
	return this.m_pBackNode == nil
}