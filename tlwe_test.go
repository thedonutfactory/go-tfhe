package tfhe

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	params512_1  = NewTLweParams(512, 1, 0., 1.)
	params512_2  = NewTLweParams(512, 2, 0., 1.)
	params1024_1 = NewTLweParams(1024, 1, 0., 1.)
	params1024_2 = NewTLweParams(1024, 2, 0., 1.)
	params2048_1 = NewTLweParams(2048, 1, 0., 1.)
	params2048_2 = NewTLweParams(2048, 2, 0., 1.)

	key512_1  = newRandomKey(params512_1)
	key512_2  = newRandomKey(params512_2)
	key1024_1 = newRandomKey(params1024_1)
	key1024_2 = newRandomKey(params1024_2)
	key2048_1 = newRandomKey(params2048_1)
	key2048_2 = newRandomKey(params2048_2)

	allParameters = []*TLweParams{params512_1, params512_2, params1024_1, params1024_2, params2048_1, params2048_2}
	allRandomKeys = []*TLweKey{key512_1, key512_2, key1024_1, key1024_2, key2048_1, key2048_2}
)

const toler double = 1e-8

func newRandomKey(params *TLweParams) *TLweKey {
	key := NewTLweKey(params)
	N := params.N
	k := params.k

	for i := int32(0); i < k; i++ {
		for j := int32(0); j < N; j++ {
			key.key[i].Coefs[j] = rand.Int31() % 2
		}
	}
	return key
}

/*
* Definition of the function fillRandom
* Fills a TLweSample with random Torus32 values (uniform distribution)
 */
func fillTLweRandom(result *TLweSample, params *TLweParams) {
	k := params.k
	N := params.N
	for i := int32(0); i <= k; i++ {
		for j := int32(0); j < N; j++ {
			result.A[i].CoefsT[j] = UniformTorus32Dist()
		}
	}
	result.CurrentVariance = 0.2
}

/*
* Definition of the function copySample
* Copies a TLweSample
 */
func copyTLweSample(result *TLweSample, sample *TLweSample, params *TLweParams) {
	k := params.k
	N := params.N

	for i := int32(0); i <= k; i++ {
		for j := int32(0); j < N; j++ {
			result.A[i].CoefsT[j] = sample.A[i].CoefsT[j]
		}
	}
	result.CurrentVariance = sample.CurrentVariance
}

/*
Testing the function tLweKeyGen
* EXPORT void tLweKeyGen(TLweKey* result)
*
* This function generates a random TLwe key for the given parameters
* The TLwe key for the result must be allocated and initialized
* (this means that the parameters are already in the result)
*/
//TEST_F(TLweTest, tLweKeyGen) {
func TestTLweKeyGen(t *testing.T) {
	assert := assert.New(t)
	for _, params := range allParameters {
		// Generating the key
		key := NewTLweKey(params)
		tLweKeyGen(key)
		assert.EqualValues(params, key.params)

		N := key.params.N
		k := key.params.k
		s := key.key

		//verify that the key is binary and kind-of random
		var count int32 = 0
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.True(s[i].Coefs[j] == 0 || s[i].Coefs[j] == 1)
				count += s[i].Coefs[j]
			}
		}
		assert.LessOrEqual(count, k*N-20)       // <=
		assert.GreaterOrEqual(count, int32(20)) // >=
	}
}

/*
	Testing the functions tLweSymEncryptT, tLwePhase, tLweSymDecryptT
	* EXPORT void tLweSymEncryptT(TLweSample* result, Torus32 message, double alpha, const TLweKey* key)
	* EXPORT void tLwePhase(TorusPolynomial* phase, const TLweSample* sample, const TLweKey* key)
	* EXPORT Torus32 tLweSymDecryptT(const TLweSample* sample, const TLweKey* key, int32_t Msize)
	*
	* This functions encrypt and decrypt a random Torus32 message by using the given key
*/
//TEST_F (TLweTest, tLweSymEncryptPhaseDecryptT) {
func TestTLweSymEncryptPhaseDecryptT(t *testing.T) {
	assert := assert.New(t)

	//TODO: parallelization
	const (
		NB_SAMPLES = 10
		M          = 8
		alpha      = 1. / (10. * M)
	)
	all_keys1024 := []*TLweKey{key1024_1, key1024_2}

	for _, key := range all_keys1024 {
		params := key.params
		N := params.N
		k := params.k
		samples := NewTLweSampleArray(NB_SAMPLES, params)
		phase := NewTorusPolynomial(N)
		var decrypt Torus32

		//verify correctness of the decryption
		for trial := 0; trial < NB_SAMPLES; trial++ {
			// The message is a Torus32
			message := ModSwitchToTorus32(rand.Int31()%M, M)

			// Encrypt and decrypt
			tLweSymEncryptT(&samples[trial], message, alpha, key)
			decrypt = tLweSymDecryptT(&samples[trial], key, M)
			//ILA: Testing APPROX correct decryption
			//the absolute value of the difference between message and decrypt is <= than toler
			assert.LessOrEqual(math.Abs(T32tod(message-decrypt)), toler)
			assert.LessOrEqual(math.Abs(T32tod(message-decrypt)), toler)

			// ILA: It is really necessary? phase used in decrypt!!!
			// Phase
			TLwePhase(phase, &samples[trial], key)
			// Testing phase
			dmessage := T32tod(message)
			dphase := T32tod(phase.CoefsT[0])
			assert.LessOrEqual(absfrac(dmessage-dphase), 10.*alpha) //ILA: why absfrac?
			//assert.EqualValues(alpha*alpha, samples[trial].CurrentVariance)
		}

		// Verify that samples are random enough (all coordinates different)
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				testset := make(map[Torus32]bool)
				//set<Torus32> testset
				for trial := 0; trial < NB_SAMPLES; trial++ {
					//testset.insert(samples[trial].a[i].CoefsT[j])
					testset[samples[trial].A[i].CoefsT[j]] = true
				}
				assert.GreaterOrEqual(double(len(testset)), 0.9*NB_SAMPLES) // >=
			}
		}
	}
}

/*
Testing the functions tLweSymEncrypt, tLwePhase, tLweApproxPhase, tLweSymDecrypt
* EXPORT void tLweSymEncrypt(TLweSample* result, TorusPolynomial* message, double alpha, const TLweKey* key)
* EXPORT void tLwePhase(TorusPolynomial* phase, const TLweSample* sample, const TLweKey* key)
* EXPORT void tLweApproxPhase(TorusPolynomial* message, const TorusPolynomial* phase, int32_t Msize, int32_t N)
* EXPORT void tLweSymDecrypt(TorusPolynomial* result, const TLweSample* sample, const TLweKey* key, int32_t Msize)
*
* This functions encrypt and decrypt a random TorusPolynomial message by using the given key
*/
//TEST_F (TLweTest, tLweSymEncryptPhaseDecrypt) {
func TestTLweSymEncryptPhaseDecrypt(t *testing.T) {
	assert := assert.New(t)
	const (
		NB_SAMPLES = 10
		M          = 8
		alpha      = 1. / (10. * M)
	)
	all_keys1024 := []*TLweKey{key1024_1, key1024_2}

	for _, key := range all_keys1024 {
		params := key.params
		N := params.N
		k := params.k

		samples := NewTLweSampleArray(NB_SAMPLES, params)
		message := NewTorusPolynomial(N)
		phase := NewTorusPolynomial(N)
		approxphase := NewTorusPolynomial(N)
		decrypt := NewTorusPolynomial(N)

		//verify correctness of the decryption
		for trial := 0; trial < NB_SAMPLES; trial++ {
			for j := int32(0); j < N; j++ {
				message.CoefsT[j] = ModSwitchToTorus32(rand.Int31()%M, M)
			}

			// Encrypt and Decrypt
			TLweSymEncrypt(&samples[trial], message, alpha, key)
			TLweSymDecrypt(decrypt, &samples[trial], key, M)
			//ILA: Testing APPROX correct decryption
			assert.LessOrEqual(torusPolynomialNormInftyDist(message, decrypt), toler)
			// for (int32_t j = 0; j < N; ++j) ASSERT_EQ(message.CoefsT[j],decrypt.CoefsT[j])

			// ILA: It is really necessary? phase and ApproxPhase used in decrypt!!!
			// Phase and ApproxPhase
			TLwePhase(phase, &samples[trial], key)
			TLweApproxPhase(approxphase, phase, M, N)
			// Testing Phase and ApproxPhase
			for j := int32(0); j < N; j++ {
				dmessage := T32tod(message.CoefsT[j])
				dphase := T32tod(phase.CoefsT[j])
				dapproxphase := T32tod(approxphase.CoefsT[j])
				assert.LessOrEqual(absfrac(dmessage-dphase), 10.*alpha)   // ILA: why absfrac?
				assert.LessOrEqual(absfrac(dmessage-dapproxphase), alpha) // ILA verify
			}

			//assert.EqualValues(alpha*alpha, samples[trial].CurrentVariance)
		}

		// Verify that samples are random enough (all coordinates different)
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				testset := make(map[Torus32]bool)
				for trial := 0; trial < NB_SAMPLES; trial++ {
					fmt.Println(samples[trial].A[i].CoefsT[j])
					testset[samples[trial].A[i].CoefsT[j]] = true
				}
				assert.GreaterOrEqual(double(len(testset)), 0.9*NB_SAMPLES) // >=
			}
		}
	}
}

/* **********************************
Arithmetic operations on TLwe samples
********************************** */

/*
Testing the function tLweClear
* EXPORT void tLweClear(TLweSample* result, const TLweParams* params)
*
* tLweClear sets the TLweSample to (0,0)
*/
//TEST_F(TLweTest, tLweClear) {
func TestTLweClear(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample := NewTLweSample(params)

		// Generate a random TLweSample and then set it to (0,0)
		fillTLweRandom(sample, params)
		tLweClear(sample, params)

		// Verify that the sample as been correctly set to (0,0)
		for i := int32(0); i <= k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(0, sample.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(0., sample.CurrentVariance)
	}
}

/*
Testing the function tLweCopy
* EXPORT void tLweCopy(TLweSample* result, const TLweSample* sample, const TLweParams* params)
*
* tLweCopy sets the (TLweSample) result equl to a given (TLweSample) sample
*/
func TestTLweCopy(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample := NewTLweSample(params)
		result := NewTLweSample(params)

		// Generate a random TLweSample and copy it
		fillTLweRandom(sample, params)
		tLweCopy(result, sample, params)

		// Verify that the sample as been correctly copied
		for i := int32(0); i <= k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(result.A[i].CoefsT[j], sample.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(result.CurrentVariance, sample.CurrentVariance)
	}
}

/*
Testing the function tLweNoiselessTrivial
* EXPORT void tLweNoiselessTrivial(TLweSample* result, const TorusPolynomial* mu, const TLweParams* params)
*
* tLweNoiselessTrivial sets the TLweSample to (0,mu)
*/
//TEST_F(TLweTest, tLweNoiselessTrivial) {
func TestTLweNoiselessTrivial(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		message := NewTorusPolynomial(N)
		for j := int32(0); j < N; j++ {
			message.CoefsT[j] = UniformTorus32Dist()
		}
		sample := NewTLweSample(params)

		// Generate a random TLweSample and set it to (0,mu)
		fillTLweRandom(sample, params)
		tLweNoiselessTrivial(sample, message, params)

		// Verify that the sample as been correctly set
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(0, sample.A[i].CoefsT[j])
			}
		}
		for j := int32(0); j < N; j++ {
			assert.EqualValues(message.CoefsT[j], sample.B.CoefsT[j])
		}
		assert.EqualValues(0., sample.CurrentVariance)
	}
}

/*
Testing the function tLweAddTo
* EXPORT void tLweAddTo(TLweSample* result, const TLweSample* sample, const TLweParams* params)
*
* tLweAddTo computes result = result + sample
*/
func TestTLweAddTo(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample1 := NewTLweSample(params)
		sample2 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)

		// Generate two random TLweSample and adds the second to the first
		fillTLweRandom(sample1, params)
		fillTLweRandom(sample2, params)
		copyTLweSample(sample1copy, sample1, params)
		tLweAddTo(sample1, sample2, params)

		// Verify if the operation was correctly executed
		for i := int32(0); i < k; i++ {
			// torusPolynomialAddTo(sample1copy.a[i], sample2.a[i])
			// Test equality between sample1copy.a[i] and sample1.a[i]
			for j := int32(0); j < N; j++ {
				assert.EqualValues(sample1copy.A[i].CoefsT[j]+sample2.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance+sample2.CurrentVariance, sample1.CurrentVariance)
	}
}

/*
Testing the function tLweSubTo
* EXPORT void tLweSubTo(TLweSample* result, const TLweSample* sample, const TLweParams* params)
*
* tLweSubTo computes result = result - sample
*/
func TestTLweSubTo(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample1 := NewTLweSample(params)
		sample2 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)

		// Generate two random TLweSample and subs the second to the first
		fillTLweRandom(sample1, params)
		fillTLweRandom(sample2, params)
		copyTLweSample(sample1copy, sample1, params)
		tLweSubTo(sample1, sample2, params)

		// Verify if the operation was correctly executed
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(sample1copy.A[i].CoefsT[j]-sample2.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance+sample2.CurrentVariance, sample1.CurrentVariance)
	}
}

/*
Testing the function tLweAddMulTo
* EXPORT void tLweAddMulTo(TLweSample* result, int32_t p, const TLweSample* sample, const TLweParams* params)
*
* tLweAddMulTo computes result = result + p.sample
*/
func TestTLweAddMulTo(t *testing.T) {
	assert := assert.New(t)
	const p int32 = 3
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample1 := NewTLweSample(params)
		sample2 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)

		// Generate two random TLweSample and adds the second multiplied by an integer p to the first
		fillTLweRandom(sample1, params)
		fillTLweRandom(sample2, params)
		copyTLweSample(sample1copy, sample1, params)
		tLweAddMulTo(sample1, p, sample2, params)

		// Verify if the operation was correctly executed
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(sample1copy.A[i].CoefsT[j]+p*sample2.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance+double(p*p)*sample2.CurrentVariance, sample1.CurrentVariance)

	}
}

/*
Testing the function tLweSubMulTo
* EXPORT void tLweSubMulTo(TLweSample* result, int32_t p, const TLweSample* sample, const TLweParams* params)
*
* tLweSubMulTo computes result = result - p.sample
*/
func TestTLweSubMulTo(t *testing.T) {
	assert := assert.New(t)
	const p int32 = 3
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		sample1 := NewTLweSample(params)
		sample2 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)

		// Generate two random TLweSample and subs the second multiplied by an integer p to the first
		fillTLweRandom(sample1, params)
		fillTLweRandom(sample2, params)
		copyTLweSample(sample1copy, sample1, params)
		tLweSubMulTo(sample1, p, sample2, params)

		// Verify if the operation was correctly executed
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.EqualValues(sample1copy.A[i].CoefsT[j]-p*sample2.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance+double(p*p)*sample2.CurrentVariance, sample1.CurrentVariance)
	}
}

/** result += (0,x) */
//EXPORT void tLweAddTTo(TLweSample* result, const Torus32 x, const TLweParams* params)
func TestTLweAddTTo(t *testing.T) {
	assert := assert.New(t)
	x := UniformTorus32Dist()
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		pos := rand.Int31() % params.k
		sample1 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)
		fillTLweRandom(sample1, params)
		copyTLweSample(sample1copy, sample1, params)
		tLweAddTTo(sample1, pos, x, params)
		// Verify if the operation was correctly executed
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				if i == pos && j == 0 {
					assert.EqualValues(sample1copy.A[i].CoefsT[j]+x, sample1.A[i].CoefsT[j])
				} else {
					assert.EqualValues(sample1copy.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
				}
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance, sample1.CurrentVariance)
	}
}

/** result += p*(0,x) */
//EXPORT void tLweAddRTTo(TLweSample* result, const IntPolynomial* p, const Torus32 x, const TLweParams* params)
func TestTLweAddRTTo(t *testing.T) {
	assert := assert.New(t)
	x := UniformTorus32Dist()
	for _, key := range allRandomKeys {
		params := key.params
		N := params.N
		k := params.k
		pos := rand.Int31() % params.k
		sample1 := NewTLweSample(params)
		sample1copy := NewTLweSample(params)
		p := NewIntPolynomial(N)
		fillTLweRandom(sample1, params)
		for i := int32(0); i < N; i++ {
			p.Coefs[i] = UniformTorus32Dist() % 1000
		}
		copyTLweSample(sample1copy, sample1, params)
		tLweAddRTTo(sample1, pos, p, x, params)
		// Verify if the operation was correctly executed
		for i := int32(0); i <= k; i++ {
			for j := int32(0); j < N; j++ {
				if i != pos {
					assert.EqualValues(sample1copy.A[i].CoefsT[j], sample1.A[i].CoefsT[j])
				} else {
					assert.EqualValues(sample1copy.A[i].CoefsT[j]+p.Coefs[j]*x, sample1.A[i].CoefsT[j])
				}
			}
		}
		assert.EqualValues(sample1copy.CurrentVariance, sample1.CurrentVariance)
	}
}
