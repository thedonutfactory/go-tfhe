package tfhe

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

type LweKey struct {
	params *LweParams
	key    []int32
}

//this structure contains Lwe parameters
//this structure is constant (cannot be modified once initialized):
//the pointer to the param can be passed directly
//to all the Lwe keys that use these params.
type LweParams struct {
	N        int32
	alphaMin double //le plus petit bruit tq sur
	alphaMax double //le plus gd bruit qui permet le déchiffrement
}

type LweSample struct {
	A               []Torus32 // the n coefs of the mask
	B               Torus32
	CurrentVariance double
}

//func (s *LweSample) B() Torus32 {
//	return s.A[len(s.A)]
//}

func NewLweParams(n int32, min, max double) *LweParams {
	return &LweParams{n, min, max}
}

func NewLweKey(params *LweParams) *LweKey {
	return &LweKey{params: params, key: make([]int32, params.N)}
}

func NewLweSample(params *LweParams) *LweSample {
	return &LweSample{A: make([]Torus32, params.N), B: 0, CurrentVariance: 0}
}

func NewLweSampleArray(size int32, params *LweParams) (arr []*LweSample) {
	arr = make([]*LweSample, size)
	for i := int32(0); i < size; i++ {
		arr[i] = NewLweSample(params)
	}
	return
}

/** generate a new unititialized ciphertext array of length nbelems */
func NewGateBootstrappingCiphertextArray(nbelems int, params *TFheGateBootstrappingParameterSet) (arr []*LweSample) {
	return NewLweSampleArray(int32(nbelems), params.InOutParams)
}

/**
 * This function generates a random Lwe key for the given parameters.
 * The Lwe key for the result must be allocated and initialized
 * (this means that the parameters are already in the result)
 */
func LweKeyGen(result *LweKey) {
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}

	z := make([]int32, result.params.N)
	// Generate some random numbers from standard normal distribution
	for i := range z {
		z[i] = Torus32(math.Round(dist.Rand()))
	}
	result.key = z
}

// variablize for use with test mocking
var LweSymEncrypt = LweSymEncryptImpl

/**
 * This function encrypts message by using key, with stdev alpha
 * The Lwe sample for the result must be allocated and initialized
 * (this means that the parameters are already in the result)
 */
func LweSymEncryptImpl(result *LweSample, message Torus32, alpha double, key *LweKey) {
	result.B = gaussian32(message, alpha)
	for i := 0; i < int(key.params.N); i++ {
		result.A[i] = UniformTorus32Dist()
		result.B += result.A[i] * key.key[i]
	}
	result.CurrentVariance = alpha * alpha
}

/*
 * This function encrypts a message by using key and a given noise value
 */
func LweSymEncryptWithExternalNoise(result *LweSample, message Torus32, noise double, alpha double, key *LweKey) {
	result.B = message + int32(noise)
	for i := 0; i < int(key.params.N); i++ {
		result.A[i] = UniformTorus32Dist()
		result.B += result.A[i] * key.key[i]
	}

	result.CurrentVariance = alpha * alpha
}

/**
 * This function computes the phase of sample by using key : phi = b - a.s
 */
func LwePhase(sample *LweSample, key *LweKey) Torus32 {
	var axs Torus32 = 0
	for i := 0; i < int(key.params.N); i++ {
		axs += sample.A[i] * key.key[i]
	}
	return sample.B - axs
}

/**
 * This function computes the decryption of sample by using key
 * The constant Msize indicates the message space and is used to approximate the phase
 */
func LweSymDecrypt(sample *LweSample, key *LweKey, Msize int32) Torus32 {
	phi := LwePhase(sample, key)
	return approxPhase(phi, Msize)
}

//Arithmetic operations on Lwe samples
/** result = (0,0) */
func LweClear(result *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] = 0
	}
	result.B = 0
	result.CurrentVariance = 0.
}

/** result = sample */
func LweCopy(result *LweSample, sample *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] = sample.A[i]
	}
	result.B = sample.B
	result.CurrentVariance = sample.CurrentVariance
}

/** result = -sample */
func LweNegate(result *LweSample, sample *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] = -sample.A[i]
	}
	result.B = -sample.B
	result.CurrentVariance = sample.CurrentVariance
}

/** result = (0,mu) */
func LweNoiselessTrivial(result *LweSample, mu Torus32, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] = 0
	}
	result.B = mu
	result.CurrentVariance = 0.
}

/** result = result + sample */
func LweAddTo(result *LweSample, sample *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] += sample.A[i]
	}
	result.B += sample.B
	result.CurrentVariance += sample.CurrentVariance
}

/** result = result - sample */
func LweSubTo(result *LweSample, sample *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] -= sample.A[i]
	}
	result.B -= sample.B
	result.CurrentVariance += sample.CurrentVariance
}

/** result = result + p.sample */
func LweAddMulTo(result *LweSample, p int32, sample *LweSample, params *LweParams) {
	for i := 0; i < int(params.N); i++ {
		result.A[i] += p * sample.A[i]
	}
	result.B += p * sample.B
	result.CurrentVariance += float64(p*p) * sample.CurrentVariance
}

/** result = result - p.sample */
func LweSubMulTo(result *LweSample, p int32, sample *LweSample, params *LweParams) {
	for i := int32(0); i < params.N; i++ {
		result.A[i] -= p * sample.A[i]
	}
	result.B -= p * sample.B
	result.CurrentVariance += float64(p*p) * sample.CurrentVariance
}

func tLweExtractLweSampleIndex(result *LweSample, x *TLweSample, index int32, params *LweParams, rparams *TLweParams) {
	N := rparams.N
	k := rparams.K
	Assert(params.N == k*N)

	for i := int32(0); i < k; i++ {
		for j := int32(0); j <= index; j++ {
			result.A[i*N+j] = x.A[i].CoefsT[index-j]
		}
		for j := index + 1; j < N; j++ {
			result.A[i*N+j] = -x.A[i].CoefsT[N+index-j]
		}
	}
	result.B = x.B().CoefsT[index]
}

func tLweExtractLweSample(result *LweSample, x *TLweSample, params *LweParams, rparams *TLweParams) {
	tLweExtractLweSampleIndex(result, x, 0, params, rparams)
}

func tLweExtractKey(result *LweKey, key *TLweKey) {
	N := key.params.N
	k := key.params.K
	Assert(result.params.N == k*N)
	for i := int32(0); i < k; i++ {
		for j := int32(0); j < N; j++ {
			result.key[i*N+j] = key.key[i].Coefs[j]
		}
	}
}
