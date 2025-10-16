package poly

// MulPoly returns p0 * p1.
func (e *Evaluator) MulPoly(p0, p1 Poly) Poly {
	pOut := e.NewPoly()
	e.MulPolyAssign(p0, p1, pOut)
	return pOut
}

// MulPolyAssign computes pOut = p0 * p1.
// This uses FFT-based multiplication for efficiency.
func (e *Evaluator) MulPolyAssign(p0, p1, pOut Poly) {
	// Transform both polynomials to frequency domain
	fp0 := e.ToFourierPoly(p0)
	fp1 := e.ToFourierPoly(p1)

	// Multiply in frequency domain (element-wise complex multiplication)
	e.MulFourierPolyAssign(fp0, fp1, fp0)

	// Transform back to time domain
	e.ToPolyAssignUnsafe(fp0, pOut)
}

// MulAddPolyAssign computes pOut += p0 * p1.
func (e *Evaluator) MulAddPolyAssign(p0, p1, pOut Poly) {
	fp0 := e.ToFourierPoly(p0)
	fp1 := e.ToFourierPoly(p1)
	e.MulFourierPolyAssign(fp0, fp1, fp0)
	e.ToPolyAddAssignUnsafe(fp0, pOut)
}

// MulSubPolyAssign computes pOut -= p0 * p1.
func (e *Evaluator) MulSubPolyAssign(p0, p1, pOut Poly) {
	fp0 := e.ToFourierPoly(p0)
	fp1 := e.ToFourierPoly(p1)
	e.MulFourierPolyAssign(fp0, fp1, fp0)
	e.ToPolySubAssignUnsafe(fp0, pOut)
}
