package main

import (
	"fmt"
)

func branchedIfElse(input float64) float64 {
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

func unBranchedIfElse(input float64) float64 {
	// do stuff..regardless
	finalValue1 := 2 + 4 + 8*input

	// do other stuff..regardless
	finalValue2 := 1 + input

	// result mask
	resultMask := 0. //= input > 0 // 1 if should use finalValue1.
	if input > 0 {
		resultMask = 1.
	}

	return finalValue1*resultMask + finalValue2*(1-resultMask)
}

func ifElseUnbranched(input float64, ifCondition func(input float64) int, trueBranch func(input float64) float64, falseBranch func(input float64) float64) float64 {
	resultMask := ifCondition(input)
	return trueBranch(input)*float64(resultMask) + falseBranch(input)*float64(1-resultMask)
}

func main() {
	inp := 42.
	fmt.Printf("branchedIfElse: %f\n", branchedIfElse(inp))
	fmt.Printf("unBranchedIfElse: %f\n", unBranchedIfElse(inp))

	ifCondition := func(input float64) int {
		// result mask
		resultMask := 0 //= input > 0 // 1 if should use finalValue1.
		if input > 0 {
			resultMask = 1
		}
		return resultMask
	}

	trueBranch := func(input float64) float64 {
		return 2 + 4 + 8*input
	}
	falseBranch := func(input float64) float64 {
		return 1 + input
	}

	fmt.Printf("unBranchedIfElse: %f\n", ifElseUnbranched(ifCondition, trueBranch, falseBranch))

}
