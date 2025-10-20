// Package vec implements vector operations acting on slices.
//
// Operations usually take two forms: for example,
//   - Add(v0, v1) adds v0, v1, allocates a new vector to store the result and returns it.
//   - AddAssign(v0, v1, vOut) adds v0, v1 and writes the result to pre-allocated vOut without returning.
//
// Note that in most cases, v0, v1, and vOut can overlap.
// However, for operations that cannot, InPlace methods are implemented separately.
//
// For performance reasons, most functions in this package don't implement bound checks.
// If length mismatch happens, it may panic or produce wrong results.
package vec

import (
	"github.com/thedonutfactory/go-tfhe/math/num"
)

// Equals returns if two vectors are equal.
func Equals[T comparable](v0, v1 []T) bool {
	if len(v0) != len(v1) {
		return false
	}

	for i := range v0 {
		if v0[i] != v1[i] {
			return false
		}
	}
	return true
}

// Fill fills vector with x.
func Fill[T any](v []T, x T) {
	for i := range v {
		v[i] = x
	}
}

// Cast casts vector v of type []T1 to []T2.
func Cast[T1, T2 num.Real](v []T1) []T2 {
	vOut := make([]T2, len(v))
	CastAssign(v, vOut)
	return vOut
}

// CastAssign casts v of type []T1 to vOut of type []T2.
func CastAssign[T1, T2 num.Real](v []T1, vOut []T2) {
	for i := range vOut {
		vOut[i] = T2(v[i])
	}
}

// Rotate rotates v l times to the right.
// If l < 0, then it rotates the vector l times to the left.
// If Abs(l) > len(s), it may panic.
func Rotate[T any](v []T, l int) []T {
	vOut := make([]T, len(v))
	RotateAssign(v, l, vOut)
	return vOut
}

// RotateAssign rotates v l times to the right and writes it to vOut.
// If l < 0, then it rotates the vector l times to the left.
//
// v and vOut should not overlap. For rotating a slice inplace,
// use [vec.RotateInPlace].
func RotateAssign[T any](v []T, l int, vOut []T) {
	if l < 0 {
		l = len(v) - ((-l) % len(v))
	} else {
		l %= len(v)
	}

	CopyAssign(v, vOut[l:])
	CopyAssign(v[len(v)-l:], vOut[:l])
}

// RotateInPlace rotates v l times to the right in-place.
// If l < 0, then it rotates the vector l times to the left.
func RotateInPlace[T any](v []T, l int) {
	if l < 0 {
		l = len(v) - ((-l) % len(v))
	} else {
		l %= len(v)
	}

	ReverseInPlace(v)
	ReverseInPlace(v[:l])
	ReverseInPlace(v[l:])
}

// Reverse reverses v.
func Reverse[T any](v []T) []T {
	vOut := make([]T, len(v))
	ReverseAssign(v, vOut)
	return vOut
}

// ReverseAssign reverse v and writes it to vOut.
//
// v and vOut should not overlap. For reversing a slice inplace,
// use [vec.ReverseInPlace].
func ReverseAssign[T any](v, vOut []T) {
	for i := range vOut {
		vOut[len(vOut)-i-1] = v[i]
	}
}

// ReverseInPlace reverses v in-place.
func ReverseInPlace[T any](v []T) {
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
}

// BitReverseInPlace reorders v into bit-reversal order in-place.
func BitReverseInPlace[T any](v []T) {
	var bit, j int
	for i := 1; i < len(v); i++ {
		bit = len(v) >> 1
		for j >= bit {
			j -= bit
			bit >>= 1
		}
		j += bit
		if i < j {
			v[i], v[j] = v[j], v[i]
		}
	}
}

// Copy returns a copy of v.
func Copy[T any](v []T) []T {
	if v == nil {
		return nil
	}
	return append(make([]T, 0, len(v)), v...)
}

// CopyAssign copies v0 to v1.
func CopyAssign[T any](v0, v1 []T) {
	copy(v1, v0)
}

// Dot returns the dot product of two vectors.
func Dot[T num.Number](v0, v1 []T) T {
	var res T
	for i := range v0 {
		res += v0[i] * v1[i]
	}
	return res
}

// Add returns v0 + v1.
func Add[T num.Number](v0, v1 []T) []T {
	vOut := make([]T, len(v0))
	AddAssign(v0, v1, vOut)
	return vOut
}

// Sub returns v0 - v1.
func Sub[T num.Number](v0, v1 []T) []T {
	vOut := make([]T, len(v0))
	SubAssign(v0, v1, vOut)
	return vOut
}

// Neg returns -v0.
func Neg[T num.Number](v0 []T) []T {
	vOut := make([]T, len(v0))
	NegAssign(v0, vOut)
	return vOut
}

// NegAssign computes vOut = -v0.
func NegAssign[T num.Number](v0, vOut []T) {
	for i := range vOut {
		vOut[i] = -v0[i]
	}
}

// ScalarMul returns c * v0.
func ScalarMul[T num.Number](v0 []T, c T) []T {
	vOut := make([]T, len(v0))
	ScalarMulAssign(v0, c, vOut)
	return vOut
}

// ElementWiseMul returns v0 * v1, where * is an elementwise multiplication.
func ElementWiseMul[T num.Number](v0, v1 []T) []T {
	vOut := make([]T, len(v0))
	ElementWiseMulAssign(v0, v1, vOut)
	return vOut
}

// CmplxToFloat4 converts a complex128 vector to
// float-4 representation used in fourier polynomials.
//
// Namely, it converts
//
//	[(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
//
// to
//
//	[(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
//
// The length of the input vector should be multiple of 4.
func CmplxToFloat4(v []complex128) []float64 {
	vOut := make([]float64, 2*len(v))
	CmplxToFloat4Assign(v, vOut)
	return vOut
}

// CmplxToFloat4Assign converts a complex128 vector to
// float-4 representation used in fourier polynomials and writes it to vOut.
//
// Namely, it converts
//
//	[(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
//
// to
//
//	[(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
//
// The length of the input vector should be multiple of 4,
// and the length of vOut should be 2 times of the length of v.
func CmplxToFloat4Assign(v []complex128, vOut []float64) {
	for i, j := 0, 0; i < len(v); i, j = i+4, j+8 {
		vOut[j+0] = real(v[i+0])
		vOut[j+1] = real(v[i+1])
		vOut[j+2] = real(v[i+2])
		vOut[j+3] = real(v[i+3])

		vOut[j+4] = imag(v[i+0])
		vOut[j+5] = imag(v[i+1])
		vOut[j+6] = imag(v[i+2])
		vOut[j+7] = imag(v[i+3])
	}
}

// Float4ToCmplx converts a float-4 complex vector to
// naturally represented complex128 vector.
//
// Namely, it converts
//
//	[(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
//
// to
//
//	[(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
//
// The length of the input vector should be multiple of 8.
func Float4ToCmplx(v []float64) []complex128 {
	vOut := make([]complex128, len(v)/2)
	Float4ToCmplxAssign(v, vOut)
	return vOut
}

// Float4ToCmplxAssign converts a float-4 complex vector to
// naturally represented complex128 vector and writes it to vOut.
//
// Namely, it converts
//
//	[(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
//
// to
//
//	[(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
//
// The length of the input vector should be multiple of 8,
// and the length of vOut should be half of the length of v.
func Float4ToCmplxAssign(v []float64, vOut []complex128) {
	for i, j := 0, 0; i < len(v); i, j = i+8, j+4 {
		vOut[j+0] = complex(v[i+0], v[i+4])
		vOut[j+1] = complex(v[i+1], v[i+5])
		vOut[j+2] = complex(v[i+2], v[i+6])
		vOut[j+3] = complex(v[i+3], v[i+7])
	}
}
