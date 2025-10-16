package poly

import (
	"math"
	"unsafe"

	"github.com/thedonutfactory/go-tfhe/params"
)

// ToFourierPoly transforms Poly to FourierPoly.
func (e *Evaluator) ToFourierPoly(p Poly) FourierPoly {
	fpOut := NewFourierPoly(e.degree)
	e.ToFourierPolyAssign(p, fpOut)
	return fpOut
}

// ToFourierPolyAssign transforms Poly to FourierPoly and writes it to fpOut.
func (e *Evaluator) ToFourierPolyAssign(p Poly, fpOut FourierPoly) {
	convertPolyToFourierPolyAssign(p.Coeffs, fpOut.Coeffs)
	fftInPlace(fpOut.Coeffs, e.tw)
}

// ToPoly transforms FourierPoly to Poly.
func (e *Evaluator) ToPoly(fp FourierPoly) Poly {
	pOut := NewPoly(e.degree)
	e.ToPolyAssign(fp, pOut)
	return pOut
}

// ToPolyAssign transforms FourierPoly to Poly and writes it to pOut.
func (e *Evaluator) ToPolyAssign(fp FourierPoly, pOut Poly) {
	e.buffer.fpInv.CopyFrom(fp)
	ifftInPlace(e.buffer.fpInv.Coeffs, e.twInv)
	floatModQInPlace(e.buffer.fpInv.Coeffs, e.q)
	convertFourierPolyToPolyAssign(e.buffer.fpInv.Coeffs, pOut.Coeffs)
}

// ToPolyAssignUnsafe transforms FourierPoly to Poly and writes it to pOut.
// This method modifies fp directly, so use it only if you don't need fp after.
func (e *Evaluator) ToPolyAssignUnsafe(fp FourierPoly, pOut Poly) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolyAssign(fp.Coeffs, pOut.Coeffs)
}

// ToPolyAddAssignUnsafe transforms FourierPoly to Poly and adds it to pOut.
// This method modifies fp directly.
func (e *Evaluator) ToPolyAddAssignUnsafe(fp FourierPoly, pOut Poly) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolyAddAssign(fp.Coeffs, pOut.Coeffs)
}

// ToPolySubAssignUnsafe transforms FourierPoly to Poly and subtracts it from pOut.
// This method modifies fp directly.
func (e *Evaluator) ToPolySubAssignUnsafe(fp FourierPoly, pOut Poly) {
	ifftInPlace(fp.Coeffs, e.twInv)
	floatModQInPlace(fp.Coeffs, e.q)
	convertFourierPolyToPolySubAssign(fp.Coeffs, pOut.Coeffs)
}

// convertPolyToFourierPolyAssign converts and folds p to fpOut.
// This splits the polynomial into two halves and interleaves them for SIMD efficiency.
func convertPolyToFourierPolyAssign(p []params.Torus, fpOut []float64) {
	N := len(p)

	// Process 4 elements at a time for SIMD efficiency
	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		q0 := (*[4]params.Torus)(unsafe.Pointer(&p[ii]))
		q1 := (*[4]params.Torus)(unsafe.Pointer(&p[ii+N/2]))
		fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

		// First half (real parts)
		fqOut[0] = float64(int32(q0[0]))
		fqOut[1] = float64(int32(q0[1]))
		fqOut[2] = float64(int32(q0[2]))
		fqOut[3] = float64(int32(q0[3]))

		// Second half (imaginary parts)
		fqOut[4] = float64(int32(q1[0]))
		fqOut[5] = float64(int32(q1[1]))
		fqOut[6] = float64(int32(q1[2]))
		fqOut[7] = float64(int32(q1[3]))
	}
}

// floatModQInPlace computes coeffs mod Q in place.
func floatModQInPlace(coeffs []float64, Q float64) {
	N := len(coeffs)

	for i := 0; i < N; i += 8 {
		c := (*[8]float64)(unsafe.Pointer(&coeffs[i]))

		c[0] = math.Round(c[0] - Q*math.Round(c[0]/Q))
		c[1] = math.Round(c[1] - Q*math.Round(c[1]/Q))
		c[2] = math.Round(c[2] - Q*math.Round(c[2]/Q))
		c[3] = math.Round(c[3] - Q*math.Round(c[3]/Q))

		c[4] = math.Round(c[4] - Q*math.Round(c[4]/Q))
		c[5] = math.Round(c[5] - Q*math.Round(c[5]/Q))
		c[6] = math.Round(c[6] - Q*math.Round(c[6]/Q))
		c[7] = math.Round(c[7] - Q*math.Round(c[7]/Q))
	}
}

// convertFourierPolyToPolyAssign converts and unfolds fp to pOut.
func convertFourierPolyToPolyAssign(fp []float64, pOut []params.Torus) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] = params.Torus(int64(fq[0]))
		qOut0[1] = params.Torus(int64(fq[1]))
		qOut0[2] = params.Torus(int64(fq[2]))
		qOut0[3] = params.Torus(int64(fq[3]))

		qOut1[0] = params.Torus(int64(fq[4]))
		qOut1[1] = params.Torus(int64(fq[5]))
		qOut1[2] = params.Torus(int64(fq[6]))
		qOut1[3] = params.Torus(int64(fq[7]))
	}
}

// convertFourierPolyToPolyAddAssign converts and unfolds fp and adds it to pOut.
func convertFourierPolyToPolyAddAssign(fp []float64, pOut []params.Torus) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] += params.Torus(int64(fq[0]))
		qOut0[1] += params.Torus(int64(fq[1]))
		qOut0[2] += params.Torus(int64(fq[2]))
		qOut0[3] += params.Torus(int64(fq[3]))

		qOut1[0] += params.Torus(int64(fq[4]))
		qOut1[1] += params.Torus(int64(fq[5]))
		qOut1[2] += params.Torus(int64(fq[6]))
		qOut1[3] += params.Torus(int64(fq[7]))
	}
}

// convertFourierPolyToPolySubAssign converts and unfolds fp and subtracts it from pOut.
func convertFourierPolyToPolySubAssign(fp []float64, pOut []params.Torus) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]params.Torus)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] -= params.Torus(int64(fq[0]))
		qOut0[1] -= params.Torus(int64(fq[1]))
		qOut0[2] -= params.Torus(int64(fq[2]))
		qOut0[3] -= params.Torus(int64(fq[3]))

		qOut1[0] -= params.Torus(int64(fq[4]))
		qOut1[1] -= params.Torus(int64(fq[5]))
		qOut1[2] -= params.Torus(int64(fq[6]))
		qOut1[3] -= params.Torus(int64(fq[7]))
	}
}

// butterfly performs FFT butterfly operation.
func butterfly(uR, uI, vR, vI, wR, wI float64) (float64, float64, float64, float64) {
	vwR := vR*wR - vI*wI
	vwI := vR*wI + vI*wR
	return uR + vwR, uI + vwI, uR - vwR, uI - vwI
}

// fftInPlace performs in-place FFT on coeffs using twiddle factors tw.
// This is optimized for SIMD processing of 4 complex numbers at a time.
func fftInPlace(coeffs []float64, tw []complex128) {
	N := len(coeffs)
	wIdx := 0

	// First stage
	wReal := real(tw[wIdx])
	wImag := imag(tw[wIdx])
	wIdx++
	for j := 0; j < N/2; j += 8 {
		u := (*[8]float64)(unsafe.Pointer(&coeffs[j]))
		v := (*[8]float64)(unsafe.Pointer(&coeffs[j+N/2]))

		u[0], u[4], v[0], v[4] = butterfly(u[0], u[4], v[0], v[4], wReal, wImag)
		u[1], u[5], v[1], v[5] = butterfly(u[1], u[5], v[1], v[5], wReal, wImag)
		u[2], u[6], v[2], v[6] = butterfly(u[2], u[6], v[2], v[6], wReal, wImag)
		u[3], u[7], v[3], v[7] = butterfly(u[3], u[7], v[3], v[7], wReal, wImag)
	}

	// Middle stages
	t := N / 2
	for m := 2; m <= N/16; m <<= 1 {
		t >>= 1
		for i := 0; i < m; i++ {
			j1 := 2 * i * t
			j2 := j1 + t

			wReal := real(tw[wIdx])
			wImag := imag(tw[wIdx])
			wIdx++

			for j := j1; j < j2; j += 8 {
				u := (*[8]float64)(unsafe.Pointer(&coeffs[j]))
				v := (*[8]float64)(unsafe.Pointer(&coeffs[j+t]))

				u[0], u[4], v[0], v[4] = butterfly(u[0], u[4], v[0], v[4], wReal, wImag)
				u[1], u[5], v[1], v[5] = butterfly(u[1], u[5], v[1], v[5], wReal, wImag)
				u[2], u[6], v[2], v[6] = butterfly(u[2], u[6], v[2], v[6], wReal, wImag)
				u[3], u[7], v[3], v[7] = butterfly(u[3], u[7], v[3], v[7], wReal, wImag)
			}
		}
	}

	// Second-to-last stage
	for j := 0; j < N; j += 8 {
		wReal := real(tw[wIdx])
		wImag := imag(tw[wIdx])
		wIdx++

		uvReal := (*[4]float64)(unsafe.Pointer(&coeffs[j]))
		uvImag := (*[4]float64)(unsafe.Pointer(&coeffs[j+4]))

		uvReal[0], uvImag[0], uvReal[2], uvImag[2] = butterfly(uvReal[0], uvImag[0], uvReal[2], uvImag[2], wReal, wImag)
		uvReal[1], uvImag[1], uvReal[3], uvImag[3] = butterfly(uvReal[1], uvImag[1], uvReal[3], uvImag[3], wReal, wImag)
	}

	// Last stage
	for j := 0; j < N; j += 8 {
		wReal0 := real(tw[wIdx])
		wImag0 := imag(tw[wIdx])
		wReal1 := real(tw[wIdx+1])
		wImag1 := imag(tw[wIdx+1])
		wIdx += 2

		uvReal := (*[4]float64)(unsafe.Pointer(&coeffs[j]))
		uvImag := (*[4]float64)(unsafe.Pointer(&coeffs[j+4]))

		uvReal[0], uvImag[0], uvReal[1], uvImag[1] = butterfly(uvReal[0], uvImag[0], uvReal[1], uvImag[1], wReal0, wImag0)
		uvReal[2], uvImag[2], uvReal[3], uvImag[3] = butterfly(uvReal[2], uvImag[2], uvReal[3], uvImag[3], wReal1, wImag1)
	}
}

// invButterfly performs inverse FFT butterfly operation.
func invButterfly(uR, uI, vR, vI, wR, wI float64) (float64, float64, float64, float64) {
	uR, uI, vR, vI = uR+vR, uI+vI, uR-vR, uI-vI
	vwR := vR*wR - vI*wI
	vwI := vR*wI + vI*wR
	return uR, uI, vwR, vwI
}

// ifftInPlace performs in-place inverse FFT on coeffs using twiddle factors twInv.
func ifftInPlace(coeffs []float64, twInv []complex128) {
	N := len(coeffs)
	wIdx := 0

	// First stage (reverse of last FFT stage)
	for j := 0; j < N; j += 8 {
		wReal0 := real(twInv[wIdx])
		wImag0 := imag(twInv[wIdx])
		wReal1 := real(twInv[wIdx+1])
		wImag1 := imag(twInv[wIdx+1])
		wIdx += 2

		uvReal := (*[4]float64)(unsafe.Pointer(&coeffs[j]))
		uvImag := (*[4]float64)(unsafe.Pointer(&coeffs[j+4]))

		uvReal[0], uvImag[0], uvReal[1], uvImag[1] = invButterfly(uvReal[0], uvImag[0], uvReal[1], uvImag[1], wReal0, wImag0)
		uvReal[2], uvImag[2], uvReal[3], uvImag[3] = invButterfly(uvReal[2], uvImag[2], uvReal[3], uvImag[3], wReal1, wImag1)
	}

	// Second stage
	for j := 0; j < N; j += 8 {
		wReal := real(twInv[wIdx])
		wImag := imag(twInv[wIdx])
		wIdx++

		uvReal := (*[4]float64)(unsafe.Pointer(&coeffs[j]))
		uvImag := (*[4]float64)(unsafe.Pointer(&coeffs[j+4]))

		uvReal[0], uvImag[0], uvReal[2], uvImag[2] = invButterfly(uvReal[0], uvImag[0], uvReal[2], uvImag[2], wReal, wImag)
		uvReal[1], uvImag[1], uvReal[3], uvImag[3] = invButterfly(uvReal[1], uvImag[1], uvReal[3], uvImag[3], wReal, wImag)
	}

	// Middle stages
	t := 8
	for m := N / 16; m >= 2; m >>= 1 {
		for i := 0; i < m; i++ {
			j1 := 2 * i * t
			j2 := j1 + t

			wReal := real(twInv[wIdx])
			wImag := imag(twInv[wIdx])
			wIdx++

			for j := j1; j < j2; j += 8 {
				u := (*[8]float64)(unsafe.Pointer(&coeffs[j]))
				v := (*[8]float64)(unsafe.Pointer(&coeffs[j+t]))

				u[0], u[4], v[0], v[4] = invButterfly(u[0], u[4], v[0], v[4], wReal, wImag)
				u[1], u[5], v[1], v[5] = invButterfly(u[1], u[5], v[1], v[5], wReal, wImag)
				u[2], u[6], v[2], v[6] = invButterfly(u[2], u[6], v[2], v[6], wReal, wImag)
				u[3], u[7], v[3], v[7] = invButterfly(u[3], u[7], v[3], v[7], wReal, wImag)
			}
		}
		t <<= 1
	}

	// Last stage with scaling
	scale := float64(N / 2)
	wReal := real(twInv[wIdx])
	wImag := imag(twInv[wIdx])
	for j := 0; j < N/2; j += 8 {
		u := (*[8]float64)(unsafe.Pointer(&coeffs[j]))
		v := (*[8]float64)(unsafe.Pointer(&coeffs[j+N/2]))

		u[0], u[4], v[0], v[4] = invButterfly(u[0], u[4], v[0], v[4], wReal, wImag)
		u[1], u[5], v[1], v[5] = invButterfly(u[1], u[5], v[1], v[5], wReal, wImag)
		u[2], u[6], v[2], v[6] = invButterfly(u[2], u[6], v[2], v[6], wReal, wImag)
		u[3], u[7], v[3], v[7] = invButterfly(u[3], u[7], v[3], v[7], wReal, wImag)

		u[0] /= scale
		u[1] /= scale
		u[2] /= scale
		u[3] /= scale

		u[4] /= scale
		u[5] /= scale
		u[6] /= scale
		u[7] /= scale

		v[0] /= scale
		v[1] /= scale
		v[2] /= scale
		v[3] /= scale

		v[4] /= scale
		v[5] /= scale
		v[6] /= scale
		v[7] /= scale
	}
}
