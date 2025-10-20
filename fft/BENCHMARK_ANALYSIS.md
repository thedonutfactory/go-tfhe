# FFT Benchmark Analysis and Optimization Opportunities

## Benchmark Results Summary (Apple M3 Pro)

### Core Operations Performance
| Operation | Time (ns/op) | Allocations | Notes |
|-----------|--------------|-------------|-------|
| IFFT1024 | 3,214 | 0 B/op, 0 allocs/op | ✅ Zero-allocation hot path |
| FFT1024 | 3,223 | 0 B/op, 0 allocs/op | ✅ Zero-allocation hot path |
| PolyMul1024 | 10,647 | 0 B/op, 0 allocs/op | ✅ Zero-allocation hot path |
| FFT Roundtrip | 6,496 | 0 B/op, 0 allocs/op | ✅ Zero-allocation hot path |

### Initialization
| Operation | Time (ns/op) | Allocations | Notes |
|-----------|--------------|-------------|-------|
| NewFFTProcessor | 10,497 | 62,720 B, 9 allocs | ⚠️ Expensive - reuse processors! |

### Batch Operations (per batch)
| Batch Size | IFFT (ns/op) | FFT (ns/op) | Allocations |
|------------|--------------|-------------|-------------|
| 1 poly | 3,733 | 4,337 | 1 alloc |
| 2 polys | 7,370 | 6,888 | 1 alloc |
| 4 polys | 14,538 | 13,961 | 1 alloc |
| 8 polys | 29,109 | 27,500 | 1 alloc |
| 16 polys | 67,822 | 55,022 | 1 alloc |
| 32 polys | 115,298 | 109,952 | 1 alloc |

### Array vs Slice Performance
| Operation | Array (ns/op) | Slice (ns/op) | Difference |
|-----------|---------------|---------------|------------|
| IFFT | 3,202 | 3,325 | +3.8% slower |
| PolyMul | 10,749 | 11,193 | +4.1% slower + 1 alloc |

### Processor Reuse vs Recreation
| Strategy | Time (ns/op) | Allocations | Speedup |
|----------|--------------|-------------|---------|
| Reuse processor | 8,018 | 0 B/op, 0 allocs | ✅ Baseline |
| Create per call | 17,250 | 62,720 B, 9 allocs | ❌ 2.15x slower |

## Key Findings

### ✅ Strengths
1. **Zero-allocation hot path**: Core FFT operations have perfect memory efficiency
2. **Consistent performance**: No worst-case input patterns (zeros, random, sequential all similar)
3. **Linear scaling**: Batch operations scale linearly with batch size
4. **SIMD optimization working**: ~3.2µs per 1024-point FFT is excellent

### ⚠️ Optimization Opportunities

#### 1. Processor Reuse (CRITICAL)
**Finding**: Creating a new processor for each operation is **2.15x slower** and allocates 62KB.

**Recommendation**: 
- Always reuse `FFTProcessor` instances
- Consider processor pooling for concurrent scenarios
- Document the importance of reuse in API

```go
// ❌ DON'T DO THIS
for i := 0; i < n; i++ {
    proc := fft.NewFFTProcessor(1024)
    result := proc.PolyMul1024(&a, &b)
}

// ✅ DO THIS INSTEAD
proc := fft.NewFFTProcessor(1024)
for i := 0; i < n; i++ {
    result := proc.PolyMul1024(&a, &b)
}
```

#### 2. Batch Operations Allocations
**Finding**: Each batch operation allocates a new result slice.

**Current code** (lines 241-246):
```go
func (p *FFTProcessor) BatchIFFT1024(inputs [][1024]params.Torus) [][1024]float64 {
    results := make([][1024]float64, len(inputs))  // New allocation every call
    for i := range inputs {
        results[i] = p.IFFT1024(&inputs[i])
    }
    return results
}
```

**Recommendation**: Add variants that accept pre-allocated result buffers:
```go
func (p *FFTProcessor) BatchIFFT1024InPlace(inputs [][1024]params.Torus, results [][1024]float64) {
    for i := range inputs {
        results[i] = p.IFFT1024(&inputs[i])
    }
}
```

#### 3. Slice-based API Overhead
**Finding**: Slice-based functions are 4% slower and allocate memory.

**Recommendation**: 
- Prioritize array-based API for performance-critical code
- Document the performance difference
- Consider adding variants that write to caller-provided buffers:

```go
func (p *FFTProcessor) IFFT1024ToBuf(input *[1024]params.Torus, output []float64) {
    result := p.IFFT1024(input)
    copy(output, result[:])
}
```

#### 4. Parallel Batch Processing
**Finding**: Batch operations process sequentially. For large batches, this could benefit from parallelization.

**Recommendation**: For batch sizes > 8, use goroutines:
```go
func (p *FFTProcessor) BatchIFFT1024Parallel(inputs [][1024]params.Torus) [][1024]float64 {
    results := make([][1024]float64, len(inputs))
    if len(inputs) <= 8 {
        // Sequential for small batches
        for i := range inputs {
            results[i] = p.IFFT1024(&inputs[i])
        }
        return results
    }
    
    // Parallel for large batches
    var wg sync.WaitGroup
    for i := range inputs {
        wg.Add(1)
        go func(idx int) {
            defer wg.Done()
            results[idx] = p.IFFT1024(&inputs[idx])
        }(i)
    }
    wg.Wait()
    return results
}
```

**Note**: Each goroutine would need its own `FFTProcessor` to be thread-safe, or use sync.Pool.

## Optimization Priority

1. **HIGH**: Document processor reuse pattern (easy, big impact)
2. **MEDIUM**: Add zero-allocation batch operation variants
3. **MEDIUM**: Add parallel batch processing for large batches
4. **LOW**: Profile and optimize the FFTInPlace/IFFTInPlace calls (already very fast)
5. **LOW**: Consider cache-friendly memory layout for very large batches

## Profiling Next Steps

To identify micro-optimizations within the FFT itself:

```bash
# CPU profiling
go test -bench=BenchmarkPolyMul1024 -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Look for:
# - Hot loops in FFTInPlace/IFFTInPlace
# - Conversion overhead (CmplxToFloat4, Float4ToCmplx)
# - Cache misses in large data structures
```

## Conclusion

The current FFT implementation is **already very well optimized**:
- Zero allocations in hot path
- Efficient SIMD usage
- ~3.2µs for 1024-point FFT is excellent

The main optimization opportunity is at the **API usage level**: ensuring processors are reused and providing zero-allocation variants for batch operations.




