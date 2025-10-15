package tlwe

import (
	"math/rand"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// TLWELv0 represents a Level 0 TLWE ciphertext
type TLWELv0 struct {
	P []params.Torus // Length is N+1, where last element is b
}

// NewTLWELv0 creates a new TLWE Level 0 ciphertext
func NewTLWELv0() *TLWELv0 {
	n := params.GetTLWELv0().N
	return &TLWELv0{
		P: make([]params.Torus, n+1),
	}
}

// B returns the b component of the TLWE ciphertext
func (t *TLWELv0) B() params.Torus {
	n := params.GetTLWELv0().N
	return t.P[n]
}

// SetB sets the b component of the TLWE ciphertext
func (t *TLWELv0) SetB(val params.Torus) {
	n := params.GetTLWELv0().N
	t.P[n] = val
}

// EncryptF64 encrypts a float64 value with TLWE Level 0
func (t *TLWELv0) EncryptF64(p float64, alpha float64, key []params.Torus) *TLWELv0 {
	rng := rand.New(rand.NewSource(rand.Int63()))
	n := params.GetTLWELv0().N

	var innerProduct params.Torus
	for i := 0; i < n; i++ {
		randU32 := params.Torus(rng.Uint32())
		innerProduct += key[i] * randU32
		t.P[i] = randU32
	}

	b := utils.GaussianF64(p, alpha, rng)
	t.SetB(innerProduct + b)
	return t
}

// EncryptBool encrypts a boolean value with TLWE Level 0
func (t *TLWELv0) EncryptBool(pBool bool, alpha float64, key []params.Torus) *TLWELv0 {
	var p float64
	if pBool {
		p = 0.125
	} else {
		p = -0.125
	}
	return t.EncryptF64(p, alpha, key)
}

// DecryptBool decrypts a TLWE Level 0 ciphertext to a boolean
func (t *TLWELv0) DecryptBool(key []params.Torus) bool {
	n := params.GetTLWELv0().N
	var innerProduct params.Torus
	for i := 0; i < n; i++ {
		innerProduct += t.P[i] * key[i]
	}

	resTorus := int32(t.P[n] - innerProduct)
	return resTorus >= 0
}

// Add adds two TLWE Level 0 ciphertexts
func (t *TLWELv0) Add(other *TLWELv0) *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = t.P[i] + other.P[i]
	}
	return result
}

// Sub subtracts two TLWE Level 0 ciphertexts
func (t *TLWELv0) Sub(other *TLWELv0) *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = t.P[i] - other.P[i]
	}
	return result
}

// Neg negates a TLWE Level 0 ciphertext
func (t *TLWELv0) Neg() *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = 0 - t.P[i]
	}
	return result
}

// Mul multiplies two TLWE Level 0 ciphertexts (element-wise)
func (t *TLWELv0) Mul(other *TLWELv0) *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = t.P[i] * other.P[i]
	}
	return result
}

// AddMul adds a TLWE ciphertext multiplied by a constant
func (t *TLWELv0) AddMul(other *TLWELv0, multiplier params.Torus) *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = t.P[i] + (other.P[i] * multiplier)
	}
	return result
}

// SubMul subtracts a TLWE ciphertext multiplied by a constant
func (t *TLWELv0) SubMul(other *TLWELv0, multiplier params.Torus) *TLWELv0 {
	result := NewTLWELv0()
	for i := range result.P {
		result.P[i] = t.P[i] - (other.P[i] * multiplier)
	}
	return result
}

// TLWELv1 represents a Level 1 TLWE ciphertext
type TLWELv1 struct {
	P []params.Torus // Length is N+1, where last element is b
}

// NewTLWELv1 creates a new TLWE Level 1 ciphertext
func NewTLWELv1() *TLWELv1 {
	n := params.GetTLWELv1().N
	return &TLWELv1{
		P: make([]params.Torus, n+1),
	}
}

// SetB sets the b component of the TLWE Level 1 ciphertext
func (t *TLWELv1) SetB(val params.Torus) {
	n := params.GetTLWELv1().N
	t.P[n] = val
}

// EncryptF64 encrypts a float64 value with TLWE Level 1
func (t *TLWELv1) EncryptF64(p float64, alpha float64, key []params.Torus) *TLWELv1 {
	rng := rand.New(rand.NewSource(rand.Int63()))
	n := params.GetTLWELv1().N

	var innerProduct params.Torus
	for i := 0; i < n; i++ {
		randU32 := params.Torus(rng.Uint32())
		innerProduct += key[i] * randU32
		t.P[i] = randU32
	}

	b := utils.GaussianF64(p, alpha, rng)
	t.SetB(innerProduct + b)
	return t
}

// EncryptBool encrypts a boolean value with TLWE Level 1
func (t *TLWELv1) EncryptBool(pBool bool, alpha float64, key []params.Torus) *TLWELv1 {
	var p float64
	if pBool {
		p = 0.125
	} else {
		p = -0.125
	}
	return t.EncryptF64(p, alpha, key)
}

// DecryptBool decrypts a TLWE Level 1 ciphertext to a boolean
func (t *TLWELv1) DecryptBool(key []params.Torus) bool {
	n := params.GetTLWELv1().N
	var innerProduct params.Torus
	for i := 0; i < n; i++ {
		innerProduct += t.P[i] * key[i]
	}

	resTorus := int32(t.P[len(key)] - innerProduct)
	return resTorus >= 0
}
