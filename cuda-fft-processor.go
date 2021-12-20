package tfhe

import "github.com/mjibson/go-dsp/fft"

type CudaFftProcessor struct{}

func (*CudaFftProcessor) executeReverseInt(a []int32) []complex128 {
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
	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = z[2*i+1]
	}
	return res
}

func (*CudaFftProcessor) executeReverseTorus32(a []Torus32) []complex128 {
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
	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = z[2*i+1]
	}
	return res
}

func (*CudaFftProcessor) executeDirectTorus32(a []complex128) []Torus32 {
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

	res := make([]Torus32, N)
	res[0] = int32(int64(real(z[0]) * _1sN * _2p32))
	for i := 1; i < N; i++ {
		res[i] = -int32(int64(real(z[N-i]) * _1sN * _2p32))
	}
	return res
}
