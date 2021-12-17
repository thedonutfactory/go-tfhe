package tfhe

import (
	"github.com/mjibson/go-dsp/fft"
)

func Mulfft(a []complex128) []complex128 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return fft.FFT(a)
}

func Mult(a, b []complex128) []complex128 {
	n := max(len(a), len(b))
	c := make([]complex128, n)
	for i := 0; i < n; i++ {
		c[i] = a[i] * b[i]
	}
	return c
}

// Max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func CastComplex(arr []int32) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		res[i] = complex(float64(v), 0.)
	}
	return
}

func CastTorus(arr []complex128) (res []int32) {
	res = make([]int32, len(arr))
	for i, v := range arr {
		res[i] = int32(real(v))
	}
	return
}

func Multiply(a, b []int32) []int32 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
		b = append(b, 0)
	}
	x := fft.FFT(CastComplex(a))
	y := fft.FFT(CastComplex(b))
	c := Mult(x, y)
	z := fft.IFFT(c)

	res := make([]int32, n)
	for i := range res {
		t := real(z[i] + (z[n+i] * -1))
		res[i] = int32(int64(t))
	}
	return res
}

func executeReverseTorus32(a []Torus32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func executeReverseInt(a []int32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func executeDirectTorus32(a []complex128) (res []Torus32) {
	Ns2 := len(a) * 2
	res = castTorus(fft.FFT(a))
	for i := 0; i < int(Ns2); i++ {
		res = append(res, 0)
	}
	return

}

/**
 * FFT functions
 */

func intPolynomialIfft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	result.coefsC = executeReverseInt(p.Coefs)
}

func torusPolynomialIfft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	result.coefsC = executeReverseTorus32(p.Coefs)
}

func torusPolynomialFft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	result.Coefs = executeDirectTorus32(p.coefsC)
}
