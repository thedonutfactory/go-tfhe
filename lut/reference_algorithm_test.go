package lut

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// TestReferenceAlgorithmStepByStep traces the reference algorithm step by step
func TestReferenceAlgorithmStepByStep(t *testing.T) {
	messageModulus := 2
	polyDegree := params.GetTRGSWLv1().N // 1024
	lookUpTableSize := 2 * polyDegree    // 2048

	t.Log("=== Reference Algorithm for NOT Function ===\n")
	t.Logf("Parameters: messageModulus=%d, N=%d, LUTSize=%d\n", messageModulus, polyDegree, lookUpTableSize)

	notFunc := func(x int) int { return 1 - x }

	// Step 1: Create raw LUT
	t.Log("Step 1: Fill raw LUT")
	lutRaw := make([]params.Torus, lookUpTableSize)

	for x := 0; x < messageModulus; x++ {
		start := divRound(x*lookUpTableSize, messageModulus)
		end := divRound((x+1)*lookUpTableSize, messageModulus)

		output := notFunc(x)
		var encodedOutput params.Torus
		if output == 0 {
			encodedOutput = utils.F64ToTorus(-0.125) // 0.875
		} else {
			encodedOutput = utils.F64ToTorus(0.125)
		}

		t.Logf("  Message %d: NOT(%d)=%d → encode to %.3f", x, x, output, utils.TorusToF64(encodedOutput))
		t.Logf("    Fill indices [%d, %d)", start, end)

		for i := start; i < end; i++ {
			lutRaw[i] = encodedOutput
		}
	}

	t.Log("\n  Check key positions in raw LUT:")
	checkPos := []int{0, 256, 512, 768, 1024, 1280, 1536, 1792}
	for _, pos := range checkPos {
		t.Logf("    lutRaw[%4d] = %.3f", pos, utils.TorusToF64(lutRaw[pos]))
	}

	// Step 2: Rotate by offset
	offset := divRound(lookUpTableSize, 2*messageModulus)
	t.Logf("\nStep 2: Rotate by offset=%d", offset)

	rotated := make([]params.Torus, lookUpTableSize)
	for i := 0; i < lookUpTableSize; i++ {
		srcIdx := (i + offset) % lookUpTableSize
		rotated[i] = lutRaw[srcIdx]
	}

	t.Log("  Check key positions after rotation:")
	for _, pos := range checkPos {
		srcPos := (pos + offset) % lookUpTableSize
		t.Logf("    rotated[%4d] = lutRaw[%4d] = %.3f", pos, srcPos, utils.TorusToF64(rotated[pos]))
	}

	// Step 3: Negate tail
	negateStart := lookUpTableSize - offset
	t.Logf("\nStep 3: Negate indices [%d, %d)", negateStart, lookUpTableSize)

	for i := negateStart; i < lookUpTableSize; i++ {
		rotated[i] = -rotated[i]
	}

	t.Log("  Check key positions after negation:")
	for _, pos := range checkPos {
		neg := ""
		if pos >= negateStart {
			neg = " (negated)"
		}
		t.Logf("    rotated[%4d] = %.3f%s", pos, utils.TorusToF64(rotated[pos]), neg)
	}

	// Step 4: Store first N coefficients
	t.Logf("\nStep 4: Store first N=%d coefficients in polynomial", polyDegree)

	result := NewLookUpTable()
	for i := 0; i < polyDegree; i++ {
		result.Poly.B[i] = rotated[i]
		result.Poly.A[i] = 0
	}

	t.Log("\n  Final LUT key positions:")
	checkPosFinal := []int{0, 256, 512, 768}
	for _, pos := range checkPosFinal {
		t.Logf("    LUT.Poly.B[%4d] = %.3f", pos, utils.TorusToF64(result.Poly.B[pos]))
	}

	// Compare with actual generator
	t.Log("\n  Comparing with GenLookUpTable:")
	gen := NewGenerator(2)
	actualLUT := gen.GenLookUpTable(notFunc)

	matches := 0
	for i := 0; i < polyDegree; i++ {
		if result.Poly.B[i] == actualLUT.Poly.B[i] {
			matches++
		}
	}
	t.Logf("    Matching coefficients: %d / %d", matches, polyDegree)

	if matches != polyDegree {
		t.Log("\n  First 10 mismatches:")
		count := 0
		for i := 0; i < polyDegree && count < 10; i++ {
			if result.Poly.B[i] != actualLUT.Poly.B[i] {
				t.Logf("      [%d]: manual=%.3f, actual=%.3f",
					i, utils.TorusToF64(result.Poly.B[i]), utils.TorusToF64(actualLUT.Poly.B[i]))
				count++
			}
		}
	}

	// Now verify this gives correct results for ideal inputs
	t.Log("\n  Verification with ideal encoded inputs:")

	falseIdeal := utils.F64ToTorus(-0.125) // 0.875
	trueIdeal := utils.F64ToTorus(0.125)

	falseMS := gen.ModSwitch(falseIdeal)
	trueMS := gen.ModSwitch(trueIdeal)

	t.Logf("    false (0.875) → ModSwitch=%d → extract from LUT[%d]", falseMS, falseMS%polyDegree)
	t.Logf("      Value: %.3f, Expected: %.3f (NOT(false)=true)",
		utils.TorusToF64(result.Poly.B[falseMS%polyDegree]), 0.125)

	t.Logf("    true (0.125) → ModSwitch=%d → extract from LUT[%d]", trueMS, trueMS%polyDegree)
	t.Logf("      Value: %.3f, Expected: %.3f (NOT(true)=false)",
		utils.TorusToF64(result.Poly.B[trueMS%polyDegree]), 0.875)
}
