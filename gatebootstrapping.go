package tfhe

import "math"

type TFheGateBootstrappingParameterSet struct {
	ks_t        int32
	ks_basebit  int32
	InOutParams *LweParams
	tgswParams  *TGswParams
}

type TFheGateBootstrappingCloudKeySet struct {
	params *TFheGateBootstrappingParameterSet
	bk     *LweBootstrappingKey
	// bkFFT  *LweBootstrappingKeyFFT
}

type TFheGateBootstrappingSecretKeySet struct {
	Params  *TFheGateBootstrappingParameterSet
	LweKey  *LweKey
	tgswKey *TGswKey
	Cloud   *TFheGateBootstrappingCloudKeySet
}

func NewTFheGateBootstrappingParameterSet(ks_t, ks_basebit int32, in_out_params *LweParams, tgsw_params *TGswParams) *TFheGateBootstrappingParameterSet {
	return &TFheGateBootstrappingParameterSet{
		ks_t:        ks_t,
		ks_basebit:  ks_basebit,
		InOutParams: in_out_params,
		tgswParams:  tgsw_params,
	}
}

func NewTFheGateBootstrappingCloudKeySet(params *TFheGateBootstrappingParameterSet, bk *LweBootstrappingKey) *TFheGateBootstrappingCloudKeySet {
	return &TFheGateBootstrappingCloudKeySet{
		params: params,
		bk:     bk,
	}
}

func NewTFheGateBootstrappingSecretKeySet(params *TFheGateBootstrappingParameterSet, bk *LweBootstrappingKey, lwe_key *LweKey, tgsw_key *TGswKey) *TFheGateBootstrappingSecretKeySet {
	return &TFheGateBootstrappingSecretKeySet{
		Params:  params,
		LweKey:  lwe_key,
		tgswKey: tgsw_key,
		Cloud:   NewTFheGateBootstrappingCloudKeySet(params, bk),
	}
}

func Default80bitGateBootstrappingParameters() *TFheGateBootstrappingParameterSet {
	// These are the historic parameter set provided in 2016,
	// that were analyzed against lattice enumeration attacks
	// Currently (in 2020), the security of these parameters is estimated to lambda = **80-bit security**
	// (w.r.t bkz-sieve model, + hybrid attack model)
	const (
		N          = 1024
		k          = 1
		n          = 500
		bk_l       = 2
		bk_Bgbit   = 10
		ks_basebit = 2
		ks_length  = 8
		ks_stdev   = 2.44e-5  //standard deviation
		bk_stdev   = 7.18e-9  //standard deviation
		max_stdev  = 0.012467 //max standard deviation for a 1/4 msg space
	)

	params_in := NewLweParams(n, ks_stdev, max_stdev)
	params_accum := NewTLweParams(N, k, bk_stdev, max_stdev)
	params_bk := NewTGswParams(bk_l, bk_Bgbit, params_accum)

	return NewTFheGateBootstrappingParameterSet(ks_length, ks_basebit, params_in, params_bk)
}

// this is the default and recommended parameter set
func Default128bitGateBootstrappingParameters() *TFheGateBootstrappingParameterSet {
	// These are the parameter set provided in CGGI2019.
	// Currently (in 2020), the security of these parameters is estimated to lambda = 129-bit security
	// (w.r.t bkz-sieve model, + additional hybrid attack models)
	const (
		N          = 1024
		k          = 1
		n          = 630
		bk_l       = 3
		bk_Bgbit   = 7
		ks_basebit = 2
		ks_length  = 8
		max_stdev  = 0.012467 //max standard deviation for a 1/4 msg space
	)

	ks_stdev := math.Pow(2, -15) //standard deviation
	bk_stdev := math.Pow(2, -25) //standard deviation

	params_in := NewLweParams(n, ks_stdev, max_stdev)
	params_accum := NewTLweParams(N, k, bk_stdev, max_stdev)
	params_bk := NewTGswParams(bk_l, bk_Bgbit, params_accum)

	return NewTFheGateBootstrappingParameterSet(ks_length, ks_basebit, params_in, params_bk)
}

func NewDefaultGateBootstrappingParameters(minimum_lambda int32) *TFheGateBootstrappingParameterSet {
	if minimum_lambda > 128 {
		panic("Sorry, for now, the parameters are only implemented for 80bit and 128bit of security!")
	}

	if minimum_lambda > 80 && minimum_lambda <= 128 {
		return Default128bitGateBootstrappingParameters()
	}

	if minimum_lambda > 0 && minimum_lambda <= 80 {
		return Default80bitGateBootstrappingParameters()
	}

	panic("the requested security parameter must be positive (currently, 80 and 128-bits are supported)")
}

func NewRandomGateBootstrappingSecretKeyset(params *TFheGateBootstrappingParameterSet) *TFheGateBootstrappingSecretKeySet {
	lwe_key := NewLweKey(params.InOutParams)
	LweKeyGen(lwe_key)
	tgsw_key := NewTGswKey(params.tgswParams)
	TGswKeyGen(tgsw_key)
	bk := NewLweBootstrappingKey(params.ks_t, params.ks_basebit, params.InOutParams,
		params.tgswParams)
	tfhe_createLweBootstrappingKey(bk, lwe_key, tgsw_key)
	return NewTFheGateBootstrappingSecretKeySet(params, bk, lwe_key, tgsw_key)
}

/** encrypts a boolean */
func BootsSymEncrypt(result *LweSample, message int32, key *TFheGateBootstrappingSecretKeySet) {
	_1s8 := ModSwitchToTorus32(1, 8)
	var mu Torus32 = -_1s8
	if message != 0 {
		mu = _1s8
	}
	//Torus32 mu = message ? _1s8 : -_1s8;
	alpha := key.Params.InOutParams.alphaMin //TODO: specify noise
	LweSymEncrypt(result, mu, alpha, key.LweKey)
}

/** decrypts a boolean */
func BootsSymDecrypt(sample *LweSample, key *TFheGateBootstrappingSecretKeySet) int32 {
	mu := LwePhase(sample, key.LweKey)
	/*
		if mu != 0 {
			return 1
		} else {
			return 0
		}
	*/

	if mu > 0 {
		return 1
	} else {
		return 0
	}

	//return (mu > 0 ? 1 : 0); //we have to do that because of the C binding
}
