package tfhe

import (
	"testing"

	set "github.com/badgerodon/collections/set"
	"github.com/stretchr/testify/assert"
)

var (
	dimensions               = [...]int{500, 750, 1024, 2000}
	powers_of_two_dimensions = [...]int{512, 1024, 2048}
)

func TestTorusPolynomialUniform(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	const NB_TRIALS int = 10
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(NB_TRIALS, int32(N))
		for i := 0; i < NB_TRIALS; i++ {
			torusPolynomialUniform(&pols[i])
		}
		for j := 0; j < N; j++ {
			testset := set.New()
			for i := 0; i < NB_TRIALS; i++ {
				testset.Insert(pols[i].CoefsT[j])
			}
			assert.GreaterOrEqual(float64(testset.Len()), 0.9*float64(NB_TRIALS))
		}
	}
}

//  TorusPolynomial = 0

func TestTorusPolynomialClear(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pol := NewTorusPolynomial(int32(N))
		torusPolynomialUniform(pol)
		torusPolynomialClear(pol)
		for j := 0; j < N; j++ {
			assert.EqualValues(0, pol.CoefsT[j])
		}
	}
}

//  TorusPolynomial = TorusPolynomial
func TestTorusPolynomialCopy(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pol := NewTorusPolynomial(int32(N))
		polc := NewTorusPolynomial(int32(N))
		torusPolynomialUniform(pol)
		torusPolynomialUniform(polc)
		pol0 := pol.CoefsT[0]
		pol1 := pol.CoefsT[1]
		TorusPolynomialCopy(polc, pol)
		//check that the copy is in the right direction
		assert.EqualValues(pol0, polc.CoefsT[0])
		assert.EqualValues(pol1, polc.CoefsT[1])
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(pol.CoefsT[j], polc.CoefsT[j])
		}
	}
}

//  TorusPolynomial + TorusPolynomial
func TestTorusPolynomialAdd(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		polc := &pols[4]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialAdd(polc, pola, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(polc.CoefsT[j], pola.CoefsT[j]+polb.CoefsT[j])
		}
	}
}

//  TorusPolynomial += TorusPolynomial
func TestTorusPolynomialAddTo(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialAddTo(pola, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(pola.CoefsT[j], polacopy.CoefsT[j]+polbcopy.CoefsT[j])
		}
	}
}

//  TorusPolynomial - TorusPolynomial
func TestTorusPolynomialSub(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		polc := &pols[4]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialSub(polc, pola, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(polc.CoefsT[j], pola.CoefsT[j]-polb.CoefsT[j])
		}
	}
}

//  TorusPolynomial -= TorusPolynomial
//EXPORT void torusPolynomialSubTo(TorusPolynomial* result, const TorusPolynomial* poly2)
func TestTorusPolynomialSubTo(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialSubTo(pola, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(pola.CoefsT[j], polacopy.CoefsT[j]-polbcopy.CoefsT[j])
		}
	}
}

//  TorusPolynomial + p*TorusPolynomial
func TestTorusPolynomialAddMulZ(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorus32Dist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		polc := &pols[4]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialAddMulZ(polc, pola, p, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(polc.CoefsT[j], pola.CoefsT[j]+p*polb.CoefsT[j])
		}
	}
}

//  TorusPolynomial += p*TorusPolynomial
//EXPORT void torusPolynomialAddMulZTo(TorusPolynomial* result, const int32_t p, const TorusPolynomial* poly2)
func TestTorusPolynomialAddMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorus32Dist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialAddMulZTo(pola, p, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(pola.CoefsT[j], polacopy.CoefsT[j]+p*polbcopy.CoefsT[j])
		}
	}
}

//  TorusPolynomial - p*TorusPolynomial
//EXPORT void torusPolynomialSubMulZ(TorusPolynomial* result, const TorusPolynomial* poly1, int32_t p, const TorusPolynomial* poly2)
func TestTorusPolynomialSubMulZ(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorus32Dist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		polc := &pols[4]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialSubMulZ(polc, pola, p, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(polc.CoefsT[j], pola.CoefsT[j]-p*polb.CoefsT[j])
		}
	}
}

//  TorusPolynomial -= p*TorusPolynomial
//EXPORT void torusPolynomialSubMulZTo(TorusPolynomial* result, const int32_t p, const TorusPolynomial* poly2)
func TestTorusPolynomialSubMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorus32Dist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, int32(N))
		pola := &pols[0]
		polacopy := &pols[1]
		polb := &pols[2]
		polbcopy := &pols[3]
		torusPolynomialUniform(pola)
		torusPolynomialUniform(polb)
		TorusPolynomialCopy(polacopy, pola)
		TorusPolynomialCopy(polbcopy, polb)
		TorusPolynomialSubMulZTo(pola, p, polb)
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(polbcopy.CoefsT[j], polb.CoefsT[j])
			assert.EqualValues(pola.CoefsT[j], polacopy.CoefsT[j]-p*polbcopy.CoefsT[j])
		}
	}
}

func anticyclicGet(tab []int32, a int32, N int32) int32 {
	agood := ((a % (2 * N)) + (2 * N)) % (2 * N)
	if agood < N {
		return tab[agood]
	} else {
		return -tab[agood-N]
	}
}

func randomSmallInts(tab []int32, bound, N int32) {
	for j := int32(0); j < N; j++ {
		tab[j] = (UniformTorus32Dist() % bound)
	}
}

func intTabCopy(dest, tab []int32, N int32) {
	for j := int32(0); j < N; j++ {
		dest[j] = tab[j]
	}
}

//  TorusPolynomial = X^a * TorusPolynomial
//EXPORT void torusPolynomialMulByXai(TorusPolynomial* result, int32_t a, const TorusPolynomial* bk)
func TestTorusPolynomialMulByXai(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 50
	for _, N := range dimensions {
		for trial := 0; trial < NB_TRIALS; trial++ {
			//TODO: parallelization
			a := (UniformTorus32Dist() % 1000000) - 500000
			ai := ((a % int32(2*N)) + int32(2*N)) % int32(2*N)
			pols := NewTorusPolynomialArray(4, int32(N))
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			torusPolynomialUniform(pola)
			torusPolynomialUniform(polb)
			TorusPolynomialCopy(polacopy, pola)
			TorusPolynomialMulByXai(polb, ai, pola)
			//check equality
			for j := int32(0); j < int32(N); j++ {
				assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
				assert.EqualValues(polb.CoefsT[j], anticyclicGet(polacopy.CoefsT, j-ai, int32(N)))
			}
		}
	}
}

//  intPolynomial = (X^ai-1) * intPolynomial
//EXPORT void intPolynomialMulByXaiMinusOne(IntPolynomial* result, int32_t a, const IntPolynomial* bk)
func TestIntPolynomialMulByXaiMinusOne(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 50
	for _, N := range dimensions {
		for trial := 0; trial < NB_TRIALS; trial++ {
			//TODO: parallelization
			a := (UniformTorus32Dist() % 1000000) - 500000
			ai := ((a % int32(2*N)) + int32(2*N)) % int32(2*N)
			pols := NewIntPolynomialArray(3, int32(N))
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			//fill the polynomial with random Coefs
			for j := 0; j < N; j++ {
				pola.Coefs[j] = UniformTorus32Dist()
				polb.Coefs[j] = UniformTorus32Dist()
			}
			intPolynomialCopy(polacopy, pola)
			intPolynomialMulByXaiMinusOne(polb, ai, pola)
			//check equality
			for j := int32(0); j < int32(N); j++ {
				assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
				assert.EqualValues(polb.Coefs[j],
					anticyclicGet(polacopy.Coefs, j-ai, int32(N))-anticyclicGet(polacopy.Coefs, j, int32(N)))
			}
		}
	}
}

//  TorusPolynomial = (X^ai-1) * TorusPolynomial
//EXPORT void torusPolynomialMulByXaiMinusOne(TorusPolynomial* result, int32_t a, const TorusPolynomial* bk)
func TestTorusPolynomialMulByXaiMinusOne(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 50
	for _, N := range dimensions {
		for trial := 0; trial < NB_TRIALS; trial++ {
			//TODO: parallelization
			a := (UniformTorus32Dist() % 1000000) - 500000
			ai := ((a % int32(2*N)) + int32(2*N)) % int32(2*N)
			pols := NewTorusPolynomialArray(3, int32(N))
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			//fill the polynomial with random Coefs
			torusPolynomialUniform(pola)
			torusPolynomialUniform(polb)
			TorusPolynomialCopy(polacopy, pola)
			TorusPolynomialMulByXaiMinusOne(polb, ai, pola)
			//check equality
			for j := int32(0); j < int32(N); j++ {
				assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
				assert.EqualValues(polb.CoefsT[j],
					anticyclicGet(polacopy.CoefsT, j-ai, int32(N))-anticyclicGet(polacopy.CoefsT, j, int32(N)))
			}
		}
	}
}

//  Norme Euclidienne d'un IntPolynomial
//EXPORT double intPolynomialNormSq2(const IntPolynomial* poly)
func TestIntPolynomialNormSq2(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 50
	for _, N := range dimensions {
		pols := NewIntPolynomialArray(2, int32(N))
		a := &pols[0]
		acopy := &pols[1]
		for trial := 0; trial < NB_TRIALS; trial++ {
			var norm2 double = 0
			for j := 0; j < N; j++ {
				r := (UniformTorus32Dist() % 1000) - 500
				a.Coefs[j] = r
				acopy.Coefs[j] = r
				norm2 += double(r * r)
			}
			value := intPolynomialNormSq2(a)
			assert.EqualValues(norm2, value)
			for j := 0; j < N; j++ {
				assert.EqualValues(a.Coefs[j], acopy.Coefs[j])
			}
		}
	}
}

// This is the naive external multiplication of an integer polynomial
// with a torus polynomial. (this function should yield exactly the same
// result as the karatsuba or fft version, but should be slower)
//EXPORT void torusPolynomialMultNaive(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2)
func TestPolynomialMultNaive(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 5
	for _, N := range dimensions {
		ipols := NewIntPolynomialArray(2, int32(N))
		tpols := NewTorusPolynomialArray(3, int32(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < NB_TRIALS; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, int32(N))
			intTabCopy(acopy.Coefs, a.Coefs, int32(N))
			torusPolynomialMultNaive(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus32 = 0
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// This is the karatsuba external multiplication of an integer polynomial
// with a torus polynomial.
// WARNING: for karatsuba, N must be a power of 2
//EXPORT void torusPolynomialMultKaratsuba(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2)
func TestTorusPolynomialMultKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 5
	for _, N := range powers_of_two_dimensions {
		ipols := NewIntPolynomialArray(2, int32(N))
		tpols := NewTorusPolynomialArray(3, int32(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < NB_TRIALS; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, int32(N))
			intTabCopy(acopy.Coefs, a.Coefs, int32(N))
			torusPolynomialMultKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus32 = 0
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// result += poly1 * poly2 (via karatsuba)
// WARNING: N must be a power of 2 to use this function. Else, the
// behaviour is unpredictable
//EXPORT void torusPolynomialAddMulRKaratsuba(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2)
func TestTorusPolynomialAddMulRKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 5
	for _, N := range powers_of_two_dimensions {
		ipols := NewIntPolynomialArray(2, int32(N))
		tpols := NewTorusPolynomialArray(4, int32(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		ccopy := &tpols[3]
		for trial := 0; trial < NB_TRIALS; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			TorusPolynomialCopy(ccopy, c)
			randomSmallInts(a.Coefs, 100000, int32(N))
			intTabCopy(acopy.Coefs, a.Coefs, int32(N))
			torusPolynomialAddMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus32 = ccopy.CoefsT[j]
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// result -= poly1 * poly2 (via karatsuba)
// WARNING: N must be a power of 2 to use this function. Else, the
// behaviour is unpredictable
//EXPORT void torusPolynomialAddMulRKaratsuba(TorusPolynomial* result, const IntPolynomial* poly1, const TorusPolynomial* poly2)
func TestTorusPolynomialSubMulRKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	NB_TRIALS := 5
	for _, N := range powers_of_two_dimensions {
		ipols := NewIntPolynomialArray(2, int32(N))
		tpols := NewTorusPolynomialArray(4, int32(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		ccopy := &tpols[3]
		for trial := 0; trial < NB_TRIALS; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			TorusPolynomialCopy(ccopy, c)
			randomSmallInts(a.Coefs, 100000, int32(N))
			intTabCopy(acopy.Coefs, a.Coefs, int32(N))
			torusPolynomialSubMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus32 = ccopy.CoefsT[j]
				for k := 0; k < N; k++ {
					r -= bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}
