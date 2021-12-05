package tfhe

func tfheMuxRotate(result *TLweSample, accum *TLweSample, bki *TGswSample, barai int,
	bkParams *TGswParams) {
	// ACC = BKi*[(X^barai-1)*ACC]+ACC
	// temp = (X^barai-1)*ACC
	TLweMulByXaiMinusOne(result, barai, accum, bkParams.TlweParams)
	// temp *= BKi
	TGswExternMulToTLwe(result, bki, bkParams)
	// ACC += temp
	TLweAddTo(result, accum, bkParams.TlweParams)
}

/**
* multiply the accumulator by X^sum(bara_i.s_i)
* @param accum the TLWE sample to multiply
* @param bk An array of n TGSW samples where bk_i encodes s_i
* @param bara An array of n coefficients between 0 and 2N-1
* @param bkParams The parameters of bk
 */
func tfheBlindRotate(accum *TLweSample, bk []*TGswSample, bara []int, n int, bkParams *TGswParams) {
	temp := NewTLweSample(bkParams.TlweParams)
	temp2 := temp
	temp3 := accum
	for i := int(0); i < n; i++ {
		barai := bara[i]
		if barai == 0 {
			continue //indeed, this is an easy case!
		}
		tfheMuxRotate(temp2, temp3, bk[i], barai, bkParams)
		//swap(temp2, temp3)
		temp3, temp2 = temp2, temp3

	}
	if temp3 != accum {
		TLweCopy(accum, temp3, bkParams.TlweParams)
	}
}

/**
* result = LWE(v_p) where p=barb-sum(bara_i.s_i) mod 2N
* @param result the output LWE sample
* @param v a 2N-elt anticyclic function (represented by a TorusPolynomial)
* @param bk An array of n TGSW samples where bk_i encodes s_i
* @param barb A coefficients between 0 and 2N-1
* @param bara An array of n coefficients between 0 and 2N-1
* @param bkParams The parameters of bk
 */
func tfheBlindRotateAndExtract(result *LweSample,
	v *TorusPolynomial,
	bk []*TGswSample,
	barb int,
	bara []int,
	n int,
	bkParams *TGswParams) {

	accumParams := bkParams.TlweParams
	extractParams := &accumParams.ExtractedLweparams
	N := accumParams.N
	_2N := 2 * N

	testvectbis := NewTorusPolynomial(N)
	acc := NewTLweSample(accumParams)

	if barb != 0 {
		TorusPolynomialMulByXai(testvectbis, _2N-barb, v)
	} else {
		TorusPolynomialCopy(testvectbis, v)
	}
	TLweNoiselessTrivial(acc, testvectbis, accumParams)
	tfheBlindRotate(acc, bk, bara, n, bkParams)
	tLweExtractLweSample(result, acc, extractParams, accumParams)
}

/**
* result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
* @param result The resulting LweSample
* @param bk The bootstrapping + keyswitch key
* @param mu The output message (if phase(x)>0)
* @param x The input sample
 */
func tfheBootstrapWoKS(result *LweSample, bk *LweBootstrappingKey, mu Torus32, x *LweSample) {
	bkParams := bk.BkParams
	accumParams := bk.AccumParams
	inParams := bk.InOutParams
	N := accumParams.N
	Nx2 := 2 * N
	n := inParams.N
	testvect := NewTorusPolynomial(N)
	bara := make([]int, N)
	barb := ModSwitchFromTorus(x.B, Nx2)
	for i := int(0); i < n; i++ {
		bara[i] = ModSwitchFromTorus(x.A[i], Nx2)
	}
	//the initial testvec = [mu,mu,mu,...,mu]
	for i := int32(0); i < N; i++ {
		testvect.Coefs[i] = mu
	}
	tfheBlindRotateAndExtract(result, testvect, bk.Bk, barb, bara, n, bkParams)
}

/**
* result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
* @param result The resulting LweSample
* @param bk The bootstrapping + keyswitch key
* @param mu The output message (if phase(x)>0)
* @param x The input sample
 */
func tfheBootstrap(bk *LweBootstrappingKey, mu Torus32, x *LweSample) *LweSample {
	u := NewLweSample(&bk.AccumParams.ExtractedLweparams)
	tfheBootstrapWoKS(u, bk, mu, x)
	// Key Switching
	return lweKeySwitch(bk.Ks, u)
}

func tfheCreateLweBootstrappingKey(bk *LweBootstrappingKey, keyIn *LweKey, rgswKey *TGswKey) {
	Assert(bk.BkParams == rgswKey.Params)
	Assert(bk.InOutParams == keyIn.Params)

	inOutParams := bk.InOutParams
	bkParams := bk.BkParams
	accumParams := bkParams.TlweParams
	extractParams := &accumParams.ExtractedLweparams

	//LweKeySwitchKey* ks; ///< the keyswitch key (s'.s)
	accumKey := &rgswKey.TlweKey
	extractedKey := NewLweKey(extractParams)
	tLweExtractKey(extractedKey, accumKey)
	lweCreateKeySwitchKey(bk.Ks, extractedKey, keyIn)

	//TGswSample* bk; ///< the bootstrapping key (s.s")
	kin := keyIn.Key
	alpha := accumParams.AlphaMin
	n := inOutParams.N

	for i := int32(0); i < n; i++ {
		TGswSymEncryptInt(bk.Bk[i], kin[i], alpha, rgswKey)
	}
}
