package tlwe

import (
	"github.com/thedonutfactory/go-tfhe/params"
)

// EncryptLWEMessage encrypts an integer message using general message encoding
// This is different from EncryptBool which uses ±1/8 binary encoding.
//
// For programmable bootstrapping, use this function to match the LUT encoding.
// Encoding: message → message * scale, where scale = 2^31 / messageModulus
func (t *TLWELv0) EncryptLWEMessage(message int, messageModulus int, alpha float64, key []params.Torus) *TLWELv0 {
	// Calculate scale: 2^31 / messageModulus
	scale := float64(uint64(1)<<31) / float64(messageModulus)

	// Normalize message
	message = message % messageModulus
	if message < 0 {
		message += messageModulus
	}

	// Encode: message * scale / 2^32 to get value in [0, 1)
	encodedMessage := float64(message) * scale / float64(uint64(1)<<32)

	return t.EncryptF64(encodedMessage, alpha, key)
}

// DecryptLWEMessage decrypts an integer message using general message encoding
//
// Following the reference implementation: num.DivRound(phase, scale) % messageModulus
// DivRound(a, b) rounds a/b to nearest integer
func (t *TLWELv0) DecryptLWEMessage(messageModulus int, key []params.Torus) int {
	// Calculate scale: 2^31 / messageModulus
	scale := params.Torus(uint64(1)<<31) / params.Torus(messageModulus)

	// Get phase (decrypted value with noise)
	n := params.GetTLWELv0().N
	var innerProduct params.Torus
	for i := 0; i < n; i++ {
		innerProduct += t.P[i] * key[i]
	}
	phase := t.P[n] - innerProduct

	// DivRound: (a + b/2) / b
	// For unsigned: (phase + scale/2) / scale
	decoded := int((phase + scale/2) / scale)

	message := decoded % messageModulus
	if message < 0 {
		message += messageModulus
	}

	return message
}
