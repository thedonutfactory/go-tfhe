package tfhe

import "math"

type GateBootstrappingParameterSet struct {
	KsT         int32
	KsBasebit   int32
	InOutParams *LweParams
	TgswParams  *TGswParams
}

type PublicKey struct {
	Params *GateBootstrappingParameterSet
	Bkw    *LweBootstrappingKeyWrapper
	//Bk     *LweBootstrappingKey
	//BkFFT *LweBootstrappingKeyFFT
}

type PrivateKey struct {
	Params  *GateBootstrappingParameterSet
	LweKey  *LweKey
	TgswKey *TGswKey
	//Cloud   *PublicKey
}

func NewTFheGateBootstrappingParameterSet(ksT, ksBasebit int32, inOutParams *LweParams, tgswParams *TGswParams) *GateBootstrappingParameterSet {
	return &GateBootstrappingParameterSet{
		KsT:         ksT,
		KsBasebit:   ksBasebit,
		InOutParams: inOutParams,
		TgswParams:  tgswParams,
	}
}

func NewPublicKey(params *GateBootstrappingParameterSet, bkw *LweBootstrappingKeyWrapper) *PublicKey {

	return &PublicKey{
		Params: params,
		Bkw:    bkw,
	}
}

func NewPrivateKey(params *GateBootstrappingParameterSet, bk *LweBootstrappingKeyWrapper, lweKey *LweKey, tgswKey *TGswKey) *PrivateKey {
	return &PrivateKey{
		Params:  params,
		LweKey:  lweKey,
		TgswKey: tgswKey,
		//Cloud:   NewPublicKey(params, bk),
	}
}

func Default80bitGateBootstrappingParameters() *GateBootstrappingParameterSet {
	// These are the historic parameter set provided in 2016,
	// that were analyzed against lattice enumeration attacks
	// Currently (in 2020), the security of these parameters is estimated to lambda = **80-bit security**
	// (w.r.t bkz-sieve model, + hybrid attack model)
	const (
		N         = 1024
		k         = 1
		n         = 500
		bkL       = 2
		bkBgbit   = 10
		ksBasebit = 2
		ksLength  = 8
		ksStdev   = 2.44e-5  //standard deviation
		bkStdev   = 7.18e-9  //standard deviation
		maxStdev  = 0.012467 //max standard deviation for a 1/4 msg space
	)

	paramsIn := NewLweParams(n, ksStdev, maxStdev)
	paramsAccum := NewTLweParams(N, k, bkStdev, maxStdev)
	paramsBk := NewTGswParams(bkL, bkBgbit, paramsAccum)

	return NewTFheGateBootstrappingParameterSet(ksLength, ksBasebit, paramsIn, paramsBk)
}

// this is the default and recommended parameter set
func Default128bitGateBootstrappingParameters() *GateBootstrappingParameterSet {
	// These are the parameter set provided in CGGI2019.
	// Currently (in 2020), the security of these parameters is estimated to lambda = 129-bit security
	// (w.r.t bkz-sieve model, + additional hybrid attack models)
	const (
		N         = 1024
		k         = 1
		n         = 630
		bkL       = 3
		bkBgbit   = 7
		ksBasebit = 2
		ksLength  = 8
		maxStdev  = 0.012467 //max standard deviation for a 1/4 msg space
	)

	ksStdev := math.Pow(2, -15) //standard deviation
	bkStdev := math.Pow(2, -25) //standard deviation

	paramsIn := NewLweParams(n, ksStdev, maxStdev)
	paramsAccum := NewTLweParams(N, k, bkStdev, maxStdev)
	paramsBk := NewTGswParams(bkL, bkBgbit, paramsAccum)

	return NewTFheGateBootstrappingParameterSet(ksLength, ksBasebit, paramsIn, paramsBk)
}

func NewDefaultGateBootstrappingParameters(minimumLambda int32) *GateBootstrappingParameterSet {
	if minimumLambda > 128 {
		panic("Sorry, for now, the parameters are only implemented for 80bit and 128bit of security!")
	}

	if minimumLambda > 80 && minimumLambda <= 128 {
		return Default128bitGateBootstrappingParameters()
	}

	if minimumLambda > 0 && minimumLambda <= 80 {
		return Default80bitGateBootstrappingParameters()
	}

	panic("the requested security parameter must be positive (currently, 80 and 128-bits are supported)")
}

func GenerateKeys(params *GateBootstrappingParameterSet) (*PublicKey, *PrivateKey) {
	lweKey := NewLweKey(params.InOutParams)
	LweKeyGen(lweKey)
	tgswKey := NewTGswKey(params.TgswParams)
	TGswKeyGen(tgswKey)
	bkw := NewLweBootstrappingKeyWrapper(params.KsT, params.KsBasebit, params.InOutParams, params.TgswParams, lweKey, tgswKey)
	//tfheCreateLweBootstrappingKey(bkw.Bk, lweKey, tgswKey)
	return NewPublicKey(params, bkw), NewPrivateKey(params, bkw, lweKey, tgswKey)
}

/** encrypts a boolean */
func BootsSymEncrypt(result *LweSample, message int32, key *PrivateKey) {
	_1s8 := ModSwitchToTorus32(1, 8)
	var mu Torus32 = -_1s8
	if message != 0 {
		mu = _1s8
	}
	//Torus32 mu = message ? _1s8 : -_1s8;
	alpha := key.Params.InOutParams.AlphaMin //TODO: specify noise
	LweSymEncrypt(result, mu, alpha, key.LweKey)
}

/** decrypts a boolean */
func BootsSymDecrypt(sample *LweSample, key *PrivateKey) int32 {
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
}
