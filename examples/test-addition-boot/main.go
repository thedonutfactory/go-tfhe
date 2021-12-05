package main

import (
	"fmt"
	"os"
	"time"

	t "github.com/thedonutfactory/go-tfhe"
)

func fullAdderMUX(x []*t.LweSample, y []*t.LweSample, nbBits int, key *t.PublicKey, priv *t.PrivateKey) []*t.LweSample {
	inOutParams := priv.Params.InOutParams
	sum := t.NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := t.NewLweSampleArray(2, inOutParams)
	t.BootsSymEncrypt(carry[0], 0, priv) // first carry initialized to 0
	// temps
	temp := t.NewLweSampleArray(2, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = t.Xor(x[i], y[i], key) // temp = xi XOR yi
		sum[i] = t.Xor(temp[0], carry[0], key)

		// carry = MUX(xi XOR yi, carry(i-1), xi AND yi)
		temp[1] = t.And(x[i], y[i], key) // temp1 = xi AND yi
		carry[1] = t.Mux(temp[0], carry[0], temp[1], key)

		mess1 := t.BootsSymDecrypt(temp[0], priv)
		mess2 := t.BootsSymDecrypt(carry[0], priv)
		mess3 := t.BootsSymDecrypt(temp[1], priv)
		messmux := t.BootsSymDecrypt(carry[1], priv)

		tt := mess3
		if mess1 != 0 {
			tt = mess2
		}

		if messmux != tt {
			fmt.Printf("\tError[fullAdderMUX]: %d - %f - %f - %f - %f\n", i,
				t.TorusToDouble(t.LwePhase(temp[0], priv.LweKey)),
				t.TorusToDouble(t.LwePhase(carry[0], priv.LweKey)),
				t.TorusToDouble(t.LwePhase(temp[1], priv.LweKey)),
				t.TorusToDouble(t.LwePhase(carry[1], priv.LweKey)),
			)
		}

		carry[0] = t.Copy(carry[1], key)
	}
	sum[nbBits] = t.Copy(carry[1], key)
	return sum
}

func fullAdder(x []*t.LweSample, y []*t.LweSample, nbBits int, key *t.PublicKey, priv *t.PrivateKey) []*t.LweSample {
	inOutParams := priv.Params.InOutParams
	sum := t.NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := t.NewLweSampleArray(2, inOutParams)
	t.BootsSymEncrypt(carry[0], 0, priv) // first carry initialized to 0
	// temps
	temp := t.NewLweSampleArray(3, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = t.Xor(x[i], y[i], key) // temp = xi XOR yi
		sum[i] = t.Xor(temp[0], carry[0], key)

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		temp[1] = t.And(x[i], y[i], key)        // temp1 = xi AND yi
		temp[2] = t.And(carry[0], temp[0], key) // temp2 = carry AND temp
		carry[1] = t.Xor(temp[1], temp[2], key)
		carry[0] = t.Copy(carry[1], key)
	}
	sum[nbBits] = t.Copy(carry[0], key)
	return sum
}

func comparisonMUX(x []*t.LweSample, y []*t.LweSample, nbBits int, key *t.PublicKey, priv *t.PrivateKey) *t.LweSample {

	inOutParams := priv.Params.InOutParams
	// carries
	carry := t.NewLweSampleArray(2, inOutParams)
	t.BootsSymEncrypt(carry[0], 1, priv) // first carry initialized to 1

	for i := 0; i < nbBits; i++ {
		temp := t.Xor(x[i], y[i], key) // temp = xi XOR yi
		carry[1] = t.Mux(temp, y[i], carry[0], key)
		carry[0] = t.Copy(carry[1], key)
	}
	return t.Copy(carry[0], key)
}

<<<<<<< HEAD
func fromBool(x bool) int32 {
	if !x {
=======
func fromBool(x bool) int {
	if x == false {
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		return 0
	} else {
		return 1
	}
}

func toBool(x int) bool {
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

func decryptAndDisplayResult(sum []*t.LweSample, keyset *t.PrivateKey) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := t.BootsSymDecrypt(sum[i], keyset)
		fmt.Printf("%d ", messSum)
	}
	fmt.Println("]")
}

func keys(params *t.GateBootstrappingParameterSet) (*t.PublicKey, *t.PrivateKey) {
	var pubKey *t.PublicKey
	var privKey *t.PrivateKey
	if _, err := os.Stat("private.key"); err == nil {
		fmt.Println("------ Reading keys from file ------")
		privKey, _ = t.ReadPrivKey("private.key")
		pubKey, _ = t.ReadPubKey("public.key")

	} else {
		fmt.Println("------ Key Generation ------")
		// generate the keys
		pubKey, privKey = t.GenerateKeys(params)
		t.WritePrivKey(privKey, "private.key")
		t.WritePubKey(pubKey, "public.key")
	}
	return pubKey, privKey
}

func main() {
	const (
		nbBits   = 8
		nbTrials = 1
	)
	// generate params
<<<<<<< HEAD
	var minimumLambda int32 = 100
	params := t.NewDefaultGateBootstrappingParameters(minimumLambda)
=======
	var minimumLambda int = 100
	params := tfhe.NewDefaultGateBootstrappingParameters(minimumLambda)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
	inOutParams := params.InOutParams
	pubKey, privKey := keys(params)

	for trial := 0; trial < nbTrials; trial++ {

		xBits := toBits(22)
		yBits := toBits(33)

		// generate samples
		x := t.NewLweSampleArray(nbBits, inOutParams)
		y := t.NewLweSampleArray(nbBits, inOutParams)
		for i := 0; i < nbBits; i++ {
<<<<<<< HEAD
			//t.BootsSymEncrypt(x[i], rand.Int31()%2, keyset)
			//t.BootsSymEncrypt(y[i], rand.Int31()%2, keyset)
			t.BootsSymEncrypt(x[i], int32(xBits[i]), privKey)
			t.BootsSymEncrypt(y[i], int32(yBits[i]), privKey)
=======
			//tfhe.BootsSymEncrypt(x[i], rand.Int31()%2, keyset)
			//tfhe.BootsSymEncrypt(y[i], rand.Int31()%2, keyset)
			tfhe.BootsSymEncrypt(x[i], int(xBits[i]), keyset)
			tfhe.BootsSymEncrypt(y[i], int(yBits[i]), keyset)
>>>>>>> de69509 (factored out the use of explicit int32 in favour of int)
		}
		// output sum
		//sum := t.NewLweSampleArray(nbBits+1, inOutParams)

		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits addition circuit (FA in MUX version), trial %d\n", nbBits, trial)
		start := time.Now()
		sum := fullAdderMUX(x, y, nbBits, pubKey, privKey)
		duration := time.Since(start)

		decryptAndDisplayResult(sum, privKey)
		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Printf("finished Bootstrapping %d bits addition circuit (FA in MUX version)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		var messCarry int = 0
		for i := 0; i < nbBits; i++ {
			messX := t.BootsSymDecrypt(x[i], privKey)
			messY := t.BootsSymDecrypt(y[i], privKey)
			messSum := t.BootsSymDecrypt(sum[i], privKey)

			if messSum != (messX ^ messY ^ messCarry) {
				fmt.Printf("\tVerification Error %d, %f - %f - %f\n", i,
					t.TorusToDouble(t.LwePhase(x[i], privKey.LweKey)),
					t.TorusToDouble(t.LwePhase(y[i], privKey.LweKey)),
					t.TorusToDouble(t.LwePhase(sum[i], privKey.LweKey)),
				)
			}
			if messCarry != 0 {
				messCarry = fromBool(toBool(messX) || toBool(messY))
			} else {
				messCarry = fromBool(toBool(messX) && toBool(messY))
			}
		}
		messSum := t.BootsSymDecrypt(sum[nbBits], privKey)
		if messSum != messCarry {
			fmt.Printf("\tVerification Error - %d, %d bits\n", trial, nbBits)
		}

		// evaluate the addition circuit
		fmt.Printf("Starting Bootstrapping %d bits addition circuit (FA)...trial %d\n", nbBits, trial)
		start = time.Now()
		sum = fullAdder(x, y, nbBits, pubKey, privKey)
		duration = time.Since(start)
		decryptAndDisplayResult(sum, privKey)
		fmt.Printf("finished Bootstrappings %d bits addition circuit (FA)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int = 0
			for i := 0; i < nbBits; i++ {
				messX := t.BootsSymDecrypt(x[i], privKey)
				messY := t.BootsSymDecrypt(y[i], privKey)
				messSum := t.BootsSymDecrypt(sum[i], privKey)

				if messSum != (messX ^ messY ^ messCarry) {
					fmt.Printf("\tVerification Error - %d, %d\n", trial, i)
				}

				if messCarry != 0 {
					messCarry = fromBool(toBool(messX) || toBool(messY))
				} else {
					messCarry = fromBool(toBool(messX) && toBool(messY))
				}
			}
			messSum := t.BootsSymDecrypt(sum[nbBits], privKey)
			if messSum != messCarry {
				fmt.Printf("\tVerification Error - %d, %d\n", trial, nbBits)
			}
		}

		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits comparison, trial %d\n", nbBits, trial)
		start = time.Now()
		comp := comparisonMUX(x, y, nbBits, pubKey, privKey)
		duration = time.Since(start)
		fmt.Printf("finished Bootstrappings %d bits comparison\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int = 1
			for i := 0; i < nbBits; i++ {
				messX := t.BootsSymDecrypt(x[i], privKey)
				messY := t.BootsSymDecrypt(y[i], privKey)
				if messX^messY != 0 {
					messCarry = messY
				}
			}
			messComp := t.BootsSymDecrypt(comp, privKey)
			if messComp != messCarry {
				fmt.Printf("\tVerification Error %d, %d\n", trial, nbBits)
			}
		}
	}
}
