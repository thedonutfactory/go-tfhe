package lut

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// TestEncoderDetailed provides detailed tracing of encoder behavior
func TestEncoderDetailed(t *testing.T) {
	t.Log("=== Testing Encoder with Binary Messages ===")
	enc := NewEncoder(2)

	t.Logf("MessageModulus: %d", enc.MessageModulus)
	t.Logf("Scale: %f", enc.Scale)

	// Test encoding 0
	val0 := enc.Encode(0)
	f0 := utils.TorusToF64(val0)
	t.Logf("Encode(0) = %d (%.6f in [0,1))", val0, f0)

	// Test encoding 1
	val1 := enc.Encode(1)
	f1 := utils.TorusToF64(val1)
	t.Logf("Encode(1) = %d (%.6f in [0,1))", val1, f1)

	// Test decoding
	dec0 := enc.Decode(val0)
	dec1 := enc.Decode(val1)
	t.Logf("Decode(Encode(0)) = %d", dec0)
	t.Logf("Decode(Encode(1)) = %d", dec1)

	if dec0 != 0 {
		t.Errorf("Decode(Encode(0)) = %d, want 0", dec0)
	}
	if dec1 != 1 {
		t.Errorf("Decode(Encode(1)) = %d, want 1", dec1)
	}
}

// TestLUTGenerationDetailed provides detailed tracing of LUT generation
func TestLUTGenerationDetailed(t *testing.T) {
	t.Log("=== Testing LUT Generation for Identity Function ===")

	gen := NewGenerator(2)
	t.Logf("PolyDegree: %d", gen.PolyDegree)
	t.Logf("LookUpTableSize: %d", gen.LookUpTableSize)
	t.Logf("MessageModulus: %d", gen.Encoder.MessageModulus)
	t.Logf("Scale: %f", gen.Encoder.Scale)

	identity := func(x int) int { return x }

	t.Log("\n--- Step 1: Generate LUT ---")
	lut := gen.GenLookUpTable(identity)

	t.Log("\n--- Step 2: Examine LUT Contents ---")
	t.Log("First 20 B coefficients:")
	for i := 0; i < 20 && i < gen.PolyDegree; i++ {
		val := lut.Poly.B[i]
		fval := utils.TorusToF64(val)
		t.Logf("  B[%d] = %10d (%.6f)", i, val, fval)
	}

	t.Log("\nLast 20 B coefficients:")
	for i := gen.PolyDegree - 20; i < gen.PolyDegree; i++ {
		val := lut.Poly.B[i]
		fval := utils.TorusToF64(val)
		t.Logf("  B[%d] = %10d (%.6f)", i, val, fval)
	}

	t.Log("\n--- Step 3: Check A coefficients (should be zero) ---")
	nonZeroA := 0
	for i := 0; i < gen.PolyDegree; i++ {
		if lut.Poly.A[i] != 0 {
			nonZeroA++
		}
	}
	t.Logf("Non-zero A coefficients: %d (should be 0)", nonZeroA)

	if nonZeroA > 0 {
		t.Errorf("Expected all A coefficients to be zero, found %d non-zero", nonZeroA)
	}
}

// TestLUTGenerationStepByStep traces the algorithm step by step
func TestLUTGenerationStepByStep(t *testing.T) {
	t.Log("=== Step-by-Step LUT Generation for Identity ===")

	gen := NewGenerator(2)
	messageModulus := gen.Encoder.MessageModulus

	t.Logf("Parameters:")
	t.Logf("  MessageModulus: %d", messageModulus)
	t.Logf("  PolyDegree (N): %d", gen.PolyDegree)
	t.Logf("  LookUpTableSize (2N): %d", gen.LookUpTableSize)

	// Manually trace through the algorithm
	identity := func(x int) int { return x }

	t.Log("\n--- Step 1: Create raw LUT ---")
	lutRaw := make([]params.Torus, gen.LookUpTableSize)

	for x := 0; x < messageModulus; x++ {
		start := divRound(x*gen.LookUpTableSize, messageModulus)
		end := divRound((x+1)*gen.LookUpTableSize, messageModulus)
		y := gen.Encoder.Encode(identity(x))

		t.Logf("Message %d:", x)
		t.Logf("  f(%d) = %d", x, identity(x))
		t.Logf("  Encoded: %d (%.6f)", y, utils.TorusToF64(y))
		t.Logf("  Range in LUT: [%d, %d)", start, end)

		for i := start; i < end; i++ {
			lutRaw[i] = y
		}
	}

	t.Log("\n--- Step 2: Apply offset rotation ---")
	offset := divRound(gen.LookUpTableSize, 2*messageModulus)
	t.Logf("Offset: %d", offset)

	rotated := make([]params.Torus, gen.LookUpTableSize)
	for i := 0; i < gen.LookUpTableSize; i++ {
		srcIdx := (i + offset) % gen.LookUpTableSize
		rotated[i] = lutRaw[srcIdx]
	}

	t.Log("First 10 values after rotation:")
	for i := 0; i < 10; i++ {
		t.Logf("  rotated[%d] = %d (%.6f)", i, rotated[i], utils.TorusToF64(rotated[i]))
	}

	t.Log("\n--- Step 3: Apply negacyclic property ---")
	t.Logf("Storing first N=%d values directly", gen.PolyDegree)
	t.Logf("Subtracting second N values (due to X^N = -1)")

	result := make([]params.Torus, gen.PolyDegree)
	for i := 0; i < gen.PolyDegree; i++ {
		result[i] = rotated[i]
	}
	for i := gen.PolyDegree; i < gen.LookUpTableSize; i++ {
		result[i-gen.PolyDegree] -= rotated[i]
	}

	t.Log("\nFinal B coefficients (first 10):")
	for i := 0; i < 10; i++ {
		t.Logf("  B[%d] = %d (%.6f)", i, result[i], utils.TorusToF64(result[i]))
	}

	// Compare with actual generation
	t.Log("\n--- Comparing with actual GenLookUpTable ---")
	actualLUT := gen.GenLookUpTable(identity)

	matches := 0
	for i := 0; i < gen.PolyDegree; i++ {
		if result[i] == actualLUT.Poly.B[i] {
			matches++
		}
	}

	t.Logf("Matching coefficients: %d / %d", matches, gen.PolyDegree)

	if matches != gen.PolyDegree {
		t.Log("\nFirst 10 differences:")
		count := 0
		for i := 0; i < gen.PolyDegree && count < 10; i++ {
			if result[i] != actualLUT.Poly.B[i] {
				t.Logf("  B[%d]: manual=%d, actual=%d", i, result[i], actualLUT.Poly.B[i])
				count++
			}
		}
	}
}

// TestModSwitchDetailed traces ModSwitch behavior
func TestModSwitchDetailed(t *testing.T) {
	t.Log("=== Testing ModSwitch ===")

	gen := NewGenerator(2)
	n := gen.PolyDegree
	lookUpTableSize := gen.LookUpTableSize

	t.Logf("PolyDegree (N): %d", n)
	t.Logf("LookUpTableSize (2N): %d", lookUpTableSize)

	testCases := []struct {
		name  string
		value params.Torus
		desc  string
	}{
		{"zero", 0, "0"},
		{"quarter", params.Torus(1 << 30), "1/4 of torus"},
		{"half", params.Torus(1 << 31), "1/2 of torus"},
		{"three-quarter", params.Torus(3 << 30), "3/4 of torus"},
		{"max", params.Torus(^uint32(0)), "max value"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := gen.ModSwitch(tc.value)

			// Calculate what it should be
			fVal := utils.TorusToF64(tc.value)
			expectedFloat := fVal * float64(lookUpTableSize)

			t.Logf("Input: %s (%d)", tc.desc, tc.value)
			t.Logf("  As float in [0,1): %.6f", fVal)
			t.Logf("  Scaled to [0, 2N): %.2f", expectedFloat)
			t.Logf("  ModSwitch result: %d", result)
			t.Logf("  In range [0, %d): %v", lookUpTableSize, result >= 0 && result < lookUpTableSize)

			if result < 0 || result >= lookUpTableSize {
				t.Errorf("ModSwitch result %d out of range [0, %d)", result, lookUpTableSize)
			}
		})
	}
}

// TestCompareWithReferenceTestVector compares our LUT with what a test vector should look like
func TestCompareWithReferenceTestVector(t *testing.T) {
	t.Log("=== Comparing LUT with Reference Test Vector ===")

	// A reference test vector for binary has constant 1/8 in all positions
	// This represents the identity function in TFHE
	referenceValue := utils.F64ToTorus(0.125)

	gen := NewGenerator(2)
	identity := func(x int) int { return x }
	lut := gen.GenLookUpTable(identity)

	t.Logf("Reference value (constant 1/8): %d (%.6f)", referenceValue, utils.TorusToF64(referenceValue))

	t.Log("\nComparing first 20 B coefficients:")
	matches := 0
	for i := 0; i < 20 && i < gen.PolyDegree; i++ {
		actual := lut.Poly.B[i]
		actualF := utils.TorusToF64(actual)
		refF := utils.TorusToF64(referenceValue)

		match := ""
		if actual == referenceValue {
			matches++
			match = "✓"
		} else {
			match = "✗"
		}

		t.Logf("  B[%d]: actual=%.6f, reference=%.6f %s", i, actualF, refF, match)
	}

	t.Logf("\nMatches: %d / %d", matches, 20)
}
