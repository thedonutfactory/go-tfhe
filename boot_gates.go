package tfhe

/*
 * Homomorphic bootstrapped NAND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Nand(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) - ca - cb
	NandConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, NandConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//return tfheBootstrap(bk.Bkw.Bk, MU, tempResult)
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
	// tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Or(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) + ca + cb
	OrConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func And(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca + cb
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Xor(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/4) + 2*(ca + cb)
	XorConst := ModSwitchToTorus32(1, 4)
	LweNoiselessTrivial(tempResult, XorConst, inOutParams)
	LweAddMulTo(tempResult, 2, ca, inOutParams)
	LweAddMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XNOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Xnor(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/4) + 2*(-ca-cb)
	XnorConst := ModSwitchToTorus32(-1, 4)
	LweNoiselessTrivial(tempResult, XnorConst, inOutParams)
	LweSubMulTo(tempResult, 2, ca, inOutParams)
	LweSubMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped NOT gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func Not(ca *LweSample, bk *PublicKey) *LweSample {
	inOutParams := bk.Params.InOutParams
	result := NewLweSample(inOutParams)
	LweNegate(result, ca, inOutParams)
	return result
}

/*
 * Homomorphic bootstrapped COPY gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func Copy(ca *LweSample, bk *PublicKey) *LweSample {
	inOutParams := bk.Params.InOutParams
	result := NewLweSample(inOutParams)
	LweCopy(result, ca, inOutParams)
	return result
}

/*
 * Homomorphic Trivial Constant gate (doesn't need to be bootstrapped)
 * Takes a boolean value)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func Constant(value bool, bk *PublicKey) *LweSample {
	inOutParams := bk.Params.InOutParams
	MU := ModSwitchToTorus32(1, 8)
	var muValue = -MU
	if value == true {
		muValue = MU
	}
	result := NewLweSample(inOutParams)
	LweNoiselessTrivial(result, muValue, inOutParams)
	return result
}

/*
 * Homomorphic bootstrapped NOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Nor(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca - cb
	NorConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, NorConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndNY Gate: not(a) and b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func AndNY(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca + cb
	AndNYConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndNYConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndYN Gate: a and not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func AndYN(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca - cb
	AndYNConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndYNConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrNY Gate: not(a) or b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func OrNY(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) - ca + cb
	OrNYConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrNYConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrYN Gate: a or not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func OrYN(ca, cb *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) + ca - cb
	OrYNConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrYNConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	return tfheBootstrapFFT(bk.Bkw.BkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped Mux(a,b,c) = a?b:c = a*b + not(a)*c
 * Takes in input 3 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func Mux(a, b, c *LweSample, bk *PublicKey) *LweSample {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.Params.InOutParams
	extractedParams := &bk.Params.TgswParams.TlweParams.ExtractedLweparams

	tempResult := NewLweSample(inOutParams)
	tempResult1 := NewLweSample(extractedParams)
	u1 := NewLweSample(extractedParams)
	u2 := NewLweSample(extractedParams)

	//compute "AND(a,b)": (0,-1/8) + a + b
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweAddTo(tempResult, a, inOutParams)
	LweAddTo(tempResult, b, inOutParams)
	// Bootstrap without KeySwitch
	// tfheBootstrapWoKS_FFT(u1, bk.bkFFT, MU, tempResult);
	tfheBootstrapWoKSFFT(u1, bk.Bkw.BkFFT, MU, tempResult)

	//compute "AND(not(a),c)": (0,-1/8) - a + c
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweSubTo(tempResult, a, inOutParams)
	LweAddTo(tempResult, c, inOutParams)
	// Bootstrap without KeySwitch
	//tfheBootstrapWoKS_FFT(u2, bk.bkFFT, MU, tempResult);
	tfheBootstrapWoKSFFT(u2, bk.Bkw.BkFFT, MU, tempResult)

	// Add u1=u1+u2
	MuxConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult1, MuxConst, extractedParams)
	LweAddTo(tempResult1, u1, extractedParams)
	LweAddTo(tempResult1, u2, extractedParams)
	// Key switching
	//lweKeySwitch(result, bk.bkFFT.ks, tempResult1)
	return lweKeySwitch(bk.Bkw.BkFFT.Ks, tempResult1)
}

/*
 * Homomorphic bootstrapped OR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,1/8) + ca + cb
	OrConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(temp_result, OrConst, in_out_params)
	LweAddTo(temp_result, ca, in_out_params)
	LweAddTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped AND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsAND(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,-1/8) + ca + cb
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(temp_result, AndConst, in_out_params)
	LweAddTo(temp_result, ca, in_out_params)
	LweAddTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped XOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsXOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,1/4) + 2*(ca + cb)
	XorConst := ModSwitchToTorus32(1, 4)
	LweNoiselessTrivial(temp_result, XorConst, in_out_params)
	LweAddMulTo(temp_result, 2, ca, in_out_params)
	LweAddMulTo(temp_result, 2, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped XNOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsXNOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,-1/4) + 2*(-ca-cb)
	XnorConst := ModSwitchToTorus32(-1, 4)
	LweNoiselessTrivial(temp_result, XnorConst, in_out_params)
	LweSubMulTo(temp_result, 2, ca, in_out_params)
	LweSubMulTo(temp_result, 2, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped NOT gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNOT(result, ca *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	in_out_params := bk.params.InOutParams
	LweNegate(result, ca, in_out_params)
}

/*
 * Homomorphic bootstrapped COPY gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsCOPY(result, ca *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	in_out_params := bk.params.InOutParams
	LweCopy(result, ca, in_out_params)
}

/*
 * Homomorphic Trivial Constant gate (doesn't need to be bootstrapped)
 * Takes a boolean value)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsCONSTANT(result *LweSample, value int32, bk *TFheGateBootstrappingCloudKeySet) {
	in_out_params := bk.params.InOutParams
	MU := ModSwitchToTorus32(1, 8)
	var muValue = -MU
	if value != 0 {
		muValue = MU
	}
	LweNoiselessTrivial(result, muValue, in_out_params)
}

/*
 * Homomorphic bootstrapped NOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,-1/8) - ca - cb
	NorConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(temp_result, NorConst, in_out_params)
	LweSubTo(temp_result, ca, in_out_params)
	LweSubTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped AndNY Gate: not(a) and b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsANDNY(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,-1/8) - ca + cb
	AndNYConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(temp_result, AndNYConst, in_out_params)
	LweSubTo(temp_result, ca, in_out_params)
	LweAddTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped AndYN Gate: a and not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsANDYN(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,-1/8) + ca - cb
	AndYNConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(temp_result, AndYNConst, in_out_params)
	LweAddTo(temp_result, ca, in_out_params)
	LweSubTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped OrNY Gate: not(a) or b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsORNY(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,1/8) - ca + cb
	OrNYConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(temp_result, OrNYConst, in_out_params)
	LweSubTo(temp_result, ca, in_out_params)
	LweAddTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped OrYN Gate: a or not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsORYN(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,1/8) + ca - cb
	OrYNConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(temp_result, OrYNConst, in_out_params)
	LweAddTo(temp_result, ca, in_out_params)
	LweSubTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
}

/*
 * Homomorphic bootstrapped Mux(a,b,c) = a?b:c = a*b + not(a)*c
 * Takes in input 3 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsMUX(result, a, b, c *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	extracted_params := &bk.params.tgswParams.TlweParams.extractedLweparams

	temp_result := NewLweSample(in_out_params)
	temp_result1 := NewLweSample(extracted_params)
	u1 := NewLweSample(extracted_params)
	u2 := NewLweSample(extracted_params)

	//compute "AND(a,b)": (0,-1/8) + a + b
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(temp_result, AndConst, in_out_params)
	LweAddTo(temp_result, a, in_out_params)
	LweAddTo(temp_result, b, in_out_params)
	// Bootstrap without KeySwitch
	// tfhe_bootstrap_woKS_FFT(u1, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap_woKS(u1, bk.bk, MU, temp_result)

	//compute "AND(not(a),c)": (0,-1/8) - a + c
	LweNoiselessTrivial(temp_result, AndConst, in_out_params)
	LweSubTo(temp_result, a, in_out_params)
	LweAddTo(temp_result, c, in_out_params)
	// Bootstrap without KeySwitch
	//tfhe_bootstrap_woKS_FFT(u2, bk.bkFFT, MU, temp_result);
	tfhe_bootstrap_woKS(u2, bk.bk, MU, temp_result)

	// Add u1=u1+u2
	MuxConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(temp_result1, MuxConst, extracted_params)
	LweAddTo(temp_result1, u1, extracted_params)
	LweAddTo(temp_result1, u2, extracted_params)
	// Key switching
	//lweKeySwitch(result, bk.bkFFT.ks, temp_result1)
	lweKeySwitch(result, bk.bk.ks, temp_result1)
}
