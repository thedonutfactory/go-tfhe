package main

import (
	"fmt"
	"math"

	tfhe "github.com/thedonutfactory/go-tfhe"
)

func approxEquals(a tfhe.Torus, b tfhe.Torus) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {

	const (
		N        int = 1024
		k        int = 2
		alphaMin     = 0.01
		alphaMax     = 0.071
		Msize    int = 7 // taille de l'espace des coeffs du polynome du message
		alpha        = 0.02
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
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		temp := unift.Int32() // distribution(generator);
		mu.Coefs[i] = tfhe.ModSwitchToTorus32(temp, Msize)
		//cout << mu.Coefs[i] << endl;
=======
	for i := int(0); i < N; i++ {
		temp := unift.int() // distribution(generator);
		mu.CoefsT[i] = tfhe.ModSwitchToTorus(temp, Msize)
		//cout << mu.coefsT[i] << endl;
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
	}
	// PHASE, DECRYPTION
	phi := tfhe.NewTorusPolynomial(N)
	dechif := tfhe.NewTorusPolynomial(N)

	tfhe.TLweKeyGen(key)                            // KEY GENERATION
	tfhe.TLweSymEncrypt(cipher, mu, alpha, key)     // ENCRYPTION
	tfhe.TLwePhase(phi, cipher, key)                // PHASE COMUPTATION
	tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

	fmt.Println("Test LweSymDecrypt :")
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		if dechif.Coefs[i] != mu.Coefs[i] {
			fmt.Printf("%d - %d =? %d error!!!\n", i, dechif.Coefs[i], mu.Coefs[i])
=======
	for i := int(0); i < N; i++ {
		if dechif.CoefsT[i] != mu.CoefsT[i] {
			fmt.Printf("%d - %d =? %d error!!!\n", i, dechif.CoefsT[i], mu.CoefsT[i])
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		}
	}
	fmt.Println("----------------------")

	phiT := tfhe.NewTorusPolynomial(N)
	for trial := 1; trial < 1000; trial++ {
		muT := tfhe.ModSwitchToTorus(unift.int(), Msize)
		var dechifT tfhe.Torus = 0

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
<<<<<<< HEAD
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		mu0.Coefs[i] = tfhe.ModSwitchToTorus32(temp, Msize)
	}
	mu1 := tfhe.NewTorusPolynomial(N)
	for i := int32(0); i < N; i++ {
		temp := unift.Int32()
		mu1.Coefs[i] = tfhe.ModSwitchToTorus32(temp, Msize)
=======
	for i := int(0); i < N; i++ {
		temp := unift.int()
		mu0.CoefsT[i] = tfhe.ModSwitchToTorus(temp, Msize)
	}
	mu1 := tfhe.NewTorusPolynomial(N)
	for i := int(0); i < N; i++ {
		temp := unift.int()
		mu1.CoefsT[i] = tfhe.ModSwitchToTorus(temp, Msize)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
	}
	var p int = 1
	poly := tfhe.NewIntPolynomial(N)
	for i := int(0); i < N; i++ {
		poly.Coefs[i] = unift.int()
	}

	var decInt int = 0
	var muInt int = 0

	for trial := 1; trial < 2; trial++ {

		tfhe.TLweSymEncrypt(cipher0, mu0, alpha, key) // ENCRYPTION
		tfhe.TLweSymEncrypt(cipher1, mu1, alpha, key) // ENCRYPTION

		// cipher = cipher0 + cipher1
		tfhe.TLweCopy(cipher, cipher0, params)
		tfhe.TLweAddTo(cipher, cipher1, params)
		tfhe.TorusPolynomialAdd(mu, mu0, mu1)           // mu = mu0 + mu1
		tfhe.TLweSymDecrypt(dechif, cipher, key, Msize) // DECRYPTION

		fmt.Printf("Test tLweAddTo Trial : %d\n", trial)
<<<<<<< HEAD
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.Coefs[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.Coefs[i], Msize)
=======
		for i := int(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus(mu.CoefsT[i], Msize)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
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
<<<<<<< HEAD
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.Coefs[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.Coefs[i], Msize)
=======
		for i := int(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus(mu.CoefsT[i], Msize)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
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
<<<<<<< HEAD
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.Coefs[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.Coefs[i], Msize)
=======
		for i := int(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus(mu.CoefsT[i], Msize)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
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
<<<<<<< HEAD
		for i := int32(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus32(dechif.Coefs[i], Msize)
			muInt = tfhe.ModSwitchFromTorus32(mu.Coefs[i], Msize)
=======
		for i := int(0); i < N; i++ {
			decInt = tfhe.ModSwitchFromTorus(dechif.CoefsT[i], Msize)
			muInt = tfhe.ModSwitchFromTorus(mu.CoefsT[i], Msize)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
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
	fmt.Println("TEST Operations TLwe with Torus messages :")
	fmt.Println("-----------------------------------------------")

	// CIPHERTEXTS
	cipherT0 := tfhe.NewTLweSample(params)
	cipherT1 := tfhe.NewTLweSample(params)
	var pT int = 1

	for trial := 1; trial < 1000; trial++ {

		// MESSAGES
		muT0 := tfhe.ModSwitchToTorus(unift.int(), Msize)
		muT1 := tfhe.ModSwitchToTorus(unift.int(), Msize)

		var muT tfhe.Torus = 0
		var dechifT tfhe.Torus = 0

		tfhe.TLweSymEncryptT(cipherT0, muT0, alpha, key) // ENCRYPTION
		tfhe.TLweSymEncryptT(cipherT1, muT1, alpha, key) // ENCRYPTION

		// cipher = cipher0 + cipher1
		tfhe.TLweCopy(cipherT, cipherT0, params)
		tfhe.TLweAddTo(cipherT, cipherT1, params)
		muT = muT0 + muT1
		dechifT = tfhe.TLweSymDecryptT(cipherT, key, Msize) // DECRYPTION

		decInt = tfhe.ModSwitchFromTorus(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus(muT, Msize)
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

		decInt = tfhe.ModSwitchFromTorus(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus(muT, Msize)
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

		decInt = tfhe.ModSwitchFromTorus(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus(muT, Msize)
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

		decInt = tfhe.ModSwitchFromTorus(dechifT, Msize)
		muInt = tfhe.ModSwitchFromTorus(muT, Msize)
		if decInt != muInt {
			fmt.Printf("Test tLweAddMulTo Trial : %d\n", trial)
			fmt.Printf("%d =? %d Error!!!\n", decInt, muInt)
			fmt.Println(cipherT.CurrentVariance)
			fmt.Println("----------------------")
		}
	}

}
