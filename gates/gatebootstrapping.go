package gates

import (
	"math"

	"github.com/thedonutfactory/go-tfhe/core"
)

type GateBootstrappingParameterSet struct {
	KsT         int32
	KsBasebit   int32
	InOutParams *core.LweParams
	TgswParams  *core.TGswParams
}

type PublicKey struct {
	Params *GateBootstrappingParameterSet
	Bkw    *core.LweBootstrappingKeyWrapper
}

type PrivateKey struct {
	Params  *GateBootstrappingParameterSet
	LweKey  *core.LweKey
	TgswKey *core.TGswKey
}

/** generate a new uninitialized ciphertext array of length nbelems */
func NewGateBootstrappingCiphertextArray(nbelems int, params *GateBootstrappingParameterSet) (arr []*core.LweSample) {
	return core.NewLweSampleArray(int32(nbelems), params.InOutParams)
}

func NewTFheGateBootstrappingParameterSet(ksT, ksBasebit int32, inOutParams *core.LweParams, tgswParams *core.TGswParams) *GateBootstrappingParameterSet {
	return &GateBootstrappingParameterSet{
		KsT:         ksT,
		KsBasebit:   ksBasebit,
		InOutParams: inOutParams,
		TgswParams:  tgswParams,
	}
}

func NewPublicKey(params *GateBootstrappingParameterSet, bkw *core.LweBootstrappingKeyWrapper) *PublicKey {

	return &PublicKey{
		Params: params,
		Bkw:    bkw,
	}
}

func NewPrivateKey(params *GateBootstrappingParameterSet, bk *core.LweBootstrappingKeyWrapper, lweKey *core.LweKey, tgswKey *core.TGswKey) *PrivateKey {
	return &PrivateKey{
		Params:  params,
		LweKey:  lweKey,
		TgswKey: tgswKey,
		//Cloud:   NewPublicKey(params, bk),
	}
}

// Default parameter set.
//
// This parameter set ensures 128-bits of security, and a probability of error is upper-bounded by
// $2^{-25}$. The secret keys generated with this parameter set are uniform binary.
// This parameter set allows to evaluate faster Boolean circuits than the `Default128bitGateBootstrappingParameters`
// one.
func NewDefaultGateBootstrappingParameters() *GateBootstrappingParameterSet {
	const (
		N         = 512
		k         = 2
		n         = 586
		bkL       = 2
		bkBgbit   = 8
		ksBasebit = 2
		ksLength  = 5
		ksStdev   = 0.000_089_761_673_968_349_98     // 2^{-13.44...}
		bkStdev   = 0.000_000_029_890_407_929_674_34 // 2^{-24.9...}
		maxStdev  = 0.012467                         //max standard deviation for a 1/4 msg space
	)

	paramsIn := core.NewLweParams(n, ksStdev, maxStdev)
	paramsAccum := core.NewTLweParams(N, k, bkStdev, maxStdev)
	paramsBk := core.NewTGswParams(bkL, bkBgbit, paramsAccum)

	return NewTFheGateBootstrappingParameterSet(ksLength, ksBasebit, paramsIn, paramsBk)
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

	paramsIn := core.NewLweParams(n, ksStdev, maxStdev)
	paramsAccum := core.NewTLweParams(N, k, bkStdev, maxStdev)
	paramsBk := core.NewTGswParams(bkL, bkBgbit, paramsAccum)

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

	paramsIn := core.NewLweParams(n, ksStdev, maxStdev)
	paramsAccum := core.NewTLweParams(N, k, bkStdev, maxStdev)
	paramsBk := core.NewTGswParams(bkL, bkBgbit, paramsAccum)

	return NewTFheGateBootstrappingParameterSet(ksLength, ksBasebit, paramsIn, paramsBk)
}

func DefaultGateBootstrappingParameters(minimumLambda int32) *GateBootstrappingParameterSet {
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

func (params *GateBootstrappingParameterSet) GenerateKeys() (*PublicKey, *PrivateKey) {
	lweKey := core.NewLweKey(params.InOutParams)
	core.LweKeyGen(lweKey)
	tgswKey := core.NewTGswKey(params.TgswParams)
	core.TGswKeyGen(tgswKey)
	bkw := core.NewLweBootstrappingKeyWrapper(params.KsT, params.KsBasebit, params.InOutParams, params.TgswParams, lweKey, tgswKey)
	//tfheCreateLweBootstrappingKey(bkw.Bk, lweKey, tgswKey)
	return NewPublicKey(params, bkw), NewPrivateKey(params, bkw, lweKey, tgswKey)
}
