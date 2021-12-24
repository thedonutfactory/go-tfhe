package main

import (
	"fmt"
	"os"
	"time"

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
	ctx := gates.NewDefaultGateBootstrappingParameters()
	pub, prv := ctx.GenerateKeys()

	// encrypt 2 8-bit ciphertexts
	x := prv.Encrypt(int8(22))
	y := prv.Encrypt(int8(33))

	start := time.Now()
	// perform homomorphic sum gate operations
	BITS := 8
	temp := ctx.Int(3)
	sum := ctx.Int(BITS + 1)
	carry := ctx.Int2()
	for i := 0; i < BITS; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = pub.Xor(x[i], y[i]) // temp = xi XOR yi
		sum[i] = pub.Xor(temp[0], carry[0])

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		temp[1] = pub.And(x[i], y[i])
		temp[2] = pub.And(carry[0], temp[0])
		carry[1] = pub.Xor(temp[1], temp[2])
		carry[0] = pub.Copy(carry[1])
	}
	sum[BITS] = pub.Copy(carry[0])

	duration := time.Since(start)
	fmt.Printf("finished Bootstrapping %d bits addition circuit\n", BITS)
	fmt.Printf("total time: %s\n", duration)

	// decrypt results
	z := prv.Decrypt(sum[:])
	fmt.Println("The sum of of x and y: ", z)
}
