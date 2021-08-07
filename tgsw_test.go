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
	fakeUid         int32
	message         []TorusPolynomial
	currentVariance double
}

// construct
func NewFakeTLweFromTLweSample(sample *TLweSample) *FakeTLwe {
	return &FakeTLwe{
		fakeUid:         FAKE_TLWE_UID,
		message:         sample.A, //NewTorusPolynomial(N),
		currentVariance: 0.,
	}
}

// Fake TGSW structure
type FakeTGsw struct {
	//TODO: parallelization
	fakeUid         int64
	message         *IntPolynomial
	currentVariance double
}

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
	seed := NewFakeTLweFromTLweSample(sample)
	for i := int32(0); i < kpl; i++ {
		for j := int32(0); j < N; j++ {
			result[i].Coefs[j] = (i+3*j+seed.message[0].CoefsT[j])%25 - 12
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

func TestTGswSymEncryptInt(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allTGswKeys {
		N := key.params.TlweParams.N
		s := NewTGswSample(key.params)

		mess := rand.Int31()%1000 - 500
		var alpha double = 3.14 // valeur pseudo aleatoire fixé

		TGswSymEncryptInt(s, mess, alpha, key)

		fs := NewFakeTGsw(s, key.params)
		assert.EqualValues(fs.message.Coefs[0], mess)
		for j := int32(1); j < N; j++ {
			assert.EqualValues(fs.message.Coefs[j], 0)
		}
		assert.EqualValues(fs.currentVariance, alpha*alpha)
	}
}

func TestTGswClear(t *testing.T) {
	assert := assert.New(t)
	for _, param := range allTGswParams {
		s := NewTGswSample(param)
		kpl := param.Kpl
		zeroPol := NewTorusPolynomial(param.TlweParams.N)

		torusPolynomialClear(zeroPol)
		TGswClear(s, param)
		for i := int32(0); i < kpl; i++ {
			si := NewFakeTLweFromTLweSample(&s.AllSample[i])
			assert.EqualValues(torusPolynomialNormInftyDist(&si.message[0], zeroPol), 0)
		}
	}
}

// Test direct Result*H donne le bon resultat
// sample: TLweSample composed by k+1 torus polynomials, each with N coefficients
// result: int32 polynomial with Nl(k+1) coefficients
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
			}
		}
	}
}

//this function will create a fixed (random-looking) TGsw sample
func fullyRandomTGsw(result *TGswSample, alpha double, params *TGswParams) {
	for _, v := range result.AllSample {
		for _, vv := range v.A {
			torusPolynomialUniform(&vv)
		}
	}
}

func TestTGswAddH(t *testing.T) {
	assert := assert.New(t)
	for _, params := range allTGswParams {
		s := NewTGswSample(params)
		stemp := NewTGswSample(params)
		kpl := params.Kpl
		l := params.l
		k := params.TlweParams.K
		N := params.TlweParams.N
		h := params.H
		alpha := 4.2 // valeur pseudo aleatoire fixé

		// make a full random TGSW
		fullyRandomTGsw(s, alpha, params)

		// copy s to stemp
		for i := int32(0); i < kpl; i++ {
			TLweCopy(&stemp.AllSample[i], &s.AllSample[i], params.TlweParams)
		}

		TGswAddH(s, params)

		//verify all coefficients
		for bloc := int32(0); bloc <= k; bloc++ {
			for i := int32(0); i < l; i++ {
				assert.EqualValues(s.BlocSample[bloc][i].CurrentVariance, stemp.BlocSample[bloc][i].CurrentVariance)
				for u := int32(0); u <= k; u++ {
					//verify that pol[bloc][i][u]=initial[bloc][i][u]+(bloc==u?hi:0)
					newpol := &s.BlocSample[bloc][i].A[u]
					oldpol := &stemp.BlocSample[bloc][i].A[u]
					var check int32 = 0
					if bloc == u {
						check = h[i]
					}
					assert.EqualValues(newpol.CoefsT[0], oldpol.CoefsT[0]+check)
					for j := int32(1); j < N; j++ {
						assert.EqualValues(newpol.CoefsT[j], oldpol.CoefsT[j])
					}
				}
			}
		}
	}
}

func TestTGswAddMuH(t *testing.T) {
	assert := assert.New(t)
	for _, params := range allTGswParams {
		s := NewTGswSample(params)
		stemp := NewTGswSample(params)
		kpl := params.Kpl
		l := params.l
		k := params.TlweParams.K
		N := params.TlweParams.N
		h := params.H
		alpha := 4.2 // valeur pseudo aleatoire fixé
		mess := newRandomIntPolynomial(N)

		// make a full random TGSW
		fullyRandomTGsw(s, alpha, params)

		// copy s to stemp
		for i := int32(0); i < kpl; i++ {
			TLweCopy(&stemp.AllSample[i], &s.AllSample[i], params.TlweParams)
		}

		TGswAddMuH(s, mess, params)

		//verify all coefficients
		for bloc := int32(0); bloc <= k; bloc++ {
			for i := int32(0); i < l; i++ {
				assert.EqualValues(s.BlocSample[bloc][i].CurrentVariance, stemp.BlocSample[bloc][i].CurrentVariance)
				for u := int32(0); u <= k; u++ {
					//verify that pol[bloc][i][u]=initial[bloc][i][u]+(bloc==u?hi*mess:0)
					newpol := &s.BlocSample[bloc][i].A[u]
					oldpol := &stemp.BlocSample[bloc][i].A[u]
					if bloc == u {
						for j := int32(0); j < N; j++ {
							assert.EqualValues(newpol.CoefsT[j], oldpol.CoefsT[j]+h[i]*mess.Coefs[j])
						}
					} else {
						for j := int32(0); j < N; j++ {
							assert.EqualValues(newpol.CoefsT[j], oldpol.CoefsT[j])
						}
					}
				}
			}
		}
	}
}

func TestAddMuIntH(t *testing.T) {
	assert := assert.New(t)
	for _, params := range allTGswParams {
		s := NewTGswSample(params)
		stemp := NewTGswSample(params)
		kpl := params.Kpl
		l := params.l
		k := params.TlweParams.K
		N := params.TlweParams.N
		h := params.H
		alpha := 4.2 // valeur pseudo aleatoire fixé
		mess := rand.Int31()*2345 - 1234

		// make a full random TGSW
		fullyRandomTGsw(s, alpha, params)

		// copy s to stemp
		for i := int32(0); i < kpl; i++ {
			TLweCopy(&stemp.AllSample[i], &s.AllSample[i], params.TlweParams)
		}

		TGswAddMuIntH(s, mess, params)

		//verify all coefficients
		for bloc := int32(0); bloc <= k; bloc++ {
			for i := int32(0); i < l; i++ {
				assert.EqualValues(s.BlocSample[bloc][i].CurrentVariance, stemp.BlocSample[bloc][i].CurrentVariance)
				for u := int32(0); u <= k; u++ {
					//verify that pol[bloc][i][u]=initial[bloc][i][u]+(bloc==u?hi*mess:0)
					newpol := &s.BlocSample[bloc][i].A[u]
					oldpol := &stemp.BlocSample[bloc][i].A[u]
					var check int32 = 0
					if bloc == u {
						check = h[i] * mess
					}
					assert.EqualValues(newpol.CoefsT[0], oldpol.CoefsT[0]+check)
					for j := int32(1); j < N; j++ {
						assert.EqualValues(newpol.CoefsT[j], oldpol.CoefsT[j])
					}
				}
			}
		}

	}
}

func TestTGswEncryptZero(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allTGswKeys {
		kpl := key.params.Kpl
		s := NewTGswSample(key.params)
		alpha := 4.2 // valeur pseudo aleatoire fixé

		// Zero polynomial
		zeroPol := NewTorusPolynomial(key.params.TlweParams.N)
		torusPolynomialClear(zeroPol)

		TGswEncryptZero(s, alpha, key)
		for i := int32(0); i < kpl; i++ {
			//FakeTLwe * si = fake(&s.AllSample[i])
			assert.EqualValues(torusPolynomialNormInftyDist(s.AllSample[i].B(), zeroPol), 0)
			assert.EqualValues(s.AllSample[i].CurrentVariance, alpha*alpha)
		}
	}
}

func TestTGswTorus32PolynomialDecompH(t *testing.T) {
	assert := assert.New(t)
	for _, param := range allTGswParams {
		N := param.TlweParams.N
		l := param.l
		Bgbit := param.Bgbit
		h := param.H

		// compute the tolerance
		var toler int32 = 0
		if Bgbit*l < 32 {
			toler = 1 << (32 - Bgbit*l)
		}
		// fmt.Printf("%d,%d,%d\n", Bgbit, l, toler)

		result := NewIntPolynomialArray(int(l), N)
		sample := NewTorusPolynomial(N)
		torusPolynomialUniform(sample)

		TGswTorus32PolynomialDecompH(result, sample, param)

		for i := int32(0); i < N; i++ {
			// recomposition
			var test Torus32 = 0
			for j := int32(0); j < l; j++ {
				test += result[j].Coefs[i] * h[j]
			}
			assert.LessOrEqual(Abs(test-sample.CoefsT[i]), toler)
		}
	}
}

func TestTGswExternProduct(t *testing.T) {
	assert := assert.New(t)
	for _, params := range allTGswParams {
		N := params.TlweParams.N
		kpl := params.Kpl

		sample := NewTGswSample(params)

		result := NewTLweSample(params.TlweParams)
		b := NewTLweSample(params.TlweParams)

		fresult := NewFakeTLweFromTLweSample(result) //fake(result);
		fb := NewFakeTLweFromTLweSample(b)
		fsamplerows := sample.AllSample //NewFakeTLweFromTLweSample(sample.AllSample) //fake(sample.allSample);
		alpha := 4.2                    // valeur pseudo aleatoire fixé

		fullyRandomTGsw(sample, alpha, params)
		torusPolynomialUniform(&fb.message[0])

		decomp := NewIntPolynomialArray(int(kpl), N)
		TGswTLweDecompH(decomp, b, params)
		expectedRes := NewTorusPolynomial(N)
		tmp := NewTorusPolynomial(N)

		torusPolynomialClear(expectedRes)
		for i := int32(0); i < kpl; i++ {
			torusPolynomialMultKaratsuba(tmp, &decomp[i], fsamplerows[i].B())
			TorusPolynomialAddTo(expectedRes, tmp)
		}
		TGswExternProduct(result, sample, b, params)
		assert.EqualValues(torusPolynomialNormInftyDist(&fresult.message[0], expectedRes), 0)
	}
}

func TestTGswMulByXaiMinusOne(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allTGswKeys {
		kpl := key.params.Kpl
		N := key.params.TlweParams.N
		for ai := int32(0); ai < 2*N; ai += 235 {
			res := NewTGswSample(key.params)
			bk := NewTGswSample(key.params)
			alpha := 4.2 // valeur pseudo aleatoire fixé
			poly := NewTorusPolynomial(N)

			//generate all rows at random
			fresrows := res.AllSample
			fbkrows := bk.AllSample
			fullyRandomTGsw(bk, alpha, key.params)
			TGswMulByXaiMinusOne(res, ai, bk, key.params)
			for i := int32(0); i < kpl; i++ {
				TorusPolynomialMulByXaiMinusOne(poly, ai, fbkrows[i].B())
				assert.EqualValues(torusPolynomialNormInftyDist(fresrows[i].B(), poly), 0)
				var check float64 = 2
				if ai == 0 {
					check = 1
				}
				assert.EqualValues(fresrows[i].CurrentVariance, check*fbkrows[i].CurrentVariance)
			}
		}
	}
}

func TestTGswExternMulToTLwe(t *testing.T) {
	assert := assert.New(t)
	for _, key := range allTGswKeys {
		params := key.params
		N := params.TlweParams.N
		kpl := params.Kpl

		sample := NewTGswSample(params)
		accum := NewTLweSample(params.TlweParams)
		faccum := NewFakeTLweFromTLweSample(accum)
		fsamplerows := sample.AllSample
		alpha := 4.2 // valeur pseudo aleatoire fixé

		fullyRandomTGsw(sample, alpha, params)
		torusPolynomialUniform(&faccum.message[0])

		decomp := NewIntPolynomialArray(int(kpl), N)
		TGswTLweDecompH(decomp, accum, params)
		expectedRes := NewTorusPolynomial(N)
		//tmp := NewTorusPolynomial(N)

		torusPolynomialClear(expectedRes)
		for i := int32(0); i < kpl; i++ {
			torusPolynomialAddMulRKaratsuba(expectedRes, &decomp[i], fsamplerows[i].B())
		}
		TGswExternMulToTLwe(accum, sample, params)
		assert.EqualValues(torusPolynomialNormInftyDist(&faccum.message[0], expectedRes), 0)
	}
}

func TestTGswNoiselessTrivial(t *testing.T) {
	assert := assert.New(t)
	for _, param := range allTGswParams {
		N := param.TlweParams.N
		res := NewTGswSample(param)
		mu := newRandomIntPolynomial(N)
		fres := NewFakeTGsw(res, param)

		TGswNoiselessTrivial(res, mu, param)

		for j := int32(0); j < N; j++ {
			assert.EqualValues(fres.message.Coefs[j], mu.Coefs[j])
		}
		assert.EqualValues(fres.currentVariance, 0.)
	}
}
