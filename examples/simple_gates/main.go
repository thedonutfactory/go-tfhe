package main

import (
	"fmt"
	"strings"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func encrypt(x bool, secretKey *key.SecretKey) *gates.Ciphertext {
	return tlwe.NewTLWELv0().EncryptBool(x, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
}

func decrypt(x *gates.Ciphertext, secretKey *key.SecretKey) bool {
	return x.DecryptBool(secretKey.KeyLv0)
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              Go-TFHE: Homomorphic Gates Example              ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	secretKey := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(secretKey)

	// Test inputs
	testCases := []struct {
		a, b bool
	}{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	for _, tc := range testCases {
		fmt.Printf("Testing inputs: A=%v, B=%v\n", tc.a, tc.b)
		fmt.Println(strings.Repeat("-", 60))

		// Encrypt inputs
		ctA := encrypt(tc.a, secretKey)
		ctB := encrypt(tc.b, secretKey)

		// Test each gate
		testGate := func(name string, gateFunc func(*gates.Ciphertext, *gates.Ciphertext, *cloudkey.CloudKey) *gates.Ciphertext, expected bool) {
			result := gateFunc(ctA, ctB, ck)
			decrypted := decrypt(result, secretKey)
			status := "✅"
			if decrypted != expected {
				status = "❌"
			}
			fmt.Printf("  %s %s: %v (expected %v)\n", status, name, decrypted, expected)
		}

		testGate("AND ", gates.AND, tc.a && tc.b)
		testGate("OR  ", gates.OR, tc.a || tc.b)
		testGate("NAND", gates.NAND, !(tc.a && tc.b))
		testGate("NOR ", gates.NOR, !(tc.a || tc.b))
		testGate("XOR ", gates.XOR, tc.a != tc.b)
		testGate("XNOR", gates.XNOR, tc.a == tc.b)

		fmt.Println()
	}

	// Test NOT gate
	fmt.Println("Testing NOT gate:")
	fmt.Println(strings.Repeat("-", 60))
	for _, val := range []bool{false, true} {
		ct := encrypt(val, secretKey)
		notCt := gates.NOT(ct)
		decrypted := decrypt(notCt, secretKey)
		expected := !val
		status := "✅"
		if decrypted != expected {
			status = "❌"
		}
		fmt.Printf("  %s NOT(%v) = %v (expected %v)\n", status, val, decrypted, expected)
	}

	fmt.Println()
	fmt.Println("✅ All gate tests complete!")
}
