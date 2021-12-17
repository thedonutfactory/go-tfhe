package main

import "fmt"

func doSomething(input int32) int32 {
	if input > 0 {
		// do stuff
		finalValue := 2 + 4 + 8*input
		return finalValue
	} else {
		// do other stuff
		finalValue := 1 + input
		return finalValue
	}
}

func doSomething2(input int32) int32 {
	// do stuff..regardless
	finalValue1 := 2 + 4 + 8*input

	// do other stuff..regardless
	finalValue2 := 1 + input

	//var resultMask bool = input > 0 // 1 if should use finalValue1.
	var resultMask int32 = 0
	if input > 0 {
		resultMask = 1
	}

	return finalValue1*resultMask + finalValue2*(1-resultMask)
}

func main() {
	fmt.Println(doSomething(32))
	fmt.Println(doSomething2(32))

	fmt.Println(doSomething(-32))
	fmt.Println(doSomething2(-32))
}
