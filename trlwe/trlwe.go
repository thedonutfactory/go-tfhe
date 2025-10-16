package trlwe

import (
	"math/rand"

	"github.com/thedonutfactory/go-tfhe/fft"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// TRLWELv1 represents a Level 1 TRLWE ciphertext
type TRLWELv1 struct {
	A []params.Torus
	B []params.Torus
}

// NewTRLWELv1 creates a new TRLWE Level 1 ciphertext
func NewTRLWELv1() *TRLWELv1 {
	n := params.GetTRLWELv1().N
	return &TRLWELv1{
		A: make([]params.Torus, n),
		B: make([]params.Torus, n),
	}
}

// EncryptF64 encrypts a vector of float64 values with TRLWE Level 1
func (t *TRLWELv1) EncryptF64(p []float64, alpha float64, key []params.Torus, plan *fft.FFTPlan) *TRLWELv1 {
	rng := rand.New(rand.NewSource(rand.Int63()))
	n := params.GetTRLWELv1().N

	// Generate random a
	for i := 0; i < n; i++ {
		t.A[i] = params.Torus(rng.Uint32())
	}

	// Add Gaussian noise to plaintext
	t.B = utils.GaussianF64Vec(p, alpha, rng)

	// Compute a * s and add to b
	var aArray [1024]params.Torus
	var keyArray [1024]params.Torus
	copy(aArray[:], t.A)
	copy(keyArray[:], key)

	polyRes := plan.Processor.PolyMul1024(&aArray, &keyArray)

	for i := 0; i < n; i++ {
		t.B[i] += polyRes[i]
	}

	return t
}

// EncryptBool encrypts a vector of boolean values with TRLWE Level 1
func (t *TRLWELv1) EncryptBool(pBool []bool, alpha float64, key []params.Torus, plan *fft.FFTPlan) *TRLWELv1 {
	pF64 := make([]float64, len(pBool))
	for i, b := range pBool {
		if b {
			pF64[i] = 0.125
		} else {
			pF64[i] = -0.125
		}
	}
	return t.EncryptF64(pF64, alpha, key, plan)
}

// DecryptBool decrypts a TRLWE Level 1 ciphertext to a vector of booleans
func (t *TRLWELv1) DecryptBool(key []params.Torus, plan *fft.FFTPlan) []bool {
	n := len(t.A)

	var aArray [1024]params.Torus
	var keyArray [1024]params.Torus
	copy(aArray[:], t.A)
	copy(keyArray[:], key)

	polyRes := plan.Processor.PolyMul1024(&aArray, &keyArray)

	result := make([]bool, n)
	for i := 0; i < n; i++ {
		value := int32(t.B[i] - polyRes[i])
		result[i] = value >= 0
	}
	return result
}

// TRLWELv1FFT represents a TRLWE Level 1 ciphertext in FFT form
type TRLWELv1FFT struct {
	A []float64
	B []float64
}

// NewTRLWELv1FFT creates a new TRLWE Level 1 FFT ciphertext from a regular TRLWE
func NewTRLWELv1FFT(trlwe *TRLWELv1, plan *fft.FFTPlan) *TRLWELv1FFT {
	var aArray [1024]params.Torus
	var bArray [1024]params.Torus
	copy(aArray[:], trlwe.A)
	copy(bArray[:], trlwe.B)

	aFFT := plan.Processor.IFFT1024(&aArray)
	bFFT := plan.Processor.IFFT1024(&bArray)

	return &TRLWELv1FFT{
		A: aFFT[:],
		B: bFFT[:],
	}
}

// NewTRLWELv1FFTDummy creates a dummy TRLWE Level 1 FFT ciphertext
func NewTRLWELv1FFTDummy() *TRLWELv1FFT {
	// FourierPoly needs 2*N for interleaved real/imaginary layout
	return &TRLWELv1FFT{
		A: make([]float64, 2*params.GetTRLWELv1().N),
		B: make([]float64, 2*params.GetTRLWELv1().N),
	}
}

// SampleExtractIndex extracts a TLWE sample from a TRLWE at index k
func SampleExtractIndex(trlwe *TRLWELv1, k int) *tlwe.TLWELv1 {
	n := params.GetTRLWELv1().N
	result := tlwe.NewTLWELv1()

	for i := 0; i < n; i++ {
		if i <= k {
			result.P[i] = trlwe.A[k-i]
		} else {
			result.P[i] = ^params.Torus(0) - trlwe.A[n+k-i]
		}
	}
	result.SetB(trlwe.B[k])

	return result
}

// SampleExtractIndex2 extracts a TLWE Lv0 sample from a TRLWE at index k
func SampleExtractIndex2(trlwe *TRLWELv1, k int) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N
	result := tlwe.NewTLWELv0()

	for i := 0; i < n; i++ {
		if i <= k {
			result.P[i] = trlwe.A[k-i]
		} else {
			result.P[i] = ^params.Torus(0) - trlwe.A[n+k-i]
		}
	}
	result.SetB(trlwe.B[k])

	return result
}
