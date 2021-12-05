package tfhe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	params500_0    = NewLweParams(500, 0., 1.)
	params500_1em5 = NewLweParams(500, 1e-5, 1.)
	keyParams500   = NewLweKey(params500_0)
)

/**
* fills the KeySwitching key array
* @param result The (n x t x base) array of samples.
*        result[i][j][k] encodes k.s[i]/base^(j+1)
* @param outKey The LWE key to encode all the output samples
* @param outAlpha The standard deviation of all output samples
* @param inKey The (binary) input key
* @param n The size of the input key
* @param t The precision of the keyswitch (technically, 1/2.base^t)
* @param basebit Log_2 of base
 */
func TestLweCreateKeySwitchKeyFromArray(tt *testing.T) {
	assert := assert.New(tt)
	// Mock LweSymEncrypt for this test
	old := LweSymEncrypt
	defer func() { LweSymEncrypt = old }()
	LweSymEncrypt = func(
		result *LweSample,
		message Torus,
		alpha double,
		key *LweKey) {
		LweNoiselessTrivial(result, message, key.Params)
		result.CurrentVariance = alpha * alpha
	}

	test := NewLweKeySwitchKey(300, 14, 2, params500_1em5)
	alpha := 1e-5
	N := test.N
	t := test.T
	basebit := test.Basebit
	base := test.Base
	inKey := make([]int32, N)
	for i := int32(0); i < N; i++ {
		if UniformTorus32Dist()%2 == 0 {
			inKey[i] = 1
		} else {
			inKey[i] = 0
		}
	}
	lweCreateKeySwitchKeyFromArray(test.Ks, key500, alpha, inKey, N, t, basebit)
	for i := int32(0); i < N; i++ {
		for j := int32(0); j < t; j++ {
			for k := int32(0); k < base; k++ {
				ksIjk := test.Ks[i][j][k]
				assert.EqualValues(alpha*alpha, ksIjk.CurrentVariance)
				//fmt.Printf("%d, %d\n", k*inKey[i]*1<<(32-(j+1)*basebit), ksIjk.B)
				assert.EqualValues(k*inKey[i]*1<<(32-(j+1)*basebit), ksIjk.B)
			}
		}
	}
}
