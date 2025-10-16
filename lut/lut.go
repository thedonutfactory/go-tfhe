// Package lut provides LookUpTable support for programmable bootstrapping.
// This enables evaluating arbitrary functions on encrypted data during bootstrapping.
package lut

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/trlwe"
)

// LookUpTable is a TRLWE ciphertext that encodes a function
// for programmable bootstrapping.
// During blind rotation, the LUT is rotated based on the encrypted value,
// effectively evaluating the function on the encrypted data.
type LookUpTable struct {
	// Polynomial encoding the function values
	Poly *trlwe.TRLWELv1
}

// NewLookUpTable creates a new lookup table
func NewLookUpTable() *LookUpTable {
	return &LookUpTable{
		Poly: trlwe.NewTRLWELv1(),
	}
}

// Copy returns a deep copy of the lookup table
func (lut *LookUpTable) Copy() *LookUpTable {
	result := NewLookUpTable()
	copy(result.Poly.A, lut.Poly.A)
	copy(result.Poly.B, lut.Poly.B)
	return result
}

// CopyFrom copies values from another lookup table
func (lut *LookUpTable) CopyFrom(other *LookUpTable) {
	copy(lut.Poly.A, other.Poly.A)
	copy(lut.Poly.B, other.Poly.B)
}

// Clear clears the lookup table (sets all coefficients to 0)
func (lut *LookUpTable) Clear() {
	n := params.GetTRGSWLv1().N
	for i := 0; i < n; i++ {
		lut.Poly.A[i] = 0
		lut.Poly.B[i] = 0
	}
}
