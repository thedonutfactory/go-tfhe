package main

import (
	"fmt"

	t "github.com/thedonutfactory/go-tfhe"
)

func NandCheck(out, in0, in1 *t.Ptxt) {
	out.Message = 1 - in0.Message*in1.Message
}

func main() {
	kNumTests := 50

	//SetSeed();  // set random seed

	pt := t.NewPtxtArray(2)
	ct := t.NewCtxtArray(2)

	fmt.Println("------ Key Generation ------")
	pub_key, pri_key := t.KeyGen()

	//fmt.Printf("%v, %v", pub_key, pri_key)

	correct := true
	for i := 0; i < kNumTests; i++ {
		pt[0].Message = 0 //uint32(rand.Int31() % t.KPtxtSpace)
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

	fmt.Println("------ Test NAND Gate ------")
	kNumTests = 4
	correct = true
	for i := 0; i < kNumTests; i++ {
		pt[0].Message = 1 //uint32(rand.Int31() % t.KPtxtSpace)
		pt[1].Message = 0 //uint32(rand.Int31() % t.KPtxtSpace)
		t.Encrypt(ct[0], pt[0], pri_key)
		t.Encrypt(ct[1], pt[1], pri_key)
		t.Nand(ct[0], ct[0], ct[1], pub_key)
		NandCheck(pt[1], pt[0], pt[1])
		t.Decrypt(pt[0], ct[0], pri_key)
		if pt[1].Message != pt[0].Message {
			correct = false
			break
		}
	}
	if correct {
		fmt.Println("PASS")
	} else {
		fmt.Println("FAIL")
	}
}
