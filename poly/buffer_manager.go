package poly

import "github.com/thedonutfactory/go-tfhe/params"

// BufferManager centralizes all polynomial operation buffers
// This provides a single, well-documented place to manage FFT, decomposition, and rotation buffers
type BufferManager struct {
	// Polynomial degree
	n int

	// === FFT Buffers ===

	// Forward/Inverse FFT working buffers
	FFT struct {
		Poly    Poly        // Time domain working buffer
		Fourier FourierPoly // Frequency domain working buffer
	}

	// === Decomposition Buffers ===

	Decomposition struct {
		// Decomposed polynomials in time domain [level]
		Poly []Poly
		// Decomposed polynomials in Fourier domain [level]
		Fourier []FourierPoly
	}

	// === Multiplication Buffers ===

	Multiplication struct {
		// Result accumulators in Fourier domain
		AccA FourierPoly
		AccB FourierPoly
		// Temporary buffer for operations
		Temp FourierPoly
	}

	// === Rotation Buffers ===

	Rotation struct {
		// Pool of polynomials for X^k multiplication
		Pool  []Poly
		InUse int // Number currently in use

		// TRLWE rotation buffers
		TRLWEPool []*struct {
			A []params.Torus
			B []params.Torus
		}
		TRLWEInUse int
	}

	// === Temporary Buffers ===

	// General-purpose temporary buffers
	Temp struct {
		Poly1 Poly
		Poly2 Poly
		Poly3 Poly
	}
}

// NewBufferManager creates a new centralized buffer manager
func NewBufferManager(n int) *BufferManager {
	l := params.GetTRGSWLv1().L

	bm := &BufferManager{n: n}

	// Initialize FFT buffers
	bm.FFT.Poly = NewPoly(n)
	bm.FFT.Fourier = NewFourierPoly(n)

	// Initialize decomposition buffers for 2*L levels (A and B components)
	bm.Decomposition.Poly = make([]Poly, l*2)
	bm.Decomposition.Fourier = make([]FourierPoly, l*2)
	for i := 0; i < l*2; i++ {
		bm.Decomposition.Poly[i] = NewPoly(n)
		bm.Decomposition.Fourier[i] = NewFourierPoly(n)
	}

	// Initialize multiplication buffers
	bm.Multiplication.AccA = NewFourierPoly(n)
	bm.Multiplication.AccB = NewFourierPoly(n)
	bm.Multiplication.Temp = NewFourierPoly(n)

	// Initialize rotation pool (4 polynomials should be enough for most operations)
	bm.Rotation.Pool = make([]Poly, 4)
	for i := 0; i < 4; i++ {
		bm.Rotation.Pool[i] = NewPoly(n)
	}
	bm.Rotation.InUse = 0

	// Initialize TRLWE rotation pool
	bm.Rotation.TRLWEPool = make([]*struct {
		A []params.Torus
		B []params.Torus
	}, 4)
	for i := 0; i < 4; i++ {
		bm.Rotation.TRLWEPool[i] = &struct {
			A []params.Torus
			B []params.Torus
		}{
			A: make([]params.Torus, n),
			B: make([]params.Torus, n),
		}
	}
	bm.Rotation.TRLWEInUse = 0

	// Initialize temporary buffers
	bm.Temp.Poly1 = NewPoly(n)
	bm.Temp.Poly2 = NewPoly(n)
	bm.Temp.Poly3 = NewPoly(n)

	return bm
}

// GetRotationBuffer returns a polynomial buffer for rotation operations
func (bm *BufferManager) GetRotationBuffer() Poly {
	if bm.Rotation.InUse >= len(bm.Rotation.Pool) {
		// Wrap around if we run out (should rarely happen)
		bm.Rotation.InUse = 0
	}
	buffer := bm.Rotation.Pool[bm.Rotation.InUse]
	bm.Rotation.InUse++
	return buffer
}

// GetTRLWEBuffer returns a TRLWE buffer (A, B components)
func (bm *BufferManager) GetTRLWEBuffer() ([]params.Torus, []params.Torus) {
	if bm.Rotation.TRLWEInUse >= len(bm.Rotation.TRLWEPool) {
		bm.Rotation.TRLWEInUse = 0
	}
	buffer := bm.Rotation.TRLWEPool[bm.Rotation.TRLWEInUse]
	bm.Rotation.TRLWEInUse++
	return buffer.A, buffer.B
}

// Reset resets all buffer indices
func (bm *BufferManager) Reset() {
	bm.Rotation.InUse = 0
	bm.Rotation.TRLWEInUse = 0
}

// MemoryUsage returns approximate memory usage in bytes
func (bm *BufferManager) MemoryUsage() int {
	n := bm.n
	l := params.GetTRGSWLv1().L

	// Poly: N * 4 bytes, FourierPoly: N * 8 * 2 bytes (complex)
	polySize := n * 4
	fourierSize := n * 8 * 2

	mem := 0

	// FFT buffers
	mem += polySize + fourierSize

	// Decomposition buffers (2*L levels)
	mem += (polySize + fourierSize) * l * 2

	// Multiplication buffers
	mem += fourierSize * 3

	// Rotation pool
	mem += polySize * len(bm.Rotation.Pool)
	mem += polySize * 2 * len(bm.Rotation.TRLWEPool) // A and B

	// Temp buffers
	mem += polySize * 3

	return mem
}
