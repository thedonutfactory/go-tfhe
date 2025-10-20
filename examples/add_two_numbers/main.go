package main

import (
	"fmt"
	"time"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/evaluator"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/lut"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Fast 8-bit Addition Using Programmable Bootstrapping         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Use Uint5 parameters for messageModulus=32
	params.CurrentSecurityLevel = params.SecurityUint5
	fmt.Printf("Security Level: %s\n", params.SecurityInfo())
	fmt.Println()

	// Generate keys
	fmt.Println("⏱️  Generating keys...")
	keyStart := time.Now()
	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := evaluator.NewEvaluator(params.GetTRGSWLv1().N)
	keyDuration := time.Since(keyStart)
	fmt.Printf("   Key generation completed in %v\n", keyDuration)
	fmt.Println()

	// Inputs
	a := uint8(42)
	b := uint8(137)
	expected := uint8(179)

	fmt.Printf("Computing: %d + %d = %d (encrypted)\n", a, b, expected)
	fmt.Println()

	// Step 1: Split into nibbles (4-bit chunks)
	aLow := int(a & 0x0F)         // Low nibble of a (bits 0-3)
	aHigh := int((a >> 4) & 0x0F) // High nibble of a (bits 4-7)
	bLow := int(b & 0x0F)         // Low nibble of b
	bHigh := int((b >> 4) & 0x0F) // High nibble of b

	fmt.Printf("Input A: %3d = 0b%04b_%04b (nibbles: high=%d, low=%d)\n", a, aHigh, aLow, aHigh, aLow)
	fmt.Printf("Input B: %3d = 0b%04b_%04b (nibbles: high=%d, low=%d)\n", b, bHigh, bLow, bHigh, bLow)
	fmt.Println()

	// Step 2: Generate lookup tables
	fmt.Println("📋 Generating lookup tables...")
	lutStart := time.Now()
	gen := lut.NewGenerator(32)

	lutSumLow := gen.GenLookUpTable(func(x int) int {
		return x % 16 // Extract lower 4 bits
	})

	lutCarryLow := gen.GenLookUpTable(func(x int) int {
		if x >= 16 {
			return 1 // Carry out
		}
		return 0
	})

	lutSumHigh := gen.GenLookUpTable(func(x int) int {
		return x % 16 // Extract lower 4 bits
	})

	lutDuration := time.Since(lutStart)
	fmt.Printf("   LUT generation: %v\n", lutDuration)
	fmt.Println()

	// Step 3: Encrypt nibbles
	fmt.Println("🔒 Encrypting nibbles...")
	encStart := time.Now()

	ctALow := tlwe.NewTLWELv0()
	ctALow.EncryptLWEMessage(aLow, 32, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	ctAHigh := tlwe.NewTLWELv0()
	ctAHigh.EncryptLWEMessage(aHigh, 32, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	ctBLow := tlwe.NewTLWELv0()
	ctBLow.EncryptLWEMessage(bLow, 32, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	ctBHigh := tlwe.NewTLWELv0()
	ctBHigh.EncryptLWEMessage(bHigh, 32, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	encDuration := time.Since(encStart)
	fmt.Printf("   Encrypted 4 nibbles in %v\n", encDuration)
	fmt.Println()

	// Step 4: Homomorphic addition of low nibbles (no bootstrap needed!)
	fmt.Println("➕ Computing encrypted addition...")
	addStart := time.Now()

	n := params.GetTLWELv0().N
	ctTempLow := tlwe.NewTLWELv0()
	for j := 0; j < n+1; j++ {
		ctTempLow.P[j] = ctALow.P[j] + ctBLow.P[j]
	}
	fmt.Println("   Step 1: Low nibbles added (homomorphic add, no bootstrap)")

	// Step 5: Bootstrap 1 - Extract low sum (mod 16)
	pbs1Start := time.Now()
	ctSumLow := eval.BootstrapLUT(ctTempLow, lutSumLow,
		cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)
	pbs1Duration := time.Since(pbs1Start)
	fmt.Printf("   Bootstrap 1: Extract low sum (mod 16) - %v\n", pbs1Duration)

	// Step 6: Bootstrap 2 - Extract carry from low nibbles
	pbs2Start := time.Now()
	ctCarry := eval.BootstrapLUT(ctTempLow, lutCarryLow,
		cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)
	pbs2Duration := time.Since(pbs2Start)
	fmt.Printf("   Bootstrap 2: Extract carry bit - %v\n", pbs2Duration)

	// Step 7: Add high nibbles + carry (homomorphic)
	ctTempHigh := tlwe.NewTLWELv0()
	for j := 0; j < n+1; j++ {
		ctTempHigh.P[j] = ctAHigh.P[j] + ctBHigh.P[j] + ctCarry.P[j]
	}
	fmt.Println("   Step 2: High nibbles + carry added (homomorphic add, no bootstrap)")

	// Step 8: Bootstrap 3 - Extract high sum (mod 16)
	pbs3Start := time.Now()
	ctSumHigh := eval.BootstrapLUT(ctTempHigh, lutSumHigh,
		cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)
	pbs3Duration := time.Since(pbs3Start)
	fmt.Printf("   Bootstrap 3: Extract high sum (mod 16) - %v\n", pbs3Duration)

	addDuration := time.Since(addStart)
	fmt.Println()

	// Step 9: Decrypt results
	fmt.Println("🔓 Decrypting result...")
	decStart := time.Now()

	sumLow := ctSumLow.DecryptLWEMessage(32, secretKey.KeyLv0)
	sumHigh := ctSumHigh.DecryptLWEMessage(32, secretKey.KeyLv0)

	decDuration := time.Since(decStart)
	fmt.Printf("   Decrypted nibbles in %v\n", decDuration)
	fmt.Println()

	// Step 10: Combine nibbles into final result
	result := uint8(sumLow | (sumHigh << 4))

	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("RESULTS")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Input A:    %3d = 0b%04b_%04b\n", a, aHigh, aLow)
	fmt.Printf("Input B:    %3d = 0b%04b_%04b\n", b, bHigh, bLow)
	fmt.Printf("Result:     %3d = 0b%04b_%04b (nibbles: high=%d, low=%d)\n",
		result, sumHigh, sumLow, sumHigh, sumLow)
	fmt.Printf("Expected:   %3d\n", expected)
	fmt.Println()

	if result == expected {
		fmt.Println("✅ SUCCESS! Result is correct!")
	} else {
		fmt.Printf("❌ FAILURE! Expected %d, got %d\n", expected, result)
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("PERFORMANCE SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("Key Generation:  %v\n", keyDuration)
	fmt.Printf("LUT Generation:  %v\n", lutDuration)
	fmt.Printf("Encryption:      %v (4 nibbles)\n", encDuration)
	fmt.Printf("Addition:        %v (3 bootstraps)\n", addDuration)
	fmt.Printf("  - Bootstrap 1: %v (low sum)\n", pbs1Duration)
	fmt.Printf("  - Bootstrap 2: %v (carry)\n", pbs2Duration)
	fmt.Printf("  - Bootstrap 3: %v (high sum)\n", pbs3Duration)
	fmt.Printf("Decryption:      %v\n", decDuration)
	fmt.Println()

}
