// Package poly implements polynomial and its operations.
package poly

import (
	"math"

	"github.com/thedonutfactory/go-tfhe/math/num"
	"github.com/thedonutfactory/go-tfhe/math/vec"
)

// Poly is a polynomial over Z_Q[X]/(X^N + 1).
type Poly[T num.Integer] struct {
	Coeffs []T
}

// NewPoly creates a polynomial with degree N with empty coefficients.
//
// Panics when N is not a power of two, or when N is smaller than MinDegree or larger than MaxDegree.
func NewPoly[T num.Integer](N int) Poly[T] {
	switch {
	case !num.IsPowerOfTwo(N):
		panic("degree not power of two")
	case N < MinDegree:
		panic("degree smaller than MinDegree")
	}

	return Poly[T]{Coeffs: make([]T, N)}
}

// From creates a new polynomial from given coefficient slice.
// The given slice is copied, and extended to degree N.
func From[T num.Integer](coeffs []T, N int) Poly[T] {
	p := NewPoly[T](N)
	vec.CopyAssign(coeffs, p.Coeffs)
	return p
}

// Copy returns a copy of the polynomial.
func (p Poly[T]) Copy() Poly[T] {
	return Poly[T]{Coeffs: vec.Copy(p.Coeffs)}
}

// CopyFrom copies p0 to p.
func (p *Poly[T]) CopyFrom(p0 Poly[T]) {
	vec.CopyAssign(p0.Coeffs, p.Coeffs)
}

// Degree returns the degree of the polynomial.
// This is equivalent with length of coefficients.
func (p Poly[T]) Degree() int {
	return len(p.Coeffs)
}

// Clear clears all the coefficients to zero.
func (p Poly[T]) Clear() {
	vec.Fill(p.Coeffs, 0)
}

// Equals checks if p0 is equal with p.
func (p Poly[T]) Equals(p0 Poly[T]) bool {
	return vec.Equals(p.Coeffs, p0.Coeffs)
}

// FourierPoly is a fourier transformed polynomial over C[X]/(X^N/2 + 1).
// This corresponds to a polynomial over Z_Q[X]/(X^N + 1).
type FourierPoly struct {
	// Coeffs is represented as float-4 complex vector
	// for efficient computation.
	//
	// Namely,
	//
	//	[(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
	//
	// is represented as
	//
	//	[(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
	//
	Coeffs []float64
}

// NewFourierPoly creates a fourier polynomial with degree N/2 with empty coefficients.
//
// Panics when N is not a power of two, or when N is smaller than MinDegree.
func NewFourierPoly(N int) FourierPoly {
	switch {
	case !num.IsPowerOfTwo(N):
		panic("degree not power of two")
	case N < MinDegree:
		panic("degree smaller than MinDegree")
	}

	return FourierPoly{Coeffs: make([]float64, N)}
}

// Degree returns the (doubled) degree of the polynomial.
func (p FourierPoly) Degree() int {
	return len(p.Coeffs)
}

// Copy returns a copy of the polynomial.
func (p FourierPoly) Copy() FourierPoly {
	return FourierPoly{Coeffs: vec.Copy(p.Coeffs)}
}

// CopyFrom copies p0 to p.
func (p *FourierPoly) CopyFrom(p0 FourierPoly) {
	vec.CopyAssign(p0.Coeffs, p.Coeffs)
}

// Clear clears all the coefficients to zero.
func (p FourierPoly) Clear() {
	vec.Fill(p.Coeffs, 0)
}

// Equals checks if p0 is equal with p.
// Note that due to floating point errors,
// this function may return false even if p0 and p are equal.
func (p FourierPoly) Equals(p0 FourierPoly) bool {
	return vec.Equals(p.Coeffs, p0.Coeffs)
}

// Approx checks if p0 is approximately equal with p,
// with a difference smaller than eps.
func (p FourierPoly) Approx(p0 FourierPoly, eps float64) bool {
	for i := range p.Coeffs {
		if math.Abs(p.Coeffs[i]-p0.Coeffs[i]) > eps {
			return false
		}
	}
	return true
}
