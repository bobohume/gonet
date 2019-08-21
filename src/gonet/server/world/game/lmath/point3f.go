package lmath

import (
	"math"
	"gonet/base"
)

const(
	POINT_EPSILON = (1e-4)
)

type(
	Point3F struct {
		X float32
		Y float32
		Z float32
	}

	IPoint3F interface {
		Set(float32 , float32, float32)
		SetF([] float32)
		SetMin(Point3F)
		SetMax(Point3F)
		Interpolate(Point3F, Point3F, float32)
		Zero()
		IsZero() bool
		Len() float32
		LenSquared() float32
		MagnitudeSafe() float32
		Equal(Point3F) bool
		Neg()
		Normalize()
		NormalizeSafe()
		NormalizeF(float32)
		Convolve(Point3F)
		ConvolveInverse(Point3F)

		Add(Point3F) *Point3F
		Sub(Point3F) *Point3F
		Mul(Point3F) *Point3F
		MulF(float32) *Point3F
		Div(f float32) *Point3F

		Cross(b Point3F) Point3F
		Dot(p Point3F) float32

		Perp2D(u  Point3F) float32
		Dot2D(u  Point3F) float32

		ToF32() []*float32
		ToF() []float32
	}
)

func (this *Point3F) Set(x float32, y float32, z float32){
	this.X, this.Y, this.Z = x, y, z
}

func (this *Point3F) SetF(f []float32){
	this.X, this.Y, this.Z = f[0], f[1], f[2]
}

func (this *Point3F) SetMin(p Point3F){
	this.X, this.Y, this.Z = float32(math.Min(float64(this.X), float64(p.X))), float32(math.Min(float64(this.Y), float64(p.Y))), float32(math.Min(float64(this.Z), float64(p.Z)))
}

func (this *Point3F) SetMax(p Point3F){
	this.X, this.Y, this.Z = float32(math.Max(float64(this.X), float64(p.X))), float32(math.Max(float64(this.Y), float64(p.Y))), float32(math.Max(float64(this.Z), float64(p.Z)))
}

func (this *Point3F) Interpolate(from Point3F, to Point3F, factor float32){
	base.Assert(factor >= 0.0 && factor <= 1.0, "Out of bound interpolation factor")
	inverse := 1.0 - factor
	this.X = from.X * inverse + to.X * factor
	this.Y = from.Y * inverse + to.Y * factor
	this.Z = from.Z * inverse + to.Z * factor
}

func (this *Point3F) Zero(){
	this.X, this.Y, this.Z = 0, 0, 0
}

func (this *Point3F) IsZero() bool{
	return ((this.X * this.X) <= POINT_EPSILON) &&((this.Y * this.Y) <= POINT_EPSILON) && ((this.Z * this.Z) <= POINT_EPSILON)
}

func (this *Point3F) Len() float32{
	return float32(math.Sqrt(float64(this.X* this.X + this.Y * this.Y + this.Z * this.Z)))
}

func (this *Point3F) LenSquared() float32{
	if this.IsZero(){
		return 0.0
	}else{
		return this.Len()
	}
}

func (this *Point3F) MagnitudeSafe() float32{
	return this.X * this.X + this.Y * this.Y + this.Z * this.Z
}

func (this *Point3F) Equal(p Point3F) bool{
	return ((math.Abs(float64(this.X - p.X)) < POINT_EPSILON) &&
		(math.Abs(float64(this.Y - p.Y)) < POINT_EPSILON) &&
		(math.Abs(float64(this.Z - p.Z)) < POINT_EPSILON))
}

func (this *Point3F) Neg(){
	this.X, this.Y, this.Z = -this.X, -this.Y, -this.Z
}

func (this *Point3F) Normalize(){
	squared := this.X * this.X + this.Y * this.Y + this.Z * this.Z
	if squared != 0{
		factor := 1.0 / float32(math.Sqrt(float64(squared)))
		this.X *= factor
		this.Y *= factor
		this.Z *= factor
	}else{
		this.X = 0.0
		this.Y = 0.0
		this.Z = 1.0
	}
}

func (this *Point3F)  NormalizeSafe(){
	vmag := this.MagnitudeSafe()
	if vmag > POINT_EPSILON{
		*this = *this.MulF(1.0 / vmag)
	}
}

func (this *Point3F) NormalizeF(f float32){
	squared := this.X * this.X + this.Y * this.Y + this.Z * this.Z
	if squared != 0{
		factor := f / float32(math.Sqrt(float64(squared)))
		this.X *= factor
		this.Y *= factor
		this.Z *= factor
	}else{
		this.X = 0.0
		this.Y = 0.0
		this.Z = 1.0
	}
}

func (this *Point3F) Convolve(p Point3F){
	this.X *= p.X
	this.Y *= p.Y
	this.Z *= p.Z
}

func (this *Point3F) ConvolveInverse(p Point3F){
	this.X /= p.X
	this.Y /= p.Y
	this.Z /= p.Z
}


func (this *Point3F) Add(p Point3F)  *Point3F{
	return &Point3F{this.X + p.X, this.Y + p.Y, this.Z + p.Z}
}

func (this *Point3F) Sub(p Point3F)  *Point3F{
	return &Point3F{this.X - p.X, this.Y - p.Y, this.Z - p.Z}
}

func (this *Point3F) MulF(f float32) *Point3F{
	return &Point3F{this.X * f, this.Y * f, this.Z * f}
}

func (this *Point3F) Mul(p Point3F)  *Point3F{
	return &Point3F{this.X * p.X, this.Y * p.Y, this.Z * p.Z}
}

func (this *Point3F) Div(f float32) *Point3F{
	return &Point3F{this.X / f, this.Y / f, this.Z / f}
}

func (this *Point3F) ToF32() []*float32{
	return  []*float32{&this.X, &this.Y, &this.Z}
}

func (this *Point3F) ToF() []float32{
	return  []float32{this.X, this.Z, this.Y}
}

func (this *Point3F) Cross(p Point3F) Point3F{
	return Point3F{(this.Y * p.Z) - (this.Z * p.Y), (this.Z * p.X) - (this.X * p.Z), (this.X * p.Y) - (this.Y * p.X)}
}

func (this *Point3F) Dot(p Point3F) float32 {
	return this.X * p.X + this.Y * p.Y +this.Z * p.Z
}

func CrossFFF(a []*float32, b []*float32, c []*float32){
	p1 := Point3F{*a[0], *a[1], *a[2]}
	p2 := Point3F{*b[0], *b[1], *b[2]}
	pos := p1.Cross(p2)
	p3 := pos.ToF32()
	for i,v := range p3{
		*c[i] = *v
	}
}

func DotPP(p1 Point3F, p2 Point3F) float32{
	return (p1.X*p2.X + p1.Y*p2.Y + p1.Z*p2.Z)
}

func (this *Point3F) Perp2D(u  Point3F) float32 {
	return this.Z*u.X - this.X*u.Z
}

func (this *Point3F) Dot2D(u Point3F) float32 {
	return this.X*u.X + this.Z*u.Z
}
