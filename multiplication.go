package tfhe

import (
	"math/big"
)

func torusPolynomialMultNaivePlainAux(result []Torus, poly1 []int64, poly2 []Torus, N int) {
	_2Nm1 := 2*N - 1
	var ri Torus
	for i := 0; i < N; i++ {
		ri = 0
		for j := 0; j <= i; j++ {
			ri += poly1[j] * poly2[i-j]
		}
		result[i] = ri
	}
	for i := N; i < _2Nm1; i++ {
		ri = 0
		for j := i - N + 1; j < N; j++ {
			ri += poly1[j] * poly2[i-j]
		}
		result[i] = ri
	}
}

func torusPolynomialMultNaiveAux(result []Torus, poly1 []int64, poly2 []Torus, N int) {
	var ri Torus
	for i := 0; i < N; i++ {
		ri = 0
		for j := 0; j <= i; j++ {
			ri += poly1[j] * poly2[i-j]
		}
		for j := i + 1; j < N; j++ {
			ri -= poly1[j] * poly2[N+i-j]
		}
		result[i] = ri
	}
}

/**
 * This is the naive external multiplication of an integer polynomial
 * with a torus polynomial. (this function should yield exactly the same
 * result as the karatsuba or fft version)
 */
func torusPolynomialMultNaive(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	Assert(result != poly2)
	Assert(poly2.N == N && result.N == N)
	torusPolynomialMultNaiveAux(result.CoefsT, poly1.Coefs, poly2.CoefsT, N)
}

/**
 * This function multiplies 2 polynomials (an integer poly and a torus poly) by using Karatsuba
 * The karatsuba function is torusPolynomialMultKaratsuba: it takes in input two polynomials and multiplies them
 * To do that, it uses the auxiliary function karatsubaAux, which is recursive ad which works with
 * the vectors containing the coefficients of the polynomials (primitive types)
 */

// A and B of size = size
// R of size = 2*size-1

func karatsubaAux(R []Torus, A []int64, B []Torus, size int) {
	h := size / 2
	sm1 := size - 1

	//we stop the karatsuba recursion at h=4, because on my machine,
	//it seems to be optimal
	if h <= 4 {
		torusPolynomialMultNaivePlainAux(R, A, B, size)
		return
	}

	//we split the polynomials in 2
	Atemp := make([]int64, h)
	Btemp := make([]Torus, h)
	Rtemp := make([]Torus, size)
	//Note: in the above line, I have put size instead of sm1 so that buf remains aligned on a power of 2
	for i := 0; i < h; i++ {
		Atemp[i] = A[i] + A[h+i]
	}
	for i := 0; i < h; i++ {
		Btemp[i] = B[i] + B[h+i]
	}

	// Karatsuba recursivly
	karatsubaAux(R, A, B, h) // (R[0],R[2*h-2]), (A[0],A[h-1]), (B[0],B[h-1])
	// karatsubaAux(R+size, A+h, B+h, h, buf) // (R[2*h],R[4*h-2]), (A[h],A[2*h-1]), (B[h],B[2*h-1])
	karatsubaAux(R[size:], A[h:], B[h:], h)
	karatsubaAux(Rtemp, Atemp, Btemp, h)
	R[sm1] = 0 //this one needs to be set manually
	for i := 0; i < sm1; i++ {
		Rtemp[i] -= R[i] + R[size+i]
	}
	for i := 0; i < sm1; i++ {
		R[h+i] += Rtemp[i]
	}
}

func Mul(x, y []*big.Int) (res []*big.Int) {
	res = make([]*big.Int, len(x))
	for i, _ := range x {
		res[i] = big.NewInt(0).Mul(x[i], y[i])
	}
	return
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialMultKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	R := make([]Torus, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	//result.CoefsT = karatsubaAux(poly1.Coefs, poly2.CoefsT)
	karatsubaAux(R, poly1.Coefs, poly2.CoefsT, N)

	// reduction mod X^N+1
	for i := 0; i < N-1; i++ {
		result.CoefsT[i] = R[i] - R[N+i]
	}
	result.CoefsT[N-1] = R[N-1]
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialAddMulRKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	R := make([]Torus, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	karatsubaAux(R, poly1.Coefs, poly2.CoefsT, N)
	//R := karatsubaAux(poly1.Coefs, poly2.CoefsT)

	// reduction mod X^N+1
	for i := 0; i < N-1; i++ {
		result.CoefsT[i] += R[i] - R[N+i]
	}
	result.CoefsT[N-1] += R[N-1]
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialSubMulRKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	R := make([]Torus, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	karatsubaAux(R, poly1.Coefs, poly2.CoefsT, N)
	//R = karatsubaAux(poly1.Coefs, poly2.CoefsT)

	// reduction mod X^N+1
	for i := 0; i < N-1; i++ {
		result.CoefsT[i] -= R[i] - R[N+i]
	}
	result.CoefsT[N-1] -= R[N-1]
}
