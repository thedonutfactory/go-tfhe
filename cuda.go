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

	//TorusPolynomialMulR(result, poly1, poly2)

	//a := []C.float{-1, 2, 4, 0, 5, 3, 6, 2, 1}
	//b := []C.float{3, 0, 2, 3, 4, 5, 4, 7, 2}
	//var c []C.float = make([]C.float, 9)

	a := convert(poly1.Coefs)
	b := convert(poly2.CoefsT)
	c := make([]C.long, len(a))
	//Maxmul(poly1.Coefs, poly2.CoefsT, result.CoefsT, 32)
	MultiplyWithCuda(a, b, c, len(a))
	result.CoefsT = convertToTorus(c)
	//fmt.Println(poly1.Coefs)
	fmt.Println(a, b)
	fmt.Println(result.CoefsT)
}

func convert(ar []int) []C.long {
	rval := make([]C.long, len(ar))
	var v int
	var i int
	for i, v = range ar {
		rval[i] = C.long(v)
	}
	return rval
}

func convertToTorus(ar []C.long) []int {
	newar := make([]int, len(ar))
	var v C.long
	var i int
	for i, v = range ar {
		newar[i] = int(v)
	}
	return newar
}
