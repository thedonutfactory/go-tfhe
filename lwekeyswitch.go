package tfhe

import (
	"gonum.org/v1/gonum/stat/distuv"
)

type LweKeySwitchKey struct {
	n         int        ///< length of the input key: s'
	t         int        ///< decomposition length
	basebit   int        ///< log_2(base)
	base      int        ///< decomposition base: a power of 2
	outParams *LweParams ///< params of the output key s
	//ks0Raw    []LweSample // array which contains all Lwe samples of size nlbase
	//ks1Raw    [][]*LweSample   // of size nl points to a ks0Raw array whose cells are spaced at base positions
	ks [][][]*LweSample ///< the keyswitch elements: a n.l.base matrix
	// of size n points to ks1 an array whose cells are spaced by ell positions
}

func NewLweKeySwitchKey(n, t, basebit int, outParams *LweParams) *LweKeySwitchKey {

	base := 1 << basebit
	ks0Raw := NewLweSampleArray(n*t*base, outParams)
	ks := make([][][]*LweSample, n)
	var c int = 0
	for i := 0; i < n; i++ {
		ks[i] = make([][]*LweSample, n)
		for j := 0; j < t; j++ {
			ks[i][j] = make([]*LweSample, t)
			for k := 0; k < base; k++ {
				ks[i][j][k] = ks0Raw[c]
				c++
			}
		}
	}
	return &LweKeySwitchKey{
		n:         n,
		t:         t,
		basebit:   basebit,
		base:      base,
		outParams: outParams,
		ks:        ks,
	}
}

func NewLweKeySwitchKeyArray(size, n, t, basebit int, outParams *LweParams) (ksk []LweKeySwitchKey) {
	ksk = make([]LweKeySwitchKey, size)
	for i := 0; i < size; i++ {
		ksk = append(ksk, *NewLweKeySwitchKey(n, t, basebit, outParams))
	}
	return
}

/*
Renormalization of KS
 * compute the error of the KS that has been generated and translate the ks to recenter the gaussian in 0
*/
func renormalizeKSkey(ks *LweKeySwitchKey, outKey *LweKey, inKey []int64) {
	n := ks.n
	basebit := ks.basebit
	t := ks.t
	base := 1 << basebit
	var err Torus

	// compute the average error
	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			for h := 1; h < base; h++ { // pas le terme en 0
				// compute the phase
				phase := LwePhase(ks.ks[i][j][h], outKey)
				// compute the error
				x := (inKey[i] * int64(h)) * (1 << (64 - (j+1)*basebit))
				tempErr := phase - x
				// sum all errors
				err += tempErr
			}
		}
	}
	nb := n * t * (base - 1)
	err = Dtot32(T32tod(err) / double(nb))

	// relinearize
	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			for h := 1; h < base; h++ { // pas le terme en 0
				ks.ks[i][j][h].B -= err
			}
		}
	}
}

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
func lweCreateKeySwitchKeyFromArray(result [][][]*LweSample,
	outKey *LweKey, outAlpha double,
	inKey []int64, n, t, basebit int) {

	base := 1 << basebit // base=2 in [CGGI16]

	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			for k := 0; k < base; k++ {
				x := (inKey[i] * int64(k)) * (1 << (64 - (j+1)*basebit))
				LweSymEncrypt(result[i][j][k], x, outAlpha, outKey)
				//fmt.Printf("i,j,k,ki,x,phase=%d,%d,%d,%d,%d,%d\n", i, j, k, inKey[i], x, LwePhase(result[i][j][k], outKey))
			}
		}
	}
}

/**
 * translates the message of the result sample by -sum(a[i].s[i]) where s is the secret
 * embedded in ks.
 * @param result the LWE sample to translate by -sum(ai.si).
 * @param ks The (n x t x base) key switching key
 *        ks[i][j][k] encodes k.s[i]/base^(j+1)
 * @param params The common LWE parameters of ks and result
 * @param ai The input torus array
 * @param n The size of the input key
 * @param t The precision of the keyswitch (technically, 1/2.base^t)
 * @param basebit Log_2 of base
 */
func lweKeySwitchTranslateFromArray(result *LweSample,
	ks [][][]*LweSample, params *LweParams,
	ai []Torus,
	n, t, basebit int64) {

	base := 1 << basebit
	precOffset := int64(1 << (64 - (1 + basebit*t)))
	mask := int64(base - 1)

	for i := int64(0); i < n; i++ {
		aibar := int64(ai[i]) + precOffset
		for j := int64(0); j < t; j++ {
			aij := (aibar >> (64 - (j+1)*basebit)) & mask
			if aij != 0 {
				LweSubTo(result, ks[i][j][aij], params)
			}
		}
	}
}

/*
Create the key switching key: normalize the error in the beginning
 * chose a random vector of gaussian noises (same size as ks)
 * recenter the noises
 * generate the ks by creating noiseless encryprions and then add the noise
*/
func lweCreateKeySwitchKey(result *LweKeySwitchKey, inKey *LweKey, outKey *LweKey) {

	n := result.n
	t := result.t
	basebit := result.basebit
	base := 1 << basebit
	alpha := outKey.params.alphaMin
	sizeks := n * t * (base - 1)
	var err double = 0

	dist := distuv.Uniform{
		Min: 0,
		Max: alpha,
	}
	// chose a random vector of gaussian noises
	noise := make([]double, sizeks)
	for i := 0; i < sizeks; i++ {
		noise[i] = dist.Rand()
		err += noise[i]
	}
	// recenter the noises
	err = err / double(sizeks)
	for i := 0; i < sizeks; i++ {
		noise[i] -= err
	}

	// generate the ks
	var index int = 0
	for i := 0; i < n; i++ {
		for j := 0; j < t; j++ {
			// term h=0 as trivial encryption of 0 (it will not be used in the KeySwitching)
			LweNoiselessTrivial(result.ks[i][j][0], 0, outKey.params)

			for h := 1; h < base; h++ { // pas le terme en 0
				mess := (inKey.key[i] * int64(h)) * (1 << (64 - (j+1)*basebit))
				LweSymEncryptWithExternalNoise(result.ks[i][j][h], mess, noise[index], alpha, outKey)
				index += 1
			}
		}
	}
}

//sample=(a',b')
func lweKeySwitch(result *LweSample, ks *LweKeySwitchKey, sample *LweSample) {
	params := ks.outParams
	n := ks.n
	basebit := ks.basebit
	t := ks.t

	LweNoiselessTrivial(result, sample.B, params)
	lweKeySwitchTranslateFromArray(result,
		ks.ks, params,
		sample.A, int64(n), int64(t), int64(basebit))
}
