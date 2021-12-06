package main

import (
	"fmt"

	tfhe "github.com/thedonutfactory/tfhe/go-tfhe"
)

func approxEquals(a tfhe.Torus, b tfhe.Torus) bool {
	return tfhe.Abs(a-b) < 10
}

func main() {

	params := tfhe.NewLweParams(512, 0.2, 0.5) //les deux alpha mis un peu au hasard
	n := params.N
	key := tfhe.NewLweKey(params)
	cipher := tfhe.NewLweSample(params)
	mu := tfhe.Dtot32(0.5)
	//Attention, 1<<30 correspond au message 0.25!! Ila: t'as raison!
	alpha := 0.0625
	var Msize int64 = 2

	tfhe.LweKeyGen(key)
	tfhe.LweSymEncrypt(cipher, mu, alpha, key)
	fmt.Print("a = [")
	for i := int64(0); i < n-1; i++ {
		fmt.Printf("%d, ", cipher.A[i])
		//cout << T32tod(cipher->a[i]) << ", ";
	}
	fmt.Printf("%d] \n", cipher.A[n-1])
	//cout << T32tod(cipher->a[n - 1]) << "]" << endl;
	fmt.Printf("b = %f \n", tfhe.T32tod(cipher.B))
	//cout << "b = " << T32tod(cipher->b) << endl;

	phi := tfhe.LwePhase(cipher, key)
	fmt.Printf("phi = %f \n", tfhe.T32tod(phi))
	//cout << "phi = " << T32tod(phi) << endl;
	message := tfhe.LweSymDecrypt(cipher, key, Msize)
	fmt.Printf("message = %f \n", tfhe.T32tod(message))
	//cout << "message = " << T32tod(message) << endl;

	//lwe crash test
	var failures int64 = 0
	var trials int64 = 1000
	for i := int64(0); i < trials; i++ {
		input := tfhe.Dtot32(float64((i % 3) / 3.))
		tfhe.LweKeyGen(key)
		tfhe.LweSymEncrypt(cipher, input, 0.047, key) // Ila: attention au niveau de bruit!!! à voir (0.06 n'est pas le bon je crois, 0.047 marche parfaitement)
		phi = tfhe.LwePhase(cipher, key)
		decrypted := tfhe.LweSymDecrypt(cipher, key, 3)
		if !approxEquals(input, decrypted) {
			fmt.Errorf("WARNING: the msg %f gave phase %f and was incorrectly decrypted to %f \n", tfhe.T32tod(input), tfhe.T32tod(phi), tfhe.T32tod(decrypted))
			failures++
		}
	}
	fmt.Printf("There were %d failures out of %d trials \n", failures, trials)
	fmt.Println("(it might be normal)")

}
