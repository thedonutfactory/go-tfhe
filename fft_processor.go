package tfhe

import (
	"github.com/mjibson/go-dsp/fft"
	//"github.com/takatoh/fft"
)

var fftProcessor *DefaultFFTProcessor = NewDefaultFFTProcessor(1024)

type FFTProcessor interface {
	executeReverseTorus32(a []Torus32) (res []complex128)
	executeReverseInt(a []int32) (res []complex128)
	executeDirectTorus32(a []complex128) (res []Torus32)
}

type DefaultFFTProcessor struct {
	_2N int32
	N   int32
	Ns2 int32
}

type LagrangeHalfCPolynomial struct {
	coefsC []complex128
}

func NewDefaultFFTProcessor(N int32) *DefaultFFTProcessor {
	return &DefaultFFTProcessor{
		_2N: 2 * N,
		N:   N,
		Ns2: N / 2,
	}
}

func NewLagrangeHalfCPolynomial(n int32) *LagrangeHalfCPolynomial {
	//Assert(n == 1024)
	return &LagrangeHalfCPolynomial{coefsC: make([]complex128, n/2)}
}

func LagrangeHalfCPolynomialClear(p *LagrangeHalfCPolynomial) {
	p.coefsC = make([]complex128, len(p.coefsC))
}

func LagrangeHalfCPolynomialSetTorusConstant(result *LagrangeHalfCPolynomial, mu Torus32) {
	muc := complex(T32tod(mu), 0.)
	for j := 0; j < len(result.coefsC); j++ {
		result.coefsC[j] = muc
	}
}

func LagrangeHalfCPolynomialAddTorusConstant(result *LagrangeHalfCPolynomial, mu Torus32) {
	muc := complex(T32tod(mu), 0.)
	for j := 0; j < len(result.coefsC); j++ {
		result.coefsC[j] += muc
	}
}

/*
EXPORT void LagrangeHalfCPolynomialSetXaiMinusOne(LagrangeHalfCPolynomial* result, const int32_t ai) {
    LagrangeHalfCPolynomial_IMPL* result1 = (LagrangeHalfCPolynomial_IMPL*) result;
    const int32_t Ns2 = result1->proc->Ns2;
    const int32_t _2N = result1->proc->_2N;
    const cplx* omegaxminus1 = result1->proc->omegaxminus1;
    for (int32_t i=0; i<Ns2; i++)
	result1->coefsC[i]=omegaxminus1[((2*i+1)*ai)%_2N];
}
*/

/** termwise multiplication in Lagrange space */
func LagrangeHalfCPolynomialMul(result *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(result.coefsC); j++ {
		result.coefsC[j] = a.coefsC[j] * b.coefsC[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func LagrangeHalfCPolynomialAddMul(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.coefsC); j++ {
		accum.coefsC[j] += a.coefsC[j] * b.coefsC[j]
	}
}

/** termwise multiplication and addTo in Lagrange space */
func LagrangeHalfCPolynomialSubMul(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial, b *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.coefsC); j++ {
		accum.coefsC[j] += a.coefsC[j] * b.coefsC[j]
	}
}

func LagrangeHalfCPolynomialAddTo(accum *LagrangeHalfCPolynomial, a *LagrangeHalfCPolynomial) {
	for j := 0; j < len(accum.coefsC); j++ {
		accum.coefsC[j] += a.coefsC[j]
	}
}

func castComplex(arr []int32) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		res[i] = complex(float64(v), 0)
	}
	return
}

func castTorus(arr []complex128) (res []int32) {
	res = make([]int32, len(arr))
	for i, v := range arr {
		res[i] = int32(real(v))
	}
	return
}

func (p *DefaultFFTProcessor) executeReverseTorus32(a []Torus32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	//res = fft.IFFT(y)
	return
}

func (p *DefaultFFTProcessor) executeReverseInt(a []int32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func (p *DefaultFFTProcessor) executeDirectTorus32(a []complex128) (res []Torus32) {
	res = castTorus(fft.FFT(a))
	for i := 0; i < 512; i++ {
		res = append(res, 0)
	}
	return
}

/*
func (p *DefaultFFTProcessor) executeReverseInt(a []int32) (res []complex128) {
	N := int(p.N)
	Ns2 := int(p.Ns2)
	_2N := p._2N
	//rev_out_cplx := make([]complex128, N)
	rev_in := make([]complex128, _2N)

	//cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible

	for i := 0; i < N; i++ {
		rev_in[i] = complex(double(a[i])/2., 0.)
	}
	for i := 0; i < N; i++ {
		rev_in[N+i] = -rev_in[i]
	}
	//fftw_execute(rev_p);

	// FFT
	y := fft.FFT(rev_in)
	// IFFT
	rev_out_cplx := fft.IFFT(y)

	res = make([]complex128, N)
	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}
	for i := 0; i <= Ns2; i++ {
		r := math.Abs(real(rev_out_cplx[2*i]))
		if r >= 1e-20 {
			fmt.Printf("Error %f >= %f\n", r, 1e-20)
		}
		//Assert(r < 1e-20)
	}
	return
}

func (p *DefaultFFTProcessor) executeDirectTorus32(a []complex128) (res []Torus32) {
	N := int(p.N)
	Ns2 := int(p.Ns2)

	_2p32 := double(int64(1) << 32)
	_1sN := double(1) / double(N)
	//cplx* in_cplx = (cplx*) in; //fftw_complex and cplx are layout-compatible

	in_cplx := make([]complex128, N+1)

	for i := 0; i <= Ns2; i++ {
		in_cplx[2*i] = 0
	}
	for i := 0; i < Ns2; i++ {
		in_cplx[2*i+1] = a[i]
	}

	out := fft.IFFT(in_cplx)

	//fftw_execute(p)

	res = make([]Torus32, N)
	for i := 0; i < N; i++ {
		res[i] = Torus32(int64(real(out[i]) * _1sN * _2p32))
	}
	//pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
	for i := 0; i < N; i++ {
		r := math.Abs(real(out[N+i]) + real(out[i]))
		if r >= 1e-20 {
			fmt.Printf("Error %f >= %f\n", r, 1e-20)
		}
		//Assert(r < 1e-20)
	}
	return
}

func (p *DefaultFFTProcessor) executeReverseTorus32(a []Torus32) (res []complex128) {
	var _2pm33 double = 1. / double(int64(1)<<33)
	N := int(p.N)
	Ns2 := int(p.Ns2)
	_2N := p._2N
	rev_in := make([]complex128, _2N)
	//int32_t* aa = (int32_t*) a;
	//cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
	for i := 0; i < N; i++ {
		rev_in[i] = complex(double(a[i])*_2pm33, 0.)
	}
	for i := 0; i < N; i++ {
		rev_in[N+i] = -rev_in[i]
	}
	// FFT
	y := fft.FFT(rev_in)
	// IFFT
	rev_out_cplx := fft.IFFT(y)

	res = make([]complex128, N)
	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}
	for i := 0; i <= Ns2; i++ {
		r := cmplx.Abs(rev_out_cplx[2*i])
		//r := math.Abs(real(rev_out_cplx[2*i]))
		if r >= 1e-20 {
			fmt.Printf("Error %f >= %f\n", r, 1e-20)
		}
		//Assert(r < 1e-20)
	}
	return
}
*/

/*
func (p *DefaultFFTProcessor) check_alternate_real(x []complex128) {
	for i := 0; i < p._2N; i++ {
		assert(fabs(imag_inout[i])<1e-8);
	}
	for i := 0; i < p.N; i++ {
		assert(fabs(real_inout[i]+real_inout[N+i])<1e-9);
	}
}
*/

/*
func (p *DefaultFFTProcessor) executeReverseInt(a []int32) (res []complex128) {
	N := int(p.N)
	x := make([]complex128, 2*N)
	for i := 0; i < N; i++ {
		x[i] = complex(double(a[i]/2.), 0.)
	}
	for i := 0; i < N; i++ {
		x[N+i] = complex(-real(x[i]), 0.)
	}

	//check_alternate_real();
	//fft_transform_reverse(tables_reverse,real_inout,imag_inout);

	res = fft.IFFT(x)
	/*
		for i := 0; i < N; i += 2 {
			res[i] = y[i]
			//res_dbl[i]=real_inout[i+1];
			//res_dbl[i+1]=imag_inout[i+1];
		}

			for (int32_t i=0; i<Ns2; i++) {
				assert(abs(cplx(real_inout[2*i+1],imag_inout[2*i+1])-res[i])<1e-20);
			}
			check_conjugate_cplx();
*/
//	return
//}

/*
func (p *DefaultFFTProcessor) executeReverseTorus32(a []Torus32) (res []complex128) {
	var _2pm33 double = 1. / double(int64(1)<<33)
	//int32_t* aa = (int32_t*) a;

	N := int(p.N)
	x := make([]complex128, 2*N)
	for i := 0; i < N; i++ {
		x[i] = complex(double(a[i])*_2pm33, 0.)
	}
	for i := 0; i < N; i++ {
		x[N+i] = complex(-real(x[i]), 0.)
	}

	//check_alternate_real();
	//fft_transform_reverse(tables_reverse,real_inout,imag_inout);

	y := fft.IFFT(x)

	res = make([]complex128, p.Ns2)
	for i := int32(0); i < p.Ns2; i++ {
		res[i] = y[2*i+1] //complex128( real_inout[2*i+1], imag_inout[2*i+1])
	}
	//check_conjugate_cplx();

	return
}

func (p *DefaultFFTProcessor) executeDirectTorus32(a []complex128) (res []Torus32) {
	_2p32 := double(int64(1) << 32)
	_1sN := double(1) / double(p.N)
	//double* a_dbl=(double*) a;

	x := make([]complex128, 2*p.N)

	for i := int32(0); i < p.N; i++ {
		x[2*i] = complex(0., 0.)
	}

	for i := 0; i < int(p.Ns2); i++ {
		x[2*i+1] = a[i]
	}

	for i := 0; i < int(p.Ns2); i++ {
		x[int(p._2N)-1-2*i] = a[i]
	}

	//for (int32_t i=0; i<N; i++) real_inout[2*i]=0;
	//for (int32_t i=0; i<N; i++) imag_inout[2*i]=0;

	//for (int32_t i=0; i<Ns2; i++) real_inout[2*i+1]=real(a[i]);
	//for (int32_t i=0; i<Ns2; i++) imag_inout[2*i+1]=imag(a[i]);

	//for (int32_t i=0; i<Ns2; i++) real_inout[_2N-1-2*i]=real(a[i]);
	//for (int32_t i=0; i<Ns2; i++) imag_inout[_2N-1-2*i]=-imag(a[i]);

	//cmplx.Abs()

	//fft_transform(tables_direct,real_inout,imag_inout);
	y := fft.FFT(x)

	res = make([]Torus32, p.N)
	for i := int32(0); i < p.N; i++ {
		// res[i] = Torus32(int64(real(y[i]) * _1sN * _2p32))
		res[i] = Torus32(int64(math.Mod(cmplx.Abs(y[i])*_1sN, 1.) * _2p32))
		//res[i] = int32(cmplx.Abs(y[i]) * _1sN * _2p32)
	}
	//pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
	//check_alternate_real();
	return
}
*/
