package main

import (
	"fmt"
	"time"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Traditional 8-bit Addition (Bit-by-Bit Ripple Carry)         ║")
	fmt.Println("║  Using Standard Boolean Gates (NO Programmable Bootstrap)     ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Use default 128-bit security for binary operations
	params.CurrentSecurityLevel = params.Security128Bit
	fmt.Printf("Security Level: %s\n", params.SecurityInfo())
	fmt.Println()

	// Generate keys
	fmt.Println("⏱️  Generating keys...")
	keyStart := time.Now()
	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	keyDuration := time.Since(keyStart)
	fmt.Printf("   Key generation completed in %v\n", keyDuration)
	fmt.Println()

	// Test case: 42 + 137 = 179
	a := uint8(42)
	b := uint8(137)
	expected := uint8(179)

	fmt.Printf("Computing: %d + %d = %d (encrypted)\n", a, b, expected)
	fmt.Println()

	// Encrypt the two 8-bit numbers as bits
	fmt.Println("🔒 Encrypting inputs (16 bits total)...")
	encryptStart := time.Now()

	ctA := make([]*tlwe.TLWELv0, 8)
	ctB := make([]*tlwe.TLWELv0, 8)

	for i := 0; i < 8; i++ {
		bitA := (a >> i) & 1
		bitB := (b >> i) & 1

		ctA[i] = tlwe.NewTLWELv0().EncryptBool(bitA == 1, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
		ctB[i] = tlwe.NewTLWELv0().EncryptBool(bitB == 1, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
	}

	encryptDuration := time.Since(encryptStart)
	fmt.Printf("   Encryption completed in %v\n", encryptDuration)
	fmt.Println()

	// Perform 8-bit ripple-carry addition
	fmt.Println("➕ Computing 8-bit addition using ripple-carry adder...")
	fmt.Println("   (Using full adders with XOR, AND, OR gates)")
	fmt.Println()

	addStart := time.Now()

	ctSum := make([]*tlwe.TLWELv0, 8)
	var ctCarry *tlwe.TLWELv0

	// Initialize carry to 0 (false)
	ctCarry = tlwe.NewTLWELv0().EncryptBool(false, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	gateCount := 0

	// Bit-by-bit addition with full adders
	for i := 0; i < 8; i++ {
		fmt.Printf("   Processing bit %d...\n", i)

		// Full Adder:
		// sum[i] = a[i] XOR b[i] XOR carry
		// carry_out = (a[i] AND b[i]) OR (carry AND (a[i] XOR b[i]))

		// Step 1: XOR of a and b
		xorAB := gates.XOR(ctA[i], ctB[i], cloudKey)
		gateCount++

		// Step 2: Sum bit = xorAB XOR carry
		ctSum[i] = gates.XOR(xorAB, ctCarry, cloudKey)
		gateCount++

		// Step 3: Compute carry out
		// carry_out = (a AND b) OR (carry AND xorAB)
		andAB := gates.AND(ctA[i], ctB[i], cloudKey)
		gateCount++

		andCarryXor := gates.AND(ctCarry, xorAB, cloudKey)
		gateCount++

		ctCarry = gates.OR(andAB, andCarryXor, cloudKey)
		gateCount++

		fmt.Printf("      (5 gates: 2 XOR, 2 AND, 1 OR)\n")
	}

	addDuration := time.Since(addStart)

	fmt.Println()
	fmt.Printf("   ✅ Addition completed in %v\n", addDuration)
	fmt.Printf("   📊 Total gates used: %d\n", gateCount)
	fmt.Printf("   📊 Bootstraps performed: ~%d (approx %d per gate)\n", gateCount, gateCount)
	fmt.Println()

	// Decrypt and verify
	fmt.Println("🔓 Decrypting result...")
	decryptStart := time.Now()

	var result uint8
	for i := 0; i < 8; i++ {
		bit := ctSum[i].DecryptBool(secretKey.KeyLv0)
		if bit {
			result |= (1 << i)
		}
	}

	decryptDuration := time.Since(decryptStart)
	fmt.Printf("   Decryption completed in %v\n", decryptDuration)
	fmt.Println()

	// Display results
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("RESULTS:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Input A:        %d (0b%08b)\n", a, a)
	fmt.Printf("Input B:        %d (0b%08b)\n", b, b)
	fmt.Printf("Expected Sum:   %d (0b%08b)\n", expected, expected)
	fmt.Printf("Computed Sum:   %d (0b%08b)\n", result, result)
	fmt.Println()

	if result == expected {
		fmt.Println("✅ SUCCESS! Result matches expected value!")
	} else {
		fmt.Println("❌ FAILURE! Result does not match expected value!")
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("PERFORMANCE SUMMARY:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Key Generation:  %v\n", keyDuration)
	fmt.Printf("Encryption:      %v (16 bits)\n", encryptDuration)
	fmt.Printf("Addition:        %v (%d gates)\n", addDuration, gateCount)
	fmt.Printf("Decryption:      %v (8 bits)\n", decryptDuration)
	fmt.Printf("Total Time:      %v\n", keyDuration+encryptDuration+addDuration+decryptDuration)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("METHOD COMPARISON:")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Traditional (this example):\n")
	fmt.Printf("  • Operations: %d boolean gates (XOR, AND, OR)\n", gateCount)
	fmt.Printf("  • Bootstraps: ~%d (1 per gate)\n", gateCount)
	fmt.Printf("  • Time: %v\n", addDuration)
	fmt.Println()
	fmt.Printf("PBS-based (add_two_numbers_fast):\n")
	fmt.Printf("  • Operations: 4 programmable bootstraps (nibble-based)\n")
	fmt.Printf("  • Bootstraps: 4 (processes 4 bits at once)\n")
	fmt.Printf("  • Time: ~230ms (estimated with Uint5 params)\n")
	fmt.Println()
	fmt.Printf("Speedup: ~%.1fx faster with PBS! 🚀\n", float64(addDuration.Milliseconds())/230.0)
	fmt.Println()

	fmt.Println("💡 KEY INSIGHT:")
	fmt.Println("   Traditional: 40 operations processing 1 bit at a time")
	fmt.Println("   PBS Method:   4 operations processing 4 bits at once")
	fmt.Println("   Result: 10x fewer operations, significantly faster!")
}
