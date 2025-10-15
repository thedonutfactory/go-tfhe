package gates_test

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/gates"
	"github.com/thedonutfactory/go-tfhe/key"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// Test helpers
func encrypt(t *testing.T, val bool, sk *key.SecretKey) *gates.Ciphertext {
	return tlwe.NewTLWELv0().EncryptBool(val, params.GetTLWELv0().ALPHA, sk.KeyLv0)
}

func decrypt(t *testing.T, ct *gates.Ciphertext, sk *key.SecretKey) bool {
	return ct.DecryptBool(sk.KeyLv0)
}

// TestNAND tests the NAND gate
func TestNAND(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, true},
		{false, true, true},
		{true, false, true},
		{true, true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.NAND(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("NAND(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestAND tests the AND gate
func TestAND(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, false},
		{false, true, false},
		{true, false, false},
		{true, true, true},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.AND(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("AND(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestOR tests the OR gate
func TestOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, true},
		{true, true, true},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.OR(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("OR(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestXOR tests the XOR gate
func TestXOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, true},
		{true, true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.XOR(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("XOR(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestXNOR tests the XNOR gate
func TestXNOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, true},
		{false, true, false},
		{true, false, false},
		{true, true, true},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.XNOR(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("XNOR(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestNOR tests the NOR gate
func TestNOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, true},
		{false, true, false},
		{true, false, false},
		{true, true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.NOR(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("NOR(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestANDNY tests the ANDNY gate (NOT(a) AND b)
func TestANDNY(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, false},
		{false, true, true},
		{true, false, false},
		{true, true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.ANDNY(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("ANDNY(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestANDYN tests the ANDYN gate (a AND NOT(b))
func TestANDYN(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, false},
		{false, true, false},
		{true, false, true},
		{true, true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.ANDYN(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("ANDYN(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestORNY tests the ORNY gate (NOT(a) OR b)
func TestORNY(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, true},
		{false, true, true},
		{true, false, false},
		{true, true, true},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.ORNY(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("ORNY(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestORYN tests the ORYN gate (a OR NOT(b))
func TestORYN(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		a, b, expected bool
	}{
		{false, false, true},
		{false, true, false},
		{true, false, true},
		{true, true, true},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.ORYN(ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("ORYN(%v, %v) = %v, expected %v", tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestNOT tests the NOT gate
func TestNOT(t *testing.T) {
	sk := key.NewSecretKey()

	testCases := []struct {
		a, expected bool
	}{
		{false, true},
		{true, false},
	}

	for _, tc := range testCases {
		ctA := encrypt(t, tc.a, sk)
		result := gates.NOT(ctA)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("NOT(%v) = %v, expected %v", tc.a, dec, tc.expected)
		}
	}
}

// TestCopy tests the Copy operation
func TestCopy(t *testing.T) {
	sk := key.NewSecretKey()

	testCases := []bool{false, true}

	for _, tc := range testCases {
		ctA := encrypt(t, tc, sk)
		result := gates.Copy(ctA)
		dec := decrypt(t, result, sk)

		if dec != tc {
			t.Errorf("Copy(%v) = %v, expected %v", tc, dec, tc)
		}
	}
}

// TestConstant tests the Constant operation
func TestConstant(t *testing.T) {
	sk := key.NewSecretKey()

	testCases := []bool{false, true}

	for _, tc := range testCases {
		result := gates.Constant(tc)
		dec := decrypt(t, result, sk)

		if dec != tc {
			t.Errorf("Constant(%v) = %v, expected %v", tc, dec, tc)
		}
	}
}

// TestMUX tests the MUX gate
func TestMUX(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := []struct {
		sel, a, b, expected bool
	}{
		{false, false, false, false}, // sel=0 -> b
		{false, false, true, true},   // sel=0 -> b
		{false, true, false, false},  // sel=0 -> b
		{false, true, true, true},    // sel=0 -> b
		{true, false, false, false},  // sel=1 -> a
		{true, false, true, false},   // sel=1 -> a
		{true, true, false, true},    // sel=1 -> a
		{true, true, true, true},     // sel=1 -> a
	}

	for _, tc := range testCases {
		ctSel := encrypt(t, tc.sel, sk)
		ctA := encrypt(t, tc.a, sk)
		ctB := encrypt(t, tc.b, sk)
		result := gates.MUX(ctSel, ctA, ctB, ck)
		dec := decrypt(t, result, sk)

		if dec != tc.expected {
			t.Errorf("MUX(sel=%v, a=%v, b=%v) = %v, expected %v", tc.sel, tc.a, tc.b, dec, tc.expected)
		}
	}
}

// TestBatchAND tests the batch AND operation
func TestBatchAND(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := [][2]bool{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	inputs := make([][2]*gates.Ciphertext, len(testCases))
	expected := make([]bool, len(testCases))

	for i, tc := range testCases {
		inputs[i] = [2]*gates.Ciphertext{
			encrypt(t, tc[0], sk),
			encrypt(t, tc[1], sk),
		}
		expected[i] = tc[0] && tc[1]
	}

	results := gates.BatchAND(inputs, ck)

	if len(results) != len(expected) {
		t.Fatalf("BatchAND returned %d results, expected %d", len(results), len(expected))
	}

	for i, result := range results {
		dec := decrypt(t, result, sk)
		if dec != expected[i] {
			t.Errorf("BatchAND[%d](%v AND %v) = %v, expected %v",
				i, testCases[i][0], testCases[i][1], dec, expected[i])
		}
	}
}

// TestBatchOR tests the batch OR operation
func TestBatchOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := [][2]bool{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	inputs := make([][2]*gates.Ciphertext, len(testCases))
	expected := make([]bool, len(testCases))

	for i, tc := range testCases {
		inputs[i] = [2]*gates.Ciphertext{
			encrypt(t, tc[0], sk),
			encrypt(t, tc[1], sk),
		}
		expected[i] = tc[0] || tc[1]
	}

	results := gates.BatchOR(inputs, ck)

	if len(results) != len(expected) {
		t.Fatalf("BatchOR returned %d results, expected %d", len(results), len(expected))
	}

	for i, result := range results {
		dec := decrypt(t, result, sk)
		if dec != expected[i] {
			t.Errorf("BatchOR[%d](%v OR %v) = %v, expected %v",
				i, testCases[i][0], testCases[i][1], dec, expected[i])
		}
	}
}

// TestBatchXOR tests the batch XOR operation
func TestBatchXOR(t *testing.T) {
	sk := key.NewSecretKey()
	ck := cloudkey.NewCloudKey(sk)

	testCases := [][2]bool{
		{false, false},
		{false, true},
		{true, false},
		{true, true},
	}

	inputs := make([][2]*gates.Ciphertext, len(testCases))
	expected := make([]bool, len(testCases))

	for i, tc := range testCases {
		inputs[i] = [2]*gates.Ciphertext{
			encrypt(t, tc[0], sk),
			encrypt(t, tc[1], sk),
		}
		expected[i] = tc[0] != tc[1]
	}

	results := gates.BatchXOR(inputs, ck)

	if len(results) != len(expected) {
		t.Fatalf("BatchXOR returned %d results, expected %d", len(results), len(expected))
	}

	for i, result := range results {
		dec := decrypt(t, result, sk)
		if dec != expected[i] {
			t.Errorf("BatchXOR[%d](%v XOR %v) = %v, expected %v",
				i, testCases[i][0], testCases[i][1], dec, expected[i])
		}
	}
}
