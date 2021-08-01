package tfhe

import "fmt"

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

func tabsi(count int, msg int32) {
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

func castComplex(arr []int32) (res []complex128) {
	res = make([]complex128, len(arr))
	for i, v := range arr {
		//res[i] = complex(T32tod(v), 0.)
		res[i] = complex(float64(v), 0.)
	}
	return
}

func castTorus(arr []complex128) (res []int32) {
	res = make([]int32, len(arr))
	for i, v := range arr {
		res[i] = int32(real(v)) // Dtot32(real(v)) // int32(real(v))
	}
	return
}
