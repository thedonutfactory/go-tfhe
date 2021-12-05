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

func tabsi(count int, msg int) {
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

func castComplex(arr []int) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		//res[i] = complex(T32tod(v), 0.)
		res[i] = complex(float64(v), 0.)
	}
	return
}

func castInt(arr []complex128) (res []int) {
	res = make([]int, len(arr))
	for i, v := range arr {
		res[i] = int(int(real(v)))
	}
	return
}

func castTorus2(arr []complex128) (res []int) {
	_2p32 := double(int(1) << 32)
	_1sN := double(1) / double(4)
	res = make([]int, len(arr))
	for i, v := range arr {
		t := real(v) * _2p32 * _1sN
		fmt.Printf("%f -> %f, %d\n", real(v), t, Torus(int((t))))
		res[i] = int(int64((real(v)) * _2p32 * _1sN))
	}
	return
}

func castTorus(arr []complex128) (res []Torus) {
	res = make([]int, len(arr))
	for i, v := range arr {

		//res[i] = int32(real(v))
		//res[i] = Torus32(int(real(v)))
		//res[i] = int32(int(math.Round(real(v))))
		res[i] = DoubleToTorus(real(v)) // int32(real(v))
	}
	return
}

func toBig(a []int) (res []*big.Int) {
	res = make([]*big.Int, len(a))
	for i, v := range a {
		res[i] = big.NewInt(int64(v))
	}
	return
}

func fromBig(a []*big.Int) (res []int) {
	res = make([]int, len(a))
	for i, v := range a {
		res[i] = int(v.Int64())
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
