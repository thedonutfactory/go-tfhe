package main

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

func main() {
	//in := []C.float{1.23, 4.56}
	//C.test(&in[0]) // C 1.230000 4.560000

	a := []C.long{9, -10, 7, 6}
	b := []C.long{-5, 4, 0, -2}
	var c []C.long = make([]C.long, len(a))
	MultiplyWithCuda(a, b, c, len(a))
	fmt.Println(c)
	fmt.Println(convertToTorus(c))
}

func convertToTorus(ar []C.long) []int64 {
	newar := make([]int64, len(ar))
	var v C.long
	var i int
	for i, v = range ar {
		newar[i] = int64(v)
	}
	return newar
}
