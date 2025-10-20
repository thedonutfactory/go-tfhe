package poly

import (
	"math"
	"math/cmplx"
	"math/rand"
	"testing"

	"github.com/thedonutfactory/go-tfhe/math/vec"
)

func fftInPlaceRef(coeffs, tw []complex128) {
	N := len(coeffs)

	t := N
	for m := 1; m <= N/2; m <<= 1 {
		t >>= 1
		for i := 0; i < m; i++ {
			j1 := i * t << 1
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]*tw[i]
				coeffs[j], coeffs[j+t] = U+V, U-V
			}
		}
	}
}

func invFFTInPlaceRef(coeffs, twInv []complex128) {
	N := len(coeffs)

	t := 1
	for m := N / 2; m >= 1; m >>= 1 {
		for i := 0; i < m; i++ {
			j1 := i * t << 1
			j2 := j1 + t
			for j := j1; j < j2; j++ {
				U, V := coeffs[j], coeffs[j+t]
				coeffs[j], coeffs[j+t] = U+V, (U-V)*twInv[i]
			}
		}
		t <<= 1
	}
}

func TestFFTAssembly(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	N := 64
	eps := 1e-10

	coeffs := make([]complex128, N)
	for i := 0; i < N; i++ {
		coeffs[i] = complex(r.Float64(), r.Float64())
	}
	coeffsAVX2 := vec.CmplxToFloat4(coeffs)
	coeffsAVX2Out := make([]complex128, N)

	twRef := make([]complex128, N/2)
	twInvRef := make([]complex128, N/2)
	for i := 0; i < N/2; i++ {
		e := -2 * math.Pi * float64(i) / float64(N)
		twRef[i] = cmplx.Exp(complex(0, e))
		twInvRef[i] = cmplx.Exp(-complex(0, e))
	}
	vec.BitReverseInPlace(twRef)
	vec.BitReverseInPlace(twInvRef)

	twist := make([]complex128, N)
	twistInv := make([]complex128, N)
	for i := 0; i < N; i++ {
		e := 2 * math.Pi * float64(i) / float64(4*N)
		twist[i] = cmplx.Exp(complex(0, e))
		twistInv[i] = cmplx.Exp(-complex(0, e)) / complex(float64(N), 0)
	}

	tw, twInv := genTwiddleFactors(N)

	t.Run("FFT", func(t *testing.T) {
		vec.CmplxToFloat4Assign(coeffs, coeffsAVX2)
		fftInPlace(coeffsAVX2, tw)
		vec.Float4ToCmplxAssign(coeffsAVX2, coeffsAVX2Out)

		vec.ElementWiseMulAssign(coeffs, twist, coeffs)
		fftInPlaceRef(coeffs, twRef)

		for i := 0; i < N; i++ {
			if cmplx.Abs(coeffs[i]-coeffsAVX2Out[i]) > eps {
				t.Fatalf("FFT: %v != %v", coeffs[i], coeffsAVX2Out[i])
			}
		}
	})

	t.Run("InvFFT", func(t *testing.T) {
		vec.CmplxToFloat4Assign(coeffs, coeffsAVX2)
		ifftInPlace(coeffsAVX2, twInv)
		vec.Float4ToCmplxAssign(coeffsAVX2, coeffsAVX2Out)

		invFFTInPlaceRef(coeffs, twInvRef)
		vec.ElementWiseMulAssign(coeffs, twistInv, coeffs)

		for i := 0; i < N; i++ {
			if cmplx.Abs(coeffs[i]-coeffsAVX2Out[i]) > eps {
				t.Fatalf("InvFFT: %v != %v", coeffs[i], coeffsAVX2Out[i])
			}
		}
	})
}
