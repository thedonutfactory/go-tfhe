package tfhe

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

/** This structure represents an integer polynomial modulo X^N+1 */
type IntPolynomial struct {
	N     int
	Coefs []int
}

/** This structure represents an torus polynomial modulo X^N+1 */
type TorusPolynomial struct {
	N     int32
	Coefs []Torus32
}

func NewTorusPolynomial(n int32) *TorusPolynomial {
	return &TorusPolynomial{N: n, Coefs: make([]Torus32, n)}
}

func NewTorusPolynomialArray(size int, n int) (arr []TorusPolynomial) {
	arr = make([]TorusPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = TorusPolynomial{N: n, Coefs: make([]Torus32, n)}
	}
	return
}

func NewIntPolynomial(n int) *IntPolynomial {
	return &IntPolynomial{N: n, Coefs: make([]int, n)}
}

func NewIntPolynomialArray(size int, n int) (arr []IntPolynomial) {
	arr = make([]IntPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = *NewIntPolynomial(n)
	}
	return
}

// TorusPolynomial = 0
func torusPolynomialClear(result *TorusPolynomial) {
	for i := int32(0); i < result.N; i++ {
		result.Coefs[i] = 0
	}
}

// TorusPolynomial = random
func torusPolynomialUniform(result *TorusPolynomial) {
	//x := result.Coefs
	dist := distuv.Uniform{
		Min: math.MinInt,
		Max: math.MaxInt,
	}
	for i := int32(0); i < result.N; i++ {
		result.Coefs[i] = Torus32(dist.Rand())
	}
}

// TorusPolynomial = TorusPolynomial
func TorusPolynomialCopy(result *TorusPolynomial, sample *TorusPolynomial) {
	//assert(result != sample)
	if result == sample {
		panic("result == sample")
	}
	s := sample.Coefs
	r := result.Coefs
	for i := int32(0); i < result.N; i++ {
		r[i] = s[i]
	}
}

// TorusPolynomial + TorusPolynomial
func TorusPolynomialAdd(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	Assert(result != poly1) //if it fails here, please use addTo
	Assert(result != poly2) //if it fails here, please use addTo
	r := result.Coefs
	a := poly1.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] + b[i]
	}
}

// TorusPolynomial += TorusPolynomial
func TorusPolynomialAddTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	//r := result.Coefs
	//b := poly2.Coefs
	for i := int32(0); i < poly2.N; i++ {
		result.Coefs[i] += poly2.Coefs[i]
	}
}

// TorusPolynomial - TorusPolynomial
func TorusPolynomialSub(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	//assert(result != poly1); //if it fails here, please use subTo
	//assert(result != poly2); //if it fails here, please use subTo
	if result == poly1 || result == poly2 {
		panic("result == poly1 || result == poly2")
	}
	r := result.Coefs
	a := poly1.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] - b[i]
	}
}

// TorusPolynomial -= TorusPolynomial
func TorusPolynomialSubTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	r := result.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly2.N; i++ {
		r[i] -= b[i]
	}
}

// TorusPolynomial + p*TorusPolynomial
func TorusPolynomialAddMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.Coefs
	a := poly1.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] + p*b[i]
	}
}

// TorusPolynomial += p*TorusPolynomial
func TorusPolynomialAddMulZTo(result *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly2.N; i++ {
		r[i] += p * b[i]
	}
}

// TorusPolynomial - p*TorusPolynomial
func TorusPolynomialSubMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.Coefs
	a := poly1.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] - p*b[i]
	}
}

//result= (X^{a}-1)*source
func TorusPolynomialMulByXaiMinusOne(result *TorusPolynomial, a int, source *TorusPolynomial) {
	N := source.N
	out := result.Coefs
	in := source.Coefs

	//assert(a >= 0 && a < 2 * N)
	if a < 0 || a > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if a < N {
		for i := int(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N] - in[i]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a] - in[i]
		}
	} else {
		aa := a - N
		for i := int(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N] - in[i]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa] - in[i]
		}
	}
}

//result= X^{a}*source
func TorusPolynomialMulByXai(result *TorusPolynomial, a int, source *TorusPolynomial) {
	N := source.N
	out := result.Coefs
	in := source.Coefs

	//assert(a >= 0 && a < 2 * N)
	if a < 0 || a > 2*N {
		panic("a < 0 || a > 2 * N")
	}
	//assert(result != source)
	if result == source {
		panic("result == source")
	}

	if a < N {
		for i := int(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a]
		}
	} else {
		aa := a - N
		for i := int(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa]
		}
	}
}

// TorusPolynomial -= p*TorusPolynomial
func TorusPolynomialSubMulZTo(result *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.Coefs
	b := poly2.Coefs
	for i := int32(0); i < poly2.N; i++ {
		r[i] -= p * b[i]
	}
}

// Norme Euclidienne d'un IntPolynomial
func intPolynomialNormSq2(poly *IntPolynomial) int {
	var temp1 int = 0
	for i := int(0); i < poly.N; i++ {
		temp0 := poly.Coefs[i] * poly.Coefs[i]
		temp1 += temp0
	}
	return temp1
}

// Sets to zero
func intPolynomialClear(poly *IntPolynomial) {
	for i := int(0); i < poly.N; i++ {
		poly.Coefs[i] = 0
	}
}

// Sets to zero
func intPolynomialCopy(result *IntPolynomial, source *IntPolynomial) {
	for i := int(0); i < source.N; i++ {
		result.Coefs[i] = source.Coefs[i]
	}
}

/** accum += source */
func intPolynomialAddTo(accum *IntPolynomial, source *IntPolynomial) {
	for i := int(0); i < source.N; i++ {
		accum.Coefs[i] += source.Coefs[i]
	}
}

/**  result = (X^ai-1) * source */
func intPolynomialMulByXaiMinusOne(result *IntPolynomial, ai int, source *IntPolynomial) {
	N := source.N
	out := result.Coefs
	in := source.Coefs

	//assert(ai >= 0 && ai < 2 * N)
	if ai < 0 || ai > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if ai < N {
		for i := int(0); i < ai; i++ { //sur que i-a<0
			out[i] = -in[i-ai+N] - in[i]
		}
		for i := ai; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-ai] - in[i]
		}
	} else {
		aa := ai - N
		for i := int(0); i < aa; i++ { //sur que i-a<0
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
	for i := int32(0); i < poly1.N; i++ {
		r := math.Abs(T32tod(poly1.Coefs[i] - poly2.Coefs[i]))
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
	for i := int32(1); i < N; i++ {
		r := math.Abs(TorusToDouble(poly1.Coefs[i] - poly2.Coefs[i]))
		fmt.Printf("%d, %d => %f \n", poly1.Coefs[i], poly2.Coefs[i], r)
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
	for i := int32(0); i < N; i++ {
		r := math.Abs(TorusToDouble(poly1.Coefs[i] - poly2.Coefs[i]))
		if r > norm {
			norm = r
		}
	}
	return norm
}

// Norme 2 d'un IntPolynomial
func intPolynomialNorm2sq(poly *IntPolynomial) double {
	var norm double = 0
	for i := int(0); i < poly.N; i++ {
		r := poly.Coefs[i]
		norm += double(r * r)
	}
	return norm
}

// Norme infini de la distance entre deux IntPolynomial
func intPolynomialNormInftyDist(poly1 *IntPolynomial, poly2 *IntPolynomial) double {
	var norm double = 0
	// Max between the coefficients of abs(poly1-poly2)
	for i := int(0); i < poly1.N; i++ {
		r := Abs(poly1.Coefs[i] - poly2.Coefs[i])
		if double(r) > norm {
			norm = double(r)
		}
	}
	return norm
}

func TorusPolynomialMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	result.Coefs = Multiply(poly1.Coefs, poly2.Coefs)
}

func TorusPolynomialAddMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	tmpr := NewTorusPolynomial(poly1.N)
	tmpr.Coefs = Multiply(poly1.Coefs, poly2.Coefs)
	TorusPolynomialAddTo(result, tmpr)
}

func TorusPolynomialSubMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	tmpr := NewTorusPolynomial(poly1.N)
	tmpr.Coefs = Multiply(poly1.Coefs, poly2.Coefs)
	TorusPolynomialSubTo(result, tmpr)
}

func torusPolynomialMultFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialMulR(result, poly1, poly2)
}

func torusPolynomialAddMulRFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialAddMulR(result, poly1, poly2)
}

func torusPolynomialSubMulRFFT(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	TorusPolynomialSubMulR(result, poly1, poly2)
}
