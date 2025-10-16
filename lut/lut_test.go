package lut

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

func TestLookUpTableBasic(t *testing.T) {
	// Test creation and basic operations
	lut := NewLookUpTable()

	if lut == nil {
		t.Fatal("NewLookUpTable returned nil")
	}

	if lut.Poly == nil {
		t.Fatal("LookUpTable polynomial is nil")
	}

	// Test clear
	lut.Poly.B[0] = 123
	lut.Clear()
	if lut.Poly.B[0] != 0 {
		t.Error("Clear did not clear the polynomial")
	}
}

func TestLookUpTableCopy(t *testing.T) {
	lut1 := NewLookUpTable()
	lut1.Poly.B[0] = 42
	lut1.Poly.A[0] = 17

	// Test Copy
	lut2 := lut1.Copy()
	if lut2.Poly.B[0] != 42 || lut2.Poly.A[0] != 17 {
		t.Error("Copy did not copy values correctly")
	}

	// Modify original and ensure copy is unchanged
	lut1.Poly.B[0] = 99
	if lut2.Poly.B[0] != 42 {
		t.Error("Copy is not independent of original")
	}

	// Test CopyFrom
	lut3 := NewLookUpTable()
	lut3.CopyFrom(lut1)
	if lut3.Poly.B[0] != 99 {
		t.Error("CopyFrom did not copy values correctly")
	}
}

func TestEncoder(t *testing.T) {
	// Test binary encoder (message modulus = 2)
	enc := NewEncoder(2)

	// Test encoding
	val0 := enc.Encode(0)
	val1 := enc.Encode(1)

	// Values should be different
	if val0 == val1 {
		t.Error("Encoded values for 0 and 1 should be different")
	}

	// Test decoding
	if enc.Decode(val0) != 0 {
		t.Errorf("Decode(Encode(0)) = %d, want 0", enc.Decode(val0))
	}
	if enc.Decode(val1) != 1 {
		t.Errorf("Decode(Encode(1)) = %d, want 1", enc.Decode(val1))
	}

	// Test DecodeBool
	if enc.DecodeBool(val0) != false {
		t.Error("DecodeBool(Encode(0)) should be false")
	}
	if enc.DecodeBool(val1) != true {
		t.Error("DecodeBool(Encode(1)) should be true")
	}
}

func TestEncoderModular(t *testing.T) {
	// Test with message modulus = 4
	enc := NewEncoder(4)

	for i := 0; i < 4; i++ {
		encoded := enc.Encode(i)
		decoded := enc.Decode(encoded)
		if decoded != i {
			t.Errorf("Encode/Decode(%d) = %d, want %d", i, decoded, i)
		}
	}

	// Test negative wrapping
	if enc.Encode(-1) != enc.Encode(3) {
		t.Error("Negative values should wrap modulo MessageModulus")
	}

	// Test overflow wrapping
	if enc.Encode(4) != enc.Encode(0) {
		t.Error("Values >= MessageModulus should wrap")
	}
}

func TestGeneratorIdentity(t *testing.T) {
	// Test identity function (f(x) = x)
	gen := NewGenerator(4)

	identity := func(x int) int { return x }
	lut := gen.GenLookUpTable(identity)

	if lut == nil {
		t.Fatal("GenLookUpTable returned nil")
	}

	// Lookup table should be created without error
	// Detailed functional testing requires full TFHE stack
}

func TestGeneratorConstant(t *testing.T) {
	// Test constant function (f(x) = c)
	gen := NewGenerator(2)

	constantOne := func(x int) int { return 1 }
	lut := gen.GenLookUpTable(constantOne)

	if lut == nil {
		t.Fatal("GenLookUpTable returned nil")
	}

	// All values should encode to the same constant
	// Detailed verification requires full TFHE stack
}

func TestGeneratorNOT(t *testing.T) {
	// Test NOT function for binary (f(x) = 1 - x)
	gen := NewGenerator(2)

	notFunc := func(x int) int { return 1 - x }
	lut := gen.GenLookUpTable(notFunc)

	if lut == nil {
		t.Fatal("GenLookUpTable returned nil")
	}
}

func TestGeneratorCustomModulus(t *testing.T) {
	// Test with custom message modulus
	gen := NewGenerator(8)

	// Function that doubles the input mod 8
	doubleFunc := func(x int) int { return (2 * x) % 8 }
	lut := gen.GenLookUpTableCustom(doubleFunc, 8, 1.0/16.0)

	if lut == nil {
		t.Fatal("GenLookUpTableCustom returned nil")
	}
}

func TestModSwitch(t *testing.T) {
	gen := NewGenerator(2)
	n := params.GetTRGSWLv1().N

	// Test modulus switching at key points
	tests := []struct {
		name  string
		input params.Torus
	}{
		{"zero", 0},
		{"quarter", params.Torus(1 << 30)},
		{"half", params.Torus(1 << 31)},
		{"three-quarters", params.Torus(3 << 30)},
		{"max", params.Torus(^uint32(0))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.ModSwitch(tt.input)

			// Result should be in valid range
			if result < 0 || result >= 2*n {
				t.Errorf("ModSwitch(%d) = %d, out of range [0, %d)", tt.input, result, 2*n)
			}
		})
	}
}

func TestGeneratorFullControl(t *testing.T) {
	// Test GenLookUpTableFull for fine-grained control
	gen := NewGenerator(2)

	// Function that returns exact torus values
	fullFunc := func(x int) params.Torus {
		if x == 0 {
			return utils.F64ToTorus(0.0)
		}
		return utils.F64ToTorus(0.25)
	}

	lut := gen.GenLookUpTableFull(fullFunc)

	if lut == nil {
		t.Fatal("GenLookUpTableFull returned nil")
	}
}

func BenchmarkLookUpTableCreation(b *testing.B) {
	gen := NewGenerator(2)
	identity := func(x int) int { return x }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.GenLookUpTable(identity)
	}
}

func BenchmarkModSwitch(b *testing.B) {
	gen := NewGenerator(2)
	testVal := params.Torus(12345678)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = gen.ModSwitch(testVal)
	}
}

func BenchmarkEncode(b *testing.B) {
	enc := NewEncoder(2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enc.Encode(i % 2)
	}
}

func BenchmarkDecode(b *testing.B) {
	enc := NewEncoder(2)
	testVal := enc.Encode(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = enc.Decode(testVal)
	}
}
