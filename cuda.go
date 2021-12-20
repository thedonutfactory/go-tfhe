package tfhe

/*
void multiplyWithCuda(long *c, const long *a, const long *b, unsigned int size);
#cgo LDFLAGS: -L. -L./ -lmultkernel
*/
import "C"

import (
	"fmt"
)

func MultiplyWithCuda(a []C.long, b []C.long, c []C.long, size int) {
	C.multiplyWithCuda(&c[0], &a[0], &b[0], C.uint(size))
}

// GPU
func torusPolynomialMultCuda(result *TorusPolynomial, poly1 *IntPolynomial, poly2 *TorusPolynomial) {

	/*
		for i := 0; i < 10; i++ {
			a := []C.long{9, -10, 7, 6}
			b := []C.long{-5, 4, 0, -2}
			var c []C.long = make([]C.long, len(a))
			MultiplyWithCuda(a, b, c, len(a))
			fmt.Println(c)
			fmt.Println(convertToTorus(c))
		}
	*/

	a := convert(poly1.Coefs)
	b := convert(poly2.Coefs)
	c := make([]C.long, len(a))
	//Maxmul(poly1.Coefs, poly2.CoefsT, result.CoefsT, 32)
	MultiplyWithCuda(a, b, c, len(a))

	n := len(a)
	res := make([]int32, n)
	for i := range res {
		t := c[i] - c[n+i]
		result.Coefs[i] = int32(int64(t))
	}
	//return res

	//result.Coefs = convertToTorus(c)
	//fmt.Println(poly1.Coefs)
	fmt.Println(a, b)
	fmt.Println(result.Coefs)

}

func convert(ar []int32) []C.long {
	rval := make([]C.long, len(ar))
	var v int32
	var i int
	for i, v = range ar {
		rval[i] = C.long(v)
	}
	return rval
}

func convertToTorus(ar []C.long) []Torus32 {
	newar := make([]Torus32, len(ar))
	var v C.long
	var i int
	for i, v = range ar {
		newar[i] = int32(v)
	}
	return newar
}
