package tfhe

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	params500  *LweParams = &LweParams{N: 500, AlphaMin: 0, AlphaMax: 1}
	params750  *LweParams = &LweParams{N: 750, AlphaMin: 0, AlphaMax: 1}
	params1024 *LweParams = &LweParams{N: 1024, AlphaMin: 0, AlphaMax: 1}

	key500  *LweKey = newRandomLweKey(params500)
	key750  *LweKey = newRandomLweKey(params750)
	key1024 *LweKey = newRandomLweKey(params1024)

	allParams = []*LweParams{params500, params750, params1024}
	allKeys   = []*LweKey{key500, key750, key1024}
)

//this function creates a new lwekey and initializes it with random
//values. We do not use the c++11 random generator, since it gets in
//a deadlock mode on static const initializers
func newRandomLweKey(params *LweParams) *LweKey {

	key := &LweKey{Params: params, Key: make([]int32, params.N)}
	for i := int32(0); i < params.N; i++ {
		key.Key[i] = rand.Int31() % 2
	}
	return key
}

func TestLweKeyGen(t *testing.T) {
	assert := assert.New(t)

	for _, params := range allParams {
		key := &LweKey{Params: params}
		LweKeyGen(key)
		assert.Equal(params, key.Params, "Params and key.params should be the same.")
		n := key.Params.N
		s := key.Key
		//verify that the key is binary and kind-of random
		var count int = 0
		for i := int(0); i < n; i++ {
			assert.True(s[i] == 0 || s[i] == 1, "Key values should be 0 or 1.")
			count += s[i]
		}
		assert.LessOrEqual(count, n-20)
		assert.GreaterOrEqual(count, int(20))
	}
}

func TestLweSymEncryptPhaseDecrypt(t *testing.T) {
	assert := assert.New(t)

	var nbSamples int = 10
	var M int = 8
	alpha := 1.0 / (10.0 * double(M))

	for _, key := range allKeys {
		params := key.Params
		samples := NewLweSampleArray(nbSamples, params)
		// fmt.Println(samples)

		//verify correctness of the decryption
		for trial := int(0); trial < nbSamples; trial++ {
			message := ModSwitchToTorus(trial, M)
			LweSymEncrypt(samples[trial], message, alpha, key)
			phase := LwePhase(samples[trial], key)
			decrypt := LweSymDecrypt(samples[trial], key, M)
			dmessage := TorusToDouble(message)
			dphase := TorusToDouble(phase)
			assert.Equal(message, decrypt)
			assert.LessOrEqual(absfrac(dmessage-dphase), 10.*alpha)
			assert.Equal(alpha*alpha, samples[trial].CurrentVariance)
		}
		//verify that samples are random enough (all coordinates different)
		n := params.N
		for i := int(0); i < n; i++ {
			testset := make(map[Torus]bool)
			//set < Torus > testset
			for trial := int(0); trial < nbSamples; trial++ {
				testset[samples[trial].A[i]] = true
				//testset.insert(samples[trial].A[i])
			}
			assert.GreaterOrEqual(float32(len(testset)), 0.9*float32(nbSamples))
		}
	}

}

// fills a LweSample with random Torus
func fillRandom(result *LweSample, params *LweParams) {
	n := params.N
	for i := int(0); i < n; i++ {
		result.A[i] = UniformTorusDist()
	}
	result.B = UniformTorusDist()
	result.CurrentVariance = 0.2
}

//Arithmetic operations on Lwe samples
// result = (0,0)
func TestLweClear(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allKeys {
		params := key.Params
		n := params.N
		sample := NewLweSample(params)
		fillRandom(sample, params)
		LweClear(sample, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(0, sample.A[i])
		}
		assert.EqualValues(0, sample.B)
		assert.EqualValues(0., sample.CurrentVariance)
	}
}

// result = (0,mu)
func TestLweNoiselessTrivial(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allKeys {
		message := UniformTorus32Dist()
		params := key.Params
		n := params.N
		sample := NewLweSample(params)
		fillRandom(sample, params)
		LweNoiselessTrivial(sample, message, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(0, sample.A[i])
		}
		assert.EqualValues(message, sample.B)
		assert.EqualValues(0., sample.CurrentVariance)
	}
}

// copy a LweSample
func copySample(result *LweSample, sample *LweSample, params *LweParams) {
	n := params.N
	for i := int(0); i < n; i++ {
		result.A[i] = sample.A[i]
	}
	result.B = sample.B
	result.CurrentVariance = sample.CurrentVariance
}

func TestLweAddTo(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allKeys {
		params := key.Params
		n := params.N
		a := NewLweSample(params)
		b := NewLweSample(params)
		acopy := NewLweSample(params)
		fillRandom(a, params)
		fillRandom(b, params)
		copySample(acopy, a, params)
		LweAddTo(a, b, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(acopy.A[i]+b.A[i], a.A[i])
		}
		assert.EqualValues(acopy.B+b.B, a.B)
		assert.EqualValues(acopy.CurrentVariance+b.CurrentVariance, a.CurrentVariance)
	}
}

// result = result - sample
func TestLweSubTo(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allKeys {
		params := key.Params
		n := params.N
		a := NewLweSample(params)
		b := NewLweSample(params)
		acopy := NewLweSample(params)
		fillRandom(a, params)
		fillRandom(b, params)
		copySample(acopy, a, params)
		LweSubTo(a, b, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(acopy.A[i]-b.A[i], a.A[i])
		}
		assert.EqualValues(acopy.B-b.B, a.B)
		assert.EqualValues(acopy.CurrentVariance+b.CurrentVariance, a.CurrentVariance)
	}
}

func TestLweAddMulTo(t *testing.T) {
	assert := assert.New(t)
	const p int = 3
	for _, key := range allKeys {
		params := key.Params
		n := params.N
		a := NewLweSample(params)
		b := NewLweSample(params)
		acopy := NewLweSample(params)
		fillRandom(a, params)
		fillRandom(b, params)
		copySample(acopy, a, params)
		LweAddMulTo(a, p, b, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(acopy.A[i]+p*b.A[i], a.A[i])
		}
		assert.EqualValues(acopy.B+p*b.B, a.B)
		assert.EqualValues(acopy.CurrentVariance+float64(p)*float64(p)*b.CurrentVariance, a.CurrentVariance)
	}
}

func TestLweSubMulTo(t *testing.T) {
	assert := assert.New(t)
	const p int = 3
	for _, key := range allKeys {
		params := key.Params
		n := params.N
		a := NewLweSample(params)
		b := NewLweSample(params)
		acopy := NewLweSample(params)
		fillRandom(a, params)
		fillRandom(b, params)
		copySample(acopy, a, params)
		LweSubMulTo(a, p, b, params)
		for i := int(0); i < n; i++ {
			assert.EqualValues(acopy.A[i]-p*b.A[i], a.A[i])
		}
		assert.EqualValues(acopy.B-p*b.B, a.B)
		assert.EqualValues(acopy.CurrentVariance+float64(p)*float64(p)*b.CurrentVariance, a.CurrentVariance)
	}
}
