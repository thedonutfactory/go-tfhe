package main

import (
	"fmt"
	"math"

	tfhe "github.com/TheDonutFactory/go-tfhe"
)

func approxEquals(a tfhe.Torus32, b tfhe.Torus32) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {

	const (
		N        int32 = 1024
		k        int32 = 2
		alphaMin       = 0.01
		alphaMax       = 0.071
		Msize    int32 = 7 // taille de l'espace des coeffs du polynome du message
		alpha          = 0.02
	)

	unift := tfhe.NewUniform(0, Msize-1)

	// PARAMETERS
	params := tfhe.NewTLweParams(N, k, alphaMin, alphaMax) //les deux alpha mis un peu au hasard
	// KEY
	key := tfhe.NewTLweKey(params)
	// CIPHERTEXTS
	cipher := tfhe.NewTLweSample(params)
	cipherT := tfhe.NewTLweSample(params)

	//the probability that a sample with stdev alpha decrypts wrongly on
	//the a Msize message space.
	expectedErrorProba := 1. - math.Erf(1./(math.Sqrt(2.)*2.*float64(Msize)*alpha))

	fmt.Println("-------------")
	fmt.Println("WARNING:")
	fmt.Printf("All the tests below are supposed to fail with proba: %f\n", expectedErrorProba)
	fmt.Println("It is normal and it is part of the test!")
	fmt.Println("-------------")

	//MESSAGE
	mu := tfhe.NewTorusPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := unift.Int32() // distribution(generator);
		mu.CoefsT[i] = tfhe.ModSwitchToTorus32(temp, Msize)
		//cout << mu.coefsT[i] << endl;
	}
	// PHASE, DECRYPTION
	phi := tfhe.NewTorusPolynomial(N)
	dechif := tfhe.NewTorusPolynomial(N)

	tfhe.TLweKeyGen(key)                            // KEY GENERATION
	tfhe.TLweSymEncrypt(cipher, mu, alpha, key)     // ENCRYPTION
	tfhe.TLwePhase(phi, cipher, key)                // PHASE COMUPTATION
	tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

	fmt.Println("Test LweSymDecrypt :")
	for i := int32(0); i < N; i++ {
		if dechif.CoefsT[i] != mu.CoefsT[i] {
			fmt.Printf("%d - %d =? %d error!!!\n", i, dechif.CoefsT[i], mu.CoefsT[i])
		}
	}
	fmt.Println("----------------------")

	phiT := tfhe.NewTorusPolynomial(N)
	for trial := 1; trial < 1000; trial++ {
		muT := tfhe.ModSwitchToTorus32(unift.Int32(), Msize)
		var dechifT tfhe.Torus32 = 0

		tfhe.TLweSymEncryptT(cipherT, muT, alpha, key)
		tfhe.TLwePhase(phiT, cipherT, key)
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize)

		if dechifT != muT {
			fmt.Printf("Test LweSymDecryptT: trial %d\n", trial)

			fmt.Printf("%d =? %d Error!!!\n", dechifT, muT)
			fmt.Println("----------------------")
		}
	}

	// TEST ADD, SUB, LINEAR COMBINATION, POLYNOMIAL COMBINATIONS

	fmt.Println()
	fmt.Println()
	fmt.Println("----------------------")
	fmt.Println("TEST Operations TLwe :")
	fmt.Println("-------------------------")

	// CIPHERTEXTS
	cipher0 := tfhe.NewTLweSample(params)
	cipher1 := tfhe.NewTLweSample(params)

	// MESSAGES
	mu0 := tfhe.NewTorusPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		mu0.CoefsT[i] = tfhe.ModSwitchToTorus32(temp, Msize)
	}
	mu1 := tfhe.NewTorusPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		mu1.CoefsT[i] = tfhe.ModSwitchToTorus32(temp, Msize)
	}
	var p int32 = 1
	poly := tfhe.NewIntPolynomial(N)
	for i := int32(0); i < N; i++ {
		poly.Coefs[i] = unift.Int32()
	}

	var decInt int32 = 0
	var muInt int32 = 0

	for trial := 1; trial < 2; trial++ {

		tfhe.TLweSymEncrypt(cipher0, mu0, alpha, key) // ENCRYPTION
		tfhe.TLweSymEncrypt(cipher1, mu1, alpha, key) // ENCRYPTION

		// cipher = cipher0 + cipher1
		tfhe.TLweCopy(cipher, cipher0, params)
		tfhe.TLweAddTo(cipher, cipher1, params)
		tfhe.TorusPolynomialAdd(mu, mu0, mu1)           // mu = mu0 + mu1
		tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

		fmt.Printf("Test tLweAddTo Trial : %d\n", trial)
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.CoefsT[i], Msize)
			if decInt != muInt {
				fmt.Printf("%d =? %d error !!!\n", decInt, muInt)
			}
		}
		fmt.Println(cipher.CurrentVariance)
		fmt.Println("----------------------")

		// cipher = cipher0 - cipher1
		tfhe.TLweCopy(cipher, cipher0, params)
		tfhe.TLweSubTo(cipher, cipher1, params)
		tfhe.TorusPolynomialSub(mu, mu0, mu1)           // mu = mu0 - mu1
		tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

		fmt.Printf("Test tLweSubTo Trial : %d\n", trial)
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.CoefsT[i], Msize)
			if decInt != muInt {
				fmt.Printf("%d =? %d error !!!\n", decInt, muInt)
			}
		}
		fmt.Println(cipher.CurrentVariance)
		fmt.Println("----------------------")

		// cipher = cipher0 + p.cipher1
		tfhe.TLweCopy(cipher, cipher0, params)
		tfhe.TLweAddMulTo(cipher, p, cipher1, params)
		tfhe.TorusPolynomialAddMulZ(mu, mu0, p, mu1)    // mu = mu0 + p.mu1
		tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

		fmt.Printf("Test tLweAddMulTo Trial : %d\n", trial)
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.CoefsT[i], Msize)
			if decInt != muInt {
				fmt.Printf("%d =? %d error !!!\n", decInt, muInt)
			}
		}
		fmt.Println(cipher.CurrentVariance)
		fmt.Println("----------------------")

		// cipher = cipher0 - p.cipher1
		tfhe.TLweCopy(cipher, cipher0, params)
		tfhe.TLweSubMulTo(cipher, p, cipher1, params)
		tfhe.TorusPolynomialSubMulZ(mu, mu0, p, mu1)    // mu = mu0 - p.mu1
		tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

		fmt.Printf("Test tLweSubMulTo Trial : %d\n", trial)
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.CoefsT[i], Msize)
			if decInt != muInt {
				fmt.Printf("%d =? %d error !!!\n", decInt, muInt)
			}
		}
		fmt.Println(cipher.CurrentVariance)
		fmt.Println("----------------------")
	}

	// TEST ADD, SUB, LINEAR COMBINATION, POLYNOMIAL COMBINATIONS

	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------")
	fmt.Println("TEST Operations TLwe with Torus32 messages :")
	fmt.Println("-----------------------------------------------")

	// CIPHERTEXTS
	cipherT0 := tfhe.NewTLweSample(params)
	cipherT1 := tfhe.NewTLweSample(params)
	var pT int32 = 1

	for trial := 1; trial < 1000; trial++ {

		// MESSAGES
		muT0 := tfhe.ModSwitchToTorus32(unift.Int32(), Msize)
		muT1 := tfhe.ModSwitchToTorus32(unift.Int32(), Msize)

		var muT tfhe.Torus32 = 0
		var dechifT tfhe.Torus32 = 0

		tfhe.TLweSymEncryptT(cipherT0, muT0, alpha, key) // ENCRYPTION
		tfhe.TLweSymEncryptT(cipherT1, muT1, alpha, key) // ENCRYPTION

		// cipher = cipher0 + cipher1
		tfhe.TLweCopy(cipherT, cipherT0, params)
		tfhe.TLweAddTo(cipherT, cipherT1, params)
		muT = muT0 + muT1
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize) // DECRYPTION

		decInt = tfhe.ModSwitchFromTorus32(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus32(muT, Msize)
		if decInt != muInt {
			fmt.Printf("Test tLweAddTo Trial : %d\n", trial)
			fmt.Printf("%d =? %d Error!!!\n", decInt, muInt)
			fmt.Println(cipherT.CurrentVariance)
			fmt.Println("----------------------")
		}

		// cipher = cipher0 - cipher1
		tfhe.TLweCopy(cipherT, cipherT0, params)
		tfhe.TLweSubTo(cipherT, cipherT1, params)
		muT = muT0 - muT1
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize) // DECRYPTION

		decInt = tfhe.ModSwitchFromTorus32(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus32(muT, Msize)
		if decInt != muInt {
			fmt.Printf("Test tLweSubTo Trial : %d\n", trial)
			fmt.Printf("%d =? %d Error!!!\n", decInt, muInt)
			fmt.Println(cipherT.CurrentVariance)
			fmt.Println("----------------------")
		}

		// cipher = cipher0 + p.cipher1
		tfhe.TLweCopy(cipherT, cipherT0, params)
		tfhe.TLweAddMulTo(cipherT, pT, cipherT1, params)
		muT = muT0 + pT*muT1
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize) // DECRYPTION

		decInt = tfhe.ModSwitchFromTorus32(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus32(muT, Msize)
		if decInt != muInt {
			fmt.Printf("Test tLweAddMulTo Trial : %d\n", trial)
			fmt.Printf("%d =? %d Error!!!\n", decInt, muInt)
			fmt.Println(cipherT.CurrentVariance)
			fmt.Println("----------------------")
		}

		// result = result - p.sample
		tfhe.TLweCopy(cipherT, cipherT0, params)
		tfhe.TLweSubMulTo(cipherT, pT, cipherT1, params)
		muT = muT0 - pT*muT1
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize) // DECRYPTION

		decInt = tfhe.ModSwitchFromTorus32(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus32(muT, Msize)
		if decInt != muInt {
			fmt.Printf("Test tLweAddMulTo Trial : %d\n", trial)
			fmt.Printf("%d =? %d Error!!!\n", decInt, muInt)
			fmt.Println(cipherT.CurrentVariance)
			fmt.Println("----------------------")
		}
	}

}
