package fft_test

import (
	"math/rand"
	"testing"

	"github.com/thedonutfactory/go-tfhe/fft"
	"github.com/thedonutfactory/go-tfhe/params"
)

// TestFFTRoundtrip tests that IFFT followed by FFT returns the original input
func TestFFTRoundtrip(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	var input [1024]params.Torus
	for i := range input {
		input[i] = params.Torus(rng.Uint32())
	}

	freq := proc.IFFT1024(&input)
	output := proc.FFT1024(&freq)

	var maxDiff int64
	for i := 0; i < 1024; i++ {
		diff := int64(output[i]) - int64(input[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	if maxDiff >= 2 {
		t.Errorf("FFT roundtrip error too large: %d (should be < 2)", maxDiff)
		t.Logf("First 10 values:")
		for i := 0; i < 10; i++ {
			t.Logf("  [%d] in:%d out:%d diff:%d", i, input[i], output[i], int64(output[i])-int64(input[i]))
		}
	}
}

// TestFFTSimple tests FFT with simple delta function input
func TestFFTSimple(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)

	// Delta function test: single non-zero value
	var input [1024]params.Torus
	input[0] = 1000

	freq := proc.IFFT1024(&input)
	output := proc.FFT1024(&freq)

	diff := int64(output[0]) - int64(input[0])
	if diff < 0 {
		diff = -diff
	}

	if diff >= 10 {
		t.Errorf("Delta function roundtrip error: %d (should be < 10)", diff)
		t.Logf("input[0]=%d, output[0]=%d", input[0], output[0])
	}
}

// TestPolyMul1024 tests polynomial multiplication against naive implementation
func TestPolyMul1024(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	trials := 100
	for trial := 0; trial < trials; trial++ {
		var a, b [1024]params.Torus
		for i := range a {
			a[i] = params.Torus(rng.Uint32())
			// Keep b VERY small (like Rust tests use params::trgsw_lv1::BG = 64)
			b[i] = params.Torus(rng.Uint32()) % params.Torus(params.GetTRGSWLv1().BG)
		}

		fftResult := proc.PolyMul1024(&a, &b)
		naiveResult := naivePolyMul(&a, &b)

		var maxDiff int64
		for i := 0; i < 1024; i++ {
			diff := int64(fftResult[i]) - int64(naiveResult[i])
			if diff < 0 {
				diff = -diff
			}
			if diff > maxDiff {
				maxDiff = diff
			}
		}

		if maxDiff >= 2 {
			t.Errorf("Trial %d: Polynomial multiplication error too large: %d", trial, maxDiff)
			t.Logf("First 5 mismatches:")
			count := 0
			for i := 0; i < 1024 && count < 5; i++ {
				diff := int64(fftResult[i]) - int64(naiveResult[i])
				if diff < 0 {
					diff = -diff
				}
				if diff >= 2 {
					t.Logf("  [%d] FFT:%d Naive:%d Diff:%d", i, fftResult[i], naiveResult[i], int64(fftResult[i])-int64(naiveResult[i]))
					count++
				}
			}
			break
		}
	}
}

// naivePolyMul computes negacyclic polynomial multiplication naively
// a(X) * b(X) mod (X^N+1)
func naivePolyMul(a, b *[1024]params.Torus) [1024]params.Torus {
	var result [1024]params.Torus
	const N = 1024

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if i+j < N {
				result[i+j] += a[i] * b[j]
			} else {
				// Wrap around with negation (X^N = -1)
				result[i+j-N] -= a[i] * b[j]
			}
		}
	}

	return result
}

// TestIFFTSlice tests the slice-based IFFT function
func TestIFFTSlice(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	input := make([]params.Torus, 1024)
	for i := range input {
		input[i] = params.Torus(rng.Uint32())
	}

	freq := proc.IFFT(input)
	output := proc.FFT(freq)

	if len(output) != len(input) {
		t.Fatalf("Output length %d != input length %d", len(output), len(input))
	}

	var maxDiff int64
	for i := 0; i < len(input); i++ {
		diff := int64(output[i]) - int64(input[i])
		if diff < 0 {
			diff = -diff
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	if maxDiff >= 2 {
		t.Errorf("Slice FFT roundtrip error: %d (should be < 2)", maxDiff)
	}
}

// TestPolyMulSlice tests the slice-based polynomial multiplication
func TestPolyMulSlice(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	a := make([]params.Torus, 1024)
	b := make([]params.Torus, 1024)
	for i := range a {
		a[i] = params.Torus(rng.Uint32())
		b[i] = params.Torus(rng.Uint32() % 64)
	}

	result := proc.PolyMul(a, b)

	if len(result) != 1024 {
		t.Fatalf("PolyMul result length %d != 1024", len(result))
	}

	// Verify first few values against naive
	var aArr, bArr [1024]params.Torus
	copy(aArr[:], a)
	copy(bArr[:], b)
	naive := naivePolyMul(&aArr, &bArr)

	for i := 0; i < 10; i++ {
		diff := int64(result[i]) - int64(naive[i])
		if diff < 0 {
			diff = -diff
		}
		if diff >= 2 {
			t.Errorf("PolyMul[%d]: FFT=%d Naive=%d Diff=%d", i, result[i], naive[i], diff)
		}
	}
}

// TestBatchIFFT tests batch IFFT operation
func TestBatchIFFT(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	inputs := make([][1024]params.Torus, 3)
	for i := range inputs {
		for j := range inputs[i] {
			inputs[i][j] = params.Torus(rng.Uint32())
		}
	}

	results := proc.BatchIFFT1024(inputs)

	if len(results) != len(inputs) {
		t.Fatalf("BatchIFFT returned %d results, expected %d", len(results), len(inputs))
	}

	// Verify each result matches individual IFFT
	for i := range inputs {
		expected := proc.IFFT1024(&inputs[i])
		for j := 0; j < 1024; j++ {
			if results[i][j] != expected[j] {
				t.Errorf("BatchIFFT[%d][%d] = %f, individual IFFT = %f", i, j, results[i][j], expected[j])
				break
			}
		}
	}
}

// TestBatchFFT tests batch FFT operation
func TestBatchFFT(t *testing.T) {
	proc := fft.NewFFTProcessor(1024)
	rng := rand.New(rand.NewSource(42))

	inputs := make([][1024]float64, 3)
	for i := range inputs {
		for j := range inputs[i] {
			inputs[i][j] = rng.Float64() * 1000
		}
	}

	results := proc.BatchFFT1024(inputs)

	if len(results) != len(inputs) {
		t.Fatalf("BatchFFT returned %d results, expected %d", len(results), len(inputs))
	}

	// Verify each result matches individual FFT
	for i := range inputs {
		expected := proc.FFT1024(&inputs[i])
		for j := 0; j < 1024; j++ {
			if results[i][j] != expected[j] {
				t.Errorf("BatchFFT[%d][%d] = %d, individual FFT = %d", i, j, results[i][j], expected[j])
				break
			}
		}
	}
}

// BenchmarkIFFT1024 benchmarks the optimized IFFT operation
func BenchmarkIFFT1024(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)
	var input [1024]params.Torus
	for i := range input {
		input[i] = params.Torus(i * 1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = proc.IFFT1024(&input)
	}
}

// BenchmarkFFT1024 benchmarks the optimized FFT operation
func BenchmarkFFT1024(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)
	var input [1024]float64
	for i := range input {
		input[i] = float64(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = proc.FFT1024(&input)
	}
}

// BenchmarkPolyMul1024 benchmarks polynomial multiplication
func BenchmarkPolyMul1024(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)
	var a, bb [1024]params.Torus
	for i := range a {
		a[i] = params.Torus(i * 1000)
		bb[i] = params.Torus(i * 10)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = proc.PolyMul1024(&a, &bb)
	}
}

// BenchmarkFFTRoundtrip benchmarks a complete FFT roundtrip
func BenchmarkFFTRoundtrip(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)
	var input [1024]params.Torus
	for i := range input {
		input[i] = params.Torus(i * 1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		freq := proc.IFFT1024(&input)
		_ = proc.FFT1024(&freq)
	}
}

// BenchmarkNewFFTProcessor benchmarks the initialization overhead
func BenchmarkNewFFTProcessor(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fft.NewFFTProcessor(1024)
	}
}

// BenchmarkBatchIFFT1024 benchmarks batch IFFT operations
func BenchmarkBatchIFFT1024(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	// Test with different batch sizes
	for _, batchSize := range []int{1, 2, 4, 8, 16, 32} {
		b.Run(string(rune('0'+batchSize/10))+string(rune('0'+batchSize%10))+"_polys", func(b *testing.B) {
			inputs := make([][1024]params.Torus, batchSize)
			for i := range inputs {
				for j := range inputs[i] {
					inputs[i][j] = params.Torus(j * 1000)
				}
			}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = proc.BatchIFFT1024(inputs)
			}
		})
	}
}

// BenchmarkBatchFFT1024 benchmarks batch FFT operations
func BenchmarkBatchFFT1024(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	// Test with different batch sizes
	for _, batchSize := range []int{1, 2, 4, 8, 16, 32} {
		b.Run(string(rune('0'+batchSize/10))+string(rune('0'+batchSize%10))+"_polys", func(b *testing.B) {
			inputs := make([][1024]float64, batchSize)
			for i := range inputs {
				for j := range inputs[i] {
					inputs[i][j] = float64(j)
				}
			}

			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = proc.BatchFFT1024(inputs)
			}
		})
	}
}

// BenchmarkPolyMulWithDifferentMagnitudes benchmarks polynomial multiplication
// with different input magnitudes (small vs large values)
func BenchmarkPolyMulWithDifferentMagnitudes(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	b.Run("small_values", func(b *testing.B) {
		var a, bb [1024]params.Torus
		for i := range a {
			a[i] = params.Torus(i % 64)
			bb[i] = params.Torus(i % 64)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.PolyMul1024(&a, &bb)
		}
	})

	b.Run("medium_values", func(b *testing.B) {
		var a, bb [1024]params.Torus
		for i := range a {
			a[i] = params.Torus(i * 1000)
			bb[i] = params.Torus(i * 1000)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.PolyMul1024(&a, &bb)
		}
	})

	b.Run("large_values", func(b *testing.B) {
		rng := rand.New(rand.NewSource(42))
		var a, bb [1024]params.Torus
		for i := range a {
			a[i] = params.Torus(rng.Uint32())
			bb[i] = params.Torus(rng.Uint32())
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.PolyMul1024(&a, &bb)
		}
	})
}

// BenchmarkIFFT1024WithPatterns benchmarks IFFT with different input patterns
func BenchmarkIFFT1024WithPatterns(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	b.Run("zeros", func(b *testing.B) {
		var input [1024]params.Torus
		// All zeros

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT1024(&input)
		}
	})

	b.Run("delta", func(b *testing.B) {
		var input [1024]params.Torus
		input[0] = 1000

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT1024(&input)
		}
	})

	b.Run("sequential", func(b *testing.B) {
		var input [1024]params.Torus
		for i := range input {
			input[i] = params.Torus(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT1024(&input)
		}
	})

	b.Run("random", func(b *testing.B) {
		rng := rand.New(rand.NewSource(42))
		var input [1024]params.Torus
		for i := range input {
			input[i] = params.Torus(rng.Uint32())
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT1024(&input)
		}
	})
}

// BenchmarkFFT1024WithPatterns benchmarks FFT with different input patterns
func BenchmarkFFT1024WithPatterns(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	b.Run("zeros", func(b *testing.B) {
		var input [1024]float64
		// All zeros

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.FFT1024(&input)
		}
	})

	b.Run("sequential", func(b *testing.B) {
		var input [1024]float64
		for i := range input {
			input[i] = float64(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.FFT1024(&input)
		}
	})

	b.Run("random", func(b *testing.B) {
		rng := rand.New(rand.NewSource(42))
		var input [1024]float64
		for i := range input {
			input[i] = rng.Float64() * 1000
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.FFT1024(&input)
		}
	})
}

// BenchmarkSliceVsArray benchmarks slice-based vs array-based FFT operations
func BenchmarkSliceVsArray(b *testing.B) {
	proc := fft.NewFFTProcessor(1024)

	b.Run("IFFT_array", func(b *testing.B) {
		var input [1024]params.Torus
		for i := range input {
			input[i] = params.Torus(i * 1000)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT1024(&input)
		}
	})

	b.Run("IFFT_slice", func(b *testing.B) {
		input := make([]params.Torus, 1024)
		for i := range input {
			input[i] = params.Torus(i * 1000)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.IFFT(input)
		}
	})

	b.Run("PolyMul_array", func(b *testing.B) {
		var a, bb [1024]params.Torus
		for i := range a {
			a[i] = params.Torus(i * 1000)
			bb[i] = params.Torus(i * 10)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.PolyMul1024(&a, &bb)
		}
	})

	b.Run("PolyMul_slice", func(b *testing.B) {
		a := make([]params.Torus, 1024)
		bb := make([]params.Torus, 1024)
		for i := range a {
			a[i] = params.Torus(i * 1000)
			bb[i] = params.Torus(i * 10)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = proc.PolyMul(a, bb)
		}
	})
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("processor_reuse", func(b *testing.B) {
		proc := fft.NewFFTProcessor(1024)
		var input [1024]params.Torus
		for i := range input {
			input[i] = params.Torus(i * 1000)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			freq := proc.IFFT1024(&input)
			_ = proc.FFT1024(&freq)
		}
	})

	b.Run("processor_per_call", func(b *testing.B) {
		var input [1024]params.Torus
		for i := range input {
			input[i] = params.Torus(i * 1000)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			proc := fft.NewFFTProcessor(1024)
			freq := proc.IFFT1024(&input)
			_ = proc.FFT1024(&freq)
		}
	})
}
