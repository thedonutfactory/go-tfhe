# FFT Benchmarking Suite

This directory contains comprehensive benchmarking tools for the FFT implementation.

## Quick Start

### Run Benchmarks

```bash
# Quick benchmark (1s per benchmark)
./benchmark.sh -q

# Full benchmark (2s per benchmark, more accurate)
./benchmark.sh -f

# With CPU profiling
./benchmark.sh -p

# Full benchmark with profiling and save as baseline
./benchmark.sh -f -p -s
```

### Compare Performance

```bash
# 1. Create baseline before optimization
./benchmark.sh -f -s

# 2. Make your optimizations
# ... edit fft.go ...

# 3. Compare with baseline
./benchmark.sh -f -c
```

## Files

### Benchmark Code
- **`fft_test.go`**: Complete test and benchmark suite
  - Correctness tests (TestFFTRoundtrip, TestPolyMul1024, etc.)
  - Performance benchmarks (see below)

### Documentation
- **`BENCHMARK_ANALYSIS.md`**: Analysis of current performance
  - Performance metrics table
  - Key findings and strengths
  - Optimization opportunities ranked by priority
  
- **`OPTIMIZATION_GUIDE.md`**: Detailed optimization guide
  - CPU profile analysis with time breakdown
  - Prioritized optimization recommendations
  - Code examples for each optimization
  - Estimated impact and effort for each

- **`BENCHMARKING_README.md`**: This file

### Scripts
- **`benchmark.sh`**: Automated benchmark runner
  - Quick and full benchmark modes
  - CPU and memory profiling
  - Baseline comparison
  - Colorized output with key metrics

## Available Benchmarks

### Core Operations
| Benchmark | Description | Typical Time | Allocations |
|-----------|-------------|--------------|-------------|
| `BenchmarkIFFT1024` | Time → Frequency domain | ~3.2 µs | 0 |
| `BenchmarkFFT1024` | Frequency → Time domain | ~3.2 µs | 0 |
| `BenchmarkPolyMul1024` | Polynomial multiplication | ~10.8 µs | 0 |
| `BenchmarkFFTRoundtrip` | IFFT + FFT roundtrip | ~6.5 µs | 0 |

### Initialization
| Benchmark | Description | Typical Time | Allocations |
|-----------|-------------|--------------|-------------|
| `BenchmarkNewFFTProcessor` | Processor creation | ~10.5 µs | 62,720 B, 9 allocs |

### Batch Operations
| Benchmark | Description | Notes |
|-----------|-------------|-------|
| `BenchmarkBatchIFFT1024` | Batch IFFT (1-32 polys) | Tests different batch sizes |
| `BenchmarkBatchFFT1024` | Batch FFT (1-32 polys) | Tests different batch sizes |

### Input Pattern Analysis
| Benchmark | Description | Purpose |
|-----------|-------------|---------|
| `BenchmarkIFFT1024WithPatterns` | IFFT with zeros/delta/sequential/random | Detect worst-case inputs |
| `BenchmarkFFT1024WithPatterns` | FFT with zeros/sequential/random | Detect worst-case inputs |

### API Comparison
| Benchmark | Description | Purpose |
|-----------|-------------|---------|
| `BenchmarkSliceVsArray` | Compare array vs slice API | Measure overhead of slice-based API |
| `BenchmarkPolyMulWithDifferentMagnitudes` | Small/medium/large values | Check magnitude dependence |

### Memory Management
| Benchmark | Description | Purpose |
|-----------|-------------|---------|
| `BenchmarkMemoryAllocation` | Processor reuse vs recreation | Show importance of reusing processors |

## Running Specific Benchmarks

```bash
# Run only IFFT benchmarks
go test -bench=IFFT -benchmem

# Run only PolyMul benchmarks
go test -bench=PolyMul -benchmem

# Run with longer benchtime for more accuracy
go test -bench=BenchmarkPolyMul1024 -benchtime=5s -benchmem

# Run benchmarks and save results
go test -bench=. -benchmem > results.txt
```

## Profiling

### CPU Profiling

```bash
# Generate CPU profile
./benchmark.sh -p

# Or manually:
go test -bench=BenchmarkPolyMul1024 -cpuprofile=cpu.prof -benchtime=2s

# Analyze interactively
go tool pprof cpu.prof
# Commands: top, list <function>, web (requires graphviz)

# Show top functions
go tool pprof -top cpu.prof

# Show specific function
go tool pprof -list=IFFT1024 cpu.prof
```

### Memory Profiling

```bash
# Generate memory profile
./benchmark.sh -m

# Or manually:
go test -bench=BenchmarkPolyMul1024 -memprofile=mem.prof -benchtime=2s

# Analyze
go tool pprof -top mem.prof
```

### Combined Profiling

```bash
# CPU and memory profiling together
./benchmark.sh -f -p -m
```

## Comparing Implementations

### Using benchstat (Recommended)

Install benchstat:
```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Compare two implementations:
```bash
# Before optimization
git checkout main
./benchmark.sh -f -s
mv benchmark_results.txt old.txt

# After optimization  
git checkout my-optimization-branch
./benchmark.sh -f
mv benchmark_results.txt new.txt

# Compare
benchstat old.txt new.txt
```

Example output:
```
name             old time/op  new time/op  delta
PolyMul1024-12   10.8µs ± 2%  9.5µs ± 1%  -12.04%  (p=0.000 n=10+10)

name             old alloc/op  new alloc/op  delta
PolyMul1024-12    0.00B         0.00B          ~     (all equal)
```

### Manual Comparison

```bash
# Before
./benchmark.sh -f
cp benchmark_results.txt baseline.txt

# After changes
./benchmark.sh -f -c  # Compares with baseline.txt
```

## Interpreting Results

### Good Signs ✓
- **Zero allocations** in hot path (IFFT1024, FFT1024, PolyMul1024)
- **Consistent timing** across input patterns (zeros, random, etc.)
- **Linear scaling** for batch operations
- **Low variance** in benchmark results (± 1-2%)

### Warning Signs ⚠
- Allocations appearing in hot path functions
- Large variance in timing (± 10%+)
- Non-linear scaling for batch operations
- Performance degradation with certain input patterns

### Example: Successful Optimization

```
Before:
BenchmarkPolyMul1024-12    100000    12000 ns/op    0 B/op    0 allocs/op

After:
BenchmarkPolyMul1024-12    120000    10000 ns/op    0 B/op    0 allocs/op
                           +20%      -16.7%         ✓ Still zero allocs
```

### Example: Problematic Optimization

```
Before:
BenchmarkPolyMul1024-12    100000    12000 ns/op     0 B/op    0 allocs/op

After:
BenchmarkPolyMul1024-12    150000    11000 ns/op    4096 B/op   1 allocs/op
                           +50%      -8.3%          ❌ New allocation!
```

**Why is this problematic?** The 8.3% speedup is overshadowed by the new allocation, which will:
- Increase GC pressure
- Reduce performance in tight loops
- May cause performance to degrade over time

## Advanced Profiling

### Flame Graphs

Generate flame graph (requires go-torch):
```bash
go test -bench=BenchmarkPolyMul1024 -cpuprofile=cpu.prof -benchtime=5s
go-torch cpu.prof
```

### Assembly Inspection

View generated assembly:
```bash
go build -gcflags="-S" ./fft > asm.txt 2>&1
grep -A 50 "IFFT1024" asm.txt
```

### Cache Analysis (Linux only)

```bash
perf stat -e cache-references,cache-misses,cycles,instructions \
  go test -bench=BenchmarkPolyMul1024
```

### Continuous Benchmarking

For CI/CD integration:
```bash
# In your CI script
./benchmark.sh -f > results.txt

# Compare with main branch
git checkout main
./benchmark.sh -f > baseline.txt

benchstat baseline.txt results.txt | tee comparison.txt

# Fail if performance degrades by >5%
# (implement in CI script)
```

## Tips for Optimization

1. **Always profile first**: Don't guess where the bottleneck is
2. **Measure everything**: Run benchmarks before and after every change
3. **One change at a time**: Make isolated changes so you know what helped
4. **Check correctness**: Run `go test` after every optimization
5. **Watch for regressions**: Ensure zero allocations are maintained
6. **Consider the cost**: Don't sacrifice readability for 1% speedup

## Benchmark Stability

For accurate results:
- Close other applications
- Disable CPU frequency scaling (if possible)
- Run with `-benchtime=2s` or higher for stability
- Run multiple times and compare
- Watch out for thermal throttling on laptops

## Current Performance Baseline (Apple M3 Pro)

As of the last benchmark run:
- **IFFT1024**: 3.2 µs (0 allocs)
- **FFT1024**: 3.2 µs (0 allocs)
- **PolyMul1024**: 10.8 µs (0 allocs)
- **FFT Roundtrip**: 6.5 µs (0 allocs)

These numbers represent **highly optimized** performance. Any optimization should be carefully measured and validated.

## Getting Help

- See `OPTIMIZATION_GUIDE.md` for specific optimization recommendations
- See `BENCHMARK_ANALYSIS.md` for analysis of current performance
- Check the CPU profile: `./benchmark.sh -p`
- Ask: "Is this optimization worth the added complexity?"

## Summary

This benchmarking suite provides everything needed to:
1. ✅ Measure current performance
2. ✅ Identify bottlenecks
3. ✅ Compare implementations
4. ✅ Track improvements over time
5. ✅ Prevent performance regressions

Use it to guide optimization efforts and maintain high performance.

