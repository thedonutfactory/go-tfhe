package main

import (
	"fmt"
	"math"
	"os"
	"time"

	. "github.com/thedonutfactory/go-tfhe/core"
	. "github.com/thedonutfactory/go-tfhe/gates"
	. "github.com/thedonutfactory/go-tfhe/io"
	. "github.com/thedonutfactory/go-tfhe/types"
)

func fullAdderMUX(x, y Int, key *PublicKey, priv *PrivateKey) Int {
	nbBits := len(x)
	inOutParams := priv.Params.InOutParams
	sum := NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	carry[0] = priv.BootsSymEncrypt(0) // first carry initialized to 0
	// temps
	temp := NewLweSampleArray(2, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = key.Xor(x[i], y[i]) // temp = xi XOR yi
		sum[i] = key.Xor(temp[0], carry[0])

		// carry = MUX(xi XOR yi, carry(i-1), xi AND yi)
		temp[1] = key.And(x[i], y[i]) // temp1 = xi AND yi
		carry[1] = key.Mux(temp[0], carry[0], temp[1])

		mess1 := priv.BootsSymDecrypt(temp[0])
		mess2 := priv.BootsSymDecrypt(carry[0])
		mess3 := priv.BootsSymDecrypt(temp[1])
		messmux := priv.BootsSymDecrypt(carry[1])

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

		carry[0] = key.Copy(carry[1])
	}
	sum[nbBits] = key.Copy(carry[1])
	return sum
}

func fullAdder(x, y Int, key *PublicKey, priv *PrivateKey) Int {
	nbBits := len(x)
	inOutParams := priv.Params.InOutParams
	sum := NewLweSampleArray(int32(nbBits)+1, inOutParams)
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	carry[0] = priv.BootsSymEncrypt(0) // first carry initialized to 0
	// temps
	temp := NewLweSampleArray(3, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		temp[0] = key.Xor(x[i], y[i]) // temp = xi XOR yi
		sum[i] = key.Xor(temp[0], carry[0])

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		temp[1] = key.And(x[i], y[i])        // temp1 = xi AND yi
		temp[2] = key.And(carry[0], temp[0]) // temp2 = carry AND temp
		carry[1] = key.Xor(temp[1], temp[2])
		carry[0] = key.Copy(carry[1])
	}
	sum[nbBits] = key.Copy(carry[0])
	return sum
}

func comparisonMUX(x, y Int, key *PublicKey, priv *PrivateKey) Int1 {
	nbBits := len(x)
	inOutParams := priv.Params.InOutParams
	// carries
	carry := NewLweSampleArray(2, inOutParams)
	carry[0] = priv.BootsSymEncrypt(1) // first carry initialized to 1

	for i := 0; i < nbBits; i++ {
		temp := key.Xor(x[i], y[i]) // temp = xi XOR yi
		carry[1] = key.Mux(temp, y[i], carry[0])
		carry[0] = key.Copy(carry[1])
	}
	out := NewInt1(inOutParams)
	out[0] = key.Copy(carry[0])
	return out
}

func fromBool(x bool) int {
	if !x {
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

/*
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
*/

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func toBits(val, size int) []int {
	l := make([]int, size)

	l[0] = val & 0x1
	for i := 1; i < size; i++ {
		l[i] = (val & powInt(2, i)) >> i
	}
	return l
}

func decryptAndDisplayResult(sum []*LweSample, keyset *PrivateKey) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := keyset.BootsSymDecrypt(sum[i])
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
		pubKey, privKey = params.GenerateKeys()
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
	ctx := DefaultGateBootstrappingParameters(minimumLambda)
	pubKey, privKey := keys(ctx)

	for trial := 0; trial < nbTrials; trial++ {

		x := privKey.Encrypt(int8(22))
		y := privKey.Encrypt(int8(33))

		// output sum
		//sum := NewLweSampleArray(nbBits+1, inOutParams)

		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits addition circuit (FA in MUX version), trial %d\n", nbBits, trial)
		start := time.Now()
		sum := fullAdderMUX(x, y, pubKey, privKey)
		duration := time.Since(start)

		decryptAndDisplayResult(sum, privKey)
		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Printf("finished Bootstrapping %d bits addition circuit (FA in MUX version)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		var messCarry int = 0
		for i := 0; i < nbBits; i++ {
			messX := privKey.BootsSymDecrypt(x[i])
			messY := privKey.BootsSymDecrypt(y[i])
			messSum := privKey.BootsSymDecrypt(sum[i])

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
		messSum := privKey.BootsSymDecrypt(sum[nbBits])
		if messSum != messCarry {
			fmt.Printf("\tVerification Error - %d, %d bits\n", trial, nbBits)
		}

		// evaluate the addition circuit
		fmt.Printf("Starting Bootstrapping %d bits addition circuit (FA)...trial %d\n", nbBits, trial)
		start = time.Now()
		sum = fullAdder(x, y, pubKey, privKey)
		duration = time.Since(start)
		decryptAndDisplayResult(sum, privKey)
		fmt.Printf("finished Bootstrappings %d bits addition circuit (FA)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int = 0
			for i := 0; i < nbBits; i++ {
				messX := privKey.BootsSymDecrypt(x[i])
				messY := privKey.BootsSymDecrypt(y[i])
				messSum := privKey.BootsSymDecrypt(sum[i])

				if messSum != (messX ^ messY ^ messCarry) {
					fmt.Printf("\tVerification Error - %d, %d\n", trial, i)
				}

				if messCarry != 0 {
					messCarry = fromBool(toBool(messX) || toBool(messY))
				} else {
					messCarry = fromBool(toBool(messX) && toBool(messY))
				}
			}
			messSum := privKey.BootsSymDecrypt(sum[nbBits])
			if messSum != messCarry {
				fmt.Printf("\tVerification Error - %d, %d\n", trial, nbBits)
			}
		}

		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits comparison, trial %d\n", nbBits, trial)
		start = time.Now()
		comp := comparisonMUX(x, y, pubKey, privKey)
		duration = time.Since(start)
		fmt.Printf("finished Bootstrappings %d bits comparison\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int = 1
			for i := 0; i < nbBits; i++ {
				messX := privKey.BootsSymDecrypt(x[i])
				messY := privKey.BootsSymDecrypt(y[i])
				if messX^messY != 0 {
					messCarry = messY
				}
			}
			messComp := privKey.BootsSymDecrypt(comp[0])
			if messComp != messCarry {
				fmt.Printf("\tVerification Error %d, %d\n", trial, nbBits)
			}
		}
	}
}
