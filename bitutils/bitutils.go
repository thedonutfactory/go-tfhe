package bitutils

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
)

// Convert converts a slice of bits to a number
// Bits are in little-endian order (LSB first)
func ConvertU8(bits []bool) uint8 {
	var result uint8
	for i := len(bits) - 1; i >= 0; i-- {
		result <<= 1
		if bits[i] {
			result |= 1
		}
	}
	return result
}

func ConvertU16(bits []bool) uint16 {
	var result uint16
	for i := len(bits) - 1; i >= 0; i-- {
		result <<= 1
		if bits[i] {
			result |= 1
		}
	}
	return result
}

func ConvertU32(bits []bool) uint32 {
	var result uint32
	for i := len(bits) - 1; i >= 0; i-- {
		result <<= 1
		if bits[i] {
			result |= 1
		}
	}
	return result
}

func ConvertU64(bits []bool) uint64 {
	var result uint64
	for i := len(bits) - 1; i >= 0; i-- {
		result <<= 1
		if bits[i] {
			result |= 1
		}
	}
	return result
}

// ToBits converts a number to a slice of bits
// Returns bits in little-endian order (LSB first)
func ToBits(val uint64, size int) []bool {
	vec := make([]bool, size)
	for i := 0; i < size; i++ {
		vec[i] = ((val >> i) & 1) != 0
	}
	return vec
}

// U8ToBits converts a uint8 to a slice of bits
func U8ToBits(val uint8) []bool {
	return ToBits(uint64(val), 8)
}

// U16ToBits converts a uint16 to a slice of bits
func U16ToBits(val uint16) []bool {
	return ToBits(uint64(val), 16)
}

// U32ToBits converts a uint32 to a slice of bits
func U32ToBits(val uint32) []bool {
	return ToBits(uint64(val), 32)
}

// U64ToBits converts a uint64 to a slice of bits
func U64ToBits(val uint64) []bool {
	return ToBits(val, 64)
}

// EncryptBits encrypts a slice of bits using the given secret key
func EncryptBits(bits []bool, alpha float64, key []params.Torus) []*tlwe.TLWELv0 {
	result := make([]*tlwe.TLWELv0, len(bits))
	for i, bit := range bits {
		result[i] = tlwe.NewTLWELv0().EncryptBool(bit, alpha, key)
	}
	return result
}

// DecryptBits decrypts a slice of ciphertexts to bits
func DecryptBits(ctxts []*tlwe.TLWELv0, key []params.Torus) []bool {
	result := make([]bool, len(ctxts))
	for i, ctxt := range ctxts {
		result[i] = ctxt.DecryptBool(key)
	}
	return result
}
