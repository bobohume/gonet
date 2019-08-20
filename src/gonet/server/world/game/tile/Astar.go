package tile

import (
	"gonet/base"
	"container/heap"
	"math"
)

type(
	Tile struct {
		x int
		y int
	}

	ATile struct {
		Tile
		father *ATile

		gVal float32//权值
		hVal float32//距离
		fVal float32//期望
	}

	IATile interface {
		calcGVal() float32
		calcHVal() float32
		calcFVal() float32
	}

	OpenHeap struct {
		m_Nodes base.Vector//权值队列
	}

	IOpenHeap interface {
		heap.Interface
		Init()
	}
)

//--------------astar-----------------//
func (this *ATile) calcGVal() float32 {
	if this.father != nil {
		deltaX := math.Abs(float64(this.father.x - this.x))
		deltaY := math.Abs(float64(this.father.y - this.y))
		if deltaX == 1 && deltaY == 0 {
			this.gVal = this.father.gVal + 10
		} else if deltaX == 0 && deltaY == 1 {
			this.gVal = this.father.gVal + 10
		} else if deltaX == 1 && deltaY == 1 {
			this.gVal = this.father.gVal + 14
		} else {
			//fmt.Printf("father point is invalid!")
			//this.gVal = this.father.gVal + 24
		}
	}
	return this.gVal
}

func (this *ATile) calcHVal(end *ATile) float32 {
	this.hVal = float32(math.Abs(float64(end.x-this.x)) + math.Abs(float64(end.y-this.y))) * 10
	return this.hVal
}

func (this *ATile) calcFVal(end *ATile) float32 {
	this.fVal = this.calcGVal() + this.calcHVal(end)
	return this.fVal
}

func NewATile(tile Tile, father *ATile, end *ATile) (ap *ATile) {
	ap = &ATile{tile,  father, 0, 0, 0}
	if end != nil {
		ap.calcFVal(end)
	}
	return ap
}

func (this *ATile) IsEqual(p *ATile) bool {
	if this.x == p.x && this.y == p.y{
		return true
	}
	return false
}
//--------------大小堆-----------------//
func (t *OpenHeap) Len() int{
	return t.m_Nodes.Len()
}

func (t *OpenHeap) Swap(i, j int){
	t.m_Nodes.Swap(i, j)
}

//权值比较
func (t *OpenHeap) Less(i, j int) bool{
	if t.m_Nodes.Get(i).(*ATile).fVal == t.m_Nodes.Get(j).(*ATile).fVal{
		return t.m_Nodes.Get(i).(*ATile).hVal < t.m_Nodes.Get(j).(*ATile).hVal
	}
	return  t.m_Nodes.Get(i).(*ATile).fVal < t.m_Nodes.Get(j).(*ATile).fVal
}

func (t *OpenHeap) Push(x interface{}){
	t.m_Nodes.Push_back(x)
}

func (t *OpenHeap) Pop()  interface{} {
	top := t.m_Nodes.Back()
	t.m_Nodes.Pop_back()
	return top
}