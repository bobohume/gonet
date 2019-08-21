package lmath

import (
	"math"
)

type (
	Box3F struct {
		Min Point3F///< Minimum extents of box
		Max Point3F///< Maximum extents of box
	}

	IBox3F interface {
		/// Check to see if another box is contained in this box.
		IsContained(Box3F) bool
		/// Check to see if a point is contained in this box.
		IsContainedp(Point3F) bool
		/// Check to see if another box overlaps this box.
		IsOverlapped(Box3F) bool

		Len_x() float32
		Len_y() float32
		Len_z() float32

		/// Perform an intersection operation with another box
		/// and store the results in this box.
		Intersect(Box3F)
		Intersectp(Point3F)

		/// Get the center of this box.
		///
		/// This is the average of min and max.
		GetCenter(Point3F)

		/// Collide a line against the box.
		///
		/// @param   start   Start of line.
		/// @param   end     End of line.
		/// @param   t       Value from 0.0-1.0, indicating position
		///                  along line of collision.
		/// @param   n       Normal of collision.
		CollideLineff(start *Point3F, end *Point3F, t *float32, n *Point3F) bool

		/// Collide a line against the box.
		///
		/// Returns true on collision.
		CollideLine(start *Point3F, end *Point3F) bool

		/// Collide an oriented box against the box.
		///
		/// Returns true if "oriented" box collides with us.
		/// Assumes incoming box is centered at origin of source space.
		///
		/// @param   radii   The dimension of incoming box (half x,y,z length).
		/// @param   toUs    A transform that takes incoming box into our space.
		CollideOrientedBox(radii *Point3F, toUs *MatrixF) bool

		/// Check that the box is valid.
		///
		/// Currently, this just means that min < max.
		IsValidBox() bool

		/// Return the closest point of the box, relative to the passed point.
		GetClosestPoint(refPt Point3F) Point3F

		/// Return distance of closest point on box to refPt.
		GetDistanceFromPoint(refPt Point3F) float32

		/// Extend box to include point
		Extend(p Point3F)

		SetInvalid()
		SetMaxSize()
	}
)

func (this *Box3F) IsContainedp(p Point3F) bool{
	return (p.X >= this.Min.X && p.X < this.Max.X) &&
		(p.Y >= this.Min.Y && p.Y < this.Max.Y) &&
		(p.Z >= this.Min.Z && p.Z < this.Max.Z)
}

func (this *Box3F) IsContained(b Box3F) bool{
	return this.Min.X <= b.Min.X &&
		this.Min.Y <= b.Min.Y &&
		this.Min.Z <= b.Min.Z &&
		this.Max.X >= b.Max.X &&
		this.Max.Y >= b.Max.Y &&
		this.Max.Z >= b.Max.Z
}

func (this *Box3F) IsOverlapped(b Box3F) bool{
	if b.Min.X > this.Max.X || b.Min.Y > this.Max.Y || b.Min.Z > this.Max.Z{
		return false
	}
	if b.Max.X < this.Min.X || b.Max.Y < this.Min.Y || b.Max.Z < this.Min.Z{
		return false
	}

	return true
}

func (this *Box3F) Len_x() float32{
	return  this.Max.X - this.Min.X
}

func (this *Box3F) Len_y() float32{
	return  this.Max.Y - this.Min.Y
}

func (this *Box3F) Len_z() float32{
	return  this.Max.Z - this.Min.Z
}

func (this *Box3F) Intersect(b Box3F) {
	this.Min.SetMin(b.Min)
	this.Max.SetMax(b.Max)
}

func (this *Box3F) Intersectp(b Point3F) {
	this.Min.SetMin(b)
	this.Max.SetMax(b)
}

func (this *Box3F) IsValidBox() bool {
	return this.Min.X <= this.Max.X && this.Min.Y <= this.Max.Y && this.Min.Z <= this.Max.Z
}

func (this *Box3F) SetInvalid()  {
	this.Min.Set(1e9, 1e9, 1e9)
	this.Max.Set(-1e9, -1e9, -1e9)
}

func (this *Box3F) SetMaxSize()  {
	this.Max.Set(1e9, 1e9, 1e9)
	this.Min.Set(-1e9, -1e9, -1e9)
}

func (this *Box3F) GetClosestPoint(refPt Point3F) Point3F{
	var closest Point3F
	if refPt.X <= this.Min.X{
		closest.X = this.Min.X
	}else if(refPt.X > this.Max.X){
		closest.X = this.Max.X
	}else{
		closest.X = refPt.X
	}

	if refPt.Y <= this.Min.Y{
		closest.Y = this.Min.Y
	}else if(refPt.Y > this.Max.Y){
		closest.Y = this.Max.Y
	}else{
		closest.Y = refPt.Y
	}

	if refPt.Z <= this.Min.Z{
		closest.Z = this.Min.Z
	}else if(refPt.Z > this.Max.Z){
		closest.Z = this.Max.Z
	}else{
		closest.Z = refPt.Z
	}
	return closest
}

func (this *Box3F) GetDistanceFromPoint(refPt Point3F) float32{
	var vec Point3F

	if refPt.X < this.Min.X{
		vec.X = this.Min.X - refPt.X
	}else if(refPt.X > this.Max.X){
		vec.X = refPt.X - this.Max.X
	}else{
		vec.X = 0
	}

	if refPt.Y < this.Min.Y{
		vec.Y = this.Min.Y - refPt.Y
	}else if(refPt.Y > this.Max.Y){
		vec.Y = refPt.Y - this.Max.Y
	}else{
		vec.Y = 0
	}

	if refPt.Z < this.Min.Z{
		vec.Z = this.Min.Z - refPt.Z
	}else if(refPt.Z > this.Max.Z){
		vec.Z = refPt.Z - this.Max.Z
	}else{
		vec.Z = 0
	}
	return vec.Len()
}

func (this *Box3F) Extend(p Point3F){
	if p.X < this.Min.X{
		this.Min.X = p.X
	}else if(p.X > this.Max.X){
		this.Max.X = p.X
	}
	if p.Y < this.Min.Y{
		this.Min.Y = p.Y
	}else if(p.Y > this.Max.Y){
		this.Max.Y = p.Y
	}
	if p.Z < this.Min.Z{
		this.Min.Z = p.Z
	}else if(p.Z > this.Max.Z){
		this.Max.Z = p.Z
	}
}

func (this *Box3F) GetCenter(b Point3F) {
	b.X, b.Y, b.Z = (this.Min.X + this.Max.X) * 0.5, (this.Min.Y + this.Max.Y) * 0.5, (this.Min.Z + this.Max.Z) * 0.5
}

func (this *Box3F) CollideLineff(start *Point3F, end *Point3F, t *float32, n *Point3F) bool{
	var st, et, fst, fet float32
	bmin, bmax, si, ei := this.Min.ToF32(), this.Max.ToF32(), start.ToF32(), end.ToF32()

	na := [3]Point3F {Point3F{1.0, 0.0, 0.0}, Point3F{0.0, 1.0, 0.0}, Point3F{0.0, 0.0, 1.0}}
	var finalNormal Point3F

	for i := 0; i < 3; i++{
		n_neg := false
		if *si[i] < *ei[i] {
			if *si[i] > *bmax[i] || *ei[i] < *bmin[i]{
				return false
			}

			var di float32
			di = *ei[i] - *si[i]
			if *si[i] < *bmin[i]{
				st = (*bmin[i] - *si[i]) / di
			}else{
				st = 0.0
			}

			if *ei[i] > *bmax[i]{
				et = (*bmax[i] - *si[i]) / di
			}else{
				et = 1.0
			}
			n_neg = true
		}else{
			if *ei[i] > *bmax[i] || *si[i] < *bmin[i]{
				return false
			}

			var di float32
			di = *ei[i] - *si[i]
			if *si[i] > *bmax[i]{
				st = (*bmax[i] - *si[i]) / di
			}else{
				st = 0.0
			}

			if *ei[i] < *bmin[i]{
				et = (*bmin[i] - *ei[i]) / di
			}else{
				et = 1.0
			}
		}

		if (st > fst){
			fst = st
			finalNormal = na[i]
			if n_neg{
				finalNormal.Neg()
			}
		}

		if et < fet{
			fet = et
		}

		if fet < fst{
			return false
		}
	}

	*t = fst
	*n = finalNormal
	return true
}

func (this *Box3F) CollideLine(start *Point3F, end *Point3F) bool{
	var t float32
	var normal Point3F
	return this.CollideLineff(start, end, &t, &normal)
}

// returns true if "oriented" box collides with us
// radiiB is dimension of incoming box (half x,y,z length
// toA is transform that takes incoming box into our space
// assumes incoming box is centered at origin of source space
func (this *Box3F) CollideOrientedBox(bRadii *Point3F, toA *MatrixF) bool{
	var p Point3F;
	toA.GetColumn(3,&p)
	aCenter := this.Min.Add(this.Max)
	aCenter.MulF(0.5)
	p = *p.Sub(*aCenter)
	aRadii := this.Max.Sub(this.Min)
	aRadii.MulF(0.5)

	var absXX,absXY,absXZ float32
	var absYX,absYY,absYZ float32
	var absZX,absZY,absZZ float32

    f := toA.ToF32();

	absXX = float32(math.Abs(float64(*f[0])))
	absYX = float32(math.Abs(float64(*f[1])))
	absZX = float32(math.Abs(float64(*f[2])))

	if (aRadii.X + bRadii.X * absXX + bRadii.Y * absYX + bRadii.Z * absZX - float32(math.Abs(float64(p.X)))<0.0){
		return false
	}

	absXY = float32(math.Abs(float64(*f[4])))
	absYY = float32(math.Abs(float64(*f[5])))
	absZY = float32(math.Abs(float64(*f[6])))
	if (aRadii.Y + bRadii.X * absXY +	bRadii.Y * absYY +	bRadii.Z * absZY - float32(math.Abs(float64(p.Y)))<0.0){
		return false
	}

	absXZ = float32(math.Abs(float64(*f[8])))
	absYZ = float32(math.Abs(float64(*f[9])))
	absZZ = float32(math.Abs(float64(*f[10])))
	if (aRadii.Z + bRadii.X * absXZ + bRadii.Y * absYZ +	bRadii.Z * absZZ - float32(math.Abs(float64(p.Z)))<0.0){
		return false
	}

	if (aRadii.X*absXX + aRadii.Y*absXY + aRadii.Z*absXZ + bRadii.X - float32(math.Abs(float64(p.X**f[0] + p.Y**f[4] + p.Z**f[8])))<0.0){
		return false
	}

	if (aRadii.X*absYX + aRadii.Y*absYY + aRadii.Z*absYZ + bRadii.Y - float32(math.Abs(float64(p.X**f[1] + p.Y**f[5] + p.Z**f[9])))<0.0) {
		return false
	}

	if (aRadii.X*absZX + aRadii.Y*absZY + aRadii.Z*absZZ + bRadii.Z - float32(math.Abs(float64(p.X**f[2] + p.Y**f[6] + p.Z**f[10])))<0.0) {
		return false
	}

	if (float32(math.Abs(float64(p.Z*(*f[4]) - p.Y*(*f[8])))) >
		aRadii.Y * absXZ + aRadii.Z * absXY +
			bRadii.Y * absZX + bRadii.Z * absYX){
				return false
	}

	if (float32(math.Abs(float64(p.Z**f[5] - p.Y**f[9]))) >
		aRadii.Y * absYZ + aRadii.Z * absYY +
			bRadii.X * absZX + bRadii.Z * absXX){
				return false
	}

	if (float32(math.Abs(float64(p.Z**f[6] - p.Y**f[10]))) >
		aRadii.Y * absZZ + aRadii.Z * absZY +
			bRadii.X * absYX + bRadii.Y * absXX){
				return false
	}

	if (float32(math.Abs(float64(p.X**f[8] - p.Z**f[0]))) >
		aRadii.X * absXZ + aRadii.Z * absXX +
			bRadii.Y * absZY + bRadii.Z * absYY){
				return false
	}

	if (float32(math.Abs(float64(p.X**f[9] - p.Z**f[1]))) >
		aRadii.X * absYZ + aRadii.Z * absYX +
			bRadii.X * absZY + bRadii.Z * absXY){
				return false
	}

	if (float32(math.Abs(float64(p.X**f[10] - p.Z**f[2]))) >
		aRadii.X * absZZ + aRadii.Z * absZX +
			bRadii.X * absYY + bRadii.Y * absXY){
				return false
	}

	if (float32(math.Abs(float64(p.Y**f[0] - p.X**f[4]))) >
		aRadii.X * absXY + aRadii.Y * absXX +
			bRadii.Y * absZZ + bRadii.Z * absYZ){
				return false
	}

	if (float32(math.Abs(float64(p.Y**f[1] - p.X**f[5]))) >
		aRadii.X * absYY + aRadii.Y * absYX +
			bRadii.X * absZZ + bRadii.Z * absXZ){
				return false
	}

	if (float32(math.Abs(float64(p.Y**f[2] - p.X**f[6]))) >
		aRadii.X * absZY + aRadii.Y * absZX +
			bRadii.X * absYZ + bRadii.Y * absXZ){
				return false
	}

	return true
}