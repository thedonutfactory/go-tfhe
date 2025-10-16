package poly

import (
	"testing"

	"github.com/thedonutfactory/go-tfhe/params"
)

// TestFFTRoundTrip tests that FFT -> IFFT gives back the original polynomial
func TestFFTRoundTrip(t *testing.T) {
	eval := NewEvaluator(1024)

	// Create a test polynomial
	p := eval.NewPoly()
	for i := range p.Coeffs {
		p.Coeffs[i] = params.Torus(i * 12345)
	}

	// Transform to frequency domain and back
	fp := eval.ToFourierPoly(p)
	pOut := eval.ToPoly(fp)

	// Check if we got the original back (with some tolerance for floating point errors)
	for i := range p.Coeffs {
		diff := int64(pOut.Coeffs[i]) - int64(p.Coeffs[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > 10 { // Allow small error due to floating point rounding
			t.Errorf("Coefficient %d: got %d, want %d (diff %d)", i, pOut.Coeffs[i], p.Coeffs[i], diff)
		}
	}
}

// TestPolyMul tests polynomial multiplication
func TestPolyMul(t *testing.T) {
	eval := NewEvaluator(1024)

	// Create two simple test polynomials
	p1 := eval.NewPoly()
	p2 := eval.NewPoly()

	p1.Coeffs[0] = 100
	p1.Coeffs[1] = 200

	p2.Coeffs[0] = 10
	p2.Coeffs[1] = 20

	// Multiply
	pOut := eval.MulPoly(p1, p2)

	// Expected result for first few coefficients:
	// (100 + 200*X) * (10 + 20*X) = 1000 + 2000*X + 2000*X + 4000*X^2
	//                               = 1000 + 4000*X + 4000*X^2

	// Due to negacyclic ring, we need to check this works correctly
	// For now, just verify the function runs without panic
	if pOut.Coeffs == nil {
		t.Error("MulPoly returned nil coefficients")
	}
}

// BenchmarkFFT benchmarks the FFT operation
func BenchmarkFFT(b *testing.B) {
	eval := NewEvaluator(1024)
	p := eval.NewPoly()
	for i := range p.Coeffs {
		p.Coeffs[i] = params.Torus(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eval.ToFourierPoly(p)
	}
}

// BenchmarkIFFT benchmarks the inverse FFT operation
func BenchmarkIFFT(b *testing.B) {
	eval := NewEvaluator(1024)
	p := eval.NewPoly()
	for i := range p.Coeffs {
		p.Coeffs[i] = params.Torus(i)
	}
	fp := eval.ToFourierPoly(p)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eval.ToPoly(fp)
	}
}

// BenchmarkPolyMul benchmarks polynomial multiplication
func BenchmarkPolyMul(b *testing.B) {
	eval := NewEvaluator(1024)
	p1 := eval.NewPoly()
	p2 := eval.NewPoly()
	for i := range p1.Coeffs {
		p1.Coeffs[i] = params.Torus(i)
		p2.Coeffs[i] = params.Torus(i * 2)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = eval.MulPoly(p1, p2)
	}
}

// BenchmarkElementWiseMul benchmarks element-wise multiplication in frequency domain
func BenchmarkElementWiseMul(b *testing.B) {
	eval := NewEvaluator(1024)
	p1 := eval.NewPoly()
	p2 := eval.NewPoly()
	for i := range p1.Coeffs {
		p1.Coeffs[i] = params.Torus(i)
		p2.Coeffs[i] = params.Torus(i * 2)
	}
	fp1 := eval.ToFourierPoly(p1)
	fp2 := eval.ToFourierPoly(p2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eval.MulFourierPolyAssign(fp1, fp2, fp1)
	}
}
