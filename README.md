# go-tfhe: Fully Homomorphic Encryption Library
A pure Go implementation of TFHE (Torus Fully Homomorphic Encryption) for gophers

<img width="400" height="400" alt="image" src="https://github.com/user-attachments/assets/0f22ddfc-26b7-457b-8ed6-6be89e0f5256" />

## Overview

Go-TFHE is a library for performing homomorphic operations on encrypted data. It allows you to compute on encrypted data without decrypting it, enabling privacy-preserving computation in the cloud.

> Not the language you were looking for? Check out our [rust](https://github.com/thedonutfactory/rs-tfhe) or [zig](https://github.com/thedonutfactory/zig-tfhe) sister projects

### Features

- **Multiple Parameter Profiles**: 80-bit, 110-bit, 128-bit security + Uint5 for arithmetic
- **Homomorphic Gates**: AND, OR, NAND, NOR, XOR, XNOR, NOT, MUX
- **Programmable Bootstrapping**: Evaluate arbitrary functions during bootstrapping
- **Fast Arithmetic**: 4-bootstrap nibble addition with messageModulus=32
- **N=2048 Support**: Full parity with tfhe-go reference implementation
- **Batch Operations**: Parallel processing for multiple gates
- **Optimized FFT**: Ported from tfhe-go for best performance
- **Pure Go**: No C dependencies, easy to build and deploy
- **Concurrent**: Leverages Go's goroutines for parallelization

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

## Security Levels and Parameter Profiles

Go-TFHE supports multiple parameter profiles optimized for different use cases:

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

### Uint5 Parameters - Fast Multi-Bit Arithmetic ⭐ NEW!

```go
params.CurrentSecurityLevel = params.SecurityUint5
```

- **N (LWE/Poly dimension)**: 1071/2048
- **ALPHA (noise)**: 7.1e-08 / 2.2e-17 (~700x lower noise!)
- **messageModulus**: Up to **32** (5-bit message space)
- **Polynomial degree**: **2048** (doubled)
- **Use case**: Fast multi-bit arithmetic, homomorphic addition/multiplication
- **Performance**: **~230ms for 8-bit addition** (only 4 bootstraps!)
- **Key generation**: ~5-6 seconds (slower than standard params)
- **Security**: Comparable to 80-bit, optimized for precision over maximum hardness

**Perfect for**: Arithmetic circuits, financial calculations, machine learning inference

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

## Programmable Bootstrapping

Programmable bootstrapping is an advanced feature that allows you to **evaluate arbitrary functions on encrypted data** during the bootstrapping process. This combines noise refreshing with function evaluation in a single operation.

### What is Programmable Bootstrapping?

Traditional bootstrapping refreshes a ciphertext's noise but keeps the encrypted value unchanged. Programmable bootstrapping goes further: it applies a function `f` to the encrypted value while refreshing the noise.

If you have an encryption of `x`, programmable bootstrapping gives you an encryption of `f(x)`.

### Basic Usage

```go
import (
    "github.com/thedonutfactory/go-tfhe/cloudkey"
    "github.com/thedonutfactory/go-tfhe/evaluator"
    "github.com/thedonutfactory/go-tfhe/key"
    "github.com/thedonutfactory/go-tfhe/params"
    "github.com/thedonutfactory/go-tfhe/tlwe"
)

// Generate keys
secretKey := key.NewSecretKey()
cloudKey := cloudkey.NewCloudKey(secretKey)
eval := evaluator.NewEvaluator(params.GetTRGSWLv1().N)

// Encrypt a message using LWE message encoding
// Note: Use EncryptLWEMessage (not EncryptBool) for programmable bootstrapping
ct := tlwe.NewTLWELv0()
ct.EncryptLWEMessage(1, 2, params.GetTLWELv0().ALPHA, secretKey.KeyLv0) // message 1 (true)

// Define a function to apply (e.g., NOT)
notFunc := func(x int) int { return 1 - x }

// Apply the function during bootstrapping
result := eval.BootstrapFunc(
    ct,
    notFunc,
    2, // message modulus (2 for binary)
    cloudKey.BootstrappingKey,
    cloudKey.KeySwitchingKey,
    cloudKey.DecompositionOffset,
)

// Decrypt result using LWE message decoding
output := result.DecryptLWEMessage(2, secretKey.KeyLv0) // 0 (false)
```

**Important:** Programmable bootstrapping uses general LWE message encoding (`message * scale`), not binary boolean encoding (±1/8). Always use:
- `EncryptLWEMessage()` to encrypt messages
- `DecryptLWEMessage()` to decrypt results

### Lookup Table (LUT) Reuse

For better performance when applying the same function multiple times, pre-compute the lookup table:

```go
import "github.com/thedonutfactory/go-tfhe/lut"

// Create a lookup table generator
gen := lut.NewGenerator(2) // 2 = binary messages

// Pre-compute the lookup table once
notFunc := func(x int) int { return 1 - x }
lookupTable := gen.GenLookUpTable(notFunc)

// Reuse the LUT for multiple operations
for _, ct := range ciphertexts {
    result := eval.BootstrapLUT(
        ct,
        lookupTable,
        cloudKey.BootstrappingKey,
        cloudKey.KeySwitchingKey,
        cloudKey.DecompositionOffset,
    )
    // Process result...
}
```

### Supported Functions

You can evaluate **any** function `f: {0, 1, ..., m-1} → {0, 1, ..., m-1}` where `m` is the message modulus.

**Examples:**

```go
// Identity (refresh noise without changing value)
identity := func(x int) int { return x }

// NOT (boolean negation)
not := func(x int) int { return 1 - x }

// Constant functions
alwaysTrue := func(x int) int { return 1 }
alwaysFalse := func(x int) int { return 0 }

// Multi-bit functions (with message modulus = 4)
gen := lut.NewGenerator(4)
increment := func(x int) int { return (x + 1) % 4 }
double := func(x int) int { return (2 * x) % 4 }
```

### Use Cases

1. **Noise Refresh with Transformation**: Apply a function while cleaning up noise
2. **Efficient NOT gates**: Faster than traditional NOT + bootstrap
3. **Lookup Table Evaluation**: Implement truth tables directly
4. **Multi-bit Operations**: Work with values beyond binary
5. **Custom Boolean Functions**: Implement any boolean function efficiently

### Performance Comparison

| Operation | Traditional | Programmable Bootstrap | Speedup |
|-----------|-------------|------------------------|---------|
| NOT + Bootstrap | 2 operations | 1 operation | 2x |
| Lookup Table (precomputed) | - | Single bootstrap | - |
| Function + Noise Refresh | 2 operations | 1 operation | 2x |

### Advanced: Custom Message Moduli

```go
// Work with 3-bit values (8 possible messages)
gen := lut.NewGenerator(8)

// Define a function operating on 0-7
customFunc := func(x int) int {
    // Apply any transformation
    return (x * 3 + 2) % 8
}

lookupTable := gen.GenLookUpTable(customFunc)
```

### Example: Complete Demo

See the complete working example in `examples/programmable_bootstrap/`:

```bash
cd examples/programmable_bootstrap
go run main.go
```

This example demonstrates:
- Identity function
- NOT function  
- Constant functions
- LUT reuse for efficiency
- Multi-bit message support

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
├── lut/          # Lookup tables for programmable bootstrapping
├── evaluator/    # Zero-allocation evaluator for TFHE operations
├── key/          # Key generation and management
├── gates/        # Homomorphic gate operations
└── examples/     # Example applications
```

### Key Algorithms

1. **TLWE/TRLWE Encryption**: Torus-based Learning With Errors
2. **Blind Rotation**: Core bootstrapping operation using TRGSW
3. **Programmable Bootstrapping**: Evaluate arbitrary functions during bootstrapping
4. **Key Switching**: Convert between different key spaces
5. **Gadget Decomposition**: Break down ciphertexts for external product
6. **FFT-based Polynomial Multiplication**: Efficient negacyclic convolution
7. **Lookup Table Generation**: Encode functions as test vectors

## Performance

Performance characteristics on a typical modern CPU:

| Operation | Time (128-bit) | Time (80-bit) |
|-----------|----------------|---------------|
| Key Generation | ~5-10s | ~3-5s |
| Single Gate | ~100-150ms | ~60-80ms |
| Batch (8 gates) | ~200-300ms | ~120-180ms |
| Addition (8-bit) | ~8-12s | ~5-7s |

*Note: Performance can vary based on CPU architecture and number of cores.*

## Examples

See the `examples/` directory for complete working examples:

- `add_two_numbers/` - Homomorphic addition of two 16-bit numbers
- `simple_gates/` - Test all available homomorphic gates
- `programmable_bootstrap/` - Demonstrate programmable bootstrapping with various functions

Run examples:

```bash
cd examples/add_two_numbers
go run main.go

cd examples/simple_gates
go run main.go

cd examples/programmable_bootstrap
go run main.go
```

## Key Advantages

- **Pure Go**: No C dependencies, no build tools required
- **Easy Deployment**: Single binary, cross-platform compilation
- **Simple Integration**: Standard Go modules, no CGO
- **Parallelization**: Built-in concurrency with goroutines
- **Multiple Security Levels**: 80/110/128-bit security parameters

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

- **FFT Performance**: Uses standard Go FFT library (future: custom SIMD-optimized FFT)
- **Memory Usage**: Go's garbage collector trades memory for convenience

## Future Improvements

- [ ] Add SIMD-optimized FFT using Go assembly
- [ ] Implement custom FFT with better cache locality  
- [ ] Add GPU acceleration support (Metal/CUDA)
- [ ] Optimize memory allocations and reduce GC pressure
- [ ] Add more example circuits (multiplication, comparison, sorting, etc.)
- [ ] Support for wider data types (16-bit, 32-bit operations)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

## References

- Original TFHE paper: [TFHE: Fast Fully Homomorphic Encryption over the Torus](https://eprint.iacr.org/2018/421)
- Extended FFT paper: [Fast and Error-Free Negacyclic Integer Convolution](https://eprint.iacr.org/2021/480)

