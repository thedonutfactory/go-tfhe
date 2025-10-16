# Uint Parameter Sets Status

## Production Ready ✅

| Parameter | messageModulus | Poly Degree | Status | Test Results |
|-----------|----------------|-------------|--------|--------------|
| **Uint2** | 4 | 512 | ✅ **READY** | 100% pass (Identity, Complement, Modulo) |
| **Uint3** | 8 | 1024 | ✅ **READY** | 100% pass (Identity, Complement, Modulo) |
| **Uint4** | 16 | 2048 | ✅ **READY** | 100% pass (Identity, Complement, Modulo) |
| **Uint5** | 32 | 2048 | ✅ **READY** | 100% pass (Identity, Complement, Modulo) |

## Experimental ⚠️

| Parameter | messageModulus | Poly Degree | LUTSize | Status | Test Results |
|-----------|----------------|-------------|---------|--------|--------------|
| Uint6 | 64 | 2048 | 4096 | ⚠️ **EXPERIMENTAL** | Identity ✅, Complement ❌, Modulo ❌ |
| Uint7 | 128 | 2048 | 8192 | ⚠️ **EXPERIMENTAL** | Partial failures |
| Uint8 | 256 | 2048 | 18432 | ⚠️ **EXPERIMENTAL** | Partial failures |

## Why Uint6-8 Are Experimental

Uint6-8 use **extended lookup tables** where `LookUpTableSize > PolyDegree`:
- Uint6: LookUpTableSize = 4096 = 2 × PolyDegree (polyExtendFactor = 2)
- Uint7: LookUpTableSize = 8192 = 4 × PolyDegree (polyExtendFactor = 4)
- Uint8: LookUpTableSize = 18432 = 9 × PolyDegree (polyExtendFactor = 9)

Our current LUT generation assumes `LookUpTableSize = PolyDegree`. Supporting extended LUTs requires:
1. Modified LUT generation algorithm with polyExtendFactor
2. Special blind rotation handling for extended LUTs
3. Additional testing and validation

## Recommendation

**For Production Use:**
- Use **Uint2-5** which are fully tested and reliable
- Uint5 supports messageModulus=32 which is sufficient for most applications
- For 8-bit values, use nibble-based decomposition with Uint5

**For Research/Development:**
- Uint6-8 can be explored for specific use cases
- Identity function works, suggesting basic PBS is functional
- More complex functions need additional work

## Workaround for Larger Values

Instead of Uint8 (0-255 direct), use Uint5 with byte decomposition:
```go
// Split 8-bit value into two 4-bit nibbles
low := value & 0x0F
high := (value >> 4) & 0x0F

// Encrypt with Uint5 (messageModulus=32)
// Process nibbles separately
// Combine with only 4 bootstraps!
```

This is actually **faster and more reliable** than direct Uint8!

## Future Work

To make Uint6-8 production-ready:
1. Implement extended LUT generation (polyExtendFactor > 1)
2. Update LUT generator to handle larger table sizes
3. Comprehensive testing of extended PBS
4. Performance optimization

For now, **Uint2-5 provide excellent coverage** for practical homomorphic arithmetic!
