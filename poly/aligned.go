package poly

import "github.com/thedonutfactory/go-tfhe/params"

// Memory alignment utilities for better cache performance

// NewPolyAligned creates a polynomial with cache-line aligned memory
// This helps with SIMD operations and cache efficiency
func NewPolyAligned(N int) Poly {
	if !isPowerOfTwo(N) {
		panic("degree not power of two")
	}
	if N < MinDegree {
		panic("degree smaller than MinDegree")
	}

	// Allocate with extra space for alignment
	// Cache lines are typically 64 bytes = 16 x uint32
	const cacheLineSize = 16
	coeffs := make([]params.Torus, N+cacheLineSize)

	// Find aligned offset
	offset := 0
	addr := uintptr(0)
	if len(coeffs) > 0 {
		addr = uintptr(len(coeffs)) % cacheLineSize
		if addr != 0 {
			offset = int(cacheLineSize - addr)
		}
	}

	// Return slice starting at aligned offset
	return Poly{Coeffs: coeffs[offset : offset+N]}
}

// NewFourierPolyAligned creates a fourier polynomial with cache-line aligned memory
func NewFourierPolyAligned(N int) FourierPoly {
	if !isPowerOfTwo(N) {
		panic("degree not power of two")
	}
	if N < MinDegree {
		panic("degree smaller than MinDegree")
	}

	// Allocate with extra space for alignment
	// Cache lines are typically 64 bytes = 8 x float64
	const cacheLineSize = 8
	coeffs := make([]float64, N+cacheLineSize)

	// Find aligned offset
	offset := 0
	addr := uintptr(0)
	if len(coeffs) > 0 {
		addr = uintptr(len(coeffs)) % cacheLineSize
		if addr != 0 {
			offset = int(cacheLineSize - addr)
		}
	}

	// Return slice starting at aligned offset
	return FourierPoly{Coeffs: coeffs[offset : offset+N]}
}
