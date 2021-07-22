package tfhe

type TFheGateBootstrappingParameterSet struct {
	ks_t          int
	ks_basebit    int
	in_out_params *LweParams
	tgsw_params   *TGswParams
}

type TFheGateBootstrappingCloudKeySet struct {
	params *TFheGateBootstrappingParameterSet
	bk     *LweBootstrappingKey
	// bkFFT  *LweBootstrappingKeyFFT
}

type TFheGateBootstrappingSecretKeySet struct {
	params   *TFheGateBootstrappingParameterSet
	lwe_key  *LweKey
	tgsw_key *TGswKey
	cloud    TFheGateBootstrappingCloudKeySet
}
