package poly

import (
	"math/bits"

	"github.com/thedonutfactory/go-tfhe/math/num"
)

// AddFourierPoly returns fp0 + fp1.
func (e *Evaluator[T]) AddFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.AddFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// AddFourierPolyAssign computes fpOut = fp0 + fp1.
func (e *Evaluator[T]) AddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	addCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// SubFourierPoly returns fp0 - fp1.
func (e *Evaluator[T]) SubFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.SubFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// SubFourierPolyAssign computes fpOut = fp0 - fp1.
func (e *Evaluator[T]) SubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	subCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// NegFourierPoly returns -fp0.
func (e *Evaluator[T]) NegFourierPoly(fp0 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.NegFourierPolyAssign(fp0, fpOut)
	return fpOut
}

// NegFourierPolyAssign computes fpOut = -fp0.
func (e *Evaluator[T]) NegFourierPolyAssign(fp0, fpOut FourierPoly) {
	negCmplxAssign(fp0.Coeffs, fpOut.Coeffs)
}

// FloatMulFourierPoly returns c * fp0.
func (e *Evaluator[T]) FloatMulFourierPoly(fp0 FourierPoly, c float64) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.FloatMulFourierPolyAssign(fp0, c, fpOut)
	return fpOut
}

// FloatMulFourierPolyAssign computes fpOut = c * fp0.
func (e *Evaluator[T]) FloatMulFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	floatMulCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// FloatMulAddFourierPolyAssign computes fpOut += c * fp0.
func (e *Evaluator[T]) FloatMulAddFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	floatMulAddCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// FloatMulSubFourierPolyAssign computes fpOut -= c * fp0.
func (e *Evaluator[T]) FloatMulSubFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	floatMulSubCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// CmplxMulFourierPoly returns c * fp0.
func (e *Evaluator[T]) CmplxMulFourierPoly(fp0 FourierPoly, c complex128) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.CmplxMulFourierPolyAssign(fp0, c, fpOut)
	return fpOut
}

// CmplxMulFourierPolyAssign computes fpOut = c * fp0.
func (e *Evaluator[T]) CmplxMulFourierPolyAssign(fp0 FourierPoly, c complex128, fpOut FourierPoly) {
	cmplxMulCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// CmplxMulAddFourierPolyAssign computes fpOut += c * fp0.
func (e *Evaluator[T]) CmplxMulAddFourierPolyAssign(fp0 FourierPoly, c complex128, fpOut FourierPoly) {
	cmplxMulAddCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// CmplxMulSubFourierPolyAssign computes fpOut -= c * fp0.
func (e *Evaluator[T]) CmplxMulSubFourierPolyAssign(fp0 FourierPoly, c complex128, fpOut FourierPoly) {
	cmplxMulSubCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// MulFourierPoly returns fp0 * fp1.
func (e *Evaluator[T]) MulFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.MulFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// MulFourierPolyAssign computes fpOut = fp0 * fp1.
func (e *Evaluator[T]) MulFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulAddFourierPolyAssign computes fpOut += fp0 * fp1.
func (e *Evaluator[T]) MulAddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulAddCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulSubFourierPolyAssign computes fpOut -= fp0 * fp1.
func (e *Evaluator[T]) MulSubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulSubCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// PolyMulFourierPoly returns p * fp0 as FourierPoly.
func (e *Evaluator[T]) PolyMulFourierPoly(fp0 FourierPoly, p Poly[T]) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.PolyMulFourierPolyAssign(fp0, p, fpOut)
	return fpOut
}

// PolyMulFourierPolyAssign computes fpOut = p * fp0.
func (e *Evaluator[T]) PolyMulFourierPolyAssign(fp0 FourierPoly, p Poly[T], fpOut FourierPoly) {
	e.ToFourierPolyAssign(p, e.buffer.fp)

	elementWiseMulCmplxAssign(fp0.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// PolyMulAddFourierPolyAssign computes fpOut += p * fp0.
func (e *Evaluator[T]) PolyMulAddFourierPolyAssign(fp0 FourierPoly, p Poly[T], fpOut FourierPoly) {
	e.ToFourierPolyAssign(p, e.buffer.fp)

	elementWiseMulAddCmplxAssign(fp0.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// PolyMulSubFourierPolyAssign computes fpOut -= p * fp0.
func (e *Evaluator[T]) PolyMulSubFourierPolyAssign(fp0 FourierPoly, p Poly[T], fpOut FourierPoly) {
	e.ToFourierPolyAssign(p, e.buffer.fp)

	elementWiseMulSubCmplxAssign(fp0.Coeffs, e.buffer.fp.Coeffs, fpOut.Coeffs)
}

// PermuteFourierPoly returns fp0(X^d).
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteFourierPoly(fp0 FourierPoly, d int) FourierPoly {
	if d&1 == 0 {
		panic("d not odd")
	}

	fpOut := e.NewFourierPoly()
	e.PermuteFourierPolyAssign(fp0, d, fpOut)
	return fpOut
}

// PermuteFourierPolyAssign computes fpOut = fp0(X^d).
//
// fp0 and fpOut should not overlap. For inplace permutation,
// use [*Evaluator.PermuteFourierPolyInPlace].
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteFourierPolyAssign(fp0 FourierPoly, d int, fpOut FourierPoly) {
	if d&1 == 0 {
		panic("d not odd")
	}

	revShiftBits := 64 - (num.Log2(e.degree) - 1)

	d = d & (2*e.degree - 1)
	if d%4 == 1 {
		k := (d - 1) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) | (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (d*ci - k) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) | (cj & 3)

				fpOut.Coeffs[i+0] = fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] = fp0.Coeffs[j+4]
			}
		}
	} else {
		k := (d - 3) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) | (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (-d*ci + k + 1) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) | (cj & 3)

				fpOut.Coeffs[i+0] = fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] = -fp0.Coeffs[j+4]
			}
		}
	}
}

// PermuteFourierPolyInPlace computes fp0 = fp0(X^d).
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteFourierPolyInPlace(fp0 FourierPoly, d int) {
	if d&1 == 0 {
		panic("d not odd")
	}

	e.PermuteFourierPolyAssign(fp0, d, e.buffer.fpOut)
	fp0.CopyFrom(e.buffer.fpOut)
}

// PermuteAddFourierPolyAssign computes fpOut += fp0(X^d).
//
// fp0 and fpOut should not overlap.
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteAddFourierPolyAssign(fp0 FourierPoly, d int, fpOut FourierPoly) {
	if d&1 == 0 {
		panic("d not odd")
	}

	d = d & (2*e.degree - 1)
	revShiftBits := 64 - (num.Log2(e.degree) - 1)
	if d%4 == 1 {
		k := (d - 1) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) + (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (d*ci - k) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) + (cj & 3)

				fpOut.Coeffs[i+0] += fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] += fp0.Coeffs[j+4]
			}
		}
	} else {
		k := (d - 3) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) + (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (-d*ci + k + 1) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) + (cj & 3)

				fpOut.Coeffs[i+0] += fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] += -fp0.Coeffs[j+4]
			}
		}
	}
}

// PermuteSubFourierPolyAssign computes fpOut -= fp0(X^d).
//
// fp0 and fpOut should not overlap.
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteSubFourierPolyAssign(fp0 FourierPoly, d int, fpOut FourierPoly) {
	if d&1 == 0 {
		panic("d not odd")
	}

	d = d & (2*e.degree - 1)
	revShiftBits := 64 - (num.Log2(e.degree) - 1)
	if d%4 == 1 {
		k := (d - 1) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) | (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (d*ci - k) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) | (cj & 3)

				fpOut.Coeffs[i+0] -= fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] -= fp0.Coeffs[j+4]
			}
		}
	} else {
		k := (d - 3) >> 2
		var ci, cj int
		for ii := 0; ii < e.degree; ii += 8 {
			for i := ii; i < ii+4; i++ {
				ci = ((i >> 3) << 2) | (i & 3)
				ci = int(bits.Reverse64(uint64(ci)) >> revShiftBits)
				cj = (-d*ci + k + 1) & (e.degree>>1 - 1)
				cj = int(bits.Reverse64(uint64(cj)) >> revShiftBits)
				j := ((cj >> 2) << 3) | (cj & 3)

				fpOut.Coeffs[i+0] -= fp0.Coeffs[j+0]
				fpOut.Coeffs[i+4] -= -fp0.Coeffs[j+4]
			}
		}
	}
}
