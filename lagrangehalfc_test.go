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
