package tfhe

func tfhe_MuxRotate(result *TLweSample, accum *TLweSample, bki *TGswSample, barai int32,
	bk_params *TGswParams) {
	// ACC = BKi*[(X^barai-1)*ACC]+ACC
	// temp = (X^barai-1)*ACC
	TLweMulByXaiMinusOne(result, barai, accum, bk_params.TlweParams)
	// temp *= BKi
	TGswExternMulToTLwe(result, bki, bk_params)
	// ACC += temp
	TLweAddTo(result, accum, bk_params.TlweParams)
}

/**
* multiply the accumulator by X^sum(bara_i.s_i)
* @param accum the TLWE sample to multiply
* @param bk An array of n TGSW samples where bk_i encodes s_i
* @param bara An array of n coefficients between 0 and 2N-1
* @param bk_params The parameters of bk
 */
func tfhe_blindRotate(accum *TLweSample, bk []*TGswSample, bara []int32, n int32, bk_params *TGswParams) {
	temp := NewTLweSample(bk_params.TlweParams)
	temp2 := temp
	temp3 := accum
	for i := int32(0); i < n; i++ {
		barai := bara[i]
		if barai == 0 {
			continue //indeed, this is an easy case!
		}
		tfhe_MuxRotate(temp2, temp3, bk[i], barai, bk_params)
		//swap(temp2, temp3)
		temp3, temp2 = temp2, temp3

	}
	if temp3 != accum {
		TLweCopy(accum, temp3, bk_params.TlweParams)
	}
}

/**
* result = LWE(v_p) where p=barb-sum(bara_i.s_i) mod 2N
* @param result the output LWE sample
* @param v a 2N-elt anticyclic function (represented by a TorusPolynomial)
* @param bk An array of n TGSW samples where bk_i encodes s_i
* @param barb A coefficients between 0 and 2N-1
* @param bara An array of n coefficients between 0 and 2N-1
* @param bk_params The parameters of bk
 */
func tfhe_blindRotateAndExtract(result *LweSample,
	v *TorusPolynomial,
	bk []*TGswSample,
	barb int32,
	bara []int32,
	n int32,
	bk_params *TGswParams) {

	accum_params := bk_params.TlweParams
	extract_params := &accum_params.extractedLweparams
	N := accum_params.N
	_2N := 2 * N

	testvectbis := NewTorusPolynomial(N)
	acc := NewTLweSample(accum_params)

	if barb != 0 {
		TorusPolynomialMulByXai(testvectbis, _2N-barb, v)
	} else {
		TorusPolynomialCopy(testvectbis, v)
	}
	TLweNoiselessTrivial(acc, testvectbis, accum_params)
	tfhe_blindRotate(acc, bk, bara, n, bk_params)
	tLweExtractLweSample(result, acc, extract_params, accum_params)
}

/**
* result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
* @param result The resulting LweSample
* @param bk The bootstrapping + keyswitch key
* @param mu The output message (if phase(x)>0)
* @param x The input sample
 */
func tfheBootstrapWoKS(result *LweSample, bk *LweBootstrappingKey, mu Torus32, x *LweSample) {
	bk_params := bk.bk_params
	accum_params := bk.accum_params
	in_params := bk.in_out_params
	N := accum_params.N
	Nx2 := 2 * N
	n := in_params.N
	testvect := NewTorusPolynomial(N)
	bara := make([]int32, N) // new int32_t[N];
	barb := ModSwitchFromTorus32(x.B, Nx2)
	for i := int32(0); i < n; i++ {
		bara[i] = ModSwitchFromTorus32(x.A[i], Nx2)
	}
	//the initial testvec = [mu,mu,mu,...,mu]
	for i := int32(0); i < N; i++ {
		testvect.CoefsT[i] = mu
	}
	tfhe_blindRotateAndExtract(result, testvect, bk.bk, barb, bara, n, bk_params)
}

/**
* result = LWE(mu) iff phase(x)>0, LWE(-mu) iff phase(x)<0
* @param result The resulting LweSample
* @param bk The bootstrapping + keyswitch key
* @param mu The output message (if phase(x)>0)
* @param x The input sample
 */
func tfheBootstrap(result *LweSample, bk *LweBootstrappingKey, mu Torus32, x *LweSample) {
	u := NewLweSample(&bk.accum_params.extractedLweparams)
	tfheBootstrapWoKS(u, bk, mu, x)
	// Key Switching
	lweKeySwitch(result, bk.ks, u)
}

func tfhe_createLweBootstrappingKey(bk *LweBootstrappingKey, key_in *LweKey, rgsw_key *TGswKey) {
	Assert(bk.bk_params == rgsw_key.params)
	Assert(bk.in_out_params == key_in.params)

	in_out_params := bk.in_out_params
	bk_params := bk.bk_params
	accum_params := bk_params.TlweParams
	extract_params := &accum_params.extractedLweparams

	//LweKeySwitchKey* ks; ///< the keyswitch key (s'.s)
	accum_key := &rgsw_key.TlweKey
	extracted_key := NewLweKey(extract_params)
	tLweExtractKey(extracted_key, accum_key)
	lweCreateKeySwitchKey(bk.ks, extracted_key, key_in)

	//TGswSample* bk; ///< the bootstrapping key (s.s")
	kin := key_in.key
	alpha := accum_params.alphaMin
	n := in_out_params.N
	//const int32_t kpl = bk_params.kpl;
	//const int32_t k = accum_params.k;
	//const int32_t N = accum_params.N;
	//cout << "create the bootstrapping key bk ("  << "  " << n*kpl*(k+1)*N*4 << " bytes)" << endl;
	//cout << "  with noise_stdev: " << alpha << endl;
	for i := int32(0); i < n; i++ {
		TGswSymEncryptInt(bk.bk[i], kin[i], alpha, rgsw_key)
	}

}
