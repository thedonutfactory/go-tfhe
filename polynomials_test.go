package tfhe

import (
	"testing"

	set "github.com/badgerodon/collections/set"
	"github.com/stretchr/testify/assert"
)

var (
	dimensions            = [...]int{500, 750, 1024, 2000}
	powersOfTwoDimensions = [...]int{512, 1024, 2048}
)

func TestTorusPolynomialUniform(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	const nbTrials int = 10
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(nbTrials, N)
		for i := 0; i < nbTrials; i++ {
			torusPolynomialUniform(&pols[i])
		}
		for j := 0; j < N; j++ {
			testset := set.New()
			for i := 0; i < nbTrials; i++ {
				testset.Insert(pols[i].CoefsT[j])
			}
			assert.GreaterOrEqual(float64(testset.Len()), 0.9*float64(nbTrials))
		}
	}
}

//  TorusPolynomial = 0

func TestTorusPolynomialClear(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pol := NewTorusPolynomial(N)
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
		pol := NewTorusPolynomial(N)
		polc := NewTorusPolynomial(N)
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
		pols := NewTorusPolynomialArray(5, N)
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
		pols := NewTorusPolynomialArray(5, N)
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
		pols := NewTorusPolynomialArray(5, N)
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
		pols := NewTorusPolynomialArray(5, N)
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
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, N)
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
func TestTorusPolynomialAddMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, N)
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
func TestTorusPolynomialSubMulZ(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, N)
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
func TestTorusPolynomialSubMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, N)
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

func anticyclicGet(tab []int64, a, N int64) int64 {
	agood := ((a % (2 * N)) + (2 * N)) % (2 * N)
	if agood < N {
		return tab[agood]
	} else {
		return -tab[agood-N]
	}
}

func randomSmallInts(tab []int64, bound int64, N int) {
	for j := 0; j < N; j++ {
		tab[j] = (UniformTorusDist() % bound)
	}
}

func intTabCopy(dest, tab []int64, N int) {
	for j := 0; j < N; j++ {
		dest[j] = tab[j]
	}
}

//  TorusPolynomial = X^a * TorusPolynomial
func TestTorusPolynomialMulByXai(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 50
	for _, N := range dimensions {
		for trial := 0; trial < nbTrials; trial++ {
			//TODO: parallelization
			a := (UniformTorusDist() % 1000000) - 500000
			ai := ((a % int64(2*N)) + int64(2*N)) % int64(2*N)
			pols := NewTorusPolynomialArray(4, N)
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			torusPolynomialUniform(pola)
			torusPolynomialUniform(polb)
			TorusPolynomialCopy(polacopy, pola)
			TorusPolynomialMulByXai(polb, ai, pola)
			//check equality
			for j := int64(0); j < int64(N); j++ {
				assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
				assert.EqualValues(polb.CoefsT[j], anticyclicGet(polacopy.CoefsT, j-ai, int64(N)))
			}
		}
	}
}

//  intPolynomial = (X^ai-1) * intPolynomial
func TestIntPolynomialMulByXaiMinusOne(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 50
	for _, N := range dimensions {
		for trial := 0; trial < nbTrials; trial++ {
			//TODO: parallelization
			a := (UniformTorusDist() % 1000000) - 500000
			ai := ((a % int64(2*N)) + int64(2*N)) % int64(2*N)
			pols := NewIntPolynomialArray(3, N)
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			//fill the polynomial with random Coefs
			for j := 0; j < N; j++ {
				pola.Coefs[j] = UniformTorusDist()
				polb.Coefs[j] = UniformTorusDist()
			}
			intPolynomialCopy(polacopy, pola)
			intPolynomialMulByXaiMinusOne(polb, ai, pola)
			//check equality
			for j := int64(0); j < int64(N); j++ {
				assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
				assert.EqualValues(polb.Coefs[j],
					anticyclicGet(polacopy.Coefs, j-ai, int64(N))-anticyclicGet(polacopy.Coefs, j, int64(N)))
			}
		}
	}
}

//  TorusPolynomial = (X^ai-1) * TorusPolynomial
func TestTorusPolynomialMulByXaiMinusOne(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 50
	for _, N := range dimensions {
		for trial := 0; trial < nbTrials; trial++ {
			//TODO: parallelization
			a := (UniformTorusDist() % 1000000) - 500000
			ai := ((a % int64(2*N)) + int64(2*N)) % int64(2*N)
			pols := NewTorusPolynomialArray(3, N)
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			//fill the polynomial with random Coefs
			torusPolynomialUniform(pola)
			torusPolynomialUniform(polb)
			TorusPolynomialCopy(polacopy, pola)
			TorusPolynomialMulByXaiMinusOne(polb, ai, pola)
			//check equality
			for j := int64(0); j < int64(N); j++ {
				assert.EqualValues(polacopy.CoefsT[j], pola.CoefsT[j])
				assert.EqualValues(polb.CoefsT[j],
					anticyclicGet(polacopy.CoefsT, j-ai, int64(N))-anticyclicGet(polacopy.CoefsT, j, int64(N)))
			}
		}
	}
}

//  Norme Euclidienne d'un IntPolynomial
func TestIntPolynomialNormSq2(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 50
	for _, N := range dimensions {
		pols := NewIntPolynomialArray(2, N)
		a := &pols[0]
		acopy := &pols[1]
		for trial := 0; trial < nbTrials; trial++ {
			var norm2 double = 0
			for j := 0; j < N; j++ {
				r := (UniformTorusDist() % 1000) - 500
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
func TestPolynomialMultNaive(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 5
	for _, N := range dimensions {
		ipols := NewIntPolynomialArray(2, N)
		tpols := NewTorusPolynomialArray(3, N)
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, N)
			intTabCopy(acopy.Coefs, a.Coefs, N)
			torusPolynomialMultNaive(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus = 0
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int64(j-k), int64(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// This is the karatsuba external multiplication of an integer polynomial
// with a torus polynomial.
// WARNING: for karatsuba, N must be a power of 2
func TestTorusPolynomialMultKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 5
	for _, N := range powersOfTwoDimensions {
		ipols := NewIntPolynomialArray(2, N)
		tpols := NewTorusPolynomialArray(3, N)
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, N)
			intTabCopy(acopy.Coefs, a.Coefs, N)
			torusPolynomialMultKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus = 0
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int64(j-k), int64(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// result += poly1 * poly2 (via karatsuba)
// WARNING: N must be a power of 2 to use this function. Else, the
// behaviour is unpredictable
func TestTorusPolynomialAddMulRKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 5
	for _, N := range powersOfTwoDimensions {
		ipols := NewIntPolynomialArray(2, N)
		tpols := NewTorusPolynomialArray(4, N)
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		ccopy := &tpols[3]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			TorusPolynomialCopy(ccopy, c)
			randomSmallInts(a.Coefs, 100000, N)
			intTabCopy(acopy.Coefs, a.Coefs, N)
			torusPolynomialAddMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus = ccopy.CoefsT[j]
				for k := 0; k < N; k++ {
					r += bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int64(j-k), int64(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}

// result -= poly1 * poly2 (via karatsuba)
// WARNING: N must be a power of 2 to use this function. Else, the
// behaviour is unpredictable
func TestTorusPolynomialSubMulRKaratsuba(t *testing.T) {
	assert := assert.New(t)
	//TODO: parallelization
	nbTrials := 5
	for _, N := range powersOfTwoDimensions {
		ipols := NewIntPolynomialArray(2, N)
		tpols := NewTorusPolynomialArray(4, N)
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		ccopy := &tpols[3]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			TorusPolynomialCopy(ccopy, c)
			randomSmallInts(a.Coefs, 100000, N)
			intTabCopy(acopy.Coefs, a.Coefs, N)
			torusPolynomialSubMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.CoefsT[j], b.CoefsT[j])
				var r Torus = ccopy.CoefsT[j]
				for k := 0; k < N; k++ {
					r -= bcopy.CoefsT[k] * anticyclicGet(acopy.Coefs, int64(j-k), int64(N))
				}
				assert.EqualValues(r, c.CoefsT[j])
			}
		}
	}
}
