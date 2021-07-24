package tfhe

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const FAKE_TLWE_UID int32 = 542354312          // precaution: distinguish fakes from trues
const FAKE_TGSW_UID int64 = 123444802642375465 // precaution: do not confuse fakes with trues

// Fake TLWE structure
type FakeTLwe struct {
	fake_uid        int32
	message         TorusPolynomial
	currentVariance double
	//this padding is here to make sure that FakeTLwe and TLweSample have the same size
	//char unused_padding[sizeof(TLweSample) - sizeof(int64_t) - sizeof(TorusPolynomial *) - sizeof(double)];
}

// construct
func NewFakeTLwe(N int32) *FakeTLwe {
	return &FakeTLwe{
		fake_uid:        FAKE_TLWE_UID,
		message:         *NewTorusPolynomial(N),
		currentVariance: 0.,
	}
}

func NewFakeTLweFromTLweSample(sample *TLweSample) *FakeTLwe {
	return &FakeTLwe{
		fake_uid:        FAKE_TLWE_UID,
		message:         *sample.B, //NewTorusPolynomial(N),
		currentVariance: 0.,
	}
}

// Fake TGSW structure
type FakeTGsw struct {
	//TODO: parallelization
	fakeUid         int64
	message         *IntPolynomial
	currentVariance double

	//this padding is here to make sure that FakeTLwe and TLweSample have the same size
	//char unused_padding[sizeof(TGswSample) - sizeof(int64_t) - sizeof(IntPolynomial *) - sizeof(double)];
}

/*
func NewFakeTGsw(N int32) *FakeTGsw {
	return &FakeTGsw{
		fakeUid:         FAKE_TGSW_UID,
		message:         NewIntPolynomial(N),
		currentVariance: 0,
	}
}



inline void fake_init_TGswSample(TGswSample *ptr, const TGswParams *params) {
        int32_t N = params->tlwe_params->N;
        FakeTGsw *arr = (FakeTGsw *) ptr;
        new(arr) FakeTGsw(N);
    }
*/

func NewFakeTGsw(sample *TGswSample, params *TGswParams) *FakeTGsw {
	var alpha double = 4.2 // valeur pseudo aleatoire fixé
	return &FakeTGsw{
		fakeUid:         FAKE_TGSW_UID,
		message:         NewIntPolynomial(params.TlweParams.N),
		currentVariance: alpha * alpha,
	}
}

func (fake *FakeTGsw) setMessageVariance(mess int32, variance double) {
	intPolynomialClear(fake.message)
	fake.message.Coefs[0] = mess
	fake.currentVariance = variance
}

//this function creates a fixed (random-looking) result,
//whose seed is the sample
func FakeTGswTLweDecompH(result []IntPolynomial, sample *TLweSample, params *TGswParams) {
	kpl := params.Kpl
	N := params.TlweParams.N
	//const FakeTLwe *seed = fake(sample);
	seed := NewFakeTLweFromTLweSample(sample)
	for i := int32(0); i < kpl; i++ {
		for j := int32(0); j < N; j++ {
			result[i].Coefs[j] = (i+3*j+seed.message.CoefsT[j])%25 - 12
		}
	}
}

// we use the function rand because in the "const static" context the uniformly random generator doesn't work!
func newRandomTGswKey(params *TGswParams) *TGswKey {
	key := NewTGswKey(params)
	N := params.TlweParams.N
	k := params.TlweParams.K
	for i := int32(0); i < k; i++ {
		for j := int32(0); j < N; j++ {
			key.key[i].Coefs[j] = rand.Int31() % 2
		}
	}
	return key
}

// we use the function rand because in the "const static" context the uniformly random generator doesn't work!
func newRandomIntPolynomial(N int32) *IntPolynomial {
	poly := NewIntPolynomial(N)

	for i := int32(0); i < N; i++ {
		poly.Coefs[i] = rand.Int31()%10 - 5
	}
	return poly
}

/*
* Parameters and keys (for N=512,1024,2048 and k=1,2)
 */
var (
	tGswParams512_1  = NewTGswParams(4, 8, NewTLweParams(512, 1, 0., 1.))
	tGswParams512_2  = NewTGswParams(3, 10, NewTLweParams(512, 2, 0., 1.))
	tGswParams1024_1 = NewTGswParams(3, 10, NewTLweParams(1024, 1, 0., 1.))
	tGswParams1024_2 = NewTGswParams(4, 8, NewTLweParams(1024, 2, 0., 1.))
	tGswParams2048_1 = NewTGswParams(4, 8, NewTLweParams(2048, 1, 0., 1.))
	tGswParams2048_2 = NewTGswParams(3, 10, NewTLweParams(2048, 2, 0., 1.))

	tGswKey512_1  = newRandomTGswKey(tGswParams512_1)
	tGswKey512_2  = newRandomTGswKey(tGswParams512_2)
	tGswKey1024_1 = newRandomTGswKey(tGswParams1024_1)
	tGswKey1024_2 = newRandomTGswKey(tGswParams1024_2)
	tGswKey2048_1 = newRandomTGswKey(tGswParams2048_1)
	tGswKey2048_2 = newRandomTGswKey(tGswParams2048_2)

	allTGswParams     = []*TGswParams{tGswParams512_1, tGswParams512_2, tGswParams1024_1, tGswParams1024_2, tGswParams2048_1, tGswParams2048_2}
	allTGswParams1024 = []*TGswParams{tGswParams512_1, tGswParams512_2, tGswParams1024_1, tGswParams1024_2, tGswParams2048_1, tGswParams2048_2}

	allTGswKeys     = []*TGswKey{tGswKey512_1, tGswKey512_2, tGswKey1024_1, tGswKey1024_2, tGswKey2048_1, tGswKey2048_2}
	allTGswKeys1024 = []*TGswKey{tGswKey1024_1, tGswKey1024_2}
)

/*
*  Testing the function tGswKeyGen
* EXPORT void tLweKeyGen(TLweKey* result);
*
* This function generates a random TLwe key for the given parameters
* The TLwe key for the result must be allocated and initialized
* (this means that the parameters are already in the result)
 */
//EXPORT void tGswKeyGen(TGswKey* result);
//TEST_F(TGswTest, tGswKeyGen) {
func TestTGswKeyGen(t *testing.T) {
	assert := assert.New(t)
	for _, param := range allTGswParams {
		key := NewTGswKey(param)
		k := param.TlweParams.K
		N := param.TlweParams.N

		TGswKeyGen(key)
		for i := int32(0); i < k; i++ {
			for j := int32(0); j < N; j++ {
				assert.True(key.key[i].Coefs[j] == 0 || key.key[i].Coefs[j] == 1)
			}
		}
	}
}

//EXPORT void tGswSymEncrypt(TGswSample* result, const IntPolynomial* message, double alpha, const TGswKey* key);
//TEST_F(TGswFakeTest, tGswSymEncrypt) {
func TestTGswSymEncrypt(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allTGswKeys {
		N := key.params.TlweParams.N
		s := NewTGswSample(key.params)
		mess := newRandomIntPolynomial(N)
		var alpha double = 4.2 // valeur pseudo aleatoire fixé

		TGswSymEncrypt(s, mess, alpha, key)
		fs := NewFakeTGsw(s, key.params)
		for j := int32(0); j < N; j++ {
			assert.EqualValues(fs.message.Coefs[j], mess.Coefs[j])
		}
		assert.EqualValues(fs.currentVariance, alpha*alpha)
	}
}

//EXPORT void tGswTLweDecompH(IntPolynomial* result, const TLweSample* sample,const TGswParams* params);
// Test direct Result*H donne le bon resultat
// sample: TLweSample composed by k+1 torus polynomials, each with N coefficients
// result: int32_t polynomial with Nl(k+1) coefficients
func TestTGswTLweDecompH(t *testing.T) {
	assert := assert.New(t)
	for _, param := range allTGswParams {
		N := param.TlweParams.N
		k := param.TlweParams.K
		Bgbit := param.Bgbit
		l := param.l
		kpl := param.Kpl
		h := param.H

		//compute the tolerance
		var toler int32 = 0
		if Bgbit*l < 32 {
			toler = 1 << (32 - Bgbit*l)
		}
		//printf("%d,%d,%d\n",Bgbit,l,toler)

		result := NewIntPolynomialArray(int(kpl), N)
		sample := NewTLweSample(param.TlweParams)

		// sample randomly generated
		for bloc := int32(0); bloc <= k; bloc++ {
			torusPolynomialUniform(&sample.A[bloc])
		}

		TGswTLweDecompH(result, sample, param)

		for bloc := int32(0); bloc <= k; bloc++ {
			for i := int32(0); i < N; i++ {
				var test Torus32 = 0
				for j := int32(0); j < l; j++ {
					test += result[bloc*l+j].Coefs[i] * h[j]
				}
				assert.LessOrEqual(Abs(test-sample.A[bloc].CoefsT[i]), toler)
				//ASSERT_LE(abs(test-sample.a[bloc].coefsT[i]), toler) //exact or approx decomposition
			}
		}

	}
}
