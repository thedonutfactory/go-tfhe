package tfhe

import (
	"github.com/mjibson/go-dsp/fft"
)

// the default fft processor
var fftProc FftProcessorInterface = &GoDspFftProcessor{}

type FftProcessorInterface interface {
	executeReverseInt(a []int32) []complex128
	executeReverseTorus32(a []Torus32) []complex128
	executeDirectTorus32(a []complex128) []Torus32
}

func intPolynomialIfft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	result.Coefs = fftProc.executeReverseInt(p.Coefs)
}

func torusPolynomialIfft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	result.Coefs = fftProc.executeReverseTorus32(p.Coefs)
}

func torusPolynomialFft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	result.Coefs = fftProc.executeDirectTorus32(p.Coefs)
}

func Multiply(a, b []int32) []int32 {
	N := int32(len(a))
	poly1 := NewIntPolynomial(N)
	poly1.Coefs = a
	poly2 := NewTorusPolynomial(N)
	poly2.Coefs = b
	result := NewTorusPolynomial(N)
	j, k, l := NewLagrangeHalfCPolynomial(N), NewLagrangeHalfCPolynomial(N), NewLagrangeHalfCPolynomial(N)
	intPolynomialIfft(j, poly1)
	torusPolynomialIfft(k, poly2)
	LagrangeHalfCPolynomialMul(l, j, k)
	torusPolynomialFft(result, l)
	return result.Coefs
}

func MultiplyRef(a, b []int32) []int32 {
	n := len(a)
	a = Pad(a)
	b = Pad(b)
	x := fft.FFT(CastComplex(a))
	y := fft.FFT(CastComplex(b))
	c := Mult(x, y)
	z := fft.IFFT(c)

	res := make([]int32, n)
	for i := range res {
		t := real(z[i] - z[n+i])
		res[i] = int32(int64(t))
	}
	return res
}

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

func addTo(a, b []complex128) []complex128 {
	n := max(len(a), len(b))
	c := make([]complex128, n)
	for i := 0; i < n; i++ {
		a[i] += b[i]
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

func Pad(a []int32) []int32 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return a
}

func AddTo(a, b []int32) []int32 {
	n := len(a)
	a = Pad(a)
	b = Pad(b)
	x := fft.FFT(CastComplex(a))
	y := fft.FFT(CastComplex(b))
	c := addTo(x, y)
	z := fft.IFFT(c)

	res := make([]int32, n)
	for i := range res {
		t := real(z[i] - z[n+i])
		res[i] = int32(int64(t))
	}
	return res
}
