package evaluator

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// PrepareNAND prepares a NAND input for bootstrapping (zero-allocation)
func (e *Evaluator) PrepareNAND(a, b *tlwe.TLWELv0) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N
	result := tlwe.NewTLWELv0()

	// NAND: -(a + b) + 1/8
	for i := 0; i < n; i++ {
		result.P[i] = -(a.P[i] + b.P[i])
	}
	result.P[n] = -(a.P[n] + b.P[n]) + utils.F64ToTorus(0.125)

	return result
}

// PrepareAND prepares an AND input for bootstrapping
func (e *Evaluator) PrepareAND(a, b *tlwe.TLWELv0) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N
	result := tlwe.NewTLWELv0()

	// AND: (a + b) - 1/8
	for i := 0; i < n; i++ {
		result.P[i] = a.P[i] + b.P[i]
	}
	result.P[n] = a.P[n] + b.P[n] + utils.F64ToTorus(-0.125)

	return result
}

// PrepareOR prepares an OR input for bootstrapping
func (e *Evaluator) PrepareOR(a, b *tlwe.TLWELv0) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N
	result := tlwe.NewTLWELv0()

	// OR: (a + b) + 1/8
	for i := 0; i < n; i++ {
		result.P[i] = a.P[i] + b.P[i]
	}
	result.P[n] = a.P[n] + b.P[n] + utils.F64ToTorus(0.125)

	return result
}

// PrepareXOR prepares an XOR input for bootstrapping
func (e *Evaluator) PrepareXOR(a, b *tlwe.TLWELv0) *tlwe.TLWELv0 {
	n := params.GetTLWELv0().N
	result := tlwe.NewTLWELv0()

	// XOR: (a + 2*b) + 1/4
	for i := 0; i < n; i++ {
		result.P[i] = a.P[i] + 2*b.P[i]
	}
	result.P[n] = a.P[n] + 2*b.P[n] + utils.F64ToTorus(0.25)

	return result
}
