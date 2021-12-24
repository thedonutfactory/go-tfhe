package fft

import (
	fftw "github.com/cpmech/gosl/fun/fftw"
	. "github.com/thedonutfactory/go-tfhe/types"
)

type FftwProcessor struct {
	revIntInout []complex128
	revIntPlan  *fftw.Plan1d

	revTorusInout []complex128
	revTorusPlan  *fftw.Plan1d

	dirTorusInout []complex128
	dirTorusPlan  *fftw.Plan1d

	N int
}

func NewFftwProcessor(N int) *FftwProcessor {

	revIntInout := make([]complex128, 2*N)
	revTorusInout := make([]complex128, 2*N)
	dirTorusInout := make([]complex128, 2*N)

	return &FftwProcessor{
		N:             N,
		revIntInout:   revIntInout,
		revIntPlan:    fftw.NewPlan1d(revIntInout, true, false),
		revTorusInout: revTorusInout,
		revTorusPlan:  fftw.NewPlan1d(revTorusInout, true, false),
		dirTorusInout: dirTorusInout,
		dirTorusPlan:  fftw.NewPlan1d(dirTorusInout, false, false),
	}
}

func (p *FftwProcessor) executeReverseInt(a []int32) []complex128 {
	//N := len(a)
	Ns2 := p.N / 2
	//_2N := N * 2
	//cplxInout := make([]complex128, _2N)
	for i := 0; i < p.N; i++ {
		p.revIntInout[i] = complex(float64(a[i])/2., 0.)
	}
	for i := 0; i < p.N; i++ {
		p.revIntInout[p.N+i] = complex(-(float64(a[i]) / 2), 0.)
	}
	//plan := fftw.NewPlan1d(p.revIntInout, true, false)
	//p.revIntPlan = fftw.NewPlan1d(p.revIntInout, true, false)

	p.revIntPlan.Execute()
	//defer p.revIntPlan.Free()

	//z := fft.FFT(cplxInout)
	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = p.revIntInout[2*i+1]
	}
	return res
}

func (p *FftwProcessor) executeReverseTorus32(a []Torus32) []complex128 {
	//N := len(a)
	Ns2 := p.N / 2
	//ß_2N := N * 2
	//cplxInout := make([]complex128, _2N)
	for i := 0; i < p.N; i++ {
		t := float64(a[i]) / float64(int64(1)<<33)
		p.revTorusInout[i] = complex(t, 0.)
	}
	for i := 0; i < p.N; i++ {
		t := float64(a[i]) / float64(int64(1)<<33)
		p.revTorusInout[p.N+i] = complex(-t, 0.)
	}
	//z := fft.FFT(cplxInout)
	//plan := fftw.NewPlan1d(p.revTorusInout, true, false)
	//p.revTorusPlan = fftw.NewPlan1d(p.revIntInout, true, false)
	p.revTorusPlan.Execute()
	//defer p.revTorusPlan.Free()

	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = p.revTorusInout[2*i+1]
	}
	return res
}

func (p *FftwProcessor) executeDirectTorus32(a []complex128) []Torus32 {
	//N := len(a) * 2
	Ns2 := p.N / 2
	_2N := p.N * 2
	_2p32 := float64(int64(1) << 32)
	_1sN := float64(1) / float64(p.N)

	//cplxInout := make([]complex128, N*2)
	for i := 0; i < p.N; i++ {
		p.dirTorusInout[2*i] = complex(0., 0.)
	}
	for i := 0; i < Ns2; i++ {
		p.dirTorusInout[2*i+1] = complex(real(a[i]), imag(a[i]))
	}
	for i := 0; i < Ns2; i++ {
		p.dirTorusInout[_2N-1-2*i] = complex(real(a[i]), -imag(a[i]))
	}
	//z := fft.FFT(cplxInout)
	//plan := fftw.NewPlan1d(cplxInout, false, false)
	//p.dirTorusPlan = fftw.NewPlan1d(p.revIntInout, false, false)
	p.dirTorusPlan.Execute()
	//defer p.dirTorusPlan.Free()

	res := make([]Torus32, p.N)

	for i := 0; i < p.N; i++ {
		res[i] = Torus32(int64(real(p.dirTorusInout[i]) * _1sN * _2p32))
	}

	/*
		res[0] = int32(int64(real(z[0]) * _1sN * _2p32))
		for i := 1; i < N; i++ {
			res[i] = -int32(int64(real(z[N-i]) * _1sN * _2p32))
		}
	*/
	return res
}
