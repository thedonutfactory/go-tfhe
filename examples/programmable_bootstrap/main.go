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
	fmt.Println("=== Programmable Bootstrapping Demo ===")
	fmt.Println()

	// Use 80-bit security for faster demo
	params.CurrentSecurityLevel = params.Security80Bit
	fmt.Printf("Security Level: %s\n", params.SecurityInfo())
	fmt.Println()

	// Generate keys
	fmt.Println("Generating keys...")
	startKey := time.Now()
	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	fmt.Printf("Key generation took: %v\n", time.Since(startKey))
	fmt.Println()

	// Create evaluator
	eval := evaluator.NewEvaluator(params.GetTRGSWLv1().N)

	// Example 1: Identity function
	fmt.Println("Example 1: Identity Function (f(x) = x)")
	fmt.Println("This refreshes noise while preserving the value")
	identity := func(x int) int { return x }
	demoFunction(eval, secretKey, cloudKey, identity, "identity", 0, 1)

	// Example 2: NOT function
	fmt.Println("\nExample 2: NOT Function (f(x) = 1 - x)")
	fmt.Println("This flips the bit during bootstrapping")
	notFunc := func(x int) int { return 1 - x }
	demoFunction(eval, secretKey, cloudKey, notFunc, "NOT", 0, 1)

	// Example 3: Constant function
	fmt.Println("\nExample 3: Constant Function (f(x) = 1)")
	fmt.Println("This always returns 1, regardless of input")
	constantOne := func(x int) int { return 1 }
	demoFunction(eval, secretKey, cloudKey, constantOne, "constant(1)", 0, 1)

	// Example 4: AND with constant (simulation)
	fmt.Println("\nExample 4: Constant Function (f(x) = 0)")
	fmt.Println("This always returns 0")
	constantZero := func(x int) int { return 0 }
	demoFunction(eval, secretKey, cloudKey, constantZero, "constant(0)", 0, 1)

	// Example 5: LUT reuse demonstration
	fmt.Println("\nExample 5: Lookup Table Reuse")
	fmt.Println("Pre-compute LUT once, use multiple times for efficiency")
	demoLUTReuse(eval, secretKey, cloudKey)

	// Example 6: Multi-bit messages (4 values)
	fmt.Println("\nExample 6: Multi-bit Messages (2-bit values)")
	demoMultiBit(eval, secretKey, cloudKey)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nNote: Programmable bootstrapping uses general LWE message encoding")
	fmt.Println("(message * scale), not binary boolean encoding (±1/8).")
	fmt.Println("Use EncryptLWEMessage() for encryption and DecryptLWEMessage() for decryption.")
}

func demoFunction(eval *evaluator.Evaluator, secretKey *key.SecretKey, cloudKey *cloudkey.CloudKey,
	f func(int) int, name string, inputs ...int) {

	for i, input := range inputs {
		// Encrypt input using LWE message encoding
		ct := tlwe.NewTLWELv0()
		ct.EncryptLWEMessage(input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

		// Apply programmable bootstrap
		start := time.Now()
		result := eval.BootstrapFunc(
			ct,
			f,
			2, // binary (message modulus = 2)
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)
		elapsed := time.Since(start)

		// Decrypt using LWE message decoding
		output := result.DecryptLWEMessage(2, secretKey.KeyLv0)

		fmt.Printf("  Input %d: %d → %s(%d) = %d (took %v)\n",
			i+1, input, name, input, output, elapsed)
	}
}

func demoLUTReuse(eval *evaluator.Evaluator, secretKey *key.SecretKey, cloudKey *cloudkey.CloudKey) {
	// Pre-compute lookup table for NOT function
	gen := lut.NewGenerator(2)
	notFunc := func(x int) int { return 1 - x }

	fmt.Println("  Pre-computing NOT lookup table...")
	start := time.Now()
	lookupTable := gen.GenLookUpTable(notFunc)
	lutTime := time.Since(start)
	fmt.Printf("  LUT generation took: %v\n", lutTime)

	// Apply to multiple inputs using the same LUT
	inputs := []int{0, 1, 0, 1, 0}

	var totalBootstrapTime time.Duration
	for i, input := range inputs {
		ct := tlwe.NewTLWELv0()
		ct.EncryptLWEMessage(input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

		start := time.Now()
		result := eval.BootstrapLUT(
			ct,
			lookupTable,
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)
		elapsed := time.Since(start)
		totalBootstrapTime += elapsed

		output := result.DecryptLWEMessage(2, secretKey.KeyLv0)
		fmt.Printf("  Input %d: %d → NOT(%d) = %d (took %v)\n",
			i+1, input, input, output, elapsed)
	}

	avgTime := totalBootstrapTime / time.Duration(len(inputs))
	fmt.Printf("  Average bootstrap time: %v\n", avgTime)
	fmt.Println("  ✓ LUT reuse avoids recomputing the lookup table!")
}

func demoMultiBit(eval *evaluator.Evaluator, secretKey *key.SecretKey, cloudKey *cloudkey.CloudKey) {
	// Use 2-bit messages (values 0, 1, 2, 3)
	messageModulus := 4

	// Function that increments by 1 (mod 4)
	increment := func(x int) int { return (x + 1) % 4 }

	fmt.Println("  Testing increment function: f(x) = (x + 1) mod 4")

	// Test a few values
	testInputs := []int{0, 1, 2, 3}

	for _, input := range testInputs {
		ct := tlwe.NewTLWELv0()
		ct.EncryptLWEMessage(input, messageModulus, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

		start := time.Now()
		result := eval.BootstrapFunc(
			ct,
			increment,
			messageModulus,
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)
		elapsed := time.Since(start)

		output := result.DecryptLWEMessage(messageModulus, secretKey.KeyLv0)
		expected := increment(input)

		status := "✓"
		if output != expected {
			status = "✗"
		}

		fmt.Printf("  increment(%d) = %d (expected %d) %s (took %v)\n",
			input, output, expected, status, elapsed)
	}

	fmt.Println("  ✓ Framework supports arbitrary message moduli!")
}
