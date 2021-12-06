package tfhe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const NB_BITS = 4

// generate params
var (
	minimumLambda int = 100
	// generate the secret keyset
	keyset            = NewRandomGateBootstrappingSecretKeyset(NewDefaultGateBootstrappingParameters(minimumLambda))
	ops    Operations = &CipheredOperations{bk: keyset.Cloud}
)

func to4Bits(val int) []int {
	l := make([]int, NB_BITS)

	l[0] = val & 0x1
	l[1] = (val & 0x2) >> 1
	l[2] = (val & 0x4) >> 2
	l[3] = (val & 0x8) >> 3
	return l
}

func toCiphertext(nums ...int) [][]*LweSample {
	r := make([][]*LweSample, 0)
	nbBits := NB_BITS
	for _, num := range nums {
		xBits := to4Bits(num)
		x := NewLweSampleArray(nbBits, keyset.Params.InOutParams)
		for i := 0; i < nbBits; i++ {
			BootsSymEncrypt(x[i], int64(xBits[i]), keyset)
		}
		r = append(r, x)
	}
	return r
}

// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

func toPlaintext(ciphers ...[]*LweSample) []int {
	arr := make([]int, len(ciphers))
	for ci, c := range ciphers {
		var current int = 0
		for i := 0; i < len(c); i++ {
			message := BootsSymDecrypt(c[i], keyset)
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

func decryptAndDisplayResult(sum []*LweSample, tt *testing.T) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := BootsSymDecrypt(sum[i], keyset)
		fmt.Printf("%d ", messSum)
	}
	fmt.Print("]")
}

func TestToPlaintext(tt *testing.T) {
	assert := assert.New(tt)
	value := 7
	v := toCiphertext(value)
	v1 := toPlaintext(v[0])
	assert.EqualValues(value, v1[0])
}

func TestCompareBit(tt *testing.T) {
	//assert := assert.New(tt)

	a := NewLweSample(keyset.Params.InOutParams)
	BootsSymEncrypt(a, 1, keyset)

	b := NewLweSample(keyset.Params.InOutParams)
	BootsSymEncrypt(a, 1, keyset)

	carry := NewLweSample(keyset.Params.InOutParams)
	tmp := NewLweSample(keyset.Params.InOutParams)

	result := ops.CompareBit(a, b, carry, tmp)

	carryBit := BootsSymDecrypt(carry, keyset)
	bBit := BootsSymDecrypt(b, keyset)

	fmt.Println(carryBit)
	fmt.Println(bBit)
	fmt.Println(result)
}

func TestEqual(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 3)
	result := ops.Equals(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)
	equalityBit := BootsSymDecrypt(result[0], keyset)
	assert.EqualValues(1, equalityBit)
}

func TestNotEqual(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 3)
	result := ops.Equals(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)
	equalityBit := BootsSymDecrypt(result[0], keyset)
	assert.EqualValues(0, equalityBit)
}

func TestMinimum(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 4)
	result := ops.Minimum(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 3 -> 0011
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestMaximum(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4, 1)
	result := ops.Maximum(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 4 -> 0100
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestAddition(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 1)
	result := ops.Add(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 3 -> 0011
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestSubtraction(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(5, 2)
	result := ops.Sub(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 3 -> 0011
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(3, r[0])
}

func TestMultiply(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 2)
	result := ops.Mul(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 6 -> 0110
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(6, r[0])
}

func TestGte(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(3, 1)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be true(1) -> 0001
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestGteCheckFalse(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1, 5)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be false(0) -> 0000
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(0, r[0])
}

func TestGteCheckEquality(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2, 2)
	result := ops.Gte(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be true(1) -> 0001
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestShiftLeft(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.ShiftLeft(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0001 << 1 = 0010
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestShiftLeftByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.ShiftLeft(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0001 << 2 = 0100
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestShiftRight(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.ShiftRight(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0100 >> 1 = 0010
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestShiftRightByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.ShiftRight(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0100 >> 2 = 0001
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(1, r[0])
}

func TestUshiftLeft(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.UshiftLeft(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0001 << 1 = 0010
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestUshiftLeftByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(1)
	result := ops.UshiftLeft(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0001 << 2 = 0100
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(4, r[0])
}

func TestUshiftRight(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.UshiftRight(v[0], 1, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0100 >> 1 = 0010
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(2, r[0])
}

func TestUshiftRightByTwo(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(4)
	result := ops.UshiftRight(v[0], 2, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 0100 >> 2 = 0001
	assert.EqualValues(1, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))
}

func TestDivide(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(6, 3)
	result := ops.Div(v[0], v[1], NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 2 -> 0010
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[3], keyset))
}

func TestPow(tt *testing.T) {
	assert := assert.New(tt)
	v := toCiphertext(2)
	result := ops.Pow(v[0], 3, NB_BITS)
	decryptAndDisplayResult(result, tt)

	// result should be 8 -> 1000
	assert.EqualValues(0, BootsSymDecrypt(result[0], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[1], keyset))
	assert.EqualValues(0, BootsSymDecrypt(result[2], keyset))
	assert.EqualValues(1, BootsSymDecrypt(result[3], keyset))

	r := toPlaintext(result)
	assert.EqualValues(8, r[0])
}

func TestAnd(tt *testing.T) {
	assert := assert.New(tt)
	a, b := 2, 3
	v := toCiphertext(a, b)
	result := ops.And(v[0], v[1], NB_BITS)
	r := toPlaintext(result)
	assert.EqualValues(a&b, r[0])
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
