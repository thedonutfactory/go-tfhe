package tfhe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const NB_BITS = 4

// generate params
var (
	minimumLambda int32 = 100
	// generate the secret keyset
	keyset             = NewRandomGateBootstrappingSecretKeyset(NewDefaultGateBootstrappingParameters(minimumLambda))
	ops    *Operations = &Operations{bk: keyset.Cloud}
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
	nbBits := int32(NB_BITS)
	for _, num := range nums {
		xBits := to4Bits(num)
		x := NewLweSampleArray(nbBits, keyset.Params.InOutParams)
		for i := int32(0); i < nbBits; i++ {
			BootsSymEncrypt(x[i], int32(xBits[i]), keyset)
		}
		r = append(r, x)
	}
	return r
}

func decryptAndDisplayResult(sum []*LweSample, tt *testing.T) {
	fmt.Print("[ ")
	for i := len(sum) - 1; i >= 0; i-- {
		messSum := BootsSymDecrypt(sum[i], keyset)
		fmt.Printf("%d ", messSum)
	}
	fmt.Print("]")
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
}
