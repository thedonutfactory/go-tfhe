package main

import (
	"fmt"
	"time"

	tfhe "github.com/TheDonutFactory/go-tfhe"
)

func fullAdderMUX(sum []*tfhe.LweSample, x []*tfhe.LweSample, y []*tfhe.LweSample, nbBits int, keyset *tfhe.TFheGateBootstrappingSecretKeySet) {
	inOutParams := keyset.Params.InOutParams
	// carries
	carry := tfhe.NewLweSampleArray(2, inOutParams)
	tfhe.BootsSymEncrypt(carry[0], 0, keyset) // first carry initialized to 0
	// temps
	temp := tfhe.NewLweSampleArray(2, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		tfhe.BootsXOR(temp[0], x[i], y[i], keyset.Cloud) // temp = xi XOR yi
		tfhe.BootsXOR(sum[i], temp[0], carry[0], keyset.Cloud)

		// carry = MUX(xi XOR yi, carry(i-1), xi AND yi)
		tfhe.BootsAND(temp[1], x[i], y[i], keyset.Cloud) // temp1 = xi AND yi
		tfhe.BootsMUX(carry[1], temp[0], carry[0], temp[1], keyset.Cloud)

		mess1 := tfhe.BootsSymDecrypt(temp[0], keyset)
		mess2 := tfhe.BootsSymDecrypt(carry[0], keyset)
		mess3 := tfhe.BootsSymDecrypt(temp[1], keyset)
		messmux := tfhe.BootsSymDecrypt(carry[1], keyset)

		t := mess3
		if mess1 != 0 {
			t = mess2
		}

		if messmux != t {
			fmt.Printf("\tError[fullAdderMUX]: %d - %f - %f - %f - %f\n", i,
				tfhe.T32tod(tfhe.LwePhase(temp[0], keyset.LweKey)),
				tfhe.T32tod(tfhe.LwePhase(carry[0], keyset.LweKey)),
				tfhe.T32tod(tfhe.LwePhase(temp[1], keyset.LweKey)),
				tfhe.T32tod(tfhe.LwePhase(carry[1], keyset.LweKey)),
			)
		}

		tfhe.BootsCOPY(carry[0], carry[1], keyset.Cloud)
	}
	tfhe.BootsCOPY(sum[nbBits], carry[1], keyset.Cloud)
}

func fullAdder(sum []*tfhe.LweSample, x []*tfhe.LweSample, y []*tfhe.LweSample, nbBits int, keyset *tfhe.TFheGateBootstrappingSecretKeySet) {
	inOutParams := keyset.Params.InOutParams
	// carries
	carry := tfhe.NewLweSampleArray(2, inOutParams)
	tfhe.BootsSymEncrypt(carry[0], 0, keyset) // first carry initialized to 0
	// temps
	temp := tfhe.NewLweSampleArray(3, inOutParams)

	for i := 0; i < nbBits; i++ {
		//sumi = xi XOR yi XOR carry(i-1)
		tfhe.BootsXOR(temp[0], x[i], y[i], keyset.Cloud) // temp = xi XOR yi
		tfhe.BootsXOR(sum[i], temp[0], carry[0], keyset.Cloud)

		// carry = (xi AND yi) XOR (carry(i-1) AND (xi XOR yi))
		tfhe.BootsAND(temp[1], x[i], y[i], keyset.Cloud)        // temp1 = xi AND yi
		tfhe.BootsAND(temp[2], carry[0], temp[0], keyset.Cloud) // temp2 = carry AND temp
		tfhe.BootsXOR(carry[1], temp[1], temp[2], keyset.Cloud)
		tfhe.BootsCOPY(carry[0], carry[1], keyset.Cloud)
	}
	tfhe.BootsCOPY(sum[nbBits], carry[0], keyset.Cloud)
}

func comparisonMUX(comp *tfhe.LweSample, x []*tfhe.LweSample, y []*tfhe.LweSample, nbBits int, keyset *tfhe.TFheGateBootstrappingSecretKeySet) {

	inOutParams := keyset.Params.InOutParams
	// carries
	carry := tfhe.NewLweSampleArray(2, inOutParams)
	tfhe.BootsSymEncrypt(carry[0], 1, keyset) // first carry initialized to 1
	// temps
	temp := tfhe.NewLweSample(inOutParams)

	for i := 0; i < nbBits; i++ {
		tfhe.BootsXOR(temp, x[i], y[i], keyset.Cloud) // temp = xi XOR yi
		tfhe.BootsMUX(carry[1], temp, y[i], carry[0], keyset.Cloud)
		tfhe.BootsCOPY(carry[0], carry[1], keyset.Cloud)
	}
	tfhe.BootsCOPY(comp, carry[0], keyset.Cloud)
}

func fromBool(x bool) int32 {
	if x == false {
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

func decryptAndDisplayResult(sum []*tfhe.LweSample, keyset *tfhe.TFheGateBootstrappingSecretKeySet) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := tfhe.BootsSymDecrypt(sum[i], keyset)
		fmt.Printf("%d ", messSum)
	}
	fmt.Println("]")
}

func main() {
	const (
		nbBits   = 8
		nbTrials = 1
	)
	// generate params
	var minimumLambda int32 = 100
	params := tfhe.NewDefaultGateBootstrappingParameters(minimumLambda)
	inOutParams := params.InOutParams
	// generate the secret keyset
	keyset := tfhe.NewRandomGateBootstrappingSecretKeyset(params)

	for trial := 0; trial < nbTrials; trial++ {

		xBits := toBits(22)
		yBits := toBits(33)

		// generate samples
		x := tfhe.NewLweSampleArray(nbBits, inOutParams)
		y := tfhe.NewLweSampleArray(nbBits, inOutParams)
		for i := 0; i < nbBits; i++ {
			//tfhe.BootsSymEncrypt(x[i], rand.Int31()%2, keyset)
			//tfhe.BootsSymEncrypt(y[i], rand.Int31()%2, keyset)
			tfhe.BootsSymEncrypt(x[i], int32(xBits[i]), keyset)
			tfhe.BootsSymEncrypt(y[i], int32(yBits[i]), keyset)
		}
		// output sum
		sum := tfhe.NewLweSampleArray(nbBits+1, inOutParams)

		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits addition circuit (FA in MUX version), trial %d\n", nbBits, trial)
		start := time.Now()
		fullAdderMUX(sum, x, y, nbBits, keyset)
		duration := time.Since(start)

		decryptAndDisplayResult(sum, keyset)
		// Formatted string, such as "2h3m0.5s" or "4.503μs"
		fmt.Printf("finished Bootstrapping %d bits addition circuit (FA in MUX version)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		var messCarry int32 = 0
		for i := 0; i < nbBits; i++ {
			messX := tfhe.BootsSymDecrypt(x[i], keyset)
			messY := tfhe.BootsSymDecrypt(y[i], keyset)
			messSum := tfhe.BootsSymDecrypt(sum[i], keyset)

			if messSum != (messX ^ messY ^ messCarry) {
				fmt.Printf("\tVerification Error %d, %f - %f - %f\n", i,
					tfhe.T32tod(tfhe.LwePhase(x[i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(y[i], keyset.LweKey)),
					tfhe.T32tod(tfhe.LwePhase(sum[i], keyset.LweKey)),
				)
			}

			if messCarry != 0 {
				messCarry = fromBool(toBool(messX) || toBool(messY))
			} else {
				messCarry = fromBool(toBool(messX) && toBool(messY))
			}

			//messCarry = messCarry ? (messX || messY) : (messX && messY);
		}
		messSum := tfhe.BootsSymDecrypt(sum[nbBits], keyset)
		if messSum != messCarry {
			fmt.Printf("\tVerification Error - %d, %d bits\n", trial, nbBits)
		}

		// evaluate the addition circuit
		fmt.Printf("Starting Bootstrapping %d bits addition circuit (FA)...trial %d\n", nbBits, trial)
		start = time.Now()
		fullAdder(sum, x, y, nbBits, keyset)
		duration = time.Since(start)
		decryptAndDisplayResult(sum, keyset)
		fmt.Printf("finished Bootstrappings %d bits addition circuit (FA)\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int32 = 0
			for i := 0; i < nbBits; i++ {
				messX := tfhe.BootsSymDecrypt(x[i], keyset)
				messY := tfhe.BootsSymDecrypt(y[i], keyset)
				messSum := tfhe.BootsSymDecrypt(sum[i], keyset)

				if messSum != (messX ^ messY ^ messCarry) {
					fmt.Printf("\tVerification Error - %d, %d\n", trial, i)
				}

				if messCarry != 0 {
					messCarry = fromBool(toBool(messX) || toBool(messY))
				} else {
					messCarry = fromBool(toBool(messX) && toBool(messY))
				}
				//messCarry = messCarry ? (messX || messY) : (messX && messY);
			}
			messSum := tfhe.BootsSymDecrypt(sum[nbBits], keyset)
			if messSum != messCarry {
				fmt.Printf("\tVerification Error - %d, %d\n", trial, nbBits)
			}
		}

		comp := tfhe.NewLweSample(inOutParams)
		// evaluate the addition circuit
		fmt.Printf("starting Bootstrapping %d bits comparison, trial %d\n", nbBits, trial)
		start = time.Now()
		comparisonMUX(comp, x, y, nbBits, keyset)
		duration = time.Since(start)
		fmt.Printf("finished Bootstrappings %d bits comparison\n", nbBits)
		fmt.Printf("total time: %s\n", duration)

		// verification
		{
			var messCarry int32 = 1
			for i := 0; i < nbBits; i++ {
				messX := tfhe.BootsSymDecrypt(x[i], keyset)
				messY := tfhe.BootsSymDecrypt(y[i], keyset)

				if messX^messY != 0 {
					messCarry = messY
				}

				//messCarry = (messX ^ messY) ? messY : messCarry;
			}
			messComp := tfhe.BootsSymDecrypt(comp, keyset)
			if messComp != messCarry {
				fmt.Printf("\tVerification Error %d, %d\n", trial, nbBits)
			}
		}
	}
}
