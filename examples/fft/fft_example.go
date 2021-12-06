package main

import (
	"fmt"

	"github.com/takatoh/fft"
)

func main33() {
	x0 := []float64{
		5,
		32,
		38,
		-33,
		-19,
		-10,
		1,
		-8,
		-20,
		10,
		-1,
		4,
		11,
		-1,
		-7,
		-2,
	}
	n := len(x0)
	x := make([]complex128, n)
	for k := 0; k < n; k++ {
		x[k] = complex(x0[k], 0.0)
	}

	y := fft.FFT(x, n)
	z := fft.IFFT(y, n)

	fmt.Println(" K   DATA  FOURIER TRANSFORM  INVERSE TRANSFORM")
	for k := 0; k < n; k++ {
		fmt.Printf("%2d %6.1f  %8.3f%8.3f   %8.3f%8.3f\n",
			k, x0[k], real(y[k]), imag(y[k]), real(z[k]), imag(z[k]))
	}
}
