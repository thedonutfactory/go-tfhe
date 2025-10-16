package poly

import "unsafe"

// AddFourierPoly returns fp0 + fp1.
func (e *Evaluator) AddFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.AddFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// AddFourierPolyAssign computes fpOut = fp0 + fp1.
func (e *Evaluator) AddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	addCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// SubFourierPoly returns fp0 - fp1.
func (e *Evaluator) SubFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.SubFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// SubFourierPolyAssign computes fpOut = fp0 - fp1.
func (e *Evaluator) SubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	subCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulFourierPoly returns fp0 * fp1.
func (e *Evaluator) MulFourierPoly(fp0, fp1 FourierPoly) FourierPoly {
	fpOut := e.NewFourierPoly()
	e.MulFourierPolyAssign(fp0, fp1, fpOut)
	return fpOut
}

// MulFourierPolyAssign computes fpOut = fp0 * fp1.
// This is element-wise complex multiplication in the frequency domain.
func (e *Evaluator) MulFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulAddFourierPolyAssign computes fpOut += fp0 * fp1.
func (e *Evaluator) MulAddFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulAddCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// MulSubFourierPolyAssign computes fpOut -= fp0 * fp1.
func (e *Evaluator) MulSubFourierPolyAssign(fp0, fp1, fpOut FourierPoly) {
	elementWiseMulSubCmplxAssign(fp0.Coeffs, fp1.Coeffs, fpOut.Coeffs)
}

// FloatMulFourierPolyAssign computes fpOut = c * fp0.
func (e *Evaluator) FloatMulFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	floatMulCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// FloatMulAddFourierPolyAssign computes fpOut += c * fp0.
func (e *Evaluator) FloatMulAddFourierPolyAssign(fp0 FourierPoly, c float64, fpOut FourierPoly) {
	floatMulAddCmplxAssign(fp0.Coeffs, c, fpOut.Coeffs)
}

// addCmplxAssign computes vOut = v0 + v1.
func addCmplxAssign(v0, v1, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] = w0[0] + w1[0]
		wOut[1] = w0[1] + w1[1]
		wOut[2] = w0[2] + w1[2]
		wOut[3] = w0[3] + w1[3]

		wOut[4] = w0[4] + w1[4]
		wOut[5] = w0[5] + w1[5]
		wOut[6] = w0[6] + w1[6]
		wOut[7] = w0[7] + w1[7]
	}
}

// subCmplxAssign computes vOut = v0 - v1.
func subCmplxAssign(v0, v1, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] = w0[0] - w1[0]
		wOut[1] = w0[1] - w1[1]
		wOut[2] = w0[2] - w1[2]
		wOut[3] = w0[3] - w1[3]

		wOut[4] = w0[4] - w1[4]
		wOut[5] = w0[5] - w1[5]
		wOut[6] = w0[6] - w1[6]
		wOut[7] = w0[7] - w1[7]
	}
}

// floatMulCmplxAssign computes vOut = c * v0.
func floatMulCmplxAssign(v0 []float64, c float64, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] = c * w0[0]
		wOut[1] = c * w0[1]
		wOut[2] = c * w0[2]
		wOut[3] = c * w0[3]

		wOut[4] = c * w0[4]
		wOut[5] = c * w0[5]
		wOut[6] = c * w0[6]
		wOut[7] = c * w0[7]
	}
}

// floatMulAddCmplxAssign computes vOut += c * v0.
func floatMulAddCmplxAssign(v0 []float64, c float64, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] += c * w0[0]
		wOut[1] += c * w0[1]
		wOut[2] += c * w0[2]
		wOut[3] += c * w0[3]

		wOut[4] += c * w0[4]
		wOut[5] += c * w0[5]
		wOut[6] += c * w0[6]
		wOut[7] += c * w0[7]
	}
}

// elementWiseMulCmplxAssign computes vOut = v0 * v1 (element-wise complex multiplication).
// This is the key operation for polynomial multiplication in the frequency domain.
func elementWiseMulCmplxAssign(v0, v1, vOut []float64) {
	var vOutR, vOutI float64

	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		// Complex multiplication: (a + bi)(c + di) = (ac - bd) + (ad + bc)i
		// Real part stored in first 4 floats, imaginary in last 4
		vOutR = w0[0]*w1[0] - w0[4]*w1[4]
		vOutI = w0[0]*w1[4] + w0[4]*w1[0]
		wOut[0], wOut[4] = vOutR, vOutI

		vOutR = w0[1]*w1[1] - w0[5]*w1[5]
		vOutI = w0[1]*w1[5] + w0[5]*w1[1]
		wOut[1], wOut[5] = vOutR, vOutI

		vOutR = w0[2]*w1[2] - w0[6]*w1[6]
		vOutI = w0[2]*w1[6] + w0[6]*w1[2]
		wOut[2], wOut[6] = vOutR, vOutI

		vOutR = w0[3]*w1[3] - w0[7]*w1[7]
		vOutI = w0[3]*w1[7] + w0[7]*w1[3]
		wOut[3], wOut[7] = vOutR, vOutI
	}
}

// elementWiseMulAddCmplxAssign computes vOut += v0 * v1.
func elementWiseMulAddCmplxAssign(v0, v1, vOut []float64) {
	var vOutR, vOutI float64

	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		vOutR = wOut[0] + (w0[0]*w1[0] - w0[4]*w1[4])
		vOutI = wOut[4] + (w0[0]*w1[4] + w0[4]*w1[0])
		wOut[0], wOut[4] = vOutR, vOutI

		vOutR = wOut[1] + (w0[1]*w1[1] - w0[5]*w1[5])
		vOutI = wOut[5] + (w0[1]*w1[5] + w0[5]*w1[1])
		wOut[1], wOut[5] = vOutR, vOutI

		vOutR = wOut[2] + (w0[2]*w1[2] - w0[6]*w1[6])
		vOutI = wOut[6] + (w0[2]*w1[6] + w0[6]*w1[2])
		wOut[2], wOut[6] = vOutR, vOutI

		vOutR = wOut[3] + (w0[3]*w1[3] - w0[7]*w1[7])
		vOutI = wOut[7] + (w0[3]*w1[7] + w0[7]*w1[3])
		wOut[3], wOut[7] = vOutR, vOutI
	}
}

// elementWiseMulSubCmplxAssign computes vOut -= v0 * v1.
func elementWiseMulSubCmplxAssign(v0, v1, vOut []float64) {
	var vOutR, vOutI float64

	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		vOutR = wOut[0] - (w0[0]*w1[0] - w0[4]*w1[4])
		vOutI = wOut[4] - (w0[0]*w1[4] + w0[4]*w1[0])
		wOut[0], wOut[4] = vOutR, vOutI

		vOutR = wOut[1] - (w0[1]*w1[1] - w0[5]*w1[5])
		vOutI = wOut[5] - (w0[1]*w1[5] + w0[5]*w1[1])
		wOut[1], wOut[5] = vOutR, vOutI

		vOutR = wOut[2] - (w0[2]*w1[2] - w0[6]*w1[6])
		vOutI = wOut[6] - (w0[2]*w1[6] + w0[6]*w1[2])
		wOut[2], wOut[6] = vOutR, vOutI

		vOutR = wOut[3] - (w0[3]*w1[3] - w0[7]*w1[7])
		vOutI = wOut[7] - (w0[3]*w1[7] + w0[7]*w1[3])
		wOut[3], wOut[7] = vOutR, vOutI
	}
}
