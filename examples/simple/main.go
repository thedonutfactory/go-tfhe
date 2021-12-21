package main

import (
	"fmt"
	"os"

	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/io"
)

func keys(params *gates.GateBootstrappingParameterSet) (*gates.PublicKey, *gates.PrivateKey) {
	var pubKey *gates.PublicKey
	var privKey *gates.PrivateKey
	if _, err := os.Stat("private.key"); err == nil {
		fmt.Println("------ Reading keys from file ------")
		privKey, _ = io.ReadPrivKey("private.key")
		pubKey, _ = io.ReadPubKey("public.key")

	} else {
		fmt.Println("------ Key Generation ------")
		// generate the keys
		pubKey, privKey = params.GenerateKeys()
		io.WritePrivKey(privKey, "private.key")
		io.WritePubKey(pubKey, "public.key")
	}
	return pubKey, privKey
}

func main() {
	// generate public and private keys
	ctx := gates.DefaultGateBootstrappingParameters(100)
	pub, prv := ctx.GenerateKeys()

	// encrypt 2 8-bit ciphertexts
	x := prv.Encrypt(int8(22))
	y := prv.Encrypt(int8(33))

	// create ciphertext variables
	temp := ctx.Int8()
	sum := ctx.Int8()
	carry := ctx.Int8()

	// perform homomorphic gate operations
	temp[0] = pub.Xor(x[0], y[0])
	sum[0] = pub.Xor(temp[0], carry[0])
	temp[1] = pub.And(x[0], y[0])
	temp[2] = pub.And(carry[0], temp[0])
	carry[1] = pub.Xor(temp[1], temp[2])
	carry[0] = pub.Copy(carry[1])

	// decrypt results
	z := prv.Decrypt(sum[:])
	fmt.Println("The result is: ", z)
}
