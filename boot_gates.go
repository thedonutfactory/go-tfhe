package tfhe

/*
 * Homomorphic bootstrapped NAND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNAND(result, ca, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
	MU := ModSwitchToTorus32(1, 8)
	in_out_params := bk.params.InOutParams
	temp_result := NewLweSample(in_out_params)

	//compute: (0,1/8) - ca - cb
	NandConst := ModSwitchToTorus32(1, 8)
	LweNoiselessTrivial(temp_result, NandConst, in_out_params)
	LweSubTo(temp_result, ca, in_out_params)
	LweSubTo(temp_result, cb, in_out_params)

	//if the phase is positive, the result is 1/8
	//if the phase is positive, else the result is -1/8
	tfhe_bootstrap(result, bk.bk, MU, temp_result)
	// tfhe_bootstrap_FFT(result, bk.bkFFT, MU, temp_result)
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
