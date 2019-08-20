package lmath

import (
	"gonet/base"
	"math"
)

type(
	Point2F struct {
		X float32
		Y float32
	}

	IPoint2F interface {
		Set(float32 , float32)
		SetF([] float32)
		SetMin(Point2F)
		SetMax(Point2F)
		Interpolate(Point2F, Point2F, float32)
		Zero()
		IsZero() bool
		Len() float32
		LenSquared() float32
		MagnitudeSafe() float32
		Equal(Point2F) bool
		Neg()
		Normalize()
		NormalizeSafe()
		NormalizeF(float32)
		Convolve(Point2F)
		ConvolveInverse(Point2F)

		Add(Point2F) *Point2F
		Sub(Point2F) *Point2F
		Mul(Point2F) *Point2F
		MulF(float32) *Point2F
		Div(f float32) *Point2F

		Cross(p Point2F) float32
		Dot(p Point2F) float32

		ToF32() []*float32
		ToF() []float32
	}
)

func (this *Point2F) Set(x float32, y float32){
	this.X, this.Y = x, y
}

func (this *Point2F) SetF(f []float32){
	this.X, this.Y = f[0], f[1]
}

func (this *Point2F) SetMin(p Point2F){
	this.X, this.Y = float32(math.Min(float64(this.X), float64(p.X))), float32(math.Min(float64(this.Y), float64(p.Y)))
}

func (this *Point2F) SetMax(p Point2F){
	this.X, this.Y = float32(math.Max(float64(this.X), float64(p.X))), float32(math.Max(float64(this.Y), float64(p.Y)))
}

func (this *Point2F) Interpolate(from Point2F, to Point2F, factor float32){
	base.Assert(factor >= 0.0 && factor <= 1.0, "Out of bound interpolation factor")
	inverse := 1.0 - factor
	this.X = from.X * inverse + to.X * factor
	this.Y = from.Y * inverse + to.Y * factor
}

func (this *Point2F) Zero(){
	this.X, this.Y = 0, 0
}

func (this *Point2F) IsZero() bool{
	return ((this.X * this.X) <= POINT_EPSILON) &&((this.Y * this.Y) <= POINT_EPSILON)
}

func (this *Point2F) Len() float32{
	return float32(math.Sqrt(float64(this.X* this.X + this.Y * this.Y)))
}

func (this *Point2F) LenSquared() float32{
	if this.IsZero(){
		return 0.0
	}else{
		return this.Len()
	}
}

func (this *Point2F) MagnitudeSafe() float32{
	return this.X * this.X + this.Y * this.Y
}

func (this *Point2F) Equal(p Point2F) bool{
	return ((math.Abs(float64(this.X - p.X)) < POINT_EPSILON) &&
		(math.Abs(float64(this.Y - p.Y)) < POINT_EPSILON))
}

func (this *Point2F) Neg(){
	this.X, this.Y = -this.X, -this.Y
}

func (this *Point2F) Normalize(){
	squared := this.X * this.X + this.Y * this.Y
	if squared != 0{
		factor := 1.0 / float32(math.Sqrt(float64(squared)))
		this.X *= factor
		this.Y *= factor
	}else{
		this.X = 0.0
		this.Y = 0.0
	}
}

func (this *Point2F)  NormalizeSafe(){
	vmag := this.MagnitudeSafe()
	if vmag > POINT_EPSILON{
		*this = *this.MulF(1.0 / vmag)
	}
}

func (this *Point2F) NormalizeF(f float32){
	squared := this.X * this.X + this.Y * this.Y
	if squared != 0{
		factor := f / float32(math.Sqrt(float64(squared)))
		this.X *= factor
		this.Y *= factor
	}else{
		this.X = 0.0
		this.Y = 0.0
	}
}

func (this *Point2F) Convolve(p Point2F){
	this.X *= p.X
	this.Y *= p.Y
}

func (this *Point2F) ConvolveInverse(p Point2F){
	this.X /= p.X
	this.Y /= p.Y
}


func (this *Point2F) Add(p Point2F)  *Point2F{
	return &Point2F{this.X + p.X, this.Y + p.Y}
}

func (this *Point2F) Sub(p Point2F)  *Point2F{
	return &Point2F{this.X - p.X, this.Y - p.Y}
}

func (this *Point2F) MulF(f float32) *Point2F{
	return &Point2F{this.X * f, this.Y * f}
}

func (this *Point2F) Mul(p Point2F)  *Point2F{
	return &Point2F{this.X * p.X, this.Y * p.Y}
}

func (this *Point2F) Div(f float32) *Point2F{
	return &Point2F{this.X / f, this.Y / f}
}

func (this *Point2F) ToF32() []*float32{
	return  []*float32{&this.X, &this.Y}
}

func (this *Point2F) ToF() []float32{
	return  []float32{this.X, this.Y}
}

func (this *Point2F) Cross(p Point2F) float32{
	return this.X * p.Y - this.Y * p.X
}

func (this *Point2F) Dot(p Point2F) float32{
	return this.X * p.X + this.Y * p.Y
}
