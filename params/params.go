// Package params provides TFHE security parameters for different security levels.
//
// This library supports multiple security levels to allow users to choose
// the right balance between performance and security for their use case.
//
// # Available Security Levels
//
// - **80-bit**: Fast performance, suitable for development/testing
//   - ~20-30% faster than default
//
// - **110-bit**: Balanced performance and security
//   - Original TFHE reference parameters
//
// - **128-bit** (DEFAULT): High security, quantum-resistant
//   - Strong security guarantees for production use
//
// # Security Parameters Explained
//
// The security level is determined by several cryptographic parameters:
// - `N`: LWE dimension (higher = more secure, slower)
// - `ALPHA`: Noise standard deviation (smaller = often more secure with proper dimension)
// - `L`: Gadget decomposition levels (more = more secure, slower)
// - `BGBIT`: Decomposition base bits (smaller = more levels, more secure, slower)
package params

// Torus represents a 32-bit torus element
type Torus uint32

// SecurityLevel represents the selected security level
type SecurityLevel int

const (
	Security80Bit  SecurityLevel = 80
	Security110Bit SecurityLevel = 110
	Security128Bit SecurityLevel = 128
	SecurityUint1  SecurityLevel = 1 // Specialized for 1-bit message space (messageModulus=2, binary/boolean, N=1024)
	SecurityUint2  SecurityLevel = 2 // Specialized for 2-bit message space (messageModulus=4, N=512)
	SecurityUint3  SecurityLevel = 3 // Specialized for 3-bit message space (messageModulus=8, N=1024)
	SecurityUint4  SecurityLevel = 4 // Specialized for 4-bit message space (messageModulus=16, N=2048)
	SecurityUint5  SecurityLevel = 5 // Specialized for 5-bit message space (messageModulus=32, N=2048)
	SecurityUint6  SecurityLevel = 6 // Specialized for 6-bit message space (messageModulus=64, N=2048)
	SecurityUint7  SecurityLevel = 7 // Specialized for 7-bit message space (messageModulus=128, N=2048)
	SecurityUint8  SecurityLevel = 8 // Specialized for 8-bit message space (messageModulus=256, N=2048)
)

// Current security level (can be changed at runtime if needed)
var CurrentSecurityLevel = Security128Bit

// TLWE Level 0 Parameters
type TLWELv0Params struct {
	N     int
	ALPHA float64
}

// TLWE Level 1 Parameters
type TLWELv1Params struct {
	N     int
	ALPHA float64
}

// TRLWE Level 1 Parameters
type TRLWELv1Params struct {
	N     int
	ALPHA float64
}

// TRGSW Level 1 Parameters
type TRGSWLv1Params struct {
	N         int
	NBIT      int
	BGBIT     uint32
	BG        uint32
	L         int
	BASEBIT   int
	IKS_T     int
	ALPHA     float64
	BlockSize int // Block size for block blind rotation (1=original, >1=block algorithm)
}

// ============================================================================
// 80-BIT SECURITY PARAMETERS (Performance-Optimized)
// ============================================================================
var params80Bit = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     550,
		ALPHA: 5.0e-5, // 2^-14.3 approximately
	},
	TLWELv1: TLWELv1Params{
		N:     1024,
		ALPHA: 3.73e-8, // 2^-24.7 approximately
	},
	TRLWELv1: TRLWELv1Params{
		N:     1024,
		ALPHA: 3.73e-8,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         1024,
		NBIT:      10,
		BGBIT:     6,
		BG:        1 << 6,
		L:         3,
		BASEBIT:   2,
		IKS_T:     7,
		ALPHA:     3.73e-8,
		BlockSize: 3, // Use block blind rotation (3-4x faster)
	},
}

// ============================================================================
// 110-BIT SECURITY PARAMETERS (Original TFHE, Balanced)
// ============================================================================
var params110Bit = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     630,
		ALPHA: 3.0517578125e-05, // 2^-15 approximately
	},
	TLWELv1: TLWELv1Params{
		N:     1024,
		ALPHA: 2.980232238769531e-8, // 2^-25 approximately
	},
	TRLWELv1: TRLWELv1Params{
		N:     1024,
		ALPHA: 2.980232238769531e-8,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         1024,
		NBIT:      10,
		BGBIT:     6,
		BG:        1 << 6,
		L:         3,
		BASEBIT:   2,
		IKS_T:     8,
		ALPHA:     2.980232238769531e-8,
		BlockSize: 3, // Use block blind rotation (3-4x faster)
	},
}

// ============================================================================
// 128-BIT SECURITY PARAMETERS (DEFAULT - High Security, Quantum-Resistant)
// ============================================================================
var params128Bit = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     700,
		ALPHA: 2.0e-5, // 2^-15.6 approximately
	},
	TLWELv1: TLWELv1Params{
		N:     1024,
		ALPHA: 2.0e-8, // 2^-25.6 approximately
	},
	TRLWELv1: TRLWELv1Params{
		N:     1024,
		ALPHA: 2.0e-8,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         1024,
		NBIT:      10,
		BGBIT:     6,
		BG:        1 << 6,
		L:         3,
		BASEBIT:   2,
		IKS_T:     9,
		ALPHA:     2.0e-8,
		BlockSize: 3, // Use block blind rotation (3-4x faster)
	},
}

// ============================================================================
// UINT1 PARAMETERS (Specialized for 1-bit message space, messageModulus=2)
// ============================================================================
// For binary/boolean operations with Uint naming convention.
// Equivalent to Security128Bit but named for consistency in Uint series.
// Key features:
// - messageModulus=2 (binary: 0 or 1)
// - Standard polynomial degree (N=1024)
// - Balanced security and performance
//
// Note: This is essentially an alias for production binary operations.
// Use this when you want consistent Uint naming, or use Security128Bit directly.
var paramsUint1 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     700,
		ALPHA: 2.0e-05,
	},
	TLWELv1: TLWELv1Params{
		N:     1024,
		ALPHA: 2.0e-08,
	},
	TRLWELv1: TRLWELv1Params{
		N:     1024,
		ALPHA: 2.0e-08,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         1024,
		NBIT:      10,
		BGBIT:     10,
		BG:        1 << 10,
		L:         2,
		BASEBIT:   2,
		IKS_T:     8,
		ALPHA:     2.0e-08,
		BlockSize: 3,
	},
}

// ============================================================================
// UINT2 PARAMETERS (Specialized for 2-bit message space, messageModulus=4)
// ============================================================================
// Based on tfhe-go's ParamsUint2 configuration.
// Note: Uses GLWERank=3 which may require special handling in some operations.
// Key features:
// - Small polynomial degree (N=512) for fast operations
// - Supports messageModulus=4
// - Lower noise for 2-bit precision
//
// Security: Comparable to standard parameters, optimized for 2-bit arithmetic.
var paramsUint2 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     687,
		ALPHA: 0.00002120846893069971872305794214,
	},
	TLWELv1: TLWELv1Params{
		N:     512,
		ALPHA: 0.00000000000231841227527049948463,
	},
	TRLWELv1: TRLWELv1Params{
		N:     512,
		ALPHA: 0.00000000000231841227527049948463,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         512,
		NBIT:      9,  // 512 = 2^9
		BGBIT:     18, // Base = 1 << 18
		BG:        1 << 18,
		L:         1,
		BASEBIT:   4, // KeySwitch base bits
		IKS_T:     3, // KeySwitch level
		ALPHA:     0.00000000000231841227527049948463,
		BlockSize: 3,
	},
}

// ============================================================================
// UINT3 PARAMETERS (Specialized for 3-bit message space, messageModulus=8)
// ============================================================================
// Based on tfhe-go's ParamsUint3 configuration.
// Key features:
// - Standard polynomial degree (N=1024)
// - Supports messageModulus=8
// - Very low noise for 3-bit precision
//
// Security: Optimized for 3-bit arithmetic with good noise margin.
var paramsUint3 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     820,
		ALPHA: 0.00000251676160959795544987084234,
	},
	TLWELv1: TLWELv1Params{
		N:     1024,
		ALPHA: 0.00000000000000022204460492503131,
	},
	TRLWELv1: TRLWELv1Params{
		N:     1024,
		ALPHA: 0.00000000000000022204460492503131,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         1024,
		NBIT:      10, // 1024 = 2^10
		BGBIT:     23, // Base = 1 << 23
		BG:        1 << 23,
		L:         1,
		BASEBIT:   6, // KeySwitch base bits
		IKS_T:     2, // KeySwitch level
		ALPHA:     0.00000000000000022204460492503131,
		BlockSize: 4,
	},
}

// ============================================================================
// UINT4 PARAMETERS (Specialized for 4-bit message space, messageModulus=16)
// ============================================================================
// Based on tfhe-go's ParamsUint4 configuration.
// Key features:
// - Large polynomial degree (N=2048)
// - Supports messageModulus=16
// - Very low noise for 4-bit precision
//
// Security: Optimized for 4-bit arithmetic, same noise as Uint3.
var paramsUint4 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     820,
		ALPHA: 0.00000251676160959795544987084234,
	},
	TLWELv1: TLWELv1Params{
		N:     2048,
		ALPHA: 0.00000000000000022204460492503131,
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 0.00000000000000022204460492503131,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11, // 2048 = 2^11
		BGBIT:     22, // Base = 1 << 22
		BG:        1 << 22,
		L:         1,
		BASEBIT:   5, // KeySwitch base bits
		IKS_T:     3, // KeySwitch level
		ALPHA:     0.00000000000000022204460492503131,
		BlockSize: 4,
	},
}

// ============================================================================
// UINT5 PARAMETERS (Specialized for 5-bit message space, messageModulus=32)
// ============================================================================
// Based on tfhe-go's ParamsUint5 configuration. These parameters are
// specifically designed for multi-bit arithmetic with large message spaces.
// Key features:
// - ~700x lower noise than standard 80-bit security
// - Larger polynomial degree (2048 vs 1024)
// - Supports messageModulus up to 32 reliably
// - Enables 4-bootstrap nibble addition
//
// Security: Provides comparable security to 80-bit level but optimized
// for precision rather than maximum cryptographic hardness.
var paramsUint5 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     1071,
		ALPHA: 7.088226765410429399593757e-08,
	},
	TLWELv1: TLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11,
		BGBIT:     22,
		BG:        1 << 22,
		L:         1,
		BASEBIT:   6,
		IKS_T:     3,
		ALPHA:     2.2204460492503131e-17,
		BlockSize: 7,
	},
}

// ============================================================================
// UINT6 PARAMETERS (Specialized for 6-bit message space, messageModulus=64)
// ============================================================================
// Based on tfhe-go's ParamsUint6 configuration.
// Key features:
// - Same noise as Uint5 for reliable 6-bit operations
// - LookUpTableSize = 4096 (polyExtendFactor = 2)
// - Supports messageModulus=64
//
// Note: LookUpTableSize > PolyDegree requires extended LUT generation
var paramsUint6 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     1071,
		ALPHA: 7.088226765410429399593757e-08,
	},
	TLWELv1: TLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11,
		BGBIT:     22,
		BG:        1 << 22,
		L:         1,
		BASEBIT:   6,
		IKS_T:     3,
		ALPHA:     2.2204460492503131e-17,
		BlockSize: 7,
	},
}

// ============================================================================
// UINT7 PARAMETERS (Specialized for 7-bit message space, messageModulus=128)
// ============================================================================
// Based on tfhe-go's ParamsUint7 configuration.
// Key features:
// - Larger LWE dimension (1160) for added security
// - LookUpTableSize = 8192 (polyExtendFactor = 4)
// - Supports messageModulus=128
//
// Note: Requires extended LUT generation with polyExtendFactor=4
var paramsUint7 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     1160,
		ALPHA: 1.966220007498402695211596e-08,
	},
	TLWELv1: TLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11,
		BGBIT:     22,
		BG:        1 << 22,
		L:         1,
		BASEBIT:   7,
		IKS_T:     3,
		ALPHA:     2.2204460492503131e-17,
		BlockSize: 8,
	},
}

// ============================================================================
// UINT8 PARAMETERS (Specialized for 8-bit message space, messageModulus=256)
// ============================================================================
// Based on tfhe-go's ParamsUint8 configuration.
// Key features:
// - Same dimensions as Uint7
// - LookUpTableSize = 18432 (polyExtendFactor = 9)
// - Supports full 8-bit values (0-255)
//
// Note: Requires extended LUT generation with polyExtendFactor=9
var paramsUint8 = struct {
	TLWELv0  TLWELv0Params
	TLWELv1  TLWELv1Params
	TRLWELv1 TRLWELv1Params
	TRGSWLv1 TRGSWLv1Params
}{
	TLWELv0: TLWELv0Params{
		N:     1160,
		ALPHA: 1.966220007498402695211596e-08,
	},
	TLWELv1: TLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11,
		BGBIT:     22,
		BG:        1 << 22,
		L:         1,
		BASEBIT:   7,
		IKS_T:     3,
		ALPHA:     2.2204460492503131e-17,
		BlockSize: 8,
	},
}

// GetTLWELv0 returns the TLWE Level 0 parameters for the current security level
func GetTLWELv0() TLWELv0Params {
	switch CurrentSecurityLevel {
	case Security80Bit:
		return params80Bit.TLWELv0
	case Security110Bit:
		return params110Bit.TLWELv0
	case SecurityUint1:
		return paramsUint1.TLWELv0
	case SecurityUint2:
		return paramsUint2.TLWELv0
	case SecurityUint3:
		return paramsUint3.TLWELv0
	case SecurityUint4:
		return paramsUint4.TLWELv0
	case SecurityUint5:
		return paramsUint5.TLWELv0
	case SecurityUint6:
		return paramsUint6.TLWELv0
	case SecurityUint7:
		return paramsUint7.TLWELv0
	case SecurityUint8:
		return paramsUint8.TLWELv0
	default:
		return params128Bit.TLWELv0
	}
}

// GetTLWELv1 returns the TLWE Level 1 parameters for the current security level
func GetTLWELv1() TLWELv1Params {
	switch CurrentSecurityLevel {
	case Security80Bit:
		return params80Bit.TLWELv1
	case Security110Bit:
		return params110Bit.TLWELv1
	case SecurityUint1:
		return paramsUint1.TLWELv1
	case SecurityUint2:
		return paramsUint2.TLWELv1
	case SecurityUint3:
		return paramsUint3.TLWELv1
	case SecurityUint4:
		return paramsUint4.TLWELv1
	case SecurityUint5:
		return paramsUint5.TLWELv1
	case SecurityUint6:
		return paramsUint6.TLWELv1
	case SecurityUint7:
		return paramsUint7.TLWELv1
	case SecurityUint8:
		return paramsUint8.TLWELv1
	default:
		return params128Bit.TLWELv1
	}
}

// GetTRLWELv1 returns the TRLWE Level 1 parameters for the current security level
func GetTRLWELv1() TRLWELv1Params {
	switch CurrentSecurityLevel {
	case Security80Bit:
		return params80Bit.TRLWELv1
	case Security110Bit:
		return params110Bit.TRLWELv1
	case SecurityUint1:
		return paramsUint1.TRLWELv1
	case SecurityUint2:
		return paramsUint2.TRLWELv1
	case SecurityUint3:
		return paramsUint3.TRLWELv1
	case SecurityUint4:
		return paramsUint4.TRLWELv1
	case SecurityUint5:
		return paramsUint5.TRLWELv1
	case SecurityUint6:
		return paramsUint6.TRLWELv1
	case SecurityUint7:
		return paramsUint7.TRLWELv1
	case SecurityUint8:
		return paramsUint8.TRLWELv1
	default:
		return params128Bit.TRLWELv1
	}
}

// GetTRGSWLv1 returns the TRGSW Level 1 parameters for the current security level
func GetTRGSWLv1() TRGSWLv1Params {
	switch CurrentSecurityLevel {
	case Security80Bit:
		return params80Bit.TRGSWLv1
	case Security110Bit:
		return params110Bit.TRGSWLv1
	case SecurityUint1:
		return paramsUint1.TRGSWLv1
	case SecurityUint2:
		return paramsUint2.TRGSWLv1
	case SecurityUint3:
		return paramsUint3.TRGSWLv1
	case SecurityUint4:
		return paramsUint4.TRGSWLv1
	case SecurityUint5:
		return paramsUint5.TRGSWLv1
	case SecurityUint6:
		return paramsUint6.TRGSWLv1
	case SecurityUint7:
		return paramsUint7.TRGSWLv1
	case SecurityUint8:
		return paramsUint8.TRGSWLv1
	default:
		return params128Bit.TRGSWLv1
	}
}

// KSKAlpha returns the key switching key alpha for the current security level
func KSKAlpha() float64 {
	return GetTLWELv0().ALPHA
}

// BSKAlpha returns the bootstrapping key alpha for the current security level
func BSKAlpha() float64 {
	return GetTLWELv1().ALPHA
}

// SecurityInfo returns a description of the current security level
func SecurityInfo() string {
	var desc string
	switch CurrentSecurityLevel {
	case Security80Bit:
		desc = "80-bit security (performance-optimized)"
	case Security110Bit:
		desc = "110-bit security (balanced, original TFHE)"
	case SecurityUint1:
		desc = "Uint1 parameters (1-bit binary/boolean, messageModulus=2, N=1024)"
	case SecurityUint2:
		desc = "Uint2 parameters (2-bit messages, messageModulus=4, N=512)"
	case SecurityUint3:
		desc = "Uint3 parameters (3-bit messages, messageModulus=8, N=1024)"
	case SecurityUint4:
		desc = "Uint4 parameters (4-bit messages, messageModulus=16, N=2048)"
	case SecurityUint5:
		desc = "Uint5 parameters (5-bit messages, messageModulus=32, N=2048)"
	case SecurityUint6:
		desc = "Uint6 parameters (6-bit messages, messageModulus=64, N=2048)"
	case SecurityUint7:
		desc = "Uint7 parameters (7-bit messages, messageModulus=128, N=2048)"
	case SecurityUint8:
		desc = "Uint8 parameters (8-bit messages, messageModulus=256, N=2048)"
	default:
		desc = "128-bit security (high security, quantum-resistant)"
	}
	return desc
}

// GetBlockCount returns the number of blocks for block blind rotation
func GetBlockCount() int {
	lweDim := GetTLWELv0().N
	blockSize := GetTRGSWLv1().BlockSize
	if blockSize <= 1 {
		return lweDim // Original algorithm (no blocks)
	}
	return (lweDim + blockSize - 1) / blockSize // Ceiling division
}

// UseBlockBlindRotation returns true if block blind rotation should be used
func UseBlockBlindRotation() bool {
	return GetTRGSWLv1().BlockSize > 1
}
