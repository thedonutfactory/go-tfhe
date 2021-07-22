package main

import (
	"fmt"
	"math"

	tfhe "github.com/thedonutfactory/fthe/go-tfhe"
)

const (
	N             int32   = 1024
	k             int32   = 1
	alpha_min_gsw float64 = 0.
	alpha_max_gsw float64 = 0.071
	Msize         int32   = 2
	alpha         float64 = 1e-6
	l             int     = 3
	Bgbits        int32   = 10
)

// **********************************************************************************
// ********************************* MAIN *******************************************
// **********************************************************************************
func approxEquals(a tfhe.Torus32, b tfhe.Torus32) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {
	for i := 0; i < 20000; i++ {
		tfhe.UniformTorus32Dist()
	}

	//TODO: parallelization
	//static uniform_int_distribution<int32_t> unift(0, Msize - 1)

	// PARAMETERS
	rlwe_params := tfhe.NewTLweParams(N, k, alpha_min_gsw, alpha_max_gsw) //les deux alpha mis un peu au hasard
	rgsw_params := tfhe.NewTGswParams(int32(l), Bgbits, rlwe_params)

	//tGswParams512_1  = NewTGswParams(4, 8, NewTLweParams(512, 1, 0., 1.))

	// KEY
	rgsw_key := tfhe.NewTGswKey(rgsw_params)
	rlwe_key := &rgsw_key.TlweKey
	// CIPHERTEXTS
	cipherA := tfhe.NewTGswSample(rgsw_params)
	cipherB := tfhe.NewTLweSample(rlwe_params)
	cipherAB := tfhe.NewTLweSample(rlwe_params)

	//the probability that a sample with stdev alpha decrypts wrongly on
	//the a Msize message space.

	expected_error_proba := 1. - math.Erf(1./(math.Sqrt(2.)*2.*float64(Msize)*alpha))

	fmt.Println("-------------")
	fmt.Println("WARNING:")
	fmt.Printf("All the tests below are supposed to fail with proba: %f \n", expected_error_proba)
	fmt.Println("It is normal and it is part of the test!")
	fmt.Println("-------------")

	//MESSAGE RLwe
	muB := tfhe.NewTorusPolynomial(N)

	//test decompH
	fmt.Println("Test decompH on TorusPolynomial")
	muBDecH := tfhe.NewIntPolynomialArray(l, N)
	for i := int32(0); i < N; i++ {
		muB.CoefsT[i] = tfhe.UniformTorus32Dist()
	}
	tfhe.TGswTorus32PolynomialDecompH(muBDecH, muB, rgsw_params)
	for i := int32(0); i < N; i++ {
		expected := muB.CoefsT[i]
		var actual int32 = 0
		for j := 0; j < l; j++ {
			actual += muBDecH[j].Coefs[i] * rgsw_params.H[j]
			//fmt.Printf("DEBUG: l: %d, i: %d, j: %d, muBDecH[j].Coefs[i]: %d, rgsw_params->h[j]: %d \n", l, i, j, muBDecH[j].Coefs[i], rgsw_params.H[j])
		}
		//jcl
		//fmt.Printf("\t DEBUG: actual: %d, expected: %d\n", actual, expected)
		if tfhe.Abs(expected-actual) > 3 {
			fmt.Printf("decompH error %d: %d != %d\n", i, actual, expected)
		}
	}

	for i := int32(0); i < N; i++ {
		temp := tfhe.UniformInt32Dist(0, Msize-1) //unift(generator)
		muB.CoefsT[i] = tfhe.ModSwitchToTorus32(temp, Msize)
		//cout << mu.CoefsT[i] << endl
	}
	//MESSAGE RLwe
	muA := tfhe.NewIntPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := tfhe.UniformInt32Dist(0, Msize-1)
		muA.Coefs[i] = 1 - (temp % 3)
		//cout << mu.CoefsT[i] << endl
	}
	// PHASE, DECRYPTION
	dechifA := tfhe.NewIntPolynomial(N)
	dechifB := tfhe.NewTorusPolynomial(N)
	dechifAB := tfhe.NewTorusPolynomial(N)
	muAB := tfhe.NewTorusPolynomial(N)

	tfhe.TGswKeyGen(rgsw_key)                          // KEY GENERATION
	tfhe.TLweSymEncrypt(cipherB, muB, alpha, rlwe_key) // ENCRYPTION

	//decryption test tlwe
	fmt.Println("Test TLweSymDecrypt on muB:")
	fmt.Printf(" variance: %f\n", cipherB.CurrentVariance)
	tfhe.TLweSymDecrypt(dechifB, cipherB, rlwe_key, Msize) // DECRYPTION
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
	tfhe.TGswTLweDecompH(cipherBDecH, cipherB, rgsw_params)
	for p := int32(0); p <= k; p++ {
		for i := int32(0); i < N; i++ {
			expected := cipherB.A[p].CoefsT[i]
			var actual int32 = 0
			for j := int32(0); j < int32(l); j++ {
				x := int32(l)*p + j
				// fmt.Printf("DEBUG: l: %d, p: %d, j: %d, x: %d, i: %d, cipherBDecH[x].coefs[i]: %d, rgsw_params->h[j]: %d \n", l, p, j, x, i, cipherBDecH[x].Coefs[i], rgsw_params.H[j])
				actual += cipherBDecH[x].Coefs[i] * rgsw_params.H[j]
				// actual += cipherBDecH[l*int(p)+j].Coefs[i] * rgsw_params.H[j]
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
	tfhe.TGswClear(cipherA, rgsw_params)
	tfhe.TGswAddH(cipherA, rgsw_params)
	tfhe.TGswExternProduct(cipherAB, cipherA, cipherB, rgsw_params)
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
	tfhe.TLweSymDecrypt(dechifAB, cipherAB, rlwe_key, Msize) // DECRYPTION
	fmt.Println("Test LweSymDecrypt after product 3.5 H*muB:")
	fmt.Printf(" variance: %s", cipherAB.CurrentVariance)
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
	tfhe.TGswSymEncrypt(cipherA, muA, alpha, rgsw_key) // ENCRYPTION
	tfhe.TLwePhase(dechifB, &cipherA.BlocSample[k][0], rlwe_key)
	fmt.Println("manual decryption test: ")
	for i := int32(0); i < N; i++ {
		//fmt.Printf("muA->c[i]: %d, dechifB->c[i]: %d\n", muA.Coefs[i], dechifB.CoefsT[i])
		expected := muA.Coefs[i]
		actual := tfhe.ModSwitchFromTorus32(-512*dechifB.CoefsT[i], 2)
		if expected != actual {
			fmt.Printf("tgsw encryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	tfhe.TGswSymDecrypt(dechifA, cipherA, rgsw_key, int(Msize))
	fmt.Println("automatic decryption test:")
	for i := int32(0); i < N; i++ {
		expected := muA.Coefs[i]
		actual := dechifA.Coefs[i]
		if expected != actual {
			fmt.Printf("tgsw decryption error %d: %d != %d\n", i, actual, expected)
		}
	}

	tfhe.TorusPolynomialMulR(muAB, muA, muB)
	tfhe.TGswExternProduct(cipherAB, cipherA, cipherB, rgsw_params)
	tfhe.TLweSymDecrypt(dechifAB, cipherAB, rlwe_key, Msize) // DECRYPTION

	fmt.Println("Test LweSymDecrypt after product 3.5:")
	fmt.Printf(" variance: %s", cipherAB.CurrentVariance)
	for i := int32(0); i < N; i++ {
		expected := tfhe.ModSwitchFromTorus32(muAB.CoefsT[i], int(Msize))
		actual := tfhe.ModSwitchFromTorus32(dechifAB.CoefsT[i], int(Msize))
		if expected != actual {
			fmt.Printf("tlwe decryption error %d: %d != %d\n", i, actual, expected)
		}
	}
	fmt.Println("----------------------")
}
