package tfhe

/*
 * Homomorphic bootstrapped NAND gate
 * Takes in input 2 LWE samples (with message space [-1/8,1/8], noise<1/16)
 * Outputs a LWE bootstrapped sample (with message space [-1/8,1/8], noise<1/16)
 */
func BootsNAND(result *LweSample, ca *LweSample, cb *LweSample, bk *TFheGateBootstrappingCloudKeySet) {
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
