package gates

import (
	"fmt"
	"math"
	"strconv"

	. "github.com/thedonutfactory/go-tfhe/core"
	. "github.com/thedonutfactory/go-tfhe/types"
)

func (key *PrivateKey) Encrypt(message interface{}) Int {
	return key.ToCtxt(message)
}

func (key *PrivateKey) Decrypt(message Int) int {
	return key.ToPtxt(message)
}

/** encrypts a boolean */
func (key *PrivateKey) BootsSymEncrypt(message int) *LweSample {
	_1s8 := ModSwitchToTorus32(1, 8)
	var mu Torus32 = -_1s8
	if message != 0 {
		mu = _1s8
	}
	//Torus32 mu = message ? _1s8 : -_1s8;
	r := NewLweSample(key.LweKey.Params)
	alpha := key.Params.InOutParams.AlphaMin //TODO: specify noise
	LweSymEncrypt(r, mu, alpha, key.LweKey)
	return r
}

/** decrypts a boolean */
func (key *PrivateKey) BootsSymDecrypt(sample *LweSample) int {
	mu := LwePhase(sample, key.LweKey)
	if mu > 0 {
		return 1
	} else {
		return 0
	}
}

func (key *PrivateKey) ToPtxt(val Int) int {
	return key.PlainBits(val)
}

func (key *PrivateKey) ToCtxt(val interface{}) Int {
	var ctxt Int
	switch v := val.(type) {
	default:
		fmt.Printf("unexpected type %T", v)
	case uint8:
		value, ok := val.(uint8)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 8)
	case uint16:
		value, ok := val.(uint16)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 16)
	case uint32:
		value, ok := val.(uint32)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 32)
	case uint64:
		value, ok := val.(uint64)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 64)
	case int8:
		value, ok := val.(int8)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 8)
	case int16:
		value, ok := val.(int16)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 16)
	case int32:
		value, ok := val.(int32)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 32)
	case int64:
		value, ok := val.(int64)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		ctxt = key.CipherBits(int(value), 64)
	}
	return ctxt
}

func (key *PrivateKey) CipherBits(val, size int) Int {
	ctxt := NewInt(size, key.Params.InOutParams)
	ctxt[0] = key.BootsSymEncrypt(val & 0x1)
	for i := 1; i < size; i++ {
		ctxt[i] = key.BootsSymEncrypt(val & PowInt(2, i) >> i)
	}
	return ctxt
}

func (key *PrivateKey) PlainBits(val Int) int {
	binary := ""
	for i := 0; i < len(val); i++ {
		binary += strconv.Itoa(key.BootsSymDecrypt(val[i]))
	}
	output, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return int(output)
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func Bits(val, size int) []int {
	l := make([]int, size)

	l[0] = val & 0x1
	for i := 1; i < size; i++ {
		l[i] = (val & PowInt(2, i)) >> i
	}
	return l
}

func GetBits(val interface{}) []int {

	var arr []int

	switch v := val.(type) {
	default:
		fmt.Printf("unexpected type %T", v)
	case uint8:
		value, ok := val.(uint8)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 8)
	case uint16:
		value, ok := val.(uint16)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 16)
	case uint32:
		value, ok := val.(uint32)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 32)
	case uint64:
		value, ok := val.(uint64)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 64)
	case int8:
		value, ok := val.(int8)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 8)
	case int16:
		value, ok := val.(int16)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 16)
	case int32:
		value, ok := val.(int32)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 32)
	case int64:
		value, ok := val.(int64)
		if !ok {
			fmt.Printf("Unable to convert type %T", value)
		}
		arr = Bits(int(value), 64)
	}

	return arr
}
