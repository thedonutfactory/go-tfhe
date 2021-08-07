package main

import (
	"fmt"
	"math/rand"
	"time"

	tfhe "github.com/TheDonutFactory/go-tfhe"
)

func main() {

	const (
		//nb_samples int32 = 64
		//nb_trials        = 10
		nb_samples int32 = 64
		nb_trials        = 1
	)
	// generate params
	var minimum_lambda int32 = 100
	params := tfhe.NewDefaultGateBootstrappingParameters(minimum_lambda)
	in_out_params := params.InOutParams
	// generate the secret keyset
	keyset := tfhe.NewRandomGateBootstrappingSecretKeyset(params)

	for trial := 0; trial < nb_trials; trial++ {

		// generate samples
		test_in := tfhe.NewLweSampleArray(2*nb_samples, in_out_params)
		// generate inputs (64-.127)
		for i := nb_samples; i < 2*nb_samples; i++ {
			tfhe.BootsSymEncrypt(test_in[i], rand.Int31()%2, keyset)
		}
		// fake encrypt
		tfhe.BootsSymEncrypt(test_in[0], rand.Int31()%2, keyset)

		// evaluate the NAND tree
		fmt.Printf("starting bootstrapping NAND tree...trial %d\n", trial)

		start := time.Now()
		for i := nb_samples - 1; i > 0; i-- {
			tfhe.BootsNAND(test_in[i], test_in[2*i], test_in[2*i+1], keyset.Cloud)
		}
		duration := time.Since(start)

		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Println(duration)

		fmt.Println("finished bootstrappings NAND tree")
		fmt.Printf("time per bootNAND gate: %s", duration)

		// verification
		for i := nb_samples - 1; i > 0; i-- {
			mess1 := tfhe.BootsSymDecrypt(test_in[2*i], keyset)
			mess2 := tfhe.BootsSymDecrypt(test_in[2*i+1], keyset)
			out := tfhe.BootsSymDecrypt(test_in[i], keyset)

			if out != 1-(mess1*mess2) {
				fmt.Printf("Error - trial %d [ %d ]", trial, i)

				fmt.Printf("%f - %f - %f \n",
					tfhe.T32tod(tfhe.LwePhase(test_in[i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(test_in[2*i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(test_in[2*i+1], keyset.LweKey)),
				)
			}
		}
	}
}
