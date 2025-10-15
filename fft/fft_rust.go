//go:build rust
// +build rust

package fft

// #cgo LDFLAGS: -L${SRCDIR}/../fft-bridge/target/release -ltfhe_fft_bridge
// #include <stdint.h>
// #include <stdlib.h>
//
// // Opaque FFT processor handle
// typedef void* FFTProcessorHandle;
//
// // C function declarations (match rs-tfhe signatures)
// extern FFTProcessorHandle fft_processor_new();
// extern void fft_processor_free(FFTProcessorHandle processor);
// extern void ifft_1024_negacyclic(FFTProcessorHandle processor, const uint32_t* torus_in, double* freq_out);
// extern void fft_1024_negacyclic(FFTProcessorHandle processor, const double* freq_in, uint32_t* torus_out);
// extern void batch_ifft_1024_negacyclic(FFTProcessorHandle processor, const uint32_t* torus_in, double* freq_out, size_t count);
import "C"
import (
	"unsafe"

	"github.com/thedonutfactory/go-tfhe/params"
)

// FFTProcessor uses Rust's realfft/rustfft via CGO for maximum performance
type FFTProcessor struct {
	handle C.FFTProcessorHandle
}

// NewFFTProcessor creates a new Rust-backed FFT processor
func NewFFTProcessor(n int) *FFTProcessor {
	if n != 1024 {
		panic("Only N=1024 supported for now")
	}
	handle := C.fft_processor_new()
	if handle == nil {
		panic("failed to create Rust FFT processor")
	}
	return &FFTProcessor{handle: handle}
}

// Free releases the Rust FFT processor resources
func (p *FFTProcessor) Free() {
	if p.handle != nil {
		C.fft_processor_free(p.handle)
		p.handle = nil
	}
}

// IFFT1024 transforms time domain → frequency domain (Torus → float64)
// Matches pure Go signature: IFFT1024(input *[1024]params.Torus) [1024]float64
// Calls Rust: ifft_1024_negacyclic(processor, torus_in: *const u32, freq_out: *mut f64)
func (p *FFTProcessor) IFFT1024(input *[1024]params.Torus) [1024]float64 {
	var result [1024]float64

	// Call Rust IFFT: torus→freq (matches rs-tfhe signature)
	C.ifft_1024_negacyclic(
		p.handle,
		(*C.uint32_t)(unsafe.Pointer(&input[0])),
		(*C.double)(unsafe.Pointer(&result[0])),
	)

	return result
}

// FFT1024 transforms frequency domain → time domain (float64 → Torus)
// Matches pure Go signature: FFT1024(input *[1024]float64) [1024]params.Torus
// Calls Rust: fft_1024_negacyclic(processor, freq_in: *const f64, torus_out: *mut u32)
func (p *FFTProcessor) FFT1024(input *[1024]float64) [1024]params.Torus {
	var result [1024]params.Torus

	// Call Rust FFT: freq→torus (matches rs-tfhe signature)
	C.fft_1024_negacyclic(
		p.handle,
		(*C.double)(unsafe.Pointer(&input[0])),
		(*C.uint32_t)(unsafe.Pointer(&result[0])),
	)

	return result
}

// PolyMul1024 performs negacyclic polynomial multiplication using Rust FFT
func (p *FFTProcessor) PolyMul1024(a, b *[1024]params.Torus) [1024]params.Torus {
	// Forward FFT: torus→freq
	aFFT := p.IFFT1024(a)
	bFFT := p.IFFT1024(b)

	// Complex multiplication with 0.5 scaling
	var resultFFT [1024]float64
	halfN := 512
	for i := 0; i < halfN; i++ {
		aRe := aFFT[i]
		aIm := aFFT[i+halfN]
		bRe := bFFT[i]
		bIm := bFFT[i+halfN]

		// Complex multiply: (a_re + i*a_im) * (b_re + i*b_im)
		// Result scaled by 0.5 for negacyclic convolution
		resultFFT[i] = (aRe*bRe - aIm*bIm) * 0.5
		resultFFT[i+halfN] = (aRe*bIm + aIm*bRe) * 0.5
	}

	// Inverse FFT: freq→torus
	return p.FFT1024(&resultFFT)
}

// IFFT transforms time domain (N values) → frequency domain (N values)
// Convenience wrapper for slices
func (p *FFTProcessor) IFFT(input []params.Torus) []float64 {
	var arr [1024]params.Torus
	copy(arr[:], input)
	result := p.IFFT1024(&arr)
	return result[:]
}

// FFT transforms frequency domain (N values) → time domain (N values)
// Convenience wrapper for slices
func (p *FFTProcessor) FFT(input []float64) []params.Torus {
	var arr [1024]float64
	copy(arr[:], input)
	result := p.FFT1024(&arr)
	return result[:]
}

// PolyMul performs negacyclic polynomial multiplication for variable-length vectors
func (p *FFTProcessor) PolyMul(a, b []params.Torus) []params.Torus {
	if len(a) == 1024 && len(b) == 1024 {
		var aArr [1024]params.Torus
		var bArr [1024]params.Torus
		copy(aArr[:], a)
		copy(bArr[:], b)
		result := p.PolyMul1024(&aArr, &bArr)
		return result[:]
	}
	panic("PolyMul only supports 1024-element inputs")
}

// BatchIFFT1024 transforms multiple polynomials at once (Torus → float64)
func (p *FFTProcessor) BatchIFFT1024(inputs [][1024]params.Torus) [][1024]float64 {
	results := make([][1024]float64, len(inputs))
	for i := range inputs {
		results[i] = p.IFFT1024(&inputs[i])
	}
	return results
}

// BatchFFT1024 transforms multiple frequency-domain representations at once (float64 → Torus)
func (p *FFTProcessor) BatchFFT1024(inputs [][1024]float64) [][1024]params.Torus {
	results := make([][1024]params.Torus, len(inputs))
	for i := range inputs {
		results[i] = p.FFT1024(&inputs[i])
	}
	return results
}

// FFTPlan provides FFT planning and execution with Rust backend
type FFTPlan struct {
	Processor *FFTProcessor
	N         int
}

// NewFFTPlan creates a new FFT plan for the given polynomial size
func NewFFTPlan(n int) *FFTPlan {
	return &FFTPlan{
		Processor: NewFFTProcessor(n),
		N:         n,
	}
}
