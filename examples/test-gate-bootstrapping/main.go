package main

import (
	"fmt"
	"math/rand"
	"time"

	tfhe "github.com/TheDonutFactory/go-tfhe"
)

func main() {

	const (
		nb_samples int32 = 64
		nb_trials        = 10
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
		// generate inputs (64-->127)
		for i := nb_samples; i < 2*nb_samples; i++ {
			tfhe.BootsSymEncrypt(test_in[i], rand.Int31()%2, keyset)
		}
		// fake encrypt
		tfhe.BootsSymEncrypt(test_in[0], rand.Int31()%2, keyset)

		// evaluate the NAND tree
		fmt.Printf("starting bootstrapping NAND tree...trial %d", trial)

		start := time.Now()
		for i := nb_samples - 1; i > 0; i-- {
			tfhe.BootsNAND(test_in[i], test_in[2*i], test_in[2*i+1], keyset.Cloud)
		}
		duration := time.Since(start)

		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Println(duration)

		fmt.Println("finished bootstrappings NAND tree")
		//cout << "time per bootNAND gate (microsecs)... " << (end - begin) / double(nb_samples-1) << endl

	}

}
