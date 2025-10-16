package trgsw

import (
	"math"
	"sync"

	"github.com/thedonutfactory/go-tfhe/fft"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/poly"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// TRGSWLv1 represents a Level 1 TRGSW ciphertext
type TRGSWLv1 struct {
	TRLWE []*trlwe.TRLWELv1
}

// NewTRGSWLv1 creates a new TRGSW Level 1 ciphertext
func NewTRGSWLv1() *TRGSWLv1 {
	l := params.GetTRGSWLv1().L
	trlweArray := make([]*trlwe.TRLWELv1, l*2)
	for i := range trlweArray {
		trlweArray[i] = trlwe.NewTRLWELv1()
	}
	return &TRGSWLv1{
		TRLWE: trlweArray,
	}
}

// EncryptTorus encrypts a torus value with TRGSW Level 1
func (t *TRGSWLv1) EncryptTorus(p params.Torus, alpha float64, key []params.Torus, plan *fft.FFTPlan) *TRGSWLv1 {
	l := params.GetTRGSWLv1().L
	bg := float64(params.GetTRGSWLv1().BG)
	n := params.GetTRGSWLv1().N

	// Calculate p_f64 values
	pF64 := make([]float64, l)
	for i := 0; i < l; i++ {
		pF64[i] = 1.0 / math.Pow(bg, float64(i+1))
	}
	pTorus := utils.F64ToTorusVec(pF64)
	plainZero := make([]float64, n)

	// Encrypt all TRLWE samples
	for i := range t.TRLWE {
		t.TRLWE[i] = trlwe.NewTRLWELv1().EncryptF64(plainZero, alpha, key, plan)
	}

	// Add the gadget decomposition
	for i := 0; i < l; i++ {
		t.TRLWE[i].A[0] += p * pTorus[i]
		t.TRLWE[i+l].B[0] += p * pTorus[i]
	}

	return t
}

// TRGSWLv1FFT represents a TRGSW Level 1 ciphertext in FFT form
type TRGSWLv1FFT struct {
	TRLWEFFT []TRLWELv1FFT
}

// TRLWELv1FFT represents a TRLWE Level 1 ciphertext in FFT form
type TRLWELv1FFT struct {
	A poly.FourierPoly
	B poly.FourierPoly
}

// NewTRGSWLv1FFT creates a new TRGSW Level 1 FFT ciphertext from a regular TRGSW
func NewTRGSWLv1FFT(trgsw *TRGSWLv1, polyEval *poly.Evaluator) *TRGSWLv1FFT {
	trlweFFTArray := make([]TRLWELv1FFT, len(trgsw.TRLWE))
	for i, t := range trgsw.TRLWE {
		trlweFFTArray[i] = TRLWELv1FFT{
			A: polyEval.ToFourierPoly(poly.Poly{Coeffs: t.A}),
			B: polyEval.ToFourierPoly(poly.Poly{Coeffs: t.B}),
		}
	}
	return &TRGSWLv1FFT{
		TRLWEFFT: trlweFFTArray,
	}
}

// NewTRGSWLv1FFTDummy creates a dummy TRGSW Level 1 FFT ciphertext
func NewTRGSWLv1FFTDummy(polyEval *poly.Evaluator) *TRGSWLv1FFT {
	l := params.GetTRGSWLv1().L
	trlweFFTArray := make([]TRLWELv1FFT, l*2)
	for i := range trlweFFTArray {
		trlweFFTArray[i] = TRLWELv1FFT{
			A: polyEval.NewFourierPoly(),
			B: polyEval.NewFourierPoly(),
		}
	}
	return &TRGSWLv1FFT{
		TRLWEFFT: trlweFFTArray,
	}
}

// CloudKeyData contains the data needed from CloudKey for trgsw operations
type CloudKeyData interface {
	GetDecompositionOffset() params.Torus
	GetBlindRotateTestvec() *trlwe.TRLWELv1
	GetBootstrappingKey() []*TRGSWLv1FFT
}

// ExternalProductWithFFT performs external product with FFT optimization
// This version uses pre-allocated buffers for maximum zero-allocation performance
func ExternalProductWithFFT(trgswFFT *TRGSWLv1FFT, trlweIn *trlwe.TRLWELv1, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	l := params.GetTRGSWLv1().L

	// Use decomposition buffer pool (zero-allocation)
	decompositionInPlace(trlweIn, decompositionOffset, polyEval)

	// Clear accumulation buffers (reuse existing buffers)
	polyEval.ClearBuffer("fpAcc")
	polyEval.ClearBuffer("fpBcc")

	// For each decomposition level
	for i := 0; i < l*2; i++ {
		// decFFT is already in buffer.decompFFT[i] from decompositionInPlace
		// Accumulate in frequency domain (multiply-add)
		polyEval.MulAddFourierPolyAssignBuffered(i, trgswFFT.TRLWEFFT[i].A, "fpAcc")
		polyEval.MulAddFourierPolyAssignBuffered(i, trgswFFT.TRLWEFFT[i].B, "fpBcc")
	}

	// Get pooled TRLWE buffer for result
	resultA, resultB := polyEval.GetTRLWEBuffer()

	// Transform back to time domain using pooled buffer
	polyEval.BufferToPolyAssign("fpAcc", resultA)
	polyEval.BufferToPolyAssign("fpBcc", resultB)

	return &trlwe.TRLWELv1{A: resultA, B: resultB}
}

// decompositionInPlace performs gadget decomposition directly into evaluator buffers (zero-allocation)
func decompositionInPlace(trlweIn *trlwe.TRLWELv1, decompositionOffset params.Torus, polyEval *poly.Evaluator) {
	l := params.GetTRGSWLv1().L
	n := params.GetTRGSWLv1().N

	offset := decompositionOffset
	bgbit := params.GetTRGSWLv1().BGBIT
	mask := params.Torus((1 << bgbit) - 1)
	halfBG := params.Torus(1 << (bgbit - 1))

	// Decompose directly into buffers
	for i := 0; i < l*2; i++ {
		buf := polyEval.GetDecompBuffer(i)
		for j := 0; j < n; j++ {
			buf.Coeffs[j] = 0
		}
	}

	for j := 0; j < n; j++ {
		tmp0 := trlweIn.A[j] + offset
		tmp1 := trlweIn.B[j] + offset
		for i := 0; i < l; i++ {
			polyEval.GetDecompBuffer(i).Coeffs[j] = ((tmp0 >> (32 - (uint32(i)+1)*bgbit)) & mask) - halfBG
		}
		for i := 0; i < l; i++ {
			polyEval.GetDecompBuffer(i + l).Coeffs[j] = ((tmp1 >> (32 - (uint32(i)+1)*bgbit)) & mask) - halfBG
		}
	}

	// Transform all decomposition levels to frequency domain
	for i := 0; i < l*2; i++ {
		polyEval.ToFourierPolyInBuffer(*polyEval.GetDecompBuffer(i), i)
	}
}

// CMUX performs controlled MUX operation (zero-allocation version using TRLWE pool)
// if cond == 0 then in1 else in2
func CMUX(in1, in2 *trlwe.TRLWELv1, cond *TRGSWLv1FFT, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	n := params.GetTRGSWLv1().N

	// Get TRLWE buffer from pool for difference computation
	tmpA, tmpB := polyEval.GetTRLWEBuffer()
	for i := 0; i < n; i++ {
		tmpA[i] = in2.A[i] - in1.A[i]
		tmpB[i] = in2.B[i] - in1.B[i]
	}
	tmp := &trlwe.TRLWELv1{A: tmpA, B: tmpB}

	// External product (uses internal buffers for zero-alloc in hot path)
	tmp2 := ExternalProductWithFFT(cond, tmp, decompositionOffset, polyEval)

	// Add in1 to result (reuse tmp2)
	for i := 0; i < n; i++ {
		tmp2.A[i] += in1.A[i]
		tmp2.B[i] += in1.B[i]
	}

	return tmp2
}

// BlindRotate performs blind rotation for bootstrapping (optimized with buffer pool)
func BlindRotate(src *tlwe.TLWELv0, blindRotateTestvec *trlwe.TRLWELv1, bootstrappingKey []*TRGSWLv1FFT, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	n := params.GetTRGSWLv1().N
	nBit := params.GetTRGSWLv1().NBIT

	// Reset rotation pool for this operation
	polyEval.ResetRotationPool()

	bTilda := 2*n - ((int(src.B()) + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))

	// Initial rotation using buffer pool
	resultA := polyEval.PolyMulWithXK(blindRotateTestvec.A, bTilda)
	resultB := polyEval.PolyMulWithXK(blindRotateTestvec.B, bTilda)
	result := &trlwe.TRLWELv1{A: resultA, B: resultB}

	tlweLv0N := params.GetTLWELv0().N
	for i := 0; i < tlweLv0N; i++ {
		aTilda := int((src.P[i] + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))

		// Use buffer pool for rotation
		res2A := polyEval.PolyMulWithXK(result.A, aTilda)
		res2B := polyEval.PolyMulWithXK(result.B, aTilda)
		res2 := &trlwe.TRLWELv1{A: res2A, B: res2B}

		result = CMUX(result, res2, bootstrappingKey[i], decompositionOffset, polyEval)
	}

	return result
}

// evaluatorPool is a pool of evaluators for parallel operations
var evaluatorPool = sync.Pool{
	New: func() interface{} {
		return poly.NewEvaluator(params.GetTRGSWLv1().N)
	},
}

// BatchBlindRotate performs multiple blind rotations in parallel (zero-allocation)
func BatchBlindRotate(srcs []*tlwe.TLWELv0, blindRotateTestvec *trlwe.TRLWELv1, bootstrappingKey []*TRGSWLv1FFT, decompositionOffset params.Torus) []*trlwe.TRLWELv1 {
	results := make([]*trlwe.TRLWELv1, len(srcs))
	var wg sync.WaitGroup

	for i, src := range srcs {
		wg.Add(1)
		go func(idx int, s *tlwe.TLWELv0) {
			defer wg.Done()
			// Get evaluator from pool (reuse instead of allocate)
			polyEval := evaluatorPool.Get().(*poly.Evaluator)
			defer evaluatorPool.Put(polyEval)

			results[idx] = BlindRotate(s, blindRotateTestvec, bootstrappingKey, decompositionOffset, polyEval)
		}(i, src)
	}

	wg.Wait()
	return results
}

// polyMulWithXKInPlace multiplies a polynomial by X^k in-place (zero-allocation)
func polyMulWithXKInPlace(a []params.Torus, k int, result []params.Torus) {
	n := len(a)
	k = k % (2 * n) // Normalize k to [0, 2N)

	if k == 0 {
		copy(result, a)
		return
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

// IdentityKeySwitching performs identity key switching
func IdentityKeySwitching(src *tlwe.TLWELv1, keySwitchingKey []*tlwe.TLWELv0) *tlwe.TLWELv0 {
	n := params.GetTRGSWLv1().N
	basebit := params.GetTRGSWLv1().BASEBIT
	base := 1 << basebit
	iksT := params.GetTRGSWLv1().IKS_T

	result := tlwe.NewTLWELv0()
	tlweLv0N := params.GetTLWELv0().N
	result.P[tlweLv0N] = src.P[len(src.P)-1]

	precOffset := params.Torus(1 << (32 - (1 + basebit*iksT)))

	for i := 0; i < n; i++ {
		aBar := src.P[i] + precOffset
		for j := 0; j < iksT; j++ {
			k := (aBar >> (32 - (j+1)*basebit)) & params.Torus((1<<basebit)-1)
			if k != 0 {
				idx := (base * iksT * i) + (base * j) + int(k)
				for x := 0; x < len(result.P); x++ {
					result.P[x] -= keySwitchingKey[idx].P[x]
				}
			}
		}
	}

	return result
}
