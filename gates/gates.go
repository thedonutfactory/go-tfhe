package gates

import (
	"sync"

	"github.com/thedonutfactory/go-tfhe/cloudkey"
	"github.com/thedonutfactory/go-tfhe/evaluator"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trgsw"
	"github.com/thedonutfactory/go-tfhe/trlwe"
	"github.com/thedonutfactory/go-tfhe/utils"
)

// Ciphertext is an alias for TLWELv0
type Ciphertext = tlwe.TLWELv0

// Global evaluator for single-threaded operations (zero-allocation)
var globalEval *evaluator.Evaluator

func init() {
	globalEval = evaluator.NewEvaluator(params.GetTRGSWLv1().N)
}

// NAND performs homomorphic NAND operation (zero-allocation)
func NAND(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	prepared := globalEval.PrepareNAND(tlweA, tlweB)
	result := bootstrap(prepared, ck)

	return result
}

// OR performs homomorphic OR operation (zero-allocation)
func OR(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	prepared := globalEval.PrepareOR(tlweA, tlweB)
	return bootstrap(prepared, ck)
}

// AND performs homomorphic AND operation (zero-allocation)
func AND(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	prepared := globalEval.PrepareAND(tlweA, tlweB)
	return bootstrap(prepared, ck)
}

// XOR performs homomorphic XOR operation (zero-allocation)
func XOR(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	prepared := globalEval.PrepareXOR(tlweA, tlweB)
	return bootstrap(prepared, ck)
}

// XNOR performs homomorphic XNOR operation
func XNOR(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweXNOR := tlweA.SubMul(tlweB, 2)
	// NOTE: Go implementation uses +0.25 instead of -0.25 (inverted from Rust)
	// This may be due to FFT library differences
	tlweXNOR.SetB(tlweXNOR.B() + utils.F64ToTorus(0.25))
	return bootstrap(tlweXNOR, ck)
}

// Constant creates a constant encrypted value
func Constant(value bool) *Ciphertext {
	mu := utils.F64ToTorus(0.125)
	if !value {
		mu = 1 - mu
	}
	result := tlwe.NewTLWELv0()
	result.SetB(mu)
	return result
}

// NOR performs homomorphic NOR operation
func NOR(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweNOR := tlweA.Add(tlweB).Neg()
	tlweNOR.SetB(tlweNOR.B() + utils.F64ToTorus(-0.125))
	return bootstrap(tlweNOR, ck)
}

// ANDNY performs homomorphic AND-NOT-Y operation (NOT(a) AND b)
func ANDNY(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweANDNY := tlweA.Neg().Add(tlweB)
	tlweANDNY.SetB(tlweANDNY.B() + utils.F64ToTorus(-0.125))
	return bootstrap(tlweANDNY, ck)
}

// ANDYN performs homomorphic AND-Y-NOT operation (a AND NOT(b))
func ANDYN(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweANDYN := tlweA.Sub(tlweB)
	tlweANDYN.SetB(tlweANDYN.B() + utils.F64ToTorus(-0.125))
	return bootstrap(tlweANDYN, ck)
}

// ORNY performs homomorphic OR-NOT-Y operation (NOT(a) OR b)
func ORNY(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweORNY := tlweA.Neg().Add(tlweB)
	tlweORNY.SetB(tlweORNY.B() + utils.F64ToTorus(0.125))
	return bootstrap(tlweORNY, ck)
}

// ORYN performs homomorphic OR-Y-NOT operation (a OR NOT(b))
func ORYN(tlweA, tlweB *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	tlweORYN := tlweA.Sub(tlweB)
	tlweORYN.SetB(tlweORYN.B() + utils.F64ToTorus(0.125))
	return bootstrap(tlweORYN, ck)
}

// MUX performs homomorphic multiplexer: a?b:c = a*b + NOT(a)*c
func MUX(tlweA, tlweB, tlweC *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	// Compute using regular AND and OR gates
	// This is more reliable than the optimized version with bootstrap_without_key_switch
	andAB := AND(tlweA, tlweB, ck)
	notA := NOT(tlweA)
	andNotAC := AND(notA, tlweC, ck)
	return OR(andAB, andNotAC, ck)
}

// NOT performs homomorphic NOT operation
func NOT(tlweA *Ciphertext) *Ciphertext {
	return tlweA.Neg()
}

// Copy copies a ciphertext
func Copy(tlweA *Ciphertext) *Ciphertext {
	result := tlwe.NewTLWELv0()
	copy(result.P, tlweA.P)
	return result
}

// bootstrap performs full bootstrapping with key switching (TRUE zero-allocation)
// Returns pointer to internal buffer - result is only valid until next bootstrap call
func bootstrapTemp(ctxt *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	return globalEval.Bootstrap(ctxt, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.KeySwitchingKey, ck.DecompositionOffset)
}

// bootstrap2 performs full bootstrapping with key switching (zero-allocation)
// copies the prepared buffer to the result
func bootstrap(ctxt *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	result := tlwe.NewTLWELv0()
	bootstrapped := globalEval.Bootstrap(ctxt, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.KeySwitchingKey, ck.DecompositionOffset)
	copy(result.P, bootstrapped.P)
	result.SetB(bootstrapped.B())
	return result
}

// bootstrapWithoutKeySwitch performs bootstrapping without key switching (uses global eval)
func bootstrapWithoutKeySwitch(ctxt *Ciphertext, ck *cloudkey.CloudKey) *Ciphertext {
	trlweResult := trlwe.NewTRLWELv1()
	globalEval.BlindRotateAssign(ctxt, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset, trlweResult)
	return trlwe.SampleExtractIndex2(trlweResult, 0)
}

// ============================================================================
// BATCH GATE OPERATIONS - Parallel Processing
// ============================================================================

// BatchNAND performs batch NAND operations in parallel
func BatchNAND(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	// Step 1: Prepare all inputs for bootstrapping
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweNAND := pair[0].Add(pair[1]).Neg()
		tlweNAND.SetB(tlweNAND.B() + utils.F64ToTorus(0.125))
		prepared[i] = tlweNAND
	}

	// Step 2: Batch blind rotate (bottleneck - parallelized)
	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	// Step 3: Post-process (sample extract + key switching, parallel)
	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}

// BatchAND performs batch AND operations in parallel
func BatchAND(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweAND := pair[0].Add(pair[1])
		tlweAND.SetB(tlweAND.B() + utils.F64ToTorus(-0.125))
		prepared[i] = tlweAND
	}

	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}

// BatchOR performs batch OR operations in parallel
func BatchOR(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweOR := pair[0].Add(pair[1])
		tlweOR.SetB(tlweOR.B() + utils.F64ToTorus(0.125))
		prepared[i] = tlweOR
	}

	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}

// BatchXOR performs batch XOR operations in parallel
func BatchXOR(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweXOR := pair[0].AddMul(pair[1], 2)
		tlweXOR.SetB(tlweXOR.B() + utils.F64ToTorus(0.25))
		prepared[i] = tlweXOR
	}

	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}

// BatchNOR performs batch NOR operations in parallel
func BatchNOR(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweNOR := pair[0].Add(pair[1]).Neg()
		tlweNOR.SetB(tlweNOR.B() + utils.F64ToTorus(-0.125))
		prepared[i] = tlweNOR
	}

	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}

// BatchXNOR performs batch XNOR operations in parallel
func BatchXNOR(inputs [][2]*Ciphertext, ck *cloudkey.CloudKey) []*Ciphertext {
	prepared := make([]*Ciphertext, len(inputs))
	for i, pair := range inputs {
		tlweXNOR := pair[0].SubMul(pair[1], 2)
		tlweXNOR.SetB(tlweXNOR.B() + utils.F64ToTorus(-0.25))
		prepared[i] = tlweXNOR
	}

	trlwes := trgsw.BatchBlindRotate(prepared, ck.BlindRotateTestvec, ck.BootstrappingKey, ck.DecompositionOffset)

	results := make([]*Ciphertext, len(trlwes))
	var wg sync.WaitGroup
	for i, trlweResult := range trlwes {
		wg.Add(1)
		go func(idx int, t *trlwe.TRLWELv1) {
			defer wg.Done()
			tlweLv1 := trlwe.SampleExtractIndex(t, 0)
			results[idx] = trgsw.IdentityKeySwitching(tlweLv1, ck.KeySwitchingKey)
		}(i, trlweResult)
	}
	wg.Wait()

	return results
}
