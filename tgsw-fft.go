package tfhe

type TGswSampleFFT struct {
	AllSample  []*TLweSampleFFT   ///< TLweSample* allSample; (k+1)l TLwe Sample
	BlocSample [][]*TLweSampleFFT ///< optional access to the different size l blocks
	// double currentVariance
	K int32
	L int32
}

func InitNewTGswSampleFFT(params *TGswParams, allSamples []*TLweSampleFFT) *TGswSampleFFT {
	var c int = 0
	blocSamples := make([][]*TLweSampleFFT, params.TlweParams.K+1)
	for i := range blocSamples {
		blocSamples[i] = make([]*TLweSampleFFT, params.L)
		for j := range blocSamples[i] {
			blocSamples[i][j] = allSamples[c]
			c++
		}
	}
	return &TGswSampleFFT{
		AllSample:  allSamples,
		BlocSample: blocSamples,
		K:          params.TlweParams.K,
		L:          params.L,
	}
}

func NewTGswSampleFFT(params *TGswParams) *TGswSampleFFT {
	k := params.TlweParams.K
	l := params.L
	kpl := params.Kpl
	allSamples := make([]*TLweSampleFFT, kpl)
	for i := range allSamples {
		allSamples[i] = NewTLweSampleFFT(params.TlweParams)
	}
	var c int = 0
	blocSamples := make([][]*TLweSampleFFT, k+1)
	for i := range blocSamples {
		blocSamples[i] = make([]*TLweSampleFFT, l)
		for j := range blocSamples[i] {
			blocSamples[i][j] = allSamples[c]
			c++
		}
	}
	return &TGswSampleFFT{
		AllSample:  allSamples,
		BlocSample: blocSamples,
		K:          k,
		L:          l,
	}
}

func NewTGswSampleFFTArray(size int32, params *TGswParams) (arr []*TGswSampleFFT) {
	arr = make([]*TGswSampleFFT, size)
	for i := int32(0); i < size; i++ {
		arr[i] = NewTGswSampleFFT(params)
	}
	return
}

func init_TGswSampleFFT(obj *TGswSampleFFT, params *TGswParams) {
	k := params.TlweParams.K
	l := params.L
	all_samples := NewTLweSampleFFTArray((k+1)*l, params.TlweParams)
	obj = InitNewTGswSampleFFT(params, all_samples)
}

// For all the kpl TLWE samples composing the TGSW sample
// It computes the inverse FFT of the coefficients of the TLWE sample
func tGswToFFTConvert(result *TGswSampleFFT, source *TGswSample, params *TGswParams) {
	kpl := params.Kpl

	for p := int32(0); p < kpl; p++ {
		tLweToFFTConvert(result.AllSample[p], &source.AllSample[p], params.TlweParams)
	}
}

// For all the kpl TLWE samples composing the TGSW sample
// It computes the FFT of the coefficients of the TLWEfft sample
func tGswFromFFTConvert(result *TGswSample, source *TGswSampleFFT, params *TGswParams) {
	kpl := params.Kpl
	for p := int32(0); p < kpl; p++ {
		tLweFromFFTConvert(&result.AllSample[p], source.AllSample[p], params.TlweParams)
	}
}

// result = result + H
func tGswFFTAddH(result *TGswSampleFFT, params *TGswParams) {
	k := params.TlweParams.K
	l := params.L

	for j := int32(0); j < l; j++ {
		hj := params.H[j]
		for i := int32(0); i <= k; i++ {
			LagrangeHalfCPolynomialAddTorusConstant(result.BlocSample[i][j].A[i], hj)
		}
	}
}

// result = list of TLWE (0,0)
func tGswFFTClear(result *TGswSampleFFT, params *TGswParams) {
	kpl := params.Kpl

	for p := int32(0); p < kpl; p++ {
		tLweFFTClear(result.AllSample[p], params.TlweParams)
	}
}

// External product (*): accum = gsw (*) accum
func tGswFFTExternMulToTLwe(accum *TLweSample, gsw *TGswSampleFFT, params *TGswParams) {
	tlwe_params := params.TlweParams
	k := tlwe_params.K
	l := params.L
	kpl := params.Kpl
	N := tlwe_params.N
	//TODO attention, improve these new/delete...
	deca := NewIntPolynomialArray(int(kpl), N)         //decomposed accumulator
	decaFFT := NewLagrangeHalfCPolynomialArray(kpl, N) //fft version
	tmpa := NewTLweSampleFFT(tlwe_params)

	for i := int32(0); i <= k; i++ {
		TGswTorus32PolynomialDecompH(deca[i*l:], &accum.A[i], params)
	}
	for p := int32(0); p < kpl; p++ {
		fftProc.intPolynomialIfft(decaFFT[p], &deca[p])
	}

	tLweFFTClear(tmpa, tlwe_params)
	for p := int32(0); p < kpl; p++ {
		tLweFFTAddMulRTo(tmpa, decaFFT[p], gsw.AllSample[p], tlwe_params)
	}
	tLweFromFFTConvert(accum, tmpa, tlwe_params)
}
