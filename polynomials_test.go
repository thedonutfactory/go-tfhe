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
		pols := NewTorusPolynomialArray(nbTrials, int(N))
		for i := 0; i < nbTrials; i++ {
			torusPolynomialUniform(&pols[i])
		}
		for j := 0; j < N; j++ {
			testset := set.New()
			for i := 0; i < nbTrials; i++ {
				testset.Insert(pols[i].Coefs[j])
			}
			assert.GreaterOrEqual(float64(testset.Len()), 0.9*float64(nbTrials))
		}
	}
}

//  TorusPolynomial = 0

func TestTorusPolynomialClear(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pol := NewTorusPolynomial(int(N))
		torusPolynomialUniform(pol)
		torusPolynomialClear(pol)
		for j := 0; j < N; j++ {
			assert.EqualValues(0, pol.Coefs[j])
		}
	}
}

//  TorusPolynomial = TorusPolynomial
func TestTorusPolynomialCopy(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pol := NewTorusPolynomial(int(N))
		polc := NewTorusPolynomial(int(N))
		torusPolynomialUniform(pol)
		torusPolynomialUniform(polc)
		pol0 := pol.Coefs[0]
		pol1 := pol.Coefs[1]
		TorusPolynomialCopy(polc, pol)
		//check that the copy is in the right direction
		assert.EqualValues(pol0, polc.Coefs[0])
		assert.EqualValues(pol1, polc.Coefs[1])
		//check equality
		for j := 0; j < N; j++ {
			assert.EqualValues(pol.Coefs[j], polc.Coefs[j])
		}
	}
}

//  TorusPolynomial + TorusPolynomial
func TestTorusPolynomialAdd(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(polc.Coefs[j], pola.Coefs[j]+polb.Coefs[j])
		}
	}
}

//  TorusPolynomial += TorusPolynomial
func TestTorusPolynomialAddTo(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(pola.Coefs[j], polacopy.Coefs[j]+polbcopy.Coefs[j])
		}
	}
}

//  TorusPolynomial - TorusPolynomial
func TestTorusPolynomialSub(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(polc.Coefs[j], pola.Coefs[j]-polb.Coefs[j])
		}
	}
}

//  TorusPolynomial -= TorusPolynomial
//EXPORT void torusPolynomialSubTo(TorusPolynomial* result, const TorusPolynomial* poly2)
func TestTorusPolynomialSubTo(t *testing.T) {
	assert := assert.New(t)
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(pola.Coefs[j], polacopy.Coefs[j]-polbcopy.Coefs[j])
		}
	}
}

//  TorusPolynomial + p*TorusPolynomial
func TestTorusPolynomialAddMulZ(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(polc.Coefs[j], pola.Coefs[j]+p*polb.Coefs[j])
		}
	}
}

//  TorusPolynomial += p*TorusPolynomial
func TestTorusPolynomialAddMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, int(N))
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
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(pola.Coefs[j], polacopy.Coefs[j]+p*polbcopy.Coefs[j])
		}
	}
}

//  TorusPolynomial - p*TorusPolynomial
func TestTorusPolynomialSubMulZ(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(5, int(N))
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
			assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(polc.Coefs[j], pola.Coefs[j]-p*polb.Coefs[j])
		}
	}
}

//  TorusPolynomial -= p*TorusPolynomial
func TestTorusPolynomialSubMulZTo(t *testing.T) {
	assert := assert.New(t)
	p := UniformTorusDist()
	for _, N := range dimensions {
		pols := NewTorusPolynomialArray(4, int(N))
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
			assert.EqualValues(polbcopy.Coefs[j], polb.Coefs[j])
			assert.EqualValues(pola.Coefs[j], polacopy.Coefs[j]-p*polbcopy.Coefs[j])
		}
	}
}

func anticyclicGet(tab []int, a int, N int) int {
	agood := ((a % (2 * N)) + (2 * N)) % (2 * N)
	if agood < N {
		return tab[agood]
	} else {
		return -tab[agood-N]
	}
}

func randomSmallInts(tab []int, bound, N int) {
	for j := int(0); j < N; j++ {
		tab[j] = (UniformTorusDist() % bound)
	}
}

func intTabCopy(dest, tab []int, N int) {
	for j := int(0); j < N; j++ {
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
			ai := ((a % int(2*N)) + int(2*N)) % int(2*N)
			pols := NewTorusPolynomialArray(4, int(N))
			pola := &pols[0]
			polacopy := &pols[1]
			polb := &pols[2]
			torusPolynomialUniform(pola)
			torusPolynomialUniform(polb)
			TorusPolynomialCopy(polacopy, pola)
			TorusPolynomialMulByXai(polb, ai, pola)
			//check equality
			for j := int32(0); j < int32(N); j++ {
				assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
				assert.EqualValues(polb.Coefs[j], anticyclicGet(polacopy.Coefs, j-ai, int32(N)))
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
			ai := ((a % int(2*N)) + int(2*N)) % int(2*N)
			pols := NewIntPolynomialArray(3, int(N))
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
			for j := int(0); j < int(N); j++ {
				assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
				assert.EqualValues(polb.Coefs[j],
					anticyclicGet(polacopy.Coefs, j-ai, int(N))-anticyclicGet(polacopy.Coefs, j, int(N)))
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
			ai := ((a % int(2*N)) + int(2*N)) % int(2*N)
			pols := NewTorusPolynomialArray(3, int(N))
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
				assert.EqualValues(polacopy.Coefs[j], pola.Coefs[j])
				assert.EqualValues(polb.Coefs[j],
					anticyclicGet(polacopy.Coefs, j-ai, int32(N))-anticyclicGet(polacopy.Coefs, j, int32(N)))
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
		pols := NewIntPolynomialArray(2, int(N))
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
		ipols := NewIntPolynomialArray(2, int(N))
		tpols := NewTorusPolynomialArray(3, int(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, int(N))
			intTabCopy(acopy.Coefs, a.Coefs, int(N))
			torusPolynomialMultNaive(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.Coefs[j], b.Coefs[j])
				var r Torus32 = 0
				for k := 0; k < N; k++ {
					r += bcopy.Coefs[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.Coefs[j])
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
		ipols := NewIntPolynomialArray(2, int(N))
		tpols := NewTorusPolynomialArray(3, int(N))
		a := &ipols[0]
		acopy := &ipols[1]
		b := &tpols[0]
		bcopy := &tpols[1]
		c := &tpols[2]
		for trial := 0; trial < nbTrials; trial++ {
			torusPolynomialUniform(b)
			torusPolynomialUniform(c)
			TorusPolynomialCopy(bcopy, b)
			randomSmallInts(a.Coefs, 100000, int(N))
			intTabCopy(acopy.Coefs, a.Coefs, int(N))
			torusPolynomialMultKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.Coefs[j], b.Coefs[j])
				var r Torus32 = 0
				for k := 0; k < N; k++ {
					r += bcopy.Coefs[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.Coefs[j])
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
		ipols := NewIntPolynomialArray(2, int(N))
		tpols := NewTorusPolynomialArray(4, int(N))
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
			randomSmallInts(a.Coefs, 100000, int(N))
			intTabCopy(acopy.Coefs, a.Coefs, int(N))
			torusPolynomialAddMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.Coefs[j], b.Coefs[j])
				var r Torus32 = ccopy.Coefs[j]
				for k := 0; k < N; k++ {
					r += bcopy.Coefs[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.Coefs[j])
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
		ipols := NewIntPolynomialArray(2, int(N))
		tpols := NewTorusPolynomialArray(4, int(N))
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
			randomSmallInts(a.Coefs, 100000, int(N))
			intTabCopy(acopy.Coefs, a.Coefs, int(N))
			torusPolynomialSubMulRKaratsuba(c, a, b)
			for j := 0; j < N; j++ {
				assert.EqualValues(acopy.Coefs[j], a.Coefs[j])
				assert.EqualValues(bcopy.Coefs[j], b.Coefs[j])
				var r Torus32 = ccopy.Coefs[j]
				for k := 0; k < N; k++ {
					r -= bcopy.Coefs[k] * anticyclicGet(acopy.Coefs, int32(j-k), int32(N))
				}
				assert.EqualValues(r, c.Coefs[j])
			}
		}
	}
}
