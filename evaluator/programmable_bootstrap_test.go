package evaluator

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/lut"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

func TestProgrammableBootstrapIdentity(t *testing.T) {
	// Use 80-bit security for faster testing
	oldSecurityLevel := params.CurrentSecurityLevel
	params.CurrentSecurityLevel = params.Security80Bit
	defer func() { params.CurrentSecurityLevel = oldSecurityLevel }()

	// Generate keys
	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)

	// Create evaluator
	eval := NewEvaluator(params.GetTRGSWLv1().N)

	// Test identity function: f(x) = x
	identity := func(x int) int { return x }

	// Test with both 0 and 1
	testCases := []struct {
		name  string
		input int
		want  int
	}{
		{"identity(0)", 0, 0},
		{"identity(1)", 1, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt input using LWE message encoding (not binary encoding!)
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(tc.input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			// Apply programmable bootstrap
			result := eval.BootstrapFunc(
				ct,
				identity,
				2, // binary message modulus
				cloudKey.BootstrappingKey,
				cloudKey.KeySwitchingKey,
				cloudKey.DecompositionOffset,
			)

			// Decrypt and verify using LWE message decoding
			decrypted := result.DecryptLWEMessage(2, secretKey.KeyLv0)
			if decrypted != tc.want {
				t.Errorf("identity(%d) = %d, want %d", tc.input, decrypted, tc.want)
			}
		})
	}
}

func TestProgrammableBootstrapNOT(t *testing.T) {
	oldSecurityLevel := params.CurrentSecurityLevel
	params.CurrentSecurityLevel = params.Security80Bit
	defer func() { params.CurrentSecurityLevel = oldSecurityLevel }()

	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := NewEvaluator(params.GetTRGSWLv1().N)

	// Test NOT function: f(x) = 1 - x
	notFunc := func(x int) int { return 1 - x }

	testCases := []struct {
		name  string
		input int
		want  int
	}{
		{"NOT(0)", 0, 1},
		{"NOT(1)", 1, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(tc.input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			result := eval.BootstrapFunc(
				ct,
				notFunc,
				2,
				cloudKey.BootstrappingKey,
				cloudKey.KeySwitchingKey,
				cloudKey.DecompositionOffset,
			)

			decrypted := result.DecryptLWEMessage(2, secretKey.KeyLv0)
			if decrypted != tc.want {
				t.Errorf("NOT(%d) = %d, want %d", tc.input, decrypted, tc.want)
			}
		})
	}
}

func TestProgrammableBootstrapConstant(t *testing.T) {
	oldSecurityLevel := params.CurrentSecurityLevel
	params.CurrentSecurityLevel = params.Security80Bit
	defer func() { params.CurrentSecurityLevel = oldSecurityLevel }()

	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := NewEvaluator(params.GetTRGSWLv1().N)

	// Test constant function: f(x) = 1 (always returns 1)
	constantOne := func(x int) int { return 1 }

	testCases := []struct {
		name  string
		input int
	}{
		{"constant(0)", 0},
		{"constant(1)", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ct := tlwe.NewTLWELv0()
			ct.EncryptLWEMessage(tc.input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

			result := eval.BootstrapFunc(
				ct,
				constantOne,
				2,
				cloudKey.BootstrappingKey,
				cloudKey.KeySwitchingKey,
				cloudKey.DecompositionOffset,
			)

			// Should always decrypt to 1
			decrypted := result.DecryptLWEMessage(2, secretKey.KeyLv0)
			if decrypted != 1 {
				t.Errorf("constant(%d) = %d, want 1", tc.input, decrypted)
			}
		})
	}
}

func TestBootstrapLUTReuse(t *testing.T) {
	// Test that we can reuse a lookup table for multiple encryptions
	oldSecurityLevel := params.CurrentSecurityLevel
	params.CurrentSecurityLevel = params.Security80Bit
	defer func() { params.CurrentSecurityLevel = oldSecurityLevel }()

	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := NewEvaluator(params.GetTRGSWLv1().N)
	gen := lut.NewGenerator(2)

	// Pre-compute lookup table for NOT function
	notFunc := func(x int) int { return 1 - x }
	lookupTable := gen.GenLookUpTable(notFunc)

	// Apply to multiple inputs using the same LUT
	inputs := []int{0, 1, 0, 1, 0}

	for i, input := range inputs {
		ct := tlwe.NewTLWELv0()
		ct.EncryptLWEMessage(input, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

		// Use pre-computed LUT
		result := eval.BootstrapLUT(
			ct,
			lookupTable,
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)

		decrypted := result.DecryptLWEMessage(2, secretKey.KeyLv0)
		expected := 1 - input

		if decrypted != expected {
			t.Errorf("test %d: NOT(%d) = %d, want %d", i, input, decrypted, expected)
		}
	}
}

func TestModSwitch(t *testing.T) {
	gen := lut.NewGenerator(2)
	n := params.GetTRGSWLv1().N

	// Test that ModSwitch returns values in valid range
	tests := []params.Torus{
		0,
		1 << 30,
		1 << 31,
		3 << 30,
		params.Torus(^uint32(0)),
	}

	for _, val := range tests {
		result := gen.ModSwitch(val)
		if result < 0 || result >= n {
			t.Errorf("ModSwitch(%d) = %d, out of range [0, %d)", val, result, n)
		}
	}
}

// Benchmark programmable bootstrapping performance
func BenchmarkProgrammableBootstrap(b *testing.B) {
	params.CurrentSecurityLevel = params.Security80Bit

	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := NewEvaluator(params.GetTRGSWLv1().N)

	// Create input ciphertext
	ct := tlwe.NewTLWELv0()
	ct.EncryptLWEMessage(1, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	// Identity function
	identity := func(x int) int { return x }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eval.BootstrapFunc(
			ct,
			identity,
			2,
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)
	}
}

// Benchmark LUT reuse
func BenchmarkBootstrapLUT(b *testing.B) {
	params.CurrentSecurityLevel = params.Security80Bit

	secretKey := key.NewSecretKey()
	cloudKey := cloudkey.NewCloudKey(secretKey)
	eval := NewEvaluator(params.GetTRGSWLv1().N)
	gen := lut.NewGenerator(2)

	// Pre-compute LUT
	identity := func(x int) int { return x }
	lookupTable := gen.GenLookUpTable(identity)

	// Create input ciphertext
	ct := tlwe.NewTLWELv0()
	ct.EncryptLWEMessage(1, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eval.BootstrapLUT(
			ct,
			lookupTable,
			cloudKey.BootstrappingKey,
			cloudKey.KeySwitchingKey,
			cloudKey.DecompositionOffset,
		)
	}
}
