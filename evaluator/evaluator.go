// Package evaluator provides a zero-allocation TFHE evaluator
// Following tfhe-go's architecture exactly for maximum performance
package evaluator

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/poly"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trgsw"
	"github.com/thedonutfactory/go-tfhe/trlwe"
)

// Evaluator performs TFHE operations with zero allocations
// This follows tfhe-go's architecture exactly
type Evaluator struct {
	// PolyEvaluator for polynomial operations
	PolyEvaluator *poly.Evaluator

	// Decomposer for gadget decomposition
	Decomposer *poly.Decomposer

	// Centralized buffer pool for all operations
	Buffers *BufferPool
}

// NewEvaluator creates a new zero-allocation evaluator
func NewEvaluator(n int) *Evaluator {
	l := params.GetTRGSWLv1().L

	return &Evaluator{
		PolyEvaluator: poly.NewEvaluator(n),
		Decomposer:    poly.NewDecomposer(n, l*2), // 2*L levels for A and B
		Buffers:       NewBufferPool(n),
	}
}

// newEvaluationBuffer creates pre-allocated buffers

// ShallowCopy creates a copy with new buffers (safe for concurrent use)
func (e *Evaluator) ShallowCopy() *Evaluator {
	return &Evaluator{
		PolyEvaluator: e.PolyEvaluator.ShallowCopy(),
		Decomposer:    poly.NewDecomposer(e.PolyEvaluator.Degree(), len(e.Decomposer.GetPolyDecomposedBuffer(1))),
		Buffers:       NewBufferPool(e.PolyEvaluator.Degree()),
	}
}

// ExternalProductAssign computes external product and writes to ctOut
// This is the zero-allocation version following tfhe-go exactly
func (e *Evaluator) ExternalProductAssign(ctFourierGGSW *trgsw.TRGSWLv1FFT, ctIn *trlwe.TRLWELv1, decompositionOffset params.Torus, ctOut *trlwe.TRLWELv1) {
	l := params.GetTRGSWLv1().L
	bgbit := params.GetTRGSWLv1().BGBIT

	// Decompose ctIn into pre-allocated buffers
	polyDecomposed := e.Decomposer.GetPolyDecomposedBuffer(l * 2)
	polyFourierDecomposed := e.Decomposer.GetPolyFourierDecomposedBuffer(l * 2)

	// Decompose A
	poly.DecomposePolyAssign(ctIn.A, int(bgbit), l, decompositionOffset, polyDecomposed[:l])
	// Decompose B
	poly.DecomposePolyAssign(ctIn.B, int(bgbit), l, decompositionOffset, polyDecomposed[l:l*2])

	// Transform to Fourier domain
	for i := 0; i < l*2; i++ {
		e.PolyEvaluator.ToFourierPolyAssign(polyDecomposed[i], polyFourierDecomposed[i])
	}

	// Clear accumulation buffers
	e.Buffers.ExternalProduct.FourierA.Clear()
	e.Buffers.ExternalProduct.FourierB.Clear()

	// Accumulate external product in Fourier domain
	for i := 0; i < l*2; i++ {
		e.PolyEvaluator.MulAddFourierPolyAssign(polyFourierDecomposed[i], ctFourierGGSW.TRLWEFFT[i].A, e.Buffers.ExternalProduct.FourierA)
		e.PolyEvaluator.MulAddFourierPolyAssign(polyFourierDecomposed[i], ctFourierGGSW.TRLWEFFT[i].B, e.Buffers.ExternalProduct.FourierB)
	}

	// Transform back to time domain (write directly to output)
	e.PolyEvaluator.ToPolyAssignUnsafe(e.Buffers.ExternalProduct.FourierA, poly.Poly{Coeffs: ctOut.A})
	e.PolyEvaluator.ToPolyAssignUnsafe(e.Buffers.ExternalProduct.FourierB, poly.Poly{Coeffs: ctOut.B})
}

// CMuxAssign computes ctOut = ct0 + ctCond * (ct1 - ct0)
// Following tfhe-go's pattern exactly
func (e *Evaluator) CMuxAssign(ctCond *trgsw.TRGSWLv1FFT, ct0, ct1 *trlwe.TRLWELv1, decompositionOffset params.Torus, ctOut *trlwe.TRLWELv1) {
	n := params.GetTRGSWLv1().N

	// First copy ct0 to output
	copy(ctOut.A, ct0.A)
	copy(ctOut.B, ct0.B)

	// Compute ct1 - ct0 into buffer.ctCMux
	for i := 0; i < n; i++ {
		e.Buffers.CMUX.Temp.A[i] = ct1.A[i] - ct0.A[i]
		e.Buffers.CMUX.Temp.B[i] = ct1.B[i] - ct0.B[i]
	}

	// External product into pre-allocated buffer
	e.ExternalProductAssign(ctCond, e.Buffers.CMUX.Temp, decompositionOffset, e.Buffers.ExternalProduct.Result)

	// Add to output: ctOut = ct0 + ctCond * (ct1 - ct0)
	for i := 0; i < n; i++ {
		ctOut.A[i] += e.Buffers.ExternalProduct.Result.A[i]
		ctOut.B[i] += e.Buffers.ExternalProduct.Result.B[i]
	}
}

// BlindRotateAssign performs blind rotation and writes to ctOut
// Zero-allocation version following tfhe-go
func (e *Evaluator) BlindRotateAssign(ctIn *tlwe.TLWELv0, testvec *trlwe.TRLWELv1, bsk []*trgsw.TRGSWLv1FFT, decompositionOffset params.Torus, ctOut *trlwe.TRLWELv1) {
	n := params.GetTRGSWLv1().N
	nBit := params.GetTRGSWLv1().NBIT
	tlweLv0N := params.GetTLWELv0().N

	// Initial rotation into buffer.ctAcc1
	bTilda := 2*n - ((int(ctIn.B()) + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))
	poly.PolyMulWithXKInPlace(testvec.A, bTilda, e.Buffers.BlindRotation.Accumulator1.A)
	poly.PolyMulWithXKInPlace(testvec.B, bTilda, e.Buffers.BlindRotation.Accumulator1.B)

	// Iterate through LWE coefficients
	for i := 0; i < tlweLv0N; i++ {
		aTilda := int((ctIn.P[i] + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))

		// Rotate into buffer.ctAcc2
		poly.PolyMulWithXKInPlace(e.Buffers.BlindRotation.Accumulator1.A, aTilda, e.Buffers.BlindRotation.Accumulator2.A)
		poly.PolyMulWithXKInPlace(e.Buffers.BlindRotation.Accumulator1.B, aTilda, e.Buffers.BlindRotation.Accumulator2.B)

		// CMux: ctAcc1 = ctAcc1 + bsk[i] * (ctAcc2 - ctAcc1)
		e.CMuxAssign(bsk[i], e.Buffers.BlindRotation.Accumulator1, e.Buffers.BlindRotation.Accumulator2, decompositionOffset, e.Buffers.BlindRotation.Accumulator1)
	}

	// Copy result to output
	copy(ctOut.A, e.Buffers.BlindRotation.Accumulator1.A)
	copy(ctOut.B, e.Buffers.BlindRotation.Accumulator1.B)
}

// BootstrapAssign performs full bootstrapping (blind rotate + key switch)
// Zero-allocation version - writes to ctOut
func (e *Evaluator) BootstrapAssign(ctIn *tlwe.TLWELv0, testvec *trlwe.TRLWELv1, bsk []*trgsw.TRGSWLv1FFT, ksk []*tlwe.TLWELv0, decompositionOffset params.Torus, ctOut *tlwe.TLWELv0) {
	// Blind rotate
	e.BlindRotateAssign(ctIn, testvec, bsk, decompositionOffset, e.Buffers.BlindRotation.Rotated)

	// Sample extract
	trlwe.SampleExtractIndexAssign(e.Buffers.BlindRotation.Rotated, 0, e.Buffers.Bootstrap.ExtractedLWE)

	// Key switch - writes directly to ctOut (zero-allocation!)
	trgsw.IdentityKeySwitchingAssign(e.Buffers.Bootstrap.ExtractedLWE, ksk, ctOut)
}

// Bootstrap performs full bootstrapping and returns result using buffer pool
// Returns pointer to buffer pool - valid until 4 more bootstrap calls
func (e *Evaluator) Bootstrap(ctIn *tlwe.TLWELv0, testvec *trlwe.TRLWELv1, bsk []*trgsw.TRGSWLv1FFT, ksk []*tlwe.TLWELv0, decompositionOffset params.Torus) *tlwe.TLWELv0 {
	// Get result buffer from pool (round-robin)
	result := e.Buffers.GetNextResult()
	e.BootstrapAssign(ctIn, testvec, bsk, ksk, decompositionOffset, result)
	return result
}

// ResetBuffers resets all buffer pool indices
func (e *Evaluator) ResetBuffers() {
	e.Buffers.Reset()
}
