package tfhe

import (
	"github.com/takatoh/fft"
)

var fftProcessor DefaultFFTProcessor = NewDefaultFFTProcessor(1024)

type FFTProcessor interface {
	//FFT_Processor_fftw(const int32_t N);
	executeReverseInt(res []complex128, a []int32)
	executeReverseTorus32(res []complex128, a []Torus32)
	executeDirectTorus32(res []Torus32, a []complex128)
}

type DefaultFFTProcessor struct {
	_2N     int32
	N       int32
	Ns2     int32
	rev_in  []double
	rev_out []complex128
	in      []complex128
	out     []double
	//omegaxminus1 []complex128
}

func NewDefaultFFTProcessor(N int32) DefaultFFTProcessor {

	return DefaultFFTProcessor{
		_2N:     2 * N,
		N:       N,
		Ns2:     N / 2,
		rev_in:  make([]double, 2*N),
		out:     make([]double, 2*N),
		rev_out: make([]complex128, N+1),
		in:      make([]complex128, N+1),
		//omegaxminus1: make([]complex128, 2*N),
	}
	/*
		for x := int32(0); x < 2*N; x++ {
			res.omegaxminus1[x] = complex128(math.Cos(x* int32(math.Pi/N)-1., -math.Sin(x*int32(math.Pi/N))))
			 // instead of cos(x*M_PI/N)-1. + sin(x*M_PI/N) * I
			//exp(i.x.pi/N)-1
		}
	*/
}

type LagrangeHalfCPolynomial struct {
	coefsC []complex128
	//cplx* coefsC;
	//FFT_Processor_fftw* proc;

}

func NewLangrangeHalfCPolynomial(n int32) *LagrangeHalfCPolynomial {
	//Assert(n == 1024)
	return &LagrangeHalfCPolynomial{coefsC: make([]complex128, n/2)}
}

func (p *DefaultFFTProcessor) executeReverseTorus32(res []complex128, a []Torus32) {
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
	y := fft.FFT(rev_in, N)
	// IFFT
	rev_out_cplx := fft.IFFT(y, N)

	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}
	for i := 0; i <= Ns2; i++ {
		//Assert(math.Abs(real(rev_out_cplx[2*i])) < 1e-20)
	}
}

func (p *DefaultFFTProcessor) executeReverseInt(res []complex128, a []int32) {
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
	y := fft.FFT(rev_in, N)
	// IFFT
	rev_out_cplx := fft.IFFT(y, N)

	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}
	for i := 0; i <= Ns2; i++ {
		//Assert(math.Abs(real(rev_out_cplx[2*i])) < 1e-20)
	}
}

func (p *DefaultFFTProcessor) executeDirectTorus32(res []Torus32, a []complex128) {
	N := int(p.N)
	Ns2 := int(p.Ns2)

	_2p32 := double(int64(1) << 32)
	_1sN := double(1) / double(N)
	//cplx* in_cplx = (cplx*) in; //fftw_complex and cplx are layout-compatible

	in_cplx := make([]complex128, N)

	for i := 0; i <= Ns2; i++ {
		in_cplx[2*i] = 0
	}
	for i := 0; i < Ns2; i++ {
		in_cplx[2*i+1] = a[i]
	}

	out := fft.IFFT(in_cplx, N)

	//fftw_execute(p)
	for i := 0; i < N; i++ {
		res[i] = Torus32(int64(real(out[i]) * _1sN * _2p32))
	}
	//pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
	for i := 0; i < N; i++ {
		//Assert(math.Abs(real(out[N+i])+real(out[i])) < 1e-20)
	}

}
