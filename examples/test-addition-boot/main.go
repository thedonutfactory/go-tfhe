package main

import (
	"fmt"
	"os"
	"time"

	t "github.com/thedonutfactory/go-tfhe"
)

func fullAdder(sum []*t.Ctxt, x []*t.Ctxt, y []*t.Ctxt, nbBits int, pubKey *t.PubKey, privKey *t.PriKey) {
	// carries
	carry := t.NewCtxtArray(2)
	t.Encrypt(carry[0], t.NewPtxtInit(0), privKey)
	// temps
	temp := t.NewCtxtArray(3)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		t.Xor(temp[0], x[i], y[i], pubKey)
		t.Xor(sum[i], temp[0], carry[0], pubKey)

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		t.And(temp[1], x[i], y[i], pubKey)        // temp1 = xi AND yi
		t.And(temp[2], carry[0], temp[0], pubKey) // temp2 = carry AND temp
		t.Xor(carry[1], temp[1], temp[2], pubKey)
		t.Copy(carry[0], carry[1])
	}
	t.Copy(sum[nbBits], carry[0])
}

func verifyAdder(sum, x, y []int, nbBits int) {
	carry := []int{0, 0}
	temp := []int{0, 0, 0}
	for i := 0; i < nbBits; i++ {
		temp[0] = x[i] ^ y[i]
		sum[i] = temp[0] ^ carry[0]

		temp[1] = x[i] & y[i]
		temp[2] = carry[0] & temp[0]
		carry[1] = temp[1] ^ temp[2]
		carry[0] = carry[1]
	}
	sum[nbBits] = carry[0]
}

func fromBool(x bool) uint32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

func toBool(x uint32) bool {
	if x == 0 {
		return false
	} else {
		return true
	}
}

func toBits(val int) []int {
	l := make([]int, 8)

	l[0] = val & 0x1
	l[1] = (val & 0x2) >> 1
	l[2] = (val & 0x4) >> 2
	l[3] = (val & 0x8) >> 3

	l[4] = (val & 0x16) >> 4
	l[5] = (val & 0x32) >> 5
	l[6] = (val & 0x64) >> 6
	l[7] = (val & 0x128) >> 7

	return l
}

func decryptAndDisplayResult(sum []*t.Ctxt, pri_key *t.PriKey) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := t.NewPtxt()
		t.Decrypt(messSum, sum[i], pri_key)
		fmt.Printf("%d ", messSum.Message)
	}
	fmt.Println("]")
}

func main() {
	const (
		nbBits   = 4
		nbTrials = 1
	)
	var (
		pub_key *t.PubKey
		pri_key *t.PriKey
	)

	if _, err := os.Stat("private.key"); err == nil {
		pri_key, _ = t.ReadPrivKey("private.key")
		pub_key, _ = t.ReadPubKey("public.key")

	} else {
		fmt.Println("------ Key Generation ------")
		pub_key, pri_key = t.KeyGen()
		t.WritePrivKey(pri_key, "private.key")
		t.WritePubKey(pub_key, "public.key")
	}

	for trial := 0; trial < nbTrials; trial++ {

		xBits := toBits(2)
		yBits := toBits(3)

		vsum := make([]int, nbBits+1)
		verifyAdder(vsum, xBits, yBits, nbBits)
		fmt.Println("Expected:", vsum)

		// generate samples
		x := t.NewCtxtArray(nbBits)
		y := t.NewCtxtArray(nbBits)
		for i := 0; i < nbBits; i++ {
			//tfhe.BootsSymEncrypt(x[i], rand.Int31()%2, keyset)
			//tfhe.BootsSymEncrypt(y[i], rand.Int31()%2, keyset)
			//tfhe.BootsSymEncrypt(x[i], int32(xBits[i]), keyset)
			//tfhe.BootsSymEncrypt(y[i], int32(yBits[i]), keyset)

			ptxt := t.NewPtxtInit(uint32(xBits[i]))
			t.Encrypt(x[i], ptxt, pri_key)

			ptxt2 := t.NewPtxtInit(uint32(yBits[i]))
			t.Encrypt(y[i], ptxt2, pri_key)
		}
		// output sum
		sum := t.NewCtxtArray(nbBits + 1)

		// evaluate the addition circuit
		fmt.Printf("Starting Bootstrapping %d bits addition circuit (FA)...trial %d\n", nbBits, trial)
		start := time.Now()
		fullAdder(sum, x, y, nbBits, pub_key, pri_key)
		duration := time.Since(start)
		fmt.Print("Actual: ")
		decryptAndDisplayResult(sum, pri_key)
		fmt.Printf("finished Bootstrappings %d bits addition circuit (FA)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry uint32 = 0
			messX, messY, messSum := t.NewPtxt(), t.NewPtxt(), t.NewPtxt()
			for i := 0; i < nbBits; i++ {
				t.Decrypt(messX, x[i], pri_key)
				t.Decrypt(messY, y[i], pri_key)
				t.Decrypt(messSum, sum[i], pri_key)

				if messSum.Message != (messX.Message ^ messY.Message ^ messCarry) {
					fmt.Printf("\tVerification Error - trial: %d, bit: %d\n", trial, i)
				}

				if messCarry != 0 {
					messCarry = fromBool(toBool(messX.Message) || toBool(messY.Message))
				} else {
					messCarry = fromBool(toBool(messX.Message) && toBool(messY.Message))
				}
				//messCarry = messCarry ? (messX || messY) : (messX && messY);
			}
			t.Decrypt(messSum, sum[nbBits], pri_key)
			if messSum.Message != messCarry {
				fmt.Printf("\tVerification Error - %d, %d\n", trial, nbBits)
			}
		}

	}
}
