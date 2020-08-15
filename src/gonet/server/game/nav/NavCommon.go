//
// Copyright (c) 2009-2010 Mikko Mononen memon@inside.org
//
// This software is provided 'as-is', without any express or implied
// warranty.  In no event will the authors be held liable for any damages
// arising from the use of this software.
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would be
//    appreciated but is not required.
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
// 3. This notice may not be removed or altered from any source distribution.
//

package detour

import (
	"math"
	"reflect"
	"unsafe"
)

/**
@defgroup detour Detour

Members in this module are used to create, manipulate, and query navigation
meshes.

@note This is a summary list of members.  Use the index or search
feature to find minor members.
*/

/// @name General helper functions
/// @{

/// Used to ignore a function parameter.  VS complains about unused parameters
/// and this silences the warning.
///  @param [in] _ Unused parameter
func DtIgnoreUnused(interface{}) {}

/// Swaps the values of the two parameters.
///  @param[in,out]	a	Value A
///  @param[in,out]	b	Value B
func DtSwapFloat32(a, b *float32) { t := *a; *a = *b; *b = t }
func DtSwapUInt32(a, b *uint32)   { t := *a; *a = *b; *b = t }
func DtSwapInt32(a, b *int32)     { t := *a; *a = *b; *b = t }
func DtSwapUInt16(a, b *uint16)   { t := *a; *a = *b; *b = t }
func DtSwapInt16(a, b *int16)     { t := *a; *a = *b; *b = t }

/// Returns the minimum of two values.
///  @param[in]		a	Value A
///  @param[in]		b	Value B
///  @return The minimum of the two values.
func DtMinFloat32(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}
func DtMinUInt32(a, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}
func DtMinInt32(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}
func DtMinUInt16(a, b uint16) uint16 {
	if a < b {
		return a
	} else {
		return b
	}
}
func DtMinInt16(a, b int16) int16 {
	if a < b {
		return a
	} else {
		return b
	}
}

/// Returns the maximum of two values.
///  @param[in]		a	Value A
///  @param[in]		b	Value B
///  @return The maximum of the two values.
func DtMaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}
func DtMaxUInt32(a, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}
func DtMaxInt32(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}
func DtMaxUInt16(a, b uint16) uint16 {
	if a > b {
		return a
	} else {
		return b
	}
}
func DtMaxInt16(a, b int16) int16 {
	if a > b {
		return a
	} else {
		return b
	}
}

func DtMaxUInt8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func DtMaxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

/// Returns the absolute value.
///  @param[in]		a	The value.
///  @return The absolute value of the specified value.
func DtAbsFloat32(a float32) float32 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
func DtAbsInt32(a int32) int32 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}
func DtAbsInt16(a int16) int16 {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

/// Returns the square of the value.
///  @param[in]		a	The value.
///  @return The square of the value.
func DtSqrFloat32(a float32) float32 { return a * a }
func DtSqrUInt32(a uint32) uint32    { return a * a }
func DtSqrInt32(a int32) int32       { return a * a }
func DtSqrUInt16(a uint16) uint16    { return a * a }
func DtSqrInt16(a int16) int16       { return a * a }

/// Clamps the value to the specified range.
///  @param[in]		v	The value to clamp.
///  @param[in]		mn	The minimum permitted return value.
///  @param[in]		mx	The maximum permitted return value.
///  @return The value, clamped to the specified range.
func DtClampFloat32(v, mn, mx float32) float32 {
	if v < mn {
		return mn
	} else {
		if v > mx {
			return mx
		} else {
			return v
		}
	}
}
func DtClampUInt32(v, mn, mx uint32) uint32 {
	if v < mn {
		return mn
	} else {
		if v > mx {
			return mx
		} else {
			return v
		}
	}
}
func DtClampInt32(v, mn, mx int32) int32 {
	if v < mn {
		return mn
	} else {
		if v > mx {
			return mx
		} else {
			return v
		}
	}
}
func DtClampUInt16(v, mn, mx uint16) uint16 {
	if v < mn {
		return mn
	} else {
		if v > mx {
			return mx
		} else {
			return v
		}
	}
}
func DtClampInt16(v, mn, mx int16) int16 {
	if v < mn {
		return mn
	} else {
		if v > mx {
			return mx
		} else {
			return v
		}
	}
}

/// @}
/// @name Vector helper functions.
/// @{

/// Derives the cross product of two vectors. (@p v1 x @p v2)
///  @param[out]	dest	The cross product. [(x, y, z)]
///  @param[in]		v1		A Vector [(x, y, z)]
///  @param[in]		v2		A vector [(x, y, z)]
func DtVcross(dest, v1, v2 []float32) {
	dest[0] = v1[1]*v2[2] - v1[2]*v2[1]
	dest[1] = v1[2]*v2[0] - v1[0]*v2[2]
	dest[2] = v1[0]*v2[1] - v1[1]*v2[0]
}

/// Derives the dot product of two vectors. (@p v1 . @p v2)
///  @param[in]		v1	A Vector [(x, y, z)]
///  @param[in]		v2	A vector [(x, y, z)]
/// @return The dot product.
func DtVdot(v1, v2 []float32) float32 {
	return v1[0]*v2[0] + v1[1]*v2[1] + v1[2]*v2[2]
}

/// Performs a scaled vector addition. (@p v1 + (@p v2 * @p s))
///  @param[out]	dest	The result vector. [(x, y, z)]
///  @param[in]		v1		The base vector. [(x, y, z)]
///  @param[in]		v2		The vector to scale and add to @p v1. [(x, y, z)]
///  @param[in]		s		The amount to scale @p v2 by before adding to @p v1.
func DtVmad(dest, v1, v2 []float32, s float32) {
	dest[0] = v1[0] + v2[0]*s
	dest[1] = v1[1] + v2[1]*s
	dest[2] = v1[2] + v2[2]*s
}

/// Performs a linear interpolation between two vectors. (@p v1 toward @p v2)
///  @param[out]	dest	The result vector. [(x, y, x)]
///  @param[in]		v1		The starting vector.
///  @param[in]		v2		The destination vector.
///	 @param[in]		t		The interpolation factor. [Limits: 0 <= value <= 1.0]
func DtVlerp(dest, v1, v2 []float32, t float32) {
	dest[0] = v1[0] + (v2[0]-v1[0])*t
	dest[1] = v1[1] + (v2[1]-v1[1])*t
	dest[2] = v1[2] + (v2[2]-v1[2])*t
}

/// Performs a vector addition. (@p v1 + @p v2)
///  @param[out]	dest	The result vector. [(x, y, z)]
///  @param[in]		v1		The base vector. [(x, y, z)]
///  @param[in]		v2		The vector to add to @p v1. [(x, y, z)]
func DtVadd(dest, v1, v2 []float32) {
	dest[0] = v1[0] + v2[0]
	dest[1] = v1[1] + v2[1]
	dest[2] = v1[2] + v2[2]
}

/// Performs a vector subtraction. (@p v1 - @p v2)
///  @param[out]	dest	The result vector. [(x, y, z)]
///  @param[in]		v1		The base vector. [(x, y, z)]
///  @param[in]		v2		The vector to subtract from @p v1. [(x, y, z)]
func DtVsub(dest, v1, v2 []float32) {
	dest[0] = v1[0] - v2[0]
	dest[1] = v1[1] - v2[1]
	dest[2] = v1[2] - v2[2]
}

/// Scales the vector by the specified value. (@p v * @p t)
///  @param[out]	dest	The result vector. [(x, y, z)]
///  @param[in]		v		The vector to scale. [(x, y, z)]
///  @param[in]		t		The scaling factor.
func DtVscale(dest, v []float32, t float32) {
	dest[0] = v[0] * t
	dest[1] = v[1] * t
	dest[2] = v[2] * t
}

/// Selects the minimum value of each element from the specified vectors.
///  @param[in,out]	mn	A vector.  (Will be updated with the result.) [(x, y, z)]
///  @param[in]	v	A vector. [(x, y, z)]
func DtVmin(mn, v []float32) {
	mn[0] = DtMinFloat32(mn[0], v[0])
	mn[1] = DtMinFloat32(mn[1], v[1])
	mn[2] = DtMinFloat32(mn[2], v[2])
}

/// Selects the maximum value of each element from the specified vectors.
///  @param[in,out]	mx	A vector.  (Will be updated with the result.) [(x, y, z)]
///  @param[in]		v	A vector. [(x, y, z)]
func DtVmax(mx, v []float32) {
	mx[0] = DtMaxFloat32(mx[0], v[0])
	mx[1] = DtMaxFloat32(mx[1], v[1])
	mx[2] = DtMaxFloat32(mx[2], v[2])
}

/// Sets the vector elements to the specified values.
///  @param[out]	dest	The result vector. [(x, y, z)]
///  @param[in]		x		The x-value of the vector.
///  @param[in]		y		The y-value of the vector.
///  @param[in]		z		The z-value of the vector.
func DtVset(dest []float32, x, y, z float32) {
	dest[0] = x
	dest[1] = y
	dest[2] = z
}

/// Performs a vector copy.
///  @param[out]	dest	The result. [(x, y, z)]
///  @param[in]		a		The vector to copy. [(x, y, z)]
func DtVcopy(dest, a []float32) {
	dest[0] = a[0]
	dest[1] = a[1]
	dest[2] = a[2]
}

/// Derives the scalar length of the vector.
///  @param[in]		v The vector. [(x, y, z)]
/// @return The scalar length of the vector.
func DtVlen(v []float32) float32 {
	return DtMathSqrtf(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

/// Derives the square of the scalar length of the vector. (len * len)
///  @param[in]		v The vector. [(x, y, z)]
/// @return The square of the scalar length of the vector.
func DtVlenSqr(v []float32) float32 {
	return v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
}

/// Returns the distance between two points.
///  @param[in]		v1	A point. [(x, y, z)]
///  @param[in]		v2	A point. [(x, y, z)]
/// @return The distance between the two points.
func DtVdist(v1, v2 []float32) float32 {
	dx := v2[0] - v1[0]
	dy := v2[1] - v1[1]
	dz := v2[2] - v1[2]
	return DtMathSqrtf(dx*dx + dy*dy + dz*dz)
}

/// Returns the square of the distance between two points.
///  @param[in]		v1	A point. [(x, y, z)]
///  @param[in]		v2	A point. [(x, y, z)]
/// @return The square of the distance between the two points.
func DtVdistSqr(v1, v2 []float32) float32 {
	dx := v2[0] - v1[0]
	dy := v2[1] - v1[1]
	dz := v2[2] - v1[2]
	return dx*dx + dy*dy + dz*dz
}

/// Derives the distance between the specified points on the xz-plane.
///  @param[in]		v1	A point. [(x, y, z)]
///  @param[in]		v2	A point. [(x, y, z)]
/// @return The distance between the point on the xz-plane.
///
/// The vectors are projected onto the xz-plane, so the y-values are ignored.
func DtVdist2D(v1, v2 []float32) float32 {
	dx := v2[0] - v1[0]
	dz := v2[2] - v1[2]
	return DtMathSqrtf(dx*dx + dz*dz)
}

/// Derives the square of the distance between the specified points on the xz-plane.
///  @param[in]		v1	A point. [(x, y, z)]
///  @param[in]		v2	A point. [(x, y, z)]
/// @return The square of the distance between the point on the xz-plane.
func DtVdist2DSqr(v1, v2 []float32) float32 {
	dx := v2[0] - v1[0]
	dz := v2[2] - v1[2]
	return dx*dx + dz*dz
}

/// Normalizes the vector.
///  @param[in,out]	v	The vector to normalize. [(x, y, z)]
func DtVnormalize(v []float32) {
	d := 1.0 / DtMathSqrtf(DtSqrFloat32(v[0])+DtSqrFloat32(v[1])+DtSqrFloat32(v[2]))
	v[0] *= d
	v[1] *= d
	v[2] *= d
}

var thr float32 = DtSqrFloat32(1.0 / 16384.0)

/// Performs a 'sloppy' colocation check of the specified points.
///  @param[in]		p0	A point. [(x, y, z)]
///  @param[in]		p1	A point. [(x, y, z)]
/// @return True if the points are considered to be at the same location.
///
/// Basically, this function will return true if the specified points are
/// close enough to eachother to be considered colocated.
func DtVequal(p0, p1 []float32) bool {
	d := DtVdistSqr(p0, p1)
	return d < thr
}

/// Derives the dot product of two vectors on the xz-plane. (@p u . @p v)
///  @param[in]		u		A vector [(x, y, z)]
///  @param[in]		v		A vector [(x, y, z)]
/// @return The dot product on the xz-plane.
///
/// The vectors are projected onto the xz-plane, so the y-values are ignored.
func DtVdot2D(u, v []float32) float32 {
	return u[0]*v[0] + u[2]*v[2]
}

/// Derives the xz-plane 2D perp product of the two vectors. (uz*vx - ux*vz)
///  @param[in]		u		The LHV vector [(x, y, z)]
///  @param[in]		v		The RHV vector [(x, y, z)]
/// @return The dot product on the xz-plane.
///
/// The vectors are projected onto the xz-plane, so the y-values are ignored.
func DtVperp2D(u, v []float32) float32 {
	return u[2]*v[0] - u[0]*v[2]
}

/// @}
/// @name Computational geometry helper functions.
/// @{

/// Derives the signed xz-plane area of the triangle ABC, or the relationship of line AB to point C.
///  @param[in]		a		Vertex A. [(x, y, z)]
///  @param[in]		b		Vertex B. [(x, y, z)]
///  @param[in]		c		Vertex C. [(x, y, z)]
/// @return The signed xz-plane area of the triangle.
func DtTriArea2D(a, b, c []float32) float32 {
	abx := b[0] - a[0]
	abz := b[2] - a[2]
	acx := c[0] - a[0]
	acz := c[2] - a[2]
	return acx*abz - abx*acz
}

/// Determines if two axis-aligned bounding boxes overlap.
///  @param[in]		amin	Minimum bounds of box A. [(x, y, z)]
///  @param[in]		amax	Maximum bounds of box A. [(x, y, z)]
///  @param[in]		bmin	Minimum bounds of box B. [(x, y, z)]
///  @param[in]		bmax	Maximum bounds of box B. [(x, y, z)]
/// @return True if the two AABB's overlap.
/// @see dtOverlapBounds
func DtOverlapQuantBounds(amin, amax, bmin, bmax []uint16) bool {
	return !(amin[0] > bmax[0] || amax[0] < bmin[0] || amin[1] > bmax[1] || amax[1] < bmin[1] || amin[2] > bmax[2] || amax[2] < bmin[2])
}

/// Determines if two axis-aligned bounding boxes overlap.
///  @param[in]		amin	Minimum bounds of box A. [(x, y, z)]
///  @param[in]		amax	Maximum bounds of box A. [(x, y, z)]
///  @param[in]		bmin	Minimum bounds of box B. [(x, y, z)]
///  @param[in]		bmax	Maximum bounds of box B. [(x, y, z)]
/// @return True if the two AABB's overlap.
/// @see dtOverlapQuantBounds
func DtOverlapBounds(amin, amax, bmin, bmax []float32) bool {
	return !(amin[0] > bmax[0] || amax[0] < bmin[0] || amin[1] > bmax[1] || amax[1] < bmin[1] || amin[2] > bmax[2] || amax[2] < bmin[2])
}

/// Derives the closest point on a triangle from the specified reference point.
///  @param[out]	closest	The closest point on the triangle.
///  @param[in]		p		The reference point from which to test. [(x, y, z)]
///  @param[in]		a		Vertex A of triangle ABC. [(x, y, z)]
///  @param[in]		b		Vertex B of triangle ABC. [(x, y, z)]
///  @param[in]		c		Vertex C of triangle ABC. [(x, y, z)]
func DtClosestPtPointTriangle(closest, p, a, b, c []float32) {
	// Check if P in vertex region outside A
	ab := [3]float32{}
	ac := [3]float32{}
	ap := [3]float32{}
	DtVsub(ab[:], b, a)
	DtVsub(ac[:], c, a)
	DtVsub(ap[:], p, a)
	d1 := DtVdot(ab[:], ap[:])
	d2 := DtVdot(ac[:], ap[:])
	if d1 <= 0.0 && d2 <= 0.0 {
		// barycentric coordinates (1,0,0)
		DtVcopy(closest, a)
		return
	}

	// Check if P in vertex region outside B
	bp := [3]float32{}
	DtVsub(bp[:], p, b)
	d3 := DtVdot(ab[:], bp[:])
	d4 := DtVdot(ac[:], bp[:])
	if d3 >= 0.0 && d4 <= d3 {
		// barycentric coordinates (0,1,0)
		DtVcopy(closest, b)
		return
	}

	// Check if P in edge region of AB, if so return projection of P onto AB
	vc := d1*d4 - d3*d2
	if vc <= 0.0 && d1 >= 0.0 && d3 <= 0.0 {
		// barycentric coordinates (1-v,v,0)
		v := d1 / (d1 - d3)
		closest[0] = a[0] + v*ab[0]
		closest[1] = a[1] + v*ab[1]
		closest[2] = a[2] + v*ab[2]
		return
	}

	// Check if P in vertex region outside C
	cp := [3]float32{}
	DtVsub(cp[:], p, c)
	d5 := DtVdot(ab[:], cp[:])
	d6 := DtVdot(ac[:], cp[:])
	if d6 >= 0.0 && d5 <= d6 {
		// barycentric coordinates (0,0,1)
		DtVcopy(closest, c)
		return
	}

	// Check if P in edge region of AC, if so return projection of P onto AC
	vb := d5*d2 - d1*d6
	if vb <= 0.0 && d2 >= 0.0 && d6 <= 0.0 {
		// barycentric coordinates (1-w,0,w)
		w := d2 / (d2 - d6)
		closest[0] = a[0] + w*ac[0]
		closest[1] = a[1] + w*ac[1]
		closest[2] = a[2] + w*ac[2]
		return
	}

	// Check if P in edge region of BC, if so return projection of P onto BC
	va := d3*d6 - d5*d4
	if va <= 0.0 && (d4-d3) >= 0.0 && (d5-d6) >= 0.0 {
		// barycentric coordinates (0,1-w,w)
		w := (d4 - d3) / ((d4 - d3) + (d5 - d6))
		closest[0] = b[0] + w*(c[0]-b[0])
		closest[1] = b[1] + w*(c[1]-b[1])
		closest[2] = b[2] + w*(c[2]-b[2])
		return
	}

	// P inside face region. Compute Q through its barycentric coordinates (u,v,w)
	denom := 1.0 / (va + vb + vc)
	v := vb * denom
	w := vc * denom
	closest[0] = a[0] + ab[0]*v + ac[0]*w
	closest[1] = a[1] + ab[1]*v + ac[1]*w
	closest[2] = a[2] + ab[2]*v + ac[2]*w
}

var EPS float32 = 1e-4

/// Derives the y-axis height of the closest point on the triangle from the specified reference point.
///  @param[in]		p		The reference point from which to test. [(x, y, z)]
///  @param[in]		a		Vertex A of triangle ABC. [(x, y, z)]
///  @param[in]		b		Vertex B of triangle ABC. [(x, y, z)]
///  @param[in]		c		Vertex C of triangle ABC. [(x, y, z)]
///  @param[out]	h		The resulting height.
func DtClosestHeightPointTriangle(p, a, b, c []float32, h *float32) bool {
	v0 := [3]float32{}
	v1 := [3]float32{}
	v2 := [3]float32{}
	DtVsub(v0[:], c, a)
	DtVsub(v1[:], b, a)
	DtVsub(v2[:], p, a)

	dot00 := DtVdot2D(v0[:], v0[:])
	dot01 := DtVdot2D(v0[:], v1[:])
	dot02 := DtVdot2D(v0[:], v2[:])
	dot11 := DtVdot2D(v1[:], v1[:])
	dot12 := DtVdot2D(v1[:], v2[:])

	// Compute barycentric coordinates
	invDenom := 1.0 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	// The (sloppy) epsilon is needed to allow to get height of points which
	// are interpolated along the edges of the triangles.
	//	static const float EPS = 1e-4f;

	// If point lies inside the triangle, return interpolated ycoord.
	if u >= -EPS && v >= -EPS && (u+v) <= 1+EPS {
		*h = a[1] + v0[1]*u + v1[1]*v
		return true
	}

	return false
}

func DtIntersectSegmentPoly2D(p0, p1, verts []float32, nverts int, tmin, tmax *float32, segMin, segMax *int) bool {
	*tmin = 0
	*tmax = 1
	*segMin = -1
	*segMax = -1

	dir := [3]float32{}
	DtVsub(dir[:], p1, p0)

	for i, j := 0, nverts-1; i < nverts; j, i = i, i+1 {
		edge := [3]float32{}
		diff := [3]float32{}
		DtVsub(edge[:], verts[i*3:], verts[j*3:])
		DtVsub(diff[:], p0, verts[j*3:])
		n := DtVperp2D(edge[:], diff[:])
		d := DtVperp2D(dir[:], edge[:])
		if math.Abs(float64(d)) < 0.00000001 {
			// S is nearly parallel to this edge
			if n < 0 {
				return false
			} else {
				continue
			}
		}
		t := n / d
		if d < 0 {
			// segment S is entering across this edge
			if t > *tmin {
				*tmin = t
				*segMin = j
				// S enters after leaving polygon
				if *tmin > *tmax {
					return false
				}
			}
		} else {
			// segment S is leaving across this edge
			if t < *tmax {
				*tmax = t
				*segMax = j
				// S leaves before entering polygon
				if *tmax < *tmin {
					return false
				}
			}
		}
	}

	return true
}

func vperpXZ(a, b []float32) float32 {
	return a[0]*b[2] - a[2]*b[0]
}

func DtIntersectSegSeg2D(ap, aq, bp, bq []float32, s, t *float32) bool {
	u := [3]float32{}
	v := [3]float32{}
	w := [3]float32{}
	DtVsub(u[:], aq, ap)
	DtVsub(v[:], bq, bp)
	DtVsub(w[:], ap, bp)
	d := vperpXZ(u[:], v[:])
	if math.Abs(float64(d)) < 1e-6 {
		return false
	}
	*s = vperpXZ(v[:], w[:]) / d
	*t = vperpXZ(u[:], w[:]) / d
	return true
}

/// Determines if the specified point is inside the convex polygon on the xz-plane.
///  @param[in]		pt		The point to check. [(x, y, z)]
///  @param[in]		verts	The polygon vertices. [(x, y, z) * @p nverts]
///  @param[in]		nverts	The number of vertices. [Limit: >= 3]
/// @return True if the point is inside the polygon.
func DtPointInPolygon(pt, verts []float32, nverts int) bool {
	var i, j int
	c := false
	for i, j = 0, nverts-1; i < nverts; j, i = i, i+1 {
		vi := verts[i*3:]
		vj := verts[j*3:]
		if ((vi[2] > pt[2]) != (vj[2] > pt[2])) &&
			(pt[0] < (vj[0]-vi[0])*(pt[2]-vi[2])/(vj[2]-vi[2])+vi[0]) {
			c = !c
		}
	}
	return c
}

func DtDistancePtSegSqr2D(pt, p, q []float32, t *float32) float32 {
	pqx := q[0] - p[0]
	pqz := q[2] - p[2]
	dx := pt[0] - p[0]
	dz := pt[2] - p[2]
	d := pqx*pqx + pqz*pqz
	*t = pqx*dx + pqz*dz
	if d > 0 {
		*t /= d
	}
	if *t < 0 {
		*t = 0
	} else if *t > 1 {
		*t = 1
	}
	dx = p[0] + (*t)*pqx - pt[0]
	dz = p[2] + (*t)*pqz - pt[2]
	return dx*dx + dz*dz
}

func DtDistancePtPolyEdgesSqr(pt, verts []float32, nverts int, ed, et []float32) bool {
	var i, j int
	c := false
	for i, j = 0, nverts-1; i < nverts; j, i = i, i+1 {
		vi := verts[i*3:]
		vj := verts[j*3:]
		if ((vi[2] > pt[2]) != (vj[2] > pt[2])) &&
			(pt[0] < (vj[0]-vi[0])*(pt[2]-vi[2])/(vj[2]-vi[2])+vi[0]) {
			c = !c
		}
		ed[j] = DtDistancePtSegSqr2D(pt, vj, vi, &et[j])
	}
	return c
}

/// Derives the centroid of a convex polygon.
///  @param[out]	tc		The centroid of the polgyon. [(x, y, z)]
///  @param[in]		idx		The polygon indices. [(vertIndex) * @p nidx]
///  @param[in]		nidx	The number of indices in the polygon. [Limit: >= 3]
///  @param[in]		verts	The polygon vertices. [(x, y, z) * vertCount]
func DtCalcPolyCenter(tc []float32, idx []uint16, nidx int, verts []float32) {
	tc[0] = 0.0
	tc[1] = 0.0
	tc[2] = 0.0
	for j := 0; j < nidx; j++ {
		v := verts[idx[j]*3:]
		tc[0] += v[0]
		tc[1] += v[1]
		tc[2] += v[2]
	}
	s := 1.0 / float32(nidx)
	tc[0] *= s
	tc[1] *= s
	tc[2] *= s
}

func projectPoly(axis, poly []float32, npoly int, rmin, rmax *float32) {
	*rmax = DtVdot2D(axis, poly)
	*rmin = *rmax
	for i := 1; i < npoly; i++ {
		d := DtVdot2D(axis, poly[i*3:])
		*rmin = DtMinFloat32(*rmin, d)
		*rmax = DtMaxFloat32(*rmax, d)
	}
}

func overlapRange(amin, amax, bmin, bmax, eps float32) bool {
	return !((amin+eps) > bmax || (amax-eps) < bmin)
}

/// Determines if the two convex polygons overlap on the xz-plane.
///  @param[in]		polya		Polygon A vertices.	[(x, y, z) * @p npolya]
///  @param[in]		npolya		The number of vertices in polygon A.
///  @param[in]		polyb		Polygon B vertices.	[(x, y, z) * @p npolyb]
///  @param[in]		npolyb		The number of vertices in polygon B.
/// @return True if the two polygons overlap.
func DtOverlapPolyPoly2D(polya []float32, npolya int, polyb []float32, npolyb int) bool {
	for i, j := 0, npolya-1; i < npolya; j, i = i, i+1 {
		va := polya[j*3:]
		vb := polya[i*3:]
		n := [3]float32{vb[2] - va[2], 0, -(vb[0] - va[0])}
		var amin, amax, bmin, bmax float32
		projectPoly(n[:], polya, npolya, &amin, &amax)
		projectPoly(n[:], polyb, npolyb, &bmin, &bmax)
		if !overlapRange(amin, amax, bmin, bmax, EPS) {
			// Found separating axis
			return false
		}
	}
	for i, j := 0, npolyb-1; i < npolyb; j, i = i, i+1 {
		va := polyb[j*3:]
		vb := polyb[i*3:]
		n := [3]float32{vb[2] - va[2], 0, -(vb[0] - va[0])}
		var amin, amax, bmin, bmax float32
		projectPoly(n[:], polya, npolya, &amin, &amax)
		projectPoly(n[:], polyb, npolyb, &bmin, &bmax)
		if !overlapRange(amin, amax, bmin, bmax, EPS) {
			// Found separating axis
			return false
		}
	}
	return true
}

/// @}
/// @name Miscellanious functions.
/// @{

func DtNextPow2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func DtIlog2(v uint32) uint32 {
	var r, shift, temp uint32

	if v > 0xffff {
		temp = 1
	} else {
		temp = 0
	}
	r = temp << 4
	v >>= r

	if v > 0xff {
		temp = 1
	} else {
		temp = 0
	}
	shift = temp << 3
	v >>= shift
	r |= shift

	if v > 0xf {
		temp = 1
	} else {
		temp = 0
	}
	shift = temp << 2
	v >>= shift
	r |= shift

	if v > 0x3 {
		temp = 1
	} else {
		temp = 0
	}
	shift = temp << 1
	v >>= shift
	r |= shift
	r |= (v >> 1)
	return r
}

func DtAlign4(x int) int { return (x + 3) & ^3 }

func DtOppositeTile(side int) int { return (side + 4) & 0x7 }

func DtSwapByte(a, b *uint8) {
	tmp := *a
	*a = *b
	*b = tmp
}

func DtSwapEndianUInt16(v *uint16) {
	x0 := (*uint8)(unsafe.Pointer(v))
	x1 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 1))
	DtSwapByte(x0, x1)
}

func DtSwapEndianInt16(v *int16) {
	x0 := (*uint8)(unsafe.Pointer(v))
	x1 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 1))
	DtSwapByte(x0, x1)
}

func DtSwapEndianUInt32(v *uint32) {
	x0 := (*uint8)(unsafe.Pointer(v))
	x1 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 1))
	x2 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 2))
	x3 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 3))
	DtSwapByte(x0, x3)
	DtSwapByte(x1, x2)
}

func DtSwapEndianInt32(v *int32) {
	x0 := (*uint8)(unsafe.Pointer(v))
	x1 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 1))
	x2 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 2))
	x3 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 3))
	DtSwapByte(x0, x3)
	DtSwapByte(x1, x2)
}

func DtSwapEndianFloat32(v *float32) {
	x0 := (*uint8)(unsafe.Pointer(v))
	x1 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 1))
	x2 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 2))
	x3 := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + 3))
	DtSwapByte(x0, x3)
	DtSwapByte(x1, x2)
}

// Returns a random point in a convex polygon.
// Adapted from Graphics Gems article.
func DtRandomPointInConvexPoly(pts []float32, npts int, areas []float32, s, t float32, out []float32) {
	areasum := float32(0.0)
	for i := 2; i < npts; i++ {
		areas[i] = DtTriArea2D(pts[0:], pts[(i-1)*3:], pts[i*3:])
		areasum += DtMaxFloat32(float32(0.001), areas[i])
	}
	// Find sub triangle weighted by area.
	thr := s * areasum
	acc := float32(0.0)
	u := float32(1.0)
	tri := npts - 1
	for i := 2; i < npts; i++ {
		dacc := areas[i]
		if thr >= acc && thr < (acc+dacc) {
			u = (thr - acc) / dacc
			tri = i
			break
		}
		acc += dacc
	}

	v := DtMathSqrtf(t)

	a := 1 - v
	b := (1 - u) * v
	c := u * v
	pa := pts[0:]
	pb := pts[(tri-1)*3:]
	pc := pts[tri*3:]

	out[0] = a*pa[0] + b*pb[0] + c*pc[0]
	out[1] = a*pa[1] + b*pb[1] + c*pc[1]
	out[2] = a*pa[2] + b*pb[2] + c*pc[2]
}

///////////////////////////////////////////////////////////////////////////

// This section contains detailed documentation for members that don't have
// a source file. It reduces clutter in the main section of the header.

/**

@fn float dtTriArea2D(const float* a, const float* b, const float* c)
@par

The vertices are projected onto the xz-plane, so the y-values are ignored.

This is a low cost function than can be used for various purposes.  Its main purpose
is for point/line relationship testing.

In all cases: A value of zero indicates that all vertices are collinear or represent the same point.
(On the xz-plane.)

When used for point/line relationship tests, AB usually represents a line against which
the C point is to be tested.  In this case:

A positive value indicates that point C is to the left of line AB, looking from A toward B.<br/>
A negative value indicates that point C is to the right of lineAB, looking from A toward B.

When used for evaluating a triangle:

The absolute value of the return value is two times the area of the triangle when it is
projected onto the xz-plane.

A positive return value indicates:

<ul>
<li>The vertices are wrapped in the normal Detour wrap direction.</li>
<li>The triangle's 3D face normal is in the general up direction.</li>
</ul>

A negative return value indicates:

<ul>
<li>The vertices are reverse wrapped. (Wrapped opposite the normal Detour wrap direction.)</li>
<li>The triangle's 3D face normal is in the general down direction.</li>
</ul>

*/

func Memset(mem uintptr, val uint8, size int) {
	var dst []byte
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&dst)))
	sliceHeader.Cap = size
	sliceHeader.Len = size
	sliceHeader.Data = mem
	for i := 0; i < size; i++ {
		dst[i] = val
	}
}

func SliceSizeFromPointer(p, start unsafe.Pointer, eleSize uintptr) uint32 {
	return uint32((uintptr(p) - uintptr(start)) / eleSize)
}

func DtAssert(x bool) {
	if bool(x) == false {
		panic("DtAssert")
	}
}