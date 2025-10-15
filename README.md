# Go-TFHE: Fully Homomorphic Encryption Library in Go

A pure Go implementation of TFHE (Torus Fully Homomorphic Encryption)

## Overview

Go-TFHE is a library for performing homomorphic operations on encrypted data. It allows you to compute on encrypted data without decrypting it, enabling privacy-preserving computation in the cloud.

### Features

- ✅ **Multiple Security Levels**: Choose between 80-bit, 110-bit, or 128-bit security
- ✅ **Homomorphic Gates**: AND, OR, NAND, NOR, XOR, XNOR, NOT, MUX
- ✅ **Batch Operations**: Parallel processing for multiple gates
- ✅ **Bootstrapping**: Noise reduction using blind rotation
- ✅ **Pure Go**: No C dependencies, easy to build and deploy
- ✅ **Concurrent**: Leverages Go's goroutines for parallelization

## Installation

```bash
go get github.com/thedonutfactory/go-tfhe
```

## Quick Start

### Simple Example: Homomorphic AND

```go
package main

import (
    "fmt"
    "github.com/thedonutfactory/go-tfhe/gates"
    "github.com/thedonutfactory/go-tfhe/key"
    "github.com/thedonutfactory/go-tfhe/params"
    "github.com/thedonutfactory/go-tfhe/tlwe"
)

func main() {
    // Generate keys
    secretKey := key.NewSecretKey()
    cloudKey := key.NewCloudKey(secretKey)

    // Encrypt inputs
    a := true
    b := false
    ctA := tlwe.NewTLWELv0().EncryptBool(a, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
    ctB := tlwe.NewTLWELv0().EncryptBool(b, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)

    // Compute homomorphic AND
    ctResult := gates.AND(ctA, ctB, cloudKey)

    // Decrypt result
    result := ctResult.DecryptBool(secretKey.KeyLv0)
    fmt.Printf("%v AND %v = %v\n", a, b, result) // Output: true AND false = false
}
```

### Example: Homomorphic Addition

```go
package main

import (
    "fmt"
    "github.com/thedonutfactory/go-tfhe/bitutils"
    "github.com/thedonutfactory/go-tfhe/gates"
    "github.com/thedonutfactory/go-tfhe/key"
)

func FullAdder(serverKey *key.CloudKey, a, b, cin *gates.Ciphertext) (*gates.Ciphertext, *gates.Ciphertext) {
    aXorB := gates.XOR(a, b, serverKey)
    aAndB := gates.AND(a, b, serverKey)
    aXorBAndC := gates.AND(aXorB, cin, serverKey)
    
    sum := gates.XOR(aXorB, cin, serverKey)
    carry := gates.OR(aAndB, aXorBAndC, serverKey)
    
    return sum, carry
}

func main() {
    secretKey := key.NewSecretKey()
    cloudKey := key.NewCloudKey(secretKey)

    // Encrypt two 8-bit numbers
    a := uint8(42)
    b := uint8(17)
    
    aBits := bitutils.U8ToBits(a)
    bBits := bitutils.U8ToBits(b)
    
    ctA := bitutils.EncryptBits(aBits, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
    ctB := bitutils.EncryptBits(bBits, params.GetTLWELv0().ALPHA, secretKey.KeyLv0)
    
    // Homomorphic addition
    result := make([]*gates.Ciphertext, 8)
    carry := gates.Constant(false)
    for i := 0; i < 8; i++ {
        result[i], carry = FullAdder(cloudKey, ctA[i], ctB[i], carry)
    }
    
    // Decrypt result
    resultBits := bitutils.DecryptBits(result, secretKey.KeyLv0)
    sum := bitutils.ConvertU8(resultBits)
    
    fmt.Printf("%d + %d = %d\n", a, b, sum) // Output: 42 + 17 = 59
}
```

## Security Levels

Go-TFHE supports three security levels:

### 128-bit Security (Default) - Recommended for Production

```go
params.CurrentSecurityLevel = params.Security128Bit
```

- **N (LWE dimension)**: 700/1024
- **ALPHA (noise)**: 2.0e-5 / 2.0e-8
- **Use case**: Production systems, high-security applications
- **Performance**: ~100-150ms per gate (pure Go implementation)

### 110-bit Security - Balanced

```go
params.CurrentSecurityLevel = params.Security110Bit
```

- **N (LWE dimension)**: 630/1024
- **ALPHA (noise)**: 3.05e-5 / 2.98e-8
- **Use case**: General purpose, original TFHE parameters
- **Performance**: ~20-30% faster than 128-bit

### 80-bit Security - Development/Testing

```go
params.CurrentSecurityLevel = params.Security80Bit
```

- **N (LWE dimension)**: 550/1024
- **ALPHA (noise)**: 5.0e-5 / 3.73e-8
- **Use case**: Development, testing, prototyping
- **Performance**: ~30-40% faster than 128-bit
- **Warning**: Not recommended for production

## Available Gates

### Basic Gates
- `AND(a, b, key)` - Homomorphic AND
- `OR(a, b, key)` - Homomorphic OR
- `NAND(a, b, key)` - Homomorphic NAND
- `NOR(a, b, key)` - Homomorphic NOR
- `XOR(a, b, key)` - Homomorphic XOR
- `XNOR(a, b, key)` - Homomorphic XNOR
- `NOT(a)` - Homomorphic NOT (no bootstrapping needed)

### Advanced Gates
- `MUX(a, b, c, key)` - Homomorphic multiplexer (a ? b : c)
- `ANDNY(a, b, key)` - NOT(a) AND b
- `ANDYN(a, b, key)` - a AND NOT(b)
- `ORNY(a, b, key)` - NOT(a) OR b
- `ORYN(a, b, key)` - a OR NOT(b)

### Batch Operations

Process multiple gates in parallel for better performance:

```go
// Prepare inputs
inputs := [][2]*gates.Ciphertext{
    {ctA1, ctB1},
    {ctA2, ctB2},
    {ctA3, ctB3},
    {ctA4, ctB4},
}

// Batch AND (4 gates computed in parallel)
results := gates.BatchAND(inputs, cloudKey)

// Also available: BatchOR, BatchNAND, BatchNOR, BatchXOR, BatchXNOR
```

Expected speedup: 4-8x on multi-core systems.

## Architecture

### Core Components

```
go-tfhe/
├── params/       # Security parameters for different levels
├── utils/        # Utility functions (torus conversions, etc.)
├── bitutils/     # Bit manipulation and conversion
├── tlwe/         # TLWE (Torus Learning With Errors) encryption
├── trlwe/        # TRLWE (Ring variant of TLWE)
├── trgsw/        # TRGSW (GSW-based encryption) with FFT
├── fft/          # FFT operations for polynomial multiplication
├── key/          # Key generation and management
├── gates/        # Homomorphic gate operations
└── examples/     # Example applications
```

### Key Algorithms

1. **TLWE/TRLWE Encryption**: Torus-based Learning With Errors
2. **Blind Rotation**: Core bootstrapping operation using TRGSW
3. **Key Switching**: Convert between different key spaces
4. **Gadget Decomposition**: Break down ciphertexts for external product
5. **FFT-based Polynomial Multiplication**: Efficient negacyclic convolution

## Performance

Performance characteristics on a typical modern CPU:

| Operation | Time (128-bit) | Time (80-bit) |
|-----------|----------------|---------------|
| Key Generation | ~5-10s | ~3-5s |
| Single Gate | ~100-150ms | ~60-80ms |
| Batch (8 gates) | ~200-300ms | ~120-180ms |
| Addition (8-bit) | ~8-12s | ~5-7s |

*Note: Times are for pure Go implementation. The Rust version with hand-optimized assembly is ~3-5x faster.*

## Examples

See the `examples/` directory for complete working examples:

- `add_two_numbers/` - Homomorphic addition of two 16-bit numbers
- `simple_gates/` - Test all available homomorphic gates

Run examples:

```bash
cd examples/add_two_numbers
go run main.go

cd examples/simple_gates
go run main.go
```

## Comparison with Rust Implementation

| Feature | Go Implementation | Rust Implementation |
|---------|-------------------|---------------------|
| Pure Language | ✅ Yes | ❌ No (uses C++/ASM) |
| Easy Build | ✅ Yes | ⚠️ Requires build tools |
| Performance | ~100-150ms/gate | ~30-50ms/gate |
| Parallelization | ✅ Goroutines | ✅ Rayon |
| Security Levels | ✅ 80/110/128-bit | ✅ 80/110/128-bit |

## Building from Source

```bash
git clone https://github.com/thedonutfactory/go-tfhe
cd go-tfhe
go build ./...
```

## Testing

```bash
go test ./...
```

## Limitations

- **Performance**: Pure Go is slower than hand-optimized assembly in Rust version
- **FFT Implementation**: Uses standard Go FFT library (no SIMD optimizations)
- **Memory**: Higher memory usage compared to Rust due to GC overhead

## Future Improvements

- [ ] Add SIMD optimizations using Go assembly
- [ ] Implement custom FFT with better cache locality
- [ ] Add GPU acceleration support
- [ ] Optimize memory allocations
- [ ] Add more example circuits (multiplication, comparison, etc.)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Same license as the original Rust implementation.

## References

- Original TFHE paper: [TFHE: Fast Fully Homomorphic Encryption over the Torus](https://eprint.iacr.org/2018/421)
- Rust implementation: [rs-tfhe](https://github.com/thedonutfactory/rs-tfhe)

## Acknowledgments

This is a port of the Rust TFHE implementation. All credit for the original design and algorithms goes to the original authors.

