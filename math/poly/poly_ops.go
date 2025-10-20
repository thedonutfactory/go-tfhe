package poly

import (
	"github.com/thedonutfactory/go-tfhe/math/vec"
)

// AddPoly returns p0 + p1.
func (e *Evaluator[T]) AddPoly(p0, p1 Poly[T]) Poly[T] {
	pOut := e.NewPoly()
	e.AddPolyAssign(p0, p1, pOut)
	return pOut
}

// AddPolyAssign computes pOut = p0 + p1.
func (e *Evaluator[T]) AddPolyAssign(p0, p1, pOut Poly[T]) {
	vec.AddAssign(p0.Coeffs, p1.Coeffs, pOut.Coeffs)
}

// SubPoly returns p0 - p1.
func (e *Evaluator[T]) SubPoly(p0, p1 Poly[T]) Poly[T] {
	pOut := e.NewPoly()
	e.SubPolyAssign(p0, p1, pOut)
	return pOut
}

// SubPolyAssign computes pOut = p0 - p1.
func (e *Evaluator[T]) SubPolyAssign(p0, p1, pOut Poly[T]) {
	vec.SubAssign(p0.Coeffs, p1.Coeffs, pOut.Coeffs)
}

// NegPoly returns pOut = -p0.
func (e *Evaluator[T]) NegPoly(p0 Poly[T]) Poly[T] {
	pOut := e.NewPoly()
	e.NegPolyAssign(p0, pOut)
	return pOut
}

// NegPolyAssign computes pOut = -p0.
func (e *Evaluator[T]) NegPolyAssign(p0, pOut Poly[T]) {
	vec.NegAssign(p0.Coeffs, pOut.Coeffs)
}

// ScalarMulPoly returns c * p0.
func (e *Evaluator[T]) ScalarMulPoly(p0 Poly[T], c T) Poly[T] {
	pOut := e.NewPoly()
	e.ScalarMulPolyAssign(p0, c, pOut)
	return pOut
}

// ScalarMulPolyAssign computes pOut = c * p0.
func (e *Evaluator[T]) ScalarMulPolyAssign(p0 Poly[T], c T, pOut Poly[T]) {
	vec.ScalarMulAssign(p0.Coeffs, c, pOut.Coeffs)
}

// ScalarMulAddPolyAssign computes pOut += c * p0.
func (e *Evaluator[T]) ScalarMulAddPolyAssign(p0 Poly[T], c T, pOut Poly[T]) {
	vec.ScalarMulAddAssign(p0.Coeffs, c, pOut.Coeffs)
}

// ScalarMulSubPolyAssign computes pOut -= c * p0.
func (e *Evaluator[T]) ScalarMulSubPolyAssign(p0 Poly[T], c T, pOut Poly[T]) {
	vec.ScalarMulSubAssign(p0.Coeffs, c, pOut.Coeffs)
}

// FourierPolyMulPoly returns p0 * fp.
func (e *Evaluator[T]) FourierPolyMulPoly(p0 Poly[T], fp FourierPoly) Poly[T] {
	pOut := e.NewPoly()
	e.FourierPolyMulPolyAssign(p0, fp, pOut)
	return pOut
}

// FourierPolyMulPolyAssign computes pOut = p0 * fp.
func (e *Evaluator[T]) FourierPolyMulPolyAssign(p0 Poly[T], fp FourierPoly, pOut Poly[T]) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolyAssignUnsafe(e.buffer.fp, pOut)
}

// FourierPolyMulAddPolyAssign computes pOut += p0 * fp.
func (e *Evaluator[T]) FourierPolyMulAddPolyAssign(p0 Poly[T], fp FourierPoly, pOut Poly[T]) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolyAddAssignUnsafe(e.buffer.fp, pOut)
}

// FourierPolyMulSubPolyAssign computes pOut -= p0 * fp.
func (e *Evaluator[T]) FourierPolyMulSubPolyAssign(p0 Poly[T], fp FourierPoly, pOut Poly[T]) {
	e.ToFourierPolyAssign(p0, e.buffer.fp)
	e.MulFourierPolyAssign(e.buffer.fp, fp, e.buffer.fp)
	e.ToPolySubAssignUnsafe(e.buffer.fp, pOut)
}

// MonomialMulPoly returns X^d * p0.
func (e *Evaluator[T]) MonomialMulPoly(p0 Poly[T], d int) Poly[T] {
	pOut := e.NewPoly()
	e.MonomialMulPolyAssign(p0, d, pOut)
	return pOut
}

// MonomialMulPolyAssign computes pOut = X^d * p0.
//
// p0 and pOut should not overlap. For inplace multiplication,
// use [*Evaluator.MonomialMulPolyInPlace].
func (e *Evaluator[T]) MonomialMulPolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	switch k := d & (2*e.degree - 1); {
	case e.degree <= k:
		for i, ii := 0, -k+2*e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = p0.Coeffs[ii]
		}
		for i, ii := k-e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = -p0.Coeffs[ii]
		}
	case 0 <= k && k < e.degree:
		for i, ii := 0, -k+e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = -p0.Coeffs[ii]
		}
		for i, ii := k, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = p0.Coeffs[ii]
		}
	case -e.degree <= k && k < 0:
		for i, ii := 0, -k; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = p0.Coeffs[ii]
		}
		for i, ii := k+e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = -p0.Coeffs[ii]
		}
	case k < -e.degree:
		for i, ii := 0, -k-e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = -p0.Coeffs[ii]
		}
		for i, ii := k+2*e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] = p0.Coeffs[ii]
		}
	}
}

// MonomialMulPolyInPlace computes p0 = X^d * p0.
func (e *Evaluator[T]) MonomialMulPolyInPlace(p0 Poly[T], d int) {
	kk := d & (e.degree - 1)
	vec.RotateInPlace(p0.Coeffs, kk)

	switch k := d & (2*e.degree - 1); {
	case e.degree <= k:
		for i := kk; i < e.degree; i++ {
			p0.Coeffs[i] = -p0.Coeffs[i]
		}
	case 0 <= k && k < e.degree:
		for i := 0; i < kk; i++ {
			p0.Coeffs[i] = -p0.Coeffs[i]
		}
	case -e.degree <= k && k < 0:
		for i := e.degree + kk; i < e.degree; i++ {
			p0.Coeffs[i] = -p0.Coeffs[i]
		}
	case k < -e.degree:
		for i := 0; i < e.degree+kk; i++ {
			p0.Coeffs[i] = -p0.Coeffs[i]
		}
	}
}

// MonomialMulAddPolyAssign computes pOut += X^d * p0.
//
// p0 and pOut should not overlap.
func (e *Evaluator[T]) MonomialMulAddPolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	switch k := d & (2*e.degree - 1); {
	case e.degree <= k:
		for i, ii := 0, -k+2*e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += p0.Coeffs[ii]
		}
		for i, ii := k-e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += -p0.Coeffs[ii]
		}
	case 0 <= k && k < e.degree:
		for i, ii := 0, -k+e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += -p0.Coeffs[ii]
		}
		for i, ii := k, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += p0.Coeffs[ii]
		}
	case -e.degree <= k && k < 0:
		for i, ii := 0, -k; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += p0.Coeffs[ii]
		}
		for i, ii := k+e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += -p0.Coeffs[ii]
		}
	case k < -e.degree:
		for i, ii := 0, -k-e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += -p0.Coeffs[ii]
		}
		for i, ii := k+2*e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] += p0.Coeffs[ii]
		}
	}
}

// MonomialMulSubPolyAssign computes pOut -= X^d * p0.
//
// p0 and pOut should not overlap.
func (e *Evaluator[T]) MonomialMulSubPolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	switch k := d & (2*e.degree - 1); {
	case e.degree <= k:
		for i, ii := 0, -k+2*e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= p0.Coeffs[ii]
		}
		for i, ii := k-e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= -p0.Coeffs[ii]
		}
	case 0 <= k && k < e.degree:
		for i, ii := 0, -k+e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= -p0.Coeffs[ii]
		}
		for i, ii := k, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= p0.Coeffs[ii]
		}
	case -e.degree <= k && k < 0:
		for i, ii := 0, -k; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= p0.Coeffs[ii]
		}
		for i, ii := k+e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= -p0.Coeffs[ii]
		}
	case k < -e.degree:
		for i, ii := 0, -k-e.degree; ii < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= -p0.Coeffs[ii]
		}
		for i, ii := k+2*e.degree, 0; i < e.degree; i, ii = i+1, ii+1 {
			pOut.Coeffs[i] -= p0.Coeffs[ii]
		}
	}
}

// PermutePoly returns p0(X^d).
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermutePoly(p0 Poly[T], d int) Poly[T] {
	if d&1 == 0 {
		panic("d not odd")
	}

	pOut := e.NewPoly()
	e.PermutePolyAssign(p0, d, pOut)
	return pOut
}

// PermutePolyAssign computes pOut = p0(X^d).
//
// p0 and pOut should not overlap. For inplace permutation,
// use [*Evaluator.PermutePolyInPlace].
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermutePolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	if d&1 == 0 {
		panic("d not odd")
	}

	for i := 0; i < e.degree; i++ {
		j := (d * i) & (2*e.degree - 1)
		if j < e.degree {
			pOut.Coeffs[j] = p0.Coeffs[i]
		} else {
			pOut.Coeffs[j-e.degree] = -p0.Coeffs[i]
		}
	}
}

// PermutePolyInPlace computes p0 = p0(X^d).
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermutePolyInPlace(p0 Poly[T], d int) {
	e.PermutePolyAssign(p0, d, e.buffer.pOut)
	p0.CopyFrom(e.buffer.pOut)
}

// PermuteAddPolyAssign computes pOut += p0(X^d).
//
// p0 and pOut should not overlap.
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteAddPolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	if d&1 == 0 {
		panic("d not odd")
	}

	for i := 0; i < e.degree; i++ {
		j := (d * i) & (2*e.degree - 1)
		if j < e.degree {
			pOut.Coeffs[j] += p0.Coeffs[i]
		} else {
			pOut.Coeffs[j-e.degree] -= p0.Coeffs[i]
		}
	}
}

// PermuteSubPolyAssign computes pOut -= p0(X^d).
//
// p0 and pOut should not overlap.
//
// Panics when d is not odd.
// This is because the permutation is not bijective when d is even.
func (e *Evaluator[T]) PermuteSubPolyAssign(p0 Poly[T], d int, pOut Poly[T]) {
	if d&1 == 0 {
		panic("d not odd")
	}

	for i := 0; i < e.degree; i++ {
		j := (d * i) & (2*e.degree - 1)
		if j < e.degree {
			pOut.Coeffs[j] -= p0.Coeffs[i]
		} else {
			pOut.Coeffs[j-e.degree] += p0.Coeffs[i]
		}
	}
}
