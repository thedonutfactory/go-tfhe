package params_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/evaluator"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/lut"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// TestAllUintParameters tests all Uint parameter sets with programmable bootstrapping
func TestAllUintParameters(t *testing.T) {
	testCases := []struct {
		name           string
		secLevel       params.SecurityLevel
		messageModulus int
		skipReason     string // If non-empty, test will be skipped
	}{
		{"Uint1", params.SecurityUint1, 2, ""},
		{"Uint2", params.SecurityUint2, 4, ""},
		{"Uint3", params.SecurityUint3, 8, ""},
		{"Uint4", params.SecurityUint4, 16, ""},
		{"Uint5", params.SecurityUint5, 32, ""},
		{"Uint6", params.SecurityUint6, 64, "Extended LUT (polyExtendFactor=2) not fully implemented"},
		{"Uint7", params.SecurityUint7, 128, "Extended LUT (polyExtendFactor=4) not fully implemented"},
		{"Uint8", params.SecurityUint8, 256, "Extended LUT (polyExtendFactor=9) not fully implemented"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skipReason != "" {
				t.Skipf("Skipping %s: %s", tc.name, tc.skipReason)
				return
			}
			testUintParameterSet(t, tc.secLevel, tc.name, tc.messageModulus)
		})
	}
}

func testUintParameterSet(t *testing.T, secLevel params.SecurityLevel, name string, messageModulus int) {
	params.CurrentSecurityLevel = secLevel

	t.Logf("Testing %s with messageModulus=%d, N=%d", name, messageModulus, params.GetTRGSWLv1().N)

	// Generate keys
	keyStart := time.Now()
	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := evaluator.NewEvaluator(params.GetTRGSWLv1().N)
	keyDuration := time.Since(keyStart)
	t.Logf("Key generation: %v", keyDuration)

	gen := lut.NewGenerator(messageModulus)

	// Test identity function on a subset of values
	t.Run("Identity", func(t *testing.T) {
		lutId := gen.GenLookUpTable(func(x int) int { return x })

		// Test first few and last few values
		testValues := getTestValues(messageModulus)

		for _, x := range testValues {
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(x, messageModulus, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			ctResult := eval.BootstrapLUT(ct, lutId, cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)

			result := ctResult.DecryptLWEMessage(messageModulus, secretKey.KeyLv0)

			if result != x {
				t.Errorf("identity(%d) = %d, want %d", x, result, x)
			}
		}
	})

	// Test NOT-like function (complement)
	t.Run("Complement", func(t *testing.T) {
		lutComplement := gen.GenLookUpTable(func(x int) int {
			return (messageModulus - 1) - x
		})

		testValues := getTestValues(messageModulus)

		for _, x := range testValues {
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(x, messageModulus, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			ctResult := eval.BootstrapLUT(ct, lutComplement, cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)

			result := ctResult.DecryptLWEMessage(messageModulus, secretKey.KeyLv0)
			expected := (messageModulus - 1) - x

			if result != expected {
				t.Errorf("complement(%d) = %d, want %d", x, result, expected)
			}
		}
	})

	// Test modulo function
	t.Run("Modulo", func(t *testing.T) {
		modValue := messageModulus / 2
		lutMod := gen.GenLookUpTable(func(x int) int {
			return x % modValue
		})

		testValues := getTestValues(messageModulus)

		for _, x := range testValues {
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(x, messageModulus, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			ctResult := eval.BootstrapLUT(ct, lutMod, cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)

			result := ctResult.DecryptLWEMessage(messageModulus, secretKey.KeyLv0)
			expected := x % modValue

			if result != expected {
				t.Errorf("(%d %% %d) = %d, want %d", x, modValue, result, expected)
			}
		}
	})
}

// getTestValues returns a subset of test values to keep tests fast
// Tests first 3, middle value, and last 3 values
func getTestValues(max int) []int {
	if max <= 8 {
		// Small modulus: test all values
		result := make([]int, max)
		for i := 0; i < max; i++ {
			result[i] = i
		}
		return result
	}

	// Large modulus: test subset
	return []int{
		0, 1, 2, // First few
		max / 2,                   // Middle
		max - 3, max - 2, max - 1, // Last few
	}
}

// BenchmarkUintParameters benchmarks key generation for all Uint parameter sets
func BenchmarkUintParameters(b *testing.B) {
	paramSets := []struct {
		name     string
		secLevel params.SecurityLevel
	}{
		{"Uint1", params.SecurityUint1},
		{"Uint2", params.SecurityUint2},
		{"Uint3", params.SecurityUint3},
		{"Uint4", params.SecurityUint4},
		{"Uint5", params.SecurityUint5},
		{"Uint6", params.SecurityUint6},
		{"Uint7", params.SecurityUint7},
		{"Uint8", params.SecurityUint8},
	}

	for _, ps := range paramSets {
		b.Run(fmt.Sprintf("KeyGen/%s", ps.name), func(b *testing.B) {
			params.CurrentSecurityLevel = ps.secLevel
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				secretKey := key.NewSecretKey()
				_ = cloudkey.NewCloudKey(secretKey)
			}
		})
	}
}

// BenchmarkPBS benchmarks programmable bootstrapping for each Uint parameter set
func BenchmarkPBS(b *testing.B) {
	paramSets := []struct {
		name           string
		secLevel       params.SecurityLevel
		messageModulus int
	}{
		{"Uint1", params.SecurityUint1, 2},
		{"Uint2", params.SecurityUint2, 4},
		{"Uint3", params.SecurityUint3, 8},
		{"Uint4", params.SecurityUint4, 16},
		{"Uint5", params.SecurityUint5, 32},
		{"Uint6", params.SecurityUint6, 64},
		{"Uint7", params.SecurityUint7, 128},
		{"Uint8", params.SecurityUint8, 256},
	}

	for _, ps := range paramSets {
		b.Run(ps.name, func(b *testing.B) {
			params.CurrentSecurityLevel = ps.secLevel

			secretKey := key.NewSecretKey()
			cloudKey := cloudkey.NewCloudKey(secretKey)
			eval := evaluator.NewEvaluator(params.GetTRGSWLv1().N)

			gen := lut.NewGenerator(ps.messageModulus)
			lutId := gen.GenLookUpTable(func(x int) int { return x })

			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(1, ps.messageModulus, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = eval.BootstrapLUT(ct, lutId, cloudKey.BootstrappingKey, cloudKey.KeySwitchingKey, cloudKey.DecompositionOffset)
			}
		})
	}
}

// TestUintParameterProperties verifies parameter properties
func TestUintParameterProperties(t *testing.T) {
	testCases := []struct {
		name           string
		secLevel       params.SecurityLevel
		expectedN      int
		expectedLweN   int
		messageModulus int
	}{
		{"Uint1", params.SecurityUint1, 1024, 700, 2},
		{"Uint2", params.SecurityUint2, 512, 687, 4},
		{"Uint3", params.SecurityUint3, 1024, 820, 8},
		{"Uint4", params.SecurityUint4, 2048, 820, 16},
		{"Uint5", params.SecurityUint5, 2048, 1071, 32},
		{"Uint6", params.SecurityUint6, 2048, 1071, 64},
		{"Uint7", params.SecurityUint7, 2048, 1160, 128},
		{"Uint8", params.SecurityUint8, 2048, 1160, 256},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params.CurrentSecurityLevel = tc.secLevel

			n := params.GetTRGSWLv1().N
			lweN := params.GetTLWELv0().N

			if n != tc.expectedN {
				t.Errorf("Polynomial degree: got %d, want %d", n, tc.expectedN)
			}

			if lweN != tc.expectedLweN {
				t.Errorf("LWE dimension: got %d, want %d", lweN, tc.expectedLweN)
			}

			// Verify other parameters are set
			if params.GetTLWELv0().ALPHA == 0 {
				t.Error("LWE noise not set")
			}

			if params.GetTRGSWLv1().BG == 0 {
				t.Error("TRGSW base not set")
			}

			t.Logf("%s: N=%d, LWE_N=%d, messageModulus=%d", tc.name, n, lweN, tc.messageModulus)
		})
	}
}
