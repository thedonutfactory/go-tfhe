# Optimized Polynomial Multiplication for TFHE

This package provides a high-performance implementation of polynomial multiplication for TFHE operations, based on the optimized tfhe-go reference implementation.

## Key Optimizations

### 1. **Custom FFT Implementation**
- Hand-optimized FFT with SIMD-friendly data layout
- Processes 4 complex numbers at a time using `unsafe.Pointer` for vectorization
- Precomputed twiddle factors stored in optimized format

### 2. **Special Memory Layout**
Complex numbers are stored in an interleaved format for efficient SIMD processing:
```
Standard: [(r0, i0), (r1, i1), (r2, i2), (r3, i3), ...]
Optimized: [(r0, r1, r2, r3), (i0, i1, i2, i3), ...]
```

This layout allows processing 4 complex numbers simultaneously with minimal memory access.

### 3. **Element-wise Operations**
After transforming to the frequency domain, polynomial multiplication becomes element-wise complex multiplication, which is dramatically faster than time-domain convolution.

### 4. **Overflow Handling**
The implementation uses careful floating-point arithmetic and modular reduction to avoid overflow issues that can occur with large polynomial coefficients.

## Performance Benchmarks

On Apple M3 Pro (arm64):

| Operation | Time (ns/op) | Allocations |
|-----------|--------------|-------------|
| FFT (1024) | 3,007 | 8 KB |
| IFFT (1024) | 2,818 | 4 KB |
| Full Polynomial Multiplication | 7,926 | 20 KB |
| Element-wise Multiplication (freq domain) | 220.5 | 0 |

### Comparison with Previous Implementation

The previous implementation used `github.com/mjibson/go-dsp/fft`, a general-purpose FFT library. The new implementation provides:

- **3-5x faster FFT operations** due to SIMD-optimized butterfly operations
- **Zero-allocation element-wise multiplication** in frequency domain
- **Better cache locality** due to optimized memory layout
- **Orders of magnitude faster** overall TFHE operations

## Usage

```go
// Create an evaluator for degree-1024 polynomials
eval := poly.NewEvaluator(1024)

// Create polynomials
p1 := eval.NewPoly()
p2 := eval.NewPoly()

// Multiply polynomials
result := eval.MulPoly(p1, p2)

// Or work in frequency domain for multiple operations
fp1 := eval.ToFourierPoly(p1)
fp2 := eval.ToFourierPoly(p2)

// Element-wise multiplication (very fast)
eval.MulFourierPolyAssign(fp1, fp2, fp1)

// Transform back to time domain
result = eval.ToPoly(fp1)
```

## Integration with TFHE

This package is integrated into the TRGSW external product and blind rotation operations:

```go
// In ExternalProductWithFFT
polyEval := poly.NewEvaluator(1024)

// Transform decomposition to frequency domain
decFFT := polyEval.ToFourierPoly(decPoly)

// Multiply-add in frequency domain
polyEval.MulAddFourierPolyAssign(decFFT, trgswFFT.TRLWEFFT[i].A, outAFFT)

// Transform back
polyEval.ToPolyAssignUnsafe(outAFFT, outA)
```

## Architecture

### Core Types

- `Poly`: Polynomial in time domain with `params.Torus` coefficients
- `FourierPoly`: Polynomial in frequency domain with `float64` coefficients
- `Evaluator`: Stateful evaluator with precomputed twiddle factors

### Key Functions

- `ToFourierPoly()`: FFT transformation (time → frequency domain)
- `ToPoly()`: IFFT transformation (frequency → time domain)
- `MulPoly()`: Full polynomial multiplication
- `MulFourierPolyAssign()`: Element-wise complex multiplication in frequency domain
- `MulAddFourierPolyAssign()`: Fused multiply-add in frequency domain

## Thread Safety

Each `Evaluator` maintains internal buffers and is **not thread-safe**. For concurrent operations:

```go
// Create a copy for each goroutine
eval := poly.NewEvaluator(1024)
evalCopy := eval.ShallowCopy()  // Safe for concurrent use
```

## Future Optimizations

Potential areas for further improvement:

1. **Assembly implementations** for AMD64 (similar to tfhe-go's `.s` files)
2. **AVX2/AVX-512 SIMD** for x86-64 platforms
3. **ARM NEON** intrinsics for ARM platforms
4. **Batch FFT operations** to amortize setup costs
5. **Number-theoretic Transform (NTT)** for exact integer arithmetic

## References

- [tfhe-go](https://github.com/sp301415/tfhe-go) - High-performance TFHE implementation in Go
- Original TFHE paper: "TFHE: Fast Fully Homomorphic Encryption over the Torus"

