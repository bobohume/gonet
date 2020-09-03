package lmath

import (
	"math"
)

type(
	MatrixF [16]float32

	IMatrixF interface {
		Set(Point3F) *MatrixF/// Initialize matrix to rotate about origin by e.
		SetP(Point3F, Point3F) *MatrixF/// Initialize matrix to rotate about p by e.
		SetCrossProduct(Point3F) *MatrixF/// Initialize matrix with a cross product of p.
		SetTensorProduct(Point3F, Point3F) *MatrixF/// Initialize matrix with a tensor product of p.
		IsAffine() bool///< Check to see if this is an affine matrix.
		IsIdentity() bool///< Checks for identity matrix.
		Identity() *MatrixF/// Make this an identity matrix.
		Inverse() *MatrixF/// Invert m.
		AffineInverse() *MatrixF ///< Take inverse of matrix assuming it is affine (rotation, scale, sheer, translation only).
		Transpose() *MatrixF///< Swap rows and columns.
		Scale(Point3F) *MatrixF///< M * Matrix(p) -> M
		GetScale() *Point3F///< Return scale assuming scale was applied via mat.scale(s).
		ToPoint() *Point3F
		/// Compute the inverse of the matrix.
		///
		/// Computes inverse of full 4x4 matrix. Returns false and performs no inverse if
		/// the determinant is 0.
		///
		/// Note: In most cases you want to use the normal inverse function.  This method should
		///       be used if the matrix has something other than (0,0,0,1) in the bottom row.
		//FullInverse() bool
		TransposeTo([]*float32)   /// Swaps rows and columns into matrix.
		Normalize()  /// Normalize the matrix.

		/// Copy the requested column into a Point4F.
		GetColumn(int, *Point3F)

		/// Set the specified column from a Point3F.
		SetColumn(int, *Point3F)

		/// Copy the specified row into a Point4F.
		GetRow(int, *Point3F)

		/// Set the specified row from a Point3F.
		SetRow(int, *Point3F)

		/// Get the position of the matrix.
		GetPosition() Point3F

		/// Set the position of the matrix.
		SetPosition(*Point3F)

		Mulm(*MatrixF) *MatrixF///< M * a -> M
		Mulmm(*MatrixF, *MatrixF) *MatrixF///< a * b -> M

		Mulf(float32)  *MatrixF///< M * a -> M
		Mulmf(*MatrixF, float32) *MatrixF///< a * b -> M

		Mulp(*Point3F) ///< M * p -> p (assume w = 1.0f)
		Mulpp(*Point3F, *Point3F)///< M * p -> d (assume w = 1.0f)
		Mulb(*Box3F)
		ToF32() []*float32
	}
)

func Idx(i int, j int) int{
	return i + j * 4
}

func (this *MatrixF) Set(p Point3F) *MatrixF{
	M_matF_set_euler_C(p.ToF32(), this.ToF32())
	return this
}

func (this *MatrixF) SetP(e Point3F, p Point3F) *MatrixF{
	M_matF_set_euler_point_C(e.ToF32(), p.ToF32(), this.ToF32())
	return this
}

func (this *MatrixF) SetCrossProduct(p Point3F) *MatrixF{
	this[4], this[2], this[9] = p.Z, p.Y, p.Z
	this[1] = -(this[4])
	this[8] = -(this[2])
	this[6] = -(this[9])
	this[0], this[3], this[5], this[7], this[10], this[11], this[12], this[13], this[14] = 0.0,  0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	this[15] = 1
	return this
}

func (this *MatrixF) SetTensorProduct(p Point3F, q Point3F) *MatrixF{
	this[0] = p.X * q.X
	this[1] = p.X * q.Y
	this[2] = p.X * q.Z
	this[4] = p.Y * q.X
	this[5] = p.Y * q.Y
	this[6] = p.Y * q.Z
	this[8] = p.Z * q.X
	this[9] = p.Z * q.Y
	this[10] = p.Z * q.Z
	this[3], this[7], this[11], this[12], this[13], this[14] = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	this[15] = 1.0
	return this
}

func (this *MatrixF) IsIdentity() bool{
	return  this[0]  == 1.0 &&
		this[1]  == 0.0 &&
		this[2]  == 0.0 &&
		this[3]  == 0.0 &&
		this[4]  == 0.0 &&
		this[5]  == 1.0 &&
		this[6]  == 0.0 &&
		this[7]  == 0.0 &&
		this[8]  == 0.0 &&
		this[9]  == 0.0 &&
		this[10] == 1.0 &&
		this[11] == 0.0 &&
		this[12] == 0.0 &&
		this[13] == 0.0 &&
		this[14] == 0.0 &&
		this[15] == 1.0
}

func (this *MatrixF) Identity() *MatrixF{
	this[0]  = 1.0
	this[1]  = 0.0
	this[2]  = 0.0
	this[3]  = 0.0
	this[4]  = 0.0
	this[5]  = 1.0
	this[6]  = 0.0
	this[7]  = 0.0
	this[8]  = 0.0
	this[9]  = 0.0
	this[10] = 1.0
	this[11] = 0.0
	this[12] = 0.0
	this[13] = 0.0
	this[14] = 0.0
	this[15] = 1.0
	return this
}

func (this *MatrixF) Inverse() *MatrixF{
	M_matF_identity_C(this.ToF32())
	return this
}

func (this *MatrixF) AffineInverse() *MatrixF{
	M_matF_affineInverse_C(this.ToF32())
	return this
}

func (this *MatrixF) Transpose() *MatrixF{
	M_matF_transpose_C(this.ToF32())
	return this
}

func (this *MatrixF) Scale(p Point3F) *MatrixF{
	M_matF_scale_C(this.ToF32(), p.ToF32())
	return this
}

func (this *MatrixF) GetScale() *Point3F{
	var scale Point3F
	scale.X = float32(math.Sqrt(float64(this[0]*this[0] + this[4] * this[4] + this[8] * this[8])))
	scale.Y = float32(math.Sqrt(float64(this[1]*this[1] + this[5] * this[5] + this[9] * this[9])))
	scale.Z = float32(math.Sqrt(float64(this[2]*this[2] + this[6] * this[6] + this[10] * this[10])))
	return &scale
}

func (this *MatrixF) Normalize(){
	M_matF_normalize_C(this.ToF32())
}

func (this *MatrixF) Mulm(a *MatrixF) *MatrixF{
	tempThis := *this
	Default_matF_x_matF_C(tempThis.ToF32(), a.ToF32(), this.ToF32())
	return this
}

func (this *MatrixF) Mulmm(a *MatrixF, b *MatrixF) *MatrixF{
	Default_matF_x_matF_C(a.ToF32(), b.ToF32(), this.ToF32())
	return this
}

func (this *MatrixF) Mulf(a float32) *MatrixF{
	for i, _ := range this{
		this[i] *= a
	}
	return this
}

func (this *MatrixF) Mulmf(a *MatrixF, b float32) *MatrixF{
	*this = *a
	this.Mulf(b)
	return this
}

func (this *MatrixF) Mulp(p *Point3F){
	var d Point3F
	M_matF_x_point3F_C(this.ToF32(), p.ToF32(), d.ToF32())
	*p = d
}

func (this *MatrixF) Mulpp(p *Point3F, d *Point3F){
	M_matF_x_point3F_C(this.ToF32(), p.ToF32(), d.ToF32())
}

func (this *MatrixF) Mulb(b *Box3F){
	M_matF_x_box3F_C(this.ToF32(), b.Min.ToF32(), b.Max.ToF32())
}

func (this *MatrixF) GetRow(row int, cptr *Point3F){
	row *= 4
	cptr.X = this[row]
	row++
	cptr.Y = this[row]
	row++
	cptr.Z = this[row]
}

func (this *MatrixF) SetRow(row int, cptr *Point3F){
	row *= 4
	this[row] = cptr.X
	row++
	this[row] = cptr.Y
	row++
	this[row] = cptr.Z
}

func (this *MatrixF) GetPosition() Point3F{
	var pos Point3F
	this.GetColumn(3, &pos)
	return pos
}

func (this *MatrixF) SetPosition(pos *Point3F){
	this.SetColumn(3, pos)
}

func (this *MatrixF) SetColumn(col int, cptr *Point3F){
	this[col] = cptr.X
	this[col+4] = cptr.Y
	this[col+8] = cptr.Z
}

func (this *MatrixF) GetColumn(col int, cptr *Point3F){
	cptr.X = this[col]
	cptr.Y = this[col+4]
	cptr.Z = this[col+8]
}

func (this *MatrixF) ToPoint() *Point3F{
	mat := this.ToF32()
	var r Point3F
	r.X = float32(math.Sin(float64(*mat[Idx(2, 1)])))

	if math.Cos(float64(r.X)) != 0.0{
		r.Y = float32(math.Tan(float64(-*mat[Idx(2, 0)]) / float64(*mat[Idx(2, 2)])))
		r.Z = float32(math.Tan(float64(-*mat[Idx(0, 1)]) / float64(*mat[Idx(1, 1)])))
	}else{
		r.Y = 0.0
		r.Z = float32(math.Tan(float64(*mat[Idx(1, 0)]) / float64(*mat[Idx(0, 0)])))
	}

	return &r
}

/*func (this *MatrixF) FullInverse(){
}*/

func (this *MatrixF) IsAffine() bool{
	// An affine transform is defined by the following structure
	//
	// [ X X X P ]
	// [ X X X P ]
	// [ X X X P ]
	// [ 0 0 0 1 ]
	//
	//  Where X is an orthonormal 3x3 submatrix and P is an arbitrary translation
	//  We'll check in the following order:
	//   1: [3][3] must be 1
	//   2: Shear portion must be zero
	//   3: Dot products of rows and columns must be zero
	//   4: Length of rows and columns must be 1
	//

	if  this[Idx(3,3)] != 1.0{
		return false
	}

	if (this[Idx(0,3)] != 0.0 ||
		this[Idx(1,3)] != 0.0 ||
		this[Idx(2,3)] != 0.0){
			return false
	}

	var one, two, three Point3F
	this.GetColumn(0, &one)
	this.GetColumn(1, &two)
	this.GetColumn(2, &three)
	if (DotPP(one, two)   > 0.0001 ||
		DotPP(one, three) > 0.0001 ||
		DotPP(two, three) > 0.0001){
		return false
	}

	if (math.Sqrt(float64(1.0 - one.LenSquared())) > 0.0001 ||
		math.Sqrt(float64(1.0 - two.LenSquared())) > 0.0001 ||
		math.Sqrt(float64(1.0 - three.LenSquared())) > 0.0001){
		return false
	}

	this.GetRow(0, &one)
	this.GetRow(1, &two)
	this.GetRow(2, &three)
	if (DotPP(one, two)   > 0.0001 ||
		DotPP(one, three) > 0.0001 ||
		DotPP(two, three) > 0.0001){
		return false
	}


	if (math.Abs(float64(1.0 - one.LenSquared())) > 0.0001 ||
		math.Abs(float64(1.0 - two.LenSquared())) > 0.0001 ||
		math.Abs(float64(1.0 - three.LenSquared())) > 0.0001){
		return false
	}

	// We're ok.
	return true;
}

func (this *MatrixF) TransposeTo(matrix []*float32){
	*matrix[Idx(0,0)] = this[Idx(0,0)]
	*matrix[Idx(0,1)] = this[Idx(1,0)]
	*matrix[Idx(0,2)] = this[Idx(2,0)]
	*matrix[Idx(0,3)] = this[Idx(3,0)]
	*matrix[Idx(1,0)] = this[Idx(0,1)]
	*matrix[Idx(1,1)] = this[Idx(1,1)]
	*matrix[Idx(1,2)] = this[Idx(2,1)]
	*matrix[Idx(1,3)] = this[Idx(3,1)]
	*matrix[Idx(2,0)] = this[Idx(0,2)]
	*matrix[Idx(2,1)] = this[Idx(1,2)]
	*matrix[Idx(2,2)] = this[Idx(2,2)]
	*matrix[Idx(2,3)] = this[Idx(3,2)]
	*matrix[Idx(3,0)] = this[Idx(0,3)]
	*matrix[Idx(3,1)] = this[Idx(1,3)]
	*matrix[Idx(3,2)] = this[Idx(2,3)]
	*matrix[Idx(3,3)] = this[Idx(3,3)]
}


func (this *MatrixF) ToF32() []*float32{
	return  []*float32{&this[0], &this[1], &this[2], &this[3], &this[4], &this[5], &this[6], &this[7], &this[8], &this[9], &this[10], &this[11], &this[12], &this[13], &this[14], &this[15] }
}

func F16ToF32(M [16]float32) []*float32{
	return  []*float32{&M[0], &M[1], &M[2], &M[3], &M[4], &M[5], &M[6], &M[7], &M[8], &M[9], &M[10], &M[11], &M[12], &M[13], &M[14], &M[15] }
}