package tfhe

import (
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

var fp1024Nayuki *NayukiFFTProcessor = NewNayukiFFTProcessor(1024)

//var fp1024Nayuki *NayukiFFTProcessor = NewNayukiFFTProcessor(4)

type FFTProcessor interface {
	executeReverseTorus(a []Torus) (res []complex128)
	executeReverseInt(a []int) (res []complex128)
	executeDirectTorus(a []complex128) (res []Torus)
}

type NayukiFFTProcessor struct {
	_2N int
	N   int
	Ns2 int

	realInout     []double
	imagInout     []double
	tablesDirect  *FftTables
	tablesReverse *FftTablesUint

	omegaxminus1 []complex128
}

// Private data structure
type FftTables struct {
	n           int
	bitReversed []int
	cosTable    []double
	sinTable    []double
}

func NewNayukiFFTProcessor(N int) *NayukiFFTProcessor {
	p := &NayukiFFTProcessor{
		_2N:           2 * N,
		N:             N,
		Ns2:           N / 2,
		realInout:     make([]double, 2*N),
		imagInout:     make([]double, 2*N),
		tablesDirect:  fftInit(2 * N),
		tablesReverse: fftInitReverse(2 * N),
		omegaxminus1:  make([]complex128, 2*N),
	}

	for x := int(0); x < 2*N; x++ {
		p.omegaxminus1[x] = complex(math.Cos(double(x)*math.Pi/double(N))-1., math.Sin(double(x)*math.Pi/double(N)))
		// instead of cos(x*M_PI/N)-1. + sin(x*M_PI/N) * 1i
	}
	return p
}

func (p *NayukiFFTProcessor) checkAlternateReal() {
	if debug {
		for i := int(0); i < p._2N; i++ {
			Assert(math.Abs(p.imagInout[i]) < 1e-8)
		}
		for i := int(0); i < p.N; i++ {
			Assert(math.Abs(p.realInout[i]+p.realInout[p.N+i]) < 1e-9)
		}
	}
}

func (p *NayukiFFTProcessor) checkConjugateCplx() {
	if debug {
		for i := int(0); i < p.N; i++ {
			Assert(math.Abs(p.realInout[2*i])+math.Abs(p.imagInout[2*i]) < 1e-20)
		}
		for i := int(0); i < p.Ns2; i++ {
			a := p.imagInout[2*i+1]
			b := p.imagInout[p._2N-1-2*i]
			toler := 1e-20
			if math.Abs(a-b) >= toler {
				fmt.Printf("Error %f (%f + %f) >= %.20f", a, b, math.Abs(a-b), toler)
			}
			Assert(math.Abs(a-b) < toler)
		}
	}
}

func (p *NayukiFFTProcessor) executeReverseTorus(a []Torus) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func (p *NayukiFFTProcessor) executeReverseInt(a []int) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func (p *NayukiFFTProcessor) executeDirectTorus(a []complex128) (res []Torus) {
	res = castTorus(fft.FFT(a))
	for i := 0; i < int(p.Ns2); i++ {
		res = append(res, 0)
	}
	return

}

/**
 * FFT functions
 */

func intPolynomialIfft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	result.coefsC = fp1024Nayuki.executeReverseInt(p.Coefs)
}

func torusPolynomialIfft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	result.coefsC = fp1024Nayuki.executeReverseTorus(p.CoefsT)
}

func torusPolynomialFft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	result.CoefsT = fp1024Nayuki.executeDirectTorus(p.coefsC)
}

func fftInit(n int) *FftTables {
	// Check size argument
	if n <= 0 || (n&(n-1)) != 0 {
		return nil // Error: Size is not a power of 2
	}

	tables := &FftTables{
		n:           n,
		bitReversed: make([]int, n),
		cosTable:    make([]double, n/2),
		sinTable:    make([]double, n/2),
	}

	// Precompute values and store to tables
	levels := floorLog2(int(n))
	for i := 0; i < n; i++ {
		tables.bitReversed[i] = reverseBits(int(i), uint(levels))
	}
	for i := 0; i < n/2; i++ {
		var angle double = 2. * math.Pi * double(i) / double(n)
		tables.cosTable[i] = math.Cos(angle)
		tables.sinTable[i] = math.Sin(angle)
	}
	return tables
}

// Performs a forward FFT in place on the given arrays. The length is given by the tables struct.
func (p *NayukiFFTProcessor) fftTransform(tbl *FftTables, real []double, imag []double) {
	n := tbl.n

	// Bit-reversed addressing permutation
	bitreversed := tbl.bitReversed
	for i := 0; i < n; i++ {
		j := bitreversed[i]
		if i < int(j) {
			tp0re := real[i]
			tp0im := imag[i]
			tp1re := real[j]
			tp1im := imag[j]
			real[i] = tp1re
			imag[i] = tp1im
			real[j] = tp0re
			imag[j] = tp0im
		}
	}

	// Cooley-Tukey decimation-in-time radix-2 FFT
	costable := tbl.cosTable
	sintable := tbl.sinTable
	for size := int(2); size <= n; size *= 2 {
		halfsize := size / 2
		tablestep := n / size
		for i := 0; i < n; i += size {
			j := i
			for k := 0; j < i+halfsize; k += tablestep {
				tpre := real[j+halfsize]*costable[k] + imag[j+halfsize]*sintable[k]
				tpim := -real[j+halfsize]*sintable[k] + imag[j+halfsize]*costable[k]
				real[j+halfsize] = real[j] - tpre
				imag[j+halfsize] = imag[j] - tpim
				real[j] += tpre
				imag[j] += tpim
				j++
			}
		}
		if size == n { // Prevent overflow in 'size *= 2'
			break
		}
	}
}

type FftTablesUint struct {
	n           int
	bitReversed []int
	trigTables  []double
}

// Returns sin(2 * pi * i / n), for n that is a multiple of 4.
func accurateSine(i int, n int) double {
	if n%4 != 0 {
		return 0.
	} else {
		var neg int = 0 // Boolean
		// Reduce to full cycle
		i %= n
		// Reduce to half cycle
		if i >= n/2 {
			neg = 1
			i -= n / 2
		}
		// Reduce to quarter cycle
		if i >= n/4 {
			i = n/2 - i
		}
		// Reduce to eighth cycle
		var val double
		if i*8 < n {
			val = math.Sin(float64(2) * math.Pi * float64(i) / float64(n))
		} else {
			val = math.Cos(float64(2) * math.Pi * float64(n/4-i) / float64(n))
		}
		// Apply sign
		if neg == 0 {
			return -val
		} else {
			return val
		}
		//return neg ? -val : val;
	}
}

// Returns the largest i such that 2^i <= n.
func floorLog2(n int) int {
	var result int = 0
	for ; n > 1; n /= 2 {
		result++
	}
	return result
}

// Returns the bit reversal of the n-bit unsigned integer x.
func reverseBits(x int, n uint) int {
	var result int = 0
	for i := uint(0); i < n; i++ {
		result = (result << 1) | (x & 1)
		x >>= 1
	}
	return result
}

// Returns a pointer to an opaque structure of FFT tables. n must be a power of 2 and n >= 4.
func fftInitReverse(n int) *FftTablesUint {
	// Check size argument
	if n < 4 || n > math.MaxInt || (n&(n-1)) != 0 {
		return nil // Error: Size is too small or is not a power of 2
	}

	tables := &FftTablesUint{
		n:           n,
		bitReversed: make([]int, n),
		//trigTables:  make([]double, n-4),
		trigTables: make([]double, n-4),
		//trigTables: make([]double, n*2),
	}

	// Precompute bit reversal table
	levels := floorLog2(n)
	for i := 0; i < n; i++ {
		tables.bitReversed[i] = reverseBits(i, uint(levels))
	}

	// Precompute the packed trigonometric table for each FFT internal level
	var k int = 0
	for size := int(8); size <= n; size *= 2 {
		for i := 0; i < size/2; i += 4 {
			for j := 0; j < 4; j++ {
				tables.trigTables[k] = accurateSine(i+j+size/4, size) // Cosine
				k++
			}
			k = 0
			for j := 0; j < 4; j++ {
				tables.trigTables[k] = -accurateSine(i+j, size) // Sine
				k++
			}
		}
		if size == n {
			break
		}
	}
	return tables
}

// This is a C implementation that models the x86-64 AVX implementation.
func (p *NayukiFFTProcessor) fftTransformReverse(tbl *FftTablesUint, real []double, imag []double) {
	//struct FftTables *tbl = (struct FftTables *)tables;
	n := tbl.n

	// Bit-reversed addressing permutation
	bitreversed := tbl.bitReversed
	for i := 0; i < n; i++ {
		j := bitreversed[i]
		if i < j {
			tp0re := real[i]
			tp0im := imag[i]
			tp1re := real[j]
			tp1im := imag[j]
			real[i] = tp1re
			imag[i] = tp1im
			real[j] = tp0re
			imag[j] = tp0im
		}
	}

	// Size 2 merge (special)
	if n >= 2 {
		for i := 0; i < n; i += 2 {
			tpre := real[i]
			tpim := imag[i]
			real[i] += real[i+1]
			imag[i] += imag[i+1]
			real[i+1] = tpre - real[i+1]
			imag[i+1] = tpim - imag[i+1]
		}
	}

	// Size 4 merge (special)
	if n >= 4 {
		for i := 0; i < n; i += 4 {
			// Even indices
			tpre := real[i]
			tpim := imag[i]
			real[i] += real[i+2]
			imag[i] += imag[i+2]
			real[i+2] = tpre - real[i+2]
			imag[i+2] = tpim - imag[i+2]
			// Odd indices
			tpre = real[i+1]
			tpim = imag[i+1]
			real[i+1] -= imag[i+3]
			imag[i+1] += real[i+3]
			tpre += imag[i+3]
			tpim -= real[i+3]
			real[i+3] = tpre
			imag[i+3] = tpim
		}
	}

	// Size 8 and larger merges (general)
	trigtables := tbl.trigTables
	// for size := int(8); size <= n; size <<= 1 {
	for size := 8; size <= len(trigtables); size <<= 1 {
		halfsize := size >> 1
		for i := 0; i < n; i += size {
			var j int = 0
			var off int = 0
			for ; j < halfsize; j += 4 {
				for k := 0; k < 4; k++ { // To simulate x86 AVX 4-vectors
					vi := i + j + k // Vector index
					ti := off + k   // Table index
					re := real[vi+halfsize]
					im := imag[vi+halfsize]
					tpre := re*trigtables[ti] + im*trigtables[ti+4]
					tpim := im*trigtables[ti] - re*trigtables[ti+4]
					real[vi+halfsize] = real[vi] - tpre
					imag[vi+halfsize] = imag[vi] - tpim
					real[vi] += tpre
					imag[vi] += tpim
				}
				off += 8
			}
		}
		if size == n {
			break
		}
		//trigtables += size;
		trigtables = trigtables[size:]
	}

}
