// Package poly implements optimized polynomial operations for TFHE.
// Based on the high-performance implementation from tfhe-go.
package poly

import (
	"github.com/thedonutfactory/go-tfhe/params"
)

const (
	// MinDegree is the minimum degree of polynomial that Evaluator can handle.
	// Set to 2^4 because SIMD operations handle 4 values at a time.
	MinDegree = 1 << 4

	// splitLogBound denotes the maximum bits for polynomial multiplication.
	// This ensures failure rate less than 2^-284.
	splitLogBound = 48
)

// Poly is a polynomial over Z_Q[X]/(X^N + 1).
type Poly struct {
	Coeffs []params.Torus
}

// NewPoly creates a polynomial with degree N.
func NewPoly(N int) Poly {
	if !isPowerOfTwo(N) {
		panic("degree not power of two")
	}
	if N < MinDegree {
		panic("degree smaller than MinDegree")
	}
	return Poly{Coeffs: make([]params.Torus, N)}
}

// Degree returns the degree of the polynomial.
func (p Poly) Degree() int {
	return len(p.Coeffs)
}

// Copy returns a copy of the polynomial.
func (p Poly) Copy() Poly {
	coeffsCopy := make([]params.Torus, len(p.Coeffs))
	copy(coeffsCopy, p.Coeffs)
	return Poly{Coeffs: coeffsCopy}
}

// Clear clears all coefficients to zero.
func (p Poly) Clear() {
	for i := range p.Coeffs {
		p.Coeffs[i] = 0
	}
}

// FourierPoly is a fourier transformed polynomial over C[X]/(X^N/2 + 1).
// This corresponds to a polynomial over Z_Q[X]/(X^N + 1).
//
// Coeffs are represented as float-4 complex vector for efficient computation:
// [(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
// instead of standard [(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
type FourierPoly struct {
	Coeffs []float64
}

// NewFourierPoly creates a fourier polynomial with degree N.
func NewFourierPoly(N int) FourierPoly {
	if !isPowerOfTwo(N) {
		panic("degree not power of two")
	}
	if N < MinDegree {
		panic("degree smaller than MinDegree")
	}
	return FourierPoly{Coeffs: make([]float64, N)}
}

// Degree returns the degree of the polynomial.
func (p FourierPoly) Degree() int {
	return len(p.Coeffs)
}

// Copy returns a copy of the polynomial.
func (p FourierPoly) Copy() FourierPoly {
	coeffsCopy := make([]float64, len(p.Coeffs))
	copy(coeffsCopy, p.Coeffs)
	return FourierPoly{Coeffs: coeffsCopy}
}

// CopyFrom copies p0 to p.
func (p *FourierPoly) CopyFrom(p0 FourierPoly) {
	copy(p.Coeffs, p0.Coeffs)
}

// Clear clears all coefficients to zero.
func (p FourierPoly) Clear() {
	for i := range p.Coeffs {
		p.Coeffs[i] = 0
	}
}

// isPowerOfTwo checks if n is a power of two.
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// log2 returns the base-2 logarithm of n.
func log2(n int) int {
	if n <= 0 {
		panic("log2 of non-positive number")
	}
	log := 0
	for n > 1 {
		n >>= 1
		log++
	}
	return log
}
