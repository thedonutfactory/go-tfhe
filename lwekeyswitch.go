package tfhe

import (
	"gonum.org/v1/gonum/stat/distuv"
)

type LweKeySwitchKey struct {
	n         int32      ///< length of the input key: s'
	t         int32      ///< decomposition length
	basebit   int32      ///< log_2(base)
	base      int32      ///< decomposition base: a power of 2
	outParams *LweParams ///< params of the output key s
	//ks0Raw    []LweSample // array which contains all Lwe samples of size nlbase
	//ks1Raw    [][]*LweSample   // of size nl points to a ks0_raw array whose cells are spaced at base positions
	ks [][][]*LweSample ///< the keyswitch elements: a n.l.base matrix
	// of size n points to ks1 an array whose cells are spaced by ell positions
}

func NewLweKeySwitchKey(n, t, basebit int32, outParams *LweParams) *LweKeySwitchKey {

	base := int32(1 << basebit)
	ks0_raw := NewLweSampleArray(n*t*base, outParams)

	//ks1_raw := make([][]*LweSample, n*t) //new LweSample*[n*t];

	ks := make([][][]*LweSample, n)

	/*
		for i := int32(0); i < n; i++ {
			ks[i] = make([][]*LweSample, n)
			for j := int32(0); j < t; j++ {
				ks[i][j] = make([]*LweSample, t)
				for h := int32(0); h < base; h++ {
					ks[i][j][h] = &ks0_raw[i*base]
				}
			}
		}
	*/

	// N = 300
	// t = 14
	// base = 4
	var c int = 0
	for i := int32(0); i < n; i++ {
		ks[i] = make([][]*LweSample, n)
		for j := int32(0); j < t; j++ {
			ks[i][j] = make([]*LweSample, t)
			for k := int32(0); k < base; k++ {
				ks[i][j][k] = ks0_raw[c]
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
		//ks0Raw:    ks0_raw,
		ks: ks,
	}
}

/**
 * LweKeySwitchKey constructor function

func NewLweKeySwitchKey(n, t, basebit int32, outParams *LweParams) *LweKeySwitchKey {
	base := int32(1 << basebit)
	ks0Raw := NewLweSampleArray(n*t*base, outParams)
	return &LweKeySwitchKey{n: n, t: t, basebit: basebit, outParams: outParams, ks0Raw: ks0Raw}
	// new(obj) LweKeySwitchKey(n,t,basebit,out_params, ks0_raw)
}
*/

func new_LweKeySwitchKey_array(size, n, t, basebit int32, outParams *LweParams) (ksk []LweKeySwitchKey) {
	//LweKeySwitchKey* obj = alloc_LweKeySwitchKey_array(nbelts);
	//init_LweKeySwitchKey_array(nbelts, obj, n,t,basebit,outParams);
	//return obj;

	ksk = make([]LweKeySwitchKey, size)
	for i := int32(0); i < size; i++ {
		ksk = append(ksk, *NewLweKeySwitchKey(n, t, basebit, outParams))
	}
	return
}

/*
Renormalization of KS
 * compute the error of the KS that has been generated and translate the ks to recenter the gaussian in 0
*/
func renormalizeKSkey(ks *LweKeySwitchKey, out_key *LweKey, in_key []int32) {
	n := int32(ks.n)
	basebit := ks.basebit
	t := int32(ks.t)
	base := int32(1 << basebit)

	//var phase, err Torus32
	//var temp_err Torus32
	//var err Torus32 = 0
	// double err_norm = 0;
	var err Torus32

	// compute the average error
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < t; j++ {
			for h := int32(1); h < base; h++ { // pas le terme en 0
				// compute the phase
				phase := LwePhase(ks.ks[i][j][h], out_key)
				// compute the error
				x := (in_key[i] * h) * (1 << (32 - (j+1)*basebit))
				temp_err := phase - x
				// sum all errors
				err += temp_err
			}
		}
	}
	nb := n * t * (base - 1)
	err = Dtot32(T32tod(err) / double(nb))

	// relinearize
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < t; j++ {
			for h := int32(1); h < base; h++ { // pas le terme en 0
				ks.ks[i][j][h].B -= err
			}
		}
	}
}

/**
 * fills the KeySwitching key array
 * @param result The (n x t x base) array of samples.
 *        result[i][j][k] encodes k.s[i]/base^(j+1)
 * @param out_key The LWE key to encode all the output samples
 * @param out_alpha The standard deviation of all output samples
 * @param in_key The (binary) input key
 * @param n The size of the input key
 * @param t The precision of the keyswitch (technically, 1/2.base^t)
 * @param basebit Log_2 of base
 */
func lweCreateKeySwitchKeyFromArray(result [][][]*LweSample,
	out_key *LweKey, out_alpha double,
	in_key []int32, n, t, basebit int32) {

	base := int32(1 << basebit) // base=2 in [CGGI16]

	for i := int32(0); i < n; i++ {
		for j := int32(0); j < t; j++ {
			for k := int32(0); k < base; k++ {
				x := (in_key[i] * k) * (1 << (32 - (j+1)*basebit))
				LweSymEncrypt(result[i][j][k], x, out_alpha, out_key)
				//fmt.Printf("i,j,k,ki,x,phase=%d,%d,%d,%d,%d,%d\n", i, j, k, in_key[i], x, LwePhase(result[i][j][k], out_key))
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
	ai []Torus32,
	n, t, basebit int32) {

	base := 1 << basebit
	prec_offset := int32(1 << (32 - (1 + basebit*t)))
	mask := int32(base - 1)

	for i := int32(0); i < n; i++ {
		aibar := ai[i] + prec_offset
		for j := int32(0); j < t; j++ {
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
func lweCreateKeySwitchKey(result *LweKeySwitchKey, in_key *LweKey, out_key *LweKey) {

	n := result.n
	t := result.t
	basebit := result.basebit
	base := int32(1 << basebit)
	alpha := out_key.params.alphaMin
	sizeks := n * t * (base - 1)
	//n_out = out_key.params.n

	var err double = 0

	dist := distuv.Uniform{
		Min: 0,
		Max: alpha,
	}
	// chose a random vector of gaussian noises
	noise := make([]double, sizeks)
	for i := int32(0); i < sizeks; i++ {
		//normal_distribution<double> distribution(0.,alpha);
		//noise[i] = distribution(generator)
		noise[i] = dist.Rand()
		err += noise[i]
	}
	// recenter the noises
	err = err / double(sizeks)
	for i := int32(0); i < sizeks; i++ {
		noise[i] -= err
	}

	// generate the ks
	var index int = 0
	for i := int32(0); i < n; i++ {
		for j := int32(0); j < t; j++ {

			// term h=0 as trivial encryption of 0 (it will not be used in the KeySwitching)
			LweNoiselessTrivial(result.ks[i][j][0], 0, out_key.params)
			//lweSymEncrypt(&result.ks[i][j][0],0,alpha,out_key)

			for h := int32(1); h < base; h++ { // pas le terme en 0
				/*
				   // noiseless encryption
				   result.ks[i][j][h].b = (in_key.key[i]*h)*(1<<(32-(j+1)*basebit))
				   for (int_t p = 0; p < n_out; ++p) {
				       result.ks[i][j][h].a[p] = uniformTorus32_distrib(generator)
				       result.ks[i][j][h].b += result.ks[i][j][h].a[p] * out_key.key[p]
				   }
				   // add the noise
				   result.ks[i][j][h].b += Dtot32(noise[index])
				*/
				mess := (in_key.key[i] * h) * (1 << (32 - (j+1)*basebit))
				LweSymEncryptWithExternalNoise(result.ks[i][j][h], mess, noise[index], alpha, out_key)
				index += 1
			}
		}
	}

	// delete[] noise;
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
		sample.A, n, t, basebit)
}
