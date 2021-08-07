package main

import (
	"fmt"
	"math/rand"
	"time"

	tfhe "github.com/thedonutfactory/go-tfhe"
)

func main() {

	const (
<<<<<<< HEAD
		//nbSamples int32 = 64
		//nbTrials        = 10
		nbSamples int32 = 64
		nbTrials        = 1
=======
		//nb_samples int32 = 64
		//nb_trials        = 10
		nb_samples int32 = 64
		nb_trials        = 1
>>>>>>> caa54e1 (Completed initial bootstrapped gates - not FFT optimized)
	)
	// generate params
	var minimumLambda int32 = 100
	params := tfhe.NewDefaultGateBootstrappingParameters(minimumLambda)
	inOutParams := params.InOutParams
	// generate the secret keyset
	keyset := tfhe.NewRandomGateBootstrappingSecretKeyset(params)

	for trial := 0; trial < nbTrials; trial++ {

		// generate samples
<<<<<<< HEAD
		testIn := tfhe.NewLweSampleArray(2*nbSamples, inOutParams)
		// generate inputs (64-.127)
		for i := nbSamples; i < 2*nbSamples; i++ {
			tfhe.BootsSymEncrypt(testIn[i], rand.Int31()%2, keyset)
=======
		test_in := tfhe.NewLweSampleArray(2*nb_samples, in_out_params)
		// generate inputs (64-.127)
		for i := nb_samples; i < 2*nb_samples; i++ {
			tfhe.BootsSymEncrypt(test_in[i], rand.Int31()%2, keyset)
>>>>>>> caa54e1 (Completed initial bootstrapped gates - not FFT optimized)
		}
		// fake encrypt
		tfhe.BootsSymEncrypt(testIn[0], rand.Int31()%2, keyset)

		// evaluate the NAND tree
		fmt.Printf("starting bootstrapping NAND tree...trial %d\n", trial)

		start := time.Now()
		for i := nbSamples - 1; i > 0; i-- {
			tfhe.BootsNAND(testIn[i], testIn[2*i], testIn[2*i+1], keyset.Cloud)
		}
		duration := time.Since(start)

		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Println(duration)

		fmt.Println("finished bootstrappings NAND tree")
		fmt.Printf("time per bootNAND gate: %s", duration)

		// verification
<<<<<<< HEAD
		for i := nbSamples - 1; i > 0; i-- {
			mess1 := tfhe.BootsSymDecrypt(testIn[2*i], keyset)
			mess2 := tfhe.BootsSymDecrypt(testIn[2*i+1], keyset)
			out := tfhe.BootsSymDecrypt(testIn[i], keyset)
=======
		for i := nb_samples - 1; i > 0; i-- {
			mess1 := tfhe.BootsSymDecrypt(test_in[2*i], keyset)
			mess2 := tfhe.BootsSymDecrypt(test_in[2*i+1], keyset)
			out := tfhe.BootsSymDecrypt(test_in[i], keyset)
>>>>>>> caa54e1 (Completed initial bootstrapped gates - not FFT optimized)

			if out != 1-(mess1*mess2) {
				fmt.Printf("Error - trial %d [ %d ]", trial, i)

				fmt.Printf("%f - %f - %f \n",
<<<<<<< HEAD
					tfhe.T32tod(tfhe.LwePhase(testIn[i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(testIn[2*i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(testIn[2*i+1], keyset.LweKey)),
=======
					tfhe.T32tod(tfhe.LwePhase(test_in[i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(test_in[2*i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(test_in[2*i+1], keyset.LweKey)),
>>>>>>> caa54e1 (Completed initial bootstrapped gates - not FFT optimized)
				)
			}
		}
	}
}
