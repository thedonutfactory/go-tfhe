//go:build !(amd64 && !purego)

package poly

import "unsafe"

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

// negCmplxAssign computes vOut = -v0.
func negCmplxAssign(v0, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] = -w0[0]
		wOut[1] = -w0[1]
		wOut[2] = -w0[2]
		wOut[3] = -w0[3]

		wOut[4] = -w0[4]
		wOut[5] = -w0[5]
		wOut[6] = -w0[6]
		wOut[7] = -w0[7]
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

// floatMulSubCmplxAssign computes vOut -= c * v0.
func floatMulSubCmplxAssign(v0 []float64, c float64, vOut []float64) {
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] -= c * w0[0]
		wOut[1] -= c * w0[1]
		wOut[2] -= c * w0[2]
		wOut[3] -= c * w0[3]

		wOut[4] -= c * w0[4]
		wOut[5] -= c * w0[5]
		wOut[6] -= c * w0[6]
		wOut[7] -= c * w0[7]
	}
}

// cmplxMulCmplxAssign computes vOut = c * v0.
func cmplxMulCmplxAssign(v0 []float64, c complex128, vOut []float64) {
	cR, cI := real(c), imag(c)
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] = w0[0]*cR - w0[4]*cI
		wOut[1] = w0[1]*cR - w0[5]*cI
		wOut[2] = w0[2]*cR - w0[6]*cI
		wOut[3] = w0[3]*cR - w0[7]*cI

		wOut[4] = w0[0]*cI + w0[4]*cR
		wOut[5] = w0[1]*cI + w0[5]*cR
		wOut[6] = w0[2]*cI + w0[6]*cR
		wOut[7] = w0[3]*cI + w0[7]*cR
	}
}

// cmplxMulAddCmplxAssign computes vOut += c * v0.
func cmplxMulAddCmplxAssign(v0 []float64, c complex128, vOut []float64) {
	cR, cI := real(c), imag(c)
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] += w0[0]*cR - w0[4]*cI
		wOut[1] += w0[1]*cR - w0[5]*cI
		wOut[2] += w0[2]*cR - w0[6]*cI
		wOut[3] += w0[3]*cR - w0[7]*cI

		wOut[4] += w0[0]*cI + w0[4]*cR
		wOut[5] += w0[1]*cI + w0[5]*cR
		wOut[6] += w0[2]*cI + w0[6]*cR
		wOut[7] += w0[3]*cI + w0[7]*cR
	}
}

// cmplxMulSubCmplxAssign computes vOut -= c * v0.
func cmplxMulSubCmplxAssign(v0 []float64, c complex128, vOut []float64) {
	cR, cI := real(c), imag(c)
	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

		wOut[0] -= w0[0]*cR - w0[4]*cI
		wOut[1] -= w0[1]*cR - w0[5]*cI
		wOut[2] -= w0[2]*cR - w0[6]*cI
		wOut[3] -= w0[3]*cR - w0[7]*cI

		wOut[4] -= w0[0]*cI + w0[4]*cR
		wOut[5] -= w0[1]*cI + w0[5]*cR
		wOut[6] -= w0[2]*cI + w0[6]*cR
		wOut[7] -= w0[3]*cI + w0[7]*cR
	}
}

// elementWiseMulCmplxAssign computes vOut = v0 * v1.
func elementWiseMulCmplxAssign(v0, v1, vOut []float64) {
	var vOutR, vOutI float64

	for i := 0; i < len(vOut); i += 8 {
		w0 := (*[8]float64)(unsafe.Pointer(&v0[i]))
		w1 := (*[8]float64)(unsafe.Pointer(&v1[i]))
		wOut := (*[8]float64)(unsafe.Pointer(&vOut[i]))

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
