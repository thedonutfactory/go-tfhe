package main

import (
	"fmt"
	"math"
)

type Torus = int64
type double = float64

// https://cp-algorithms.com/algebra/fft.html#toc-tgt-1
func swap(a, b interface{}) {
	b, a = a, b
}

func fft(a []complex128, invert bool) []complex128 {
	n := len(a)

	var j int = 0
	for i := 1; i < n; i++ {
		bit := n >> 1
		for ; j&bit > 0; bit >>= 1 {
			j ^= bit
		}
		j ^= bit

		if i < j {
			swap(a[i], a[j])
		}
	}

	for len := 2; len <= n; len <<= 1 {

		var ang float64
		if invert {
			ang = 2. * math.Pi / float64(len) * -1.
		} else {
			ang = 2. * math.Pi / float64(len)
		}

		wlen := complex(math.Cos(ang), math.Sin(ang))
		for i := 0; i < n; i += len {
			w := complex(1, 0.)
			for j := 0; j < len/2; j++ {
				u := a[i+j]
				v := a[i+j+len/2] * w
				a[i+j] = u + v
				a[i+j+len/2] = u - v
				w *= wlen
			}
		}
	}

	if invert {
		cn := complex(float64(n), 0)
		for i := 0; i < len(a); i++ {
			//a[i] = a[i] / complex(float64(n), 0)
			a[i] /= cn
		}

		//for (cd & x : a)
		//    x /= n;
	}
	return a
}

func mulfft(a []complex128) []complex128 {
	n := len(a)
	for i := 0; i < n; i++ {
		a = append(a, 0)
	}
	return fft(a, true)
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

func multiply(a, b []int64) []int64 {
	//vector<cd> fa(a.begin(), a.end()), fb(b.begin(), b.end());

	fa := []complex128{
		complex(float64(a[0]), 0.),
		complex(float64(a[len(a)-1]), 0.),
	}

	fb := []complex128{
		complex(float64(b[0]), 0.),
		complex(float64(b[len(b)-1]), 0.),
	}

	var n int = 1
	for n < len(a)+len(b) {
		n <<= 1
	}

	fa = resize(fa, n)
	fb = resize(fb, n)

	fft(fa, false)
	fft(fb, false)
	for i := 0; i < n; i++ {
		fa[i] *= fb[i]
	}
	fft(fa, true)

	result := make([]int64, n)
	for i := 0; i < n; i++ {
		result[i] = int64(int(math.Round(real(fa[i]))))
		//result[i] = Dtot32(real(fa[i])) // int(math.Round(real(fa[i])))
	}
	return result
}

const two32 int64 = int64(1) << 32 // 2^32
var two32Double double = math.Pow(2, 32)

// from double to Torus - float64 to int64 conversion
func Dtot32(d double) Torus {
	return Torus(math.Round(math.Mod(d, 1) * math.Pow(2, 32)))
}

// from Torus to double
func T32tod(x Torus) double {
	return double(x) / math.Pow(2, 32)
}

func main() {

	/*
		0, -909722663
		0, 1748652883
		0, 1571540080
		0, 2136454616
	*/

	a := []int64{-1865008400, 470211269, -689632771, 1115438162}
	b := []int64{156091742, 1899894088, -1210297292, -1557125705}

	c := multiply(a, b)
	fmt.Print("Vector c:\n")
	for i := 0; i < len(c); i++ {
		fmt.Println(c[i])
	}

}
