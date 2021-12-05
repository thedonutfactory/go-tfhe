package tfhe

import (
	"gonum.org/v1/gonum/stat/distuv"
)

type LweKeySwitchKey struct {
	N         int32      ///< length of the input key: s'
	T         int32      ///< decomposition length
	Basebit   int32      ///< log_2(base)
	Base      int32      ///< decomposition base: a power of 2
	OutParams *LweParams ///< params of the output key s
	//ks0Raw    []LweSample // array which contains all Lwe samples of size nlbase
	//ks1Raw    [][]*LweSample   // of size nl points to a ks0Raw array whose cells are spaced at base positions
	Ks [][][]*LweSample ///< the keyswitch elements: a n.l.base matrix
	// of size n points to ks1 an array whose cells are spaced by ell positions
}

func NewLweKeySwitchKey(n, t, basebit int, outParams *LweParams) *LweKeySwitchKey {

	base := int32(1 << basebit)
	ks := make([][][]*LweSample, n)
	for i := int32(0); i < n; i++ {
		ks[i] = make([][]*LweSample, t)
		for j := int32(0); j < t; j++ {
			ks[i][j] = make([]*LweSample, base)
			for k := int32(0); k < base; k++ {
				ks[i][j][k] = NewLweSample(outParams) //ks0Raw[c]
			}
		}
	}
	return &LweKeySwitchKey{
		N:         n,
		T:         t,
		Basebit:   basebit,
		Base:      base,
		OutParams: outParams,
		Ks:        ks,
	}
}

func NewLweKeySwitchKeyArray(size, n, t, basebit int, outParams *LweParams) (ksk []LweKeySwitchKey) {
	ksk = make([]LweKeySwitchKey, size)
	for i := int(0); i < size; i++ {
		ksk = append(ksk, *NewLweKeySwitchKey(n, t, basebit, outParams))
	}
	return
}

/*
Renormalization of KS
 * compute the error of the KS that has been generated and translate the ks to recenter the gaussian in 0
*/
func renormalizeKSkey(ks *LweKeySwitchKey, outKey *LweKey, inKey []int32) {
	n := int32(ks.N)
	basebit := ks.Basebit
	t := int32(ks.T)
	base := int32(1 << basebit)
	var err Torus32

	// compute the average error
	for i := int(0); i < n; i++ {
		for j := int(0); j < t; j++ {
			for h := int(1); h < base; h++ { // pas le terme en 0
				// compute the phase
				phase := LwePhase(ks.Ks[i][j][h], outKey)
				// compute the error
				x := (inKey[i] * h) * (1 << (32 - (j+1)*basebit))
				tempErr := phase - x
				// sum all errors
				err += tempErr
			}
		}
	}
	nb := n * t * (base - 1)
	err = DoubleToTorus(TorusToDouble(err) / double(nb))

	// relinearize
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < t; j++ {
			for h := int32(1); h < base; h++ { // pas le terme en 0
				ks.Ks[i][j][h].B -= err
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
	inKey []int, n, t, basebit int) {

	base := int(1 << basebit) // base=2 in [CGGI16]

	for i := int(0); i < n; i++ {
		for j := int(0); j < t; j++ {
			for k := int(0); k < base; k++ {
				x := (inKey[i] * k) * (1 << (32 - (j+1)*basebit))
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
	n, t, basebit int) {

	base := 1 << basebit
	precOffset := int(1 << (32 - (1 + basebit*t)))
	mask := int(base - 1)

	for i := int(0); i < n; i++ {
		aibar := ai[i] + precOffset
		for j := int(0); j < t; j++ {
			aij := (aibar >> (32 - (j+1)*basebit)) & mask
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

	n := result.N
	t := result.T
	basebit := result.Basebit
	base := int32(1 << basebit)
	alpha := outKey.Params.AlphaMin
	sizeks := n * t * (base - 1)
	var err double = 0

	dist := distuv.Uniform{
		Min: 0,
		Max: alpha,
	}
	// chose a random vector of gaussian noises
	noise := make([]double, sizeks)
	for i := int(0); i < sizeks; i++ {
		noise[i] = dist.Rand()
		err += noise[i]
	}
	// recenter the noises
	err = err / double(sizeks)
	for i := int(0); i < sizeks; i++ {
		noise[i] -= err
	}

	// generate the ks
	var index int = 0
	for i := int(0); i < n; i++ {
		for j := int(0); j < t; j++ {
			// term h=0 as trivial encryption of 0 (it will not be used in the KeySwitching)
			LweNoiselessTrivial(result.Ks[i][j][0], 0, outKey.Params)

			for h := int32(1); h < base; h++ { // pas le terme en 0
				mess := (inKey.Key[i] * h) * (1 << (32 - (j+1)*basebit))
				LweSymEncryptWithExternalNoise(result.Ks[i][j][h], mess, noise[index], alpha, outKey)
				index += 1
			}
		}
	}
}

//sample=(a',b')
func lweKeySwitch(ks *LweKeySwitchKey, sample *LweSample) *LweSample {
	params := ks.OutParams
	n := ks.N
	basebit := ks.Basebit
	t := ks.T
	result := NewLweSample(ks.OutParams)
	LweNoiselessTrivial(result, sample.B, params)
	lweKeySwitchTranslateFromArray(result, ks.Ks, params, sample.A, n, t, basebit)
	return result
}
