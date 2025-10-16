package lut

import (
	"math"
	"testing"

	"github.com/thedonutfactory/go-tfhe/utils"
)

// TestAnalyzeLUTLayout analyzes the LUT layout for different functions
func TestAnalyzeLUTLayout(t *testing.T) {
	gen := NewGenerator(2)
	n := gen.PolyDegree

	t.Log("=== Analyzing LUT Layouts ===\n")

	// Analyze what positions correspond to which inputs
	t.Log("Step 1: Understanding input encoding and ModSwitch mapping")

	// For binary TFHE:
	// - Input 0 (false) encodes to -1/8 = 7/8 = 0.875
	// - Input 1 (true) encodes to 1/8 = 0.125

	falseEncoded := utils.F64ToTorus(-0.125) // = 0.875 in unsigned
	trueEncoded := utils.F64ToTorus(0.125)

	t.Logf("Encoded values:")
	t.Logf("  false: %d (%.6f)", falseEncoded, utils.TorusToF64(falseEncoded))
	t.Logf("  true:  %d (%.6f)", trueEncoded, utils.TorusToF64(trueEncoded))

	// What do these map to via ModSwitch?
	falseModSwitch := gen.ModSwitch(falseEncoded)
	trueModSwitch := gen.ModSwitch(trueEncoded)

	t.Logf("\nModSwitch values (out of [0, %d)):", 2*n)
	t.Logf("  ModSwitch(false) = %d", falseModSwitch)
	t.Logf("  ModSwitch(true)  = %d", trueModSwitch)

	t.Logf("\nAfter blind rotation by -ModSwitch, coefficient 0 comes from:")
	t.Logf("  For false input: LUT[%d %% %d] = LUT[%d]", falseModSwitch, n, falseModSwitch%n)
	t.Logf("  For true input:  LUT[%d %% %d] = LUT[%d]", trueModSwitch, n, trueModSwitch%n)

	// Analyze different functions
	functions := []struct {
		name string
		f    func(int) int
	}{
		{"identity", func(x int) int { return x }},
		{"NOT", func(x int) int { return 1 - x }},
		{"constant_0", func(x int) int { return 0 }},
		{"constant_1", func(x int) int { return 1 }},
	}

	for _, fn := range functions {
		t.Logf("\n--- Function: %s ---", fn.name)
		lut := gen.GenLookUpTable(fn.f)

		// Check what values are at the key positions
		falseLUTIdx := falseModSwitch % n
		trueLUTIdx := trueModSwitch % n

		falseLUTVal := lut.Poly.B[falseLUTIdx]
		trueLUTVal := lut.Poly.B[trueLUTIdx]

		t.Logf("LUT[%d] = %d (%.6f) - will be extracted for false input",
			falseLUTIdx, falseLUTVal, utils.TorusToF64(falseLUTVal))
		t.Logf("LUT[%d] = %d (%.6f) - will be extracted for true input",
			trueLUTIdx, trueLUTVal, utils.TorusToF64(trueLUTVal))

		// What should these be?
		expectedForFalse := fn.f(0)
		expectedForTrue := fn.f(1)

		var expectedFalseVal, expectedTrueVal float64
		if expectedForFalse == 0 {
			expectedFalseVal = 0.875 // -1/8
		} else {
			expectedFalseVal = 0.125 // 1/8
		}
		if expectedForTrue == 0 {
			expectedTrueVal = 0.875
		} else {
			expectedTrueVal = 0.125
		}

		t.Logf("\nExpected:")
		t.Logf("  %s(false) = %d → should encode to %.6f", fn.name, expectedForFalse, expectedFalseVal)
		t.Logf("  %s(true) = %d → should encode to %.6f", fn.name, expectedForTrue, expectedTrueVal)

		actualFalseVal := utils.TorusToF64(falseLUTVal)
		actualTrueVal := utils.TorusToF64(trueLUTVal)

		falseMatch := (actualFalseVal-expectedFalseVal < 0.01) || (actualFalseVal-expectedFalseVal > 0.99)
		trueMatch := (actualTrueVal-expectedTrueVal < 0.01) || (actualTrueVal-expectedTrueVal > 0.99)

		t.Logf("\nMatches:")
		t.Logf("  False input: %v (actual=%.6f, expected=%.6f)", falseMatch, actualFalseVal, expectedFalseVal)
		t.Logf("  True input:  %v (actual=%.6f, expected=%.6f)", trueMatch, actualTrueVal, expectedTrueVal)
	}
}

// TestLUTRegionMapping tests which regions of the LUT correspond to which inputs
func TestLUTRegionMapping(t *testing.T) {
	gen := NewGenerator(2)
	n := gen.PolyDegree

	t.Log("=== LUT Region Mapping Analysis ===\n")

	// Create a simple test: assign different values to different regions
	// and see what we get for different inputs

	t.Log("Creating test LUT with distinct regions:")
	testLUT := NewLookUpTable()

	// Fill first quarter with value A
	valA := utils.F64ToTorus(0.1)
	for i := 0; i < n/4; i++ {
		testLUT.Poly.B[i] = valA
		testLUT.Poly.A[i] = 0
	}

	// Fill second quarter with value B
	valB := utils.F64ToTorus(0.3)
	for i := n / 4; i < n/2; i++ {
		testLUT.Poly.B[i] = valB
		testLUT.Poly.A[i] = 0
	}

	// Fill third quarter with value C
	valC := utils.F64ToTorus(0.5)
	for i := n / 2; i < 3*n/4; i++ {
		testLUT.Poly.B[i] = valC
		testLUT.Poly.A[i] = 0
	}

	// Fill fourth quarter with value D
	valD := utils.F64ToTorus(0.7)
	for i := 3 * n / 4; i < n; i++ {
		testLUT.Poly.B[i] = valD
		testLUT.Poly.A[i] = 0
	}

	t.Logf("Region mapping:")
	t.Logf("  [0, %d): value A = %.3f", n/4, 0.1)
	t.Logf("  [%d, %d): value B = %.3f", n/4, n/2, 0.3)
	t.Logf("  [%d, %d): value C = %.3f", n/2, 3*n/4, 0.5)
	t.Logf("  [%d, %d): value D = %.3f", 3*n/4, n, 0.7)

	// Now check where false and true map to
	falseEncoded := utils.F64ToTorus(-0.125)
	trueEncoded := utils.F64ToTorus(0.125)

	falseModSwitch := gen.ModSwitch(falseEncoded)
	trueModSwitch := gen.ModSwitch(trueEncoded)

	falseLUTIdx := falseModSwitch % n
	trueLUTIdx := trueModSwitch % n

	t.Logf("\nInput mappings:")
	t.Logf("  false (0.875) → ModSwitch=%d → LUT[%d]", falseModSwitch, falseLUTIdx)
	t.Logf("  true (0.125) → ModSwitch=%d → LUT[%d]", trueModSwitch, trueLUTIdx)

	t.Logf("  false maps to region: %s", getRegion(falseLUTIdx, n))
	t.Logf("  true maps to region: %s", getRegion(trueLUTIdx, n))
}

func getRegion(idx, n int) string {
	if idx < n/4 {
		return "A (first quarter)"
	} else if idx < n/2 {
		return "B (second quarter)"
	} else if idx < 3*n/4 {
		return "C (third quarter)"
	} else {
		return "D (fourth quarter)"
	}
}

// TestCompareWithReferenceEncoding compares our encoding with reference
func TestCompareWithReferenceEncoding(t *testing.T) {
	t.Log("=== Comparing Encoding Schemes ===\n")

	gen := NewGenerator(2)
	n := gen.PolyDegree

	// Reference TFHE test vector for identity is constant 0.125
	// This means: no matter what rotation, we always get 0.125
	// But that can't give us different outputs for different inputs!
	//
	// The key insight: the INPUT ciphertext already encodes the value.
	// The test vector for GATES doesn't evaluate a function - it refreshes noise.
	//
	// For programmable bootstrap, we WANT different outputs for different inputs.

	t.Log("Key insight:")
	t.Log("  Standard bootstrap (for gates): input is PRE-PROCESSED, test vector is constant")
	t.Log("  Programmable bootstrap: test vector encodes the function")

	t.Log("\nFor NOT function:")
	t.Log("  We want: NOT(false=0) = true=1, NOT(true=1) = false=0")
	t.Log("  So LUT should have:")

	falseEncoded := utils.F64ToTorus(-0.125) // 0.875
	trueEncoded := utils.F64ToTorus(0.125)

	falseMS := gen.ModSwitch(falseEncoded)
	trueMS := gen.ModSwitch(trueEncoded)

	t.Logf("    Position %d (for false input): value for NOT(false)=true = 0.125", falseMS%n)
	t.Logf("    Position %d (for true input): value for NOT(true)=false = 0.875", trueMS%n)

	// Generate NOT LUT and check
	notFunc := func(x int) int { return 1 - x }
	notLUT := gen.GenLookUpTable(notFunc)

	actualFalsePos := notLUT.Poly.B[falseMS%n]
	actualTruePos := notLUT.Poly.B[trueMS%n]

	t.Logf("\nActual NOT LUT:")
	t.Logf("    Position %d: %.6f (expected 0.125 for true)", falseMS%n, utils.TorusToF64(actualFalsePos))
	t.Logf("    Position %d: %.6f (expected 0.875 for false)", trueMS%n, utils.TorusToF64(actualTruePos))

	// Check if they match
	falseOK := math.Abs(utils.TorusToF64(actualFalsePos)-0.125) < 0.01
	trueOK := math.Abs(utils.TorusToF64(actualTruePos)-0.875) < 0.01

	if !falseOK || !trueOK {
		t.Logf("\n⚠️  Mismatch detected!")
		t.Logf("  Position for false input: %v", falseOK)
		t.Logf("  Position for true input: %v", trueOK)
	} else {
		t.Logf("\n✓ LUT correctly encoded!")
	}
}
