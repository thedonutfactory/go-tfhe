package main

import (
	"fmt"
	"math"

	tfhe "github.com/thedonutfactory/go-tfhe"
)

const (
	N           int     = 1024
	k           int     = 1
	alphaMinGsw float64 = 0.
	alphaMaxGsw float64 = 0.071
	Msize       int     = 2
	alpha       float64 = 1e-6
	l           int     = 3
	Bgbits      int     = 10
)

// **********************************************************************************
// ********************************* MAIN *******************************************
// **********************************************************************************
func approxEquals(a tfhe.Torus, b tfhe.Torus) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {

	//TODO: parallelization
	unift := tfhe.NewUniform(0, Msize-1)

	// PARAMETERS
	rlweParams := tfhe.NewTLweParams(N, k, alphaMinGsw, alphaMaxGsw) //les deux alpha mis un peu au hasard
	rgswParams := tfhe.NewTGswParams(int(l), Bgbits, rlweParams)

	// KEY
	rgswKey := tfhe.NewTGswKey(rgswParams)
	rlweKey := &rgswKey.TlweKey
	// CIPHERTEXTS
	cipherA := tfhe.NewTGswSample(rgswParams)
	cipherB := tfhe.NewTLweSample(rlweParams)
	cipherAB := tfhe.NewTLweSample(rlweParams)

	//the probability that a sample with stdev alpha decrypts wrongly on
	//the a Msize message space.

	expectedErrorProba := 1. - math.Erf(1./(math.Sqrt(2.)*2.*float64(Msize)*alpha))

	fmt.Println("-------------")
	fmt.Println("WARNING:")
	fmt.Printf("All the tests below are supposed to fail with proba: %f \n", expectedErrorProba)
	fmt.Println("It is normal and it is part of the test!")
	fmt.Println("-------------")

	//MESSAGE RLwe
	muB := tfhe.NewTorusPolynomial(N)

	//test decompH
	fmt.Println("Test decompH on TorusPolynomial")
	muBDecH := tfhe.NewIntPolynomialArray(l, N)
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		muB.Coefs[i] = unift.Int32() //tfhe.UniformTorus32Dist()
	}
	tfhe.TGswTorus32PolynomialDecompH(muBDecH, muB, rgswParams)
	for i := int32(0); i < N; i++ {
		expected := muB.Coefs[i]
		var actual int32 = 0
=======
	for i := int(0); i < N; i++ {
		muB.CoefsT[i] = unift.int() //tfhe.UniformTorusDist()
	}
	tfhe.TGswTorusPolynomialDecompH(muBDecH, muB, rgswParams)
	for i := int(0); i < N; i++ {
		expected := muB.CoefsT[i]
		var actual int = 0
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		for j := 0; j < l; j++ {
			actual += muBDecH[j].Coefs[i] * rgswParams.H[j]
			//fmt.Printf("DEBUG: l: %d, i: %d, j: %d, muBDecH[j].Coefs[i]: %d, rgswParams->h[j]: %d \n", l, i, j, muBDecH[j].Coefs[i], rgswParams.H[j])
		}
		//fmt.Printf("\t DEBUG: actual: %d, expected: %d\n", actual, expected)
		if tfhe.Abs(expected-actual) > 3 {
			fmt.Printf("decompH error %d: %d != %d\n", i, actual, expected)
		}
	}

<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		muB.Coefs[i] = tfhe.ModSwitchToTorus32(temp, Msize)
		// fmt.Println(muB.Coefs[i])
=======
	for i := int(0); i < N; i++ {
		temp := unift.int()
		muB.CoefsT[i] = tfhe.ModSwitchToTorus(temp, Msize)
		// fmt.Println(muB.CoefsT[i])
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
	}
	//MESSAGE RLwe
	muA := tfhe.NewIntPolynomial(N)
	for i := int(0); i < N; i++ {
		temp := unift.int()
		muA.Coefs[i] = 1 - (temp % 3)
		// fmt.Println(muA.Coefs[i])
	}
	// PHASE, DECRYPTION
	dechifA := tfhe.NewIntPolynomial(N)
	dechifB := tfhe.NewTorusPolynomial(N)
	dechifAB := tfhe.NewTorusPolynomial(N)
	muAB := tfhe.NewTorusPolynomial(N)

	tfhe.TGswKeyGen(rgswKey)                          // KEY GENERATION
	tfhe.TLweSymEncrypt(cipherB, muB, alpha, rlweKey) // ENCRYPTION

	//decryption test tlwe
	fmt.Println("Test TLweSymDecrypt on muB:")
	fmt.Printf(" variance: %f\n", cipherB.CurrentVariance)
	tfhe.TLweSymDecrypt(dechifB, cipherB, rlweKey, Msize) // DECRYPTION
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muB.Coefs[i], Msize)
		actual := tfhe.ModSwitchFromTorus32(dechifB.Coefs[i], Msize)
=======
	for i := int(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus(muB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus(dechifB.CoefsT[i], int(Msize))
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	//test decompH on tLwe
	fmt.Println("Test decompH on TLwe(muB)")
	cipherBDecH := tfhe.NewIntPolynomialArray(l*(int(k)+1), N)
	tfhe.TGswTLweDecompH(cipherBDecH, cipherB, rgswParams)
<<<<<<< HEAD
	for p := int32(0); p <= k; p++ {
		for i := int32(0); i < N; i++ {
			expected := cipherB.A[p].Coefs[i]
			var actual int32 = 0
			for j := int32(0); j < int32(l); j++ {
				x := int32(l)*p + j
=======
	for p := int(0); p <= k; p++ {
		for i := int(0); i < N; i++ {
			expected := cipherB.A[p].CoefsT[i]
			var actual int = 0
			for j := int(0); j < int(l); j++ {
				x := int(l)*p + j
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
				actual += cipherBDecH[x].Coefs[i] * rgswParams.H[j]
			}
			// fails when p == 1, the array is not being populated properly
			if tfhe.Abs(expected-actual) > 3 {
				fmt.Printf("decompH error (p,i)=(%d,%d): %d != %d\n", p, i, actual, expected)
			}
			//jcl
			//fmt.Printf("\t DEBUG: actual: %d, expected: %d\n", actual, expected)
			expected2 := tfhe.ModSwitchFromTorus(expected, int(Msize))
			actual2 := tfhe.ModSwitchFromTorus(actual, int(Msize))
			if expected2 != actual2 {
				fmt.Printf("modswitch error %d: %d != %d\n", i, actual2, expected2)
			}
		}
	}

	//test externProduct with H
	tfhe.TGswClear(cipherA, rgswParams)
	tfhe.TGswAddH(cipherA, rgswParams)
	//cipherB.DebugTLweSample()
	//cipherA.DebugTGswSample(rgswParams)
	tfhe.TGswExternProduct(cipherAB, cipherA, cipherB, rgswParams)
	cipherAB.DebugTLweSample()
	fmt.Println("Test cipher after product 3.5 H*muB:")
<<<<<<< HEAD
	for p := int32(0); p <= k; p++ {
		for i := int32(0); i < N; i++ {
			expected := cipherB.A[p].Coefs[i]
			actual := cipherAB.A[p].Coefs[i]
=======
	for p := int(0); p <= k; p++ {
		for i := int(0); i < N; i++ {
			expected := cipherB.A[p].CoefsT[i]
			actual := cipherAB.A[p].CoefsT[i]
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
			if tfhe.Abs(expected-actual) > 10 {
				fmt.Printf("decompH error (p,i)=(%d,%d): %d != %d\n", p, i, actual, expected)
			}
			expected2 := tfhe.ModSwitchFromTorus(expected, int(Msize))
			actual2 := tfhe.ModSwitchFromTorus(actual, int(Msize))
			if expected2 != actual2 {
				fmt.Printf("modswitch error %d: %d != %d\n", i, actual2, expected2)
			}
		}
	}
	tfhe.TLweSymDecrypt(dechifAB, cipherAB, rlweKey, Msize) // DECRYPTION
	fmt.Println("Test LweSymDecrypt after product 3.5 H*muB:")
	fmt.Printf(" variance: %f", cipherAB.CurrentVariance)
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muB.Coefs[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifAB.Coefs[i], int(Msize))
=======
	for i := int(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus(muB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus(dechifAB.CoefsT[i], int(Msize))
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}
	fmt.Println("----------------------")

	//decryption test tgsw
	fmt.Println("decryption test tgsw:")
	tfhe.TGswSymEncrypt(cipherA, muA, alpha, rgswKey) // ENCRYPTION
	tfhe.TLwePhase(dechifB, cipherA.BlocSample[k][0], rlweKey)
	fmt.Println("manual decryption test: ")
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		//fmt.Printf("muA->c[i]: %d, dechifB->c[i]: %d\n", muA.Coefs[i], dechifB.Coefs[i])
		expected := muA.Coefs[i]
		actual := tfhe.ModSwitchFromTorus32(-512*dechifB.Coefs[i], 2)
=======
	for i := int(0); i < N; i++ {
		//fmt.Printf("muA->c[i]: %d, dechifB->c[i]: %d\n", muA.Coefs[i], dechifB.CoefsT[i])
		expected := muA.Coefs[i]
		actual := tfhe.ModSwitchFromTorus(-512*dechifB.CoefsT[i], 2)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		if expected != actual {
			fmt.Printf("tgsw encryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	tfhe.TGswSymDecrypt(dechifA, cipherA, rgswKey, int(Msize))
	fmt.Println("automatic decryption test:")
	for i := int(0); i < N; i++ {
		expected := muA.Coefs[i]
		actual := dechifA.Coefs[i]
		if expected != actual {
			fmt.Printf("tgsw decryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	tfhe.TorusPolynomialMulR(muAB, muA, muB)
	tfhe.TGswExternProduct(cipherAB, cipherA, cipherB, rgswParams)
	tfhe.TLweSymDecrypt(dechifAB, cipherAB, rlweKey, Msize) // DECRYPTION

	fmt.Println("Test LweSymDecrypt after product 3.5:")
	fmt.Printf(" variance: %f", cipherAB.CurrentVariance)
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muAB.Coefs[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifAB.Coefs[i], int(Msize))
=======
	for i := int(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus(muAB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus(dechifAB.CoefsT[i], int(Msize))
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}
	fmt.Println("----------------------")
}
