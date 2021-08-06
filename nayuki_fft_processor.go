package tfhe

import (
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

var fp1024_nayuki *NayukiFFTProcessor = NewNayukiFFTProcessor(1024)

//var fp1024_nayuki *NayukiFFTProcessor = NewNayukiFFTProcessor(4)

type FFTProcessor interface {
	executeReverseTorus32(a []Torus32) (res []complex128)
	executeReverseInt(a []int32) (res []complex128)
	executeDirectTorus32(a []complex128) (res []Torus32)
}

type NayukiFFTProcessor struct {
	_2N int32
	N   int32
	Ns2 int32

	real_inout     []double
	imag_inout     []double
	tables_direct  *FftTables
	tables_reverse *FftTablesUint64

	omegaxminus1 []complex128
}

// Private data structure
type FftTables struct {
	n            int64
	bit_reversed []int64
	cos_table    []double
	sin_table    []double
}

func NewNayukiFFTProcessor(N int32) *NayukiFFTProcessor {
	p := &NayukiFFTProcessor{
		_2N:            2 * N,
		N:              N,
		Ns2:            N / 2,
		real_inout:     make([]double, 2*N),
		imag_inout:     make([]double, 2*N),
		tables_direct:  fft_init(int64(2 * N)),
		tables_reverse: fft_init_reverse(int64(2 * N)),
		omegaxminus1:   make([]complex128, 2*N),
	}

	for x := int32(0); x < 2*N; x++ {
		p.omegaxminus1[x] = complex(math.Cos(double(x)*math.Pi/double(N))-1., math.Sin(double(x)*math.Pi/double(N)))
		// instead of cos(x*M_PI/N)-1. + sin(x*M_PI/N) * 1i
	}
	return p
}

func (p *NayukiFFTProcessor) check_alternate_real() {
	if debug {
		for i := int32(0); i < p._2N; i++ {
			Assert(math.Abs(p.imag_inout[i]) < 1e-8)
		}
		for i := int32(0); i < p.N; i++ {
			Assert(math.Abs(p.real_inout[i]+p.real_inout[p.N+i]) < 1e-9)
		}
	}
}

func (p *NayukiFFTProcessor) check_conjugate_cplx() {
	if debug {
		for i := int32(0); i < p.N; i++ {
			Assert(math.Abs(p.real_inout[2*i])+math.Abs(p.imag_inout[2*i]) < 1e-20)
		}
		for i := int32(0); i < p.Ns2; i++ {
			a := p.imag_inout[2*i+1]
			b := p.imag_inout[p._2N-1-2*i]
			toler := 1e-20
			if math.Abs(a-b) >= toler {
				fmt.Printf("Error %f (%f + %f) >= %.20f", a, b, math.Abs(a-b), toler)
			}
			Assert(math.Abs(a-b) < toler)
		}
	}
}

func (p *NayukiFFTProcessor) executeReverseTorus32(a []Torus32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func (p *NayukiFFTProcessor) executeReverseInt(a []int32) (res []complex128) {
	res = fft.IFFT(castComplex(a))
	return
}

func (p *NayukiFFTProcessor) executeDirectTorus32(a []complex128) (res []Torus32) {
	res = castTorus(fft.FFT(a))
	for i := 0; i < int(p.Ns2); i++ {
		res = append(res, 0)
	}
	return

}

/*
void FFT_Processor_fftw::execute_direct_Torus32(Torus32* res, const cplx* a) {
    static const double _2p32 = double(INT64_C(1)<<32);
    static const double _1sN = double(1)/double(N);
    cplx* in_cplx = (cplx*) in; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<=Ns2; i++) in_cplx[2*i]=0;
    for (int32_t i=0; i<Ns2; i++) in_cplx[2*i+1]=a[i];
    fftw_execute(p);
    for (int32_t i=0; i<N; i++) res[i]=Torus32(int64_t(out[i]*_1sN*_2p32));
    //pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
    for (int32_t i=0; i<N; i++) assert(fabs(out[N+i]+out[i])<1e-20);
}
*/

//////////
/*
func (p *NayukiFFTProcessor) executeReverseInt2(a []int32) (res []complex128) {
	//double* res_dbl=(double*) res;
	N := int(p.N)
	Ns2 := int(p.Ns2)
	_2N := int(p._2N)
	//real_inout := p.real_inout
	//imag_inout := p.imag_inout
	res_dbl := complexToFloatSlice(res)
	for i := 0; i < N; i++ {
		p.real_inout[i] = float64(a[i]) / 2.
	}
	for i := 0; i < N; i++ {
		p.real_inout[N+i] = -p.real_inout[i]
	}
	for i := 0; i < _2N; i++ {
		p.imag_inout[i] = 0
	}
	p.check_alternate_real()
	p.fft_transform_reverse(p.tables_reverse, p.real_inout, p.imag_inout)
	for i := 0; i < N; i += 2 {
		res_dbl[i] = p.real_inout[i+1]
		res_dbl[i+1] = p.imag_inout[i+1]
	}
	for i := 0; i < Ns2; i++ {
		Assert(cmplx.Abs(complex(p.real_inout[2*i+1], p.imag_inout[2*i+1])-res[i]) < 1e-20)
	}
	p.check_conjugate_cplx()
	return
}

func (p *NayukiFFTProcessor) executeReverseTorus322(a []Torus32) (res []complex128) {
	N := int(p.N)
	Ns2 := int(p.Ns2)
	_2N := int(p._2N)
	//real_inout := p.real_inout
	//imag_inout := p.imag_inout
	var _2pm33 double = 1. / double(int64(1)<<33)
	//int32_t* aa = (int32_t*) a;
	for i := 0; i < N; i++ {
		p.real_inout[i] = float64(a[i]) * _2pm33
	}
	for i := 0; i < N; i++ {
		p.real_inout[N+i] = -p.real_inout[i]
	}
	for i := 0; i < _2N; i++ {
		p.imag_inout[i] = 0
	}
	p.check_alternate_real()
	p.fft_transform_reverse(p.tables_reverse, p.real_inout, p.imag_inout)
	res = make([]complex128, N)
	for i := 0; i < Ns2; i++ {
		res[i] = complex(p.real_inout[2*i+1], p.imag_inout[2*i+1])
	}
	p.check_conjugate_cplx()
	return
}

func (p *NayukiFFTProcessor) executeDirectTorus322(a []complex128) (res []Torus32) {
	N := int(p.N)
	Ns2 := int(p.Ns2)
	_2N := int(p._2N)
	//real_inout := p.real_inout
	//imag_inout := p.imag_inout
	var _2p32 double = double(int64(1) << 32)
	var _1sN double = double(1) / double(N)
	//double* a_dbl=(double*) a;
	for i := 0; i < N; i++ {
		p.real_inout[2*i] = 0
	}
	for i := 0; i < N; i++ {
		p.imag_inout[2*i] = 0
	}
	for i := 0; i < Ns2; i++ {
		p.real_inout[2*i+1] = real(a[i])
	}
	for i := 0; i < Ns2; i++ {
		p.imag_inout[2*i+1] = imag(a[i])
	}
	for i := 0; i < Ns2; i++ {
		p.real_inout[_2N-1-2*i] = real(a[i])
	}
	for i := 0; i < Ns2; i++ {
		p.imag_inout[_2N-1-2*i] = -imag(a[i])
	}

	if debug {
		for i := 0; i < N; i++ {
			Assert(p.real_inout[2*i] == 0)
		}
		for i := 0; i < N; i++ {
			Assert(p.imag_inout[2*i] == 0)
		}
		for i := 0; i < Ns2; i++ {
			Assert(p.real_inout[2*i+1] == real(a[i]))
		}
		for i := 0; i < Ns2; i++ {
			Assert(p.imag_inout[2*i+1] == imag(a[i]))
		}
		for i := 0; i < Ns2; i++ {
			Assert(p.real_inout[_2N-1-2*i] == real(a[i]))
		}
		for i := 0; i < Ns2; i++ {
			Assert(p.imag_inout[_2N-1-2*i] == -imag(a[i]))
		}
		p.check_conjugate_cplx()
	}

	p.fft_transform(p.tables_direct, p.real_inout, p.imag_inout)
	res = make([]Torus32, N)
	for i := 0; i < N; i++ {
		res[i] = Torus32(int64(p.real_inout[i] * _1sN * _2p32))
	}
	//pas besoin du fmod... Torus32(int64_t(fmod(rev_out[i]*_1sN,1.)*_2p32));
	p.check_alternate_real()
	return
}
*/

/**
 * FFT functions
 */

func IntPolynomial_ifft(result *LagrangeHalfCPolynomial, p *IntPolynomial) {
	result.coefsC = fp1024_nayuki.executeReverseInt(p.Coefs)
}

func TorusPolynomial_ifft(result *LagrangeHalfCPolynomial, p *TorusPolynomial) {
	result.coefsC = fp1024_nayuki.executeReverseTorus32(p.CoefsT)
}

func TorusPolynomial_fft(result *TorusPolynomial, p *LagrangeHalfCPolynomial) {
	result.CoefsT = fp1024_nayuki.executeDirectTorus32(p.coefsC)
}

func fft_init(n int64) *FftTables {
	// Check size argument
	if n <= 0 || (n&(n-1)) != 0 {
		return nil // Error: Size is not a power of 2
	}
	/*
		if (n / 2 > SIZE_MAX / sizeof(double) || n > SIZE_MAX / sizeof(size_t)){
			return NULL;  // Error: Size is too large, which makes memory allocation impossible
		}
	*/

	// Allocate structure
	/*
		struct FftTables *tables = malloc(sizeof(struct FftTables));
		if (tables == NULL)
			return NULL;
		tables.n = n;

		// Allocate arrays
		tables.bit_reversed = malloc(n * sizeof(size_t));
		tables.cos_table = malloc(n / 2 * sizeof(double));
		tables.sin_table = malloc(n / 2 * sizeof(double));
		if (tables.bit_reversed == NULL || tables.cos_table == NULL || tables.sin_table == NULL) {
			free(tables.bit_reversed);
			free(tables.cos_table);
			free(tables.sin_table);
			free(tables);
			return NULL;
		}
	*/

	tables := &FftTables{
		n:            n,
		bit_reversed: make([]int64, n),
		cos_table:    make([]double, n/2),
		sin_table:    make([]double, n/2),
	}

	// Precompute values and store to tables
	//size_t i;
	levels := floor_log2(int64(n))
	for i := int64(0); i < n; i++ {
		tables.bit_reversed[i] = int64(reverse_bits(int64(i), uint32(levels)))
	}
	for i := int64(0); i < n/2; i++ {
		var angle double = 2. * math.Pi * double(i) / double(n)
		tables.cos_table[i] = math.Cos(angle)
		tables.sin_table[i] = math.Sin(angle)
	}
	return tables
}

// Performs a forward FFT in place on the given arrays. The length is given by the tables struct.
func (p *NayukiFFTProcessor) fft_transform(tbl *FftTables, real []double, imag []double) {
	n := tbl.n

	// Bit-reversed addressing permutation
	bitreversed := tbl.bit_reversed
	for i := int64(0); i < n; i++ {
		j := bitreversed[i]
		if i < int64(j) {
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
	costable := tbl.cos_table
	sintable := tbl.sin_table
	for size := int64(2); size <= n; size *= 2 {
		halfsize := size / 2
		tablestep := n / size
		for i := int64(0); i < n; i += size {
			j := i
			for k := int64(0); j < i+halfsize; k += tablestep {
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

// Returns the largest i such that 2^i <= n.
/*
func floor_log2(n int32) int32 {
	var result int32 = 0
	for ; n > 1; n /= 2 {
		result++
	}
	return result
}

// Returns the bit reversal of the n-bit unsigned integer x.
func reverse_bits(x int32, n uint32) int32 {
	var result int32 = 0
	for i := uint32(0); i < n; i++ {
		x >>= 1
		result = (result << 1) | (x & 1)
	}
	return int32(result)
}
*/

/////////////////////////////////////////////////////

type FftTablesUint64 struct {
	n            int64
	bit_reversed []int64
	trig_tables  []double
}

// Returns sin(2 * pi * i / n), for n that is a multiple of 4.
func accurate_sine(i int64, n int64) double {
	if n%4 != 0 {
		return 0.
	} else {
		var neg int32 = 0 // Boolean
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
func floor_log2(n int64) int32 {
	var result int32 = 0
	for ; n > 1; n /= 2 {
		result++
	}
	return result
}

// Returns the bit reversal of the n-bit unsigned integer x.
func reverse_bits(x int64, n uint32) int64 {
	var result int64 = 0
	for i := uint32(0); i < n; i++ {
		result = (result << 1) | (x & 1)
		x >>= 1
	}
	return result
}

// Returns a pointer to an opaque structure of FFT tables. n must be a power of 2 and n >= 4.
func fft_init_reverse(n int64) *FftTablesUint64 {
	// Check size argument
	if n < 4 || n > math.MaxUint32 || (n&(n-1)) != 0 {
		return nil // Error: Size is too small or is not a power of 2
	}

	//if (n - 4 > SIZE_MAX / sizeof(double) / 2 || n > SIZE_MAX / sizeof(size_t))
	//	return NULL;  // Error: Size is too large, which makes memory allocation impossible

	tables := &FftTablesUint64{
		n:            n,
		bit_reversed: make([]int64, n),
		//trig_tables:  make([]double, n-4),
		trig_tables: make([]double, n-4),
		//trig_tables: make([]double, n*2),
	}

	// Precompute bit reversal table
	levels := floor_log2(n)
	for i := int64(0); i < n; i++ {
		tables.bit_reversed[i] = reverse_bits(i, uint32(levels))
	}

	// Precompute the packed trigonometric table for each FFT internal level
	var k int64 = 0
	for size := int64(8); size <= n; size *= 2 {
		for i := int64(0); i < size/2; i += 4 {
			for j := int64(0); j < 4; j++ {
				tables.trig_tables[k] = accurate_sine(i+j+size/4, size) // Cosine
				k++
			}
			k = 0
			for j := int64(0); j < 4; j++ {
				tables.trig_tables[k] = -accurate_sine(i+j, size) // Sine
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
func (p *NayukiFFTProcessor) fft_transform_reverse(tbl *FftTablesUint64, real []double, imag []double) {
	//struct FftTables *tbl = (struct FftTables *)tables;
	n := tbl.n

	// Bit-reversed addressing permutation
	bitreversed := tbl.bit_reversed
	for i := int64(0); i < n; i++ {
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
		for i := int64(0); i < n; i += 2 {
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
		for i := int64(0); i < n; i += 4 {
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
	trigtables := tbl.trig_tables
	// for size := int64(8); size <= n; size <<= 1 {
	for size := int64(8); size <= int64(len(trigtables)); size <<= 1 {
		halfsize := size >> 1
		for i := int64(0); i < n; i += size {
			var j int64 = 0
			var off int64 = 0
			for ; j < halfsize; j += 4 {
				for k := int64(0); k < 4; k++ { // To simulate x86 AVX 4-vectors
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
