package bitutils_test

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/bitutils"
)

func TestU8ToBitsAndBack(t *testing.T) {
	testCases := []uint8{0, 1, 5, 42, 127, 255}

	for _, val := range testCases {
		bits := bitutils.U8ToBits(val)
		result := bitutils.ConvertU8(bits)

		if result != val {
			t.Errorf("U8: %d -> bits -> %d", val, result)
		}
	}
}

func TestU16ToBitsAndBack(t *testing.T) {
	testCases := []uint16{0, 1, 5, 42, 127, 255, 402, 304, 706, 65535}

	for _, val := range testCases {
		bits := bitutils.U16ToBits(val)
		result := bitutils.ConvertU16(bits)

		if result != val {
			t.Errorf("U16: %d -> bits -> %d", val, result)
		}
	}
}

func TestU32ToBitsAndBack(t *testing.T) {
	testCases := []uint32{0, 1, 42, 1000000, 4294967295}

	for _, val := range testCases {
		bits := bitutils.U32ToBits(val)
		result := bitutils.ConvertU32(bits)

		if result != val {
			t.Errorf("U32: %d -> bits -> %d", val, result)
		}
	}
}

func TestU64ToBitsAndBack(t *testing.T) {
	testCases := []uint64{0, 1, 42, 1000000, 18446744073709551615}

	for _, val := range testCases {
		bits := bitutils.U64ToBits(val)
		result := bitutils.ConvertU64(bits)

		if result != val {
			t.Errorf("U64: %d -> bits -> %d", val, result)
		}
	}
}

func TestToBitsLSBFirst(t *testing.T) {
	// Verify that bits are in LSB-first order
	bits := bitutils.U8ToBits(5) // 5 = 0b00000101

	// LSB first: [true, false, true, false, false, false, false, false]
	expected := []bool{true, false, true, false, false, false, false, false}

	if len(bits) != len(expected) {
		t.Fatalf("Bit length = %d, expected %d", len(bits), len(expected))
	}

	for i := range bits {
		if bits[i] != expected[i] {
			t.Errorf("U8ToBits(5)[%d] = %v, expected %v", i, bits[i], expected[i])
		}
	}
}
