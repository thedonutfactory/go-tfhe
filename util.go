package tfhe

func Assert(condition bool) {
	if condition == false {
		panic("Assertion error")
	}
}
