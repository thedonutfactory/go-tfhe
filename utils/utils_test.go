package utils_test

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

func TestF64ToTorus(t *testing.T) {
	testCases := []struct {
		input    float64
		expected params.Torus
	}{
		{0.0, 0},
		{0.125, 536870912},   // 2^29
		{-0.125, 3758096384}, // 2^32 - 2^29
		{0.25, 1073741824},   // 2^30
		{0.5, 2147483648},    // 2^31
	}

	for _, tc := range testCases {
		result := utils.F64ToTorus(tc.input)
		if result != tc.expected {
			t.Errorf("F64ToTorus(%f) = %d (0x%08x), expected %d (0x%08x)",
				tc.input, result, result, tc.expected, tc.expected)
		}
	}
}

func TestF64ToTorusVec(t *testing.T) {
	input := []float64{0.0, 0.125, 0.25}
	expected := []params.Torus{0, 536870912, 1073741824}

	result := utils.F64ToTorusVec(input)

	if len(result) != len(expected) {
		t.Fatalf("F64ToTorusVec length = %d, expected %d", len(result), len(expected))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("F64ToTorusVec[%d] = %d, expected %d", i, result[i], expected[i])
		}
	}
}
