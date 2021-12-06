package main

import (
	"fmt"
)

func ExamplePlanC2C() {
	N := 8
	data := make([]complex64, N)

	n := []int{N}
	plan := PlanC2C(n, data, data, FORWARD, ESTIMATE)
	defer plan.Destroy()

	data[0] = 1
	fmt.Println(data)
	plan.Execute()
	fmt.Println(data)
	plan.Execute()
	fmt.Println(data)

	// Output:
	// [(1+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i)]
	// [(1+0i) (1+0i) (1+0i) (1+0i) (1+0i) (1+0i) (1+0i) (1+0i)]
	// [(8+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i) (0+0i)]
}
