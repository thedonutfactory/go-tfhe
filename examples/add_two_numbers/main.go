package main

import (
	"fmt"
	"time"

	"github.com/thedonutfactory/go-tfhe/bitutils"
	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// FullAdder implements a full adder circuit
// Returns (sum, carry)
func FullAdder(serverKey *cloudkey.CloudKey, ctA, ctB, ctC *gates.Ciphertext) (*gates.Ciphertext, *gates.Ciphertext) {
	aXorB := gates.XOR(ctA, ctB, serverKey)
	aAndB := gates.AND(ctA, ctB, serverKey)
	aXorBAndC := gates.AND(aXorB, ctC, serverKey)

	// sum = (a xor b) xor c
	ctSum := gates.XOR(aXorB, ctC, serverKey)
	// carry = (a and b) or ((a xor b) and c)
	ctCarry := gates.OR(aAndB, aXorBAndC, serverKey)

	return ctSum, ctCarry
}

// Add performs homomorphic addition of two encrypted numbers
func Add(serverKey *cloudkey.CloudKey, a, b []*gates.Ciphertext, cin *gates.Ciphertext) ([]*gates.Ciphertext, *gates.Ciphertext) {
	if len(a) != len(b) {
		panic("Cannot add two numbers with different number of bits!")
	}

	result := make([]*gates.Ciphertext, len(a))
	carry := cin

	for i := 0; i < len(a); i++ {
		sum, c := FullAdder(serverKey, a[i], b[i], carry)
		carry = c
		result[i] = sum
	}

	return result, carry
}

func encrypt(x bool, secretKey *key.SecretKey) *gates.Ciphertext {
	return tlwe.NewTLWELv0().EncryptBool(x, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
}

func decrypt(x *gates.Ciphertext, secretKey *key.SecretKey) bool {
	return x.DecryptBool(secretKey.KeyLv0)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Go-TFHE: Homomorphic Addition Example           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	secretKey := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(secretKey)

	// inputs
	a := uint16(402)
	b := uint16(304)

	fmt.Printf("Input A: %d\n", a)
	fmt.Printf("Input B: %d\n", b)
	fmt.Printf("Expected Sum: %d\n", a+b)
	fmt.Println()

	aPt := bitutils.U16ToBits(a)
	bPt := bitutils.U16ToBits(b)

	// Encrypt inputs
	c1 := make([]*gates.Ciphertext, len(aPt))
	c2 := make([]*gates.Ciphertext, len(bPt))
	for i := range aPt {
		c1[i] = encrypt(aPt[i], secretKey)
		c2[i] = encrypt(bPt[i], secretKey)
	}
	cin := encrypt(false, secretKey)

	fmt.Println("Starting homomorphic addition...")
	start := time.Now()

	// ----------------- SERVER SIDE -----------------
	// Use the server public key to add the a and b ciphertexts
	c3, cout := Add(ck, c1, c2, cin)
	// -------------------------------------------------

	elapsed := time.Since(start)
	const bits uint16 = 16
	const addGatesCount uint16 = 5
	const numOps uint16 = 1
	tryNum := bits * addGatesCount * numOps
	execMsPerGate := float64(elapsed.Milliseconds()) / float64(tryNum)

	fmt.Println()
	fmt.Printf("✅ Computation complete!\n")
	fmt.Printf("⏱️  Per gate: %.2f ms\n", execMsPerGate)
	fmt.Printf("⏱️  Total: %d ms\n", elapsed.Milliseconds())
	fmt.Println()

	// Decrypt results
	r1 := make([]bool, len(c3))
	for i := range c3 {
		r1[i] = decrypt(c3[i], secretKey)
	}

	carryPt := decrypt(cout, secretKey)

	// Convert bits to integers
	s := bitutils.ConvertU16(r1)

	fmt.Println("Results:")
	fmt.Printf("  A: %d\n", a)
	fmt.Printf("  B: %d\n", b)
	fmt.Printf("  Sum: %d\n", s)
	fmt.Printf("  Carry: %v\n", carryPt)
	fmt.Println()

	if s == a+b {
		fmt.Println("✅ SUCCESS: Homomorphic addition produced correct result!")
	} else {
		fmt.Printf("❌ FAILURE: Expected %d, got %d\n", a+b, s)
	}
}
