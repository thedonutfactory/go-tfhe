package main

import (
	"fmt"
	"os"
	"time"

	. "github.com/thedonutfactory/go-tfhe/core"
	. "github.com/thedonutfactory/go-tfhe/gates"
	. "github.com/thedonutfactory/go-tfhe/io"
	. "github.com/thedonutfactory/go-tfhe/types"
)

func fullAdderMUX(x []*LweSample, y []*LweSample, nbBits int, key *PublicKey, priv *PrivateKey) []*LweSample {
	inOutParams := priv.Params.InOutParams
	sum := NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	BootsSymEncrypt(carry[0], 0, priv) // first carry initialized to 0
	// temps
	temp := NewLweSampleArray(2, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = Xor(x[i], y[i], key) // temp = xi XOR yi
		sum[i] = Xor(temp[0], carry[0], key)

		// carry = MUX(xi XOR yi, carry(i-1), xi AND yi)
		temp[1] = And(x[i], y[i], key) // temp1 = xi AND yi
		carry[1] = Mux(temp[0], carry[0], temp[1], key)

		mess1 := BootsSymDecrypt(temp[0], priv)
		mess2 := BootsSymDecrypt(carry[0], priv)
		mess3 := BootsSymDecrypt(temp[1], priv)
		messmux := BootsSymDecrypt(carry[1], priv)

		tt := mess3
		if mess1 != 0 {
			tt = mess2
		}

		if messmux != tt {
			fmt.Printf("\tError[fullAdderMUX]: %d - %f - %f - %f - %f\n", i,
				TorusToDouble(LwePhase(temp[0], priv.LweKey)),
				TorusToDouble(LwePhase(carry[0], priv.LweKey)),
				TorusToDouble(LwePhase(temp[1], priv.LweKey)),
				TorusToDouble(LwePhase(carry[1], priv.LweKey)),
			)
		}

		carry[0] = Copy(carry[1], key)
	}
	sum[nbBits] = Copy(carry[1], key)
	return sum
}

func fullAdder(x []*LweSample, y []*LweSample, nbBits int, key *PublicKey, priv *PrivateKey) []*LweSample {
	inOutParams := priv.Params.InOutParams
	sum := NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	BootsSymEncrypt(carry[0], 0, priv) // first carry initialized to 0
	// temps
	temp := NewLweSampleArray(3, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = Xor(x[i], y[i], key) // temp = xi XOR yi
		sum[i] = Xor(temp[0], carry[0], key)

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		temp[1] = And(x[i], y[i], key)        // temp1 = xi AND yi
		temp[2] = And(carry[0], temp[0], key) // temp2 = carry AND temp
		carry[1] = Xor(temp[1], temp[2], key)
		carry[0] = Copy(carry[1], key)
	}
	sum[nbBits] = Copy(carry[0], key)
	return sum
}

func comparisonMUX(x []*LweSample, y []*LweSample, nbBits int, key *PublicKey, priv *PrivateKey) *LweSample {

	inOutParams := priv.Params.InOutParams
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	BootsSymEncrypt(carry[0], 1, priv) // first carry initialized to 1

	for i := 0; i < nbBits; i++ {
		temp := Xor(x[i], y[i], key) // temp = xi XOR yi
		carry[1] = Mux(temp, y[i], carry[0], key)
		carry[0] = Copy(carry[1], key)
	}
	return Copy(carry[0], key)
}

func fromBool(x bool) int32 {
	if !x {
		return 0
	} else {
		return 1
	}
}

func toBool(x int32) bool {
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

func decryptAndDisplayResult(sum []*LweSample, keyset *PrivateKey) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := BootsSymDecrypt(sum[i], keyset)
		fmt.Printf("%d ", messSum)
	}
	fmt.Println("]")
}

func keys(params *GateBootstrappingParameterSet) (*PublicKey, *PrivateKey) {
	var pubKey *PublicKey
	var privKey *PrivateKey
	if _, err := os.Stat("private.key"); err == nil {
		fmt.Println("------ Reading keys from file ------")
		privKey, _ = ReadPrivKey("private.key")
		pubKey, _ = ReadPubKey("public.key")

	} else {
		fmt.Println("------ Key Generation ------")
		// generate the keys
		pubKey, privKey = GenerateKeys(params)
		WritePrivKey(privKey, "private.key")
		WritePubKey(pubKey, "public.key")
	}
	return pubKey, privKey
}

func main() {
	const (
		nbBits   = 8
		nbTrials = 1
	)
	// generate params
	var minimumLambda int32 = 100
	params := NewDefaultGateBootstrappingParameters(minimumLambda)
	inOutParams := params.InOutParams
	pubKey, privKey := keys(params)

	for trial := 0; trial < nbTrials; trial++ {

		xBits := toBits(22)
		yBits := toBits(33)

		// generate samples
		x := NewLweSampleArray(nbBits, inOutParams)
		y := NewLweSampleArray(nbBits, inOutParams)
		for i := 0; i < nbBits; i++ {
			//BootsSymEncrypt(x[i], rand.Int31()%2, keyset)
			//BootsSymEncrypt(y[i], rand.Int31()%2, keyset)
			BootsSymEncrypt(x[i], int32(xBits[i]), privKey)
			BootsSymEncrypt(y[i], int32(yBits[i]), privKey)
		}
		// output sum
		//sum := NewLweSampleArray(nbBits+1, inOutParams)

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
		var messCarry int32 = 0
		for i := 0; i < nbBits; i++ {
			messX := BootsSymDecrypt(x[i], privKey)
			messY := BootsSymDecrypt(y[i], privKey)
			messSum := BootsSymDecrypt(sum[i], privKey)

			if messSum != (messX ^ messY ^ messCarry) {
				fmt.Printf("\tVerification Error %d, %f - %f - %f\n", i,
					TorusToDouble(LwePhase(x[i], privKey.LweKey)),
					TorusToDouble(LwePhase(y[i], privKey.LweKey)),
					TorusToDouble(LwePhase(sum[i], privKey.LweKey)),
				)
			}
			if messCarry != 0 {
				messCarry = fromBool(toBool(messX) || toBool(messY))
			} else {
				messCarry = fromBool(toBool(messX) && toBool(messY))
			}
		}
		messSum := BootsSymDecrypt(sum[nbBits], privKey)
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
			var messCarry int32 = 0
			for i := 0; i < nbBits; i++ {
				messX := BootsSymDecrypt(x[i], privKey)
				messY := BootsSymDecrypt(y[i], privKey)
				messSum := BootsSymDecrypt(sum[i], privKey)

				if messSum != (messX ^ messY ^ messCarry) {
					fmt.Printf("\tVerification Error - %d, %d\n", trial, i)
				}

				if messCarry != 0 {
					messCarry = fromBool(toBool(messX) || toBool(messY))
				} else {
					messCarry = fromBool(toBool(messX) && toBool(messY))
				}
			}
			messSum := BootsSymDecrypt(sum[nbBits], privKey)
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
			var messCarry int32 = 1
			for i := 0; i < nbBits; i++ {
				messX := BootsSymDecrypt(x[i], privKey)
				messY := BootsSymDecrypt(y[i], privKey)
				if messX^messY != 0 {
					messCarry = messY
				}
			}
			messComp := BootsSymDecrypt(comp, privKey)
			if messComp != messCarry {
				fmt.Printf("\tVerification Error %d, %d\n", trial, nbBits)
			}
		}
	}
}
