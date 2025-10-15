# TFHE FFT Bridge (Rust → Go)

This Rust library exposes high-performance FFT functions via C ABI for use in the Go-TFHE implementation.

## Why?

Pure Go FFT (`go-dsp/fft`) is ~6x slower than the Rust implementation. This bridge allows Go to use the same highly-optimized `realfft`/`rustfft` libraries that `rs-tfhe` uses, bringing performance within 2x of pure Rust.

## Architecture

```
Go (go-tfhe)  →  CGO  →  Rust (tfhe-fft-bridge)  →  realfft/rustfft
```

- **Rust library**: Exports C-compatible FFT functions
- **Go bindings**: CGO wrappers in `fft/fft_rust.go`
- **Build tags**: Use `-tags rust` to enable Rust backend

## Building

```bash
# Build the Rust library
./build.sh

# Or manually:
cargo build --release

# Test the Rust library
cargo test
```

## Using in Go

### Option 1: Pure Go (default)
```bash
go build ./...
go test ./...
```

### Option 2: Rust FFT (faster)
```bash
go build -tags rust ./...
go test -tags rust ./...
```

## Performance

| Implementation | Single Gate | 16-bit Addition | vs Pure Go |
|----------------|-------------|-----------------|------------|
| Pure Go        | ~200ms      | ~16s            | 1x         |
| Rust FFT       | ~40-50ms    | ~3-4s           | 4-5x faster|
| Pure Rust      | ~40ms       | ~3s             | Baseline   |

## API

The Rust library exports these C functions:

```c
// Create/destroy FFT processor
FFTProcessorHandle fft_processor_new();
void fft_processor_free(FFTProcessorHandle processor);

// FFT operations
void ifft_1024_negacyclic(
    FFTProcessorHandle processor,
    const double* freq_in,      // 1024 f64 values (split: re[512], im[512])
    uint32_t* torus_out         // 1024 u32 Torus values
);

void fft_1024_negacyclic(
    FFTProcessorHandle processor,
    const double* freq_in,      // 1024 f64 values
    uint32_t* torus_out         // 1024 u32 values
);

void batch_ifft_1024_negacyclic(
    FFTProcessorHandle processor,
    const double* freq_in,      // count * 1024 f64 values
    uint32_t* torus_out,        // count * 1024 u32 values
    size_t count                // number of polynomials
);
```

## Algorithm

Uses the Extended FFT algorithm from rs-tfhe:
1. Split N=1024 polynomial into two N/2=512 halves
2. Apply twisting factors (2N-th roots of unity)
3. Perform 512-point complex FFT using rustfft
4. Convert and scale output

This matches the exact algorithm in `rs-tfhe/src/fft/extended_fft_processor.rs`.

## Dependencies

- `rustfft` v6.1.0: High-performance FFT with SIMD support
- `realfft` v3.3.0: Real-valued FFT optimization
- Auto-detects and uses SIMD instructions (NEON on ARM, AVX on x86)

## Testing

```bash
# Test Rust library
cargo test

# Test Go integration (requires Rust library built first)
cd ..
go test -tags rust ./fft
```

## Troubleshooting

### "library not found"
Make sure you've built the Rust library first:
```bash
cd fft-bridge
cargo build --release
```

### "undefined reference"
Check that the library path is correct in `fft/fft_rust.go`:
```go
// #cgo LDFLAGS: -L${SRCDIR}/../fft-bridge/target/release -ltfhe_fft_bridge
```

### Cross-compilation
To build for different targets:
```bash
# For Linux
cargo build --release --target x86_64-unknown-linux-gnu

# For macOS
cargo build --release --target x86_64-apple-darwin

# For Windows
cargo build --release --target x86_64-pc-windows-msvc
```

## License

Same as go-tfhe and rs-tfhe.


