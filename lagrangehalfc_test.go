package tfhe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TEST(LagrangeHalfcTest, fftIsBijective) {
func TestFftIsBijective(t *testing.T) {
	assert := assert.New(t)

	NBTRIALS := 10
	toler := 1e-9
	var N int32 = 1024
	for trials := 0; trials < NBTRIALS; trials++ {
		a := NewTorusPolynomial(N)
		acopy := NewTorusPolynomial(N)
		b := NewTorusPolynomial(N)
		afft := NewLagrangeHalfCPolynomial(N)
		torusPolynomialUniform(a)
		torusPolynomialCopy(acopy, a)
		TorusPolynomial_ifft(afft, a)
		TorusPolynomial_fft(b, afft)
		assert.EqualValues(torusPolynomialNormInftyDist(a, acopy), 0)
		assert.LessOrEqual(torusPolynomialNormInftyDist(a, b), toler)
	}
}

//TEST(LagrangeHalfcTest, LagrangeHalfCPolynomialClear) {
func TestLagrangeHalfCPolynomialClear(t *testing.T) {
	assert := assert.New(t)
	NBTRIALS := 10
	var N int32 = 1024
	for trials := 0; trials < NBTRIALS; trials++ {
		a := NewTorusPolynomial(N)
		zero := NewTorusPolynomial(N)
		afft := NewLagrangeHalfCPolynomial(N)
		LagrangeHalfCPolynomialClear(afft)
		torusPolynomialUniform(a)
		torusPolynomialClear(zero)
		TorusPolynomial_fft(a, afft)
		assert.EqualValues(torusPolynomialNormInftyDist(zero, a), 0)
	}
}

/** sets to this torus32 constant */
//EXPORT void LagrangeHalfCPolynomialSetTorusConstant(LagrangeHalfCPolynomial* result, const Torus32 mu);
func TestLagrangeHalfCPolynomialSetTorusConstant(t *testing.T) {
	assert := assert.New(t)
	NBTRIALS := 10
	var N int32 = 1024
	for trials := 0; trials < NBTRIALS; trials++ {
		mu := UniformTorus32Dist()
		a := NewTorusPolynomial(N)
		cste := NewTorusPolynomial(N)
		afft := NewLagrangeHalfCPolynomial(N)
		torusPolynomialUniform(a)

		//tested function
		LagrangeHalfCPolynomialSetTorusConstant(afft, mu)
		TorusPolynomial_fft(a, afft)

		//expected result
		torusPolynomialClear(cste)
		cste.CoefsT[0] = mu

		assert.EqualValues(torusPolynomialNormInftyDist(cste, a), 0)
	}
}

//EXPORT void LagrangeHalfCPolynomialAddTorusConstant(LagrangeHalfCPolynomial* result, const Torus32 cst);
func TestLagrangeHalfCPolynomialAddTorusConstant(t *testing.T) {
	assert := assert.New(t)
	NBTRIALS := 10
	var N int32 = 1024
	toler := 1e-9
	for trials := 0; trials < NBTRIALS; trials++ {
		mu := UniformTorus32Dist()
		a := NewTorusPolynomial(N)
		aPlusCste := NewTorusPolynomial(N)
		b := NewTorusPolynomial(N)
		afft := NewLagrangeHalfCPolynomial(N)

		torusPolynomialUniform(a)
		TorusPolynomial_ifft(afft, a)
		LagrangeHalfCPolynomialAddTorusConstant(afft, mu)
		TorusPolynomial_fft(b, afft)

		torusPolynomialCopy(aPlusCste, a)
		aPlusCste.CoefsT[0] += mu

		assert.LessOrEqual(torusPolynomialNormInftyDist(aPlusCste, b), toler)
	}
}
