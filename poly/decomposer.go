package poly

import "github.com/thedonutfactory/go-tfhe/params"

// Decomposer performs gadget decomposition with pre-allocated buffers
// This achieves zero-allocation decomposition operations
type Decomposer struct {
	buffer decompositionBuffer
}

// decompositionBuffer contains pre-allocated buffers for decomposition
type decompositionBuffer struct {
	// polyDecomposed is the pre-allocated buffer for polynomial decomposition
	polyDecomposed []Poly
	// polyFourierDecomposed is the pre-allocated buffer for Fourier-domain decomposition
	polyFourierDecomposed []FourierPoly
}

// NewDecomposer creates a new Decomposer with buffers for up to maxLevel decomposition levels
func NewDecomposer(N int, maxLevel int) *Decomposer {
	polyDecomposed := make([]Poly, maxLevel)
	polyFourierDecomposed := make([]FourierPoly, maxLevel)

	for i := 0; i < maxLevel; i++ {
		polyDecomposed[i] = NewPoly(N)
		polyFourierDecomposed[i] = NewFourierPoly(N)
	}

	return &Decomposer{
		buffer: decompositionBuffer{
			polyDecomposed:        polyDecomposed,
			polyFourierDecomposed: polyFourierDecomposed,
		},
	}
}

// GetPolyDecomposedBuffer returns the decomposition buffer for polynomial
func (d *Decomposer) GetPolyDecomposedBuffer(level int) []Poly {
	if level > len(d.buffer.polyDecomposed) {
		panic("decomposition level exceeds buffer size")
	}
	return d.buffer.polyDecomposed[:level]
}

// GetPolyFourierDecomposedBuffer returns the Fourier decomposition buffer
func (d *Decomposer) GetPolyFourierDecomposedBuffer(level int) []FourierPoly {
	if level > len(d.buffer.polyFourierDecomposed) {
		panic("decomposition level exceeds buffer size")
	}
	return d.buffer.polyFourierDecomposed[:level]
}

// DecomposePolyAssign decomposes polynomial p into decomposedOut using gadget decomposition
// This writes directly to the provided buffer (zero-allocation)
func DecomposePolyAssign(p []params.Torus, bgbit, level int, offset params.Torus, decomposedOut []Poly) {
	n := len(p)
	mask := params.Torus((1 << bgbit) - 1)
	halfBG := params.Torus(1 << (bgbit - 1))

	for j := 0; j < n; j++ {
		tmp := p[j] + offset
		for i := 0; i < level; i++ {
			decomposedOut[i].Coeffs[j] = ((tmp >> (32 - (uint32(i)+1)*uint32(bgbit))) & mask) - halfBG
		}
	}
}
