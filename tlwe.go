package tfhe

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

type TLweParams struct {
	N                  int       ///< a power of 2: degree of the polynomials
	K                  int       ///< number of polynomials in the mask
	alphaMin           double    ///< minimal noise s.t. the sample is secure
	alphaMax           double    ///< maximal noise s.t. we can decrypt
	extractedLweparams LweParams ///< lwe params if one extracts
}

type TLweKey struct {
	params *TLweParams     ///< the parameters of the key
	key    []IntPolynomial ///< the key (i.e k binary polynomials)
}

type TLweSample struct {
	A []TorusPolynomial ///< array of length k+1: mask + right term
	//B               *TorusPolynomial  ///< alias of a[k] to get the right term
	CurrentVariance double ///< avg variance of the sample
	K               int
}

// alias of a[k] to get the right term
func (s *TLweSample) B() *TorusPolynomial {
	return &s.A[s.K]
}

func NewTLweParams(N, k int, alphaMin, alphaMax double) *TLweParams {
	return &TLweParams{
		N:                  N,
		K:                  k,
		alphaMin:           alphaMin,
		alphaMax:           alphaMax,
		extractedLweparams: LweParams{N * k, alphaMin, alphaMax},
	}
}

func NewTLweSample(params *TLweParams) *TLweSample {
	avar := NewTorusPolynomialArray(int(params.K)+1, params.N)
	return &TLweSample{
		K: params.K,
		A: avar,
		//B:               &avar[params.K],
		CurrentVariance: 0.,
	}
}

func NewTLweSampleArray(size int, params *TLweParams) (arr []*TLweSample) {
	arr = make([]*TLweSample, size)
	for i := 0; i < size; i++ {
		arr[i] = NewTLweSample(params)
	}
	return
}

func NewTLweKey(params *TLweParams) *TLweKey {
	return &TLweKey{params: params, key: NewIntPolynomialArray(int(params.K), params.N)}
}

func (s *TLweSample) DebugTLweSample() {
	tabs(3, "TLweSample {")
	for i := 0; i < len(s.A); i++ {
		tabs(4, fmt.Sprintf("a.coefs [%d] [", i))
		// for j := 0; j < len(s.A[i].CoefsT); j++ {
		for j := 0; j < 5; j++ {
			v := s.A[i].CoefsT[j]
			if v != 0 {
				tabsi(5, v)
			}
		}
		tabs(4, "]")
	}

	tabs(4, "b.coefs [")
	// for i := 0; i < len(s.B().CoefsT); i++ {
	for i := 0; i < 5; i++ {
		v := s.B().CoefsT[i]
		if v != 0 {
			tabsi(5, v)
		}

	}
	tabs(4, "]")
	tabs(3, "}")
}

func TLweKeyGen(result *TLweKey) {
	N := result.params.N
	k := result.params.K
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
	}
	for i := 0; i < k; i++ {
		for j := 0; j < N; j++ {
			result.key[i].Coefs[j] = Torus(math.Round(dist.Rand()))
		}
	}
}

/*create an homogeneous tlwe sample*/
func tLweSymEncryptZero(result *TLweSample, alpha double, key *TLweKey) {
	N := key.params.N
	k := key.params.K

	for j := 0; j < N; j++ {
		result.B().CoefsT[j] = gaussian32(0, alpha)
	}

	for i := 0; i < k; i++ {
		torusPolynomialUniform(&result.A[i])
		TorusPolynomialAddMulR(result.B(), &key.key[i], &result.A[i])
	}

	result.CurrentVariance = alpha * alpha
}

func TLweSymEncrypt(result *TLweSample, message *TorusPolynomial, alpha double, key *TLweKey) {
	tLweSymEncryptZero(result, alpha, key)
	for j := 0; j < key.params.N; j++ {
		result.B().CoefsT[j] += message.CoefsT[j]
	}
}

/**
 * encrypts a constant message
 */
func TLweSymEncryptT(result *TLweSample, message Torus, alpha double, key *TLweKey) {
	tLweSymEncryptZero(result, alpha, key)
	result.B().CoefsT[0] += message
}

/**
 * This function computes the phase of sample by using key : phi = b - a.s
 */
func TLwePhase(phase *TorusPolynomial, sample *TLweSample, key *TLweKey) {
	TorusPolynomialCopy(phase, sample.B()) // phi = b
	for i := 0; i < key.params.K; i++ {
		TorusPolynomialSubMulR(phase, &key.key[i], &sample.A[i])
	}
}

/**
 * This function computes the approximation of the phase
 * à revoir, surtout le Msize
 */
func TLweApproxPhase(message *TorusPolynomial, phase *TorusPolynomial, Msize, N int) {
	for i := 0; i < N; i++ {
		message.CoefsT[i] = approxPhase(phase.CoefsT[i], Msize)
	}
}

func TLweSymDecrypt(result *TorusPolynomial, sample *TLweSample, key *TLweKey, Msize int) {
	TLwePhase(result, sample, key)
	TLweApproxPhase(result, result, Msize, key.params.N)
}

func TLweSymDecryptT(sample *TLweSample, key *TLweKey, Msize int) Torus {
	phase := NewTorusPolynomial(key.params.N)
	TLwePhase(phase, sample, key)
	result := approxPhase(phase.CoefsT[0], Msize)
	return result
}

//Arithmetic operations on TLwe samples
/** result = (0,0) */
func TLweClear(result *TLweSample, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		torusPolynomialClear(&result.A[i])
	}
	torusPolynomialClear(result.B())
	result.CurrentVariance = 0.
}

/** result = sample */
func TLweCopy(result *TLweSample, sample *TLweSample, params *TLweParams) {
	for i := 0; i <= params.K; i++ {
		for j := 0; j < params.N; j++ {
			result.A[i].CoefsT[j] = sample.A[i].CoefsT[j]
		}
	}
	result.CurrentVariance = sample.CurrentVariance
}

/** result = (0,mu) */
func TLweNoiselessTrivial(result *TLweSample, mu *TorusPolynomial, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		torusPolynomialClear(&result.A[i])
	}
	TorusPolynomialCopy(result.B(), mu)
	result.CurrentVariance = 0.
}

/** result = (0,mu) where mu is constant*/
func TLweNoiselessTrivialT(result *TLweSample, mu Torus, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		torusPolynomialClear(&result.A[i])
	}
	torusPolynomialClear(result.B())
	result.B().CoefsT[0] = mu
	result.CurrentVariance = 0.
}

/** result = result + sample */
func TLweAddTo(result *TLweSample, sample *TLweSample, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		TorusPolynomialAddTo(&result.A[i], &sample.A[i])
	}
	TorusPolynomialAddTo(result.B(), sample.B())
	result.CurrentVariance += sample.CurrentVariance
}

/** result = result - sample */
func TLweSubTo(result *TLweSample, sample *TLweSample, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		TorusPolynomialSubTo(&result.A[i], &sample.A[i])
	}
	TorusPolynomialSubTo(result.B(), sample.B())
	result.CurrentVariance += sample.CurrentVariance
}

/** result = result + p.sample */
func TLweAddMulTo(result *TLweSample, p int64, sample *TLweSample, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		TorusPolynomialAddMulZTo(&result.A[i], p, &sample.A[i])
	}
	TorusPolynomialAddMulZTo(result.B(), p, sample.B())
	result.CurrentVariance += double(p*p) * sample.CurrentVariance
}

/** result = result - p.sample */
func TLweSubMulTo(result *TLweSample, p int64, sample *TLweSample, params *TLweParams) {
	for i := 0; i < params.K; i++ {
		TorusPolynomialSubMulZTo(&result.A[i], p, &sample.A[i])
	}
	TorusPolynomialSubMulZTo(result.B(), p, sample.B())
	result.CurrentVariance += double(p*p) * sample.CurrentVariance
}

/** result = result + p.sample */
func TLweAddMulRTo(result *TLweSample, p *IntPolynomial, sample *TLweSample, params *TLweParams) {
	for i := 0; i <= params.K; i++ {
		TorusPolynomialAddMulR(&result.A[i], p, &sample.A[i])
	}
	result.CurrentVariance += double(intPolynomialNormSq2(p)) * sample.CurrentVariance
}

//mult externe de X^ai-1 par bki
func TLweMulByXaiMinusOne(result *TLweSample, ai int64, bk *TLweSample, params *TLweParams) {
	for i := 0; i <= params.K; i++ {
		TorusPolynomialMulByXaiMinusOne(&result.A[i], ai, &bk.A[i])
	}
}

/** result += (0,x) */
func tLweAddTTo(result *TLweSample, pos int, x Torus, params *TLweParams) {
	result.A[pos].CoefsT[0] += x
}

/** result += p*(0,x) */
func tLweAddRTTo(result *TLweSample, pos int, p *IntPolynomial, x Torus, params *TLweParams) {
	for i := 0; i < params.N; i++ {
		result.A[pos].CoefsT[i] += p.Coefs[i] * x
	}
}
