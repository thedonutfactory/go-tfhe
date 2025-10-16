# Examples Guide

This directory contains demonstrations of go-tfhe capabilities, from basic gates to advanced programmable bootstrapping.

## Available Examples

### 1. `simple_gates/` - Boolean Gate Demonstrations

**What it does**: Tests all basic homomorphic boolean gates

**Gates demonstrated**:
- AND, OR, NAND, NOR
- XOR, XNOR
- NOT

**Run**: `make run-gates`

**Time**: ~10 seconds (tests all gates on all input combinations)

**Best for**: Learning basic TFHE gates

---

### 2. `add_two_numbers/` - Traditional 8-bit Addition

**What it does**: 8-bit addition using traditional ripple-carry adder

**Method**: Bit-by-bit with boolean gates (XOR, AND, OR)

**Operations**: 40 gates (5 per bit × 8 bits)

**Run**: `make run-add`

**Time**: ~1.1 seconds

**Best for**: Understanding traditional TFHE approach and why PBS is revolutionary

**Example output**:
```
Computing: 42 + 137 = 179
Operations: 40 boolean gates
Time: ~1.1s
✅ SUCCESS!
```

---

### 3. `add_8bit_pbs/` - Fast 8-bit Addition with PBS ⭐

**What it does**: 8-bit addition using Programmable Bootstrapping (PBS)

**Method**: Nibble-based (processes 4 bits at once)

**Operations**: 3 programmable bootstraps

**Parameters**: `SecurityUint5` (messageModulus=32, N=2048)

**Run**: `make run-add-pbs`

**Time**: ~230ms

**Best for**: Seeing the dramatic PBS performance advantage

**Example output**:
```
Computing: 42 + 137 = 179
Input A:  42 = 0b0010_1010 (nibbles: high=2, low=10)
Input B: 137 = 0b1000_1001 (nibbles: high=8, low=9)

Steps:
  1. Encrypt nibbles (4 nibbles)
  2. Add low nibbles (homomorphic, no bootstrap)
  3. Bootstrap 1: Extract low sum (mod 16)
  4. Bootstrap 2: Extract carry bit
  5. Add high nibbles + carry (homomorphic)
  6. Bootstrap 3: Extract high sum (mod 16)
  7. Combine nibbles

Result: 179 = 0b1011_0011
Time: ~230ms (3 bootstraps)
✅ SUCCESS!

Speedup: 4.8x faster than traditional method! 🚀
```

---

### 4. `programmable_bootstrap/` - PBS Feature Demonstrations

**What it does**: Comprehensive PBS feature demonstrations

**Features shown**:
- Identity function (noise refresh)
- NOT function (bit flip)
- Constant functions
- LUT reuse
- Multi-bit messages (messageModulus=4)

**Run**: `make run-pbs`

**Time**: ~2-3 seconds (multiple demonstrations)

**Best for**: Learning PBS concepts and usage patterns

---

## Comparison Table

| Example | Method | Operations | Time | Speedup | Use Case |
|---------|--------|-----------|------|---------|----------|
| `simple_gates` | Boolean gates | Varies | ~10s | Baseline | Learn gates |
| `add_two_numbers` | Ripple-carry | 40 gates | ~1.1s | 1x | Traditional method |
| **`add_8bit_pbs`** | **PBS nibbles** | **3 PBS** | **~230ms** | **4.8x** ⭐ | **Fast arithmetic** |
| `programmable_bootstrap` | PBS | Varies | ~2-3s | N/A | Learn PBS |

## Quick Start

```bash
# 1. Start simple - learn the gates
make run-gates

# 2. See traditional approach
make run-add

# 3. See the PBS revolution!
make run-add-pbs

# 4. Explore PBS features
make run-pbs
```

## Understanding the Speedup

### Traditional Method (`add_two_numbers`)
```
Process: Bit-by-bit ripple carry
- For each bit (0-7):
  - XOR(a, b)      → 1 bootstrap
  - XOR(ab, c)     → 1 bootstrap  
  - AND(a, b)      → 1 bootstrap
  - AND(c, ab)     → 1 bootstrap
  - OR(...)        → 1 bootstrap
Total: 5 gates × 8 bits = 40 bootstraps
```

### PBS Method (`add_8bit_pbs`)
```
Process: Nibble-based with programmable bootstrapping
- Split into 4-bit chunks (nibbles)
- Add low nibbles → PBS extract sum & carry    (2 bootstraps)
- Add high nibbles → PBS extract sum           (1 bootstrap)
Total: 3 bootstraps

Why faster?
- Processes 4 bits at once instead of 1 bit
- LUTs encode multiple operations in single bootstrap
- 13x fewer bootstraps = massive speedup!
```

## Recommended Learning Path

1. **Start**: `simple_gates` - Understand basic operations
2. **Traditional**: `add_two_numbers` - See how addition works bit-by-bit
3. **Modern**: `add_8bit_pbs` - See the PBS advantage
4. **Deep Dive**: `programmable_bootstrap` - Explore PBS features

## Extending These Examples

### Build Your Own Operations

Using `add_8bit_pbs` as a template, you can create:

**8-bit Subtraction:**
```go
// Use LUT: f(x) = x % 16 for difference
// Handle borrow instead of carry
```

**8-bit Multiplication:**
```go
// Decompose into nibbles
// Use shift-and-add algorithm
// ~12-16 bootstraps
```

**8-bit Comparison:**
```go
// LUT: f(x) = x >= threshold ? 1 : 0
// Single bootstrap per nibble
```

## Parameter Selection for Examples

| Example | Parameter Used | Reason |
|---------|---------------|---------|
| `simple_gates` | `Security128Bit` | Binary operations |
| `add_two_numbers` | `Security128Bit` | Binary operations |
| `add_8bit_pbs` | **`SecurityUint5`** | Needs messageModulus=32 |
| `programmable_bootstrap` | `Security80Bit` | Faster demo |

## Performance Notes

All times are approximate and depend on hardware:
- Measured on: Modern CPU (2020+)
- Key generation: One-time cost
- Bootstrap times: Consistent per operation
- Can be parallelized for multiple operations

## Next Steps

After running the examples:
1. Read `PARAMETER_GUIDE.md` for parameter selection
2. Check `README.md` for API documentation
3. See `FINAL_STATUS.md` for complete library status
4. Build your own homomorphic applications!

---

**Start exploring: `make run-gates`** 🚀

