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
