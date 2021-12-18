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

func executeReverseInt(res []complex128, a []int32) {
	N := len(a)
	Ns2 := N / 2
	_2N := N * 2
	cplxInout := make([]complex128, _2N)
	for i := 0; i < N; i++ {
		cplxInout[i] = complex(float64(a[i])/2., 0.)
	}
	for i := 0; i < N; i++ {
		cplxInout[N+i] = complex(-(float64(a[i]) / 2), 0.)
	}
	z := fft.FFT(cplxInout)
	for i := 0; i < Ns2; i++ {
		res[i] = z[2*i+1]
	}
}

func executeReverseTorus32(res []complex128, a []Torus32) {
	N := len(a)
	Ns2 := N / 2
	_2N := N * 2
	cplxInout := make([]complex128, _2N)
	for i := 0; i < N; i++ {
		t := float64(a[i]) / float64(int64(1)<<33)
		cplxInout[i] = complex(t, 0.)
	}
	for i := 0; i < N; i++ {
		t := float64(a[i]) / float64(int64(1)<<33)
		cplxInout[N+i] = complex(-t, 0.)
	}
	z := fft.FFT(cplxInout)
	for i := 0; i < Ns2; i++ {
		res[i] = z[2*i+1]
	}
}

func executeDirectTorus32(res []Torus32, a []complex128) {
	N := len(a) * 2
	Ns2 := N / 2
	_2N := N * 2
	_2p32 := float64(int64(1) << 32)
	_1sN := float64(1) / double(N)

	cplxInout := make([]complex128, N*2)
	for i := 0; i < N; i++ {
		cplxInout[2*i] = complex(0., 0.)
	}
	for i := 0; i < Ns2; i++ {
		cplxInout[2*i+1] = complex(real(a[i]), imag(a[i]))
	}
	for i := 0; i < Ns2; i++ {
		cplxInout[_2N-1-2*i] = complex(real(a[i]), -imag(a[i]))
	}
	z := fft.FFT(cplxInout)
	res[0] = int32(int64(real(z[0]) * _1sN * _2p32))
	for i := 1; i < N; i++ {
		res[i] = -int32(int64(real(z[N-i]) * _1sN * _2p32))
	}
}

/**
 * FFT functions
 */

func intPolynomialIfft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	executeReverseInt(result.coefsC, p.Coefs)
}

func torusPolynomialIfft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	executeReverseTorus32(result.coefsC, p.Coefs)
}

func torusPolynomialFft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	executeDirectTorus32(result.Coefs, p.coefsC)
}
