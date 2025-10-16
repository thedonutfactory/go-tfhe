package poly

import "github.com/thedonutfactory/go-tfhe/params"

// ============================================================================
// UNIFIED BUFFER METHODS
// ============================================================================
// All buffer pool operations consolidated in one place for clarity.
// These methods operate on poly.Evaluator.buffer (evaluationBuffer struct).

// ============================================================================
// FOURIER BUFFER OPERATIONS
// ============================================================================

// ClearBuffer clears a named Fourier buffer (sets all coefficients to zero)
func (e *Evaluator) ClearBuffer(name string) {
	switch name {
	case "fpAcc":
		e.buffer.fpAcc.Clear()
	case "fpBcc":
		e.buffer.fpBcc.Clear()
	case "fpDiff":
		e.buffer.fpDiff.Clear()
	case "fpMul1":
		e.buffer.fpMul1.Clear()
	case "fpMul2":
		e.buffer.fpMul2.Clear()
	default:
		panic("unknown buffer name: " + name)
	}
}

// MulAddFourierPolyAssignBuffered performs fpOut += decompFFT[idx] * fp
// using the pre-allocated decomposition buffer
func (e *Evaluator) MulAddFourierPolyAssignBuffered(idx int, fp FourierPoly, bufferName string) {
	var fpOut *FourierPoly
	switch bufferName {
	case "fpAcc":
		fpOut = &e.buffer.fpAcc
	case "fpBcc":
		fpOut = &e.buffer.fpBcc
	default:
		panic("unknown buffer name: " + bufferName)
	}

	// Use the pre-computed FFT from decomposition buffer
	e.MulAddFourierPolyAssign(e.buffer.decompFFT[idx], fp, *fpOut)
}

// BufferToPolyAssign converts a buffer from frequency domain to time domain
// and writes directly to the output slice (zero-allocation)
func (e *Evaluator) BufferToPolyAssign(bufferName string, out []params.Torus) {
	var fp *FourierPoly
	switch bufferName {
	case "fpAcc":
		fp = &e.buffer.fpAcc
	case "fpBcc":
		fp = &e.buffer.fpBcc
	case "fpDiff":
		fp = &e.buffer.fpDiff
	default:
		panic("unknown buffer name: " + bufferName)
	}

	// Use unsafe conversion to avoid allocation
	pOut := Poly{Coeffs: out}
	e.ToPolyAssignUnsafe(*fp, pOut)
}

// ============================================================================
// DECOMPOSITION BUFFER OPERATIONS
// ============================================================================

// GetDecompBuffer returns the i-th decomposition buffer for direct write
func (e *Evaluator) GetDecompBuffer(i int) *Poly {
	if i >= len(e.buffer.decompBuffer) {
		panic("decomposition buffer index out of range")
	}
	return &e.buffer.decompBuffer[i]
}

// GetDecompFFTBuffer returns the i-th decomposition FFT buffer
func (e *Evaluator) GetDecompFFTBuffer(i int) *FourierPoly {
	if i >= len(e.buffer.decompFFT) {
		panic("decomposition FFT buffer index out of range")
	}
	return &e.buffer.decompFFT[i]
}

// ToFourierPolyInBuffer transforms a poly to fourier and stores in buffer
func (e *Evaluator) ToFourierPolyInBuffer(p Poly, bufferIdx int) {
	if bufferIdx >= len(e.buffer.decompFFT) {
		panic("buffer index out of range")
	}
	e.ToFourierPolyAssign(p, e.buffer.decompFFT[bufferIdx])
}

// CopyToDecompBuffer copies a polynomial into the decomposition buffer
func (e *Evaluator) CopyToDecompBuffer(src []params.Torus, bufferIdx int) {
	if bufferIdx >= len(e.buffer.decompBuffer) {
		panic("buffer index out of range")
	}
	copy(e.buffer.decompBuffer[bufferIdx].Coeffs, src)
}

// ============================================================================
// ROTATION POOL OPERATIONS
// ============================================================================

// GetRotationBuffer returns a rotation buffer from the pool
// Uses round-robin allocation to avoid conflicts
func (e *Evaluator) GetRotationBuffer() []params.Torus {
	buf := e.buffer.rotationPool[e.buffer.rotationIdx].Coeffs
	e.buffer.rotationIdx = (e.buffer.rotationIdx + 1) % len(e.buffer.rotationPool)
	return buf
}

// ResetRotationPool resets the rotation buffer pool index
// Call this at the start of a new operation to ensure clean state
func (e *Evaluator) ResetRotationPool() {
	e.buffer.rotationIdx = 0
}

// PolyMulWithXK multiplies a polynomial by X^k using a pooled buffer (zero-allocation)
func (e *Evaluator) PolyMulWithXK(a []params.Torus, k int) []params.Torus {
	result := e.GetRotationBuffer()
	PolyMulWithXKInPlace(a, k, result)
	return result
}

// PolyMulWithXKInPlace multiplies polynomial by X^k in the ring Z[X]/(X^N+1)
// This is the core rotation operation used throughout TFHE
func PolyMulWithXKInPlace(a []params.Torus, k int, result []params.Torus) {
	n := len(a)
	k = k % (2 * n) // Normalize k to [0, 2N)

	if k == 0 {
		copy(result, a)
		return
	}

	if k < 0 {
		k += 2 * n
	}

	if k < n {
		// Positive rotation: coefficients shift right, wrap with negation
		for i := 0; i < n-k; i++ {
			result[i+k] = a[i]
		}
		for i := n - k; i < n; i++ {
			result[i+k-n] = ^params.Torus(0) - a[i]
		}
	} else {
		// Rotation >= n: all coefficients get negated
		k -= n
		for i := 0; i < n-k; i++ {
			result[i+k] = ^params.Torus(0) - a[i]
		}
		for i := n - k; i < n; i++ {
			result[i+k-n] = a[i]
		}
	}
}

// PolyMulWithXKDirect multiplies by X^k and writes to provided buffer (zero-allocation)
func (e *Evaluator) PolyMulWithXKDirect(a []params.Torus, k int, result []params.Torus) {
	PolyMulWithXKInPlace(a, k, result)
}

// ============================================================================
// TRLWE POOL OPERATIONS
// ============================================================================

// GetTRLWEBuffer returns a TRLWE buffer from the pool
// Returns (A, B) slices that can be used to construct a TRLWE
func (e *Evaluator) GetTRLWEBuffer() ([]params.Torus, []params.Torus) {
	buf := &e.buffer.trlwePool[e.buffer.trlweIdx]
	e.buffer.trlweIdx = (e.buffer.trlweIdx + 1) % len(e.buffer.trlwePool)
	return buf.A, buf.B
}

// ResetTRLWEPool resets the TRLWE pool index
func (e *Evaluator) ResetTRLWEPool() {
	e.buffer.trlweIdx = 0
}

// ClearTRLWEBuffer clears a TRLWE buffer
func (e *Evaluator) ClearTRLWEBuffer(a, b []params.Torus) {
	for i := range a {
		a[i] = 0
		b[i] = 0
	}
}
