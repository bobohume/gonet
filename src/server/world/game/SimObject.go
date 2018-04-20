package game

import "sync/atomic"

type (
	SimHandle struct {
		SimId int
		Type int
	}

	ISimHandle interface {
		Empty() bool
		Clear()
	}

	SimObject struct {
		Pos	Position
		SimId	int
		Type int
		pMap *Map
	}

	ISimObject interface {
		Init(int)
		Reset()
		GetHandle() SimHandle

		SetAtTitle(pMap *Map, x int, y int, noOldArea bool)
		GetTitleDistance(pos Position) float32
		GetTitleDistance1(x float32, y float32) float32
		GetMap() *Map
		GetTite() *Tile
		GetArea() *Area
		generateEntiyId()
	}
)

var(
	g_SimId int32
)

func (this *SimHandle) Empty() bool{
	return this.SimId <= 0 || this.Type == 0
}

func (this *SimHandle) Clear(){
	this.SimId, this.Type = 0, 0
}

func (this *SimObject) generateEntiyId() {
	this.SimId = int(atomic.AddInt32(&g_SimId, 1))
}

func (this *SimObject) Init(iType int) {
	this.Type = iType
	this.generateEntiyId()
}

func (this *SimObject) Reset(iType int){
	this.pMap = nil
	this.Pos = Position{0, 0}
}

func (this *SimObject) GetHandle() SimHandle{
	return SimHandle{this.SimId, this.Type}
}

func (this *SimObject) GetTitleDistance(pos Position) float32{
	return this.Pos.TitleDistance(pos)
}

func (this *SimObject) GetTitleDistance1(x float32, y float32) float32{
	return this.GetTitleDistance(Position{x, y})
}

func (this *SimObject) GetMap() *Map{
	return this.pMap
}