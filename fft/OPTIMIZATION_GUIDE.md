# FFT Optimization Guide - Detailed Analysis

## Executive Summary

**Current Performance**: ~3.2µs per FFT/IFFT, ~10.8µs per PolyMul
- ✅ Zero allocations in hot path
- ✅ SIMD optimizations working well  
- ✅ Already 8-10x faster than naive implementation

**Key Finding**: The implementation is already highly optimized. The main optimization opportunities are:
1. **API-level optimizations** (processor reuse, batch operations)
2. **Potential SIMD improvements** for conversion loops
3. **Minor micro-optimizations** in conversion code

## CPU Profile Analysis (PolyMul1024)

### Time Distribution (2.30s total, 224,378 iterations)

| Component | Time (ms) | % of Total | % of PolyMul | Notes |
|-----------|-----------|------------|--------------|-------|
| **IFFT1024** | 1,280 | 55.7% | 60.4% | Time → Freq domain |
| └─ Input conversion (Torus→complex) | 260 | 11.3% | 12.3% | Lines 126-129 |
| └─ CmplxToFloat4Assign | 80 | 3.5% | 3.8% | Line 133 |
| └─ FFTInPlace (SIMD) | 600 | 26.1% | 28.3% | Line 137 ⚡ |
| └─ Float4ToCmplxAssign | 140 | 6.1% | 6.6% | Line 140 |
| └─ Output scaling | 170 | 7.4% | 8.0% | Lines 143-145 |
| **FFT1024** | 640 | 27.8% | 30.2% | Freq → Time domain |
| └─ Input scaling | 50 | 2.2% | 2.4% | Lines 158-159 |
| └─ CmplxToFloat4Assign | 60 | 2.6% | 2.8% | Line 163 |
| └─ IFFTInPlace (SIMD) | 390 | 17.0% | 18.4% | Line 167 ⚡ |
| └─ Float4ToCmplxAssign | 60 | 2.6% | 2.8% | Line 170 |
| └─ Output conversion (complex→Torus) | 80 | 3.5% | 3.8% | Lines 174-181 |
| **Complex multiplication** | 200 | 8.7% | 9.4% | Lines 212-221 |
| **Other overhead** | 10 | 0.4% | 0.5% | Function calls, etc. |

### Key Insights

1. **SIMD FFT operations** (fftInPlace + ifftInPlace): **46.7%** of total time
   - Already optimized with assembly
   - Very difficult to improve further

2. **Conversion operations** (Float4 ↔ Complex): **15.0%** of total time
   - Potential for SIMD optimization
   - Currently using scalar loops

3. **Type conversions** (Torus → Complex, Complex → Torus): **17.6%** of total time
   - Potential for optimization
   - Currently using scalar loops with type casts

4. **Complex multiplication**: **9.4%** of total time
   - Simple arithmetic, hard to optimize

## Optimization Recommendations

### 🔴 Priority 1: High Impact, Easy Implementation

#### 1.1 Add Processor Pool for Concurrent Use

**Impact**: Prevents recreating processors (2.15x speedup)

```go
// Add to fft.go
var processorPool = sync.Pool{
    New: func() interface{} {
        return NewFFTProcessor(1024)
    },
}

// GetProcessor returns a processor from the pool
func GetProcessor() *FFTProcessor {
    return processorPool.Get().(*FFTProcessor)
}

// PutProcessor returns a processor to the pool
func PutProcessor(p *FFTProcessor) {
    processorPool.Put(p)
}
```

**Usage**:
```go
proc := fft.GetProcessor()
defer fft.PutProcessor(proc)
result := proc.PolyMul1024(&a, &b)
```

#### 1.2 Document Processor Reuse Pattern

Add to package documentation and README:

```go
// ⚠️ PERFORMANCE CRITICAL: Always reuse FFTProcessor instances
// Creating a new processor for each operation is 2.15x slower
// and allocates 62KB per call.
//
// Good:
//   proc := fft.NewFFTProcessor(1024)
//   for i := 0; i < n; i++ {
//       result := proc.PolyMul1024(&a, &b)
//   }
//
// Bad:
//   for i := 0; i < n; i++ {
//       proc := fft.NewFFTProcessor(1024)  // ❌ Don't do this!
//       result := proc.PolyMul1024(&a, &b)
//   }
```

### 🟡 Priority 2: Medium Impact, Moderate Effort

#### 2.1 Optimize Input Conversion Loop (IFFT1024, Lines 126-129)

**Current** (260ms, 12.3% of PolyMul time):
```go
for i := 0; i < N2; i++ {
    inRe := float64(int32(input[i]))
    inIm := float64(int32(input[i+N2]))
    p.fourierBuffer[i] = complex(inRe, inIm)
}
```

**Potential optimization**: Vectorize this loop with SIMD
- Convert multiple Torus values at once
- Could use assembly or Go 1.21+ SIMD intrinsics

**Estimated impact**: 2-3x speedup on this section = ~1.5% total speedup

#### 2.2 Optimize Output Scaling Loop (IFFT1024, Lines 143-145)

**Current** (170ms, 8.0% of PolyMul time):
```go
for i := 0; i < N2; i++ {
    p.resultBuffer[i] = real(p.fourierBuffer[i]) * 2.0
    p.resultBuffer[i+N2] = imag(p.fourierBuffer[i]) * 2.0
}
```

**Potential optimization**: Combine with Float4ToCmplxAssign
- Do the scaling during conversion, not after
- Saves one pass over the data (better cache efficiency)

**Proposed**:
```go
// Modify vec.Float4ToCmplxAssign to accept a scale factor
func Float4ToCmplxAndScale(float4 []float64, cmplx []complex128, scale float64)

// Then in IFFT1024:
vec.Float4ToCmplxAndScale(p.float4Buffer, p.fourierBuffer, 2.0)
for i := 0; i < N2; i++ {
    p.resultBuffer[i] = real(p.fourierBuffer[i])
    p.resultBuffer[i+N2] = imag(p.fourierBuffer[i])
}
```

**Estimated impact**: ~3-5% total speedup

#### 2.3 Add Zero-Allocation Batch Operations

**Current**:
```go
func (p *FFTProcessor) BatchIFFT1024(inputs [][1024]params.Torus) [][1024]float64 {
    results := make([][1024]float64, len(inputs))  // Allocates on every call
    for i := range inputs {
        results[i] = p.IFFT1024(&inputs[i])
    }
    return results
}
```

**Add**:
```go
// BatchIFFT1024InPlace performs batch IFFT without allocating result slice
func (p *FFTProcessor) BatchIFFT1024InPlace(inputs [][1024]params.Torus, results [][1024]float64) {
    if len(results) < len(inputs) {
        panic("results slice too small")
    }
    for i := range inputs {
        results[i] = p.IFFT1024(&inputs[i])
    }
}
```

**Estimated impact**: Eliminates batch allocations (see benchmark results)

### 🟢 Priority 3: Lower Impact, More Effort

#### 3.1 Optimize Complex Multiplication Loop (Lines 212-221)

**Current** (200ms, 9.4% of PolyMul time):
```go
for i := 0; i < N2; i++ {
    ar := aFFT[i]
    ai := aFFT[i+N2]
    br := bFFT[i]
    bi := bFFT[i+N2]
    
    // Complex multiply: (ar + i*ai) * (br + i*bi) * 0.5
    resultFFT[i] = (ar*br - ai*bi) * 0.5
    resultFFT[i+N2] = (ar*bi + ai*br) * 0.5
}
```

**Potential optimization**: SIMD complex multiplication
- Process 4-8 complex numbers at once
- Could use assembly or SIMD intrinsics

**Estimated impact**: 2-3x speedup on this section = ~0.5-1% total speedup

#### 3.2 Investigate Cache Optimization

Profile cache misses with:
```bash
perf stat -e cache-misses,cache-references go test -bench=BenchmarkPolyMul1024
```

Potential optimizations:
- Align buffers to cache line boundaries
- Prefetch data in conversion loops
- Rearrange data layout for better spatial locality

**Estimated impact**: 1-3% total speedup (highly hardware-dependent)

#### 3.3 Parallel Batch Processing

For large batches (>8 polynomials), parallelize:

```go
func (p *FFTProcessor) BatchIFFT1024Parallel(inputs [][1024]params.Torus, numWorkers int) [][1024]float64 {
    if len(inputs) <= 8 || numWorkers <= 1 {
        return p.BatchIFFT1024(inputs) // Sequential for small batches
    }
    
    // Create processor pool
    processors := make([]*FFTProcessor, numWorkers)
    for i := range processors {
        processors[i] = NewFFTProcessor(1024)
    }
    
    results := make([][1024]float64, len(inputs))
    
    // Distribute work
    var wg sync.WaitGroup
    chunkSize := (len(inputs) + numWorkers - 1) / numWorkers
    
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            proc := processors[workerID]
            start := workerID * chunkSize
            end := min(start+chunkSize, len(inputs))
            
            for i := start; i < end; i++ {
                results[i] = proc.IFFT1024(&inputs[i])
            }
        }(w)
    }
    
    wg.Wait()
    return results
}
```

**Estimated impact**: Near-linear speedup for large batches (32+ polynomials)

### ⚪ Priority 4: Advanced Optimizations

#### 4.1 Assembly-Optimized Conversion Functions

Write assembly implementations for:
- `torusToComplex`: Torus → complex128 conversion
- `complexToTorus`: complex128 → Torus conversion  
- `complexMulScale`: Complex multiplication with scaling

**Estimated impact**: 5-10% total speedup
**Effort**: High (requires assembly expertise)

#### 4.2 Eliminate Float4 Intermediate Format

**Current flow**:
```
Torus → Complex → Float4 → FFT → Float4 → Complex → Torus
```

**Proposed**: Make FFT work directly on complex128
```
Torus → Complex → FFT → Complex → Torus
```

**Challenges**:
- Requires modifying poly.FFTInPlace to work on complex128
- May lose SIMD benefits if Float4 format is optimal for assembly

**Estimated impact**: 5-15% total speedup (if successful)
**Effort**: Very High (requires deep understanding of assembly FFT)

## Benchmarking Checklist

After each optimization, run:

```bash
# 1. Correctness tests
go test -v ./fft

# 2. Performance benchmarks
go test -bench=. -benchmem -benchtime=2s ./fft

# 3. CPU profiling
go test -bench=BenchmarkPolyMul1024 -cpuprofile=cpu.prof -benchtime=2s
go tool pprof -top cpu.prof

# 4. Memory profiling  
go test -bench=BenchmarkPolyMul1024 -memprofile=mem.prof -benchtime=2s
go tool pprof -top mem.prof

# 5. Assembly inspection (to verify optimizations)
go build -gcflags="-S" ./fft 2>&1 | grep -A 20 "FFT1024"
```

## Estimated Total Impact

| Optimization | Effort | Impact | Priority |
|--------------|--------|--------|----------|
| Processor pool/docs | Low | 2.15x (if misused) | 🔴 High |
| Zero-alloc batch API | Low | Elim. batch allocs | 🔴 High |
| Optimize input conversion | Medium | ~1.5% | 🟡 Medium |
| Optimize output scaling | Medium | ~3-5% | 🟡 Medium |
| SIMD complex multiply | Medium | ~0.5-1% | 🟢 Low |
| Parallel batching | Medium | ~Nx (N workers) | 🟢 Low |
| Assembly conversions | High | ~5-10% | ⚪ Advanced |
| Eliminate Float4 | Very High | ~5-15% | ⚪ Advanced |

**Realistic best case**: 
- Easy optimizations: 0% (already well-optimized)
- Medium optimizations: 5-8% improvement
- Advanced optimizations: 15-25% improvement

**Recommendation**: Focus on Priority 1 and 2 optimizations first. The code is already very well optimized, so expect diminishing returns.

## Alternative Approach: Algorithm-Level Optimization

Consider these alternative strategies:

1. **Polynomial Representation**: Keep polynomials in FFT domain longer
   - Avoid repeated FFT/IFFT conversions
   - Cache FFT results when possible

2. **Batch Processing**: Process multiple polynomials together
   - Better CPU utilization
   - Amortize setup costs

3. **Approximate Computation**: For TFHE applications
   - Check if slightly lower precision is acceptable
   - Use faster approximate FFT algorithms

4. **Hardware Acceleration**: 
   - GPU implementation for large batches
   - Custom FPGA/ASIC for production systems

These are application-level optimizations that may provide better returns than micro-optimizations.

