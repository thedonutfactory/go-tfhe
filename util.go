package tfhe

import (
	"fmt"
	"math/big"
)

func Assert(condition bool) {
	if condition == false {
		panic("Assertion error")
	}
}

func tabs(count int, msg string) {
	for i := 0; i < count; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("%s\n", msg)
}

func tabsi(count int, msg int64) {
	for i := 0; i < count; i++ {
		fmt.Printf("\t")
	}
	fmt.Printf("%d\n", msg)
}

func floatToComplexSlice(arr []float64) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		res[i] = complex(v, 0)
	}
	return
}

func complexToFloatSlice(arr []complex128) (res []float64) {
	res = make([]float64, len(arr))
	for i, v := range arr {
		res[i] = real(v)
	}
	return
}

func castComplex(arr []int64) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		//res[i] = complex(T32tod(v), 0.)
		res[i] = complex(float64(v), 0.)
	}
	return
}

func castInt(arr []complex128) (res []int64) {
	res = make([]int64, len(arr))
	for i, v := range arr {
		res[i] = int64(int(real(v)))
	}
	return
}

func castTorus2(arr []complex128) (res []int64) {
	_2p32 := double(int(1) << 32)
	_1sN := double(1) / double(4)
	res = make([]int64, len(arr))
	for i, v := range arr {
		t := real(v) * _2p32 * _1sN
		fmt.Printf("%f -> %f, %d\n", real(v), t, Torus(int((t))))
		res[i] = int64(int64((real(v)) * _2p32 * _1sN))
	}
	return
}

func castTorus(arr []complex128) (res []int64) {
	res = make([]int64, len(arr))
	for i, v := range arr {

		//res[i] = int64(real(v))
		//res[i] = Torus(int(real(v)))
		//res[i] = int64(int(math.Round(real(v))))
		res[i] = Dtot32(real(v)) // int64(real(v))
	}
	return
}

func toBig(a []int64) (res []*big.Int) {
	res = make([]*big.Int, len(a))
	for i, v := range a {
		res[i] = big.NewInt(int64(v))
	}
	return
}

func fromBig(a []*big.Int) (res []int64) {
	res = make([]int64, len(a))
	for i, v := range a {
		res[i] = int64(v.Int64())
	}
	return
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Swap(a, b int) {
	b, a = a, b
}
