package tfhe

import (
	"math/big"
)

func torusPolynomialMultNaive_plain_aux(result []Torus32, poly1 []int32, poly2 []Torus32, N int32) {
	_2Nm1 := 2*N - 1
	var ri Torus32
	for i := int32(0); i < N; i++ {
		ri = 0
		for j := int32(0); j <= i; j++ {
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

func torusPolynomialMultNaive_aux(result []Torus32, poly1 []int32, poly2 []Torus32, N int32) {
	var ri Torus32
	for i := int32(0); i < N; i++ {
		ri = 0
		for j := int32(0); j <= i; j++ {
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
	torusPolynomialMultNaive_aux(result.CoefsT, poly1.Coefs, poly2.CoefsT, N)
}

/**
 * This function multiplies 2 polynomials (an integer poly and a torus poly) by using Karatsuba
 * The karatsuba function is torusPolynomialMultKaratsuba: it takes in input two polynomials and multiplies them
 * To do that, it uses the auxiliary function Karatsuba_aux, which is recursive ad which works with
 * the vectors containing the coefficients of the polynomials (primitive types)
 */

// A and B of size = size
// R of size = 2*size-1
/*
func Karatsuba_aux(R []Torus32, A []int32, B []Torus32, size int32, buf []byte) {
	h := size / 2
	sm1 := size - 1

	//we stop the karatsuba recursion at h=4, because on my machine,
	//it seems to be optimal
	if h <= 4 {
		torusPolynomialMultNaive_plain_aux(R, A, B, size)
		return
	}

	//we split the polynomials in 2
	Atemp := make([]int32, h)
	//int32_t* Atemp = (int32_t*) buf
	//buf += h*sizeof(int32_t)

	//Torus32* Btemp = (Torus32*) buf
	//buf += h*sizeof(Torus32)
	Btemp := make([]Torus32, h)

	//Torus32* Rtemp = (Torus32*) buf
	//buf += size*sizeof(Torus32)
	Rtemp := make([]Torus32, size)
	//Note: in the above line, I have put size instead of sm1 so that buf remains aligned on a power of 2

	for i := int32(0); i < h; i++ {
		Atemp[i] = A[i] + A[h+i]
	}
	for i := int32(0); i < h; i++ {
		Btemp[i] = B[i] + B[h+i]
	}

	// Karatsuba recursivly
	Karatsuba_aux(R, A, B, h, buf)          // (R[0],R[2*h-2]), (A[0],A[h-1]), (B[0],B[h-1])
	Karatsuba_aux(R+size, A+h, B+h, h, buf) // (R[2*h],R[4*h-2]), (A[h],A[2*h-1]), (B[h],B[2*h-1])
	Karatsuba_aux(Rtemp, Atemp, Btemp, h, buf)
	R[sm1] = 0 //this one needs to be set manually
	for i := int32(0); i < sm1; i++ {
		Rtemp[i] -= R[i] + R[size+i]
	}
	for i := int32(0); i < sm1; i++ {
		R[h+i] += Rtemp[i]
	}
}
*/

func toBig(a []int32) (res []*big.Int) {
	res = make([]*big.Int, len(a))
	for i, v := range a {
		res[i] = big.NewInt(int64(v))
	}
	return
}

func fromBig(a []*big.Int) (res []int32) {
	res = make([]int32, len(a))
	for i, v := range a {
		res[i] = int32(v.Int64())
	}
	return
}

func Mul(x, y []*big.Int) (res []*big.Int) {
	res = make([]*big.Int, len(x))
	for i, _ := range x {
		res[i] = big.NewInt(0).Mul(x[i], y[i])
	}
	return
}

func Karatsuba_aux(multiplicand []int32, multiplier []Torus32) []Torus32 {
	// big.Int golang library usses karatsuba for operations
	mplcand := toBig(multiplicand)
	mplier := toBig(multiplier)
	p := Mul(mplcand, mplier)
	product := fromBig(p)
	return product

	// new double[2 * multiplicand.length];
	/*
		product := make([]Torus32, 2*len(multiplicand))
		//product := make([]Torus32, 2*len(multiplicand)-1)

		//Handle the base case where the polynomial has only one coefficient
		if len(multiplicand) == 1 {
			product[0] = karatsuba.Multiply(multiplicand[0], multiplier[0]) //multiplicand[0] * multiplier[0]
			return product
		}

		halfArraySize := len(multiplicand) / 2

		//Declare arrays to hold halved factors
		multiplicandLow := make([]int32, halfArraySize) //new double[halfArraySize];
		multiplicandHigh := make([]int32, halfArraySize)
		multipliplierLow := make([]Torus32, halfArraySize)
		multipliierHigh := make([]Torus32, halfArraySize)

		multiplicandLowHigh := make([]int32, halfArraySize)
		multipliplierLowHigh := make([]Torus32, halfArraySize)

		//Fill in the low and high arrays
		for halfSizeIndex := 0; halfSizeIndex < halfArraySize; halfSizeIndex++ {
			multiplicandLow[halfSizeIndex] = multiplicand[halfSizeIndex]
			multiplicandHigh[halfSizeIndex] = multiplicand[halfSizeIndex+halfArraySize]
			multiplicandLowHigh[halfSizeIndex] = multiplicandLow[halfSizeIndex] + multiplicandHigh[halfSizeIndex]

			multipliplierLow[halfSizeIndex] = multiplier[halfSizeIndex]
			multipliierHigh[halfSizeIndex] = multiplier[halfSizeIndex+halfArraySize]
			multipliplierLowHigh[halfSizeIndex] = multipliplierLow[halfSizeIndex] + multipliierHigh[halfSizeIndex]
		}

		//Recursively call method on smaller arrays and construct the low and high parts of the product
		productLow := Karatsuba_aux(multiplicandLow, multipliplierLow)
		productHigh := Karatsuba_aux(multiplicandHigh, multipliierHigh)
		productLowHigh := Karatsuba_aux(multiplicandLowHigh, multipliplierLowHigh)

		//Construct the middle portion of the product
		productMiddle := make([]int32, len(multiplicand)) //new double[multiplicand.length];
		for halfSizeIndex := 0; halfSizeIndex < len(multiplicand); halfSizeIndex++ {
			productMiddle[halfSizeIndex] = productLowHigh[halfSizeIndex] - productLow[halfSizeIndex] - productHigh[halfSizeIndex]
		}

		//Assemble the product from the low, middle and high parts. Start with the low and high parts of the product.
		middleOffset := len(multiplicand) / 2
		for halfSizeIndex := 0; halfSizeIndex < len(multiplicand); halfSizeIndex++ {
			product[halfSizeIndex] += productLow[halfSizeIndex]
			product[halfSizeIndex+len(multiplicand)] += productHigh[halfSizeIndex]
			product[halfSizeIndex+middleOffset] += productMiddle[halfSizeIndex]
		}
		return product
	*/
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialMultKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	//N := poly1.N
	//R := make([]Torus32, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	result.CoefsT = Karatsuba_aux(poly1.Coefs, poly2.CoefsT)
	/* R := Karatsuba_aux(poly1.Coefs, poly2.CoefsT)

	// reduction mod X^N+1
	for i := int32(0); i < N-1; i++ {
		result.CoefsT[i] = R[i] - R[N+i]
	}
	result.CoefsT[N-1] = R[N-1]
	*/
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialAddMulRKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	//R := make([]Torus32, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	//Karatsuba_aux(R, poly1.coefs, poly2.coefsT, N, buf)
	R := Karatsuba_aux(poly1.Coefs, poly2.CoefsT)

	// reduction mod X^N+1
	for i := int32(0); i < N-1; i++ {
		result.CoefsT[i] += R[i] - R[N+i]
	}
	result.CoefsT[N-1] += R[N-1]
}

// poly1, poly2 and result are polynomials mod X^N+1
func torusPolynomialSubMulRKaratsuba(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {
	N := poly1.N
	R := make([]Torus32, 2*N-1)
	//buf := make([]byte, 16*N) //that's large enough to store every tmp variables (2*2*N*4)

	// Karatsuba
	//Karatsuba_aux(R, poly1.coefs, poly2.coefsT, N, buf)
	R = Karatsuba_aux(poly1.Coefs, poly2.CoefsT)

	// reduction mod X^N+1
	for i := int32(0); i < N-1; i++ {
		result.CoefsT[i] -= R[i] - R[N+i]
	}
	result.CoefsT[N-1] -= R[N-1]
}
