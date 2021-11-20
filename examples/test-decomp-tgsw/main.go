package main

import (
	"fmt"
	"math"

	tfhe "github.com/thedonutfactory/go-tfhe"
)

const (
	N           int32   = 1024
	k           int32   = 1
	alphaMinGsw float64 = 0.
	alphaMaxGsw float64 = 0.071
	Msize       int32   = 2
	alpha       float64 = 1e-6
	l           int     = 3
	Bgbits      int32   = 10
)

// **********************************************************************************
// ********************************* MAIN *******************************************
// **********************************************************************************
func approxEquals(a tfhe.Torus32, b tfhe.Torus32) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {

	//TODO: parallelization
	unift := tfhe.NewUniform(0, Msize-1)

	// PARAMETERS
	rlweParams := tfhe.NewTLweParams(N, k, alphaMinGsw, alphaMaxGsw) //les deux alpha mis un peu au hasard
	rgswParams := tfhe.NewTGswParams(int32(l), Bgbits, rlweParams)

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
	for i := int32(0); i < N; i++ {
		muB.CoefsT[i] = unift.Int32() //tfhe.UniformTorus32Dist()
	}
	tfhe.TGswTorus32PolynomialDecompH(muBDecH, muB, rgswParams)
	for i := int32(0); i < N; i++ {
		expected := muB.CoefsT[i]
		var actual int32 = 0
		for j := 0; j < l; j++ {
			actual += muBDecH[j].Coefs[i] * rgswParams.H[j]
			//fmt.Printf("DEBUG: l: %d, i: %d, j: %d, muBDecH[j].Coefs[i]: %d, rgswParams->h[j]: %d \n", l, i, j, muBDecH[j].Coefs[i], rgswParams.H[j])
		}
		//fmt.Printf("\t DEBUG: actual: %d, expected: %d\n", actual, expected)
		if tfhe.Abs(expected-actual) > 3 {
			fmt.Printf("decompH error %d: %d != %d\n", i, actual, expected)
		}
	}

	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		muB.CoefsT[i] = tfhe.ModSwitchToTorus32(temp, Msize)
		// fmt.Println(muB.CoefsT[i])
	}
	//MESSAGE RLwe
	muA := tfhe.NewIntPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
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
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifB.CoefsT[i], int(Msize))
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	//test decompH on tLwe
	fmt.Println("Test decompH on TLwe(muB)")
	cipherBDecH := tfhe.NewIntPolynomialArray(l*(int(k)+1), N)
	tfhe.TGswTLweDecompH(cipherBDecH, cipherB, rgswParams)
	for p := int32(0); p <= k; p++ {
		for i := int32(0); i < N; i++ {
			expected := cipherB.A[p].CoefsT[i]
			var actual int32 = 0
			for j := int32(0); j < int32(l); j++ {
				x := int32(l)*p + j
				actual += cipherBDecH[x].Coefs[i] * rgswParams.H[j]
			}
			// fails when p == 1, the array is not being populated properly
			if tfhe.Abs(expected-actual) > 3 {
				fmt.Printf("decompH error (p,i)=(%d,%d): %d != %d\n", p, i, actual, expected)
			}
			//jcl
			//fmt.Printf("\t DEBUG: actual: %d, expected: %d\n", actual, expected)
			expected2 := tfhe.ModSwitchFromTorus32(expected, int(Msize))
			actual2 := tfhe.ModSwitchFromTorus32(actual, int(Msize))
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
	for p := int32(0); p <= k; p++ {
		for i := int32(0); i < N; i++ {
			expected := cipherB.A[p].CoefsT[i]
			actual := cipherAB.A[p].CoefsT[i]
			if tfhe.Abs(expected-actual) > 10 {
				fmt.Printf("decompH error (p,i)=(%d,%d): %d != %d\n", p, i, actual, expected)
			}
			expected2 := tfhe.ModSwitchFromTorus32(expected, int(Msize))
			actual2 := tfhe.ModSwitchFromTorus32(actual, int(Msize))
			if expected2 != actual2 {
				fmt.Printf("modswitch error %d: %d != %d\n", i, actual2, expected2)
			}
		}
	}
	tfhe.TLweSymDecrypt(dechifAB, cipherAB, rlweKey, Msize) // DECRYPTION
	fmt.Println("Test LweSymDecrypt after product 3.5 H*muB:")
	fmt.Printf(" variance: %f", cipherAB.CurrentVariance)
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifAB.CoefsT[i], int(Msize))
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
	for i := int32(0); i < N; i++ {
		//fmt.Printf("muA->c[i]: %d, dechifB->c[i]: %d\n", muA.Coefs[i], dechifB.CoefsT[i])
		expected := muA.Coefs[i]
		actual := tfhe.ModSwitchFromTorus32(-512*dechifB.CoefsT[i], 2)
		if expected != actual {
			fmt.Printf("tgsw encryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	tfhe.TGswSymDecrypt(dechifA, cipherA, rgswKey, int(Msize))
	fmt.Println("automatic decryption test:")
	for i := int32(0); i < N; i++ {
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
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muAB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifAB.CoefsT[i], int(Msize))
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}
	fmt.Println("----------------------")
}
