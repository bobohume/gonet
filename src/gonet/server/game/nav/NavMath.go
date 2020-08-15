/**
@defgroup detour Detour

Members in this module are wrappers around the standard math library
*/

package detour

import "math"

func DtMathFabsf(x float32) float32             { return float32(math.Abs(float64(x))) }
func DtMathSqrtf(x float32) float32             { return float32(math.Sqrt(float64(x))) }
func DtMathFloorf(x float32) float32            { return float32(math.Floor(float64(x))) }
func DtMathCeilf(x float32) float32             { return float32(math.Ceil(float64(x))) }
func DtMathCosf(x float32) float32              { return float32(math.Cos(float64(x))) }
func DtMathSinf(x float32) float32              { return float32(math.Sin(float64(x))) }
func DtMathAtan2f(y float32, x float32) float32 { return float32(math.Atan2(float64(y), float64(x))) }
