//go:build !(amd64 && !purego)

package poly

import (
	"math"
	"unsafe"

	"github.com/thedonutfactory/go-tfhe/math/num"
)

// convertPolyToFourierPolyAssign converts and folds p to fpOut.
func convertPolyToFourierPolyAssign[T num.Integer](p []T, fpOut []float64) {
	N := len(p)

	var z T
	switch any(z).(type) {
	case uint:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uint)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uint)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int(q0[0]))
			fqOut[1] = float64(int(q0[1]))
			fqOut[2] = float64(int(q0[2]))
			fqOut[3] = float64(int(q0[3]))

			fqOut[4] = float64(int(q1[0]))
			fqOut[5] = float64(int(q1[1]))
			fqOut[6] = float64(int(q1[2]))
			fqOut[7] = float64(int(q1[3]))
		}
	case uintptr:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uintptr)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uintptr)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int(q0[0]))
			fqOut[1] = float64(int(q0[1]))
			fqOut[2] = float64(int(q0[2]))
			fqOut[3] = float64(int(q0[3]))

			fqOut[4] = float64(int(q1[0]))
			fqOut[5] = float64(int(q1[1]))
			fqOut[6] = float64(int(q1[2]))
			fqOut[7] = float64(int(q1[3]))
		}
	case uint8:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uint8)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uint8)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int8(q0[0]))
			fqOut[1] = float64(int8(q0[1]))
			fqOut[2] = float64(int8(q0[2]))
			fqOut[3] = float64(int8(q0[3]))

			fqOut[4] = float64(int8(q1[0]))
			fqOut[5] = float64(int8(q1[1]))
			fqOut[6] = float64(int8(q1[2]))
			fqOut[7] = float64(int8(q1[3]))
		}
	case uint16:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uint16)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uint16)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int16(q0[0]))
			fqOut[1] = float64(int16(q0[1]))
			fqOut[2] = float64(int16(q0[2]))
			fqOut[3] = float64(int16(q0[3]))

			fqOut[4] = float64(int16(q1[0]))
			fqOut[5] = float64(int16(q1[1]))
			fqOut[6] = float64(int16(q1[2]))
			fqOut[7] = float64(int16(q1[3]))
		}
	case uint32:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uint32)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uint32)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int32(q0[0]))
			fqOut[1] = float64(int32(q0[1]))
			fqOut[2] = float64(int32(q0[2]))
			fqOut[3] = float64(int32(q0[3]))

			fqOut[4] = float64(int32(q1[0]))
			fqOut[5] = float64(int32(q1[1]))
			fqOut[6] = float64(int32(q1[2]))
			fqOut[7] = float64(int32(q1[3]))
		}
	case uint64:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			q0 := (*[4]uint64)(unsafe.Pointer(&p[ii]))
			q1 := (*[4]uint64)(unsafe.Pointer(&p[ii+N/2]))
			fqOut := (*[8]float64)(unsafe.Pointer(&fpOut[i]))

			fqOut[0] = float64(int64(q0[0]))
			fqOut[1] = float64(int64(q0[1]))
			fqOut[2] = float64(int64(q0[2]))
			fqOut[3] = float64(int64(q0[3]))

			fqOut[4] = float64(int64(q1[0]))
			fqOut[5] = float64(int64(q1[1]))
			fqOut[6] = float64(int64(q1[2]))
			fqOut[7] = float64(int64(q1[3]))
		}
	default:
		for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
			fpOut[i+0] = float64(p[ii+0])
			fpOut[i+1] = float64(p[ii+1])
			fpOut[i+2] = float64(p[ii+2])
			fpOut[i+3] = float64(p[ii+3])

			fpOut[i+4] = float64(p[ii+0+N/2])
			fpOut[i+5] = float64(p[ii+1+N/2])
			fpOut[i+6] = float64(p[ii+2+N/2])
			fpOut[i+7] = float64(p[ii+3+N/2])
		}
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
func convertFourierPolyToPolyAssign[T num.Integer](fp []float64, pOut []T) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]T)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]T)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] = T(int64(fq[0]))
		qOut0[1] = T(int64(fq[1]))
		qOut0[2] = T(int64(fq[2]))
		qOut0[3] = T(int64(fq[3]))

		qOut1[0] = T(int64(fq[4]))
		qOut1[1] = T(int64(fq[5]))
		qOut1[2] = T(int64(fq[6]))
		qOut1[3] = T(int64(fq[7]))
	}
}

// convertFourierPolyToPolyAddAssign converts and unfolds fp and adds it to pOut.
func convertFourierPolyToPolyAddAssign[T num.Integer](fp []float64, pOut []T) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]T)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]T)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] += T(int64(fq[0]))
		qOut0[1] += T(int64(fq[1]))
		qOut0[2] += T(int64(fq[2]))
		qOut0[3] += T(int64(fq[3]))

		qOut1[0] += T(int64(fq[4]))
		qOut1[1] += T(int64(fq[5]))
		qOut1[2] += T(int64(fq[6]))
		qOut1[3] += T(int64(fq[7]))
	}
}

// convertFourierPolyToPolySubAssign converts and unfolds fp and subtracts it from pOut.
func convertFourierPolyToPolySubAssign[T num.Integer](fp []float64, pOut []T) {
	N := len(fp)

	for i, ii := 0, 0; i < N; i, ii = i+8, ii+4 {
		fq := (*[8]float64)(unsafe.Pointer(&fp[i]))
		qOut0 := (*[4]T)(unsafe.Pointer(&pOut[ii]))
		qOut1 := (*[4]T)(unsafe.Pointer(&pOut[ii+N/2]))

		qOut0[0] -= T(int64(fq[0]))
		qOut0[1] -= T(int64(fq[1]))
		qOut0[2] -= T(int64(fq[2]))
		qOut0[3] -= T(int64(fq[3]))

		qOut1[0] -= T(int64(fq[4]))
		qOut1[1] -= T(int64(fq[5]))
		qOut1[2] -= T(int64(fq[6]))
		qOut1[3] -= T(int64(fq[7]))
	}
}
