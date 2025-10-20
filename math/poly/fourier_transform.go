package poly

// ToFourierPoly transforms Poly to FourierPoly.
func (e *Evaluator[T]) ToFourierPoly(p Poly[T]) FourierPoly {
	fpOut := NewFourierPoly(e.degree)
	e.ToFourierPolyAssign(p, fpOut)
	return fpOut
}

// ToFourierPolyAssign transforms Poly to FourierPoly and writes it to fpOut.
func (e *Evaluator[T]) ToFourierPolyAssign(p Poly[T], fpOut FourierPoly) {
	convertPolyToFourierPolyAssign(p.Coeffs, fpOut.Coeffs)
	fftInPlace(fpOut.Coeffs, e.tw)
}

// ToFourierPolyAddAssign transforms Poly to FourierPoly and adds it to fpOut.
func (e *Evaluator[T]) ToFourierPolyAddAssign(p Poly[T], fpOut FourierPoly) {
	convertPolyToFourierPolyAssign(p.Coeffs, e.buffer.fp.Coeffs)
	fftInPlace(e.buffer.fp.Coeffs, e.tw)
	addCmplxAssign(fpOut.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// ToFourierPolySubAssign transforms Poly to FourierPoly and subtracts it from fpOut.
func (e *Evaluator[T]) ToFourierPolySubAssign(p Poly[T], fpOut FourierPoly) {
	convertPolyToFourierPolyAssign(p.Coeffs, e.buffer.fp.Coeffs)
	fftInPlace(e.buffer.fp.Coeffs, e.tw)
	subCmplxAssign(fpOut.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// MonomialToFourierPoly transforms X^d to FourierPoly.
func (e *Evaluator[T]) MonomialToFourierPoly(d int) FourierPoly {
	fpOut := NewFourierPoly(e.degree)
	e.MonomialToFourierPolyAssign(d, fpOut)
	return fpOut
}

// MonomialToFourierPolyAssign transforms X^d to FourierPoly and writes it to fpOut.
func (e *Evaluator[T]) MonomialToFourierPolyAssign(d int, fpOut FourierPoly) {
	d &= 2*e.degree - 1
	for j, jj := 0, 0; j < e.degree; j, jj = j+8, jj+4 {
		c0 := e.twMono[(e.twMonoIdx[jj+0]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+0] = real(c0)
		fpOut.Coeffs[j+4] = imag(c0)

		c1 := e.twMono[(e.twMonoIdx[jj+1]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+1] = real(c1)
		fpOut.Coeffs[j+5] = imag(c1)

		c2 := e.twMono[(e.twMonoIdx[jj+2]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+2] = real(c2)
		fpOut.Coeffs[j+6] = imag(c2)

		c3 := e.twMono[(e.twMonoIdx[jj+3]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+3] = real(c3)
		fpOut.Coeffs[j+7] = imag(c3)
	}
}

// MonomialSubOneToFourierPoly transforms X^d-1 to FourierPoly.
//
// d should be positive.
func (e *Evaluator[T]) MonomialSubOneToFourierPoly(d int) FourierPoly {
	fpOut := NewFourierPoly(e.degree)
	e.MonomialSubOneToFourierPolyAssign(d, fpOut)
	return fpOut
}

// MonomialSubOneToFourierPolyAssign transforms X^d-1 to FourierPoly and writes it to fpOut.
func (e *Evaluator[T]) MonomialSubOneToFourierPolyAssign(d int, fpOut FourierPoly) {
	d &= 2*e.degree - 1
	for j, jj := 0, 0; j < e.degree; j, jj = j+8, jj+4 {
		c0 := e.twMono[(e.twMonoIdx[jj+0]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+0] = real(c0) - 1
		fpOut.Coeffs[j+4] = imag(c0)

		c1 := e.twMono[(e.twMonoIdx[jj+1]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+1] = real(c1) - 1
		fpOut.Coeffs[j+5] = imag(c1)

		c2 := e.twMono[(e.twMonoIdx[jj+2]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+2] = real(c2) - 1
		fpOut.Coeffs[j+6] = imag(c2)

		c3 := e.twMono[(e.twMonoIdx[jj+3]*d)&(2*e.degree-1)]
		fpOut.Coeffs[j+3] = real(c3) - 1
		fpOut.Coeffs[j+7] = imag(c3)
	}
}

// ToPoly transforms FourierPoly to Poly.
func (e *Evaluator[T]) ToPoly(fp FourierPoly) Poly[T] {
	pOut := NewPoly[T](e.degree)
	e.ToPolyAssign(fp, pOut)
	return pOut
}

// ToPolyAssign transforms FourierPoly to Poly and writes it to pOut.
func (e *Evaluator[T]) ToPolyAssign(fp FourierPoly, pOut Poly[T]) {
	e.buffer.fpInv.CopyFrom(fp)
	ifftInPlace(e.buffer.fpInv.Coeffs, e.twInv)
	floatModQInPlace(e.buffer.fpInv.Coeffs, e.q)
	convertFourierPolyToPolyAssign(e.buffer.fpInv.Coeffs, pOut.Coeffs)
}

// ToPolyAddAssign transforms FourierPoly to Poly and adds it to pOut.
func (e *Evaluator[T]) ToPolyAddAssign(fp FourierPoly, pOut Poly[T]) {
	e.buffer.fpInv.CopyFrom(fp)
	ifftInPlace(e.buffer.fpInv.Coeffs, e.twInv)
	floatModQInPlace(e.buffer.fpInv.Coeffs, e.q)
	convertFourierPolyToPolyAddAssign(e.buffer.fpInv.Coeffs, pOut.Coeffs)
}

// ToPolySubAssign transforms FourierPoly to Poly and subtracts it from pOut.
func (e *Evaluator[T]) ToPolySubAssign(fp FourierPoly, pOut Poly[T]) {
	e.buffer.fpInv.CopyFrom(fp)
	ifftInPlace(e.buffer.fpInv.Coeffs, e.twInv)
	floatModQInPlace(e.buffer.fpInv.Coeffs, e.q)
	convertFourierPolyToPolySubAssign(e.buffer.fpInv.Coeffs, pOut.Coeffs)
}

// ToPolyAssignUnsafe transforms FourierPoly to Poly and writes it to pOut.
//
// This method is slightly faster than [*Evaluator.ToPolyAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator[T]) ToPolyAssignUnsafe(fp FourierPoly, pOut Poly[T]) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolyAssign(fp.Coeffs, pOut.Coeffs)
}

// ToPolyAddAssignUnsafe transforms FourierPoly to Poly and adds it to pOut.
//
// This method is slightly faster than [*Evaluator.ToPolyAddAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator[T]) ToPolyAddAssignUnsafe(fp FourierPoly, pOut Poly[T]) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolyAddAssign(fp.Coeffs, pOut.Coeffs)
}

// ToPolySubAssignUnsafe transforms FourierPoly to Poly and subtracts it from pOut.
//
// This method is slightly faster than [*Evaluator.ToPolySubAssign], but it modifies fp directly.
// Use it only if you don't need fp after this method (e.g. fp is a buffer).
func (e *Evaluator[T]) ToPolySubAssignUnsafe(fp FourierPoly, pOut Poly[T]) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolySubAssign(fp.Coeffs, pOut.Coeffs)
}
