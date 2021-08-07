package tfhe

/*
 * Homomorphic bootstrapped NAND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNAND(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) - ca - cb
	NandConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, NandConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	tfheBootstrap(result, bk.bk, MU, tempResult)
	// tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) + ca + cb
	OrConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsAND(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca + cb
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsXOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/4) + 2*(ca + cb)
	XorConst := ModSwitchToTorus32(1, 4)
	LweNoiselessTrivial(tempResult, XorConst, inOutParams)
	LweAddMulTo(tempResult, 2, ca, inOutParams)
	LweAddMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped XNOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsXNOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/4) + 2*(-ca-cb)
	XnorConst := ModSwitchToTorus32(-1, 4)
	LweNoiselessTrivial(tempResult, XnorConst, inOutParams)
	LweSubMulTo(tempResult, 2, ca, inOutParams)
	LweSubMulTo(tempResult, 2, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped NOT gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNOT(result, ca *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	inOutParams := bk.params.InOutParams
	LweNegate(result, ca, inOutParams)
}

/*
 * Homomorphic bootstrapped COPY gate (doesn't need to be bootstrapped)
 * Takes in input 1 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsCOPY(result, ca *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	inOutParams := bk.params.InOutParams
	LweCopy(result, ca, inOutParams)
}

/*
 * Homomorphic Trivial Constant gate (doesn't need to be bootstrapped)
 * Takes a boolean value)
 * Outputs a LWE sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsCONSTANT(result *LweSample, value int32, bk *TFheGateBootstrappingCloudKeySet) {
	inOutParams := bk.params.InOutParams
	MU := ModSwitchToTorus32(1, 8)
	var muValue = -MU
	if value != 0 {
		muValue = MU
	}
	LweNoiselessTrivial(result, muValue, inOutParams)
}

/*
 * Homomorphic bootstrapped NOR gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNOR(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca - cb
	NorConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, NorConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndNY Gate: not(a) and b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsANDNY(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) - ca + cb
	AndNYConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndNYConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped AndYN Gate: a and not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsANDYN(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,-1/8) + ca - cb
	AndYNConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndYNConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrNY Gate: not(a) or b
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsORNY(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) - ca + cb
	OrNYConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrNYConst, inOutParams)
	LweSubTo(tempResult, ca, inOutParams)
	LweAddTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped OrYN Gate: a or not(b)
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsORYN(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	tempResult := NewLweSample(inOutParams)

	//compute: (0,1/8) + ca - cb
	OrYNConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult, OrYNConst, inOutParams)
	LweAddTo(tempResult, ca, inOutParams)
	LweSubTo(tempResult, cb, inOutParams)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	//tfheBootstrap_FFT(result, bk.bkFFT, MU, tempResult);
	tfheBootstrap(result, bk.bk, MU, tempResult)
}

/*
 * Homomorphic bootstrapped Mux(a,b,c) = a?b:c = a*b + not(a)*c
 * Takes in input 3 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsMUX(result, a, b, c *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	inOutParams := bk.params.InOutParams
	extracted_params := &bk.params.tgswParams.TlweParams.extractedLweparams

	tempResult := NewLweSample(inOutParams)
	tempResult1 := NewLweSample(extracted_params)
	u1 := NewLweSample(extracted_params)
	u2 := NewLweSample(extracted_params)

	//compute "AND(a,b)": (0,-1/8) + a + b
	AndConst := ModSwitchToTorus32(-1, 8)
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweAddTo(tempResult, a, inOutParams)
	LweAddTo(tempResult, b, inOutParams)
	// Bootstrap without KeySwitch
	// tfheBootstrapWoKS_FFT(u1, bk.bkFFT, MU, tempResult);
	tfheBootstrapWoKS(u1, bk.bk, MU, tempResult)

	//compute "AND(not(a),c)": (0,-1/8) - a + c
	LweNoiselessTrivial(tempResult, AndConst, inOutParams)
	LweSubTo(tempResult, a, inOutParams)
	LweAddTo(tempResult, c, inOutParams)
	// Bootstrap without KeySwitch
	//tfheBootstrapWoKS_FFT(u2, bk.bkFFT, MU, tempResult);
	tfheBootstrapWoKS(u2, bk.bk, MU, tempResult)

	// Add u1=u1+u2
	MuxConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(tempResult1, MuxConst, extracted_params)
	LweAddTo(tempResult1, u1, extracted_params)
	LweAddTo(tempResult1, u2, extracted_params)
	// Key switching
	//lweKeySwitch(result, bk.bkFFT.ks, tempResult1)
	lweKeySwitch(result, bk.bk.ks, tempResult1)
}
