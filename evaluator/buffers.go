package evaluator

import (
	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/thedonutfactory/go-tfhe/poly"
	"github.com/thedonutfactory/go-tfhe/tlwe"
	"github.com/thedonutfactory/go-tfhe/trlwe"
)

// BufferPool is a centralized buffer management system for zero-allocation TFHE operations.
// All buffers are pre-allocated once during initialization and reused throughout computation.
//
// Memory Layout Overview:
//   - Polynomial buffers: ~200 KB (FFT, decomposition, rotation)
//   - Ciphertext buffers: ~50 KB (TRLWE, LWE intermediates)
//   - Total: ~250 KB per evaluator instance
//
// Thread Safety:
//   - Each evaluator has its own buffer pool
//   - Use sync.Pool or create separate instances for concurrent operations
type BufferPool struct {
	// === Polynomial Operation Buffers ===
	// Used for FFT, multiplication, and decomposition

	// PolyBuffers manages polynomial and FFT operations
	PolyBuffers *poly.BufferManager

	// === Bootstrap Operation Buffers ===

	// External Product buffers (TRGSW ⊗ TRLWE)
	ExternalProduct struct {
		// Fourier domain accumulators
		FourierA poly.FourierPoly // ~8 KB
		FourierB poly.FourierPoly // ~8 KB
		// Time domain result
		Result *trlwe.TRLWELv1 // ~8 KB
	}

	// CMUX buffers (conditional multiplexer)
	CMUX struct {
		Temp *trlwe.TRLWELv1 // Difference buffer (ct1 - ct0)
	}

	// Blind Rotation buffers
	BlindRotation struct {
		Accumulator1 *trlwe.TRLWELv1 // Primary accumulator
		Accumulator2 *trlwe.TRLWELv1 // Secondary accumulator
		Rotated      *trlwe.TRLWELv1 // Rotation result
	}

	// Bootstrap buffers (full bootstrap = blind rotate + key switch)
	Bootstrap struct {
		ExtractedLWE *tlwe.TLWELv1 // After sample extraction
		KeySwitched  *tlwe.TLWELv0 // After key switching
	}

	// === Gate Operation Buffers ===

	// Gate preparation buffer (for AND, OR, XOR, etc.)
	GatePrep *tlwe.TLWELv0

	// Result pool for returning values without allocation
	// Round-robin buffer to handle compound operations (e.g., MUX)
	ResultPool struct {
		Buffers [4]*tlwe.TLWELv0 // 4 slots for compound operations
		Index   int              // Current index (0-3)
	}

	// === Block Blind Rotation Buffers (for 3-4x speedup) ===
	// Only allocated if params.UseBlockBlindRotation() == true

	BlockRotation *BlockRotationBuffers
}

// BlockRotationBuffers contains buffers for block-based blind rotation algorithm
// This provides 3-4x speedup by processing multiple LWE coefficients together
type BlockRotationBuffers struct {
	// Decomposed accumulator in Fourier domain
	// [blockSize][glweRank+1][level]
	AccFourierDecomposed [][][]poly.FourierPoly

	// Block accumulator in Fourier domain [blockSize]
	BlockFourierAcc []struct {
		A poly.FourierPoly
		B poly.FourierPoly
	}

	// Intermediate Fourier accumulator [blockSize]
	FourierAcc []struct {
		A poly.FourierPoly
		B poly.FourierPoly
	}

	// Fourier monomial for multiplication
	FourierMono poly.FourierPoly
}

// NewBufferPool creates a new centralized buffer pool for the given polynomial size.
// This allocates all buffers once during initialization (~250 KB total).
//
// Parameters:
//
//	n: Polynomial degree (typically 1024 for standard TFHE parameters)
//
// Memory allocation:
//   - Polynomial buffers: ~200 KB (managed by poly.BufferManager)
//   - Ciphertext buffers: ~50 KB (TRLWE, LWE structures)
//   - Block rotation: ~30 KB (if enabled)
func NewBufferPool(n int) *BufferPool {
	bp := &BufferPool{
		PolyBuffers: poly.NewBufferManager(n),
	}

	// Initialize external product buffers
	bp.ExternalProduct.FourierA = poly.NewFourierPoly(n)
	bp.ExternalProduct.FourierB = poly.NewFourierPoly(n)
	bp.ExternalProduct.Result = trlwe.NewTRLWELv1()

	// Initialize CMUX buffers
	bp.CMUX.Temp = trlwe.NewTRLWELv1()

	// Initialize blind rotation buffers
	bp.BlindRotation.Accumulator1 = trlwe.NewTRLWELv1()
	bp.BlindRotation.Accumulator2 = trlwe.NewTRLWELv1()
	bp.BlindRotation.Rotated = trlwe.NewTRLWELv1()

	// Initialize bootstrap buffers
	bp.Bootstrap.ExtractedLWE = tlwe.NewTLWELv1()
	bp.Bootstrap.KeySwitched = tlwe.NewTLWELv0()

	// Initialize gate preparation buffer
	bp.GatePrep = tlwe.NewTLWELv0()

	// Initialize result pool
	for i := 0; i < 4; i++ {
		bp.ResultPool.Buffers[i] = tlwe.NewTLWELv0()
	}
	bp.ResultPool.Index = 0

	// Initialize block rotation buffers if enabled
	if params.UseBlockBlindRotation() {
		bp.BlockRotation = newBlockRotationBuffers(n)
	}

	return bp
}

// newBlockRotationBuffers creates buffers for block-based blind rotation
func newBlockRotationBuffers(n int) *BlockRotationBuffers {
	blockSize := params.GetTRGSWLv1().BlockSize
	if blockSize < 1 {
		blockSize = 1
	}
	glweRank := 1 // Fixed for our parameters
	level := params.GetTRGSWLv1().L

	brb := &BlockRotationBuffers{}

	// Initialize AccFourierDecomposed[blockSize][glweRank+1][level]
	brb.AccFourierDecomposed = make([][][]poly.FourierPoly, blockSize)
	for i := 0; i < blockSize; i++ {
		brb.AccFourierDecomposed[i] = make([][]poly.FourierPoly, glweRank+1)
		for j := 0; j < glweRank+1; j++ {
			brb.AccFourierDecomposed[i][j] = make([]poly.FourierPoly, level)
			for k := 0; k < level; k++ {
				brb.AccFourierDecomposed[i][j][k] = poly.NewFourierPoly(n)
			}
		}
	}

	// Initialize BlockFourierAcc[blockSize]
	brb.BlockFourierAcc = make([]struct {
		A poly.FourierPoly
		B poly.FourierPoly
	}, blockSize)
	for i := 0; i < blockSize; i++ {
		brb.BlockFourierAcc[i].A = poly.NewFourierPoly(n)
		brb.BlockFourierAcc[i].B = poly.NewFourierPoly(n)
	}

	// Initialize FourierAcc[blockSize]
	brb.FourierAcc = make([]struct {
		A poly.FourierPoly
		B poly.FourierPoly
	}, blockSize)
	for i := 0; i < blockSize; i++ {
		brb.FourierAcc[i].A = poly.NewFourierPoly(n)
		brb.FourierAcc[i].B = poly.NewFourierPoly(n)
	}

	// Initialize Fourier monomial
	brb.FourierMono = poly.NewFourierPoly(n)

	return brb
}

// GetNextResult returns the next available result buffer from the round-robin pool.
// This allows operations to return results without allocation.
// The buffer is valid until 4 more operations are performed.
func (bp *BufferPool) GetNextResult() *tlwe.TLWELv0 {
	result := bp.ResultPool.Buffers[bp.ResultPool.Index]
	bp.ResultPool.Index = (bp.ResultPool.Index + 1) % 4
	return result
}

// Reset resets all buffer pool indices to their initial state.
// Call this when reusing an evaluator for a new computation.
func (bp *BufferPool) Reset() {
	bp.ResultPool.Index = 0
	bp.PolyBuffers.Reset()
}

// MemoryUsage returns the approximate memory usage in bytes
func (bp *BufferPool) MemoryUsage() int {
	n := params.GetTRGSWLv1().N

	// Polynomial buffers (managed by poly.BufferManager)
	polyMem := bp.PolyBuffers.MemoryUsage()

	// Ciphertext buffers
	trlweSize := 2 * n * 4  // 2 polynomials * N elements * 4 bytes
	tlweSize := (n + 1) * 4 // (N+1) elements * 4 bytes

	ciphertextMem := trlweSize*5 + // 5 TRLWE buffers
		tlweSize*5 + // 5 LWE buffers
		2*n*8 // 2 FourierPoly in ExternalProduct

	// Block rotation buffers (if enabled)
	blockMem := 0
	if bp.BlockRotation != nil {
		blockSize := params.GetTRGSWLv1().BlockSize
		level := params.GetTRGSWLv1().L
		glweRank := 1
		blockMem = blockSize * (glweRank + 1) * level * n * 8 * 2 // AccFourierDecomposed
		blockMem += blockSize * 2 * n * 8 * 2                     // BlockFourierAcc + FourierAcc
		blockMem += n * 8 * 2                                     // FourierMono
	}

	return polyMem + ciphertextMem + blockMem
}
