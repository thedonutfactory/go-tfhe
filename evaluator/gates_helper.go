package evaluator

import (
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// PrepareAND prepares input for AND gate (zero-allocation)
// Returns pointer to internal temp buffer
func (e *Evaluator) PrepareAND(ctA, ctB *tlwe.TLWELv0) *tlwe.TLWELv0 {
	ctA.AddAssign(ctB, e.Buffers.GatePrep)
	e.Buffers.GatePrep.SetB(e.Buffers.GatePrep.B() + utils.F64ToTorus(-0.125))
	return e.Buffers.GatePrep
}

// PrepareNAND prepares input for NAND gate (zero-allocation)
func (e *Evaluator) PrepareNAND(ctA, ctB *tlwe.TLWELv0) *tlwe.TLWELv0 {
	// Negate both and add
	for i := range e.Buffers.GatePrep.P {
		e.Buffers.GatePrep.P[i] = -ctA.P[i] - ctB.P[i]
	}
	e.Buffers.GatePrep.SetB(e.Buffers.GatePrep.B() + utils.F64ToTorus(0.125))
	return e.Buffers.GatePrep
}

// PrepareOR prepares input for OR gate (zero-allocation)
func (e *Evaluator) PrepareOR(ctA, ctB *tlwe.TLWELv0) *tlwe.TLWELv0 {
	ctA.AddAssign(ctB, e.Buffers.GatePrep)
	e.Buffers.GatePrep.SetB(e.Buffers.GatePrep.B() + utils.F64ToTorus(0.125))
	return e.Buffers.GatePrep
}

// PrepareXOR prepares input for XOR gate (zero-allocation)
func (e *Evaluator) PrepareXOR(ctA, ctB *tlwe.TLWELv0) *tlwe.TLWELv0 {
	for i := range e.Buffers.GatePrep.P {
		e.Buffers.GatePrep.P[i] = 2 * (ctA.P[i] + ctB.P[i])
	}
	e.Buffers.GatePrep.SetB(e.Buffers.GatePrep.B() + utils.F64ToTorus(0.25))
	return e.Buffers.GatePrep
}
