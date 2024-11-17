package gates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thedonutfactory/go-tfhe/core"
)

const NB_BITS = 8

// generate params
var (
	minimumLambda int32 = 100
	// generate the keys
	pubKey, priKey            = DefaultGateBootstrappingParameters(minimumLambda).GenerateKeys()
	ops            Operations = &CipheredOperations{Pk: pubKey}
)

func to4Bits(val int) []int {
	l := make([]int, NB_BITS)

	l[0] = val & 0x1
	l[1] = (val & 0x2) >> 1
	l[2] = (val & 0x4) >> 2
	l[3] = (val & 0x8) >> 3
	return l
}

func toCiphertext(nums ...int) [][]*core.LweSample {
	r := make([][]*core.LweSample, 0)
	nbBits := int32(NB_BITS)
	for _, num := range nums {
		xBits := to4Bits(num)
		x := core.NewLweSampleArray(nbBits, pubKey.Params.InOutParams)
		for i := int32(0); i < nbBits; i++ {
			x[i] = priKey.BootsSymEncrypt(xBits[i])
		}
		r = append(r, x)
	}
	return r
}

func toPlaintext(ciphers ...[]*core.LweSample) []int {
	arr := make([]int, len(ciphers))
	for ci, c := range ciphers {
		var current int = 0
		for i := 0; i < len(c); i++ {
			message := priKey.BootsSymDecrypt(c[i])
			fmt.Printf("%d ", message)
			if message == 1 {
				current |= (1 << i)
			}
		}
		fmt.Printf("\ncurrent = %d \n", current)
		arr[ci] = current
	}
	return arr
}

func decryptAndDisplayResult(sum []*core.LweSample) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := priKey.BootsSymDecrypt(sum[i])
		fmt.Printf("%d ", messSum)
	}
	fmt.Print("]\n")
}

func TestToPlaintext(tt *testing.T) {
	assert := assert.New(tt)
	value := 7
	v := toCiphertext(value)
	v1 := toPlaintext(v[0])
	assert.EqualValues(value, v1[0])
}

func TestCompareBit(tt *testing.T) {
	assert := assert.New(tt)
	b := priKey.BootsSymEncrypt(1)

	carry := core.NewLweSample(pubKey.Params.InOutParams)
	carryBit := priKey.BootsSymDecrypt(carry)
	bBit := priKey.BootsSymDecrypt(b)

	assert.EqualValues(0, carryBit)
	assert.EqualValues(1, bBit)
}

func TestEqual(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 3)
	result := ops.Equals(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)
	equalityBit := priKey.BootsSymDecrypt(result[0])
	assert.EqualValues(1, equalityBit)
}

func TestNotEqual(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 3)
	result := ops.Equals(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)
	equalityBit := priKey.BootsSymDecrypt(result[0])
	assert.EqualValues(0, equalityBit)
}

func TestNotEquals(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 3)
	result := ops.Equals(v[0], v[1], NB_BITS)
	result = ops.Not(result, NB_BITS)
	decryptAndDisplayResult(result)
	equalityBit := priKey.BootsSymDecrypt(result[0])
	assert.EqualValues(1, equalityBit)
}

func TestMinimum(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 4)
	result := ops.Minimum(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 3 -> 0011
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestMaximum(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4, 1)
	result := ops.Maximum(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 4 -> 0100
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestAddition(tt *testing.T) {
	assert := assert.New(tt)

	//ctx := gates.Default128bitGateBootstrappingParameters()
	//pub, prv := keys(ctx) //ctx.GenerateKeys()

	// encrypt 2 8-bit ciphertexts
	x := priKey.Encrypt(int8(2))
	y := priKey.Encrypt(int8(1))

	//v := toCiphertext(2, 1)
	result := ops.Add(x, y, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 3 -> 0011
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestSubtraction(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(5, 2)
	result := ops.Sub(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 3 -> 0011
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestSubtractionToZero(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1, 1)
	result := ops.Sub(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0 -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestSubtractionFromZero(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1, 0)
	result := ops.Sub(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 1 -> 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestGtAndSubtractionFrom(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 5, 1)
	r1 := ops.Gt(v[1], v[0], NB_BITS)
	decryptAndDisplayResult(r1)
	decryptAndDisplayResult(v[2])

	negR1 := ops.Negate(r1, NB_BITS)
	decryptAndDisplayResult(negR1)
	toPlaintext(negR1)

	result := ops.Add(v[2], negR1, NB_BITS)
	decryptAndDisplayResult(result)
	toPlaintext(result)

	// result should be 0 -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestMultiply(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 2)
	result := ops.Mul(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 6 -> 0110
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(6, r[0])
}

func TestGte(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 1)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be true(1) -> 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestGteCheckFalse(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1, 5)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be false(0) -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestGteCheckEquality(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 2)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be true(1) -> 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestGt(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 1)
	result := ops.Gt(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be true(1) -> 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestGtFalse(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 5)
	result := ops.Gt(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be false(1) -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestGtFalseWhenEqual(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 3)
	result := ops.Gt(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be false(1) -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestShiftLeft(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.ShiftLeft(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0001 << 1 = 0010
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestShiftLeftByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.ShiftLeft(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0001 << 2 = 0100
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestShiftRight(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.ShiftRight(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0100 >> 1 = 0010
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestShiftRightByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.ShiftRight(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0100 >> 2 = 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestUshiftLeft(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.UshiftLeft(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0001 << 1 = 0010
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestUshiftLeftByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.UshiftLeft(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0001 << 2 = 0100
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestUshiftRight(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.UshiftRight(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0100 >> 1 = 0010
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestUshiftRightByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.UshiftRight(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0100 >> 2 = 0001
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))
}

func TestDivide(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(6, 3)
	result := ops.Div(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 2 -> 0010
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))
}

func TestPow(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2)
	result := ops.Pow(v[0], 3, NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 8 -> 1000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(8, r[0])
}

func TestNegate(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2)
	result := ops.Negate(v[0], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 14 -> 1110
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(1, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(14, r[0])
}

func TestNegateZero(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(0)
	result := ops.Negate(v[0], NB_BITS)
	decryptAndDisplayResult(result)

	// result should be 0 -> 0000
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[0]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[1]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[2]))
	assert.EqualValues(0, priKey.BootsSymDecrypt(result[3]))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestAnd(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.And(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a&b, r[0])
}

func TestLogicalAnd(tt *testing.T) {
	assert := assert.New(tt)
	a, b := true, true
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.And(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a && b, toBool(r[0]))
}

func TestLogicalAnd2(tt *testing.T) {
	assert := assert.New(tt)
	a, b := true, false
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.And(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a && b, toBool(r[0]))
}

func TestLogicalAnd3(tt *testing.T) {
	assert := assert.New(tt)
	a, b := false, false
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.And(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a && b, toBool(r[0]))
}

func TestLogicalOr(tt *testing.T) {
	assert := assert.New(tt)
	a, b := true, true
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.Or(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a || b, toBool(r[0]))
}

func TestLogicalOr2(tt *testing.T) {
	assert := assert.New(tt)
	a, b := true, false
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.Or(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a || b, toBool(r[0]))
}

func TestLogicalOr3(tt *testing.T) {
	assert := assert.New(tt)
	a, b := false, false
	ca, cb := ops.Constant(a, NB_BITS), ops.Constant(b, NB_BITS)
	result := ops.Or(ca, cb, NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a || b, toBool(r[0]))
}

func toBool(val int) bool {
	if val == 0 {
		return false
	} else {
		return true
	}
}

func TestOr(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.Or(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a|b, r[0])
}

func TestXor(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 7, 3
	v := toCiphertext(a, b)
	result := ops.Xor(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a^b, r[0])
}

func TestNot(tt *testing.T) {
	assert := assert.New(tt)
	var a int = 4
	v := toCiphertext(a)
	result := ops.Not(v[0], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(0b1111^a, r[0])
}

func TestXnor(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.Xnor(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(0b1111^(a^b), r[0])
}

func TestNand(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.Nand(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(0b1111^(a&b), r[0])
}

func TestAndNY(tt *testing.T) {
	// Gate: not(a) and b
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.AndNY(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(^a&b, r[0])
}

func TestAndYN(tt *testing.T) {
	// Gate: a and not(b)
	assert := assert.New(tt)
	a, b := 2, 8
	v := toCiphertext(a, b)
	result := ops.AndYN(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a&^b, r[0])
}

func TestOrNY(tt *testing.T) {
	// Gate: not(a) or b
	assert := assert.New(tt)
	a, b := 2, 8
	v := toCiphertext(a, b)
	result := ops.OrNY(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(16+(^a|b), r[0])
}

func TestOrYN(tt *testing.T) {
	// Gate: a or not(b)
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.OrYN(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(16+(a|^b), r[0])
}
