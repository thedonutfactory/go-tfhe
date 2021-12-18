package tfhe

type LweBootstrappingKeyFFT struct {
	InOutParams   *LweParams       ///< paramètre de l'input et de l'output. key: s
	BkParams      *TGswParams      ///< params of the Gsw elems in bk. key: s"
	AccumParams   *TLweParams      ///< params of the accum variable key: s"
	ExtractParams *LweParams       ///< params after extraction: key: s'
	Bk            []*TGswSampleFFT ///< the bootstrapping key (s.s")
	Ks            *LweKeySwitchKey ///< the keyswitch key (s'.s)
}

func NewLweBootstrappingKeyFFT(InOutParams *LweParams, BkParams *TGswParams, AccumParams *TLweParams,
	ExtractParams *LweParams, Bk []*TGswSampleFFT, Ks *LweKeySwitchKey) *LweBootstrappingKeyFFT {
	return &LweBootstrappingKeyFFT{
		InOutParams:   InOutParams,
		BkParams:      BkParams,
		AccumParams:   AccumParams,
		ExtractParams: ExtractParams,
		Bk:            Bk,
		Ks:            Ks,
	}
}

func InitLweBootstrappingKeyFFT(bk *LweBootstrappingKey) *LweBootstrappingKeyFFT {
	in_out_params := bk.InOutParams
	bk_params := bk.BkParams
	accum_params := bk_params.TlweParams
	extract_params := &accum_params.ExtractedLweparams
	n := in_out_params.N
	t := bk.Ks.T
	basebit := bk.Ks.Basebit
	base := bk.Ks.Base
	N := extract_params.N

	ks := NewLweKeySwitchKey(N, t, basebit, in_out_params)
	// Copy the KeySwitching key
	for i := int32(0); i < N; i++ {
		for j := int32(0); j < t; j++ {
			for p := int32(0); p < base; p++ {
				LweCopy(ks.Ks[i][j][p], bk.Ks.Ks[i][j][p], in_out_params)
			}
		}
	}

	// Bootstrapping Key FFT
	bkFFT := NewTGswSampleFFTArray(n, bk_params)
	for i := int32(0); i < n; i++ {
		tGswToFFTConvert(bkFFT[i], bk.Bk[i], bk_params)
	}

	return NewLweBootstrappingKeyFFT(in_out_params, bk_params, accum_params, extract_params, bkFFT, ks)
}

func tfhe_MuxRotate_FFT(result, accum *TLweSample, bki *TGswSampleFFT, barai int32, bk_params *TGswParams) {
	// ACC = BKi*[(X^barai-1)*ACC]+ACC
	// temp = (X^barai-1)*ACC
	TLweMulByXaiMinusOne(result, barai, accum, bk_params.TlweParams)
	// temp *= BKi
	tGswFFTExternMulToTLwe(result, bki, bk_params)
	// ACC += temp
	TLweAddTo(result, accum, bk_params.TlweParams)
}

/**
 * multiply the accumulator by X^sum(bara_i.s_i)
 * @param accum the TLWE sample to multiply
 * @param bk An array of n TGSW FFT samples where bk_i encodes s_i
 * @param bara An array of n coefficients between 0 and 2N-1
 * @param bk_params The parameters of bk
 */
func tfhe_blindRotate_FFT(accum *TLweSample,
	bk []*TGswSampleFFT,
	bara []int32,
	n int32,
	bk_params *TGswParams) {

	//TGswSampleFFT* temp = new_TGswSampleFFT(bk_params);
	temp := NewTLweSample(bk_params.TlweParams)
	temp2 := temp
	temp3 := accum

	for i := int32(0); i < n; i++ {
		barai := bara[i]
		if barai == 0 {
			continue
		}
		tfhe_MuxRotate_FFT(temp2, temp3, bk[i], barai, bk_params)
		swap(temp2, temp3)
	}
	if temp3 != accum {
		TLweCopy(accum, temp3, bk_params.TlweParams)
	}
}

func swap(px, py *TLweSample) {
	tempx := *px
	tempy := *py
	*px = tempy
	*py = tempx
}

/**
 * result = LWE(v_p) where p=barb-sum(bara_i.s_i) mod 2N
 * @param result the output LWE sample
 * @param v a 2N-elt anticyclic function (represented by a TorusPolynomial)
 * @param bk An array of n TGSW FFT samples where bk_i encodes s_i
 * @param barb A coefficients between 0 and 2N-1
 * @param bara An array of n coefficients between 0 and 2N-1
 * @param bk_params The parameters of bk
 */
func tfhe_blindRotateAndExtract_FFT(result *LweSample, v *TorusPolynomial,
	bk []*TGswSampleFFT,
	barb int32,
	bara []int32,
	n int32,
	bk_params *TGswParams) {

	accum_params := bk_params.TlweParams
	extract_params := &accum_params.ExtractedLweparams
	N := accum_params.N
	_2N := 2 * N

	// Test polynomial
	testvectbis := NewTorusPolynomial(N)
	// Accumulator
	acc := NewTLweSample(accum_params)

	// testvector = X^{2N-barb}*v
	if barb != 0 {
		TorusPolynomialMulByXai(testvectbis, _2N-barb, v)
	} else {
		TorusPolynomialCopy(testvectbis, v)
	}
	TLweNoiselessTrivial(acc, testvectbis, accum_params)
	// Blind rotation
	tfhe_blindRotate_FFT(acc, bk, bara, n, bk_params)
	// Extraction
	tLweExtractLweSample(result, acc, extract_params, accum_params)
}

/**
 * result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
 * @param result The resulting LweSample
 * @param bk The bootstrapping + keyswitch key
 * @param mu The output message (if phase(x)>0)
 * @param x The input sample
 */
func tfheBootstrapWoKSFFT(result *LweSample, bk *LweBootstrappingKeyFFT, mu Torus32, x *LweSample) {

	bk_params := bk.BkParams
	accum_params := bk.AccumParams
	in_params := bk.InOutParams
	N := accum_params.N
	Nx2 := 2 * N
	n := in_params.N

	testvect := NewTorusPolynomial(N)
	bara := make([]int32, N)

	// Modulus switching
	barb := ModSwitchFromTorus32(x.B, Nx2)
	for i := int32(0); i < n; i++ {
		bara[i] = ModSwitchFromTorus32(x.A[i], Nx2)
	}

	// the initial testvec = [mu,mu,mu,...,mu]
	for i := int32(0); i < N; i++ {
		testvect.Coefs[i] = mu
	}

	// Bootstrapping rotation and extraction
	tfhe_blindRotateAndExtract_FFT(result, testvect, bk.Bk, barb, bara, n, bk_params)
}

/**
 * result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
 * @param result The resulting LweSample
 * @param bk The bootstrapping + keyswitch key
 * @param mu The output message (if phase(x)>0)
 * @param x The input sample
 */
func tfheBootstrapFFT(bk *LweBootstrappingKeyFFT, mu Torus32, x *LweSample) *LweSample {

	u := NewLweSample(&bk.AccumParams.ExtractedLweparams)

	tfheBootstrapWoKSFFT(u, bk, mu, x)
	// Key switching
	return lweKeySwitch(bk.Ks, u)
}
