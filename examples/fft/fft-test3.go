package main

import (
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

const N = 4
const Ns2 = 2

type Torus32 = int32
type double = float64

const _two32 int64 = int64(1) << 32 // 2^32
var _two32_double double = math.Pow(2, 32)

// from double to Torus32 - float64 to int32 conversion
func Dtot32(d double) Torus32 {
	return Torus32(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

// from Torus32 to double
func T32tod(x Torus32) double {
	return double(x) / math.Pow(2, 32)
}

// https://github.com/Pokelover166/FFT-polynomial-multiplication/blob/master/main.cpp

func Fft(a []complex128) []complex128 {
	n := len(a)
	if n == 1 {
		return []complex128{1, a[0]}
	}

	w := make([]complex128, n)
	for i := 0; i < n; i++ {
		alpha := 2. * math.Pi * float64(i) / float64(n) // M_PI * float64(2*i/n)
		w[i] = complex(math.Cos(alpha), math.Sin(alpha))
		//fmt.Println(w[i])
	}

	A0 := make([]complex128, n/2)
	A1 := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		A0[i] = a[i*2]   ///Even coefficients
		A1[i] = a[i*2+1] ///Odd coefficients
	}

	y0 := Fft(A0)
	y1 := Fft(A1)
	y := make([]complex128, n)
	for k := 0; k < n/2; k++ {
		y[k] = y0[k] + w[k]*y1[k]
		y[k+n/2] = y0[k] - w[k]*y1[k]

		//fmt.Printf("y0[k]: %.1f, w[k]: %.1f, y1[k]: %.1f w[k]*y1[k]: %.1f\n", y0[k], w[k], y1[k], w[k]*y1[k])
		fmt.Println(y0[k], w[k], y1[k], w[k]*y1[k])

		//fmt.Println(y[k])
		//fmt.Println(y[k+n/2])
	}
	fmt.Println()
	return y

}

func Ifft(a []complex128) []complex128 {
	n := len(a)
	if n == 1 {
		return []complex128{1, a[0]}
	}

	w := make([]complex128, n)
	for i := 0; i < n; i++ {
		alpha := 2. * math.Pi * float64(i) / float64(n) //2 * M_PI * i / n
		w[i] = complex(math.Cos(alpha), math.Sin(alpha))
	}
	A0 := make([]complex128, n/2)
	A1 := make([]complex128, n/2)
	for i := 0; i < n/2; i++ {
		A0[i] = a[i*2]
		A1[i] = a[i*2+1]
	}
	y0 := Fft(A0)
	y1 := Fft(A1)
	y := make([]complex128, n)
	for k := 0; k < n/2; k++ {
		y[k] = y0[k] + y1[k]/w[k] ///w[k]^-1
		y[k+n/2] = y0[k] - y1[k]/w[k]
	}
	return y
}

func invfft2(a []complex128) []complex128 {
	a = Ifft(a)
	n := len(a)
	for i := 0; i < n; i++ {
		a[i] = a[i] / complex(float64(n), 0)
	}
	return a
}

/*
void FFT_Processor_fftw::execute_reverse_torus32(cplx* res, const Torus32* a) {
    static const double _2pm33 = 1./double(INT64_C(1)<<33);
    int32_t* aa = (int32_t*) a;
    cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<N; i++) rev_in[i]=aa[i]*_2pm33;
    for (int32_t i=0; i<N; i++) rev_in[N+i]=-rev_in[i];
    fftw_execute(rev_p);
    for (int32_t i=0; i<Ns2; i++) res[i]=rev_out_cplx[2*i+1];
    for (int32_t i=0; i<=Ns2; i++) assert(abs(rev_out_cplx[2*i])<1e-20);
}
*/

func revTorus(a []Torus32) []complex128 {
	//N := len(a)
	//Ns2 := len(a) / 2
	_2pm33 := 1. / double(int64(1)<<33)
	rev_in := make([]complex128, N*2)

	for i := 0; i < N; i++ {
		rev_in[i] = complex(float64(a[i])*_2pm33, 0.)
	}
	for i := 0; i < N; i++ {
		rev_in[N+i] = -rev_in[i]
	}

	rev_out_cplx := fft.FFT(rev_in)

	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}

	// assert
	for i := 0; i <= Ns2; i++ {
		if math.Abs(real(rev_out_cplx[2*i])) >= 1e-20 {
			panic("err")
		}
	}

	return res
}

/*
void FFT_Processor_fftw::execute_reverse_int(cplx* res, const int* a) {
    cplx* rev_out_cplx = (cplx*) rev_out; //fftw_complex and cplx are layout-compatible
    for (int32_t i=0; i<N; i++) rev_in[i]=a[i]/2.;
    for (int32_t i=0; i<N; i++) rev_in[N+i]=-rev_in[i];
    fftw_execute(rev_p);
    for (int32_t i=0; i<Ns2; i++) res[i]=rev_out_cplx[2*i+1];
    for (int32_t i=0; i<=Ns2; i++) assert(abs(rev_out_cplx[2*i])<1e-20);
}
*/

func revInt(a []int32) []complex128 {
	//N := len(a)
	//Ns2 := len(a) / 2
	rev_in := make([]complex128, N*2)

	for i := 0; i < N; i++ {
		rev_in[i] = complex(float64(a[i])/2., 0.)
	}
	for i := 0; i < N; i++ {
		rev_in[N+i] = -rev_in[i]
	}

	rev_out_cplx := fft.FFT(rev_in)

	res := make([]complex128, Ns2)
	for i := 0; i < Ns2; i++ {
		res[i] = rev_out_cplx[2*i+1]
	}

	// assert
	for i := 0; i <= Ns2; i++ {
		if math.Abs(real(rev_out_cplx[2*i])) >= 1e-20 {
			panic("err")
		}
	}

	return res
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

func dirTorus(a []complex128) []Torus32 {
	//N := len(a)
	//Ns2 := len(a) / 2
	_2p32 := double(int64(1) << 32)
	_1sN := double(1) / double(N)

	in_cplx := make([]complex128, N+1)
	for i := 0; i <= Ns2; i++ {
		in_cplx[2*i] = 0
	}
	for i := 0; i < Ns2; i++ {
		in_cplx[2*i+1] = a[i]
	}

	out := fft.FFT(in_cplx)

	res := make([]Torus32, N)
	for i := 0; i < N; i++ {
		fmt.Printf("%f => %f \n", real(out[i]), real(out[i])*_1sN*_2p32)
		res[i] = Torus32(math.Round(real(out[i]) * _1sN * _2p32))
	}
	return res

	/*
		n := len(a)
		for i := 0; i < n; i++ {
			a = append(a, 0)
		}
		return fft.FFT(a)
	*/
}

func mulfft3(a []complex128) []complex128 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return fft.FFT(a)
}

func mult(a, b []complex128) []complex128 {
	n := max(len(a), len(b))
	c := make([]complex128, n)
	for i := 0; i < n; i++ {
		c[i] = a[i] * b[i]
	}
	return c
}

// Max returns the larger of x or y.
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func resize(a []complex128, n int) []complex128 {
	y := n - len(a)
	for i := 0; i < y; i++ {
		a = append(a, 0)
	}
	return a
}

func castComplex(arr []int32) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		res[i] = complex(float64(v), 0.)
	}
	return
}

func castTorus(arr []complex128) (res []int32) {
	_2p32 := double(int(1) << 32)
	_1sN := double(1) / double(4)
	//res[i]=Torus32(int64_t(out[i]*_1sN*_2p32))
	res = make([]int32, len(arr))
	for i, v := range arr {
		t := real(v) * _2p32 * _1sN
		fmt.Printf("%f -> %f, %d\n", real(v), t, Torus32(int((t))))
		res[i] = int32(real(v)) //int32(int(real(v))) // Dtot32(real(v)) // int32(real(v))
	}
	return
}

/*
	IntPolynomial_ifft(tmp[0], poly1)
	TorusPolynomial_ifft(tmp[1], poly2)
	LagrangeHalfCPolynomialMul(tmp[2], tmp[0], tmp[1])
	TorusPolynomial_fft(result, tmp[2])
*/

func multi(a, b []int32) []int32 {
	x := revInt(a)
	y := revTorus(b)
	c := mult(x, y)
	return dirTorus(c) //castTorus(mulfft3(c))
}

func main() {

	/*
		-909722663
		1748652883
		1571540080
		2136454616
	*/

	//a := []int32{-1865008400, 470211269, -689632771, 1115438162}
	//b := []int32{156091742, 1899894088, -1210297292, -1557125705}

	// -89, 100, -63, -20
	a := []int32{9, -10, 7, 6}
	b := []int32{-5, 4, 0, -2}

	c := multi(a, b)
	fmt.Print("Vector c:\n")
	for i := 0; i < len(c); i++ {
		fmt.Println(c[i])
	}

}
