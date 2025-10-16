package trgsw

import (
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
		pF64[i] = 1.0 / pow(bg, float64(i+1))
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
func ExternalProductWithFFT(trgswFFT *TRGSWLv1FFT, trlweIn *trlwe.TRLWELv1, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	dec := decomposition(trlweIn, decompositionOffset)

	l := params.GetTRGSWLv1().L

	// Initialize output in frequency domain
	outAFFT := polyEval.NewFourierPoly()
	outBFFT := polyEval.NewFourierPoly()

	// For each decomposition level
	for i := 0; i < l*2; i++ {
		// Convert decomposition to Poly
		decPoly := poly.Poly{Coeffs: dec[i][:]}

		// Transform to frequency domain
		decFFT := polyEval.ToFourierPoly(decPoly)

		// Accumulate in frequency domain (multiply-add)
		polyEval.MulAddFourierPolyAssign(decFFT, trgswFFT.TRLWEFFT[i].A, outAFFT)
		polyEval.MulAddFourierPolyAssign(decFFT, trgswFFT.TRLWEFFT[i].B, outBFFT)
	}

	// Transform back to time domain
	result := trlwe.NewTRLWELv1()
	outA := poly.Poly{Coeffs: result.A}
	outB := poly.Poly{Coeffs: result.B}
	polyEval.ToPolyAssignUnsafe(outAFFT, outA)
	polyEval.ToPolyAssignUnsafe(outBFFT, outB)

	return result
}

// decomposition performs gadget decomposition of a TRLWE ciphertext
func decomposition(trlweIn *trlwe.TRLWELv1, decompositionOffset params.Torus) [][1024]params.Torus {
	l := params.GetTRGSWLv1().L
	n := params.GetTRGSWLv1().N
	result := make([][1024]params.Torus, l*2)

	offset := decompositionOffset
	bgbit := params.GetTRGSWLv1().BGBIT
	mask := params.Torus((1 << bgbit) - 1)
	halfBG := params.Torus(1 << (bgbit - 1))

	for j := 0; j < n; j++ {
		tmp0 := trlweIn.A[j] + offset
		tmp1 := trlweIn.B[j] + offset
		for i := 0; i < l; i++ {
			result[i][j] = ((tmp0 >> (32 - (uint32(i)+1)*bgbit)) & mask) - halfBG
		}
		for i := 0; i < l; i++ {
			result[i+l][j] = ((tmp1 >> (32 - (uint32(i)+1)*bgbit)) & mask) - halfBG
		}
	}

	return result
}

// CMUX performs controlled MUX operation
// if cond == 0 then in1 else in2
func CMUX(in1, in2 *trlwe.TRLWELv1, cond *TRGSWLv1FFT, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	n := params.GetTRGSWLv1().N
	tmp := trlwe.NewTRLWELv1()

	for i := 0; i < n; i++ {
		tmp.A[i] = in2.A[i] - in1.A[i]
		tmp.B[i] = in2.B[i] - in1.B[i]
	}

	tmp2 := ExternalProductWithFFT(cond, tmp, decompositionOffset, polyEval)
	result := trlwe.NewTRLWELv1()

	for i := 0; i < n; i++ {
		result.A[i] = tmp2.A[i] + in1.A[i]
		result.B[i] = tmp2.B[i] + in1.B[i]
	}

	return result
}

// BlindRotate performs blind rotation for bootstrapping
func BlindRotate(src *tlwe.TLWELv0, blindRotateTestvec *trlwe.TRLWELv1, bootstrappingKey []*TRGSWLv1FFT, decompositionOffset params.Torus, polyEval *poly.Evaluator) *trlwe.TRLWELv1 {
	n := params.GetTRGSWLv1().N
	nBit := params.GetTRGSWLv1().NBIT

	bTilda := 2*n - ((int(src.B()) + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))
	result := &trlwe.TRLWELv1{
		A: polyMulWithXK(blindRotateTestvec.A, bTilda),
		B: polyMulWithXK(blindRotateTestvec.B, bTilda),
	}

	tlweLv0N := params.GetTLWELv0().N
	for i := 0; i < tlweLv0N; i++ {
		aTilda := int((src.P[i] + (1 << (31 - nBit - 1))) >> (32 - nBit - 1))
		res2 := &trlwe.TRLWELv1{
			A: polyMulWithXK(result.A, aTilda),
			B: polyMulWithXK(result.B, aTilda),
		}
		result = CMUX(result, res2, bootstrappingKey[i], decompositionOffset, polyEval)
	}

	return result
}

// BatchBlindRotate performs multiple blind rotations in parallel
func BatchBlindRotate(srcs []*tlwe.TLWELv0, blindRotateTestvec *trlwe.TRLWELv1, bootstrappingKey []*TRGSWLv1FFT, decompositionOffset params.Torus) []*trlwe.TRLWELv1 {
	results := make([]*trlwe.TRLWELv1, len(srcs))
	var wg sync.WaitGroup

	for i, src := range srcs {
		wg.Add(1)
		go func(idx int, s *tlwe.TLWELv0) {
			defer wg.Done()
			polyEval := poly.NewEvaluator(params.GetTRGSWLv1().N)
			results[idx] = BlindRotate(s, blindRotateTestvec, bootstrappingKey, decompositionOffset, polyEval)
		}(i, src)
	}

	wg.Wait()
	return results
}

// polyMulWithXK multiplies a polynomial by X^k in the ring Z[X]/(X^N+1)
func polyMulWithXK(a []params.Torus, k int) []params.Torus {
	n := params.GetTRGSWLv1().N
	result := make([]params.Torus, n)

	if k < n {
		for i := 0; i < n-k; i++ {
			result[i+k] = a[i]
		}
		for i := n - k; i < n; i++ {
			result[i+k-n] = ^params.Torus(0) - a[i]
		}
	} else {
		for i := 0; i < 2*n-k; i++ {
			result[i+k-n] = ^params.Torus(0) - a[i]
		}
		for i := 2*n - k; i < n; i++ {
			result[i-(2*n-k)] = a[i]
		}
	}

	return result
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

// Helper function for power
func pow(base, exp float64) float64 {
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}
