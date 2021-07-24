package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

/** This structure represents an integer polynomial modulo X^N+1 */
type IntPolynomial struct {
	N     int32
	Coefs []int32
}

/** This structure represents an torus polynomial modulo X^N+1 */
type TorusPolynomial struct {
	N      int32
	CoefsT []Torus32
}

func NewTorusPolynomial(n int32) *TorusPolynomial {
	return &TorusPolynomial{N: n, CoefsT: make([]Torus32, n)}
}

func NewTorusPolynomialArray(size int, n int32) (arr []TorusPolynomial) {
	arr = make([]TorusPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = TorusPolynomial{N: n, CoefsT: make([]Torus32, n)}
	}
	return
}

func NewIntPolynomial(n int32) *IntPolynomial {
	return &IntPolynomial{N: n, Coefs: make([]int32, n)}
}

func NewIntPolynomialArray(size int, n int32) (arr []IntPolynomial) {
	arr = make([]IntPolynomial, size)
	for i := 0; i < size; i++ {
		arr[i] = *NewIntPolynomial(n)
	}
	return
}

// TorusPolynomial = 0
func torusPolynomialClear(result *TorusPolynomial) {
	for i := int32(0); i < result.N; i++ {
		result.CoefsT[i] = 0
	}
}

// TorusPolynomial = random
func torusPolynomialUniform(result *TorusPolynomial) {
	//x := result.CoefsT
	dist := distuv.Uniform{
		Min: math.MinInt32,
		Max: math.MaxInt32,
	}
	for i := int32(0); i < result.N; i++ {
		result.CoefsT[i] = Torus32(dist.Rand())
	}
}

// TorusPolynomial = TorusPolynomial
func torusPolynomialCopy(result *TorusPolynomial, sample *TorusPolynomial) {
	//assert(result != sample)
	if result == sample {
		panic("result == sample")
	}
	s := sample.CoefsT
	r := result.CoefsT
	for i := int32(0); i < result.N; i++ {
		r[i] = s[i]
	}
}

// TorusPolynomial + TorusPolynomial
func torusPolynomialAdd(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	Assert(result != poly1) //if it fails here, please use addTo
	Assert(result != poly2) //if it fails here, please use addTo
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] + b[i]
	}
}

// TorusPolynomial += TorusPolynomial
func torusPolynomialAddTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	//r := result.CoefsT
	//b := poly2.CoefsT
	for i := int32(0); i < poly2.N; i++ {
		result.CoefsT[i] += poly2.CoefsT[i]
	}
}

// TorusPolynomial - TorusPolynomial
func torusPolynomialSub(result *TorusPolynomial, poly1 *TorusPolynomial, poly2 *TorusPolynomial) {
	//assert(result != poly1); //if it fails here, please use subTo
	//assert(result != poly2); //if it fails here, please use subTo
	if result == poly1 || result == poly2 {
		panic("result == poly1 || result == poly2")
	}
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] - b[i]
	}
}

// TorusPolynomial -= TorusPolynomial
func torusPolynomialSubTo(result *TorusPolynomial, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly2.N; i++ {
		r[i] -= b[i]
	}
}

// TorusPolynomial + p*TorusPolynomial
func torusPolynomialAddMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] + p*b[i]
	}
}

// TorusPolynomial += p*TorusPolynomial
func torusPolynomialAddMulZTo(result *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly2.N; i++ {
		r[i] += p * b[i]
	}
}

// TorusPolynomial - p*TorusPolynomial
func torusPolynomialSubMulZ(result *TorusPolynomial, poly1 *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.CoefsT
	a := poly1.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly1.N; i++ {
		r[i] = a[i] - p*b[i]
	}
}

//result= (X^{a}-1)*source
func torusPolynomialMulByXaiMinusOne(result *TorusPolynomial, a int32, source *TorusPolynomial) {
	N := source.N
	out := result.CoefsT
	in := source.CoefsT

	//assert(a >= 0 && a < 2 * N)
	if a < 0 || a > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if a < N {
		for i := int32(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N] - in[i]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a] - in[i]
		}
	} else {
		aa := a - N
		for i := int32(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N] - in[i]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa] - in[i]
		}
	}
}

//result= X^{a}*source
func torusPolynomialMulByXai(result *TorusPolynomial, a int32, source *TorusPolynomial) {
	N := source.N
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
		for i := int32(0); i < a; i++ { //sur que i-a<0
			out[i] = -in[i-a+N]
		}
		for i := a; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-a]
		}
	} else {
		aa := a - N
		for i := int32(0); i < aa; i++ { //sur que i-a<0
			out[i] = in[i-aa+N]
		}
		for i := aa; i < N; i++ { //sur que N>i-a>=0
			out[i] = -in[i-aa]
		}
	}
}

// TorusPolynomial -= p*TorusPolynomial
func torusPolynomialSubMulZTo(result *TorusPolynomial, p int32, poly2 *TorusPolynomial) {
	r := result.CoefsT
	b := poly2.CoefsT
	for i := int32(0); i < poly2.N; i++ {
		r[i] -= p * b[i]
	}
}

// Norme Euclidienne d'un IntPolynomial
func intPolynomialNormSq2(poly *IntPolynomial) int32 {
	var temp1 int32 = 0
	for i := int32(0); i < poly.N; i++ {
		temp0 := poly.Coefs[i] * poly.Coefs[i]
		temp1 += temp0
	}
	return temp1
}

// Sets to zero
func intPolynomialClear(poly *IntPolynomial) {
	for i := int32(0); i < poly.N; i++ {
		poly.Coefs[i] = 0
	}
}

// Sets to zero
func intPolynomialCopy(result *IntPolynomial, source *IntPolynomial) {
	for i := int32(0); i < source.N; i++ {
		result.Coefs[i] = source.Coefs[i]
	}
}

/** accum += source */
func intPolynomialAddTo(accum *IntPolynomial, source *IntPolynomial) {
	for i := int32(0); i < source.N; i++ {
		accum.Coefs[i] += source.Coefs[i]
	}
}

/**  result = (X^ai-1) * source */
func intPolynomialMulByXaiMinusOne(result *IntPolynomial, ai int32, source *IntPolynomial) {
	N := source.N
	out := result.Coefs
	in := source.Coefs

	//assert(ai >= 0 && ai < 2 * N)
	if ai < 0 || ai > 2*N {
		panic("a < 0 || a > 2 * N")
	}

	if ai < N {
		for i := int32(0); i < ai; i++ { //sur que i-a<0
			out[i] = -in[i-ai+N] - in[i]
		}
		for i := ai; i < N; i++ { //sur que N>i-a>=0
			out[i] = in[i-ai] - in[i]
		}
	} else {
		aa := ai - N
		for i := int32(0); i < aa; i++ { //sur que i-a<0
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
		r := math.Abs(T32tod(poly1.CoefsT[i] - poly2.CoefsT[i]))
		if r > norm {
			norm = r
		}
	}
	return norm
}
*/

func torusPolynomialNormInftyDist(poly1 *TorusPolynomial, poly2 *TorusPolynomial) double {
	N := poly1.N
	var norm double = 0

	// Max between the coefficients of abs(poly1-poly2)
	for i := int32(0); i < N; i++ {
		r := math.Abs(T32tod(poly1.CoefsT[i] - poly2.CoefsT[i]))
		if r > norm {
			norm = r
		}
	}
	return norm
}

// Norme 2 d'un IntPolynomial
func intPolynomialNorm2sq(poly *IntPolynomial) double {
	var norm double = 0
	for i := int32(0); i < poly.N; i++ {
		r := poly.Coefs[i]
		norm += double(r * r)
	}
	return norm
}

// Norme infini de la distance entre deux IntPolynomial
func intPolynomialNormInftyDist(poly1 *IntPolynomial, poly2 *IntPolynomial) double {
	var norm double = 0
	// Max between the coefficients of abs(poly1-poly2)
	for i := int32(0); i < poly1.N; i++ {
		r := Abs(poly1.Coefs[i] - poly2.Coefs[i])
		if double(r) > norm {
			norm = double(r)
		}
	}
	return norm
}

func LagrangeHalfCPolynomialMul(a []complex128, b []complex128, Ns2 int) (result *LagrangeHalfCPolynomial) {
	result = &LagrangeHalfCPolynomial{
		coefsC: make([]complex128, Ns2),
	}
	//rr := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		result.coefsC[i] = a[i] * b[i]
	}
	return
}

func IntPolynomial_ifft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	result.coefsC = fftProcessor.executeReverseInt(p.Coefs)
}

func TorusPolynomial_ifft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	result.coefsC = fftProcessor.executeReverseTorus32(p.CoefsT)
}

func TorusPolynomial_fft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	result.CoefsT = fftProcessor.executeDirectTorus32(p.coefsC)
}

/*
	const int32_t N = poly1->N;
    LagrangeHalfCPolynomial* tmp = new_LagrangeHalfCPolynomial_array(3,N);
    IntPolynomial_ifft(tmp+0,poly1);
    TorusPolynomial_ifft(tmp+1,poly2);
    LagrangeHalfCPolynomialMul(tmp+2,tmp+0,tmp+1);
    TorusPolynomial_fft(result, tmp+2);
    delete_LagrangeHalfCPolynomial_array(3,tmp);
*/

func TorusPolynomialMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	tmp := []*LagrangeHalfCPolynomial{
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
	}
	IntPolynomial_ifft(tmp[0], poly1)
	TorusPolynomial_ifft(tmp[1], poly2)
	tmp[2] = LagrangeHalfCPolynomialMul(tmp[0].coefsC, tmp[1].coefsC, int(N/2))
	TorusPolynomial_fft(result, tmp[2])
}

func TorusPolynomialAddMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	tmp := []*LagrangeHalfCPolynomial{
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
	}
	tmpr := NewTorusPolynomial(N)
	IntPolynomial_ifft(tmp[0], poly1)
	TorusPolynomial_ifft(tmp[1], poly2)
	tmp[2] = LagrangeHalfCPolynomialMul(tmp[0].coefsC, tmp[1].coefsC, int(N/2))
	TorusPolynomial_fft(tmpr, tmp[2])
	torusPolynomialAddTo(result, tmpr)
}

func TorusPolynomialSubMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	tmp := []*LagrangeHalfCPolynomial{
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
		NewLagrangeHalfCPolynomial(N),
	}
	tmpr := NewTorusPolynomial(N)
	IntPolynomial_ifft(tmp[0], poly1)
	TorusPolynomial_ifft(tmp[1], poly2)
	tmp[2] = LagrangeHalfCPolynomialMul(tmp[0].coefsC, tmp[1].coefsC, int(N/2))
	TorusPolynomial_fft(tmpr, tmp[2])
	torusPolynomialSubTo(result, tmpr)
}
