package poly

import "github.com/thedonutfactory/go-tfhe/math/num"

// MulPoly returns p0 * p1.
func (e *Evaluator[T]) MulPoly(p0, p1 Poly[T]) Poly[T] {
	pOut := e.NewPoly()
	e.MulPolyAssign(p0, p1, pOut)
	return pOut
}

// MulPolyAssign computes pOut = p0 * p1.
func (e *Evaluator[T]) MulPolyAssign(p0, p1, pOut Poly[T]) {
	splitBits, splitCount := splitParameters[T](e.degree)

	if splitCount == 1 {
		fp0 := e.ToFourierPoly(p0)
		fp1 := e.ToFourierPoly(p1)
		e.MulFourierPolyAssign(fp0, fp1, fp0)
		e.ToPolyAssignUnsafe(fp0, pOut)
		return
	}

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolyAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// MulAddPolyAssign computes pOut += p0 * p1.
func (e *Evaluator[T]) MulAddPolyAssign(p0, p1, pOut Poly[T]) {
	splitBits, splitCount := splitParameters[T](e.degree)

	if splitCount == 1 {
		fp0 := e.ToFourierPoly(p0)
		fp1 := e.ToFourierPoly(p1)
		e.MulFourierPolyAssign(fp0, fp1, fp0)
		e.ToPolyAddAssignUnsafe(fp0, pOut)
		return
	}

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolyAddAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// MulSubPolyAssign computes pOut -= p0 * p1.
func (e *Evaluator[T]) MulSubPolyAssign(p0, p1, pOut Poly[T]) {
	splitBits, splitCount := splitParameters[T](e.degree)

	if splitCount == 1 {
		fp0 := e.ToFourierPoly(p0)
		fp1 := e.ToFourierPoly(p1)
		e.MulFourierPolyAssign(fp0, fp1, fp0)
		e.ToPolySubAssignUnsafe(fp0, pOut)
		return
	}

	fp0Split := make([]FourierPoly, splitCount)
	fp1Split := make([]FourierPoly, splitCount)
	fpOutSplit := make([]FourierPoly, splitCount)
	for i := 0; i < splitCount; i++ {
		fp0Split[i] = e.NewFourierPoly()
		fp1Split[i] = e.NewFourierPoly()
		fpOutSplit[i] = e.NewFourierPoly()
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp0Split[i])

			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p1.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, fp1Split[i])
		}
	}

	for i := 0; i < splitCount; i++ {
		for j := 0; j < splitCount-i; j++ {
			e.MulAddFourierPolyAssign(fp0Split[i], fp1Split[j], fpOutSplit[i+j])
		}
	}

	e.ToPolySubAssignUnsafe(fpOutSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(fpOutSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] -= e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulPoly returns p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulPoly], and the result is exact unlike [*Evaluator.FourierPolyMulPoly].
func (e *Evaluator[T]) ShortFourierPolyMulPoly(p0 Poly[T], fpShort FourierPoly) Poly[T] {
	pOut := e.NewPoly()
	e.ShortFourierPolyMulPolyAssign(p0, fpShort, pOut)
	return pOut
}

// ShortFourierPolyMulPolyAssign computes pOut = p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulPolyAssign].
func (e *Evaluator[T]) ShortFourierPolyMulPolyAssign(p0 Poly[T], fpShort FourierPoly, pOut Poly[T]) {
	splitBits, splitCount := splitParametersShort[T](e.degree)

	if splitCount == 1 {
		e.ToFourierPolyAssign(p0, e.buffer.fpShortSplit[0])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[0], fpShort, e.buffer.fpShortSplit[0])
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
		return
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	}

	e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulAddPolyAssign computes pOut += p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulAddPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulAddPolyAssign].
func (e *Evaluator[T]) ShortFourierPolyMulAddPolyAssign(p0 Poly[T], fpShort FourierPoly, pOut Poly[T]) {
	splitBits, splitCount := splitParametersShort[T](e.degree)

	if splitCount == 1 {
		e.ToFourierPolyAssign(p0, e.buffer.fpShortSplit[0])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[0], fpShort, e.buffer.fpShortSplit[0])
		e.ToPolyAddAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
		return
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	}

	e.ToPolyAddAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] += e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}

// ShortFourierPolyMulSubPolyAssign computes pOut -= p0 * fpShort, under the assumption that fpShort is a short polynomial.
// (i.e., all coefficients are bounded by [ShortLogBound] bits.)
// This is faster than [*Evaluator.MulSubPolyAssign], and the result is exact unlike [*Evaluator.FourierPolyMulSubPolyAssign].
func (e *Evaluator[T]) ShortFourierPolyMulSubPolyAssign(p0 Poly[T], fpShort FourierPoly, pOut Poly[T]) {
	splitBits, splitCount := splitParametersShort[T](e.degree)

	if splitCount == 1 {
		e.ToFourierPolyAssign(p0, e.buffer.fpShortSplit[0])
		e.MulFourierPolyAssign(e.buffer.fpShortSplit[0], fpShort, e.buffer.fpShortSplit[0])
		e.ToPolySubAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
		return
	}

	if num.IsSigned[T]() {
		var splitChunk T = 1 << splitBits
		for i := 0; i < splitCount; i++ {
			var splitLow T = 1 << (i * int(splitBits))
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] / splitLow) % splitChunk
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	} else {
		var splitMask T = 1<<splitBits - 1
		for i := 0; i < splitCount; i++ {
			splitLowBits := i * int(splitBits)
			for j := 0; j < e.degree; j++ {
				e.buffer.pSplit.Coeffs[j] = (p0.Coeffs[j] >> splitLowBits) & splitMask
			}
			e.ToFourierPolyAssign(e.buffer.pSplit, e.buffer.fpShortSplit[i])
			e.MulFourierPolyAssign(e.buffer.fpShortSplit[i], fpShort, e.buffer.fpShortSplit[i])
		}
	}

	e.ToPolySubAssignUnsafe(e.buffer.fpShortSplit[0], pOut)
	for i := 1; i < splitCount; i++ {
		e.ToPolyAssignUnsafe(e.buffer.fpShortSplit[i], e.buffer.pSplit)
		splitLowBits := i * int(splitBits)
		for j := 0; j < e.degree; j++ {
			pOut.Coeffs[j] -= e.buffer.pSplit.Coeffs[j] << splitLowBits
		}
	}
}
