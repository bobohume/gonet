package lmath

import (
	"math"
	"gonet/base"
)

func M_mulDivS32_C(a int, b int, c int) int{
	return  int(int64(a) * int64(b) / int64(c))
}

func M_catmullrom_C(t float32, p0 float32, p1 float32, p2 float32, p3 float32) float32{
	return 0.5 * ((3.0*p1 - 3.0*p2 + p3 - p0)*t*t*t +  (2.0*p0 - 5.0*p1 + 4.0*p2 - p3)*t*t+  (p2-p0)*t+  2.0*p1)
}

func M_point2F_normalize_C(p []*float32) {
	factor := 1.0 / float32(math.Sqrt(float64(*p[0]**p[0] + *p[1]**p[1])))
	*p[0] *= factor
	*p[1] *= factor
}

func M_point3F_normalize_C(p []*float32){
	squared := *p[0]**p[0] + *p[1]**p[1] + *p[2]**p[2];
	// This can happen in Container::castRay -> ForceFieldBare::castRay
	//AssertFatal(squared != 0.0, "Error, zero length vector normalized!");
	if (squared != 0.0) {
		factor := 1.0 / float32(math.Sqrt(float64(squared)))
		*p[0] *= factor
		*p[1] *= factor
		*p[2] *= factor
	} else {
		*p[0] = 0.0
		*p[1] = 0.0
		*p[2] = 1.0
	}
}

func M_point3F_normalize_f_C(p []*float32, val float32){
	factor := val / float32(math.Sqrt(float64((*p[0]**p[0] + *p[1]**p[1] + *p[2]**p[2]))))
	*p[0] *= factor
	*p[1] *= factor
	*p[2] *= factor
}

func M_point3F_interpolate_C(from []*float32, to []*float32,  factor float32, result []*float32){
	inverse := 1.0 - factor
	*result[0] = *from[0] * inverse + *to[0] * factor
	*result[1] = *from[1] * inverse + *to[1] * factor
	*result[2] = *from[2] * inverse + *to[2] * factor
}

func M_quatF_set_matF_C(x float32, y float32, z float32, w float32, m []*float32){
	qidx := func(r int, c int) int{
		return r * 4 + c
	}

	xs := x * 2.0
	ys := y * 2.0
	zs := z * 2.0
	wx := w * xs
	wy := w * ys
	wz := w * zs
	xx := x * xs
	xy := x * ys
	xz := x * zs
	yy := y * ys
	yz := y * zs
	zz := z * zs

	*m[qidx(0,0)] = 1.0 - (yy + zz)
	*m[qidx(1,0)] = xy - wz
	*m[qidx(2,0)] = xz + wy
	*m[qidx(3,0)] = 0.0
	*m[qidx(0,1)] = xy + wz
	*m[qidx(1,1)] = 1.0 - (xx + zz)
	*m[qidx(2,1)] = yz - wx
	*m[qidx(3,1)] = 0.0
	*m[qidx(0,2)] = xz - wy
	*m[qidx(1,2)] = yz + wx
	*m[qidx(2,2)] = 1.0 - (xx + yy)
	*m[qidx(3,2)] = 0.0

	*m[qidx(0,3)] = 0.0
	*m[qidx(1,3)] = 0.0
	*m[qidx(2,3)] = 0.0
	*m[qidx(3,3)] = 1.0
}

func M_matF_set_euler_point_C(e []*float32, p []*float32, result []*float32){
	M_matF_set_euler_C(e, result)
	*result[3] = *p[0]
	*result[7] = *p[1]
	*result[11]= *p[2]
}

func M_matF_identity_C(m []*float32){
	*m[0] = 1.0
	*m[1] = 0.0
	*m[2] = 0.0
	*m[3] = 0.0

	*m[4] = 0.0
	*m[5] = 1.0
	*m[6] = 0.0
	*m[7] = 0.0

	*m[8] = 0.0
	*m[9] = 0.0
	*m[10] = 1.0
	*m[11] = 0.0

	*m[12] = 0.0
	*m[13] = 0.0
	*m[14] = 0.0
	*m[15] = 1.0
}

func M_matF_set_euler_C(e []*float32, result []*float32){
	AXIS_X, AXIS_Y,  AXIS_Z  := (1<<0), (1<<1), (1<<2)
	axis := 0;
	if (*e[0] != 0.0){
		axis |= AXIS_X
	}
	if (*e[1] != 0.0) {
		axis |= AXIS_Y
	}
	if (*e[2] != 0.0) {
		axis |= AXIS_Z
	}

	switch axis{
	case 0:
		M_matF_identity_C(result)
	case AXIS_X:
		sx,  cx := float32(math.Sin(float64(*e[0]))), float32(math.Cos(float64(*e[0])))
		*result[0] = 1.0
		*result[1] = 0.0
		*result[2] = 0.0
		*result[3] = 0.0

		*result[4] = 0.0
		*result[5] = cx
		*result[6] = sx
		*result[7] = 0.0

		*result[8] = 0.0
		*result[9] = -sx
		*result[10]= cx
		*result[11]= 0.0

		*result[12]= 0.0
		*result[13]= 0.0
		*result[14]= 0.0
		*result[15]= 1.0

	case AXIS_Y:
		sy,  cy := float32(math.Sin(float64(*e[1]))), float32(math.Cos(float64(*e[1])))
		*result[0] = cy
		*result[1] = 0.0
		*result[2] = -sy
		*result[3] = 0.0

		*result[4] = 0.0
		*result[5] = 1.0
		*result[6] = 0.0
		*result[7] = 0.0

		*result[8] = sy
		*result[9] = 0.0
		*result[10]= cy
		*result[11]= 0.0

		*result[12]= 0.0
		*result[13]= 0.0
		*result[14]= 0.0
		*result[15]= 1.0

	case AXIS_Z:
	// the matrix looks like this:
	//  r1 - (r4 * sin(x))     r2 + (r3 * sin(x))   -cos(x) * sin(y)
	//  -cos(x) * sin(z)       cos(x) * cos(z)      sin(x)
	//  r3 + (r2 * sin(x))     r4 - (r1 * sin(x))   cos(x) * cos(y)
	//
	// where:
	//  r1 = cos(y) * cos(z)
	//  r2 = cos(y) * sin(z)
	//  r3 = sin(y) * cos(z)
	//  r4 = sin(y) * sin(z)
		sz,  cz := float32(math.Sin(float64(*e[2]))), float32(math.Cos(float64(*e[2])))

		*result[0] = cz
		*result[1] = sz
		*result[2] = 0.0
		*result[3] = 0.0

		*result[4] = -sz
		*result[5] = cz
		*result[6] = 0.0
		*result[7] = 0.0

		*result[8] = 0.0
		*result[9] = 0.0
		*result[10]= 1.0
		*result[11]= 0.0

		*result[12]= 0.0
		*result[13]= 0.0
		*result[14]= 0.0
		*result[15]= 1.0

	default:
	// the matrix looks like this:
	//  r1 - (r4 * sin(x))     r2 + (r3 * sin(x))   -cos(x) * sin(y)
	//  -cos(x) * sin(z)       cos(x) * cos(z)      sin(x)
	//  r3 + (r2 * sin(x))     r4 - (r1 * sin(x))   cos(x) * cos(y)
	//
	// where:
	//  r1 = cos(y) * cos(z)
	//  r2 = cos(y) * sin(z)
	//  r3 = sin(y) * cos(z)
	//  r4 = sin(y) * sin(z)
		sx,  cx := float32(math.Sin(float64(*e[0]))), float32(math.Cos(float64(*e[0])))
		sy,  cy := float32(math.Sin(float64(*e[1]))), float32(math.Cos(float64(*e[1])))
		sz,  cz := float32(math.Sin(float64(*e[2]))), float32(math.Cos(float64(*e[2])))
		r1 := cy * cz;
		r2 := cy * sz;
		r3 := sy * cz;
		r4 := sy * sz;

		*result[0] = r1 - (r4 * sx)
		*result[1] = r2 + (r3 * sx)
		*result[2] = -cx * sy
		*result[3] = 0.0

		*result[4] = -cx * sz
		*result[5] = cx * cz
		*result[6] = sx
		*result[7] = 0.0

		*result[8] = r3 + (r2 * sx)
		*result[9] = r4 - (r1 * sx)
		*result[10]= cx * cy
		*result[11]= 0.0

		*result[12]= 0.0
		*result[13]= 0.0
		*result[14]= 0.0
		*result[15]= 1.0
	}
}

func M_matF_determinant_C(m []*float32) float32{
	return *m[0] * (*m[5] * *m[10] - *m[6] * *m[9])  + *m[4] * (*m[2] * *m[9]  - *m[1] * *m[10]) + *m[8] * (*m[1] * *m[6]  - *m[2] * *m[5])
}

func m_matF_x_vectorF(m []*float32, v []*float32, x *float32, y *float32, z *float32){
	*x = *m[0]* *v[0] + *m[1]* *v[1] + *m[2]* *v[2]
	*y = *m[4]* *v[0] + *m[5]* *v[1] + *m[6]* *v[2]
	*z = *m[8]* *v[0] + *m[9]* *v[1] + *m[10]* *v[2]
}

func M_matF_inverse_C(m []*float32){
	// using Cramers Rule find the Inverse
	// Minv = (1/det(M)) * adjoint(M)
	det := M_matF_determinant_C( m )
	base.Assert( det == 0.0, "MatrixF::inverse: non-singular matrix, no inverse.")

	invDet := 1.0/det
	var temp [16]float32
	temp[0] = (*m[5] * *m[10]- *m[6] * *m[9]) * invDet
	temp[1] = (*m[9] * *m[2] - *m[10]* *m[1]) * invDet
	temp[2] = (*m[1] * *m[6] - *m[2] * *m[5]) * invDet

	temp[4] = (*m[6] * *m[8] - *m[4] * *m[10])* invDet
	temp[5] = (*m[10]* *m[0] - *m[8] * *m[2]) * invDet
	temp[6] = (*m[2] * *m[4] - *m[0] * *m[6]) * invDet

	temp[8] = (*m[4] * *m[9] - *m[5] * *m[8]) * invDet
	temp[9] = (*m[8] * *m[1] - *m[9] * *m[0]) * invDet
	temp[10]= (*m[0] * *m[5] - *m[1] * *m[4]) * invDet

	*m[0] = temp[0]
	*m[1] = temp[1]
	*m[2] = temp[2]

	*m[4] = temp[4]
	*m[5] = temp[5]
	*m[6] = temp[6]

	*m[8] = temp[8]
	*m[9] = temp[9]
	*m[10] = temp[10]

	// invert the translation
	temp[0] = -*m[3]
	temp[1] = -*m[7]
	temp[2] = -*m[11]
	m_matF_x_vectorF(m, F16ToF32(temp), &temp[4], &temp[5], &temp[6])
	*m[3] = temp[4]
	*m[7] = temp[5]
	*m[11]= temp[6]
}

func M_matF_affineInverse_C(m []*float32){
	// Matrix class checks to make sure this is an affine transform before calling
	//  this function, so we can proceed assuming it is...
	var temp[16] float32
	for i, v := range m{
		temp[i] = *v
	}

	// Transpose rotation
	*m[1] = temp[4]
	*m[4] = temp[1]
	*m[2] = temp[8]
	*m[8] = temp[2]
	*m[6] = temp[9]
	*m[9] = temp[6]

	*m[3]  = -(temp[0]*temp[3] + temp[4]*temp[7] + temp[8]*temp[11])
	*m[7]  = -(temp[1]*temp[3] + temp[5]*temp[7] + temp[9]*temp[11])
	*m[11] = -(temp[2]*temp[3] + temp[6]*temp[7] + temp[10]*temp[11])
}

func swap(a *float32, b *float32){
	temp := *a
	*a = *b
	*b = temp
}

func M_matF_transpose_C(m []*float32){
	swap(m[1], m[4])
	swap(m[2], m[8])
	swap(m[3], m[12])
	swap(m[6], m[9])
	swap(m[7], m[13])
	swap(m[11],m[14])
}

func M_matF_scale_C(m []*float32, p []*float32){
	// Note, doesn't allow scaling w...
	*m[0]  *= *p[0]
	*m[1]  *= *p[1]
	*m[2]  *= *p[2]
	*m[4]  *= *p[0]
	*m[5]  *= *p[1]
	*m[6]  *= *p[2]
	*m[8]  *= *p[0]
	*m[9]  *= *p[1]
	*m[10] *= *p[2]
	*m[12] *= *p[0]
	*m[13] *= *p[1]
	*m[14] *= *p[2]
}

func M_matF_normalize_C(m []*float32){
	var col0, col1, col2 []*float32
	col0, col1, col2 = make([]*float32, 3), make([]*float32, 3), make([]*float32, 3)
	// extract columns 0 and 1
	*col0[0] = *m[0]
	*col0[1] = *m[4]
	*col0[2] = *m[8]

	*col1[0] = *m[1]
	*col1[1] = *m[5]
	*col1[2] = *m[9]

	// assure their relationsips to one another
	CrossFFF(col0, col1, col2)
	CrossFFF(col2, col0, col1)

	// assure their lengh is 1.0f
	M_point3F_normalize_C( col0 )
	M_point3F_normalize_C( col1 )
	M_point3F_normalize_C( col2 )

	// store the normalized columns
	*m[0] = *col0[0]
	*m[4] = *col0[1]
	*m[8] = *col0[2]

	*m[1] = *col1[0]
	*m[5] = *col1[1]
	*m[9] = *col1[2]

	*m[2] = *col2[0]
	*m[6] = *col2[1]
	*m[10]= *col2[2]
}

func Default_matF_x_matF_C(a []*float32, b []*float32, mresult []*float32){
	*mresult[0] = *a[0]* *b[0] + *a[1]* *b[4] + *a[2]* *b[8]  + *a[3]* *b[12]
	*mresult[1] = *a[0]* *b[1] + *a[1]* *b[5] + *a[2]* *b[9]  + *a[3]* *b[13]
	*mresult[2] = *a[0]* *b[2] + *a[1]* *b[6] + *a[2]* *b[10] + *a[3]* *b[14]
	*mresult[3] = *a[0]* *b[3] + *a[1]* *b[7] + *a[2]* *b[11] + *a[3]* *b[15]

	*mresult[4] = *a[4]* *b[0] + *a[5]* *b[4] + *a[6]* *b[8]  + *a[7]* *b[12]
	*mresult[5] = *a[4]* *b[1] + *a[5]* *b[5] + *a[6]* *b[9]  + *a[7]* *b[13]
	*mresult[6] = *a[4]* *b[2] + *a[5]* *b[6] + *a[6]* *b[10] + *a[7]* *b[14]
	*mresult[7] = *a[4]* *b[3] + *a[5]* *b[7] + *a[6]* *b[11] + *a[7]* *b[15]

	*mresult[8] = *a[8]* *b[0] + *a[9]* *b[4] + *a[10]* *b[8] + *a[11]* *b[12]
	*mresult[9] = *a[8]* *b[1] + *a[9]* *b[5] + *a[10]* *b[9] + *a[11]* *b[13]
	*mresult[10]= *a[8]* *b[2] + *a[9]* *b[6] + *a[10]* *b[10]+ *a[11]* *b[14]
	*mresult[11]= *a[8]* *b[3] + *a[9]* *b[7] + *a[10]* *b[11]+ *a[11]* *b[15]

	*mresult[12]= *a[12]* *b[0]+ *a[13]* *b[4]+ *a[14]* *b[8] + *a[15]* *b[12]
	*mresult[13]= *a[12]* *b[1]+ *a[13]* *b[5]+ *a[14]* *b[9] + *a[15]* *b[13]
	*mresult[14]= *a[12]* *b[2]+ *a[13]* *b[6]+ *a[14]* *b[10]+ *a[15]* *b[14]
	*mresult[15]= *a[12]* *b[3]+ *a[13]* *b[7]+ *a[14]* *b[11]+ *a[15]* *b[15]
}

func M_matF_x_point4F_C(m []*float32, p []*float32, presult []*float32){
	*presult[0] = *m[0]* *p[0] + *m[1]* *p[1] + *m[2]* *p[2]  + *m[3]* *p[3]
	*presult[1] = *m[4]* *p[0] + *m[5]* *p[1] + *m[6]* *p[2]  + *m[7]* *p[3]
	*presult[2] = *m[8]* *p[0] + *m[9]* *p[1] + *m[10]* *p[2] + *m[11]* *p[3]
	*presult[3] = *m[12]* *p[0]+ *m[13]* *p[1]+ *m[14]* *p[2] + *m[15]* *p[3]
}

func M_matF_x_point3F_C(m []*float32, p []*float32, presult []*float32){
	*presult[0] = *m[0]* *p[0] + *m[1]* *p[1] + *m[2]* *p[2]  + *m[3]
	*presult[1] = *m[4]* *p[0] + *m[5]* *p[1] + *m[6]* *p[2]  + *m[7]
	*presult[2] = *m[8]* *p[0] + *m[9]* *p[1] + *m[10]* *p[2] + *m[11]
}

func M_matF_x_scale_x_planeF_C(m []*float32, s []*float32, p []*float32, presult []*float32){
	// We take in a matrix, a scale factor, and a plane equation.  We want to output
	//  the resultant normal
	// We have T = m*s
	// To multiply the normal, we want Inv(Tr(m*s))
	//  Inv(Tr(ms)) = Inv(Tr(s) * Tr(m))
	//              = Inv(Tr(m)) * Inv(Tr(s))
	//
	//  Inv(Tr(s)) = Inv(s) = [ 1/x   0   0  0]
	//                        [   0 1/y   0  0]
	//                        [   0   0 1/z  0]
	//                        [   0   0   0  1]
	//
	// Since m is an affine matrix,
	//  Tr(m) = [ [       ] 0 ]
	//          [ [   R   ] 0 ]
	//          [ [       ] 0 ]
	//          [ [ x y z ] 1 ]
	//
	// Inv(Tr(m)) = [ [    -1 ] 0 ]
	//              [ [   R   ] 0 ]
	//              [ [       ] 0 ]
	//              [ [ A B C ] 1 ]
	// Where:
	//
	//  P = (x, y, z)
	//  A = -(Row(0, r) * P);
	//  B = -(Row(1, r) * P);
	//  C = -(Row(2, r) * P);

	/*var invScale MatrixF
	invScale.Identity()

	pScaleElems := invScale.ToF32()
	*pScaleElems[Idx(0, 0)] = 1.0 / *s[0];
	*pScaleElems[Idx(1, 1)] = 1.0 / *s[1];
	*pScaleElems[Idx(2, 2)] = 1.0 / *s[2];

	shear := Point3F{*m[Idx(3, 0)], *m[Idx(3, 1)],*m[Idx(3, 2)]}

	row0 := Point3F{*m[Idx(0, 0)], *m[Idx(0, 1)], *m[Idx(0, 2)]}
	row1 := Point3F{*m[Idx(1, 0)], *m[Idx(1, 1)], *m[Idx(1, 2)]}
	row2 := Point3F{*m[Idx(2, 0)], *m[Idx(2, 1)], *m[Idx(2, 2)]}

	A := -DotPP(&row0, &shear)
	B := -DotPP(&row1, &shear)
	C := -DotPP(&row2, &shear)

	var invTrMatrix MatrixF
	invTrMatrix.Identity()
	destMat := invTrMatrix.ToF32()
	*destMat[Idx(0, 0)] = *m[Idx(0, 0)]
	*destMat[Idx(1, 0)] = *m[Idx(1, 0)]
	*destMat[Idx(2, 0)] = *m[Idx(2, 0)]
	*destMat[Idx(0, 1)] = *m[Idx(0, 1)]
	*destMat[Idx(1, 1)] = *m[Idx(1, 1)]
	*destMat[Idx(2, 1)] = *m[Idx(2, 1)]
	*destMat[Idx(0, 2)] = *m[Idx(0, 2)]
	*destMat[Idx(1, 2)] = *m[Idx(1, 2)]
	*destMat[Idx(2, 2)] = *m[Idx(2, 2)]
	*destMat[Idx(0, 3)] = A
	*destMat[Idx(1, 3)] = B
	*destMat[Idx(2, 3)] = C
	invTrMatrix.Mul(invScale)

	norm := Point3F{*p[0], *p[1], *p[2]}
	point := norm.MulF(-(*p[3]))
	invTrMatrix.mulP(norm)
	norm.normalize()

	var temp MatrixF
	ff := temp.ToF32()
	for i, v := range m{
		*ff[i] = *v
	}

	point.X *= *s[0]
	point.Y *= *s[1]
	point.Z *= *s[2]
	temp.mulP(point)

	PlaneF resultPlane(point, norm)
	presult[0] = resultPlane.x
	presult[1] = resultPlane.y
	presult[2] = resultPlane.z
	presult[3] = resultPlane.d*/
}

func M_matF_x_box3F_C(m []*float32, min []*float32, max []*float32){
	var originalMin, originalMax [3]float32
	originalMin[0] = *min[0]
	originalMin[1] = *min[1]
	originalMin[2] = *min[2]
	originalMax[0] = *max[0]
	originalMax[1] = *max[1]
	originalMax[2] = *max[2]

	*min[0], *max[0] = *m[3], *m[3]
	*min[1], *max[1] = *m[7], *m[7]
	*min[2], *max[2] = *m[11], *m[11]

	row := m
	n1, n2, n3  := 0, 0, 0
	for i := 0; i < 3; i++ {
		Do_One_Row := func(j int){
			a := *row[j + n1] * originalMin[j]
			b := *row[j + n1] * originalMax[j]
			if a < b{
				*min[n2] += a
				*max[n3] += b
			}else{
				*min[n2] += b
				*max[n3] += a
			}
		}

		// Simpler addressing (avoiding things like [ecx+edi*4]) might be worthwhile (LH):
		Do_One_Row(0)
		Do_One_Row(1)
		Do_One_Row(2)
		n1 += 4
		n2++
		n3++
	}
}

func ClampF(val float32, low float32, high float32) float32{
	return float32(math.Max(math.Min(float64(val), float64(high)), float64(low)))
}

func ClampI(val int, low int, high int) int{
	return int(math.Max(math.Min(float64(val), float64(high)), float64(low)))
}

func Abs(val int) int{
	return int(math.Abs(float64(val)))
}

func Sqrt(val float32) float32{
	return  float32(math.Sqrt(float64(val)))
}

func Atan2(a float32, b float32) float32{
	return float32(math.Atan2(float64(a), float64(b)))
}

func Max(a, b int) int{
	return int(math.Max(float64(a), float64(b)))
}

func Min(a, b int) int{
	return int(math.Min(float64(a), float64(b)))
}

func IsZero(a float32) bool{
	return math.Abs(float64(a)) < 0.000001
}