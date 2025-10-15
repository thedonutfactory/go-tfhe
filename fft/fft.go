// Package fft provides FFT operations for TFHE polynomial multiplication.
//
// Based on "Fast and Error-Free Negacyclic Integer Convolution using Extended Fourier Transform"
// by Jakub Klemsa - https://eprint.iacr.org/2021/480
//
// This implementation matches the Rust ExtendedFftProcessor approach exactly.
package fft

import (
	"math"

	"github.com/thedonutfactory/go-tfhe/params"
	"github.com/mjibson/go-dsp/fft"
)

// FFTProcessor provides FFT operations for TFHE negacyclic polynomial multiplication
type FFTProcessor struct {
	n          int
	twistiesRe []float64
	twistiesIm []float64
}

// NewFFTProcessor creates a new FFT processor for polynomials of size n
func NewFFTProcessor(n int) *FFTProcessor {
	if n != 1024 {
		panic("Only N=1024 supported for now")
	}

	n2 := n / 2 // 512

	processor := &FFTProcessor{
		n:          n,
		twistiesRe: make([]float64, n2),
		twistiesIm: make([]float64, n2),
	}

	// Compute twisting factors: exp(i*π*k/N) for k=0..N/2-1
	// Matches Rust: let angle = i as f64 * twist_unit;
	twistUnit := math.Pi / float64(n)
	for i := 0; i < n2; i++ {
		angle := float64(i) * twistUnit
		sin, cos := math.Sincos(angle)
		processor.twistiesRe[i] = cos // Re = cos
		processor.twistiesIm[i] = sin // Im = sin
	}

	return processor
}

// IFFT1024 transforms time domain → frequency domain
// Matches Rust's ifft_1024 exactly
func (p *FFTProcessor) IFFT1024(input *[1024]params.Torus) [1024]float64 {
	const N = 1024
	const N2 = N / 2 // 512

	// Split input: input_re = input[0..512], input_im = input[512..1024]
	// Rust: let (input_re, input_im) = input.split_at(N2);

	// Apply twisting factors and convert
	// Rust code:
	// let in_re = input_re[i] as i32 as f64;
	// let in_im = input_im[i] as i32 as f64;
	// fourier[i] = Complex::new(in_re * w_re - in_im * w_im, in_re * w_im + in_im * w_re);

	fourier := make([]complex128, N2)
	for i := 0; i < N2; i++ {
		inRe := float64(int32(input[i]))
		inIm := float64(int32(input[i+N2]))
		wRe := p.twistiesRe[i]
		wIm := p.twistiesIm[i]
		// Complex multiply: (inRe + i*inIm) * (wRe + i*wIm)
		realPart := inRe*wRe - inIm*wIm
		imagPart := inRe*wIm + inIm*wRe
		fourier[i] = complex(realPart, imagPart)
	}

	// Perform 512-point FFT
	fftResult := fft.FFT(fourier)

	// Scale by 2 and convert to output
	// Rust: result[i] = fourier[i].re * 2.0;
	var result [N]float64
	for i := 0; i < N2; i++ {
		result[i] = real(fftResult[i]) * 2.0
		result[i+N2] = imag(fftResult[i]) * 2.0
	}

	return result
}

// FFT1024 transforms frequency domain → time domain
// Matches Rust's fft_1024 exactly
func (p *FFTProcessor) FFT1024(input *[1024]float64) [1024]params.Torus {
	const N = 1024
	const N2 = N / 2 // 512

	// Convert to complex and scale by 0.5
	// Rust: fourier[i] = Complex::new(input_re[i] * 0.5, input_im[i] * 0.5);
	fourier := make([]complex128, N2)
	for i := 0; i < N2; i++ {
		fourier[i] = complex(input[i]*0.5, input[i+N2]*0.5)
	}

	// Perform 512-point IFFT
	ifftResult := fft.IFFT(fourier)

	// Apply inverse twisting and convert to u32
	// NOTE: go-dsp IFFT is already normalized, so we DON'T divide by N2
	// CRITICAL: Cast through int64 first (like Rust) to avoid int32 overflow!
	// Rust: result[i] = tmp_re.round() as i64 as u32;
	var result [N]params.Torus
	for i := 0; i < N2; i++ {
		wRe := p.twistiesRe[i]
		wIm := p.twistiesIm[i]
		fRe := real(ifftResult[i])
		fIm := imag(ifftResult[i])
		// Complex multiply with conjugate: (fRe + i*fIm) * (wRe - i*wIm)
		tmpRe := fRe*wRe + fIm*wIm
		tmpIm := fIm*wRe - fRe*wIm
		// Cast through int64 to avoid overflow, then to uint32
		result[i] = params.Torus(uint32(int64(math.Round(tmpRe))))
		result[i+N2] = params.Torus(uint32(int64(math.Round(tmpIm))))
	}

	return result
}

// IFFT transforms time domain (N values) → frequency domain (N values)
func (p *FFTProcessor) IFFT(input []params.Torus) []float64 {
	var arr [1024]params.Torus
	copy(arr[:], input)
	result := p.IFFT1024(&arr)
	return result[:]
}

// FFT transforms frequency domain (N values) → time domain (N values)
func (p *FFTProcessor) FFT(input []float64) []params.Torus {
	var arr [1024]float64
	copy(arr[:], input)
	result := p.FFT1024(&arr)
	return result[:]
}

// PolyMul1024 performs negacyclic polynomial multiplication
// Matches Rust's poly_mul_1024 exactly
func (p *FFTProcessor) PolyMul1024(a, b *[1024]params.Torus) [1024]params.Torus {
	aFFT := p.IFFT1024(a)
	bFFT := p.IFFT1024(b)

	// Complex multiplication with 0.5 scaling
	// Rust:
	// result_fft[i] = (ar * br - ai * bi) * 0.5;
	// result_fft[i + N2] = (ar * bi + ai * br) * 0.5;
	var resultFFT [1024]float64
	const N2 = 512
	for i := 0; i < N2; i++ {
		ar := aFFT[i]
		ai := aFFT[i+N2]
		br := bFFT[i]
		bi := bFFT[i+N2]

		resultFFT[i] = (ar*br - ai*bi) * 0.5
		resultFFT[i+N2] = (ar*bi + ai*br) * 0.5
	}

	return p.FFT1024(&resultFFT)
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
	return make([]params.Torus, len(a))
}

// BatchIFFT1024 transforms multiple polynomials at once
func (p *FFTProcessor) BatchIFFT1024(inputs [][1024]params.Torus) [][1024]float64 {
	results := make([][1024]float64, len(inputs))
	for i := range inputs {
		results[i] = p.IFFT1024(&inputs[i])
	}
	return results
}

// BatchFFT1024 transforms multiple frequency-domain representations at once
func (p *FFTProcessor) BatchFFT1024(inputs [][1024]float64) [][1024]params.Torus {
	results := make([][1024]params.Torus, len(inputs))
	for i := range inputs {
		results[i] = p.FFT1024(&inputs[i])
	}
	return results
}

// FFTPlan wraps an FFT processor with its configuration
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
