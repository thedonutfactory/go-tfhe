package poly

import (
	"math"
	"math/cmplx"
)

// Evaluator computes polynomial operations over the N-th cyclotomic ring.
// This is optimized for TFHE operations with precomputed twiddle factors.
type Evaluator struct {
	// degree is the degree of polynomial that this evaluator can handle.
	degree int
	// q is a float64 value of the modulus (2^32 for Torus).
	q float64

	// tw is the twiddle factors for fourier transform.
	tw []complex128
	// twInv is the twiddle factors for inverse fourier transform.
	twInv []complex128
	// twMono is the twiddle factors for monomial fourier transform.
	twMono []complex128
	// twMonoIdx is the precomputed bit-reversed index for monomial fourier transform.
	twMonoIdx []int

	buffer evaluationBuffer
}

// evaluationBuffer is a buffer for Evaluator.
type evaluationBuffer struct {
	// fp is an intermediate FFT buffer.
	fp FourierPoly
	// fpInv is an intermediate inverse FFT buffer.
	fpInv FourierPoly
	// pSplit is a buffer for split operations.
	pSplit Poly
}

// NewEvaluator creates a new Evaluator with degree N.
func NewEvaluator(N int) *Evaluator {
	if !isPowerOfTwo(N) {
		panic("degree not power of two")
	}
	if N < MinDegree {
		panic("degree smaller than MinDegree")
	}

	// Q = 2^32 for Torus (uint32)
	Q := math.Exp2(32)

	tw, twInv := genTwiddleFactors(N / 2)

	twMono := make([]complex128, 2*N)
	for i := 0; i < 2*N; i++ {
		e := -math.Pi * float64(i) / float64(N)
		twMono[i] = cmplx.Exp(complex(0, e))
	}

	twMonoIdx := make([]int, N/2)
	twMonoIdx[0] = 2*N - 1
	for i := 1; i < N/2; i++ {
		twMonoIdx[i] = 4*i - 1
	}
	bitReverseInPlace(twMonoIdx)

	return &Evaluator{
		degree:    N,
		q:         Q,
		tw:        tw,
		twInv:     twInv,
		twMono:    twMono,
		twMonoIdx: twMonoIdx,
		buffer:    newEvaluationBuffer(N),
	}
}

// genTwiddleFactors generates twiddle factors for FFT.
func genTwiddleFactors(N int) (tw, twInv []complex128) {
	twFFT := make([]complex128, N/2)
	twInvFFT := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		e := -2 * math.Pi * float64(i) / float64(N)
		twFFT[i] = cmplx.Exp(complex(0, e))
		twInvFFT[i] = cmplx.Exp(-complex(0, e))
	}
	bitReverseInPlace(twFFT)
	bitReverseInPlace(twInvFFT)

	tw = make([]complex128, 0, N-1)
	twInv = make([]complex128, 0, N-1)

	for m, t := 1, N/2; m <= N/2; m, t = m<<1, t>>1 {
		twFold := cmplx.Exp(complex(0, 2*math.Pi*float64(t)/float64(4*N)))
		for i := 0; i < m; i++ {
			tw = append(tw, twFFT[i]*twFold)
		}
	}

	for m, t := N/2, 1; m >= 1; m, t = m>>1, t<<1 {
		twInvFold := cmplx.Exp(complex(0, -2*math.Pi*float64(t)/float64(4*N)))
		for i := 0; i < m; i++ {
			twInv = append(twInv, twInvFFT[i]*twInvFold)
		}
	}

	return tw, twInv
}

// bitReverseInPlace performs bit reversal permutation in place.
func bitReverseInPlace[T any](data []T) {
	n := len(data)
	if n <= 1 {
		return
	}

	j := 0
	for i := 0; i < n; i++ {
		if i < j {
			data[i], data[j] = data[j], data[i]
		}
		// Bit reversal
		m := n >> 1
		for m > 0 && j >= m {
			j -= m
			m >>= 1
		}
		j += m
	}
}

// newEvaluationBuffer creates a new evaluationBuffer.
func newEvaluationBuffer(N int) evaluationBuffer {
	return evaluationBuffer{
		fp:     NewFourierPoly(N),
		fpInv:  NewFourierPoly(N),
		pSplit: NewPoly(N),
	}
}

// Degree returns the degree of polynomial that the evaluator can handle.
func (e *Evaluator) Degree() int {
	return e.degree
}

// NewPoly creates a new polynomial with the same degree as the evaluator.
func (e *Evaluator) NewPoly() Poly {
	return NewPoly(e.degree)
}

// NewFourierPoly creates a new fourier polynomial with the same degree as the evaluator.
func (e *Evaluator) NewFourierPoly() FourierPoly {
	return NewFourierPoly(e.degree)
}

// ShallowCopy returns a shallow copy of this Evaluator.
// Returned Evaluator is safe for concurrent use.
func (e *Evaluator) ShallowCopy() *Evaluator {
	return &Evaluator{
		degree:    e.degree,
		q:         e.q,
		tw:        e.tw,
		twInv:     e.twInv,
		twMono:    e.twMono,
		twMonoIdx: e.twMonoIdx,
		buffer:    newEvaluationBuffer(e.degree),
	}
}
