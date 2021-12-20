package core

import (
	"github.com/thedonutfactory/go-tfhe/fft"
)

type TLweSampleFFT struct {
	A []*fft.LagrangeHalfCPolynomial ///< array of length k+1: mask + right term
	//B               *TorusPolynomial  ///< alias of a[k] to get the right term
	CurrentVariance float64 ///< avg variance of the sample
	K               int32
}

func NewTLweSampleFFT(params *TLweParams) *TLweSampleFFT {
	avar := fft.NewLagrangeHalfCPolynomialArray(params.K+1, params.N)
	return &TLweSampleFFT{
		K: params.K,
		A: avar,
		//B:               &avar[params.K],
		CurrentVariance: 0.,
	}
}

func NewTLweSampleFFTArray(size int32, params *TLweParams) []*TLweSampleFFT {
	arr := make([]*TLweSampleFFT, size)
	for i := int32(0); i < size; i++ {
		arr[i] = NewTLweSampleFFT(params)
	}
	return arr
}

/*
func InitTLweSampleFFT(obj *TLweSampleFFT, params *TLweParams) {
	//a is a table of k+1 polynomials, b is an alias for &a[k]
	k := params.K
	a := NewLagrangeHalfCPolynomialArray(k+1, params.N)
	obj = &TLweSampleFFT{
		K: params.K,
		A: a,
		//B:               &avar[params.K],
		CurrentVariance: 0.,
	}
}
*/

// Computes the inverse FFT of the coefficients of the TLWE sample
func tLweToFFTConvert(result *TLweSampleFFT, source *TLweSample, params *TLweParams) {
	k := params.K
	for i := int32(0); i <= k; i++ {
		fft.TorusPolynomialIfft(result.A[i], source.A[i].Coefs)
	}
	result.CurrentVariance = source.CurrentVariance
}

// Computes the FFT of the coefficients of the TLWEfft sample
func tLweFromFFTConvert(result *TLweSample, source *TLweSampleFFT, params *TLweParams) {
	k := params.K
	for i := int32(0); i <= k; i++ {
		result.A[i].Coefs = fft.TorusPolynomialFft(source.A[i])
	}
	result.CurrentVariance = source.CurrentVariance
}

//Arithmetic operations on TLwe samples
/** result = (0,0) */
func tLweFFTClear(result *TLweSampleFFT, params *TLweParams) {
	k := params.K
	for i := int32(0); i <= k; i++ {
		result.A[i].Clear()
	}
	result.CurrentVariance = 0.
}

// result = result + p*sample
func tLweFFTAddMulRTo(result *TLweSampleFFT, p *fft.LagrangeHalfCPolynomial, sample *TLweSampleFFT, params *TLweParams) {

	k := params.K
	for i := int32(0); i <= k; i++ {
		result.A[i].AddMul(p, sample.A[i])
	}
	//result.current_variance += sample.current_variance;
	//TODO: how to compute the variance correctly?
}
