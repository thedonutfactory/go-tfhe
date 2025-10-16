package evaluator

import (
	"github.com/thedonutfactory/go-tfhe/lut"
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trgsw"
	"github.com/thedonutfactory/go-tfhe/trlwe"
)

// BootstrapFunc performs programmable bootstrapping with a function
// The function f operates on the message space [0, messageModulus) and
// is evaluated homomorphically on the encrypted data during bootstrapping.
//
// This combines noise refreshing with arbitrary function evaluation.
func (e *Evaluator) BootstrapFunc(
	ctIn *tlwe.TLWELv0,
	f func(int) int,
	messageModulus int,
	bsk []*trgsw.TRGSWLv1FFT,
	ksk []*tlwe.TLWELv0,
	decompositionOffset params.Torus,
) *tlwe.TLWELv0 {
	// Generate lookup table from function
	generator := lut.NewGenerator(messageModulus)
	lookupTable := generator.GenLookUpTable(f)

	// Perform LUT-based bootstrapping
	return e.BootstrapLUT(ctIn, lookupTable, bsk, ksk, decompositionOffset)
}

// BootstrapFuncAssign performs programmable bootstrapping with a function (zero-allocation)
func (e *Evaluator) BootstrapFuncAssign(
	ctIn *tlwe.TLWELv0,
	f func(int) int,
	messageModulus int,
	bsk []*trgsw.TRGSWLv1FFT,
	ksk []*tlwe.TLWELv0,
	decompositionOffset params.Torus,
	ctOut *tlwe.TLWELv0,
) {
	// Generate lookup table from function
	generator := lut.NewGenerator(messageModulus)
	lookupTable := generator.GenLookUpTable(f)

	// Perform LUT-based bootstrapping
	e.BootstrapLUTAssign(ctIn, lookupTable, bsk, ksk, decompositionOffset, ctOut)
}

// BootstrapLUT performs programmable bootstrapping with a pre-computed lookup table
// The lookup table encodes the function to be evaluated during bootstrapping.
//
// This is more efficient than BootstrapFunc when the same function is used multiple times.
func (e *Evaluator) BootstrapLUT(
	ctIn *tlwe.TLWELv0,
	lut *lut.LookUpTable,
	bsk []*trgsw.TRGSWLv1FFT,
	ksk []*tlwe.TLWELv0,
	decompositionOffset params.Torus,
) *tlwe.TLWELv0 {
	result := e.Buffers.GetNextResult()
	e.BootstrapLUTAssign(ctIn, lut, bsk, ksk, decompositionOffset, result)

	copiedResult := tlwe.NewTLWELv0()
	copy(copiedResult.P, result.P)
	copiedResult.SetB(result.B())

	return copiedResult
}

func (e *Evaluator) BootstrapLUTTemp(
	ctIn *tlwe.TLWELv0,
	lut *lut.LookUpTable,
	bsk []*trgsw.TRGSWLv1FFT,
	ksk []*tlwe.TLWELv0,
	decompositionOffset params.Torus,
) *tlwe.TLWELv0 {
	result := e.Buffers.GetNextResult()
	e.BootstrapLUTAssign(ctIn, lut, bsk, ksk, decompositionOffset, result)
	return result
}

// BootstrapLUTAssign performs programmable bootstrapping with a lookup table (zero-allocation)
// This is the core implementation of programmable bootstrapping.
//
// Algorithm:
// 1. Blind rotate the lookup table based on the encrypted value
// 2. Sample extract to get an LWE ciphertext
// 3. Key switch to convert back to the original key
//
// The key insight is that we can reuse the existing BlindRotateAssign function
// by converting the LUT into a TRLWE ciphertext (test vector).
func (e *Evaluator) BootstrapLUTAssign(
	ctIn *tlwe.TLWELv0,
	lut *lut.LookUpTable,
	bsk []*trgsw.TRGSWLv1FFT,
	ksk []*tlwe.TLWELv0,
	decompositionOffset params.Torus,
	ctOut *tlwe.TLWELv0,
) {
	// Convert LUT to TRLWE format (test vector)
	// The LUT is already a TRLWE with the function encoded in the B polynomial
	testvec := lut.Poly

	// Perform blind rotation using the LUT as the test vector
	// This rotates the LUT based on the encrypted value, effectively evaluating the function
	e.BlindRotateAssign(ctIn, testvec, bsk, decompositionOffset, e.Buffers.BlindRotation.Rotated)

	// Extract the constant term as an LWE ciphertext
	// This gives us the function evaluation encrypted under the TRLWE key
	trlwe.SampleExtractIndexAssign(e.Buffers.BlindRotation.Rotated, 0, e.Buffers.Bootstrap.ExtractedLWE)

	// Key switch to convert back to the original LWE key
	trgsw.IdentityKeySwitchingAssign(e.Buffers.Bootstrap.ExtractedLWE, ksk, ctOut)
}
