package main

import (
	"fmt"
	"math"

	"github.com/mjibson/go-dsp/fft"
)

type Torus = int64
type double = float64

const _two32 int64 = int64(1) << 32 // 2^32
var _two32_double double = math.Pow(2, 32)

// from double to Torus - float64 to int64 conversion
func Dtot32(d double) Torus {
	return Torus(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

// from Torus to double
func T32tod(x Torus) double {
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

func invfft(a []complex128) []complex128 {
	return fft.IFFT(a)
}

func mulfft2(a []complex128) []complex128 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return Fft(a)
}

func mulfft(a []complex128) []complex128 {
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

func castComplex(arr []int64) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		res[i] = complex(float64(v), 0.)
	}
	return
}

func castTorus(arr []complex128) (res []int64) {
	_2p32 := double(int(1) << 32)
	_1sN := double(1) / double(4)
	//res[i]=Torus(int64_t(out[i]*_1sN*_2p32))
	res = make([]int64, len(arr))
	for i, v := range arr {
		t := real(v) * _2p32 * _1sN
		fmt.Printf("%f -> %f, %d\n", real(v), t, Torus(int((t))))
		res[i] = int64(real(v)) //int64(int(real(v))) // Dtot32(real(v)) // int64(real(v))
	}
	return
}

func multiply(a, b []int64) []int64 {
	x := mulfft(castComplex(a))
	y := mulfft(castComplex(b))
	c := mult(x, y)
	return castTorus(invfft(c))
}

func main() {

	/*
		-909722663
		1748652883
		1571540080
		2136454616
	*/

	//a := []int64{-1865008400, 470211269, -689632771, 1115438162}
	//b := []int64{156091742, 1899894088, -1210297292, -1557125705}

	a := []int64{9, -10, 7, 6}
	b := []int64{-5, 4, 0, -2}

	c := multiply(a, b)
	fmt.Print("Vector c:\n")
	for i := 0; i < len(c); i++ {
		fmt.Println(c[i])
	}

}
