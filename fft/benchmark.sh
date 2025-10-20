#!/bin/bash

# FFT Benchmark Script
# Usage: ./benchmark.sh [options]
# Options:
#   -f, --full       Run full benchmark suite
#   -q, --quick      Run quick benchmark (default)
#   -p, --profile    Run with CPU profiling
#   -m, --mem        Run with memory profiling
#   -c, --compare    Compare with baseline (requires baseline.txt)
#   -s, --save       Save results as new baseline
#   -h, --help       Show this help message

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default options
BENCHTIME="1s"
PROFILE_CPU=false
PROFILE_MEM=false
COMPARE=false
SAVE_BASELINE=false
BENCH_FILTER="."

print_help() {
    echo "FFT Benchmark Script"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -f, --full       Run full benchmark suite (2s per benchmark)"
    echo "  -q, --quick      Run quick benchmark (1s per benchmark, default)"
    echo "  -p, --profile    Run with CPU profiling"
    echo "  -m, --mem        Run with memory profiling"
    echo "  -c, --compare    Compare with baseline (requires baseline.txt)"
    echo "  -s, --save       Save results as new baseline"
    echo "  -h, --help       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -q                  # Quick benchmark"
    echo "  $0 -f -s               # Full benchmark and save as baseline"
    echo "  $0 -q -c               # Quick benchmark and compare with baseline"
    echo "  $0 -p                  # Run with CPU profiling"
    echo "  $0 -f -p -m            # Full benchmark with CPU and memory profiling"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--full)
            BENCHTIME="2s"
            shift
            ;;
        -q|--quick)
            BENCHTIME="1s"
            shift
            ;;
        -p|--profile)
            PROFILE_CPU=true
            shift
            ;;
        -m|--mem)
            PROFILE_MEM=true
            shift
            ;;
        -c|--compare)
            COMPARE=true
            shift
            ;;
        -s|--save)
            SAVE_BASELINE=true
            shift
            ;;
        -h|--help)
            print_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            print_help
            exit 1
            ;;
    esac
done

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  FFT Benchmark Suite${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Build test binary
echo -e "${YELLOW}Building test binary...${NC}"
go test -c -o fft.test

# Run benchmarks
BENCH_ARGS="-test.bench=${BENCH_FILTER} -test.benchmem -test.benchtime=${BENCHTIME}"

if [ "$PROFILE_CPU" = true ]; then
    BENCH_ARGS="${BENCH_ARGS} -test.cpuprofile=cpu.prof"
fi

if [ "$PROFILE_MEM" = true ]; then
    BENCH_ARGS="${BENCH_ARGS} -test.memprofile=mem.prof"
fi

echo -e "${YELLOW}Running benchmarks (benchtime=${BENCHTIME})...${NC}"
echo ""

if [ "$SAVE_BASELINE" = true ]; then
    ./fft.test ${BENCH_ARGS} | tee benchmark_results.txt
    cp benchmark_results.txt baseline.txt
    echo ""
    echo -e "${GREEN}✓ Results saved to baseline.txt${NC}"
else
    ./fft.test ${BENCH_ARGS} | tee benchmark_results.txt
fi

echo ""

# Compare with baseline if requested
if [ "$COMPARE" = true ]; then
    if [ ! -f baseline.txt ]; then
        echo -e "${RED}✗ No baseline.txt found. Run with -s to create one.${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Comparing with baseline...${NC}"
    echo ""
    
    # Use benchstat if available, otherwise simple comparison
    if command -v benchstat &> /dev/null; then
        benchstat baseline.txt benchmark_results.txt
    else
        echo -e "${YELLOW}Note: Install benchstat for detailed comparison:${NC}"
        echo -e "${YELLOW}  go install golang.org/x/perf/cmd/benchstat@latest${NC}"
        echo ""
        echo "Baseline:"
        grep "^Benchmark" baseline.txt | head -5
        echo ""
        echo "Current:"
        grep "^Benchmark" benchmark_results.txt | head -5
    fi
fi

# Show profile results if profiling was enabled
if [ "$PROFILE_CPU" = true ]; then
    echo ""
    echo -e "${YELLOW}CPU Profile Summary:${NC}"
    echo ""
    go tool pprof -top -nodecount=10 cpu.prof
    echo ""
    echo -e "${GREEN}✓ CPU profile saved to cpu.prof${NC}"
    echo -e "  View with: ${BLUE}go tool pprof cpu.prof${NC}"
fi

if [ "$PROFILE_MEM" = true ]; then
    echo ""
    echo -e "${YELLOW}Memory Profile Summary:${NC}"
    echo ""
    go tool pprof -top -nodecount=10 mem.prof
    echo ""
    echo -e "${GREEN}✓ Memory profile saved to mem.prof${NC}"
    echo -e "  View with: ${BLUE}go tool pprof mem.prof${NC}"
fi

# Summary
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Benchmark Complete${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Extract key metrics
echo -e "${YELLOW}Key Metrics:${NC}"
echo ""

IFFT_TIME=$(grep "^BenchmarkIFFT1024-" benchmark_results.txt | awk '{print $3}')
FFT_TIME=$(grep "^BenchmarkFFT1024-" benchmark_results.txt | awk '{print $3}')
POLYMUL_TIME=$(grep "^BenchmarkPolyMul1024-" benchmark_results.txt | awk '{print $3}')
ROUNDTRIP_TIME=$(grep "^BenchmarkFFTRoundtrip-" benchmark_results.txt | awk '{print $3}')

if [ -n "$IFFT_TIME" ]; then
    echo -e "  IFFT1024:        ${GREEN}${IFFT_TIME}${NC}"
fi

if [ -n "$FFT_TIME" ]; then
    echo -e "  FFT1024:         ${GREEN}${FFT_TIME}${NC}"
fi

if [ -n "$POLYMUL_TIME" ]; then
    echo -e "  PolyMul1024:     ${GREEN}${POLYMUL_TIME}${NC}"
fi

if [ -n "$ROUNDTRIP_TIME" ]; then
    echo -e "  FFT Roundtrip:   ${GREEN}${ROUNDTRIP_TIME}${NC}"
fi

echo ""

# Check for allocations in hot path
ZERO_ALLOC=$(grep "^BenchmarkIFFT1024-\|^BenchmarkFFT1024-\|^BenchmarkPolyMul1024-" benchmark_results.txt | grep "0 allocs/op" | wc -l | tr -d ' ')
if [ "$ZERO_ALLOC" = "3" ]; then
    echo -e "${GREEN}✓ Zero allocations in hot path!${NC}"
else
    echo -e "${RED}⚠ Warning: Allocations detected in hot path${NC}"
fi

echo ""

# Cleanup
rm -f fft.test

echo -e "${GREEN}Done!${NC}"




