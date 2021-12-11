package main

import (
	"fmt"
	"math/rand"

	t "github.com/thedonutfactory/go-tfhe"
)

func main() {
	kNumTests := 5

	//SetSeed();  // set random seed

	pt := t.NewPtxtArray(2)
	ct := t.NewCtxtArray(2)

	fmt.Println("------ Key Generation ------")
	_, pri_key := t.KeyGen()

	//fmt.Printf("%v, %v", pub_key, pri_key)

	correct := true
	for i := 0; i < kNumTests; i++ {
		pt[0].Message = uint32(rand.Int31() % t.KPtxtSpace)
		t.Encrypt(ct[0], pt[0], pri_key)
		t.Decrypt(pt[1], ct[0], pri_key)
		if pt[1].Message != pt[0].Message {
			correct = false
			//break
		}
	}
	if correct {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}
}
