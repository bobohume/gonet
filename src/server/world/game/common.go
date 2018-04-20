package game

import "math"

const(
	INVALID_DIRECTION = -1
	DOWN = 0
	DOWN_RIGHT = 1
	RIGHT = 2
	UP_RIGHT = 3
	UP = 4
	UP_LEFT = 5
	LEFT = 6
	DOWN_LEFT = 7
	DIRECTION_NUM = 8
)//direciton

const(
	SIM_ALL = iota
	SIM_PLAYER = iota
	SIM_HERO = iota
	SIM_MONSTER = iota
	SIM_PET = iota
)//simobject type

type(
	Position struct {
		X float32
		Y float32
	}

	IPosition interface {
		Equal(Position) bool
		Distance(Position) float64
		TitleDistance(Position) float32
		TitleDistanceSquare(Position) float32
	}

	AttrAddon struct {
		Index int
		Addon int
		Ratio int
	}

	UnitAttr struct {
		hp int
		mp int
		angry int
		shield int
		extra_hp int
		soil_shield int
		fire_shield int
		elec_shield int
		ice_shield int
	}
)

func (this *Position) Equal(pos Position) bool{
	return this.X == pos.X && this.Y == pos.Y
}

func (this *Position) Distance(pos Position) float64{
	return math.Sqrt((float64(pos.X - this.X) * float64(pos.X - this.X) + float64(pos.Y - this.Y) *float64(pos.Y - this.Y)))
}

func (this *Position) TitleDistance(pos Position) float32{
	return float32(math.Max(math.Abs(float64(pos.X - this.X)), math.Abs(float64(pos.Y - this.Y))))
}

func (this *Position) TitleDistanceSquare(pos Position) float32{
	return (pos.X - this.X) * (pos.X - this.X) + (pos.Y - this.Y) *(pos.Y - this.Y)
}