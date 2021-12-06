package tfhe

import (
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
	"gonum.org/v1/gonum/stat/distuv"
)

/** This structure represents an integer polynomial modulo X^N+1 */
type IntPolynomial struct {
	N     int
	Coefs []int64
}

/** This structure represents an torus polynomial modulo X^N+1 */
type TorusPolynomial struct {
	N      int
	CoefsT []Torus
}

func NewTorusPolynomial(n int) *TorusPolynomial {
	return &TorusPolynomial{N: n, CoefsT: make([]Torus, n)}
}

func NewTorusPolynomialArray(size, n int) (arr []TorusPolynomial) {
	arr = make([]TorusPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = TorusPolynomial{N: n, CoefsT: make([]Torus, n)}
	}
	return
}

func NewIntPolynomial(n int) *IntPolynomial {
	return &IntPolynomial{N: n, Coefs: make([]int64, n)}
}

func NewIntPolynomialArray(size, n int) (arr []IntPolynomial) {
	arr = make([]IntPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = *NewIntPolynomial(n)
	}
	return
}

// TorusPolynomial = 0
func torusPolynomialClear(result *TorusPolynomial) {
	for i := 0; i < result.N; i++ {
		result.CoefsT[i] = 0
	}
}

// TorusPolynomial = random
func torusPolynomialUniform(result *TorusPolynomial) {
	//x := result.CoefsT
	dist := distuv.Uniform{
		Min: math.MinInt64,
		Max: math.MaxInt64,
	}
	for i := 0; i < result.N; i++ {
		result.CoefsT[i] = Torus(dist.Rand())
	}
}

// TorusPolynomial = TorusPolynomial
func TorusPolynomialCopy(result *TorusPolynomial, sample *TorusPolynomial) {
	//assert(result != sample)
	if result == sample {
		panic("result == sample")
	}
	s := sample.CoefsT
	r := result.CoefsT
	for i := 0; i < result.N; i++ {
		r[i] = s[i]
	}
}

// TorusPolynomial + TorusPolynomial
func TorusPolynomialAdd(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	Assert(result != poly1) //if it fails here, please use addTo
	Assert(result != poly2) //if it fails here, please use addTo
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly1.N; i++ {
		r[i] = a[i] + b[i]
	}
}

// TorusPolynomial += TorusPolynomial
func TorusPolynomialAddTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	//r := result.CoefsT
	//b := poly2.CoefsT
	for i := 0; i < poly2.N; i++ {
		result.CoefsT[i] += poly2.CoefsT[i]
	}
}

// TorusPolynomial - TorusPolynomial
func TorusPolynomialSub(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	//assert(result != poly1); //if it fails here, please use subTo
	//assert(result != poly2); //if it fails here, please use subTo
	if result == poly1 || result == poly2 {
		panic("result == poly1 || result == poly2")
	}
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly1.N; i++ {
		r[i] = a[i] - b[i]
	}
}

// TorusPolynomial -= TorusPolynomial
func TorusPolynomialSubTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly2.N; i++ {
		r[i] -= b[i]
	}
}

// TorusPolynomial + p*TorusPolynomial
func TorusPolynomialAddMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int64, poly2 *TorusPolynomial) {
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly1.N; i++ {
		r[i] = a[i] + p*b[i]
	}
}

// TorusPolynomial += p*TorusPolynomial
func TorusPolynomialAddMulZTo(result *TorusPolynomial, p int64, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly2.N; i++ {
		r[i] += p * b[i]
	}
}

// TorusPolynomial - p*TorusPolynomial
func TorusPolynomialSubMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int64, poly2 *TorusPolynomial) {
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly1.N; i++ {
		r[i] = a[i] - p*b[i]
	}
}

//result= (X^{a}-1)*source
func TorusPolynomialMulByXaiMinusOne(result *TorusPolynomial, a int64, source *TorusPolynomial) {
	N := int64(source.N)
	out := result.CoefsT
	in := source.CoefsT

	//assert(a >= 0 && a < 2 * N)
	if a < 0 || a > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if a < N {
		for i := int64(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N] - in[i]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a] - in[i]
		}
	} else {
		aa := a - N
		for i := int64(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N] - in[i]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa] - in[i]
		}
	}
}

//result= X^{a}*source
func TorusPolynomialMulByXai(result *TorusPolynomial, a int64, source *TorusPolynomial) {
	N := int64(source.N)
	out := result.CoefsT
	in := source.CoefsT

	//assert(a >= 0 && a < 2 * N)
	if a < 0 || a > 2*N {
		panic("a < 0 || a > 2 * N")
	}
	//assert(result != source)
	if result == source {
		panic("result == source")
	}

	if a < N {
		for i := int64(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a]
		}
	} else {
		aa := a - N
		for i := int64(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa]
		}
	}
}

// TorusPolynomial -= p*TorusPolynomial
func TorusPolynomialSubMulZTo(result *TorusPolynomial, p int64, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := 0; i < poly2.N; i++ {
		r[i] -= p * b[i]
	}
}

// Norme Euclidienne d'un IntPolynomial
func intPolynomialNormSq2(poly *IntPolynomial) int64 {
	var temp1 int64 = 0
	for i := 0; i < poly.N; i++ {
		temp0 := poly.Coefs[i] * poly.Coefs[i]
		temp1 += temp0
	}
	return temp1
}

// Sets to zero
func intPolynomialClear(poly *IntPolynomial) {
	for i := 0; i < poly.N; i++ {
		poly.Coefs[i] = 0
	}
}

// Sets to zero
func intPolynomialCopy(result *IntPolynomial, source *IntPolynomial) {
	for i := 0; i < source.N; i++ {
		result.Coefs[i] = source.Coefs[i]
	}
}

/** accum += source */
func intPolynomialAddTo(accum *IntPolynomial, source *IntPolynomial) {
	for i := 0; i < source.N; i++ {
		accum.Coefs[i] += source.Coefs[i]
	}
}

/**  result = (X^ai-1) * source */
func intPolynomialMulByXaiMinusOne(result *IntPolynomial, ai int64, source *IntPolynomial) {
	N := int64(source.N)
	out := result.Coefs
	in := source.Coefs

	//assert(ai >= 0 && ai < 2 * N)
	if ai < 0 || ai > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if ai < N {
		for i := int64(0); i < ai; i++ { //sur que i-a<0
			out[i] = -in[i-ai+N] - in[i]
		}
		for i := ai; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-ai] - in[i]
		}
	} else {
		aa := ai - N
		for i := int64(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N] - in[i]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa] - in[i]
		}
	}
}

// Norme infini de la distance entre deux TorusPolynomial
/*
func torusPolynomialNormInftyDist(poly1 *TorusPolynomial, poly2 *TorusPolynomial) double {
	var norm double = 0
	// Max between the coefficients of abs(poly1-poly2)
	for i := 0; i < poly1.N; i++ {
		r := math.Abs(T32tod(poly1.CoefsT[i] - poly2.CoefsT[i]))
		if r > norm {
			norm = r
		}
	}
	return norm
}
*/

func torusPolynomialNormInftyDistSkipFirst(poly1 *TorusPolynomial, poly2 *TorusPolynomial) double {
	N := poly1.N
	var norm double = 0

	// Max between the coefficients of abs(poly1-poly2)
	fmt.Println("Warning, skipping 0th element in torusPolynomialNormInftyDist")
	for i := 1; i < N; i++ {
		r := math.Abs(T32tod(poly1.CoefsT[i] - poly2.CoefsT[i]))
		fmt.Printf("%d, %d => %f \n", poly1.CoefsT[i], poly2.CoefsT[i], r)
		if r > norm {
			norm = r
		}
	}
	return norm
}

func torusPolynomialNormInftyDist(poly1 *TorusPolynomial, poly2 *TorusPolynomial) double {
	N := poly1.N
	var norm double = 0

	// Max between the coefficients of abs(poly1-poly2)
	for i := 0; i < N; i++ {
		r := math.Abs(T32tod(poly1.CoefsT[i] - poly2.CoefsT[i]))
		fmt.Printf("%d, %d => %f \n", poly1.CoefsT[i], poly2.CoefsT[i], r)
		if r > norm {
			norm = r
		}
	}
	return norm
}

// Norme 2 d'un IntPolynomial
func intPolynomialNorm2sq(poly *IntPolynomial) double {
	var norm double = 0
	for i := 0; i < poly.N; i++ {
		r := poly.Coefs[i]
		norm += double(r * r)
	}
	return norm
}

// Norme infini de la distance entre deux IntPolynomial
func intPolynomialNormInftyDist(poly1 *IntPolynomial, poly2 *IntPolynomial) double {
	var norm double = 0
	// Max between the coefficients of abs(poly1-poly2)
	for i := 0; i < poly1.N; i++ {
		r := Abs(poly1.Coefs[i] - poly2.Coefs[i])
		if double(r) > norm {
			norm = double(r)
		}
	}
	return norm
}

func mulfft(a []complex128) []complex128 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return fft.FFT(a)
}

func invfft(a []complex128) []complex128 {
	return fft.IFFT(a)
}

func mult(a, b []complex128) []complex128 {
	n := Max(len(a), len(b))
	c := make([]complex128, n)
	for i := 0; i < n; i++ {
		c[i] = a[i] * b[i]
	}
	return c
}

func multiply(a, b []int64) []int64 {
	/*
		x := mulfft(castComplex(a))
		y := mulfft(castComplex(b))
		c := mult(x, y)
		return castTorus(invfft(c))
	*/
	x := revInt(a)
	y := revTorus(b)
	c := mult(x, y)
	return dirTorus(c)
}

func revTorus(a []Torus) []complex128 {
	N := len(a)
	Ns2 := len(a) / 2
	_2pm33 := 1. / double(int64(1)<<33)
	revIn := make([]complex128, len(a)*2)

	for i := 0; i < N; i++ {
		revIn[i] = complex(float64(a[i])*_2pm33, 0.)
	}
	for i := 0; i < N; i++ {
		revIn[N+i] = -revIn[i]
	}

	revOutCplx := fft.FFT(revIn)

	res := make([]complex128, len(a))
	for i := 0; i < Ns2; i++ {
		res[i] = revOutCplx[2*i+1]
	}

	return res
}

func revInt(a []int64) []complex128 {
	N := len(a)
	Ns2 := len(a) / 2
	revIn := make([]complex128, len(a)*2)

	for i := 0; i < N; i++ {
		revIn[i] = complex(float64(a[i])/2., 0.)
	}
	for i := 0; i < N; i++ {
		revIn[N+i] = -revIn[i]
	}

	revOutCplx := fft.FFT(revIn)

	res := make([]complex128, len(a))
	for i := 0; i < Ns2; i++ {
		res[i] = revOutCplx[2*i+1]
	}

	return res
}

func dirTorus(a []complex128) []Torus {
	N := len(a)
	Ns2 := len(a) / 2
	_2p32 := double(uint64(1) << 63)
	_1sN := double(1) / double(N)

	inCplx := make([]complex128, len(a)+1)
	for i := 0; i <= Ns2; i++ {
		inCplx[2*i] = 0
	}
	for i := 0; i < Ns2; i++ {
		inCplx[2*i+1] = a[i]
	}

	out := fft.FFT(inCplx)

	res := make([]Torus, N)
	for i := 0; i < N; i++ {
		res[i] = Torus(real(out[i]) * _1sN * _2p32)
	}
	return res
}

func TorusPolynomialMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	torusPolynomialMultKaratsuba(result, poly1, poly2)
	/*
		N := poly1.N
		tmp := []*LagrangeHalfCPolynomial{
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
		}
		intPolynomialIfft(tmp[0], poly1)
		torusPolynomialIfft(tmp[1], poly2)
		LagrangeHalfCPolynomialMul(tmp[2], tmp[0], tmp[1])
		torusPolynomialFft(result, tmp[2])
	*/
	//result.CoefsT = castTorus(fft.Convolve(castComplex(poly1.Coefs), castComplex(poly2.CoefsT)))
	//result.CoefsT = multiply(poly1.Coefs, poly2.CoefsT)
	//result.CoefsT = castTorus(r)
}

func TorusPolynomialAddMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	torusPolynomialAddMulRKaratsuba(result, poly1, poly2)
	/*
		N := poly1.N
		tmp := []*LagrangeHalfCPolynomial{
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
		}
		tmpr := NewTorusPolynomial(N)
		intPolynomialIfft(tmp[0], poly1)
		torusPolynomialIfft(tmp[1], poly2)
		LagrangeHalfCPolynomialMul(tmp[2], tmp[0], tmp[1])
		torusPolynomialFft(tmpr, tmp[2])
		torusPolynomialAddTo(result, tmpr)
	*/
}

func TorusPolynomialSubMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	torusPolynomialSubMulRKaratsuba(result, poly1, poly2)

	/*
		N := poly1.N
		tmp := []*LagrangeHalfCPolynomial{
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
			NewLagrangeHalfCPolynomial(N),
		}
		tmpr := NewTorusPolynomial(N)
		intPolynomialIfft(tmp[0], poly1)
		torusPolynomialIfft(tmp[1], poly2)
		LagrangeHalfCPolynomialMul(tmp[2], tmp[0], tmp[1])
		torusPolynomialFft(tmpr, tmp[2])
		torusPolynomialSubTo(result, tmpr)
	*/
}

/** multiplication via direct FFT - simple wrappers since we currently only have one implementation
TODO - implement FFT - currently placing Karatsuba implementation instead
*/
func torusPolynomialMultFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialMulR(result, poly1, poly2)
}

func torusPolynomialAddMulRFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialAddMulR(result, poly1, poly2)
}

func torusPolynomialSubMulRFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialSubMulR(result, poly1, poly2)
}
