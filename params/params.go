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
	SecurityUint5  SecurityLevel = 5 // Specialized for 5-bit message space (messageModulus=32)
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
// UINT5 PARAMETERS (Specialized for 5-bit message space, messageModulus=32)
// ============================================================================
// Based on tfhe-go's ParamsUint5 configuration. These parameters are
// specifically designed for multi-bit arithmetic with large message spaces.
// Key features:
// - ~700x lower noise than standard 80-bit security
// - Larger polynomial degree (2048 vs 1024)
// - Supports messageModulus up to 32 reliably
// - Enables 4-bootstrap nibble addition (vs 8 with standard params)
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
		N:     1071,                           // Larger LWE dimension for security
		ALPHA: 7.088226765410429399593757e-08, // ~700x lower noise for precision
	},
	TLWELv1: TLWELv1Params{
		N:     2048,                   // Double polynomial degree
		ALPHA: 2.2204460492503131e-17, // Extremely low noise
	},
	TRLWELv1: TRLWELv1Params{
		N:     2048,
		ALPHA: 2.2204460492503131e-17,
	},
	TRGSWLv1: TRGSWLv1Params{
		N:         2048,
		NBIT:      11, // 2048 = 2^11
		BGBIT:     22, // Base = 1 << 22
		BG:        1 << 22,
		L:         1, // Single decomposition level
		BASEBIT:   6, // KeySwitch base bits
		IKS_T:     3, // KeySwitch level
		ALPHA:     2.2204460492503131e-17,
		BlockSize: 7, // Optimized block size for N=1071
	},
}

// GetTLWELv0 returns the TLWE Level 0 parameters for the current security level
func GetTLWELv0() TLWELv0Params {
	switch CurrentSecurityLevel {
	case Security80Bit:
		return params80Bit.TLWELv0
	case Security110Bit:
		return params110Bit.TLWELv0
	case SecurityUint5:
		return paramsUint5.TLWELv0
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
	case SecurityUint5:
		return paramsUint5.TLWELv1
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
	case SecurityUint5:
		return paramsUint5.TRLWELv1
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
	case SecurityUint5:
		return paramsUint5.TRGSWLv1
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
	case SecurityUint5:
		desc = "Uint5 parameters (optimized for 5-bit messages, messageModulus=32)"
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
