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
		Min: INT32_MIN,
		Max: INT32_MAX,
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

/*
void FFT_Processor_fftw::execute_reverse_int(cplx* res, const int* a) {
    cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<N; i++) rev_in[i]=a[i]/2.;
    for (int32_t i=0; i<N; i++) rev_in[N+i]=-rev_in[i];
    fftw_execute(rev_p);
    for (int32_t i=0; i<Ns2; i++) res[i]=rev_out_cplx[2*i+1];
    for (int32_t i=0; i<=Ns2; i++) assert(abs(rev_out_cplx[2*i])<1e-20);
}
*/

/**
 * FFT functions
 */
/*
func IntPolynomial_ifft(p *IntPolynomial) []complex128 {
	n := p.N
	x := make([]complex128, n)
	for k := int32(0); k < n; k++ {
		x[k] = complex(float64(p.Coefs[k]), 0.0)
	}

	// FFT
	y := fft.FFT(x, int(n))

	// IFFT
	z := fft.IFFT(y, int(n))

	return z

}

func TorusPolynomial_ifft(p *TorusPolynomial) []complex128 {
	n := p.N
	x := make([]complex128, n)
	for k := int32(0); k < n; k++ {
		x[k] = complex(float64(p.CoefsT[k]), 0.0)
	}
	// FFT
	y := fft.FFT(x, int(n))
	// IFFT
	z := fft.IFFT(y, int(n))
	return z
}

func TorusPolynomial_fft(y []complex128) *TorusPolynomial {
	n := len(y)
	//n := p.N
	l := fft.IFFT(y, n)

	//for (int32_t i=0; i<N; i++) res[i]=Torus32(int64_t(out[i]*_1sN*_2p32));
	//pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
	//for (int32_t i=0; i<N; i++) assert(fabs(out[N+i]+out[i])<1e-20);
	tp := NewTorusPolynomial(int32(n))
	for i, v := range l {
		tp.CoefsT[i] = int32(real(v))
	}
	return tp
}
*/

func LagrangeHalfCPolynomialMul(a []complex128, b []complex128, Ns2 int) (result *LagrangeHalfCPolynomial) {
	//LagrangeHalfCPolynomial_IMPL* result1 = (LagrangeHalfCPolynomial_IMPL*) result;
	//Ns2 := result1->proc->Ns2;
	//cplx* aa = ((LagrangeHalfCPolynomial_IMPL*) a)->CoefsC;
	//cplx* bb = ((LagrangeHalfCPolynomial_IMPL*) b)->CoefsC;
	//cplx* rr = result1->CoefsC;
	result = &LagrangeHalfCPolynomial{
		coefsC: make([]complex128, Ns2),
	}
	//rr := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		result.coefsC[i] = a[i] * b[i]
	}
	return
}

func IntPolynomial_ifft(n int32, p *IntPolynomial) (result *LagrangeHalfCPolynomial) {
	result = NewLangrangeHalfCPolynomial(n)
	fftProcessor.executeReverseInt(result.coefsC, p.Coefs)
	return
}

func TorusPolynomial_ifft(n int32, p *TorusPolynomial) (result *LagrangeHalfCPolynomial) {
	result = NewLangrangeHalfCPolynomial(n)
	fftProcessor.executeReverseTorus32(result.coefsC, p.CoefsT)
	return
}

func TorusPolynomial_fft(n int32, p *LagrangeHalfCPolynomial) (result *TorusPolynomial) {
	result = NewTorusPolynomial(n)
	fftProcessor.executeDirectTorus32(result.CoefsT, p.coefsC)
	return
}

/*
void FFT_Processor_fftw::execute_reverse_int(cplx* res, const int* a) {
    cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<N; i++) rev_in[i]=a[i]/2.;
    for (int32_t i=0; i<N; i++) rev_in[N+i]=-rev_in[i];
    fftw_execute(rev_p);
    for (int32_t i=0; i<Ns2; i++) res[i]=rev_out_cplx[2*i+1];
    for (int32_t i=0; i<=Ns2; i++) assert(abs(rev_out_cplx[2*i])<1e-20);
}
void FFT_Processor_fftw::execute_reverse_torus32(cplx* res, const Torus32* a) {
    static const double _2pm33 = 1./double(INT64_C(1)<<33);
    int32_t* aa = (int32_t*) a;
    cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<N; i++) rev_in[i]=aa[i]*_2pm33;
    for (int32_t i=0; i<N; i++) rev_in[N+i]=-rev_in[i];
    fftw_execute(rev_p);
    for (int32_t i=0; i<Ns2; i++) res[i]=rev_out_cplx[2*i+1];
    for (int32_t i=0; i<=Ns2; i++) assert(abs(rev_out_cplx[2*i])<1e-20);
}
void FFT_Processor_fftw::execute_direct_Torus32(Torus32* res, const cplx* a) {
    static const double _2p32 = double(INT64_C(1)<<32);
    static const double _1sN = double(1)/double(N);
    cplx* in_cplx = (cplx*) in; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<=Ns2; i++) in_cplx[2*i]=0;
    for (int32_t i=0; i<Ns2; i++) in_cplx[2*i+1]=a[i];
    fftw_execute(p);
    for (int32_t i=0; i<N; i++) res[i]=Torus32(int64_t(out[i]*_1sN*_2p32));
    //pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
    for (int32_t i=0; i<N; i++) assert(fabs(out[N+i]+out[i])<1e-20);
}

EXPORT void IntPolynomial_ifft(LagrangeHalfCPolynomial* result, const IntPolynomial* p) {
    fp1024_fftw.execute_reverse_int(((LagrangeHalfCPolynomial_IMPL*)result)->CoefsC, p->Coefs);
}
EXPORT void TorusPolynomial_ifft(LagrangeHalfCPolynomial* result, const TorusPolynomial* p) {
    fp1024_fftw.execute_reverse_torus32(((LagrangeHalfCPolynomial_IMPL*)result)->CoefsC, p->CoefsT);
}
EXPORT void TorusPolynomial_fft(TorusPolynomial* result, const LagrangeHalfCPolynomial* p) {
    fp1024_fftw.execute_direct_Torus32(result->CoefsT, ((LagrangeHalfCPolynomial_IMPL*)p)->CoefsC);
}
*/

/** multiplication via direct FFT (it must know the implem of LagrangeHalfCPolynomial because of the tmp+1 notation */
/*
EXPORT void torusPolynomialMultFFT(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2) {
    const int32_t N = poly1->N;
    LagrangeHalfCPolynomial* tmp = new_LagrangeHalfCPolynomial_array(3,N);
    IntPolynomial_ifft(tmp+0,poly1);
    TorusPolynomial_ifft(tmp+1,poly2);
    LagrangeHalfCPolynomialMul(tmp+2,tmp+0,tmp+1);
    TorusPolynomial_fft(result, tmp+2);
    delete_LagrangeHalfCPolynomial_array(3,tmp);
}
EXPORT void torusPolynomialAddMulRFFT(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2) {
    const int32_t N = poly1->N;
    LagrangeHalfCPolynomial* tmp = new_LagrangeHalfCPolynomial_array(3,N);
    TorusPolynomial* tmpr = new_TorusPolynomial(N);
    IntPolynomial_ifft(tmp+0,poly1);
    TorusPolynomial_ifft(tmp+1,poly2);
    LagrangeHalfCPolynomialMul(tmp+2,tmp+0,tmp+1);
    TorusPolynomial_fft(tmpr, tmp+2);
    torusPolynomialAddTo(result, tmpr);
    delete_TorusPolynomial(tmpr);
    delete_LagrangeHalfCPolynomial_array(3,tmp);
}
EXPORT void torusPolynomialSubMulRFFT(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2) {
    const int32_t N = poly1->N;
    LagrangeHalfCPolynomial* tmp = new_LagrangeHalfCPolynomial_array(3,N);
    TorusPolynomial* tmpr = new_TorusPolynomial(N);
    IntPolynomial_ifft(tmp+0,poly1);
    TorusPolynomial_ifft(tmp+1,poly2);
    LagrangeHalfCPolynomialMul(tmp+2,tmp+0,tmp+1);
    TorusPolynomial_fft(tmpr, tmp+2);
    torusPolynomialSubTo(result, tmpr);
    delete_TorusPolynomial(tmpr);
    delete_LagrangeHalfCPolynomial_array(3,tmp);
}

*/

func TorusPolynomialMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	t0 := IntPolynomial_ifft(N, poly1)
	t1 := TorusPolynomial_ifft(N, poly2)
	t2 := LagrangeHalfCPolynomialMul(t0.coefsC, t1.coefsC, int(N/2))
	result = TorusPolynomial_fft(N, t2)
}

func TorusPolynomialAddMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	t0 := IntPolynomial_ifft(N, poly1)
	t1 := TorusPolynomial_ifft(N, poly2)
	t2 := LagrangeHalfCPolynomialMul(t0.coefsC, t1.coefsC, int(N/2))
	tmpr2 := TorusPolynomial_fft(N, t2)
	torusPolynomialAddTo(result, tmpr2)
}

func TorusPolynomialSubMulR(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	t0 := IntPolynomial_ifft(N, poly1)
	t1 := TorusPolynomial_ifft(N, poly2)
	t2 := LagrangeHalfCPolynomialMul(t0.coefsC, t1.coefsC, int(N/2))
	tmpr2 := TorusPolynomial_fft(N, t2)
	torusPolynomialSubTo(result, tmpr2)
}
