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

	a := []C.long{-1865008400, 470211269, -689632771, 1115438162}
	b := []C.long{156091742, 1899894088, -1210297292, -1557125705}
	var c []C.long = make([]C.long, len(a))
	MultiplyWithCuda(a, b, c, len(a))
	fmt.Println(c)
	fmt.Println(convertToTorus(c))
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

const two32 int64 = int64(1) << 32 // 2^32
// from double to Torus32
func dtot32(d int64) int {
	return int(int64((d - int64(d)) * two32))
}

/*
// from Torus32 to double
EXPORT double t32tod(Torus32 x) {
    return double(x)/_two32_double;
}
*/
